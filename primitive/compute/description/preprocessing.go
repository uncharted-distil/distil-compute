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

// CreateUserDatasetPipeline creates a pipeline description to capture user feature selection and
// semantic type information.
func CreateUserDatasetPipeline(name string, description string, allFeatures []*model.Variable,
	targetFeature string, selectedFeatures []string, filters []*model.Filter) (*pipeline.PipelineDescription, error) {

	// save the selected features in a set for quick lookup
	selectedSet := map[string]bool{}
	for _, v := range selectedFeatures {
		selectedSet[strings.ToLower(v)] = true
	}
	columnIndices := mapColumns(allFeatures, selectedSet)

	// create the semantic type update primitive
	updateSemanticTypes, err := createUpdateSemanticTypes(allFeatures, selectedSet)
	if err != nil {
		return nil, err
	}

	// create the feature selection primitive
	removeFeatures, err := createRemoveFeatures(allFeatures, selectedSet)
	if err != nil {
		return nil, err
	}

	// add filter primitives
	filterData := createFilterData(filters, columnIndices)

	// If neither have any content, we'll skip the template altogether.
	if len(updateSemanticTypes) == 0 && removeFeatures == nil && len(filterData) == 0 {
		return nil, nil
	}

	// create pipeline nodes for step we need to execute
	nodes := []*PipelineNode{}
	for _, v := range updateSemanticTypes {
		nodes = append(nodes, NewPipelineNode(v))
	}
	if removeFeatures != nil {
		nodes = append(nodes, NewPipelineNode(removeFeatures))
	}
	for _, f := range filterData {
		nodes = append(nodes, NewPipelineNode(f))
	}
	// mark this is a preprocessing template
	nodes = append(nodes, NewPipelineNode(NewInferenceStepData()))

	sourceNode := nodesToGraph(nodes)

	pip, err := NewPipelineBuilder(name, description, sourceNode).Compile()
	if err != nil {
		return nil, err
	}

	// Input set to arbitrary string for now
	pip.Inputs = []*pipeline.PipelineDescriptionInput{{
		Name: "dataset",
	}}
	return pip, nil
}

func createRemoveFeatures(allFeatures []*model.Variable, selectedSet map[string]bool) (*StepData, error) {
	// create a list of features to remove
	removeFeatures := []int{}
	for _, v := range allFeatures {
		if !selectedSet[strings.ToLower(v.Name)] {
			removeFeatures = append(removeFeatures, v.Index)
		}
	}

	if len(removeFeatures) == 0 {
		return nil, nil
	}

	// instantiate the feature remove primitive
	featureSelect, err := NewRemoveColumnsStep("", removeFeatures)
	if err != nil {
		return nil, err
	}
	return featureSelect, nil
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

func createUpdateSemanticTypes(allFeatures []*model.Variable, selectedSet map[string]bool) ([]*StepData, error) {
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

	semanticTypeUpdates := []*StepData{}
	for _, k := range keys {
		v := updateMap[k]
		var addKey string
		if len(v.addIndices) > 0 {
			addKey = k
		}
		add := &ColumnUpdate{
			SemanticTypes: []string{addKey},
			Indices:       v.addIndices,
		}
		var removeKey string
		if len(v.removeIndices) > 0 {
			removeKey = k
		}
		remove := &ColumnUpdate{
			SemanticTypes: []string{removeKey},
			Indices:       v.removeIndices,
		}
		semanticTypeUpdate, err := NewUpdateSemanticTypeStep("", add, remove)
		if err != nil {
			return nil, err
		}
		semanticTypeUpdates = append(semanticTypeUpdates, semanticTypeUpdate)
	}
	return semanticTypeUpdates, nil
}

func createFilterData(filters []*model.Filter, columnIndices map[string]int) []*StepData {

	// Map the fiters to pipeline primitives
	filterSteps := []*StepData{}
	for _, f := range filters {
		var filter *StepData
		inclusive := f.Mode == model.IncludeFilter
		colIndex := columnIndices[f.Key]

		switch f.Type {
		case model.NumericalFilter:
			filter = NewNumericRangeFilterStep("", colIndex, inclusive, *f.Min, *f.Max, false)
		case model.CategoricalFilter:
			filter = NewTermFilterStep("", colIndex, inclusive, f.Categories, true)
		case model.RowFilter:
			filter = NewTermFilterStep("", colIndex, inclusive, f.D3mIndices, true)
		case model.FeatureFilter, model.TextFilter:
			filter = NewTermFilterStep("", colIndex, inclusive, f.Categories, false)
		}

		filterSteps = append(filterSteps, filter)
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

	step0 := NewPipelineNode(NewDenormalizeStep())
	step1 := NewPipelineNode(NewDatasetToDataframeStep())
	// Sloth now includes the the time series loader in the primitive itself.
	// This is not a long term solution and will need updating.  The updated
	// primitive doesn't accept the time and value indices as args, so they
	// are currently unused.
	// step2 := NewPipelineNode(NewTimeSeriesLoaderStep(-1, timeIdx, valueIdx))
	step2 := NewPipelineNode(NewSlothStep())
	step0.Add(step1)
	step1.Add(step2)

	pipeline, err := NewPipelineBuilder(name, description, step0).Compile()
	if err != nil {
		return nil, err
	}
	return pipeline, nil
}

// CreateDukePipeline creates a pipeline to peform image featurization on a dataset.
func CreateDukePipeline(name string, description string) (*pipeline.PipelineDescription, error) {
	step0 := NewPipelineNode(NewDatasetToDataframeStep())
	step1 := NewPipelineNode(NewDukeStep())
	step0.Add(step1)

	pipeline, err := NewPipelineBuilder(name, description, step0).Compile()
	if err != nil {
		return nil, err
	}
	return pipeline, nil
}

// CreateSimonPipeline creates a pipeline to run semantic type inference on a dataset's
// columns.
func CreateSimonPipeline(name string, description string) (*pipeline.PipelineDescription, error) {
	step0 := NewPipelineNode(NewDatasetToDataframeStep())
	step1 := NewPipelineNode(NewSimonStep())
	step0.Add(step1)

	pipeline, err := NewPipelineBuilder(name, description, step0).Compile()
	if err != nil {
		return nil, err
	}
	return pipeline, nil
}

// CreateDataCleaningPipeline creates a pipeline to run data cleaning on a dataset.
func CreateDataCleaningPipeline(name string, description string) (*pipeline.PipelineDescription, error) {
	step0 := NewPipelineNode(NewDatasetToDataframeStep())
	step1 := NewPipelineNode(NewDataCleaningStep())
	step0.Add(step1)

	pipeline, err := NewPipelineBuilder(name, description, step0).Compile()
	if err != nil {
		return nil, err
	}
	return pipeline, nil
}

// CreateCrocPipeline creates a pipeline to run image featurization on a dataset.
func CreateCrocPipeline(name string, description string, targetColumns []string, outputLabels []string) (*pipeline.PipelineDescription, error) {
	step0 := NewPipelineNode(NewDenormalizeStep())
	step1 := NewPipelineNode(NewDatasetToDataframeStep())
	step2 := NewPipelineNode(NewCrocStep(targetColumns, outputLabels))
	step0.Add(step1)
	step1.Add(step2)

	pipeline, err := NewPipelineBuilder(name, description, step0).Compile()
	if err != nil {
		return nil, err
	}
	return pipeline, nil
}

// CreateUnicornPipeline creates a pipeline to run image clustering on a dataset.
func CreateUnicornPipeline(name string, description string, targetColumns []string, outputLabels []string) (*pipeline.PipelineDescription, error) {
	step0 := NewPipelineNode(NewDenormalizeStep())
	step1 := NewPipelineNode(NewDatasetToDataframeStep())
	step2 := NewPipelineNode(NewUnicornStep(targetColumns, outputLabels))
	step0.Add(step1)
	step1.Add(step2)

	pipeline, err := NewPipelineBuilder(name, description, step0).Compile()
	if err != nil {
		return nil, err
	}
	return pipeline, nil
}

// CreatePCAFeaturesPipeline creates a pipeline to run feature ranking on an input dataset.
func CreatePCAFeaturesPipeline(name string, description string) (*pipeline.PipelineDescription, error) {
	step0 := NewPipelineNode(NewDatasetToDataframeStep())
	step1 := NewPipelineNode(NewPCAFeaturesStep())
	step0.Add(step1)

	pipeline, err := NewPipelineBuilder(name, description, step0).Compile()
	if err != nil {
		return nil, err
	}
	return pipeline, nil
}

// CreateDenormalizePipeline creates a pipeline to run the denormalize primitive on an input dataset.
func CreateDenormalizePipeline(name string, description string) (*pipeline.PipelineDescription, error) {
	step0 := NewPipelineNode(NewDenormalizeStep())
	step1 := NewPipelineNode(NewDatasetToDataframeStep())
	step0.Add(step1)

	pipeline, err := NewPipelineBuilder(name, description, step0).Compile()
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

	nodes := []*PipelineNode{
		NewPipelineNode(NewDenormalizeStep()),
	}

	// ranking is dependent on user updated semantic types, so we need to make sure we apply
	// those to the original data
	updateSemanticTypeStep, err := createUpdateSemanticTypes(features, map[string]bool{})
	if err != nil {
		return nil, err
	}
	for _, s := range updateSemanticTypeStep {
		nodes = append(nodes, NewPipelineNode(s))
	}

	nodes = append(nodes,
		NewPipelineNode(NewDatasetToDataframeStep()),
		NewPipelineNode(NewColumnParserStep()),
		NewPipelineNode(NewTargetRankingStep(targetIdx)),
	)
	sourceNode := nodesToGraph(nodes)

	pipeline, err := NewPipelineBuilder(name, description, sourceNode).Compile()
	if err != nil {
		return nil, err
	}
	return pipeline, nil
}

// CreateGoatForwardPipeline creates a forward geocoding pipeline.
func CreateGoatForwardPipeline(name string, description string, placeCol string) (*pipeline.PipelineDescription, error) {
	step0 := NewPipelineNode(NewDenormalizeStep())
	step1 := NewPipelineNode(NewDatasetToDataframeStep())
	step2 := NewPipelineNode(NewGoatForwardStep(placeCol))
	step0.Add(step1)
	step1.Add(step2)

	pipeline, err := NewPipelineBuilder(name, description, step0).Compile()

	if err != nil {
		return nil, err
	}
	return pipeline, nil
}

// CreateGoatReversePipeline creates a forward geocoding pipeline.
func CreateGoatReversePipeline(name string, description string, lonSource string, latSource string) (*pipeline.PipelineDescription, error) {
	// insantiate the pipeline
	step0 := NewPipelineNode(NewDenormalizeStep())
	step1 := NewPipelineNode(NewDatasetToDataframeStep())
	step2 := NewPipelineNode(NewGoatReverseStep(lonSource, latSource))
	step0.Add(step1)
	step1.Add(step2)

	pipeline, err := NewPipelineBuilder(name, description, step0).Compile()
	if err != nil {
		return nil, err
	}
	return pipeline, nil
}

// CreateJoinPipeline creates a pipeline that joins two input datasets using a caller supplied column.
// Accuracy is a normalized value that controls how exact the join has to be.
func CreateJoinPipeline(name string, description string, leftJoinCol string, rightJoinCol string, accuracy float32) (*pipeline.PipelineDescription, error) {
	// compute column indices

	// instantiate the pipeline - this merges two intput streams via a single join call
	step0_0 := NewPipelineNode(NewDenormalizeStep())
	step1_0 := NewPipelineNode(NewDenormalizeStep())
	step2 := NewPipelineNode(NewJoinStep(leftJoinCol, rightJoinCol, accuracy))
	step3 := NewPipelineNode(NewDatasetToDataframeStep())
	step0_0.Add(step2)
	step1_0.Add(step2)
	step2.Add(step3)

	pipeline, err := NewPipelineBuilder(name, description, step0_0, step1_0).Compile()
	if err != nil {
		return nil, err
	}
	return pipeline, nil
}

// CreateDSBoxJoinPipeline creates a pipeline that joins two input datasets
// using caller supplied columns.
func CreateDSBoxJoinPipeline(name string, description string, leftJoinCols []string, rightJoinCols []string, accuracy float32) (*pipeline.PipelineDescription, error) {
	// compute column indices

	// instantiate the pipeline - this merges two intput streams via a single join call
	step0_0 := NewPipelineNode(NewDenormalizeStep())
	step1_0 := NewPipelineNode(NewDenormalizeStep())
	step2 := NewPipelineNode(NewDSBoxJoinStep(leftJoinCols, rightJoinCols, accuracy))
	step3 := NewPipelineNode(NewDatasetToDataframeStep())
	step0_0.Add(step2)
	step1_0.Add(step2)
	step2.Add(step3)

	pipeline, err := NewPipelineBuilder(name, description, step0_0, step1_0).Compile()
	if err != nil {
		return nil, err
	}
	return pipeline, nil
}

// CreateTimeseriesFormatterPipeline creates a time series formatter pipeline.
func CreateTimeseriesFormatterPipeline(name string, description string, mainResourceID string, fileColIndex int) (*pipeline.PipelineDescription, error) {
	step0 := NewPipelineNode(NewTimeseriesFormatterStep(mainResourceID, fileColIndex))
	step1 := NewPipelineNode(NewDatasetToDataframeStep())
	step0.Add(step1)

	pipeline, err := NewPipelineBuilder(name, description, step0).Compile()
	if err != nil {
		return nil, err
	}
	return pipeline, nil
}

func mapColumns(allFeatures []*model.Variable, selectedSet map[string]bool) map[string]int {
	colIndices := make(map[string]int)
	index := 0
	for _, f := range allFeatures {
		if selectedSet[f.Name] {
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

func nodesToGraph(nodes []*PipelineNode) *PipelineNode {
	currNode := nodes[0]
	for _, node := range nodes[1:] {
		currNode.Add(node)
		currNode = node
	}
	return nodes[0]
}
