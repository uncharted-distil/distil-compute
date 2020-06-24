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
	"github.com/uncharted-distil/distil-compute/model"
	"github.com/uncharted-distil/distil-compute/pipeline"
)

// CreateGeneralClusteringPipeline creates a pipeline that will cluster tabular data.
func CreateGeneralClusteringPipeline(name string, description string, datasetDescription *UserDatasetDescription,
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
		[]string{model.TA2IntegerType, "https://metadata.datadrivendiscovery.org/types/FloatVector", model.TA2RealType},
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

	steps = append(steps, NewEnrichDatesStep(map[string]DataRef{"inputs": &StepDataRef{attributeStep, "produce"}}, []string{"produce"}))
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

	steps = append(steps, NewHDBScanStep(map[string]DataRef{"inputs": &StepDataRef{offset, "produce"}}, []string{"produce"}))
	offset = offset + 1

	steps = append(steps, NewExtractColumnsStep(map[string]DataRef{"inputs": &StepDataRef{offset, "produce"}}, []string{"produce"}, []int{-1}))
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

// NewExtractColumnsBySemanticTypeStep extracts columns by supplied semantic types.
func NewExtractColumnsBySemanticTypeStep(inputs map[string]DataRef, outputMethods []string, semanticTypes []string) *StepData {
	return NewStepData(
		&pipeline.Primitive{
			Id:         "4503a4c6-42f7-45a1-a1d4-ed69699cf5e1",
			Version:    "0.4.0",
			Name:       "Extracts columns by semantic type",
			PythonPath: "d3m.primitives.data_transformation.extract_columns_by_semantic_types.Common",
			Digest:     "cf44b2f5af90f10ef9935496655a202bfc8a4a0fa24b8e9d733ee61f096bda87",
		},
		outputMethods,
		map[string]interface{}{"semantic_types": semanticTypes},
		inputs,
	)
}

// NewEnrichDatesStep adds extra information for date fields.
func NewEnrichDatesStep(inputs map[string]DataRef, outputMethods []string) *StepData {
	return NewStepData(
		&pipeline.Primitive{
			Id:         "b1367f5b-bab1-4dfc-a1a9-6a56430e516a",
			Version:    "0.4.0",
			Name:       "Enrich dates",
			PythonPath: "d3m.primitives.data_transformation.data_cleaning.DistilEnrichDates",
			Digest:     "ab9cd162ac1ee1416184f468da8d4786a29727ad61bbba1cf552d741438b365a",
		},
		outputMethods,
		map[string]interface{}{},
		inputs,
	)
}

// NewListEncoderStep expands a list across columns.
func NewListEncoderStep(inputs map[string]DataRef, outputMethods []string) *StepData {
	return NewStepData(
		&pipeline.Primitive{
			Id:         "67f53b00-f936-4bb4-873e-4698c4aaa37f",
			Version:    "0.4.0",
			Name:       "List encoder",
			PythonPath: "d3m.primitives.data_transformation.list_to_dataframe.DistilListEncoder",
			Digest:     "c99a3fc777bcfdebbd1f8c746e79cad71ec181d5978061b4f7cd82f6330daad6",
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
			Version:    "0.4.0",
			Name:       "Replace singeltons",
			PythonPath: "d3m.primitives.data_transformation.data_cleaning.DistilReplaceSingletons",
			Digest:     "40dfe842797d1513ad962d81c01a78af405b5a4409aaed82cd90fc4b04ac7e32",
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
			Version:    "0.4.0",
			Name:       "Categorical imputer",
			PythonPath: "d3m.primitives.data_transformation.imputer.DistilCategoricalImputer",
			Digest:     "0ad4182f53c57146b1817c6b91505103d2867fed75d8d934de66ef04705b8c9b",
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
			Version:    "0.4.0",
			Name:       "Text encoder",
			PythonPath: "d3m.primitives.data_transformation.encoder.DistilTextEncoder",
			Digest:     "67df378139975454858989b666d63a319bf7bf64001971a4a3f601e9b60ad36a",
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
			Version:    "0.4.0",
			Name:       "One-hot encoder",
			PythonPath: "d3m.primitives.data_transformation.one_hot_encoder.DistilOneHotEncoder",
			Digest:     "9ea16f751325297f9347b105c16c0526e8d1294616c3390fb38997a15418a65e",
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
			Version:    "0.4.0",
			Name:       "Binary encoder",
			PythonPath: "d3m.primitives.data_transformation.encoder.DistilBinaryEncoder",
			Digest:     "f3874916967418450b3bd5575446219bacdd9bf0679891436d97628da26135ae",
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
			Version:    "2020.6.10",
			Name:       "sklearn.impute.SimpleImputer",
			PythonPath: "d3m.primitives.data_cleaning.imputer.SKlearn",
			Digest:     "5cdf2101de052235f8231419be7e2f80190147c213d63c841bc770fdcfffa76f",
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
			Version:    "2020.6.10",
			Name:       "sklearn.impute.MissingIndicator",
			PythonPath: "d3m.primitives.data_cleaning.missing_indicator.SKlearn",
			Digest:     "f390c8e595f48df5848d919aa9db4b4c8791732b368b608320d882e383c4e4eb",
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
			Digest:     "e805ded4d975e125d257a74f8e50f003d782605137dfedbe8f5e567e3607c219",
		},
		outputMethods,
		map[string]interface{}{"required_output": "feature"},
		inputs,
	)
}

// NewExtractColumnsStep extracts columns by supplied indices.
func NewExtractColumnsStep(inputs map[string]DataRef, outputMethods []string, indices []int) *StepData {
	return NewStepData(
		&pipeline.Primitive{
			Id:         "81d7e261-e25b-4721-b091-a31cd46e99ae",
			Version:    "0.1.0",
			Name:       "Extracts columns",
			PythonPath: "d3m.primitives.data_transformation.extract_columns.Common",
			Digest:     "cf44b2f5af90f10ef9935496655a202bfc8a4a0fa24b8e9d733ee61f096bda87",
		},
		outputMethods,
		map[string]interface{}{"columns": indices},
		inputs,
	)
}