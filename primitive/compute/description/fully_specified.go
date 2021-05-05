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
	"encoding/json"
	"strings"

	"github.com/pkg/errors"
	"github.com/uncharted-distil/distil-compute/model"
	"github.com/uncharted-distil/distil-compute/pipeline"
	"github.com/uncharted-distil/distil-compute/primitive/compute"
)

// FullySpecifiedPipeline wraps a fully specified pipeline along with
// the fields which can be used to determine equivalent pipelines.
type FullySpecifiedPipeline struct {
	Pipeline         *pipeline.PipelineDescription
	EquivalentValues []interface{}
}

// MarshalSteps marshals a pipeline description into a json representation.
func MarshalSteps(step *pipeline.PipelineDescription) (string, error) {
	stepJSON, err := json.Marshal(step)
	if err != nil {
		return "", errors.Wrapf(err, "unable to marshal steps")
	}

	return string(stepJSON), nil
}

// CreateImageQueryPipeline creates a pipeline that will perform image retrieval.  The cacheLocation parameter
// is passed down to the image retrieval primitive, and is used to cache dot products across query operations.
// When a new dataset is being labelled, the cache location should be updated.
func CreateImageQueryPipeline(name string, description string, cacheLocation string) (*FullySpecifiedPipeline, error) {
	steps := []Step{
		NewDatasetToDataframeStep(map[string]DataRef{"inputs": &PipelineDataRef{0}}, []string{"produce"}),
		NewDistilColumnParserStep(
			map[string]DataRef{"inputs": &StepDataRef{0, "produce"}},
			[]string{"produce"},
			[]string{model.TA2IntegerType, model.TA2RealType, model.TA2RealVectorType},
		),
		NewExtractColumnsByStructuralTypeStep(map[string]DataRef{"inputs": &StepDataRef{1, "produce"}}, []string{"produce"},
			[]string{
				"int",
				"float",         // python type
				"numpy.float32", // numpy types
				"numpy.float64",
			}),
		NewExtractColumnsBySemanticTypeStep(
			map[string]DataRef{"inputs": &StepDataRef{2, "produce"}},
			[]string{"produce"},
			[]string{"https://metadata.datadrivendiscovery.org/types/Attribute", "https://metadata.datadrivendiscovery.org/types/PrimaryMultiKey"},
		),
		NewDatasetToDataframeStep(map[string]DataRef{"inputs": &PipelineDataRef{1}}, []string{"produce"}),
		NewDistilColumnParserStep(
			map[string]DataRef{"inputs": &StepDataRef{4, "produce"}},
			[]string{"produce"},
			[]string{model.TA2IntegerType, model.TA2RealType, model.TA2RealVectorType},
		),
		NewImageRetrievalStep(map[string]DataRef{
			"inputs":  &StepDataRef{3, "produce"},
			"outputs": &StepDataRef{5, "produce"}},
			[]string{"produce"},
			cacheLocation,
		),
		NewAddSemanticTypeStep(map[string]DataRef{"inputs": &StepDataRef{6, "produce"}},
			[]string{"produce"},
			&ColumnUpdate{
				SemanticTypes: []string{"https://metadata.datadrivendiscovery.org/types/PredictedTarget", "https://metadata.datadrivendiscovery.org/types/Score"},
				Indices:       []int{1},
			}),
		NewConstructPredictionStep(map[string]DataRef{"inputs": &StepDataRef{7, "produce"}}, []string{"produce"}, &StepDataRef{2, "produce"}),
	}

	inputs := []string{"inputs.0", "inputs.1"}
	outputs := []DataRef{&StepDataRef{len(steps) - 1, "produce"}}

	pipeline, err := NewPipelineBuilder(name, description, inputs, outputs, steps).Compile()
	if err != nil {
		return nil, err
	}

	pipelineJSON, err := MarshalSteps(pipeline)
	if err != nil {
		return nil, err
	}

	fullySpecified := &FullySpecifiedPipeline{
		Pipeline:         pipeline,
		EquivalentValues: []interface{}{pipelineJSON},
	}
	return fullySpecified, nil
}

// CreateImageFeaturizationPipeline creates a pipline that will featurize images.
func CreateImageFeaturizationPipeline(name string, description string, variables []*model.Variable) (*FullySpecifiedPipeline, error) {

	// add semantic types to variables that are images
	var imageCol *model.Variable
	vars := map[string]*model.Variable{}
	for _, v := range variables {
		if model.IsImage(v.Type) {
			imageCol = v
		}
		vars[v.Key] = v
	}
	if imageCol == nil {
		return nil, errors.Errorf("no image found in dataset variables")
	}

	steps := []Step{
		NewDenormalizeStep(map[string]DataRef{"inputs": &PipelineDataRef{0}}, []string{"produce"}),
		NewDatasetToDataframeStep(map[string]DataRef{"inputs": &StepDataRef{0, "produce"}}, []string{"produce"}),
	}
	offset := 1

	addImage := &ColumnUpdate{
		SemanticTypes: []string{model.TA2ImageType},
		Indices:       []int{imageCol.Index},
	}
	steps = append(steps, NewAddSemanticTypeStep(map[string]DataRef{"inputs": &StepDataRef{offset, "produce"}}, []string{"produce"}, addImage))
	offset++

	steps = append(steps, NewDataframeImageReaderStep(map[string]DataRef{"inputs": &StepDataRef{offset, "produce"}}, []string{"produce"}, []int{vars[model.D3MIndexFieldName].Index, imageCol.Index}))
	offset++

	steps = append(steps, NewImageTransferStep(map[string]DataRef{"inputs": &StepDataRef{offset, "produce"}}, []string{"produce"}))
	offset++
	transferStep := offset

	steps = append(steps, NewRemoveColumnsStep(map[string]DataRef{"inputs": &StepDataRef{1, "produce"}}, []string{"produce"}, []int{imageCol.Index}))
	offset++

	steps = append(steps, NewHorizontalConcatStep(
		map[string]DataRef{"left": &StepDataRef{offset, "produce"}, "right": &StepDataRef{transferStep, "produce"}},
		[]string{"produce"}, false, true))

	inputs := []string{"inputs"}
	outputs := []DataRef{&StepDataRef{len(steps) - 1, "produce"}}

	pipeline, err := NewPipelineBuilder(name, description, inputs, outputs, steps).Compile()
	if err != nil {
		return nil, err
	}

	pipelineJSON, err := MarshalSteps(pipeline)
	if err != nil {
		return nil, err
	}

	fullySpecified := &FullySpecifiedPipeline{
		Pipeline:         pipeline,
		EquivalentValues: []interface{}{pipelineJSON},
	}
	return fullySpecified, nil
}

// CreateMultiBandImageFeaturizationPipeline creates a pipline that will featurize multiband images.
func CreateMultiBandImageFeaturizationPipeline(name string, description string, variables []*model.Variable,
	numJobs int, batchSize int, poolFeatures bool) (*FullySpecifiedPipeline, error) {

	// add semantic types to variables that are images and group ids
	var grouping *model.MultiBandImageGrouping
	variableMap := map[string]*model.Variable{}
	for _, v := range variables {
		if v.IsGrouping() && model.IsMultiBandImage(v.Grouping.GetType()) {
			grouping = v.Grouping.(*model.MultiBandImageGrouping)
		}
		variableMap[v.Key] = v
	}

	if grouping == nil {
		return nil, errors.Errorf("no grouping found in dataset variables")
	}
	imageCol := variableMap[grouping.ImageCol]
	groupingCol := variableMap[grouping.IDCol]

	steps := []Step{
		NewDenormalizeStep(map[string]DataRef{"inputs": &PipelineDataRef{0}}, []string{"produce"}),
		NewDatasetToDataframeStep(map[string]DataRef{"inputs": &StepDataRef{0, "produce"}}, []string{"produce"}),
	}
	offset := 1

	addImage := &ColumnUpdate{
		SemanticTypes: []string{model.TA2ImageType},
		Indices:       []int{imageCol.Index},
	}
	steps = append(steps, NewAddSemanticTypeStep(map[string]DataRef{"inputs": &StepDataRef{offset, "produce"}}, []string{"produce"}, addImage))
	offset++

	addGrouping := &ColumnUpdate{
		SemanticTypes: []string{model.TA2GroupingKeyType},
		Indices:       []int{groupingCol.Index},
	}
	steps = append(steps, NewAddSemanticTypeStep(map[string]DataRef{"inputs": &StepDataRef{offset, "produce"}}, []string{"produce"}, addGrouping))
	offset++

	steps = append(steps, NewSatelliteImageLoaderStep(map[string]DataRef{"inputs": &StepDataRef{offset, "produce"}}, []string{"produce"}, numJobs))
	offset++

	steps = append(steps, NewRemoteSensingPretrainedStep(map[string]DataRef{"inputs": &StepDataRef{offset, "produce"}}, []string{"produce"}, batchSize, poolFeatures))

	inputs := []string{"inputs"}
	outputs := []DataRef{&StepDataRef{len(steps) - 1, "produce"}}

	pipeline, err := NewPipelineBuilder(name, description, inputs, outputs, steps).Compile()
	if err != nil {
		return nil, err
	}

	pipelineJSON, err := MarshalSteps(pipeline)
	if err != nil {
		return nil, err
	}

	fullySpecified := &FullySpecifiedPipeline{
		Pipeline:         pipeline,
		EquivalentValues: []interface{}{pipelineJSON},
	}
	return fullySpecified, nil
}

// CreateSlothPipeline creates a pipeline to peform timeseries clustering on a dataset.
func CreateSlothPipeline(name string, description string, timeColumn string, valueColumn string,
	timeseriesGrouping *model.TimeseriesGrouping, timeSeriesFeatures []*model.Variable) (*FullySpecifiedPipeline, error) {

	// get the grouping columns to create the groupings in the pipeline
	groupingIndices := []int{}
	timeCols := []int{}
	if timeseriesGrouping != nil {
		groupingSet := map[string]bool{}
		for _, subID := range timeseriesGrouping.SubIDs {
			groupingSet[strings.ToLower(subID)] = true
		}
		groupingIndices = listColumns(timeSeriesFeatures, groupingSet)
		timeCols = listColumns(timeSeriesFeatures, map[string]bool{timeseriesGrouping.XCol: true})
	}

	steps := make([]Step, 0)
	steps = append(steps, NewTimeseriesFormatterStep(map[string]DataRef{"inputs": &PipelineDataRef{0}}, []string{"produce"}, compute.DefaultResourceID, -1))

	offset := 1
	updateSemanticTypes, err := createUpdateSemanticTypes("", timeSeriesFeatures, nil, offset)
	if err != nil {
		return nil, err
	}
	steps = append(steps, updateSemanticTypes...)
	offset += len(updateSemanticTypes) - 1

	steps = append(steps, NewDatasetToDataframeStep(map[string]DataRef{"inputs": &StepDataRef{offset, "produce"}}, []string{"produce"}))
	steps = append(steps, NewGroupingFieldComposeStep(
		map[string]DataRef{"inputs": &StepDataRef{offset + 1, "produce"}},
		[]string{"produce"},
		groupingIndices,
		"_",
		"__grouping_key",
	))
	offset += 2

	// add the time indicator type to the time column
	addTime := NewAddSemanticTypeStep(map[string]DataRef{"inputs": &StepDataRef{offset, "produce"}}, []string{"produce"}, &ColumnUpdate{
		SemanticTypes: []string{model.TA2TimeType},
		Indices:       timeCols,
	})
	steps = append(steps, addTime)
	offset++

	steps = append(steps, NewSlothStep(map[string]DataRef{"inputs": &StepDataRef{offset, "produce"}}, []string{"produce"}))

	inputs := []string{"inputs"}
	outputs := []DataRef{&StepDataRef{len(steps) - 1, "produce"}}

	pipeline, err := NewPipelineBuilder(name, description, inputs, outputs, steps).Compile()
	if err != nil {
		return nil, err
	}

	pipelineJSON, err := MarshalSteps(pipeline)
	if err != nil {
		return nil, err
	}

	fullySpecified := &FullySpecifiedPipeline{
		Pipeline:         pipeline,
		EquivalentValues: []interface{}{pipelineJSON},
	}
	return fullySpecified, nil
}

// CreateDukePipeline creates a pipeline to peform image featurization on a dataset.
func CreateDukePipeline(name string, description string) (*FullySpecifiedPipeline, error) {
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

	pipelineJSON, err := MarshalSteps(pipeline)
	if err != nil {
		return nil, err
	}

	fullySpecified := &FullySpecifiedPipeline{
		Pipeline:         pipeline,
		EquivalentValues: []interface{}{pipelineJSON},
	}
	return fullySpecified, nil
}

// CreateSimonPipeline creates a pipeline to run semantic type inference on a dataset's
// columns.
func CreateSimonPipeline(name string, description string) (*FullySpecifiedPipeline, error) {
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

	pipelineJSON, err := MarshalSteps(pipeline)
	if err != nil {
		return nil, err
	}

	fullySpecified := &FullySpecifiedPipeline{
		Pipeline:         pipeline,
		EquivalentValues: []interface{}{pipelineJSON},
	}
	return fullySpecified, nil
}

// CreateDataCleaningPipeline creates a pipeline to run data cleaning on a dataset.
func CreateDataCleaningPipeline(name string, description string) (*FullySpecifiedPipeline, error) {
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

	pipelineJSON, err := MarshalSteps(pipeline)
	if err != nil {
		return nil, err
	}

	fullySpecified := &FullySpecifiedPipeline{
		Pipeline:         pipeline,
		EquivalentValues: []interface{}{pipelineJSON},
	}
	return fullySpecified, nil
}

// CreateGroupingFieldComposePipeline creates a pipeline to create a grouping key field for a dataset.
func CreateGroupingFieldComposePipeline(name string, description string, colIndices []int, joinChar string, outputName string) (*FullySpecifiedPipeline, error) {
	inputs := []string{"inputs"}
	outputs := []DataRef{&StepDataRef{1, "produce"}}

	steps := []Step{
		NewDatasetToDataframeStep(map[string]DataRef{"inputs": &PipelineDataRef{0}}, []string{"produce"}),
		NewGroupingFieldComposeStep(map[string]DataRef{"inputs": &StepDataRef{0, "produce"}}, []string{"produce"}, colIndices, joinChar, outputName),
	}

	pipeline, err := NewPipelineBuilder(name, description, inputs, outputs, steps).Compile()
	if err != nil {
		return nil, err
	}

	pipelineJSON, err := MarshalSteps(pipeline)
	if err != nil {
		return nil, err
	}

	fullySpecified := &FullySpecifiedPipeline{
		Pipeline:         pipeline,
		EquivalentValues: []interface{}{pipelineJSON},
	}
	return fullySpecified, nil
}

// CreateDataFilterPipeline creates a pipeline that will filter a dataset.
func CreateDataFilterPipeline(name string, description string, variables []*model.Variable, filters []*model.FilterSet) (*FullySpecifiedPipeline, error) {
	steps := []Step{}
	offset := 0

	steps = append(steps, NewDenormalizeStep(map[string]DataRef{"inputs": &PipelineDataRef{0}}, []string{"produce"}))
	steps = append(steps, NewDataCleaningStep(nil, nil))
	steps = append(steps, NewDatasetWrapperStep(map[string]DataRef{"inputs": &StepDataRef{offset, "produce"}}, []string{"produce"}, offset+1, ""))
	offset += 3

	updateSemanticTypes, err := createUpdateSemanticTypes("", variables, nil, offset)
	if err != nil {
		return nil, err
	}
	steps = append(steps, updateSemanticTypes...)
	offset += len(updateSemanticTypes)

	// apply filters
	featureSet := map[string]int{}
	for _, v := range variables {
		featureSet[strings.ToLower(v.Key)] = v.Index
	}
	filterData, err := filterBySet(filters, featureSet, offset)
	if err != nil {
		return nil, err
	}
	steps = append(steps, filterData...)
	offset += len(filterData)

	// drop excluded distil roles columns since we do not want them in the final output
	colsToDrop := []int{}
	for _, v := range variables {
		if model.ExcludedDistilRoles[v.DistilRole] {
			colsToDrop = append(colsToDrop, v.Index)
		}
	}
	featureSelect := NewRemoveColumnsStep(nil, nil, colsToDrop)
	wrapperRemove := NewDatasetWrapperStep(map[string]DataRef{"inputs": &StepDataRef{offset - 1, "produce"}}, []string{"produce"}, offset, "")
	steps = append(steps, featureSelect, wrapperRemove)
	offset += 2

	steps = append(steps, NewDatasetToDataframeStep(map[string]DataRef{"inputs": &StepDataRef{offset - 1, "produce"}}, []string{"produce"}))

	inputs := []string{"inputs"}
	outputs := []DataRef{&StepDataRef{offset, "produce"}}

	pipeline, err := NewPipelineBuilder(name, description, inputs, outputs, steps).Compile()
	if err != nil {
		return nil, err
	}

	pipelineJSON, err := MarshalSteps(pipeline)
	if err != nil {
		return nil, err
	}

	fullySpecified := &FullySpecifiedPipeline{
		Pipeline:         pipeline,
		EquivalentValues: []interface{}{pipelineJSON},
	}
	return fullySpecified, nil
}

// CreatePCAFeaturesPipeline creates a pipeline to run feature ranking on an input dataset.
func CreatePCAFeaturesPipeline(name string, description string) (*FullySpecifiedPipeline, error) {
	inputs := []string{"inputs"}
	outputs := []DataRef{&StepDataRef{3, "produce_metafeatures"}}

	steps := []Step{
		NewDatasetToDataframeStep(map[string]DataRef{"inputs": &PipelineDataRef{0}}, []string{"produce"}),
		NewProfilerStep(map[string]DataRef{"inputs": &StepDataRef{0, "produce"}}, []string{"produce"}),
		NewColumnParserStep(map[string]DataRef{"inputs": &StepDataRef{1, "produce"}}, []string{"produce"}, []string{}),
		NewPCAFeaturesStep(map[string]DataRef{"inputs": &StepDataRef{0, "produce"}}, []string{"produce_metafeatures"}),
	}

	pipeline, err := NewPipelineBuilder(name, description, inputs, outputs, steps).Compile()
	if err != nil {
		return nil, err
	}

	pipelineJSON, err := MarshalSteps(pipeline)
	if err != nil {
		return nil, err
	}

	fullySpecified := &FullySpecifiedPipeline{
		Pipeline:         pipeline,
		EquivalentValues: []interface{}{pipelineJSON},
	}
	return fullySpecified, nil
}

// CreateDenormalizePipeline creates a pipeline to run the denormalize primitive on an input dataset.
func CreateDenormalizePipeline(name string, description string) (*FullySpecifiedPipeline, error) {
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

	pipelineJSON, err := MarshalSteps(pipeline)
	if err != nil {
		return nil, err
	}

	fullySpecified := &FullySpecifiedPipeline{
		Pipeline:         pipeline,
		EquivalentValues: []interface{}{pipelineJSON},
	}
	return fullySpecified, nil
}

// CreateTargetRankingPipeline creates a pipeline to run feature ranking on an input dataset.
func CreateTargetRankingPipeline(name string, description string, target *model.Variable,
	features []*model.Variable, selectedFeatures map[string]bool) (*FullySpecifiedPipeline, error) {

	// ignore group and metadata variables - they are system-only variables that aren't part of the
	// dataset that is passed into the pipeline
	datasetFeatures := []*model.Variable{}
	for _, v := range features {
		if v.Grouping == nil && !model.ExcludedDistilRoles[v.DistilRole] {
			datasetFeatures = append(datasetFeatures, v)
		}
	}

	// recompute indices based on selected data subset
	columnIndices := mapColumns(datasetFeatures, selectedFeatures)

	offset := 0
	steps := []Step{
		NewDenormalizeStep(map[string]DataRef{"inputs": &PipelineDataRef{0}}, []string{"produce"}),
	}
	offset++

	// ranking is dependent on user updated semantic types, so we need to make sure we apply
	// those to the original data
	updateSemanticTypeStep, err := createUpdateSemanticTypes(target.Key, datasetFeatures, selectedFeatures, offset)
	if err != nil {
		return nil, err
	}
	steps = append(steps, updateSemanticTypeStep...)
	offset += len(updateSemanticTypeStep)

	// Remove types we don't want to run ranking on
	removeFeatures := createRemoveFeatures(datasetFeatures, selectedFeatures, offset)
	steps = append(steps, removeFeatures...)
	offset += len(removeFeatures)

	// Extract from dataset, convert column types
	steps = append(steps,
		NewDatasetToDataframeStep(map[string]DataRef{"inputs": &StepDataRef{offset - 1, "produce"}}, []string{"produce"}),
		NewColumnParserStep(
			map[string]DataRef{"inputs": &StepDataRef{offset, "produce"}},
			[]string{"produce"},
			// inlcude categorical because the parser will hash all the values into ints, which the MIRanking primitive can handle
			[]string{model.TA2IntegerType, model.TA2BooleanType, model.TA2RealType, model.TA2CategoricalType},
		),
	)
	offset += 2
	steps = append(steps, NewEnrichDatesStep(map[string]DataRef{"inputs": &StepDataRef{offset - 1, "produce"}}, []string{"produce"}, true))
	offset = offset + 1
	// Apply target ranking
	targetIdx := columnIndices[strings.ToLower(target.Key)]
	steps = append(steps, NewTargetRankingStep(map[string]DataRef{"inputs": &StepDataRef{offset - 1, "produce"}}, []string{"produce"}, targetIdx))
	offset++

	inputs := []string{"inputs"}
	outputs := []DataRef{&StepDataRef{offset - 1, "produce"}}

	pipeline, err := NewPipelineBuilder(name, description, inputs, outputs, steps).Compile()
	if err != nil {
		return nil, err
	}

	fullySpecified := &FullySpecifiedPipeline{
		Pipeline:         pipeline,
		EquivalentValues: []interface{}{name, target},
	}
	return fullySpecified, nil
}

// CreateGoatForwardPipeline creates a forward geocoding pipeline.
func CreateGoatForwardPipeline(name string, description string, placeCol *model.Variable) (*FullySpecifiedPipeline, error) {
	inputs := []string{"inputs"}
	outputs := []DataRef{&StepDataRef{2, "produce"}}

	steps := []Step{
		NewDenormalizeStep(map[string]DataRef{"inputs": &PipelineDataRef{0}}, []string{"produce"}),
		NewDatasetToDataframeStep(map[string]DataRef{"inputs": &StepDataRef{0, "produce"}}, []string{"produce"}),
		NewGoatForwardStep(map[string]DataRef{"inputs": &StepDataRef{1, "produce"}}, []string{"produce"}, placeCol.Index),
	}
	pipeline, err := NewPipelineBuilder(name, description, inputs, outputs, steps).Compile()

	if err != nil {
		return nil, err
	}

	pipelineJSON, err := MarshalSteps(pipeline)
	if err != nil {
		return nil, err
	}

	fullySpecified := &FullySpecifiedPipeline{
		Pipeline:         pipeline,
		EquivalentValues: []interface{}{pipelineJSON},
	}
	return fullySpecified, nil
}

// CreateGoatReversePipeline creates a forward geocoding pipeline.
func CreateGoatReversePipeline(name string, description string, lonSource *model.Variable, latSource *model.Variable) (*FullySpecifiedPipeline, error) {
	inputs := []string{"inputs"}
	outputs := []DataRef{&StepDataRef{2, "produce"}}

	steps := []Step{
		NewDenormalizeStep(map[string]DataRef{"inputs": &PipelineDataRef{0}}, []string{"produce"}),
		NewDatasetToDataframeStep(map[string]DataRef{"inputs": &StepDataRef{0, "produce"}}, []string{"produce"}),
		NewGoatReverseStep(map[string]DataRef{"inputs": &StepDataRef{1, "produce"}}, []string{"produce"}, lonSource.Index, latSource.Index),
	}

	pipeline, err := NewPipelineBuilder(name, description, inputs, outputs, steps).Compile()
	if err != nil {
		return nil, err
	}

	pipelineJSON, err := MarshalSteps(pipeline)
	if err != nil {
		return nil, err
	}

	fullySpecified := &FullySpecifiedPipeline{
		Pipeline:         pipeline,
		EquivalentValues: []interface{}{pipelineJSON},
	}
	return fullySpecified, nil
}

// CreateJoinPipeline creates a pipeline that joins two input datasets using a caller supplied column.
// Accuracy is a normalized value that controls how exact the join has to be.
func CreateJoinPipeline(name string, description string, leftJoinCols []*model.Variable, rightJoinCols []*model.Variable, accuracy float32) (*FullySpecifiedPipeline, error) {
	steps := make([]Step, 0)
	steps = append(steps, NewDenormalizeStep(map[string]DataRef{"inputs": &PipelineDataRef{0}}, []string{"produce"}))
	offset := 1

	// update semantic types as needed and parse vector types
	stepsRetype, err := createUpdateSemanticTypes("", leftJoinCols, nil, offset)
	if err != nil {
		return nil, err
	}
	steps = append(steps, stepsRetype...)
	offset += len(stepsRetype)
	steps = append(steps, NewDistilColumnParserStep(nil, nil, []string{model.TA2IntegerType, model.TA2BooleanType, model.TA2RealType, model.TA2RealVectorType}))
	steps = append(steps, NewDatasetWrapperStep(map[string]DataRef{"inputs": &StepDataRef{offset - 1, "produce"}}, []string{"produce"}, offset, ""))
	offsetLeft := offset + 1
	offset = offset + 2

	steps = append(steps, NewDenormalizeStep(map[string]DataRef{"inputs": &PipelineDataRef{1}}, []string{"produce"}))
	offset = offset + 1
	stepsRetype, err = createUpdateSemanticTypes("", rightJoinCols, nil, offset)
	if err != nil {
		return nil, err
	}
	steps = append(steps, stepsRetype...)
	offset += len(stepsRetype)
	steps = append(steps, NewDistilColumnParserStep(nil, nil, []string{model.TA2IntegerType, model.TA2BooleanType, model.TA2RealType, model.TA2RealVectorType}))
	steps = append(steps, NewDatasetWrapperStep(map[string]DataRef{"inputs": &StepDataRef{offset - 1, "produce"}}, []string{"produce"}, offset, ""))
	offsetRight := offset + 1
	offset = offset + 2

	// merge two intput streams via a single join call
	leftColNames := make([]string, len(leftJoinCols))
	for i := range leftJoinCols {
		leftColNames[i] = leftJoinCols[i].HeaderName
	}
	rightColNames := make([]string, len(rightJoinCols))
	for i := range rightJoinCols {
		rightColNames[i] = rightJoinCols[i].HeaderName
	}
	steps = append(steps, NewJoinStep(
		map[string]DataRef{"left": &StepDataRef{offsetLeft, "produce"}, "right": &StepDataRef{offsetRight, "produce"}},
		[]string{"produce"}, leftColNames, rightColNames, accuracy,
	))
	steps = append(steps, NewDatasetToDataframeStep(map[string]DataRef{"inputs": &StepDataRef{offset, "produce"}}, []string{"produce"}))

	// compute column indices
	inputs := []string{"left", "right"}
	outputs := []DataRef{&StepDataRef{len(steps) - 1, "produce"}}

	pipeline, err := NewPipelineBuilder(name, description, inputs, outputs, steps).Compile()
	if err != nil {
		return nil, err
	}

	pipelineJSON, err := MarshalSteps(pipeline)
	if err != nil {
		return nil, err
	}

	fullySpecified := &FullySpecifiedPipeline{
		Pipeline:         pipeline,
		EquivalentValues: []interface{}{pipelineJSON},
	}
	return fullySpecified, nil
}

// CreateDSBoxJoinPipeline creates a pipeline that joins two input datasets
// using caller supplied columns.
func CreateDSBoxJoinPipeline(name string, description string, leftJoinCols []string, rightJoinCols []string, accuracy float32) (*FullySpecifiedPipeline, error) {
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

	pipelineJSON, err := MarshalSteps(pipeline)
	if err != nil {
		return nil, err
	}

	fullySpecified := &FullySpecifiedPipeline{
		Pipeline:         pipeline,
		EquivalentValues: []interface{}{pipelineJSON},
	}
	return fullySpecified, nil
}

// CreateTimeseriesFormatterPipeline creates a time series formatter pipeline.
func CreateTimeseriesFormatterPipeline(name string, description string, resource string) (*FullySpecifiedPipeline, error) {

	inputs := []string{"inputs"}
	outputs := []DataRef{&StepDataRef{1, "produce"}}

	steps := []Step{
		NewTimeseriesFormatterStep(map[string]DataRef{"inputs": &PipelineDataRef{0}}, []string{"produce"}, resource, -1),
		NewDatasetToDataframeStep(map[string]DataRef{"inputs": &StepDataRef{0, "produce"}}, []string{"produce"}),
	}

	pipeline, err := NewPipelineBuilder(name, description, inputs, outputs, steps).Compile()
	if err != nil {
		return nil, err
	}

	pipelineJSON, err := MarshalSteps(pipeline)
	if err != nil {
		return nil, err
	}

	fullySpecified := &FullySpecifiedPipeline{
		Pipeline:         pipeline,
		EquivalentValues: []interface{}{pipelineJSON},
	}
	return fullySpecified, nil
}

// CreateDatamartDownloadPipeline creates a pipeline to download data from a datamart.
func CreateDatamartDownloadPipeline(name string, description string, searchResult string, systemIdentifier string) (*FullySpecifiedPipeline, error) {
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

	pipelineJSON, err := MarshalSteps(pipeline)
	if err != nil {
		return nil, err
	}

	fullySpecified := &FullySpecifiedPipeline{
		Pipeline:         pipeline,
		EquivalentValues: []interface{}{pipelineJSON},
	}
	return fullySpecified, nil
}

// CreateImageOutlierDetectionPipeline makes a pipeline for
// outlier detection with remote sensing data
func CreateImageOutlierDetectionPipeline(name string, description string, imageVariables []*model.Variable) (*FullySpecifiedPipeline, error) {

	cols := make([]int, len(imageVariables))
	for i, v := range imageVariables {
		cols[i] = v.Index
	}

	steps := []Step{
		NewDenormalizeStep(map[string]DataRef{"inputs": &PipelineDataRef{0}}, []string{"produce"}),
		NewDatasetToDataframeStep(map[string]DataRef{"inputs": &StepDataRef{0, "produce"}}, []string{"produce"}),
		NewDataframeImageReaderStep(map[string]DataRef{"inputs": &StepDataRef{1, "produce"}}, []string{"produce"}, cols),
		NewImageTransferStep(map[string]DataRef{"inputs": &StepDataRef{2, "produce"}}, []string{"produce"}),
		NewIsolationForestStep(map[string]DataRef{"inputs": &StepDataRef{3, "produce"}}, []string{"produce"}),
		NewConstructPredictionStep(map[string]DataRef{"inputs": &StepDataRef{4, "produce"}}, []string{"produce"}, &StepDataRef{1, "produce"}),
	}

	inputs := []string{"inputs"}
	outputs := []DataRef{&StepDataRef{len(steps) - 1, "produce"}}

	pipeline, err := NewPipelineBuilder(name, description, inputs, outputs, steps).Compile()
	if err != nil {
		return nil, err
	}

	pipelineJSON, err := MarshalSteps(pipeline)
	if err != nil {
		return nil, err
	}

	fullySpecified := &FullySpecifiedPipeline{
		Pipeline:         pipeline,
		EquivalentValues: []interface{}{pipelineJSON},
	}
	return fullySpecified, nil
}

// CreateMultiBandImageOutlierDetectionPipeline does outlier detection for multiband images
// for both prefeaturised and featurised
func CreateMultiBandImageOutlierDetectionPipeline(name string, description string, imageVariables []*model.Variable,
	prefeaturised bool, pooled bool, grouping *model.MultiBandImageGrouping, batchSize int, numJobs int) (*FullySpecifiedPipeline, error) {

	var steps []Step
	if prefeaturised {
		steps = []Step{
			NewDatasetToDataframeStep(map[string]DataRef{"inputs": &PipelineDataRef{0}}, []string{"produce"}),
			NewDistilColumnParserStep(map[string]DataRef{"inputs": &StepDataRef{0, "produce"}}, []string{"produce"},
				[]string{
					model.TA2RealType,
					model.TA2RealVectorType,
				}),
			NewExtractColumnsBySemanticTypeStep(
				map[string]DataRef{"inputs": &StepDataRef{1, "produce"}},
				[]string{"produce"},
				[]string{"https://metadata.datadrivendiscovery.org/types/Attribute"},
			),
		}
		offset := 2
		if !pooled {
			steps = append(steps, NewPrefeaturisedPoolingStep(map[string]DataRef{"inputs": &StepDataRef{offset, "produce"}}, []string{"produce"}))
			offset++
		}

		moreSteps := []Step{
			NewListEncoderStep(map[string]DataRef{"inputs": &StepDataRef{offset, "produce"}}, []string{"produce"}),
			NewEnrichDatesStep(map[string]DataRef{"inputs": &StepDataRef{offset + 1, "produce"}}, []string{"produce"}, false),
			NewExtractColumnsByStructuralTypeStep(map[string]DataRef{"inputs": &StepDataRef{offset + 2, "produce"}}, []string{"produce"},
				[]string{
					"float",         // python type
					"numpy.float32", // numpy types
					"numpy.float64",
				}),
			NewIsolationForestStep(map[string]DataRef{"inputs": &StepDataRef{offset + 3, "produce"}}, []string{"produce"}),
			NewConstructPredictionStep(map[string]DataRef{"inputs": &StepDataRef{offset + 4, "produce"}}, []string{"produce"}, &StepDataRef{1, "produce"}),
		}
		steps = append(steps, moreSteps...)
	} else {
		var imageVar *model.Variable
		var groupVar *model.Variable
		for _, v := range imageVariables {
			if v.Key == grouping.ImageCol {
				imageVar = v
			} else if v.Key == grouping.IDCol {
				groupVar = v
			}
		}
		if imageVar == nil {
			return nil, errors.Errorf("image var with name '%s' not found in supplied variables", grouping.ImageCol)
		}
		if groupVar == nil {
			return nil, errors.Errorf("grouping var with name '%s' not found in supplied variables", grouping.IDCol)
		}
		addGroup := &ColumnUpdate{
			SemanticTypes: []string{"https://metadata.datadrivendiscovery.org/types/GroupingKey"},
			Indices:       []int{groupVar.Index},
		}
		steps = []Step{
			NewDenormalizeStep(map[string]DataRef{"inputs": &PipelineDataRef{0}}, []string{"produce"}),
			NewAddSemanticTypeStep(nil, nil, addGroup),
			NewDatasetWrapperStep(map[string]DataRef{"inputs": &StepDataRef{0, "produce"}}, []string{"produce"}, 1, ""),
			NewDatasetToDataframeStep(map[string]DataRef{"inputs": &StepDataRef{2, "produce"}}, []string{"produce"}),
			NewSatelliteImageLoaderStep(map[string]DataRef{"inputs": &StepDataRef{3, "produce"}}, []string{"produce"}, numJobs),
			NewDistilColumnParserStep(map[string]DataRef{"inputs": &StepDataRef{4, "produce"}}, []string{"produce"},
				[]string{model.TA2IntegerType, model.TA2RealType, model.TA2RealVectorType}),
			NewRemoteSensingPretrainedStep(map[string]DataRef{"inputs": &StepDataRef{5, "produce"}}, []string{"produce"}, batchSize, true),
			NewExtractColumnsBySemanticTypeStep(map[string]DataRef{"inputs": &StepDataRef{6, "produce"}}, []string{"produce"},
				[]string{"http://schema.org/Float"}),
			NewIsolationForestStep(map[string]DataRef{"inputs": &StepDataRef{7, "produce"}}, []string{"produce"}),
			NewConstructPredictionStep(map[string]DataRef{"inputs": &StepDataRef{8, "produce"}}, []string{"produce"}, &StepDataRef{4, "produce"}),
		}
	}

	inputs := []string{"inputs"}
	outputs := []DataRef{&StepDataRef{len(steps) - 1, "produce"}}

	pipeline, err := NewPipelineBuilder(name, description, inputs, outputs, steps).Compile()
	if err != nil {
		return nil, err
	}

	pipelineJSON, err := MarshalSteps(pipeline)
	if err != nil {
		return nil, err
	}

	fullySpecified := &FullySpecifiedPipeline{
		Pipeline:         pipeline,
		EquivalentValues: []interface{}{pipelineJSON},
	}
	return fullySpecified, nil
}

// CreateTabularOutlierDetectionPipeline makes a pipeline for
// outlier detection
func CreateTabularOutlierDetectionPipeline(name string, description string, datasetDescription *UserDatasetDescription,
	augmentations []*UserDatasetAugmentation) (*FullySpecifiedPipeline, error) {

	steps, err := generatePrependSteps(datasetDescription, augmentations)
	if err != nil {
		return nil, err
	}
	offset := len(steps) - 1

	add := &ColumnUpdate{
		SemanticTypes: []string{"https://metadata.datadrivendiscovery.org/types/TrueTarget"},
		Indices:       []int{datasetDescription.TargetFeature.Index},
	}
	steps = append(steps, NewAddSemanticTypeStep(nil, nil, add))
	offset = offset + 1
	steps = append(steps, NewDatasetWrapperStep(map[string]DataRef{"inputs": &StepDataRef{offset - 1, "produce"}}, []string{"produce"}, offset, ""))
	offset = offset + 1
	steps = append(steps, NewDatasetToDataframeStep(map[string]DataRef{"inputs": &StepDataRef{offset, "produce"}}, []string{"produce"}))
	offset = offset + 1
	steps = append(steps, NewColumnParserStep(
		map[string]DataRef{"inputs": &StepDataRef{offset, "produce"}},
		[]string{"produce"},
		[]string{model.TA2IntegerType, model.TA2RealType, model.TA2RealVectorType},
	))
	offset = offset + 1
	parseStep := offset
	steps = append(steps, NewExtractColumnsBySemanticTypeStep(
		map[string]DataRef{"inputs": &StepDataRef{parseStep, "produce"}},
		[]string{"produce"},
		[]string{"https://metadata.datadrivendiscovery.org/types/Attribute"},
	))
	offset = offset + 1
	attributeStep := offset
	steps = append(steps, NewExtractColumnsBySemanticTypeStep(map[string]DataRef{"inputs": &StepDataRef{parseStep, "produce"}}, []string{"produce"}, []string{"https://metadata.datadrivendiscovery.org/types/Target", "https://metadata.datadrivendiscovery.org/types/TrueTarget"}))
	offset = offset + 1
	targetStep := offset
	steps = append(steps, NewEnrichDatesStep(map[string]DataRef{"inputs": &StepDataRef{attributeStep, "produce"}}, []string{"produce"}, false))
	offset = offset + 1
	steps = append(steps, NewListEncoderStep(map[string]DataRef{"inputs": &StepDataRef{offset, "produce"}}, []string{"produce"}))
	offset = offset + 1
	steps = append(steps, NewReplaceSingletonStep(map[string]DataRef{"inputs": &StepDataRef{offset, "produce"}}, []string{"produce"}))
	offset = offset + 1
	steps = append(steps, NewCategoricalImputerStep(map[string]DataRef{"inputs": &StepDataRef{offset, "produce"}}, []string{"produce"}))
	offset = offset + 1
	steps = append(steps, NewTextEncoderStep(map[string]DataRef{"inputs": &StepDataRef{offset, "produce"}, "outputs": &StepDataRef{targetStep, "produce"}}, []string{"produce"}))
	offset = offset + 1
	steps = append(steps, NewOneHotEncoderStep(map[string]DataRef{"inputs": &StepDataRef{offset, "produce"}}, []string{"produce"}))
	offset = offset + 1
	steps = append(steps, NewBinaryEncoderStep(map[string]DataRef{"inputs": &StepDataRef{offset, "produce"}}, []string{"produce"}))
	offset = offset + 1
	steps = append(steps, NewSKImputerStep(map[string]DataRef{"inputs": &StepDataRef{offset, "produce"}}, []string{"produce"}))
	offset = offset + 1
	steps = append(steps, NewSKMissingIndicatorStep(map[string]DataRef{"inputs": &StepDataRef{offset, "produce"}}, []string{"produce"}))
	offset = offset + 1
	steps = append(steps, NewIsolationForestStep(map[string]DataRef{"inputs": &StepDataRef{offset, "produce"}}, []string{"produce"}))
	offset = offset + 1
	steps = append(steps, NewConstructPredictionStep(map[string]DataRef{"inputs": &StepDataRef{offset, "produce"}}, []string{"produce"}, &StepDataRef{parseStep, "produce"}))

	inputs := []string{"inputs"}
	outputs := []DataRef{&StepDataRef{len(steps) - 1, "produce"}}

	pipeline, err := NewPipelineBuilder(name, description, inputs, outputs, steps).Compile()
	if err != nil {
		return nil, err
	}

	pipelineJSON, err := MarshalSteps(pipeline)
	if err != nil {
		return nil, err
	}

	fullySpecified := &FullySpecifiedPipeline{
		Pipeline:         pipeline,
		EquivalentValues: []interface{}{pipelineJSON},
	}
	return fullySpecified, nil
}

// CreateDatamartAugmentPipeline creates a pipeline to augment data with datamart data.
func CreateDatamartAugmentPipeline(name string, description string, searchResult string, systemIdentifier string) (*FullySpecifiedPipeline, error) {
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

	pipelineJSON, err := MarshalSteps(pipeline)
	if err != nil {
		return nil, err
	}

	fullySpecified := &FullySpecifiedPipeline{
		Pipeline:         pipeline,
		EquivalentValues: []interface{}{pipelineJSON},
	}
	return fullySpecified, nil
}
