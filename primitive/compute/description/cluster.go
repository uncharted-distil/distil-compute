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
	"github.com/pkg/errors"
	"github.com/uncharted-distil/distil-compute/model"
	"github.com/uncharted-distil/distil-compute/pipeline"
)

// CreateGeneralClusteringPipeline creates a pipeline that will cluster tabular data.
func CreateGeneralClusteringPipeline(name string, description string, datasetDescription *UserDatasetDescription,
	augmentations []*UserDatasetAugmentation, useKMeans bool) (*FullySpecifiedPipeline, error) {

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
		[]string{model.TA2IntegerType, model.TA2RealVectorType, model.TA2RealType},
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

	if useKMeans {
		steps = append(steps, NewKMeansClusteringStep(map[string]DataRef{"inputs": &StepDataRef{offset, "produce"}}, []string{"produce"}))
		offset = offset + 1
	} else {
		steps = append(steps, NewHDBScanStep(map[string]DataRef{"inputs": &StepDataRef{offset, "produce"}}, []string{"produce"}))
		offset = offset + 1

		steps = append(steps, NewExtractColumnsStep(map[string]DataRef{"inputs": &StepDataRef{offset, "produce"}}, []string{"produce"}, []int{-1}))
		offset = offset + 1
	}

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

// CreateImageClusteringPipeline creates a fully specified pipeline that will
// cluster images together, returning a column with the resulting cluster.
func CreateImageClusteringPipeline(name string, description string, imageVariables []*model.Variable, useKMeans bool) (*FullySpecifiedPipeline, error) {

	cols := make([]int, len(imageVariables))
	for i, v := range imageVariables {
		cols[i] = v.Index
	}

	var steps []Step
	if useKMeans {
		steps = []Step{
			NewDenormalizeStep(map[string]DataRef{"inputs": &PipelineDataRef{0}}, []string{"produce"}),
			NewDatasetToDataframeStep(map[string]DataRef{"inputs": &StepDataRef{0, "produce"}}, []string{"produce"}),
			NewDataframeImageReaderStep(map[string]DataRef{"inputs": &StepDataRef{1, "produce"}}, []string{"produce"}, cols),
			NewImageTransferStep(map[string]DataRef{"inputs": &StepDataRef{2, "produce"}}, []string{"produce"}),
			NewSKImputerStep(map[string]DataRef{"inputs": &StepDataRef{3, "produce"}}, []string{"produce"}),
			NewSKMissingIndicatorStep(map[string]DataRef{"inputs": &StepDataRef{4, "produce"}}, []string{"produce"}),
			NewKMeansClusteringStep(map[string]DataRef{"inputs": &StepDataRef{5, "produce"}}, []string{"produce"}),
			NewConstructPredictionStep(map[string]DataRef{"inputs": &StepDataRef{6, "produce"}}, []string{"produce"}, &StepDataRef{1, "produce"}),
		}
	} else {
		add := &ColumnUpdate{
			SemanticTypes: []string{"https://metadata.datadrivendiscovery.org/types/TrueTarget"},
			Indices:       cols,
		}

		steps = []Step{
			NewDenormalizeStep(map[string]DataRef{"inputs": &PipelineDataRef{0}}, []string{"produce"}),
			NewAddSemanticTypeStep(nil, nil, add),
			NewDatasetWrapperStep(map[string]DataRef{"inputs": &StepDataRef{0, "produce"}}, []string{"produce"}, 1, ""),
			NewDatasetToDataframeStep(map[string]DataRef{"inputs": &StepDataRef{2, "produce"}}, []string{"produce"}),
			NewDataframeImageReaderStep(map[string]DataRef{"inputs": &StepDataRef{3, "produce"}}, []string{"produce"}, cols),
			NewImageTransferStep(map[string]DataRef{"inputs": &StepDataRef{4, "produce"}}, []string{"produce"}),
			NewSKImputerStep(map[string]DataRef{"inputs": &StepDataRef{5, "produce"}}, []string{"produce"}),
			NewSKMissingIndicatorStep(map[string]DataRef{"inputs": &StepDataRef{6, "produce"}}, []string{"produce"}),
			NewHDBScanStep(map[string]DataRef{"inputs": &StepDataRef{7, "produce"}}, []string{"produce"}),
			NewExtractColumnsStep(map[string]DataRef{"inputs": &StepDataRef{8, "produce"}}, []string{"produce"}, []int{-1}),
			NewConstructPredictionStep(map[string]DataRef{"inputs": &StepDataRef{9, "produce"}}, []string{"produce"}, &StepDataRef{3, "produce"}),
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

// CreateMultiBandImageClusteringPipeline creates a fully specified pipeline that will
// cluster multiband images together, returning a column with the resulting cluster.
func CreateMultiBandImageClusteringPipeline(name string, description string,
	grouping *model.MultiBandImageGrouping, variables []*model.Variable, useKMeans bool,
	batchSize int, numJobs int) (*FullySpecifiedPipeline, error) {

	var imageVar *model.Variable
	var groupVar *model.Variable
	for _, v := range variables {
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

	var steps []Step
	if useKMeans {
		steps = []Step{
			NewDenormalizeStep(map[string]DataRef{"inputs": &PipelineDataRef{0}}, []string{"produce"}),
			NewAddSemanticTypeStep(nil, nil, addGroup),
			NewDatasetWrapperStep(map[string]DataRef{"inputs": &StepDataRef{0, "produce"}}, []string{"produce"}, 1, ""),
			NewDatasetToDataframeStep(map[string]DataRef{"inputs": &StepDataRef{2, "produce"}}, []string{"produce"}),
			NewSatelliteImageLoaderStep(map[string]DataRef{"inputs": &StepDataRef{3, "produce"}}, []string{"produce"}, numJobs),
			NewColumnParserStep(map[string]DataRef{"inputs": &StepDataRef{4, "produce"}}, []string{"produce"},
				[]string{model.TA2BooleanType, model.TA2IntegerType, model.TA2RealType, model.TA2RealVectorType}),
			NewRemoteSensingPretrainedStep(map[string]DataRef{"inputs": &StepDataRef{5, "produce"}}, []string{"produce"}, batchSize, true),
			NewSKImputerStep(map[string]DataRef{"inputs": &StepDataRef{6, "produce"}}, []string{"produce"}),
			NewSKMissingIndicatorStep(map[string]DataRef{"inputs": &StepDataRef{7, "produce"}}, []string{"produce"}),
			NewKMeansClusteringStep(map[string]DataRef{"inputs": &StepDataRef{8, "produce"}}, []string{"produce"}),
			NewConstructPredictionStep(map[string]DataRef{"inputs": &StepDataRef{9, "produce"}}, []string{"produce"}, &StepDataRef{4, "produce"}),
		}
	} else {
		addImage := &ColumnUpdate{
			SemanticTypes: []string{"https://metadata.datadrivendiscovery.org/types/TrueTarget"},
			Indices:       []int{imageVar.Index},
		}

		steps = []Step{
			NewDenormalizeStep(map[string]DataRef{"inputs": &PipelineDataRef{0}}, []string{"produce"}),
			NewAddSemanticTypeStep(nil, nil, addGroup),
			NewDatasetWrapperStep(map[string]DataRef{"inputs": &StepDataRef{0, "produce"}}, []string{"produce"}, 1, ""),
			NewDatasetToDataframeStep(map[string]DataRef{"inputs": &StepDataRef{2, "produce"}}, []string{"produce"}),
			NewSatelliteImageLoaderStep(map[string]DataRef{"inputs": &StepDataRef{3, "produce"}}, []string{"produce"}, numJobs),
			NewAddSemanticTypeStep(map[string]DataRef{"inputs": &StepDataRef{4, "produce"}}, []string{"produce"}, addImage),
			NewColumnParserStep(map[string]DataRef{"inputs": &StepDataRef{5, "produce"}}, []string{"produce"},
				[]string{model.TA2BooleanType, model.TA2IntegerType, model.TA2RealType, model.TA2RealVectorType}),
			NewRemoteSensingPretrainedStep(map[string]DataRef{"inputs": &StepDataRef{6, "produce"}}, []string{"produce"}, batchSize, true),
			NewSKImputerStep(map[string]DataRef{"inputs": &StepDataRef{7, "produce"}}, []string{"produce"}),
			NewSKMissingIndicatorStep(map[string]DataRef{"inputs": &StepDataRef{8, "produce"}}, []string{"produce"}),
			NewHDBScanStep(map[string]DataRef{"inputs": &StepDataRef{9, "produce"}}, []string{"produce"}),
			NewExtractColumnsStep(map[string]DataRef{"inputs": &StepDataRef{10, "produce"}}, []string{"produce"}, []int{-1}),
			NewConstructPredictionStep(map[string]DataRef{"inputs": &StepDataRef{11, "produce"}}, []string{"produce"}, &StepDataRef{5, "produce"}),
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

// CreatePreFeaturizedMultiBandImageClusteringPipeline creates a fully specified pipeline that will
// cluster multiband images together, returning a column with the resulting cluster.
func CreatePreFeaturizedMultiBandImageClusteringPipeline(name string, description string, variables []*model.Variable, useKMeans bool, poolFeatures bool) (*FullySpecifiedPipeline, error) {
	var steps []Step
	if useKMeans {
		if poolFeatures {
			steps = []Step{
				NewDatasetToDataframeStep(map[string]DataRef{"inputs": &PipelineDataRef{0}}, []string{"produce"}),
				NewDistilColumnParserStep(map[string]DataRef{"inputs": &StepDataRef{0, "produce"}}, []string{"produce"}, []string{model.TA2RealType}),
				NewExtractColumnsByStructuralTypeStep(map[string]DataRef{"inputs": &StepDataRef{1, "produce"}}, []string{"produce"},
					[]string{
						"float",         // python type
						"numpy.float32", // numpy types
						"numpy.float64",
					}),
				NewSKImputerStep(map[string]DataRef{"inputs": &StepDataRef{2, "produce"}}, []string{"produce"}),
				NewSKMissingIndicatorStep(map[string]DataRef{"inputs": &StepDataRef{3, "produce"}}, []string{"produce"}),
				NewKMeansClusteringStep(map[string]DataRef{"inputs": &StepDataRef{4, "produce"}}, []string{"produce"}),
				NewConstructPredictionStep(map[string]DataRef{"inputs": &StepDataRef{5, "produce"}}, []string{"produce"}, &StepDataRef{1, "produce"}),
			}
		} else {
			steps = []Step{
				NewDatasetToDataframeStep(map[string]DataRef{"inputs": &PipelineDataRef{0}}, []string{"produce"}),
				NewDistilColumnParserStep(map[string]DataRef{"inputs": &StepDataRef{0, "produce"}}, []string{"produce"}, []string{model.TA2RealType}),
				NewPrefeaturisedPoolingStep(map[string]DataRef{"inputs": &StepDataRef{1, "produce"}}, []string{"produce"}),
				NewExtractColumnsByStructuralTypeStep(map[string]DataRef{"inputs": &StepDataRef{2, "produce"}}, []string{"produce"},
					[]string{
						"float",         // python type
						"numpy.float32", // numpy types
						"numpy.float64",
					}),
				NewSKImputerStep(map[string]DataRef{"inputs": &StepDataRef{3, "produce"}}, []string{"produce"}),
				NewSKMissingIndicatorStep(map[string]DataRef{"inputs": &StepDataRef{4, "produce"}}, []string{"produce"}),
				NewKMeansClusteringStep(map[string]DataRef{"inputs": &StepDataRef{5, "produce"}}, []string{"produce"}),
				NewConstructPredictionStep(map[string]DataRef{"inputs": &StepDataRef{6, "produce"}}, []string{"produce"}, &StepDataRef{1, "produce"}),
			}
		}
	} else {
		steps = []Step{
			NewDatasetToDataframeStep(map[string]DataRef{"inputs": &PipelineDataRef{0}}, []string{"produce"}),
			NewDistilColumnParserStep(map[string]DataRef{"inputs": &StepDataRef{0, "produce"}}, []string{"produce"}, []string{model.TA2RealType}),
			NewExtractColumnsByStructuralTypeStep(map[string]DataRef{"inputs": &StepDataRef{1, "produce"}}, []string{"produce"},
				[]string{
					"float",         // python type
					"numpy.float32", // numpy types
					"numpy.float64",
				}),
			NewSKImputerStep(map[string]DataRef{"inputs": &StepDataRef{2, "produce"}}, []string{"produce"}),
			NewSKMissingIndicatorStep(map[string]DataRef{"inputs": &StepDataRef{3, "produce"}}, []string{"produce"}),
			NewHDBScanStep(map[string]DataRef{"inputs": &StepDataRef{4, "produce"}}, []string{"produce"}),
			NewExtractColumnsStep(map[string]DataRef{"inputs": &StepDataRef{5, "produce"}}, []string{"produce"}, []int{-1}),
			// Needs to be added since the input dataset doesn't have a target, and hdbscan doesn't set the target itself.  Without this being
			// set the subsequent ConstructPredictions step doesn't work.
			NewAddSemanticTypeStep(map[string]DataRef{"inputs": &StepDataRef{6, "produce"}}, []string{"produce"}, &ColumnUpdate{
				Indices:       []int{0},
				SemanticTypes: []string{"https://metadata.datadrivendiscovery.org/types/PredictedTarget"},
			}),
			NewConstructPredictionStep(map[string]DataRef{"inputs": &StepDataRef{7, "produce"}}, []string{"produce"}, &StepDataRef{1, "produce"}),
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

// NewExtractColumnsBySemanticTypeStep extracts columns by supplied semantic types.
func NewExtractColumnsBySemanticTypeStep(inputs map[string]DataRef, outputMethods []string, semanticTypes []string) *StepData {
	return NewStepData(
		&pipeline.Primitive{
			Id:         "4503a4c6-42f7-45a1-a1d4-ed69699cf5e1",
			Version:    "0.4.0",
			Name:       "Extracts columns by semantic type",
			PythonPath: "d3m.primitives.data_transformation.extract_columns_by_semantic_types.Common",
			Digest:     "88f0780f5324d4a881d5d51e29f33fdcdc6d2968acf3b927032cf2d832e10504",
		},
		outputMethods,
		map[string]interface{}{"semantic_types": semanticTypes},
		inputs,
	)
}

// NewExtractColumnsByStructuralTypeStep extracts columns by supplied semantic types.
func NewExtractColumnsByStructuralTypeStep(inputs map[string]DataRef, outputMethods []string, structuralTypes []string) *StepData {
	return NewStepData(
		&pipeline.Primitive{
			Id:         "79674d68-9b93-4359-b385-7b5f60645b06",
			Version:    "0.1.0",
			Name:       "Extracts columns by structural type",
			PythonPath: "d3m.primitives.data_transformation.extract_columns_by_structural_types.Common",
			Digest:     "46fd01dd704090885e6dd80640a647b5596fbbe597dc1a1f9e5833ae4ab27813",
		},
		outputMethods,
		map[string]interface{}{"structural_types": structuralTypes},
		inputs,
	)
}

// NewEnrichDatesStep adds extra information for date fields.
func NewEnrichDatesStep(inputs map[string]DataRef, outputMethods []string, replace bool) *StepData {
	return NewStepData(
		&pipeline.Primitive{
			Id:         "b1367f5b-bab1-4dfc-a1a9-6a56430e516a",
			Version:    "0.5.1",
			Name:       "Enrich dates",
			PythonPath: "d3m.primitives.data_transformation.enrich_dates.DistilEnrichDates",
			Digest:     "176228cabb6746d544fc045f43ed02e5ae2a3747c39b492c008b2ce473fac650",
		},
		outputMethods,
		map[string]interface{}{"replace": replace},
		inputs,
	)
}

// NewListEncoderStep expands a list across columns.
func NewListEncoderStep(inputs map[string]DataRef, outputMethods []string) *StepData {
	return NewStepData(
		&pipeline.Primitive{
			Id:         "67f53b00-f936-4bb4-873e-4698c4aaa37f",
			Version:    "0.5.1",
			Name:       "List encoder",
			PythonPath: "d3m.primitives.data_transformation.list_to_dataframe.DistilListEncoder",
			Digest:     "2ab8726ff23664fdc22949ea6d389f29a3ea07355da9c7679f10cff93c9723f6",
		},
		outputMethods,
		map[string]interface{}{},
		inputs,
	)
}

// NewReplaceSingletonStep replaces a field that has only one value with a constant.
func NewReplaceSingletonStep(inputs map[string]DataRef, outputMethods []string) *StepData {
	return NewStepData(
		&pipeline.Primitive{
			Id:         "7cacc8b6-85ad-4c8f-9f75-360e0faee2b8",
			Version:    "0.5.1",
			Name:       "Replace singeltons",
			PythonPath: "d3m.primitives.data_transformation.replace_singletons.DistilReplaceSingletons",
			Digest:     "fdb2f8ce03d00243fcdd7279e42a43242d4e05a37674abe4da9d6a1ffdec3c2f",
		},
		outputMethods,
		map[string]interface{}{},
		inputs,
	)
}

// NewCategoricalImputerStep finds missing categorical values and replaces them with an imputed value.
func NewCategoricalImputerStep(inputs map[string]DataRef, outputMethods []string) *StepData {
	return NewStepData(
		&pipeline.Primitive{
			Id:         "0a9936f3-7784-4697-82f0-2a5fcc744c16",
			Version:    "0.5.1",
			Name:       "Categorical imputer",
			PythonPath: "d3m.primitives.data_transformation.imputer.DistilCategoricalImputer",
			Digest:     "aeab342f3a9a13caccac3fb0c44e3f625dc2ee1ab892aa0762844cad181b90e5",
		},
		outputMethods,
		map[string]interface{}{},
		inputs,
	)
}

// NewTextEncoderStep adds an svm text encoder for text fields.
func NewTextEncoderStep(inputs map[string]DataRef, outputMethods []string) *StepData {
	return NewStepData(
		&pipeline.Primitive{
			Id:         "09f252eb-215d-4e0b-9a60-fcd967f5e708",
			Version:    "0.5.1",
			Name:       "Text encoder",
			PythonPath: "d3m.primitives.data_transformation.encoder.DistilTextEncoder",
			Digest:     "a9e111cf5345e24e293101ddfaf1d001377246bc0f33bcaaef1b9144f13b5563",
		},
		outputMethods,
		map[string]interface{}{},
		inputs,
	)
}

// NewOneHotEncoderStep adds a one hot encoder for categoricals of low cardinality.
func NewOneHotEncoderStep(inputs map[string]DataRef, outputMethods []string) *StepData {
	return NewStepData(
		&pipeline.Primitive{
			Id:         "d3d421cb-9601-43f0-83d9-91a9c4199a06",
			Version:    "0.5.1",
			Name:       "One-hot encoder",
			PythonPath: "d3m.primitives.data_transformation.one_hot_encoder.DistilOneHotEncoder",
			Digest:     "66818496dda8d06badfb1a71f75172b50b9288c312b91952d73c0d4c7408480f",
		},
		outputMethods,
		map[string]interface{}{"max_one_hot": 16},
		inputs,
	)
}

// NewBinaryEncoderStep adds a binary encoder for categoricals of high cardinality.
func NewBinaryEncoderStep(inputs map[string]DataRef, outputMethods []string) *StepData {
	return NewStepData(
		&pipeline.Primitive{
			Id:         "d38e2e28-9b18-4ce4-b07c-9d809cd8b915",
			Version:    "0.5.1",
			Name:       "Binary encoder",
			PythonPath: "d3m.primitives.data_transformation.encoder.DistilBinaryEncoder",
			Digest:     "2a030799a0cd220d04ae9ccb109f49941e40aa6b6705d74a8ab5136dd5dfb0de",
		},
		outputMethods,
		map[string]interface{}{"min_binary": 17},
		inputs,
	)
}

// NewSKImputerStep adds SK learn simple imputer
func NewSKImputerStep(inputs map[string]DataRef, outputMethods []string) *StepData {
	return NewStepData(
		&pipeline.Primitive{
			Id:         "d016df89-de62-3c53-87ed-c06bb6a23cde",
			Version:    "2020.12.1",
			Name:       "sklearn.impute.SimpleImputer",
			PythonPath: "d3m.primitives.data_cleaning.imputer.SKlearn",
			Digest:     "84bf94c87a745011023da7074c65e1cee1272843d5a11cce1c64c7f20d42e408",
		},
		outputMethods,
		map[string]interface{}{
			"use_semantic_types": true,
			"error_on_no_input":  false,
			"return_result":      "replace",
		},
		inputs,
	)
}

// NewSKMissingIndicatorStep adds SK learn missing indicator.
func NewSKMissingIndicatorStep(inputs map[string]DataRef, outputMethods []string) *StepData {
	return NewStepData(
		&pipeline.Primitive{
			Id:         "94c5c918-9ad5-3496-8e52-2359056e0120",
			Version:    "2020.12.1",
			Name:       "sklearn.impute.MissingIndicator",
			PythonPath: "d3m.primitives.data_cleaning.missing_indicator.SKlearn",
			Digest:     "def3920732545bbb752d15c9a34e077c7eb7c9462c5025c509decf5b667e9ebc",
		},
		outputMethods,
		map[string]interface{}{
			"use_semantic_types": true,
			"error_on_new":       false,
			"error_on_no_input":  false,
			"return_result":      "append",
		},
		inputs,
	)
}

// NewHDBScanStep adds clustering features.
func NewHDBScanStep(inputs map[string]DataRef, outputMethods []string) *StepData {
	return NewStepData(
		&pipeline.Primitive{
			Id:         "ca014488-6004-4b54-9403-5920fbe5a834",
			Version:    "1.0.2",
			Name:       "hdbscan",
			PythonPath: "d3m.primitives.clustering.hdbscan.Hdbscan",
			Digest:     "a284fdd5802ec6bb76e78f04ce30c496b42841df2a0e550e6468d709aca3a8ca",
		},
		outputMethods,
		map[string]interface{}{"required_output": "feature", "min_samples": 1, "min_cluster_size": 50, "cluster_selection_method": "leaf"},
		inputs,
	)
}

// NewSatelliteImageLoaderStep loads multi band images.
func NewSatelliteImageLoaderStep(inputs map[string]DataRef, outputMethods []string, numJobs int) *StepData {
	return NewStepData(
		&pipeline.Primitive{
			Id:         "77d20419-aeb6-44f9-8e63-349ea5b654f7",
			Version:    "0.5.1",
			Name:       "Columns satellite image loader",
			PythonPath: "d3m.primitives.data_transformation.satellite_image_loader.DistilSatelliteImageLoader",
			Digest:     "98856217d8ffe69d27754810f9d3a09ab9fcdaf5029f0f8a2e6a1bc8d7f0975c",
		},
		outputMethods,
		map[string]interface{}{"return_result": "replace", "compress_data": true, "n_jobs": numJobs},
		inputs,
	)
}

// NewRemoteSensingPretrainedStep featurizes a remote sensing column
func NewRemoteSensingPretrainedStep(inputs map[string]DataRef, outputMethods []string, batchSize int, pool bool) *StepData {
	return NewStepData(
		&pipeline.Primitive{
			Id:         "544bb61f-f354-48f5-b055-5c03de71c4fb",
			Version:    "1.0.0",
			Name:       "RSPretrained",
			PythonPath: "d3m.primitives.remote_sensing.remote_sensing_pretrained.RemoteSensingPretrained",
			Digest:     "a3e089d258d9be2fd08d40554d3f2c6ef0c7660242dbfbf5878905b5d63790d1",
		},
		outputMethods,
		map[string]interface{}{"batch_size": batchSize, "decompress_data": true, "pool_features": pool},
		inputs,
	)
}
