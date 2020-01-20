//
//   Copyright © 2019 Uncharted Software Inc.
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
	Filters          []*model.Filter
}

// UserDatasetAugmentation contains the augmentation parameters required
// for user dataset pipelines.
type UserDatasetAugmentation struct {
	SearchResult  string
	SystemID      string
	BaseDatasetID string
}

// CreateUserDatasetPipeline creates a pipeline description to capture user feature selection and
// semantic type information.
func CreateUserDatasetPipeline(name string, description string, datasetDescription *UserDatasetDescription,
	augmentations []*UserDatasetAugmentation) (*pipeline.PipelineDescription, error) {

	offset := 0

	// save the selected features in a set for quick lookup
	selectedSet := map[string]bool{}
	for _, v := range datasetDescription.SelectedFeatures {
		selectedSet[strings.ToLower(v)] = true
	}
	columnIndices := mapColumns(datasetDescription.AllFeatures, selectedSet)

	// create pipeline nodes for step we need to execute
	steps := []Step{} // add the denorm primitive

	// determine if this is a timeseries dataset
	isTimeseries := false
	groupingIndices := make([]int, 0)
	timeseriesGrouping := getTimeseriesGrouping(datasetDescription)
	if timeseriesGrouping != nil {
		isTimeseries = true
		groupingSet := map[string]bool{}

		// we need to udpate the selected set to include members of the grouped variable
		for _, subID := range timeseriesGrouping.SubIDs {
			selectedSet[strings.ToLower(subID)] = true
			groupingSet[strings.ToLower(subID)] = true
		}

		groupingIndices = listColumns(datasetDescription.AllFeatures, groupingSet)
		selectedSet[strings.ToLower(timeseriesGrouping.Properties.XCol)] = true
		selectedSet[strings.ToLower(timeseriesGrouping.Properties.YCol)] = true
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
		// need to read csv data, flatten then concat back to the original pipeline
		steps = append(steps, NewTimeseriesFormatterStep(map[string]DataRef{"inputs": dataRef}, []string{"produce"}, "", -1))
		steps = append(steps, NewColumnParserStep(nil, nil, []string{model.TA2IntegerType, model.TA2BooleanType, model.TA2RealType}))
		steps = append(steps, NewDatasetWrapperStep(map[string]DataRef{"inputs": &StepDataRef{offset, "produce"}}, []string{"produce"}, offset+1, ""))
		steps = append(steps, NewGroupingFieldComposeStep(nil, nil, groupingIndices, "-", "__grouping"))
		steps = append(steps, NewDatasetWrapperStep(map[string]DataRef{"inputs": &StepDataRef{offset + 2, "produce"}}, []string{"produce"}, offset+3, ""))
		offset += 5
	} else {
		steps = append(steps, NewDenormalizeStep(map[string]DataRef{"inputs": dataRef}, []string{"produce"}))
		steps = append(steps, NewColumnParserStep(nil, nil, []string{model.TA2IntegerType, model.TA2BooleanType, model.TA2RealType}))
		steps = append(steps, NewDatasetWrapperStep(map[string]DataRef{"inputs": &StepDataRef{offset, "produce"}}, []string{"produce"}, offset+1, ""))
		steps = append(steps, NewDataCleaningStep(nil, nil))
		steps = append(steps, NewDatasetWrapperStep(map[string]DataRef{"inputs": &StepDataRef{offset + 2, "produce"}}, []string{"produce"}, offset+3, ""))
		offset += 5
	}

	// create the semantic type update primitive
	updateSemanticTypes, err := createUpdateSemanticTypes(datasetDescription.AllFeatures, selectedSet, offset)
	if err != nil {
		return nil, err
	}
	steps = append(steps, updateSemanticTypes...)
	offset += len(updateSemanticTypes)

	// create the feature selection primitive
	removeFeatures := createRemoveFeatures(datasetDescription.AllFeatures, selectedSet, offset)
	steps = append(steps, removeFeatures...)
	offset += len(removeFeatures)

	// add filter primitives
	filterData := createFilterData(datasetDescription.Filters, columnIndices, offset)
	steps = append(steps, filterData...)
	offset += len(filterData)

	// If neither have any content, we'll skip the template altogether.
	if len(updateSemanticTypes) == 0 && removeFeatures == nil &&
		len(filterData) == 0 && augmentations == nil && !isTimeseries {
		return nil, nil
	}

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

func getTimeseriesGrouping(datasetDescription *UserDatasetDescription) *model.Grouping {
	if model.IsTimeSeries(datasetDescription.TargetFeature.Type) {
		return datasetDescription.TargetFeature.Grouping
	}
	for _, v := range datasetDescription.AllFeatures {
		if v.Grouping != nil && model.IsTimeSeries(v.Grouping.Type) {
			return v.Grouping
		}
	}

	return nil
}

func createRemoveFeatures(allFeatures []*model.Variable, selectedSet map[string]bool, offset int) []Step {
	// create a list of features to remove
	removeFeatures := []int{}
	for _, v := range allFeatures {
		if !selectedSet[strings.ToLower(v.Name)] {
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

func createUpdateSemanticTypes(allFeatures []*model.Variable, selectedSet map[string]bool, offset int) ([]Step, error) {
	// create maps of (semantic type, index list) - primitive allows for semantic types to be added to /
	// remove from multiple columns in a single operation
	updateMap := map[string]*update{}
	attributes := make([]int, 0)
	for _, v := range allFeatures {
		// empty selected set means all selected
		if len(selectedSet) == 0 || selectedSet[strings.ToLower(v.Name)] {
			addType := model.MapTA2Type(v.Type)
			if addType == "" {
				return nil, errors.Errorf("variable `%s` internal type `%s` can't be mapped to ta2", v.Name, v.Type)
			}
			removeType := model.MapTA2Type(v.OriginalType)
			if removeType == "" {
				return nil, errors.Errorf("remove variable `%s` internal type `%s` can't be mapped to ta2", v.Name, v.OriginalType)
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
		}

		// update all non target to attribute
		attributes = append(attributes, v.Index)
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

	if len(attributes) > 0 {
		attribs := &ColumnUpdate{
			SemanticTypes: []string{model.TA2AttributeType},
			Indices:       attributes,
		}
		attributeUpdate := NewAddSemanticTypeStep(nil, nil, attribs)
		wrapper := NewDatasetWrapperStep(map[string]DataRef{"inputs": &StepDataRef{offset - 1, "produce"}}, []string{"produce"}, offset, "")
		semanticTypeUpdates = append(semanticTypeUpdates, attributeUpdate, wrapper)
	}
	return semanticTypeUpdates, nil
}

func createFilterData(filters []*model.Filter, columnIndices map[string]int, offset int) []Step {

	// Map the fiters to pipeline primitives
	filterSteps := []Step{}
	for _, f := range filters {
		var filter Step
		inclusive := f.Mode == model.IncludeFilter
		colIndex := columnIndices[f.Key]

		switch f.Type {
		case model.NumericalFilter:
			filter = NewNumericRangeFilterStep(nil, nil, colIndex, inclusive, *f.Min, *f.Max, false)
			wrapper := NewDatasetWrapperStep(map[string]DataRef{"inputs": &StepDataRef{offset - 1, "produce"}}, []string{"produce"}, offset, "")
			filterSteps = append(filterSteps, filter, wrapper)
			offset += 2

		case model.CategoricalFilter:
			filter = NewTermFilterStep(nil, nil, colIndex, inclusive, f.Categories, true)
			wrapper := NewDatasetWrapperStep(map[string]DataRef{"inputs": &StepDataRef{offset - 1, "produce"}}, []string{"produce"}, offset, "")
			filterSteps = append(filterSteps, filter, wrapper)
			offset += 2

		case model.BivariateFilter:
			split := strings.Split(f.Key, ":")
			xCol := split[0]
			yCol := split[1]
			xColIndex := columnIndices[xCol]
			yColIndex := columnIndices[yCol]

			filter = NewNumericRangeFilterStep(nil, nil, xColIndex, inclusive, f.Bounds.MinX, f.Bounds.MaxX, false)
			wrapper := NewDatasetWrapperStep(map[string]DataRef{"inputs": &StepDataRef{offset - 1, "produce"}}, []string{"produce"}, offset, "")
			filterSteps = append(filterSteps, filter, wrapper)

			filter = NewNumericRangeFilterStep(nil, nil, yColIndex, inclusive, f.Bounds.MinY, f.Bounds.MaxY, false)
			wrapper = NewDatasetWrapperStep(map[string]DataRef{"inputs": &StepDataRef{offset - 1, "produce"}}, []string{"produce"}, offset, "")
			filterSteps = append(filterSteps, filter, wrapper)

			offset += 4

		case model.RowFilter:
			filter = NewTermFilterStep(nil, nil, colIndex, inclusive, f.D3mIndices, true)
			wrapper := NewDatasetWrapperStep(map[string]DataRef{"inputs": &StepDataRef{offset - 1, "produce"}}, []string{"produce"}, offset, "")
			filterSteps = append(filterSteps, filter, wrapper)
			offset += 2

		case model.FeatureFilter, model.TextFilter:
			filter = NewTermFilterStep(nil, nil, colIndex, inclusive, f.Categories, false)
			wrapper := NewDatasetWrapperStep(map[string]DataRef{"inputs": &StepDataRef{offset - 1, "produce"}}, []string{"produce"}, offset, "")
			filterSteps = append(filterSteps, filter, wrapper)
			offset += 2
		}

	}
	return filterSteps
}

func getSemanticTypeUpdates(v *model.Variable, inputIndex int, offset int) []Step {
	addType := model.MapTA2Type(v.Type)
	removeType := model.MapTA2Type(v.OriginalType)

	add := &ColumnUpdate{
		SemanticTypes: []string{addType},
		Indices:       []int{v.Index},
	}
	remove := &ColumnUpdate{
		SemanticTypes: []string{removeType},
		Indices:       []int{v.Index},
	}
	return []Step{
		NewAddSemanticTypeStep(nil, nil, add),
		NewDatasetWrapperStep(map[string]DataRef{"inputs": &StepDataRef{inputIndex, "produce"}}, []string{"produce"}, offset, ""),
		NewRemoveSemanticTypeStep(nil, nil, remove),
		NewDatasetWrapperStep(map[string]DataRef{"inputs": &StepDataRef{offset + 1, "produce"}}, []string{"produce"}, offset+2, ""),
	}
}

func mapColumns(allFeatures []*model.Variable, selectedSet map[string]bool) map[string]int {
	colIndices := make(map[string]int)
	index := 0
	for _, f := range allFeatures {
		if selectedSet[strings.ToLower(f.Name)] {
			colIndices[f.Name] = index
			index = index + 1
		}
	}

	return colIndices
}

func listColumns(allFeatures []*model.Variable, selectedSet map[string]bool) []int {
	colIndices := make([]int, 0)
	for i := 0; i < len(allFeatures); i++ {
		if selectedSet[strings.ToLower(allFeatures[i].Name)] {
			colIndices = append(colIndices, allFeatures[i].Index)
		}
	}

	return colIndices
}

func getIndex(allFeatures []*model.Variable, name string) (int, error) {
	for _, f := range allFeatures {
		if strings.EqualFold(name, f.Name) {
			return f.Index, nil
		}
	}
	return -1, errors.Errorf("can't find var '%s'", name)
}
