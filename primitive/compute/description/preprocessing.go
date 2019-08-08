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

const defaultResource = "learningData"

// UserDatasetDescription contains the basic parameters needs to generate
// the user dataset pipeline.
type UserDatasetDescription struct {
	AllFeatures      []*model.Variable
	TargetFeature    string
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
	// TODO: CSV reader is currently not working correctly, so we need to disable the timeseries
	// prepend and use the data as-is.  This is sufficient for now as the NK classification primitives
	// run off the unmodified dataset (although they can be configured to long form)
	// prend for now
	// for _, v := range datasetDescription.AllFeatures {
	// 	if v.Grouping != nil && model.IsTimeSeries(v.Grouping.Type) {
	// 		isTimeseries = true
	// 		break
	// 	}
	// }

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
		steps = append(steps, NewDatasetToDataframeStep(map[string]DataRef{"inputs": dataRef}, []string{"produce"}))
		steps = append(steps, NewDatasetToDataframeStepWithResource(map[string]DataRef{"inputs": dataRef}, []string{"produce"}, "0"))
		steps = append(steps, NewCSVReaderStep(map[string]DataRef{"inputs": &StepDataRef{offset + 1, "produce"}}, []string{"produce"}))
		steps = append(steps, NewHorizontalConcatStep(map[string]DataRef{"left": &StepDataRef{offset, "produce"}, "right": &StepDataRef{offset + 2, "produce"}}, []string{"produce"}, false, false))
		steps = append(steps, NewDataFrameFlattenStep(map[string]DataRef{"inputs": &StepDataRef{offset + 3, "produce"}}, []string{"produce"}))
		steps = append(steps, NewRemoveDuplicateColumnsStep(map[string]DataRef{"inputs": &StepDataRef{offset + 4, "produce"}}, []string{"produce"}))
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
		len(filterData) == 0 && augmentations == nil {
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

// CreateSlothPipeline creates a pipeline to peform timeseries clustering on a dataset.
func CreateSlothPipeline(name string, description string, timeColumn string, valueColumn string,
	timeSeriesFeatures []*model.Variable) (*pipeline.PipelineDescription, error) {

	// timeIdx, err := getIndex(timeSeriesFeatures, timeColumn)
	// if err != nil {
	// 	return nil, err
	// }

	// valueIdx, err := getIndex(timeSeriesFeatures, valueColumn)
	// if err != nil {
	// 	return nil, err
	// }

	inputs := []string{"inputs"}
	outputs := []DataRef{&StepDataRef{1, "produce"}}

	steps := []Step{
		NewDenormalizeStep(map[string]DataRef{"inputs": &PipelineDataRef{0}}, []string{"produce"}),
		// Sloth now includes the the time series loader in the primitive itself.
		// This is not a long term solution and will need updating.  The updated
		// primitive doesn't accept the time and value indices as args, so they
		// are currently unused.
		// step2 := NewPipelineNode(NewTimeSeriesLoaderStep(-1, timeIdx, valueIdx))
		NewSlothStep(map[string]DataRef{"inputs": &StepDataRef{0, "produce"}}, []string{"produce"}),
	}

	pipeline, err := NewPipelineBuilder(name, description, inputs, outputs, steps).Compile()
	if err != nil {
		return nil, err
	}
	return pipeline, nil
}

// CreateDukePipeline creates a pipeline to peform image featurization on a dataset.
func CreateDukePipeline(name string, description string) (*pipeline.PipelineDescription, error) {
	inputs := []string{"inputs"}
	outputs := []DataRef{&StepDataRef{1, "produce"}}

	steps := []Step{
		NewDatasetToDataframeStep(map[string]DataRef{"inputs": &PipelineDataRef{0}}, []string{"produce"}),
		NewDukeStep(map[string]DataRef{"inputs": &StepDataRef{0, "produce"}}, []string{"produce"}),
	}

	pipeline, err := NewPipelineBuilder(name, description, inputs, outputs, steps).Compile()
	if err != nil {
		return nil, err
	}
	return pipeline, nil
}

// CreateSimonPipeline creates a pipeline to run semantic type inference on a dataset's
// columns.
func CreateSimonPipeline(name string, description string) (*pipeline.PipelineDescription, error) {
	inputs := []string{"inputs"}
	outputs := []DataRef{&StepDataRef{1, "produce_metafeatures"}}

	steps := []Step{
		NewDatasetToDataframeStep(
			map[string]DataRef{"inputs": &PipelineDataRef{0}},
			[]string{"produce"},
		),
		NewSimonStep(
			map[string]DataRef{"inputs": &StepDataRef{0, "produce"}},
			[]string{"produce_metafeatures"}),
	}

	// produce metafeatures
	pipeline, err := NewPipelineBuilder(name, description, inputs, outputs, steps).Compile()
	if err != nil {
		return nil, err
	}
	return pipeline, nil
}

// CreateDataCleaningPipeline creates a pipeline to run data cleaning on a dataset.
func CreateDataCleaningPipeline(name string, description string) (*pipeline.PipelineDescription, error) {
	inputs := []string{"inputs"}
	outputs := []DataRef{&StepDataRef{2, "produce"}}

	steps := []Step{
		NewDatasetToDataframeStep(map[string]DataRef{"inputs": &PipelineDataRef{0}}, []string{"produce"}),
		NewColumnParserStep(
			map[string]DataRef{"inputs": &StepDataRef{0, "produce"}},
			[]string{"produce"},
			[]string{model.TA2IntegerType, model.TA2BooleanType, model.TA2RealType},
		),
		NewDataCleaningStep(map[string]DataRef{"inputs": &StepDataRef{1, "produce"}}, []string{"produce"}),
	}

	pipeline, err := NewPipelineBuilder(name, description, inputs, outputs, steps).Compile()
	if err != nil {
		return nil, err
	}
	return pipeline, nil
}

// CreateCrocPipeline creates a pipeline to run image featurization on a dataset.
func CreateCrocPipeline(name string, description string, targetColumns []string, outputLabels []string) (*pipeline.PipelineDescription, error) {
	inputs := []string{"inputs"}
	outputs := []DataRef{&StepDataRef{2, "produce"}}

	steps := []Step{
		NewDenormalizeStep(map[string]DataRef{"inputs": &PipelineDataRef{0}}, []string{"produce"}),
		NewDatasetToDataframeStep(map[string]DataRef{"inputs": &StepDataRef{0, "produce"}}, []string{"produce"}),
		NewCrocStep(map[string]DataRef{"inputs": &StepDataRef{1, "produce"}}, []string{"produce"}, targetColumns, outputLabels),
	}

	pipeline, err := NewPipelineBuilder(name, description, inputs, outputs, steps).Compile()
	if err != nil {
		return nil, err
	}
	return pipeline, nil
}

// CreateUnicornPipeline creates a pipeline to run image clustering on a dataset.
func CreateUnicornPipeline(name string, description string, targetColumns []string, outputLabels []string) (*pipeline.PipelineDescription, error) {
	inputs := []string{"inputs"}
	outputs := []DataRef{&StepDataRef{2, "produce"}}

	steps := []Step{
		NewDenormalizeStep(map[string]DataRef{"inputs": &PipelineDataRef{0}}, []string{"produce"}),
		NewDatasetToDataframeStep(map[string]DataRef{"inputs": &StepDataRef{0, "produce"}}, []string{"produce"}),
		NewUnicornStep(map[string]DataRef{"inputs": &StepDataRef{1, "produce"}}, []string{"produce"}, targetColumns, outputLabels),
	}

	pipeline, err := NewPipelineBuilder(name, description, inputs, outputs, steps).Compile()
	if err != nil {
		return nil, err
	}
	return pipeline, nil
}

// CreatePCAFeaturesPipeline creates a pipeline to run feature ranking on an input dataset.
func CreatePCAFeaturesPipeline(name string, description string) (*pipeline.PipelineDescription, error) {
	inputs := []string{"inputs"}
	outputs := []DataRef{&StepDataRef{1, "produce_metafeatures"}}

	steps := []Step{
		NewDatasetToDataframeStep(map[string]DataRef{"inputs": &PipelineDataRef{0}}, []string{"produce"}),
		NewPCAFeaturesStep(map[string]DataRef{"inputs": &StepDataRef{0, "produce"}}, []string{"produce_metafeatures"}),
	}

	pipeline, err := NewPipelineBuilder(name, description, inputs, outputs, steps).Compile()
	if err != nil {
		return nil, err
	}
	return pipeline, nil
}

// CreateDenormalizePipeline creates a pipeline to run the denormalize primitive on an input dataset.
func CreateDenormalizePipeline(name string, description string) (*pipeline.PipelineDescription, error) {
	inputs := []string{"inputs"}
	outputs := []DataRef{&StepDataRef{1, "produce"}}

	steps := []Step{
		NewDenormalizeStep(map[string]DataRef{"inputs": &PipelineDataRef{0}}, []string{"produce"}),
		NewDatasetToDataframeStep(map[string]DataRef{"inputs": &StepDataRef{0, "produce"}}, []string{"produce"}),
	}

	pipeline, err := NewPipelineBuilder(name, description, inputs, outputs, steps).Compile()
	if err != nil {
		return nil, err
	}
	return pipeline, nil
}

// CreateTargetRankingPipeline creates a pipeline to run feature ranking on an input dataset.
func CreateTargetRankingPipeline(name string, description string, target string, features []*model.Variable) (*pipeline.PipelineDescription, error) {

	// compute index associated with column name
	targetIdx := -1
	for _, f := range features {
		if strings.EqualFold(target, f.Name) {
			targetIdx = f.Index
			break
		}
	}
	if targetIdx < 0 {
		return nil, errors.Errorf("can't find var '%s'", name)
	}

	offset := 0
	steps := []Step{
		NewDenormalizeStep(map[string]DataRef{"inputs": &PipelineDataRef{0}}, []string{"produce"}),
	}
	offset++

	// ranking is dependent on user updated semantic types, so we need to make sure we apply
	// those to the original data
	updateSemanticTypeStep, err := createUpdateSemanticTypes(features, map[string]bool{}, offset)
	if err != nil {
		return nil, err
	}
	for _, s := range updateSemanticTypeStep {
		steps = append(steps, s)
	}

	offset += len(updateSemanticTypeStep)
	steps = append(steps,
		NewDatasetToDataframeStep(map[string]DataRef{"inputs": &StepDataRef{offset - 1, "produce"}}, []string{"produce"}),
		NewColumnParserStep(
			map[string]DataRef{"inputs": &StepDataRef{offset, "produce"}},
			[]string{"produce"},
			// inlcude categorical because the parser will hash all the values into ints, which the MIRanking primitive can handle
			[]string{model.TA2IntegerType, model.TA2BooleanType, model.TA2RealType, model.TA2CategoricalType},
		),
		NewTargetRankingStep(map[string]DataRef{"inputs": &StepDataRef{offset + 1, "produce"}}, []string{"produce"}, targetIdx),
	)
	offset += 3

	inputs := []string{"inputs"}
	outputs := []DataRef{&StepDataRef{offset - 1, "produce"}}

	pipeline, err := NewPipelineBuilder(name, description, inputs, outputs, steps).Compile()
	if err != nil {
		return nil, err
	}
	return pipeline, nil
}

// CreateGoatForwardPipeline creates a forward geocoding pipeline.
func CreateGoatForwardPipeline(name string, description string, placeCol *model.Variable) (*pipeline.PipelineDescription, error) {
	inputs := []string{"inputs"}
	outputs := []DataRef{&StepDataRef{2, "produce"}}

	steps := []Step{
		NewDenormalizeStep(map[string]DataRef{"inputs": &PipelineDataRef{0}}, []string{"produce"}),
		NewDatasetToDataframeStep(map[string]DataRef{"inputs": &StepDataRef{0, "produce"}}, []string{"produce"}),
		NewGoatForwardStep(map[string]DataRef{"inputs": &StepDataRef{0, "produce"}}, []string{"produce"}, placeCol.Index),
	}
	pipeline, err := NewPipelineBuilder(name, description, inputs, outputs, steps).Compile()

	if err != nil {
		return nil, err
	}
	return pipeline, nil
}

// CreateGoatReversePipeline creates a forward geocoding pipeline.
func CreateGoatReversePipeline(name string, description string, lonSource *model.Variable, latSource *model.Variable) (*pipeline.PipelineDescription, error) {
	inputs := []string{"inputs"}
	outputs := []DataRef{&StepDataRef{2, "produce"}}

	steps := []Step{
		NewDenormalizeStep(map[string]DataRef{"inputs": &PipelineDataRef{0}}, []string{"produce"}),
		NewDatasetToDataframeStep(map[string]DataRef{"inputs": &StepDataRef{0, "produce"}}, []string{"produce"}),
		NewGoatReverseStep(map[string]DataRef{"inputs": &StepDataRef{0, "produce"}}, []string{"produce"}, lonSource.Index, latSource.Index),
	}

	pipeline, err := NewPipelineBuilder(name, description, inputs, outputs, steps).Compile()
	if err != nil {
		return nil, err
	}
	return pipeline, nil
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

// CreateJoinPipeline creates a pipeline that joins two input datasets using a caller supplied column.
// Accuracy is a normalized value that controls how exact the join has to be.
func CreateJoinPipeline(name string, description string, leftJoinCol *model.Variable, rightJoinCol *model.Variable, accuracy float32) (*pipeline.PipelineDescription, error) {
	steps := make([]Step, 0)
	steps = append(steps, NewDenormalizeStep(map[string]DataRef{"inputs": &PipelineDataRef{0}}, []string{"produce"}))
	steps = append(steps, NewDenormalizeStep(map[string]DataRef{"inputs": &PipelineDataRef{1}}, []string{"produce"}))
	offset := 2
	offsetLeft := 0
	offsetRight := 1

	// update semantic types as needed
	if leftJoinCol.Type != leftJoinCol.OriginalType {
		stepsRetype := getSemanticTypeUpdates(leftJoinCol, 0, offset)
		steps = append(steps, stepsRetype...)
		offset += len(stepsRetype)
		offsetLeft = offset - 1
	}
	if rightJoinCol.Type != rightJoinCol.OriginalType {
		stepsRetype := getSemanticTypeUpdates(rightJoinCol, 1, offset)
		steps = append(steps, stepsRetype...)
		offset += len(stepsRetype)
		offsetRight = offset - 1
	}

	// merge two intput streams via a single join call
	steps = append(steps, NewJoinStep(
		map[string]DataRef{"left": &StepDataRef{offsetLeft, "produce"}, "right": &StepDataRef{offsetRight, "produce"}},
		[]string{"produce"},
		leftJoinCol.DisplayName, rightJoinCol.DisplayName, accuracy,
	))
	steps = append(steps, NewDatasetToDataframeStep(map[string]DataRef{"inputs": &StepDataRef{offset, "produce"}}, []string{"produce"}))

	// compute column indices
	inputs := []string{"left", "right"}
	outputs := []DataRef{&StepDataRef{len(steps) - 1, "produce"}}

	pipeline, err := NewPipelineBuilder(name, description, inputs, outputs, steps).Compile()
	if err != nil {
		return nil, err
	}
	return pipeline, nil
}

// CreateDSBoxJoinPipeline creates a pipeline that joins two input datasets
// using caller supplied columns.
func CreateDSBoxJoinPipeline(name string, description string, leftJoinCols []string, rightJoinCols []string, accuracy float32) (*pipeline.PipelineDescription, error) {
	inputs := []string{"inputs"}
	outputs := []DataRef{&StepDataRef{2, "produce"}}

	// instantiate the pipeline - this merges two intput streams via a single join call
	steps := []Step{
		NewDenormalizeStep(map[string]DataRef{"inputs": &PipelineDataRef{0}}, []string{"produce"}),
		NewDenormalizeStep(map[string]DataRef{"inputs": &PipelineDataRef{1}}, []string{"produce"}),
		NewDSBoxJoinStep(
			map[string]DataRef{"inputs": &StepDataRef{1, "produce"}},
			[]string{"produce"},
			leftJoinCols, rightJoinCols, accuracy),
		NewDatasetToDataframeStep(map[string]DataRef{"inputs": &StepDataRef{2, "produce"}}, []string{"produce"}),
	}

	pipeline, err := NewPipelineBuilder(name, description, inputs, outputs, steps).Compile()
	if err != nil {
		return nil, err
	}
	return pipeline, nil
}

// CreateTimeseriesFormatterPipeline creates a time series formatter pipeline.
func CreateTimeseriesFormatterPipeline(name string, description string, resourceId string) (*pipeline.PipelineDescription, error) {
	inputs := []string{"inputs"}
	outputs := []DataRef{&StepDataRef{5, "produce"}}

	steps := []Step{
		NewDatasetToDataframeStep(map[string]DataRef{"inputs": &PipelineDataRef{0}}, []string{"produce"}),
		NewDatasetToDataframeStepWithResource(map[string]DataRef{"inputs": &PipelineDataRef{0}}, []string{"produce"}, resourceId),
		NewCSVReaderStep(map[string]DataRef{"inputs": &StepDataRef{1, "produce"}}, []string{"produce"}),
		NewHorizontalConcatStep(map[string]DataRef{"left": &StepDataRef{0, "produce"}, "right": &StepDataRef{2, "produce"}}, []string{"produce"}, false, false),
		NewDataFrameFlattenStep(map[string]DataRef{"inputs": &StepDataRef{3, "produce"}}, []string{"produce"}),
		NewRemoveDuplicateColumnsStep(map[string]DataRef{"inputs": &StepDataRef{4, "produce"}}, []string{"produce"}),
	}

	pipeline, err := NewPipelineBuilder(name, description, inputs, outputs, steps).Compile()
	if err != nil {
		return nil, err
	}
	return pipeline, nil
}

// CreateDatamartDownloadPipeline creates a pipeline to download data from a datamart.
func CreateDatamartDownloadPipeline(name string, description string, searchResult string, systemIdentifier string) (*pipeline.PipelineDescription, error) {
	inputs := []string{"inputs"}
	outputs := []DataRef{&StepDataRef{1, "produce"}}

	steps := []Step{
		NewDatamartDownloadStep(map[string]DataRef{"inputs": &PipelineDataRef{0}}, []string{"produce"}, searchResult, systemIdentifier),
		NewDatasetToDataframeStep(map[string]DataRef{"inputs": &StepDataRef{0, "produce"}}, []string{"produce"}),
	}

	pipeline, err := NewPipelineBuilder(name, description, inputs, outputs, steps).Compile()
	if err != nil {
		return nil, err
	}
	return pipeline, nil
}

// CreateDatamartAugmentPipeline creates a pipeline to augment data with datamart data.
func CreateDatamartAugmentPipeline(name string, description string, searchResult string, systemIdentifier string) (*pipeline.PipelineDescription, error) {
	inputs := []string{"inputs"}
	outputs := []DataRef{&StepDataRef{1, "produce"}}

	steps := []Step{
		NewDatamartAugmentStep(map[string]DataRef{"inputs": &PipelineDataRef{0}}, []string{"produce"}, searchResult, systemIdentifier),
		NewDatasetToDataframeStep(map[string]DataRef{"inputs": &StepDataRef{0, "produce"}}, []string{"produce"}),
	}

	pipeline, err := NewPipelineBuilder(name, description, inputs, outputs, steps).Compile()
	if err != nil {
		return nil, err
	}
	return pipeline, nil
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

func getIndex(allFeatures []*model.Variable, name string) (int, error) {
	for _, f := range allFeatures {
		if strings.EqualFold(name, f.Name) {
			return f.Index, nil
		}
	}
	return -1, errors.Errorf("can't find var '%s'", name)
}
