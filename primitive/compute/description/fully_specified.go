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
	"strings"

	"github.com/pkg/errors"
	"github.com/uncharted-distil/distil-compute/model"
	"github.com/uncharted-distil/distil-compute/pipeline"
)

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

// CreateGroupingFieldComposePipeline creates a pipeline to create a grouping key field for a dataset.
func CreateGroupingFieldComposePipeline(name string, description string, colIndices []int, joinChar string, outputName string) (*pipeline.PipelineDescription, error) {
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
	return pipeline, nil
}

// CreateGoatForwardPipeline creates a forward geocoding pipeline.
func CreateGoatForwardPipeline(name string, description string, placeCol *model.Variable) (*pipeline.PipelineDescription, error) {
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
	return pipeline, nil
}

// CreateGoatReversePipeline creates a forward geocoding pipeline.
func CreateGoatReversePipeline(name string, description string, lonSource *model.Variable, latSource *model.Variable) (*pipeline.PipelineDescription, error) {
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
	return pipeline, nil
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
