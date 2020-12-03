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

// CreateImageQueryPipeline creates a pipeline that will perform image retrieval.
func CreateImageQueryPipeline(name string, description string) (*FullySpecifiedPipeline, error) {
	steps := []Step{
		NewDatasetToDataframeStep(map[string]DataRef{"inputs": &PipelineDataRef{0}}, []string{"produce"}),
		NewRemoveColumnsStep(
			map[string]DataRef{"inputs": &StepDataRef{0, "produce"}},
			[]string{"produce"},
			[]int{1, 2, 3, 4},
		),
		NewDistilColumnParserStep(
			map[string]DataRef{"inputs": &StepDataRef{1, "produce"}},
			[]string{"produce"},
			[]string{model.TA2IntegerType, model.TA2RealType, model.TA2RealVectorType},
		),
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
		NewImageRetrievalStep(map[string]DataRef{"inputs": &StepDataRef{3, "produce"}, "outputs": &StepDataRef{5, "produce"}}, []string{"produce"}),
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

// CreateMultiBandImageFeaturizationPipeline creates a pipline that will featurize multiband images.
func CreateMultiBandImageFeaturizationPipeline(name string, description string, variables []*model.Variable,
	numJobs int, batchSize int) (*FullySpecifiedPipeline, error) {

	// add semantic types to variables that are images and group ids
	var grouping *model.MultiBandImageGrouping
	variableMap := map[string]*model.Variable{}
	for _, v := range variables {
		if v.IsGrouping() && model.IsMultiBandImage(v.Grouping.GetType()) {
			grouping = v.Grouping.(*model.MultiBandImageGrouping)
		}
		variableMap[v.StorageName] = v
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

	steps = append(steps, NewRemoteSensingPretrainedStep(map[string]DataRef{"inputs": &StepDataRef{offset, "produce"}}, []string{"produce"}, batchSize))

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
func CreateTargetRankingPipeline(name string, description string, target string, features []*model.Variable) (*FullySpecifiedPipeline, error) {

	// compute index associated with column name
	targetIdx := -1
	for _, f := range features {
		if strings.EqualFold(target, f.StorageName) {
			targetIdx = f.Index
			break
		}
	}
	if targetIdx < 0 {
		return nil, errors.Errorf("can't find var '%s'", name)
	}

	// don't rank group features or any included metadata features
	selectedFeatures := []*model.Variable{}
	for _, v := range features {
		if v.Grouping == nil && v.DistilRole != model.VarDistilRoleMetadata {
			selectedFeatures = append(selectedFeatures, v)
		}
	}

	offset := 0
	steps := []Step{
		NewDenormalizeStep(map[string]DataRef{"inputs": &PipelineDataRef{0}}, []string{"produce"}),
	}
	offset++

	// ranking is dependent on user updated semantic types, so we need to make sure we apply
	// those to the original data
	updateSemanticTypeStep, err := createUpdateSemanticTypes(target, selectedFeatures, map[string]bool{}, offset)
	if err != nil {
		return nil, err
	}

	steps = append(steps, updateSemanticTypeStep...)

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
func CreateJoinPipeline(name string, description string, leftJoinCol *model.Variable, rightJoinCol *model.Variable, accuracy float32) (*FullySpecifiedPipeline, error) {
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
