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
	"github.com/uncharted-distil/distil-compute/pipeline"
)

// NewSimonStep creates a SIMON data classification step.  It examines an input
// dataframe, and assigns types to the columns based on the exposed metadata.
func NewSimonStep() *StepData {
	return NewStepData(
		&pipeline.Primitive{
			Id:         "d2fa8df2-6517-3c26-bafc-87b701c4043a",
			Version:    "1.2.1",
			Name:       "simon",
			PythonPath: "d3m.primitives.data_cleaning.column_type_profiler.Simon",
			Digest:     "44fe5bf57ddb776440d1e22ddd1cf3ffeef9a282a3899856db8741e07fd7325d",
		},
		[]string{"produce_metafeatures"},
	)
}

// NewSlothStep creates a Sloth timeseries clustering step.
func NewSlothStep() *StepData {
	return NewStepDataWithHyperparameters(
		&pipeline.Primitive{
			Id:         "77bf4b92-2faa-3e38-bb7e-804131243a7f",
			Version:    "2.0.2",
			Name:       "Sloth",
			PythonPath: "d3m.primitives.time_series_segmentation.cluster.Sloth",
			Digest:     "576297f6bb41056ede966722bb0ed0d73403752e0a80eacd85bd71e8ea930e8a",
		},
		[]string{"produce"},
		map[string]interface{}{
			"nclusters": 4,
		},
	)
}

// NewUnicornStep creates a unicorn image clustering step.
func NewUnicornStep(targetColumns []string, outputLabels []string) *StepData {
	return NewStepDataWithHyperparameters(
		&pipeline.Primitive{
			Id:         "475c26dc-eb2e-43d3-acdb-159b80d9f099",
			Version:    "1.1.0",
			Name:       "unicorn",
			PythonPath: "d3m.primitives.digital_image_processing.unicorn.Unicorn",
			Digest:     "8c1280cb1355115d98de08e9981ea7cb95f6952885d5b190d9db789921664020",
		},
		[]string{"produce"},
		map[string]interface{}{
			"target_columns": targetColumns,
			"output_labels":  outputLabels,
		},
	)
}

// NewPCAFeaturesStep creates a PCA-based feature ranking call that can be added to
// a pipeline.
func NewPCAFeaturesStep() *StepData {
	return NewStepData(
		&pipeline.Primitive{
			Id:         "04573880-d64f-4791-8932-52b7c3877639",
			Version:    "3.0.1",
			Name:       "PCA Features",
			PythonPath: "d3m.primitives.feature_selection.pca_features.Pcafeatures",
			Digest:     "51ae6de10bbc004ed2e0e81fa8dcf8b6972c62cec4549c1a5cd58305e70eec71",
		},
		[]string{"produce"},
	)
}

// NewTargetRankingStep creates a target ranking call that can be added to
// a pipeline. Ranking is based on mutual information between features and a selected
// target.  Returns a DataFrame containing (col_idx, col_name, score) tuples for
// each ranked feature. Features that could not be ranked are excluded
// from the returned set.
func NewTargetRankingStep(targetCol int) *StepData {
	return NewStepDataWithHyperparameters(
		&pipeline.Primitive{
			Id:         "a31b0c26-cca8-4d54-95b9-886e23df8886",
			Version:    "0.1.0",
			Name:       "Mutual Information Feature Ranking",
			PythonPath: "d3m.primitives.distil.MIRanking",
			Digest:     "5302eebf2fb8a80e9f00e7b74888aba9eb448a9c0463d9d26786dab717a62c61",
		},
		[]string{"produce"},
		map[string]interface{}{
			"target_col_index": targetCol,
		},
	)
}

// NewDukeStep creates a wrapper for the Duke dataset classifier.
func NewDukeStep() *StepData {
	return NewStepData(
		&pipeline.Primitive{
			Id:         "46612a42-6120-3559-9db9-3aa9a76eb94f",
			Version:    "1.1.6",
			Name:       "duke",
			PythonPath: "d3m.primitives.data_cleaning.labler.Duke",
			Digest:     "b40cbf3631a19ef0141fb852079330c622b00ef286e54a755e6a90fc85be5963",
		},
		[]string{"produce"},
	)
}

// NewDataCleaningStep creates a wrapper for the Punk data cleaning primitive.
func NewDataCleaningStep() *StepData {
	return NewStepData(
		&pipeline.Primitive{
			Id:         "fc6bf33a-f3e0-3496-aa47-9a40289661bc",
			Version:    "3.0.1",
			Name:       "datacleaning",
			PythonPath: "d3m.primitives.data_cleaning.data_cleaning.Datacleaning",
			Digest:     "e4fe3196b81200106c40669d72a561f7fe1d7f36a9ddb5c0d7ce87bfb59f76fd",
		},
		[]string{"produce"},
	)
}

// NewCrocStep creates a wrapper for the Croc image classifier.
func NewCrocStep(targetColumns []string, outputLabels []string) *StepData {
	return NewStepDataWithHyperparameters(
		&pipeline.Primitive{
			Id:         "404fae2a-2f0a-4c9b-9ad2-fb1528990561",
			Version:    "1.2.3",
			Name:       "croc",
			PythonPath: "d3m.primitives.digital_image_processing.croc.Croc",
			Digest:     "a0cd922401d96b8ffbfe11f5db188b6a4d410119319392932e417b706ed5ae6",
		},
		[]string{"produce"},
		map[string]interface{}{
			"target_columns": targetColumns,
			"output_labels":  outputLabels,
		},
	)
}

// NewDatasetToDataframeStep creates a primitive call that transforms an input dataset
// into a PANDAS dataframe.
func NewDatasetToDataframeStep() *StepData {
	return NewStepData(
		&pipeline.Primitive{
			Id:         "4b42ce1e-9b98-4a25-b68e-fad13311eb65",
			Version:    "0.3.0",
			Name:       "Dataset to DataFrame converter",
			PythonPath: "d3m.primitives.data_transformation.dataset_to_dataframe.Common",
			Digest:     "85b946aa6123354fe51a288c3be56aaca82e76d4071c1edc13be6f9e0e100144",
		},
		[]string{"produce"},
	)
}

// ColumnUpdate defines a set of column indices to add/remvoe
// a set of semantic types to/from.
type ColumnUpdate struct {
	Indices       []int
	SemanticTypes []string
}

// NewUpdateSemanticTypeStep adds and removes semantic data values from an input
// dataset.  An add of (1, 2), ("type a", "type b") would result in "type a" and "type b"
// being added to index 1 and 2.  Leaving the resource ID
// as the empty value allows the primitive to infer the main resource from the
// the dataset.
func NewUpdateSemanticTypeStep(resourceID string, add *ColumnUpdate, remove *ColumnUpdate) (*StepData, error) {
	hyperparams := map[string]interface{}{
		"add_columns":    add.Indices,
		"add_types":      add.SemanticTypes,
		"remove_columns": remove.Indices,
		"remove_types":   remove.SemanticTypes,
	}
	if resourceID != "" {
		hyperparams["resource_id"] = resourceID
	}

	return NewStepDataWithHyperparameters(
		&pipeline.Primitive{
			Id:         "98c79128-555a-4a6b-85fb-d4f4064c94ab",
			Version:    "0.2.0",
			Name:       "Semantic type updater",
			PythonPath: "d3m.primitives.data_transformation.update_semantic_types.DatasetCommon",
			Digest:     "85b946aa6123354fe51a288c3be56aaca82e76d4071c1edc13be6f9e0e100144",
		},
		[]string{"produce"},
		hyperparams,
	), nil
}

// NewDenormalizeStep denormalize data that is contained in multiple resource files.
func NewDenormalizeStep() *StepData {
	return NewStepData(
		&pipeline.Primitive{
			Id:         "f31f8c1f-d1c5-43e5-a4b2-2ae4a761ef2e",
			Version:    "0.2.0",
			Name:       "Denormalize datasets",
			PythonPath: "d3m.primitives.data_transformation.denormalize.Common",
			Digest:     "c39e3436373aed1944edbbc9b1cf24af5c71919d73bf0bb545cba0b685812df1",
		},
		[]string{"produce"},
	)
}

// NewColumnParserStep takes obj/string columns in a dataframe and parses them into their
// associated raw python types based on the attached d3m metadata.
func NewColumnParserStep() *StepData {
	return NewStepData(
		&pipeline.Primitive{
			Id:         "d510cb7a-1782-4f51-b44c-58f0236e47c7",
			Version:    "0.4.0",
			Name:       "Parses strings into their types",
			PythonPath: "d3m.primitives.data_transformation.column_parser.DataFrameCommon",
			Digest:     "",
		},
		[]string{"produce"},
	)
}

// NewRemoveColumnsStep removes columns from an input dataframe.  Columns
// are specified by name and the match is case insensitive.  Leaving the resource ID
// as the empty value allows the primitive to infer the main resource from the
// the dataset.
func NewRemoveColumnsStep(resourceID string, colIndices []int) (*StepData, error) {
	hyperparams := map[string]interface{}{
		"columns": colIndices,
	}
	if resourceID != "" {
		hyperparams["resource_id"] = resourceID
	}

	return NewStepDataWithHyperparameters(
		&pipeline.Primitive{
			Id:         "2eeff053-395a-497d-88db-7374c27812e6",
			Version:    "0.2.0",
			Name:       "Column remover",
			PythonPath: "d3m.primitives.data_transformation.remove_columns.DatasetCommon",
			Digest:     "85b946aa6123354fe51a288c3be56aaca82e76d4071c1edc13be6f9e0e100144",
		},
		[]string{"produce"},
		hyperparams,
	), nil
}

// NewTermFilterStep creates a primitive step that filters dataset rows based on a match against a
// term list.  The term match can be partial, or apply to whole terms only.  Leaving the resource ID
// as the empty value allows the primitive to infer the main resource from the
// the dataset.
func NewTermFilterStep(resourceID string, colindex int, inclusive bool, terms []string, matchWhole bool) *StepData {
	hyperparams := map[string]interface{}{
		"column":      colindex,
		"inclusive":   inclusive,
		"terms":       terms,
		"match_whole": matchWhole,
	}
	if resourceID != "" {
		hyperparams["resource_id"] = resourceID
	}

	return NewStepDataWithHyperparameters(
		&pipeline.Primitive{
			Id:         "622893c7-42fc-4561-a6f6-071fb85d610a",
			Version:    "0.1.0",
			Name:       "Term list dataset filter",
			PythonPath: "d3m.primitives.datasets.TermFilter",
			Digest:     "",
		},
		[]string{"produce"},
		hyperparams,
	)
}

// NewRegexFilterStep creates a primitive step that filter dataset rows based on a regex match.
// Leaving the resource ID as the empty value allows the primitive to infer the main resource from the
// the dataset.
func NewRegexFilterStep(resourceID string, colindex int, inclusive bool, regex string) *StepData {
	hyperparams := map[string]interface{}{
		"column":    colindex,
		"inclusive": inclusive,
		"regex":     regex,
	}
	if resourceID != "" {
		hyperparams["resource_id"] = resourceID
	}

	return NewStepDataWithHyperparameters(
		&pipeline.Primitive{
			Id:         "d1b4c4b7-63ba-4ee6-ab30-035157cccf22",
			Version:    "0.1.0",
			Name:       "Regex dataset filter",
			PythonPath: "d3m.primitives.datasets.RegexFilter",
			Digest:     "",
		},
		[]string{"produce"},
		hyperparams,
	)
}

// NewNumericRangeFilterStep creates a primitive step that filters dataset rows based on an
// included/excluded numeric range.  Inclusion of boundaries is controlled by the strict flag.
// Leaving the resource ID as the empty value allows the primitive to infer the main resource from the
// the dataset.
func NewNumericRangeFilterStep(resourceID string, colindex int, inclusive bool, min float64, max float64, strict bool) *StepData {
	hyperparams := map[string]interface{}{
		"column":    colindex,
		"inclusive": inclusive,
		"min":       min,
		"max":       max,
		"strict":    strict,
	}
	if resourceID != "" {
		hyperparams["resource_id"] = resourceID
	}

	return NewStepDataWithHyperparameters(
		&pipeline.Primitive{
			Id:         "8b1c1140-8c21-4f41-aeca-662b7d35aa29",
			Version:    "0.1.0",
			Name:       "Numeric range filter",
			PythonPath: "d3m.primitives.datasets.NumericRangeFilter",
			Digest:     "",
		},
		[]string{"produce"},
		hyperparams,
	)
}

// NewTimeSeriesLoaderStep creates a primitive step that reads time series values using a dataframe
// containing a file URI column.  The file URIs are expected to point to CSV files, with the
// supplied time and value indices pointing the columns in the CSV that form the series data.
// The result is a new dataframe that stores the timetamps as the column headers,
// and the accompanying values for each file as a row.  Note that the file index column is negative,
// the primitive will use the first CSV file name column if finds.
func NewTimeSeriesLoaderStep(fileColIndex int, timeColIndex int, valueColIndex int) *StepData {
	// exclude the file col index val ue in the case of a negative index so that the
	// primitive will infer the colum
	args := map[string]interface{}{
		"time_col_index":  timeColIndex,
		"value_col_index": valueColIndex,
	}
	if fileColIndex >= 0 {
		args["file_col_index"] = fileColIndex
	}

	return NewStepDataWithHyperparameters(
		&pipeline.Primitive{
			Id:         "1689aafa-16dc-4c55-8ad4-76cadcf46086",
			Version:    "0.1.0",
			Name:       "Time series loader",
			PythonPath: "d3m.primitives.distil.TimeSeriesLoader",
			Digest:     "",
		},
		[]string{"produce"},
		args,
	)
}

// NewGoatForwardStep creates a GOAT forward geocoding primitive.  A string column
// containing a place name or address is passed in, and the primitive will
// return a DataFrame containing the lat/lon coords of the place.  If location could
// not be found, the row in the data frame will be empty.
func NewGoatForwardStep(placeCol string) *StepData {
	args := map[string]interface{}{
		"target_columns": []string{placeCol},
	}
	return NewStepDataWithHyperparameters(
		&pipeline.Primitive{
			Id:         "c7c61da3-cf57-354e-8841-664853370106",
			Version:    "1.0.7",
			Name:       "Goat_forward",
			PythonPath: "d3m.primitives.data_cleaning.geocoding.Goat_forward",
			Digest:     "655c3b536ee2b87ec4607ba932650a0655400880de89bba2effee4a7f17df9f8",
		},
		[]string{"produce"},
		args,
	)
}

// NewGoatReverseStep creates a GOAT reverse geocoding primitive.  Columns
// containing lat and lon values are passed in, and the primitive will
// return a DataFrame containing the name of the place, with an
// empty value for coords that no meaningful place could be computed.
func NewGoatReverseStep(lonCol string, latCol string) *StepData {
	args := map[string]interface{}{
		"lon_col_index": lonCol,
		"lat_col_index": latCol,
	}
	return NewStepDataWithHyperparameters(
		&pipeline.Primitive{
			Id:         "f6e4880b-98c7-32f0-b687-a4b1d74c8f99",
			Version:    "1.0.7",
			Name:       "Goat_reverse",
			PythonPath: "d3m.primitives.data_cleaning.geocoding.Goat_reverse",
			Digest:     "2111b6253ac8b3765ccdc1d42b76bf34258b90ef824113d227e1b89a090259b9",
		},
		[]string{"produce"},
		args,
	)
}

// NewJoinStep creates a step that will attempt to join two datasets a key column
// from each.  This is currently a placeholder for testing/debugging only.
func NewJoinStep(leftCol string, rightCol string, accuracy float32) *StepData {
	return NewStepDataWithAll(
		&pipeline.Primitive{
			Id:         "6c3188bf-322d-4f9b-bb91-68151bf1f17f",
			Version:    "0.1.0",
			Name:       "Fuzzy Join Placeholder",
			PythonPath: "d3m.primitives.distil.FuzzyJoin",
			Digest:     "",
		},
		[]string{"produce"},
		map[string]interface{}{"left_col": leftCol, "right_col": rightCol, "accuracy": accuracy},
		[]string{"left", "right"},
	)
}

// NewDSBoxJoinStep creates a step that will attempt to join two datasets using
// key columns from each dataset.
func NewDSBoxJoinStep(leftCols []string, rightCols []string, accuracy float32) *StepData {
	joinType := "exact"
	if accuracy < 0.5 {
		joinType = "approximate"
	}
	return NewStepDataWithHyperparameters(
		&pipeline.Primitive{
			Id:         "datamart-join",
			Version:    "1.4.4",
			Name:       "Datamart Augmentation",
			PythonPath: "d3m.primitives.data_augmentation.Join.DSBOX",
			Digest:     "",
		},
		[]string{"produce"},
		map[string]interface{}{"left_col": leftCols, "right_col": rightCols, "join_type": joinType},
	)
}

// NewTimeseriesFormatterStep creates a step that will format a time series
// to the long form. The input dataset must be structured using resource
// files for time series data.
func NewTimeseriesFormatterStep(mainResID string, fileColIndex int) *StepData {
	args := map[string]interface{}{
		"main_resource_index": mainResID,
	}
	if fileColIndex >= 0 {
		args["file_col_index"] = fileColIndex
	}
	return NewStepDataWithHyperparameters(
		&pipeline.Primitive{
			Id:         "1c4aed23-f3d3-4e6b-9710-009a9bc9b694",
			Version:    "0.1.0",
			Name:       "Time series formatter",
			PythonPath: "d3m.primitives.distil.TimeSeriesFormatter",
			Digest:     "",
		},
		[]string{"produce"},
		args,
	)
}
