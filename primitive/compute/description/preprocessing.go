//
//   Copyright © 2021 Uncharted Software Inc.
//
//   Licensed under the Apache License, Version 2.0 (the "License");
//   you may not use this file except in compliance with the License.
//   You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
//   Unless required by applicable law or agreed to in writing, software
//   distributed under the License is distributed on an "AS IS" BASIS,
//   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//   See the License for the specific language governing permissions and
//   limitations under the License.

package description

import (
	"sort"
	"strings"

	"github.com/pkg/errors"
	"github.com/uncharted-distil/distil-compute/model"
	"github.com/uncharted-distil/distil-compute/pipeline"
)

// UserDatasetDescription contains the basic parameters needs to generate
// the user dataset pipeline.
type UserDatasetDescription struct {
	AllFeatures      []*model.Variable
	TargetFeature    *model.Variable
	SelectedFeatures []string
	Filters          []*model.FilterSet
}

// UserDatasetAugmentation contains the augmentation parameters required
// for user dataset pipelines.
type UserDatasetAugmentation struct {
	SearchResult  string
	SystemID      string
	BaseDatasetID string
}

// CreatePreFeaturizedDatasetPipeline creates a pipeline that acts on a pre featurized
// dataset. The created prepend is a simplified version due to the dataset already
// having all features for the end task stored on disk.
func CreatePreFeaturizedDatasetPipeline(name string, description string, datasetDescription *UserDatasetDescription,
	augmentations []*UserDatasetAugmentation) (*pipeline.PipelineDescription, error) {
	// build a lookup for selected features
	selectedSet := map[string]bool{}
	for _, v := range datasetDescription.SelectedFeatures {
		selectedSet[strings.ToLower(v)] = true
	}

	// type all features
	featureSet := map[string]int{}
	colsToDrop := []int{}
	for _, v := range datasetDescription.AllFeatures {
		if model.IsTA2Field(v.DistilRole, v.SelectedRole) {
			variableKey := strings.ToLower(v.Key)
			featureSet[variableKey] = v.Index
			if !selectedSet[variableKey] {
				if v.Index != datasetDescription.TargetFeature.Index && !model.IsIndexRole(v.SelectedRole) {
					colsToDrop = append(colsToDrop, v.Index)
				}
			}
		}
	}

	steps := []Step{}
	offset := 0

	steps = append(steps, NewDenormalizeStep(map[string]DataRef{"inputs": &PipelineDataRef{0}}, []string{"produce"}))
	steps = append(steps, NewDistilColumnParserStep(nil, nil, []string{model.TA2IntegerType, model.TA2BooleanType, model.TA2RealType, model.TA2RealVectorType}))
	steps = append(steps, NewDatasetWrapperStep(map[string]DataRef{"inputs": &StepDataRef{offset, "produce"}}, []string{"produce"}, offset+1, ""))
	steps = append(steps, NewDataCleaningStep(nil, nil))
	steps = append(steps, NewDatasetWrapperStep(map[string]DataRef{"inputs": &StepDataRef{offset + 2, "produce"}}, []string{"produce"}, offset+3, ""))
	offset += 5

	updateSemanticTypes, err := createUpdateSemanticTypes(datasetDescription.TargetFeature.Key, datasetDescription.AllFeatures, selectedSet, offset)
	if err != nil {
		return nil, err
	}
	steps = append(steps, updateSemanticTypes...)
	offset += len(updateSemanticTypes)

	// apply filters
	filterData, err := filterBySet(datasetDescription.Filters, featureSet, offset)
	if err != nil {
		return nil, err
	}
	steps = append(steps, filterData...)
	offset += len(filterData)

	// need to drop the features from the dataset
	featureSelect := NewRemoveColumnsStep(nil, nil, colsToDrop)
	wrapperRemove := NewDatasetWrapperStep(map[string]DataRef{"inputs": &StepDataRef{offset - 1, "produce"}}, []string{"produce"}, offset, "")
	steps = append(steps, featureSelect, wrapperRemove)
	offset += 2

	// mark this is a preprocessing template
	steps = append(steps, NewInferenceStepData(map[string]DataRef{"inputs": &StepDataRef{offset - 1, "produce"}}))
	offset++

	inputs := []string{"inputs"}
	outputs := []DataRef{&StepDataRef{offset - 1, "produce"}}

	pip, err := NewPipelineBuilder(name, description, inputs, outputs, steps).Compile()
	if err != nil {
		return nil, err
	}

	return pip, nil
}

// CreateUserDatasetPipeline creates a pipeline description to capture user feature selection and
// semantic type information.
func CreateUserDatasetPipeline(name string, description string, datasetDescription *UserDatasetDescription,
	augmentations []*UserDatasetAugmentation) (*pipeline.PipelineDescription, error) {

	steps, err := generatePrependSteps(datasetDescription, augmentations)
	if err != nil {
		return nil, err
	}
	offset := len(steps)

	// mark this is a preprocessing template
	steps = append(steps, NewInferenceStepData(map[string]DataRef{"inputs": &StepDataRef{offset - 1, "produce"}}))
	offset++

	inputs := []string{"inputs"}
	outputs := []DataRef{&StepDataRef{offset - 1, "produce"}}

	pip, err := NewPipelineBuilder(name, description, inputs, outputs, steps).Compile()
	if err != nil {
		return nil, err
	}

	return pip, nil
}

func generatePrependSteps(datasetDescription *UserDatasetDescription,
	augmentations []*UserDatasetAugmentation) ([]Step, error) {
	offset := 0

	// filter out group variables
	datasetFeatures := []*model.Variable{}
	for _, v := range datasetDescription.AllFeatures {
		if v.Grouping == nil && !model.ExcludedDistilRoles[v.DistilRole] {
			datasetFeatures = append(datasetFeatures, v)
		}
	}

	// save the selected features in a set for quick lookup
	selectedSet := map[string]bool{}
	for _, v := range datasetDescription.SelectedFeatures {
		selectedSet[strings.ToLower(v)] = true
	}
	columnIndices := mapColumns(datasetFeatures, selectedSet)

	// create pipeline nodes for step we need to execute
	steps := []Step{} // add the denorm primitive

	// determine if this is a timeseries dataset
	isTimeseries := false
	groupingIndices := make([]int, 0)
	timeseriesGrouping := getTimeseriesGrouping(datasetDescription)
	targetName := datasetDescription.TargetFeature.Key
	if timeseriesGrouping != nil {
		isTimeseries = true
		groupingSet := map[string]bool{}
		for _, subID := range timeseriesGrouping.SubIDs {
			groupingSet[strings.ToLower(subID)] = true
		}
		groupingIndices = listColumns(datasetFeatures, groupingSet)
		targetName = timeseriesGrouping.YCol
	}

	// augment the dataset if needed
	// need to track the initial dataref and set the offset properly
	var dataRef DataRef
	dataRef = &PipelineDataRef{0}
	if augmentations != nil {
		for i := 0; i < len(augmentations); i++ {
			steps = append(steps, NewDatamartAugmentStep(
				map[string]DataRef{"inputs": dataRef},
				[]string{"produce"},
				augmentations[i].SearchResult,
				augmentations[i].SystemID,
			))
			dataRef = &StepDataRef{offset, "produce"}
			offset++
		}
	}

	if isTimeseries {
		steps = append(steps, NewTimeseriesFormatterStep(map[string]DataRef{"inputs": dataRef}, []string{"produce"}, "", -1))
		steps = append(steps, NewGroupingFieldComposeStep(nil, nil, groupingIndices, "-", "__grouping"))
		steps = append(steps, NewDatasetWrapperStep(map[string]DataRef{"inputs": &StepDataRef{offset, "produce"}}, []string{"produce"}, offset+1, ""))
		steps = append(steps, NewColumnParserStep(nil, nil, []string{model.TA2IntegerType, model.TA2BooleanType, model.TA2RealType}))
		steps = append(steps, NewDatasetWrapperStep(map[string]DataRef{"inputs": &StepDataRef{offset + 2, "produce"}}, []string{"produce"}, offset+3, ""))

		// add the time indicator type to the time column
		addTime := NewAddSemanticTypeStep(nil, nil, &ColumnUpdate{
			SemanticTypes: []string{model.TA2TimeType},
			Indices:       []int{columnIndices[strings.ToLower(timeseriesGrouping.XCol)]},
		})
		steps = append(steps, addTime)
		steps = append(steps, NewDatasetWrapperStep(map[string]DataRef{"inputs": &StepDataRef{offset + 4, "produce"}}, []string{"produce"}, offset+5, ""))

		// add the suggested grouping key type to the id columns
		addGroupingKey := NewAddSemanticTypeStep(nil, nil, &ColumnUpdate{
			SemanticTypes: []string{model.TA2SuggestedGroupingType},
			Indices:       groupingIndices,
		})
		steps = append(steps, addGroupingKey)
		steps = append(steps, NewDatasetWrapperStep(map[string]DataRef{"inputs": &StepDataRef{offset + 6, "produce"}}, []string{"produce"}, offset+7, ""))
		offset += 9
	} else {
		steps = append(steps, NewDenormalizeStep(map[string]DataRef{"inputs": dataRef}, []string{"produce"}))
		steps = append(steps, NewDistilColumnParserStep(nil, nil, []string{model.TA2IntegerType, model.TA2BooleanType, model.TA2RealType, model.TA2RealVectorType}))
		steps = append(steps, NewDatasetWrapperStep(map[string]DataRef{"inputs": &StepDataRef{offset, "produce"}}, []string{"produce"}, offset+1, ""))
		steps = append(steps, NewDataCleaningStep(nil, nil))
		steps = append(steps, NewDatasetWrapperStep(map[string]DataRef{"inputs": &StepDataRef{offset + 2, "produce"}}, []string{"produce"}, offset+3, ""))
		offset += 5

		multiBandImageGrouping := getMultiBandImageGrouping(datasetDescription)
		if multiBandImageGrouping != nil {
			selectedSet[multiBandImageGrouping.Key] = true
			attribs := &ColumnUpdate{
				SemanticTypes: []string{model.TA2GroupingKeyType},
				Indices:       []int{multiBandImageGrouping.Index},
			}
			multiBandImageUpdate := NewAddSemanticTypeStep(nil, nil, attribs)
			multiBandImageWrapper := NewDatasetWrapperStep(map[string]DataRef{"inputs": &StepDataRef{offset - 1, "produce"}}, []string{"produce"}, offset, "")
			steps = append(steps, multiBandImageUpdate, multiBandImageWrapper)
			offset += 2
		}
	}

	// create the semantic type update primitive
	updateSemanticTypes, err := createUpdateSemanticTypes(targetName, datasetDescription.AllFeatures, selectedSet, offset)
	if err != nil {
		return nil, err
	}
	steps = append(steps, updateSemanticTypes...)
	offset += len(updateSemanticTypes)

	// create the feature selection primitive
	removeFeatures := createRemoveFeatures(datasetFeatures, selectedSet, offset)
	steps = append(steps, removeFeatures...)
	offset += len(removeFeatures)

	// add filter primitives
	filterData, err := filterBySet(datasetDescription.Filters, columnIndices, offset)
	if err != nil {
		return nil, err
	}
	steps = append(steps, filterData...)

	// If neither have any content, we'll skip the template altogether.
	if len(updateSemanticTypes) == 0 && removeFeatures == nil &&
		len(filterData) == 0 && augmentations == nil && !isTimeseries {
		return nil, nil
	}

	return steps, nil
}

func getTimeseriesGrouping(datasetDescription *UserDatasetDescription) *model.TimeseriesGrouping {
	if model.IsTimeSeries(datasetDescription.TargetFeature.Type) {
		return datasetDescription.TargetFeature.Grouping.(*model.TimeseriesGrouping)
	}
	for _, v := range datasetDescription.AllFeatures {
		if v.Grouping != nil && model.IsTimeSeries(v.Grouping.GetType()) {
			return v.Grouping.(*model.TimeseriesGrouping)
		}
	}

	return nil
}

func getMultiBandImageGrouping(datasetDescription *UserDatasetDescription) *model.Variable {
	// multiband image type identifies multi-band image dataset
	isMultiBandImage := false
	for _, v := range datasetDescription.AllFeatures {
		if model.IsMultiBandImage(v.Type) {
			isMultiBandImage = true
			break
		}
	}

	if !isMultiBandImage {
		return nil
	}

	// grouping role will identify the grouping key to use
	for _, v := range datasetDescription.AllFeatures {
		if v.DistilRole == model.VarDistilRoleGrouping {
			return v
		}
	}

	return nil
}

func createRemoveFeatures(allFeatures []*model.Variable, selectedSet map[string]bool, offset int) []Step {
	// create a list of features to remove
	removeFeatures := []int{}
	for _, v := range allFeatures {
		if !selectedSet[strings.ToLower(v.Key)] {
			removeFeatures = append(removeFeatures, v.Index)
		}
	}

	if len(removeFeatures) == 0 {
		return nil
	}

	// instantiate the feature remove primitive
	featureSelect := NewRemoveColumnsStep(nil, nil, removeFeatures)
	wrapper := NewDatasetWrapperStep(map[string]DataRef{"inputs": &StepDataRef{offset - 1, "produce"}}, []string{"produce"}, offset, "")
	return []Step{featureSelect, wrapper}
}

type update struct {
	removeIndices []int
	addIndices    []int
}

func newUpdate() *update {
	return &update{
		addIndices:    []int{},
		removeIndices: []int{},
	}
}

func createUpdateSemanticTypes(target string, allFeatures []*model.Variable, selectedSet map[string]bool, offset int) ([]Step, error) {
	// create maps of (semantic type, index list) - primitive allows for semantic types to be added to /
	// remove from multiple columns in a single operation
	updateMap := map[string]*update{}
	attributes := make([]int, 0)
	targetIndex := -1
	for _, v := range allFeatures {
		// empty selected set means all selected
		if len(selectedSet) == 0 || selectedSet[strings.ToLower(v.Key)] {
			addType := model.MapTA2Type(v.Type)
			if addType == "" {
				return nil, errors.Errorf("variable `%s` internal type `%s` can't be mapped to ta2", v.Key, v.Type)
			}
			// unknown type must not be passed to TA2
			if addType == model.TA2UnknownType {
				addType = model.TA2StringType
			}

			removeType := model.MapTA2Type(v.OriginalType)
			if removeType == "" {
				return nil, errors.Errorf("remove variable `%s` internal type `%s` can't be mapped to ta2", v.Key, v.OriginalType)
			}

			// only apply change when types are different
			if addType != removeType {
				if _, ok := updateMap[addType]; !ok {
					updateMap[addType] = newUpdate()
				}
				updateMap[addType].addIndices = append(updateMap[addType].addIndices, v.Index)

				if _, ok := updateMap[removeType]; !ok {
					updateMap[removeType] = newUpdate()
				}
				updateMap[removeType].removeIndices = append(updateMap[removeType].removeIndices, v.Index)
			}

			// update all non target to attribute
			if !model.IsIndexRole(v.SelectedRole) {
				if !strings.EqualFold(v.Key, target) {
					attributes = append(attributes, v.Index)
				} else {
					targetIndex = v.Index
				}
			}
		}
	}

	// Copy the created maps into the column update structure used by the primitive.  Force
	// alpha ordering to make debugging / testing predictable
	keys := make([]string, 0, len(updateMap))
	for k := range updateMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	semanticTypeUpdates := []Step{}
	for _, k := range keys {
		v := updateMap[k]

		var addKey string
		if len(v.addIndices) > 0 {
			addKey = k
			add := &ColumnUpdate{
				SemanticTypes: []string{addKey},
				Indices:       v.addIndices,
			}
			addUpdate := NewAddSemanticTypeStep(nil, nil, add)
			wrapper := NewDatasetWrapperStep(map[string]DataRef{"inputs": &StepDataRef{offset - 1, "produce"}}, []string{"produce"}, offset, "")
			semanticTypeUpdates = append(semanticTypeUpdates, addUpdate, wrapper)
			offset += 2
		}

		var removeKey string
		if len(v.removeIndices) > 0 {
			removeKey = k
			remove := &ColumnUpdate{
				SemanticTypes: []string{removeKey},
				Indices:       v.removeIndices,
			}
			removeUpdate := NewRemoveSemanticTypeStep(nil, nil, remove)
			wrapper := NewDatasetWrapperStep(map[string]DataRef{"inputs": &StepDataRef{offset - 1, "produce"}}, []string{"produce"}, offset, "")
			semanticTypeUpdates = append(semanticTypeUpdates, removeUpdate, wrapper)
			offset += 2
		}
	}

	// apply attribute semantic type
	if len(attributes) > 0 {
		attribs := &ColumnUpdate{
			SemanticTypes: []string{model.TA2AttributeType},
			Indices:       attributes,
		}
		attributeUpdate := NewAddSemanticTypeStep(nil, nil, attribs)
		wrapper := NewDatasetWrapperStep(map[string]DataRef{"inputs": &StepDataRef{offset - 1, "produce"}}, []string{"produce"}, offset, "")
		semanticTypeUpdates = append(semanticTypeUpdates, attributeUpdate, wrapper)
		offset += 2
	}

	// apply target semantic type if a target was found
	if targetIndex >= 0 {
		targetColumnUpdate := &ColumnUpdate{
			SemanticTypes: []string{model.TA2TargetType},
			Indices:       []int{targetIndex},
		}
		targetUpdate := NewAddSemanticTypeStep(nil, nil, targetColumnUpdate)
		wrapper := NewDatasetWrapperStep(map[string]DataRef{"inputs": &StepDataRef{offset - 1, "produce"}}, []string{"produce"}, offset, "")
		semanticTypeUpdates = append(semanticTypeUpdates, targetUpdate, wrapper)
	}

	return semanticTypeUpdates, nil
}

func filterBySet(filterSets []*model.FilterSet, columnIndices map[string]int, offset int) ([]Step, error) {
	// only support exclusive filters
	for _, fs := range filterSets {
		if !fs.IsValid() {
			return nil, errors.Errorf("invalid filter set detected")
		}

		for _, ff := range fs.FeatureFilters {
			for _, f := range ff.List {
				if f.Mode == model.IncludeFilter {
					return nil, errors.Errorf("pipeline filtering only supports exclusive filters")
				}
			}
		}
	}

	// Handle each filter set as an independent set of filters to exclude.
	// A filter set is composed of filter objects joined together using logical ANDs.
	// Each filter object contains filters acting on a single feature joined using logical ORs.
	// Since the filtes are exclusion filters, the logical joining operator is reversed.
	// Filters can be ORed together by having each filter operate on the same input
	// then combining the output (removing duplicates). Filters can be ANDed together
	// by having the output of one filter be the input of the following filter.
	steps := []Step{}

	// all filter sets act on the same input
	for _, fs := range filterSets {
		offsetFS := offset
		filterOutputs := &ListStepDataRef{[]DataRef{}}
		for _, ff := range fs.FeatureFilters {
			// add necessary filters for each filter
			stepsFeature, err := createFilterData(ff.List, columnIndices, offsetFS, offset)
			if err != nil {
				return nil, err
			}
			steps = append(steps, stepsFeature...)
			offset = offset + len(stepsFeature)
			filterOutputs.AddDataRef(&StepDataRef{offset - 1, "produce"})
		}

		// combine the outputs, removing duplicates
		steps = append(steps, NewVerticalConcatenationPrimitiveStep(map[string]DataRef{"inputs": filterOutputs}, []string{"produce"}, true))
		offset += 2
	}

	return steps, nil
}

func createFilterData(filters []*model.Filter, columnIndices map[string]int, offsetInput int, offsetStep int) ([]Step, error) {

	// Map the fiters to pipeline primitives
	filterSteps := []Step{}
	for _, f := range filters {
		var filter Step
		inclusive := f.Mode == model.IncludeFilter
		filterKey := strings.ToLower(f.Key)
		colIndex, ok := columnIndices[filterKey]
		if !ok {
			return nil, errors.Errorf("filter field %s not found in columns", filterKey)
		}

		switch f.Type {
		case model.NumericalFilter:
			filter = NewNumericRangeFilterStep(nil, nil, colIndex, inclusive, *f.Min, *f.Max, false)
			wrapper := NewDatasetWrapperStep(map[string]DataRef{"inputs": &StepDataRef{offsetInput - 1, "produce"}}, []string{"produce"}, offsetStep, "")
			filterSteps = append(filterSteps, filter, wrapper)
			offsetStep += 2

		case model.CategoricalFilter:
			filter = NewTermFilterStep(nil, nil, colIndex, inclusive, f.Categories, true)
			wrapper := NewDatasetWrapperStep(map[string]DataRef{"inputs": &StepDataRef{offsetInput - 1, "produce"}}, []string{"produce"}, offsetStep, "")
			filterSteps = append(filterSteps, filter, wrapper)
			offsetStep += 2

		case model.BivariateFilter:
			split := strings.Split(filterKey, ":")
			xCol := split[0]
			yCol := split[1]
			xColIndex := columnIndices[xCol]
			yColIndex := columnIndices[yCol]

			filter = NewNumericRangeFilterStep(nil, nil, xColIndex, inclusive, f.Bounds.MinX, f.Bounds.MaxX, false)
			wrapper := NewDatasetWrapperStep(map[string]DataRef{"inputs": &StepDataRef{offsetInput - 1, "produce"}}, []string{"produce"}, offsetStep, "")
			filterSteps = append(filterSteps, filter, wrapper)

			filter = NewNumericRangeFilterStep(nil, nil, yColIndex, inclusive, f.Bounds.MinY, f.Bounds.MaxY, false)
			wrapper = NewDatasetWrapperStep(map[string]DataRef{"inputs": &StepDataRef{offsetInput - 1, "produce"}}, []string{"produce"}, offsetStep, "")
			filterSteps = append(filterSteps, filter, wrapper)

			offsetStep += 4

		case model.RowFilter:
			filter = NewTermFilterStep(nil, nil, colIndex, inclusive, f.D3mIndices, true)
			wrapper := NewDatasetWrapperStep(map[string]DataRef{"inputs": &StepDataRef{offsetInput - 1, "produce"}}, []string{"produce"}, offsetStep, "")
			filterSteps = append(filterSteps, filter, wrapper)
			offsetStep += 2

		case model.TextFilter:
			filter = NewTermFilterStep(nil, nil, colIndex, inclusive, f.Categories, false)
			wrapper := NewDatasetWrapperStep(map[string]DataRef{"inputs": &StepDataRef{offsetInput - 1, "produce"}}, []string{"produce"}, offsetStep, "")
			filterSteps = append(filterSteps, filter, wrapper)
			offsetStep += 2

		case model.DatetimeFilter:
			filter = NewDateTimeRangeFilterStep(nil, nil, colIndex, inclusive, *f.Min, *f.Max, false)
			wrapper := NewDatasetWrapperStep(map[string]DataRef{"inputs": &StepDataRef{offsetInput - 1, "produce"}}, []string{"produce"}, offsetStep, "")
			filterSteps = append(filterSteps, filter, wrapper)
			offsetStep += 2

		case model.GeoBoundsFilter:
			// This is a filter that assumes to be working on a vector that contains 4 points defining a geographic area.  The vector
			// is defined as [x0, y0, x1, y1, x2, y2, x3, y3], where (x0, y0) is the LL corner of the bounds, and the points are ordered in a
			// clockwise fashion.

			minX := f.Bounds.MinX
			maxX := f.Bounds.MaxX
			minY := f.Bounds.MinY
			maxY := f.Bounds.MaxY
			minValues := []float64{minX, minY, minX, minY, minX, minY, minX, minY}
			maxValues := []float64{maxX, maxY, maxX, maxY, maxX, maxY, maxX, maxY}
			filter = NewVectorBoundsFilterStep(nil, nil, colIndex, inclusive, minValues, maxValues, false)
			wrapper := NewDatasetWrapperStep(map[string]DataRef{"inputs": &StepDataRef{offsetInput - 1, "produce"}}, []string{"produce"}, offsetStep, "")
			filterSteps = append(filterSteps, filter, wrapper)

			offsetStep += 2
		}
		offsetInput = offsetStep
	}
	return filterSteps, nil
}

func mapColumns(allFeatures []*model.Variable, selectedSet map[string]bool) map[string]int {
	colIndices := make(map[string]int)
	for _, f := range allFeatures {
		key := strings.ToLower(f.Key)
		if selectedSet[key] {
			colIndices[key] = f.Index
		}
	}

	return colIndices
}

func listColumns(allFeatures []*model.Variable, selectedSet map[string]bool) []int {
	colIndices := make([]int, 0)
	for i := 0; i < len(allFeatures); i++ {
		if selectedSet[strings.ToLower(allFeatures[i].Key)] {
			colIndices = append(colIndices, allFeatures[i].Index)
		}
	}

	return colIndices
}
