//
//   Copyright © 2020 Uncharted Software Inc.
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
	"github.com/uncharted-distil/distil-compute/pipeline"
	"github.com/uncharted-distil/distil-compute/primitive/compute"
)

// NewSimonStep creates a SIMON data classification step.  It examines an input
// dataframe, and assigns types to the columns based on the exposed metadata.
func NewSimonStep(inputs map[string]DataRef, outputMethods []string) *StepData {
	return NewStepData(
		&pipeline.Primitive{
			Id:         "d2fa8df2-6517-3c26-bafc-87b701c4043a",
			Version:    "1.2.3",
			Name:       "simon",
			PythonPath: "d3m.primitives.data_cleaning.column_type_profiler.Simon",
			Digest:     "75708b30bd2d1d4bf319ea995076d1c108d9e7fede60162d48c3300bbfb37a22",
		},
		outputMethods,
		map[string]interface{}{"statistical_classification": true, "p_threshold": 0.9},
		inputs,
	)
}

// NewDataframeImageReaderStep reads images for further processing.
func NewDataframeImageReaderStep(inputs map[string]DataRef, outputMethods []string, columns []int) *StepData {
	return NewStepData(
		&pipeline.Primitive{
			Id:         "8f2e51e8-da59-456d-ae29-53912b2b9f3d",
			Version:    "0.2.0",
			Name:       "Columns image reader",
			PythonPath: "d3m.primitives.data_transformation.image_reader.Common",
			Digest:     "84ed3998d43985188dabc5743654b0d2370493c72d025a177c6a80c053e9754d",
		},
		outputMethods,
		map[string]interface{}{"use_columns": columns, "return_result": "replace"},
		inputs,
	)
}

// NewImageTransferStep processes images.
func NewImageTransferStep(inputs map[string]DataRef, outputMethods []string) *StepData {
	return NewStepData(
		&pipeline.Primitive{
			Id:         "782e261e-8e23-4184-9258-5a412c9b32d4",
			Version:    "0.5.1",
			Name:       "Image Transfer",
			PythonPath: "d3m.primitives.feature_extraction.image_transfer.DistilImageTransfer",
			Digest:     "aebfc43c3278946d2ba61b2b674e2189a3db2f106a9859eca3d11174d5600845",
		},
		outputMethods,
		map[string]interface{}{},
		inputs,
	)
}

// NewKMeansClusteringStep clusters the input using a siple k-means clustering.
func NewKMeansClusteringStep(inputs map[string]DataRef, outputMethods []string) *StepData {
	return NewStepData(
		&pipeline.Primitive{
			Id:         "3b09024e-a83b-418c-8ff4-cf3d30a9609e",
			Version:    "0.5.1",
			Name:       "K means",
			PythonPath: "d3m.primitives.clustering.k_means.DistilKMeans",
			Digest:     "2e95d33622c9804911fe006581adec68220ebd41fc1f1b9c084bf5262dc0b964",
		},
		outputMethods,
		map[string]interface{}{"n_clusters": 4},
		inputs,
	)
}

// NewSlothStep creates a Sloth timeseries clustering step.
func NewSlothStep(inputs map[string]DataRef, outputMethods []string) *StepData {
	return NewStepData(
		&pipeline.Primitive{
			Id:         "77bf4b92-2faa-3e38-bb7e-804131243a7f",
			Version:    "2.0.5",
			Name:       "Sloth",
			PythonPath: "d3m.primitives.clustering.k_means.Sloth",
			Digest:     "11d6cd6f54825147b49f82ee4d4b8a2340a597efb2b1efdde8f0221f664181fb",
		},
		outputMethods,
		map[string]interface{}{"nclusters": 4},
		inputs,
	)
}

// NewPCAFeaturesStep creates a PCA-based feature ranking call that can be added to
// a pipeline.
func NewPCAFeaturesStep(inputs map[string]DataRef, outputMethods []string) *StepData {
	// since PCA has fit & produce, need to set the params from
	// set_training_data. In this case, outputs is not used.
	if inputs["inputs"] != nil && inputs["outputs"] == nil {
		inputs["outputs"] = inputs["inputs"]
	}

	return NewStepData(
		&pipeline.Primitive{
			Id:         "04573880-d64f-4791-8932-52b7c3877639",
			Version:    "3.1.2",
			Name:       "PCA Features",
			PythonPath: "d3m.primitives.feature_selection.pca_features.Pcafeatures",
			Digest:     "28a77a528b8e6a557ba8d96eba655bf0093a12c9fc355ff2263badcfada7ae3e",
		},
		outputMethods,
		map[string]interface{}{},
		inputs,
	)
}

// NewProfilerStep creates a profile primitive that infers the columns type using rules
func NewProfilerStep(inputs map[string]DataRef, outputMethods []string) *StepData {
	return NewStepData(
		&pipeline.Primitive{
			Id:         "e193afa1-b45e-4d29-918f-5bb1fa3b88a7",
			Version:    "0.2.0",
			Name:       "Determine missing semantic types for columns automatically",
			PythonPath: "d3m.primitives.schema_discovery.profiler.Common",
			Digest:     "cf5ff86d612b8b6d55687cee8d1e51522845e763aba01553915b3949b50c2076",
		},
		outputMethods,
		map[string]interface{}{},
		inputs,
	)
}

// NewTargetRankingStep creates a target ranking call that can be added to
// a pipeline. Ranking is based on mutual information between features and a selected
// target.  Returns a DataFrame containing (col_idx, col_name, score) tuples for
// each ranked feature. Features that could not be ranked are excluded
// from the returned set.
func NewTargetRankingStep(inputs map[string]DataRef, outputMethods []string, targetCol int) *StepData {
	return NewStepData(
		&pipeline.Primitive{
			Id:         "a31b0c26-cca8-4d54-95b9-886e23df8886",
			Version:    "0.5.1",
			Name:       "Mutual Information Feature Ranking",
			PythonPath: "d3m.primitives.feature_selection.mutual_info_classif.DistilMIRanking",
			Digest:     "8c39c4f44946afce53e1adf8fc70a55579f1ea9182bbfe6743bcd444c11cc58e",
		},
		outputMethods,
		map[string]interface{}{"target_col_index": targetCol},
		inputs,
	)
}

// NewDukeStep creates a wrapper for the Duke dataset classifier.
func NewDukeStep(inputs map[string]DataRef, outputMethods []string) *StepData {
	return NewStepData(
		&pipeline.Primitive{
			Id:         "46612a42-6120-3559-9db9-3aa9a76eb94f",
			Version:    "1.2.0",
			Name:       "duke",
			PythonPath: "d3m.primitives.data_cleaning.text_summarization.Duke",
			Digest:     "d726f8df18fe5b56834869dabc8c742295b95b2786fb0a5159dfcd0f93fbc3e1",
		},
		outputMethods,
		map[string]interface{}{},
		inputs,
	)
}

// NewDataCleaningStep creates a wrapper for the Punk data cleaning primitive.
func NewDataCleaningStep(inputs map[string]DataRef, outputMethods []string) *StepData {
	return NewStepData(
		&pipeline.Primitive{
			Id:         "fc6bf33a-f3e0-3496-aa47-9a40289661bc",
			Version:    "3.0.2",
			Name:       "Data cleaning",
			PythonPath: "d3m.primitives.data_cleaning.data_cleaning.Datacleaning",
			Digest:     "b29d22256c745e37a8a63672f5c93cbd9f84382e25c25f3775e3a28c8c98fd2f",
		},
		outputMethods,
		map[string]interface{}{},
		inputs,
	)
}

// NewDatamartDownloadStep creates a primitive call that downloads a dataset
// from a datamart.
func NewDatamartDownloadStep(inputs map[string]DataRef, outputMethods []string, searchResult string, systemIdentifier string) *StepData {
	// supplied_id and supplied_resource_id need to be part of search result.
	//   supplied_id: dataset id of the linked dataset
	//   supplied_resource_id: resource id of the dataset
	// searchResult is a json struct so ends with '}'
	// simply update that search result to fit in the required params
	//searchResult = strings.TrimSpace(searchResult)
	//searchResult = fmt.Sprintf(`%s, "supplied_id": "%s", "supplied_resource_id": "%s"}`,
	//	searchResult[:len(searchResult)-1],
	//	dataset,
	//	defaultResource,
	//)

	return NewStepData(
		&pipeline.Primitive{
			Id:         "9e2077eb-3e38-4df1-99a5-5e647d21331f",
			Version:    "0.1",
			Name:       "Download a dataset from Datamart",
			PythonPath: "d3m.primitives.data_augmentation.datamart_download.Common",
			Digest:     "7e92079cf5dd2052e93ad152d626fc16670f0dde0ae19433a2e8ce7bf2dc7746",
		},
		outputMethods,
		map[string]interface{}{
			"search_result":     searchResult,
			"system_identifier": systemIdentifier,
		},
		inputs,
	)
}

// NewDatamartAugmentStep creates a primitive call that augments a dataset
// with a datamart dataset.
func NewDatamartAugmentStep(inputs map[string]DataRef, outputMethods []string, searchResult string, systemIdentifier string) *StepData {
	// supplied_id and supplied_resource_id need to be part of search result.
	//   supplied_id: dataset id of the linked dataset
	//   supplied_resource_id: resource id of the dataset
	// searchResult is a json struct so ends with '}'
	// simply update that search result to fit in the required params
	//searchResult = strings.TrimSpace(searchResult)
	//searchResult = fmt.Sprintf(`%s, "supplied_id": "%s", "supplied_resource_id": "%s"}`,
	//	searchResult[:len(searchResult)-1],
	//	dataset,
	//	defaultResource,
	//)

	return NewStepData(
		&pipeline.Primitive{
			Id:         "fe0f1ac8-1d39-463a-b344-7bd498a31b91",
			Version:    "0.1",
			Name:       "Perform dataset augmentation using Datamart",
			PythonPath: "d3m.primitives.data_augmentation.datamart_augmentation.Common",
			Digest:     "24b33cabe7514971aa7dffcc34eae3fb3b4755ea713a910a0732b73036132fbb",
		},
		outputMethods,
		map[string]interface{}{
			"search_result":     searchResult,
			"system_identifier": systemIdentifier,
		},
		inputs,
	)
}

// NewDatasetToDataframeStep creates a primitive call that transforms an input dataset
// into a PANDAS dataframe.
func NewDatasetToDataframeStep(inputs map[string]DataRef, outputMethods []string) *StepData {
	return NewStepData(
		&pipeline.Primitive{
			Id:         "4b42ce1e-9b98-4a25-b68e-fad13311eb65",
			Version:    "0.3.0",
			Name:       "Extract a DataFrame from a Dataset",
			PythonPath: "d3m.primitives.data_transformation.dataset_to_dataframe.Common",
			Digest:     "aed657e5effa3e313bd0e59c7334100aa8552fc5aba762a959ce4569284a5e63",
		},
		outputMethods,
		map[string]interface{}{},
		inputs,
	)
}

// NewGroupingFieldComposeStep creates a primitive call that joins suggested grouping keys.
func NewGroupingFieldComposeStep(inputs map[string]DataRef, outputMethods []string, colIndices []int, joinChar string, outputName string) *StepData {

	return NewStepData(
		&pipeline.Primitive{
			Id:         "59db88b9-dd81-4e50-8f43-8f2af959560b",
			Version:    "0.1.1",
			Name:       "Grouping Field Compose",
			PythonPath: "d3m.primitives.data_transformation.grouping_field_compose.Common",
			Digest:     "5e47a5fe1b53950de84de9603573517f2361b3a028268b52ec3846ec24c5cd14",
		},
		outputMethods,
		map[string]interface{}{
			"columns":     colIndices,
			"join_char":   joinChar,
			"output_name": outputName,
		},
		inputs,
	)
}

// NewHorizontalConcatStep creates a primitive call that concats two data frames.
func NewHorizontalConcatStep(inputs map[string]DataRef, outputMethods []string, useIndex bool, removeSecondIndex bool) *StepData {

	return NewStepData(
		&pipeline.Primitive{
			Id:         "aff6a77a-faa0-41c5-9595-de2e7f7c4760",
			Version:    "0.2.0",
			Name:       "Concatenate two dataframes",
			PythonPath: "d3m.primitives.data_transformation.horizontal_concat.DataFrameCommon",
			Digest:     "f1e8fe6ba0456e562d9613bd5f4221e221e9cadd23c684564137b2aa14495ada",
		},
		outputMethods,
		map[string]interface{}{
			"use_index":           useIndex,
			"remove_second_index": removeSecondIndex,
		},
		inputs,
	)
}

// NewDatasetToDataframeStepWithResource creates a primitive call that transforms an input dataset
// into a PANDAS dataframe using the specified resource.
func NewDatasetToDataframeStepWithResource(inputs map[string]DataRef, outputMethods []string, resourceName string) *StepData {
	if resourceName == "" {
		resourceName = compute.DefaultResourceID
	}

	return NewStepData(
		&pipeline.Primitive{
			Id:         "4b42ce1e-9b98-4a25-b68e-fad13311eb65",
			Version:    "0.3.0",
			Name:       "Extract a DataFrame from a Dataset",
			PythonPath: "d3m.primitives.data_transformation.dataset_to_dataframe.Common",
			Digest:     "aed657e5effa3e313bd0e59c7334100aa8552fc5aba762a959ce4569284a5e63",
		},
		outputMethods,
		map[string]interface{}{
			"dataframe_resource": resourceName,
		},
		inputs,
	)
}

// NewDatasetWrapperStep creates a primitive that wraps a dataframe primitive such that it can be
// used as a datset primitive in the pipeline prepend.  The primitive to wrap is indicated using its
// index in the pipeline.    Leaving the resource ID as the empty value allows the primitive to infer
// the main resource from the dataset.
func NewDatasetWrapperStep(inputs map[string]DataRef, outputMethods []string, primitiveIndex int, resourceID string) *StepData {

	hyperparams := map[string]interface{}{
		"primitive": &PrimitiveReference{primitiveIndex},
	}
	if resourceID != "" {
		hyperparams["resources"] = []string{resourceID}
	}

	return NewStepData(
		&pipeline.Primitive{
			Id:         "5bef5738-1638-48d6-9935-72445f0eecdc",
			Version:    "0.1.0",
			Name:       "Map DataFrame resources to new resources using provided primitive",
			PythonPath: "d3m.primitives.operator.dataset_map.DataFrameCommon",
			Digest:     "851bfdeb8edda190fd76ec0d3c44ad6092c0009bb92ab41c62ddb7e68124b340",
		},
		outputMethods,
		hyperparams,
		inputs,
	)
}

// ColumnUpdate defines a set of column indices to add/remvoe
// a set of semantic types to/from.
type ColumnUpdate struct {
	Indices       []int
	SemanticTypes []string
}

// NewAddSemanticTypeStep adds semantic data values to an input
// dataset.  An add of (1, 2), ("type a", "type b") would result in "type a" and "type b"
// being added to index 1 and 2.
func NewAddSemanticTypeStep(inputs map[string]DataRef, outputMethods []string, add *ColumnUpdate) *StepData {
	return NewStepData(
		&pipeline.Primitive{
			Id:         "d7e14b12-abeb-42d8-942f-bdb077b4fd37",
			Version:    "0.1.0",
			Name:       "Add semantic types to columns",
			PythonPath: "d3m.primitives.data_transformation.add_semantic_types.Common",
			Digest:     "1c326225c8aec10a65b03cb22d7a77b26c2a360a0ed038cd03b3ffe2cb7803fe",
		},
		outputMethods,
		map[string]interface{}{
			"columns":        add.Indices,
			"semantic_types": add.SemanticTypes,
		},
		inputs,
	)
}

// NewRemoveSemanticTypeStep removes semantic data values from an input
// dataset.  A remove of (1, 2), ("type a", "type b") would result in "type a" and "type b"
// being removed from index 1 and 2.
func NewRemoveSemanticTypeStep(inputs map[string]DataRef, outputMethods []string, remove *ColumnUpdate) *StepData {
	return NewStepData(
		&pipeline.Primitive{
			Id:         "3002bc5b-fa47-4a3d-882e-a8b5f3d756aa",
			Version:    "0.1.0",
			Name:       "Remove semantic types from columns",
			PythonPath: "d3m.primitives.data_transformation.remove_semantic_types.Common",
			Digest:     "8a49eb8c0e96cb9b7f85b0a9c9e45fc7765a814d3c79f233bca25489511ee951",
		},
		outputMethods,
		map[string]interface{}{
			"columns":        remove.Indices,
			"semantic_types": remove.SemanticTypes,
		},
		inputs,
	)
}

// NewDenormalizeStep denormalize data that is contained in multiple resource files.
func NewDenormalizeStep(inputs map[string]DataRef, outputMethods []string) *StepData {
	return NewStepData(
		&pipeline.Primitive{
			Id:         "f31f8c1f-d1c5-43e5-a4b2-2ae4a761ef2e",
			Version:    "0.2.0",
			Name:       "Denormalize datasets",
			PythonPath: "d3m.primitives.data_transformation.denormalize.Common",
			Digest:     "eb853a0744ffef87a6d2689cdbb2462a3d0fd1e4645389b8665093d1c2f53604",
		},
		outputMethods,
		map[string]interface{}{},
		inputs,
	)
}

// NewCSVReaderStep reads data from csv files into a nested dataframe structure.
func NewCSVReaderStep(inputs map[string]DataRef, outputMethods []string) *StepData {
	hyperparams := map[string]interface{}{
		"return_result": "append",
	}
	return NewStepData(
		&pipeline.Primitive{
			Id:         "989562ac-b50f-4462-99cb-abef80d765b2",
			Version:    "0.1.0",
			Name:       "Columns CSV reader",
			PythonPath: "d3m.primitives.data_transformation.csv_reader.Common",
			Digest:     "a4a0cf41e1a1ba4c69527cf98bef206b1330a29a91f6037d1d45b9a45307baa4",
		},
		outputMethods,
		hyperparams,
		inputs,
	)
}

// NewDataFrameFlattenStep searches for nested dataframes and pulls them out.
func NewDataFrameFlattenStep(inputs map[string]DataRef, outputMethods []string) *StepData {
	hyperparams := map[string]interface{}{
		"return_result": "replace",
	}
	return NewStepData(
		&pipeline.Primitive{
			Id:         "1c4aed23-f3d3-4e6b-9710-009a9bc9b694",
			Version:    "0.1.0",
			Name:       "DataFrame Flatten",
			PythonPath: "d3m.primitives.data_transformation.flatten.DataFrameCommon",
			Digest:     "15235c3a8e088e33e733625704d8f208b05d4f6bfa0da60c172bbf8af31b51f5",
		},
		outputMethods,
		hyperparams,
		inputs,
	)
}

// NewConstructPredictionStep maps the dataframe index to d3m index.
func NewConstructPredictionStep(inputs map[string]DataRef, outputMethods []string, reference DataRef) *StepData {
	args := map[string]DataRef{"reference": reference}
	for k, c := range inputs {
		args[k] = c
	}

	return NewStepData(
		&pipeline.Primitive{
			Id:         "8d38b340-f83f-4877-baaa-162f8e551736",
			Version:    "0.3.0",
			Name:       "Construct pipeline predictions output",
			PythonPath: "d3m.primitives.data_transformation.construct_predictions.Common",
			Digest:     "7ecceddd6bf78f4a8b0719f1aff46fe2e549c0b4b096be035513a92bdb6510de",
		},
		outputMethods,
		nil,
		args,
	)
}

// NewColumnParserStep takes obj/string columns in a dataframe and parses them into their
// associated raw python types based on the attached d3m metadata.
func NewColumnParserStep(inputs map[string]DataRef, outputMethods []string, types []string) *StepData {
	return NewStepData(
		&pipeline.Primitive{
			Id:         "d510cb7a-1782-4f51-b44c-58f0236e47c7",
			Version:    "0.6.0",
			Name:       "Parses strings into their types",
			PythonPath: "d3m.primitives.data_transformation.column_parser.Common",
			Digest:     "6f73dc863e2cfcbed90757ab26c34ca8df23e24f9a26632f48dc228f2277dc7b",
		},
		outputMethods,
		map[string]interface{}{"parse_semantic_types": types, "parse_categorical_target_columns": true},
		inputs,
	)
}

// NewDistilColumnParserStep takes obj/string columns in a datafram and parsaer them into raw python types
// based on their metadata.  Avoids some performance issues present in the common ColumnParser but does not
// support as many data types.
func NewDistilColumnParserStep(inputs map[string]DataRef, outputMethods []string, types []string) *StepData {
	return NewStepData(
		&pipeline.Primitive{
			Id:         "e8e78214-9770-4c26-9eae-a45bd0ede91a",
			Version:    "0.5.1",
			Name:       "Column Parser",
			PythonPath: "d3m.primitives.data_transformation.column_parser.DistilColumnParser",
			Digest:     "7e1def8c114a73394ab17d0463763f0e11896941b74122681f61e5e34dd3a073",
		},
		outputMethods,
		map[string]interface{}{"parsing_semantics": types},
		inputs,
	)
}

// NewRemoveColumnsStep removes columns from an input dataframe.  Columns
// are specified by name and the match is case insensitive.
func NewRemoveColumnsStep(inputs map[string]DataRef, outputMethods []string, colIndices []int) *StepData {
	return NewStepData(
		&pipeline.Primitive{
			Id:         "3b09ba74-cc90-4f22-9e0a-0cf4f29a7e28",
			Version:    "0.1.0",
			Name:       "Removes columns",
			PythonPath: "d3m.primitives.data_transformation.remove_columns.Common",
			Digest:     "7580c8f483a1a3607c6ac4ff1218c49b0515a5f05e7eb3c956601a3fd5e9b0a2",
		},
		outputMethods,
		map[string]interface{}{
			"columns": colIndices,
		},
		inputs,
	)
}

// NewRemoveDuplicateColumnsStep removes duplicate columns from a dataframe.
func NewRemoveDuplicateColumnsStep(inputs map[string]DataRef, outputMethods []string) *StepData {
	return NewStepData(
		&pipeline.Primitive{
			Id:         "130513b9-09ca-4785-b386-37ab31d0cf8b",
			Version:    "0.1.0",
			Name:       "Removes duplicate columns",
			PythonPath: "d3m.primitives.data_transformation.remove_duplicate_columns.Common",
			Digest:     "73eb50696e7be12fffff21f8abf0b59736a9e72d5d729206dd6f59fa34022f29",
		},
		outputMethods,
		map[string]interface{}{},
		inputs,
	)
}

// NewTermFilterStep creates a primitive step that filters dataset rows based on a match against a
// term list.  The term match can be partial, or apply to whole terms only.
func NewTermFilterStep(inputs map[string]DataRef, outputMethods []string, colindex int, inclusive bool, terms []string, matchWhole bool) *StepData {
	return NewStepData(
		&pipeline.Primitive{
			Id:         "a6b27300-4625-41a9-9e91-b4338bfc219b",
			Version:    "0.1.0",
			Name:       "Term list dataset filter",
			PythonPath: "d3m.primitives.data_transformation.term_filter.Common",
			Digest:     "36ba569e75422d69f144228557d4528d8a5f5a982e9d5e3a32b2164358906363",
		},
		outputMethods,
		map[string]interface{}{
			"column":      colindex,
			"inclusive":   inclusive,
			"terms":       terms,
			"match_whole": matchWhole,
		},
		inputs,
	)
}

// NewRegexFilterStep creates a primitive step that filter dataset rows based on a regex match.
func NewRegexFilterStep(inputs map[string]DataRef, outputMethods []string, colindex int, inclusive bool, regex string) *StepData {
	hyperparams := map[string]interface{}{
		"column":    colindex,
		"inclusive": inclusive,
		"regex":     regex,
	}
	return NewStepData(
		&pipeline.Primitive{
			Id:         "cf73bb3d-170b-4ba9-9ead-3dd4b4524b61",
			Version:    "0.1.0",
			Name:       "Regex dataset filter",
			PythonPath: "d3m.primitives.data_transformation.regex_filter.Common",
			Digest:     "b07fc3bdd424fcc31f617e32ae4fef555bacfed3e7bf54ed5e5b8a5508bd1e75",
		},
		outputMethods,
		hyperparams,
		inputs,
	)
}

// NewNumericRangeFilterStep creates a primitive step that filters dataset rows based on an
// included/excluded numeric range.  Inclusion of boundaries is controlled by the strict flag.
func NewNumericRangeFilterStep(inputs map[string]DataRef, outputMethods []string, colindex int, inclusive bool, min float64, max float64, strict bool) *StepData {
	return NewStepData(
		&pipeline.Primitive{
			Id:         "8c246c78-3082-4ec9-844e-5c98fcc76f9d",
			Version:    "0.1.0",
			Name:       "Numeric range filter",
			PythonPath: "d3m.primitives.data_transformation.numeric_range_filter.Common",
			Digest:     "48420b15f69ac745be0096750d5853e71b2d364e59100e6ba9b58ece93c5dbbe",
		},
		outputMethods,
		map[string]interface{}{
			"column":    colindex,
			"inclusive": inclusive,
			"min":       min,
			"max":       max,
			"strict":    strict,
		},
		inputs,
	)
}

// NewDateTimeRangeFilterStep creates a primitive step that filters dataset rows based on an
// included/excluded date/time range.  Inclusion of boundaries is controlled by the strict flag.
// Min and Max values are a unix timestamp expressed as floats.
func NewDateTimeRangeFilterStep(inputs map[string]DataRef, outputMethods []string, colindex int, inclusive bool, min float64, max float64, strict bool) *StepData {
	return NewStepData(
		&pipeline.Primitive{
			Id:         "487e5a58-19e9-432c-ac61-fe05c024e42c",
			Version:    "0.2.0",
			Name:       "Datetime range filter",
			PythonPath: "d3m.primitives.data_transformation.datetime_range_filter.Common",
			Digest:     "6ff7ec1bba7d5241992da0a1acdb1aaad5ad3fa6fcf139ed7ab059c6eb59b78d",
		},
		outputMethods,
		map[string]interface{}{
			"column":    colindex,
			"inclusive": inclusive,
			"min":       min,
			"max":       max,
			"strict":    strict,
		},
		inputs,
	)
}

// NewVectorBoundsFilterStep creates a primitive that will allow for a vector of values to be filtered included/excluded value range.
// The input min and max ranges are specified as lists, where the i'th element of the min/max lists are applied to the i'th value of the target vectors
// as the filter.
func NewVectorBoundsFilterStep(inputs map[string]DataRef, outputMethods []string, column int, inclusive bool, min []float64, max []float64, strict bool) *StepData {
	return NewStepData(
		&pipeline.Primitive{
			Id:         "c2fa34c0-2d1b-42af-91d2-515da4a27752",
			Version:    "0.5.1",
			Name:       "Vector bounds filter",
			PythonPath: "d3m.primitives.data_transformation.vector_bounds_filter.DistilVectorBoundsFilter",
			Digest:     "face2225a76870fe1d29ddb83be06f19ecadab52532c0b7ec011f2f469c70321",
		},
		outputMethods,
		map[string]interface{}{
			"column":    column,
			"inclusive": inclusive,
			"mins":      min,
			"maxs":      max,
			"strict":    strict,
		},
		inputs,
	)
}

// NewGoatForwardStep creates a GOAT forward geocoding primitive.  A string column
// containing a place name or address is passed in, and the primitive will
// return a DataFrame containing the lat/lon coords of the place.  If location could
// not be found, the row in the data frame will be empty.
func NewGoatForwardStep(inputs map[string]DataRef, outputMethods []string, placeColIndex int) *StepData {
	args := map[string]interface{}{
		"target_columns": []int{placeColIndex},
		"rampup_timeout": 150,
	}
	return NewStepData(
		&pipeline.Primitive{
			Id:         "c7c61da3-cf57-354e-8841-664853370106",
			Version:    "1.0.8",
			Name:       "Goat_forward",
			PythonPath: "d3m.primitives.data_cleaning.geocoding.Goat_forward",
			Digest:     "23f291c052696d85087b8513f1142db39a844f76615b65aa67aa5714ac5115b4",
		},
		outputMethods,
		args,
		inputs,
	)
}

// NewGoatReverseStep creates a GOAT reverse geocoding primitive.  Columns
// containing lat and lon values are passed in, and the primitive will
// return a DataFrame containing the name of the place, with an
// empty value for coords that no meaningful place could be computed.
func NewGoatReverseStep(inputs map[string]DataRef, outputMethods []string, lonCol int, latCol int) *StepData {
	args := map[string]interface{}{
		"lon_col_index":  lonCol,
		"lat_col_index":  latCol,
		"rampup_timeout": 150,
	}
	return NewStepData(
		&pipeline.Primitive{
			Id:         "f6e4880b-98c7-32f0-b687-a4b1d74c8f99",
			Version:    "1.0.8",
			Name:       "Goat_reverse",
			PythonPath: "d3m.primitives.data_cleaning.geocoding.Goat_reverse",
			Digest:     "d1bc0fa5f55ccdfc7d29b353c87e32e2e99b7c1d2cff4483978b8e2ed39af54a",
		},
		outputMethods,
		args,
		inputs,
	)
}

// NewJoinStep creates a step that will attempt to join two datasets a key column
// from each.  This is currently a placeholder for testing/debugging only.
func NewJoinStep(inputs map[string]DataRef, outputMethods []string, leftCol string, rightCol string, accuracy float32) *StepData {
	return NewStepData(
		&pipeline.Primitive{
			Id:         "6c3188bf-322d-4f9b-bb91-68151bf1f17f",
			Version:    "0.2.1",
			Name:       "Fuzzy Join Placeholder",
			PythonPath: "d3m.primitives.data_transformation.fuzzy_join.DistilFuzzyJoin",
			Digest:     "",
		},
		outputMethods,
		map[string]interface{}{"left_col": leftCol, "right_col": rightCol, "accuracy": accuracy},
		inputs,
	)
}

// NewDSBoxJoinStep creates a step that will attempt to join two datasets using
// key columns from each dataset.
func NewDSBoxJoinStep(inputs map[string]DataRef, outputMethods []string, leftCols []string, rightCols []string, accuracy float32) *StepData {
	joinType := "exact"
	if accuracy < 0.5 {
		joinType = "approximate"
	}
	return NewStepData(
		&pipeline.Primitive{
			Id:         "datamart-join",
			Version:    "1.4.4",
			Name:       "Datamart Augmentation",
			PythonPath: "d3m.primitives.data_augmentation.Join.DSBOX",
			Digest:     "",
		},
		outputMethods,
		map[string]interface{}{"left_col": leftCols, "right_col": rightCols, "join_type": joinType},
		inputs,
	)
}

// NewTimeseriesFormatterStep creates a step that will format a time series
// to the long form. The input dataset must be structured using resource
// files for time series data.  If mainResID is empty the primitive will attempt
// to infer the main resource.  If fileColIndex < 0, the file column will also
// be inferred.
func NewTimeseriesFormatterStep(inputs map[string]DataRef, outputMethods []string, mainResID string, fileColIndex int) *StepData {
	args := map[string]interface{}{}
	if mainResID != "" {
		args["main_resource_id"] = mainResID
	}
	if fileColIndex >= 0 {
		args["file_col_index"] = fileColIndex
	}
	return NewStepData(
		&pipeline.Primitive{
			Id:         "6a1ce3ee-ee70-428b-b1ff-0490bdb23023",
			Version:    "0.5.1",
			Name:       "Time series formatter",
			PythonPath: "d3m.primitives.data_transformation.time_series_formatter.DistilTimeSeriesFormatter",
			Digest:     "c23bb73a590dbd8aee56ac191e7fee64ff7772318d5892d8ee4646ede1f5b25d",
		},
		outputMethods,
		args,
		inputs,
	)
}

// NewImageRetrievalStep creates a step that will rank images based on nearnest
// to images with the positive label.
func NewImageRetrievalStep(inputs map[string]DataRef, outputMethods []string) *StepData {
	args := map[string]interface{}{
		"reduce_dimension": 1024,
		"gem_p":            4,
	}
	return NewStepData(
		&pipeline.Primitive{
			Id:         "6dd2032c-5558-4621-9bea-ea42403682da",
			Version:    "1.0.0",
			Name:       "ImageRetrieval",
			PythonPath: "d3m.primitives.similarity_modeling.iterative_labeling.ImageRetrieval",
			Digest:     "84d9b341135142c9bda23b7885516ebc9f552a87ea458712501ac85a29e0dc58",
		},
		outputMethods,
		args,
		inputs,
	)
}

// NewIsolationForestStep returns labels for whether or not a data point is an anomoly
func NewIsolationForestStep(inputs map[string]DataRef, outputMethods []string) *StepData {
	args := map[string]interface{}{}
	return NewStepData(
		&pipeline.Primitive{
			Id:         "793f0b17-7413-4962-9f1d-0b285540b21f",
			Version:    "0.5.1",
			Name:       "Isolation Forest",
			PythonPath: "d3m.primitives.classification.isolation_forest.IsolationForestPrimitive",
			Digest:     "",
		},
		outputMethods,
		args,
		inputs,
	)
}

// NewPrefeaturisedPoolingStep takes inputs of non-pooled remote sensing data to pool it
func NewPrefeaturisedPoolingStep(inputs map[string]DataRef, outputMethods []string) *StepData {
	args := map[string]interface{}{}
	return NewStepData(
		&pipeline.Primitive{
			Id:         "825ea1fb-90b2-442c-9905-efba48872102",
			Version:    "0.5.1",
			Name:       "Prefeaturised Pooler",
			PythonPath: "d3m.primitives.remote_sensing.remote_sensing_pretrained.PrefeaturisedPooler",
			Digest:     "",
		},
		outputMethods,
		args,
		inputs,
	)
}
