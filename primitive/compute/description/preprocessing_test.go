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
	"io/ioutil"
	"testing"

	"github.com/gogo/protobuf/proto"
	"github.com/stretchr/testify/assert"
	"github.com/uncharted-distil/distil-compute/model"
)

const (
	searchResult = `
	{
		"id": "datamart.url.a3943fd7892d5d219012f889327c6661",
		"score": 12.832686,
		"metadata":
		{
			"name": "Newyork Weather Data around Airport 2016-18",
			"description": "This data contains weather information for NY city around LaGuardia Airport from 2016 to 2018; we...",
			"size": 1523693,
			"nb_rows": 24624,
			"columns":
			[
				{
					"name": "DATE",
					"structural_type": "http://schema.org/Text",
					"semantic_types": ["http://schema.org/DateTime"],
					"mean": 1495931400.0,
					"stddev": 25590011.431395352,
					"coverage":
					[
						{
							"range":
							{
								"gte": 1482850800.0,
								"lte": 1509444000.0
							}
						},
						{
							"range":
							{
								"gte": 1453096800.0,
								"lte": 1479884400.0
							}
						},
						{
							"range":
							{
								"gte": 1512388800.0,
								"lte": 1538787600.0
							}
						}
					]
				},
				{
					"name": "HOURLYSKYCONDITIONS",
					"structural_type":
					"http://schema.org/Text",
					"semantic_types": []
				},
				{
					"name": "HOURLYDRYBULBTEMPC",
					"structural_type": "http://schema.org/Float",
					"semantic_types": [],
					"mean": 14.666224009096823,
					"stddev": 9.973788193915643,
					"coverage":
					[
						{
							"range":
							{
								"gte": 9.0,
								"lte": 19.0
							}
						},
						{
							"range":
							{
								"gte": -6.1,
								"lte": 8.0
							}
						},
						{
							"range":
							{
								"gte": 20.6,
								"lte": 31.7
							}
						}
					]
				},
				{
					"name": "HOURLYRelativeHumidity",
					"structural_type": "http://schema.org/Float",
					"semantic_types": [],
					"mean": 60.70849577647823,
					"stddev": 18.42048051096981,
					"coverage":
					[
						{
							"range":
							{
								"gte": 50.0,
								"lte": 70.0
							}
						},
						{
							"range":
							{
								"gte": 26.0,
								"lte": 49.0
							}
						},
						{
							"range":
							{
								"gte": 73.0,
								"lte": 96.0
							}
						}
					]
				},
				{
					"name": "HOURLYWindSpeed",
					"structural_type": "http://schema.org/Float",
					"semantic_types": [],
					"mean": 10.68859649122807,
					"stddev": 5.539675475162907,
					"coverage":
					[
						{
							"range":
							{
								"gte": 0.0,
								"lte": 8.0
							}
						},
						{
							"range":
							{
								"gte": 16.0,
								"lte": 28.0
							}
						},
						{
							"range":
							{
								"gte": 9.0,
								"lte": 15.0
							}
						}
					]
				},
				{
					"name": "HOURLYWindDirection",
					"structural_type": "http://schema.org/Text",
					"semantic_types": []
				},
				{
					"name": "HOURLYStationPressure",
					"structural_type": "http://schema.org/Float",
					"semantic_types": ["https://metadata.datadrivendiscovery.org/types/PhoneNumber"],
					"mean": 29.90760315139694,
					"stddev": 0.24584097919742368,
					"coverage":
					[
						{
							"range":
							{
								"gte": 29.86,
								"lte": 30.12
							}
						},
						{
							"range":
							{
								"gte": 30.14,
								"lte": 30.55
							}
						},
						{
							"range":
							{
								"gte": 29.42,
								"lte": 29.84
							}
						}
					]
				}
			],
			"materialize":
			{
				 "direct_url": "https://drive.google.com/uc?export=download&id=1jRwzZwEGMICE3n6-nwmVxMD2c0QCHad4",
				 "identifier": "datamart.url"
			},
			"date": "2019-07-02T15:38:00.413962Z"},
			"augmentation":
			{
				"type": "none",
				"left_columns": [],
				"right_columns": []
			}
		}`
)

func TestCreateUserDatasetPipeline(t *testing.T) {

	variables := []*model.Variable{
		{
			Key:          "test_var_0",
			OriginalType: "ordinal",
			Type:         "categorical",
			Index:        0,
			DistilRole:   []string{model.VarDistilRoleData},
		},
		{
			Key:          "test_var_1",
			OriginalType: "categorical",
			Type:         "integer",
			Index:        1,
			DistilRole:   []string{model.VarDistilRoleData},
		},
		{
			Key:          "test_var_2",
			OriginalType: "categorical",
			Type:         "integer",
			Index:        2,
			DistilRole:   []string{model.VarDistilRoleData},
		},
		{
			Key:          "test_var_3",
			OriginalType: "categorical",
			Type:         "integer",
			Index:        3,
			DistilRole:   []string{model.VarDistilRoleData},
		},
	}

	pipeline, err := CreateUserDatasetPipeline("test_user_pipeline", "a test user pipeline",
		&UserDatasetDescription{
			AllFeatures:      variables,
			TargetFeature:    &model.Variable{Key: "test_target"},
			SelectedFeatures: []string{"test_var_0", "test_var_1", "test_var_3"},
			Filters:          nil,
		},
		nil,
	)
	assert.Equal(t, 16, len(pipeline.GetSteps()))

	pythonPath := pipeline.GetSteps()[0].GetPrimitive().GetPrimitive().GetPythonPath()
	assert.Equal(t, "d3m.primitives.data_transformation.extract_columns.Common", pythonPath)

	pythonPath = pipeline.GetSteps()[1].GetPrimitive().GetPrimitive().GetPythonPath()
	assert.Equal(t, "d3m.primitives.operator.dataset_map.DataFrameCommon", pythonPath)
	assert.Equal(t, int32(0), pipeline.GetSteps()[0].GetPrimitive().GetHyperparams()["primitive"].GetPrimitive().GetData())
	inputs := pipeline.GetSteps()[1].GetPrimitive().GetArguments()["inputs"].GetContainer().GetData()
	assert.Equal(t, "inputs.0", inputs)

	pythonPath = pipeline.GetSteps()[2].GetPrimitive().GetPrimitive().GetPythonPath()
	assert.Equal(t, "d3m.primitives.data_transformation.denormalize.Common", pythonPath)
	inputs = pipeline.GetSteps()[2].GetPrimitive().GetArguments()["inputs"].GetContainer().GetData()
	assert.Equal(t, "steps.1.produce", inputs)

	pythonPath = pipeline.GetSteps()[3].GetPrimitive().GetPrimitive().GetPythonPath()
	assert.Equal(t, "d3m.primitives.data_transformation.column_parser.DistilColumnParser", pythonPath)

	pythonPath = pipeline.GetSteps()[4].GetPrimitive().GetPrimitive().GetPythonPath()
	assert.Equal(t, "d3m.primitives.operator.dataset_map.DataFrameCommon", pythonPath)
	assert.Equal(t, int32(3), pipeline.GetSteps()[4].GetPrimitive().GetHyperparams()["primitive"].GetPrimitive().GetData())
	inputs = pipeline.GetSteps()[4].GetPrimitive().GetArguments()["inputs"].GetContainer().GetData()
	assert.Equal(t, "steps.2.produce", inputs)

	pythonPath = pipeline.GetSteps()[5].GetPrimitive().GetPrimitive().GetPythonPath()
	assert.Equal(t, "d3m.primitives.data_cleaning.data_cleaning.Datacleaning", pythonPath)

	pythonPath = pipeline.GetSteps()[6].GetPrimitive().GetPrimitive().GetPythonPath()
	assert.Equal(t, "d3m.primitives.operator.dataset_map.DataFrameCommon", pythonPath)
	assert.Equal(t, int32(5), pipeline.GetSteps()[6].GetPrimitive().GetHyperparams()["primitive"].GetPrimitive().GetData())
	inputs = pipeline.GetSteps()[6].GetPrimitive().GetArguments()["inputs"].GetContainer().GetData()
	assert.Equal(t, "steps.4.produce", inputs)

	// add semantic type integer to cols 1,3
	hyperParams := pipeline.GetSteps()[7].GetPrimitive().GetHyperparams()
	assert.Equal(t, []int64{1, 3}, ConvertToIntArray(hyperParams["columns"].GetValue().GetData().GetRaw().GetList()))
	assert.Equal(t, []string{"http://schema.org/Integer"}, ConvertToStringArray(hyperParams["semantic_types"].GetValue().GetData().GetRaw().GetList()))

	assert.Equal(t, int32(7), pipeline.GetSteps()[8].GetPrimitive().GetHyperparams()["primitive"].GetPrimitive().GetData())
	inputs = pipeline.GetSteps()[8].GetPrimitive().GetArguments()["inputs"].GetContainer().GetData()
	assert.Equal(t, "steps.6.produce", inputs)

	// remove semantic type categorical from cols 1,3
	hyperParams = pipeline.GetSteps()[9].GetPrimitive().GetHyperparams()
	assert.Equal(t, []int64{1, 3}, ConvertToIntArray(hyperParams["columns"].GetValue().GetData().GetRaw().GetList()))
	assert.Equal(t, []string{"https://metadata.datadrivendiscovery.org/types/CategoricalData"},
		ConvertToStringArray(hyperParams["semantic_types"].GetValue().GetData().GetRaw().GetList()))

	assert.Equal(t, int32(9), pipeline.GetSteps()[10].GetPrimitive().GetHyperparams()["primitive"].GetPrimitive().GetData())
	inputs = pipeline.GetSteps()[10].GetPrimitive().GetArguments()["inputs"].GetContainer().GetData()
	assert.Equal(t, "steps.8.produce", inputs)

	// add attrobute semantic type
	hyperParams = pipeline.GetSteps()[11].GetPrimitive().GetHyperparams()
	assert.Equal(t, []int64{0, 1, 3}, ConvertToIntArray(hyperParams["columns"].GetValue().GetData().GetRaw().GetList()))
	assert.Equal(t, []string{"https://metadata.datadrivendiscovery.org/types/Attribute"}, ConvertToStringArray(hyperParams["semantic_types"].GetValue().GetData().GetRaw().GetList()))

	assert.Equal(t, int32(11), pipeline.GetSteps()[12].GetPrimitive().GetHyperparams()["primitive"].GetPrimitive().GetData())
	inputs = pipeline.GetSteps()[12].GetPrimitive().GetArguments()["inputs"].GetContainer().GetData()
	assert.Equal(t, "steps.10.produce", inputs)

	// remove column from index two
	hyperParams = pipeline.GetSteps()[13].GetPrimitive().GetHyperparams()
	assert.Equal(t, []int64{2}, ConvertToIntArray(hyperParams["columns"].GetValue().GetData().GetRaw().GetList()))

	assert.Equal(t, int32(0), pipeline.GetSteps()[13].GetPrimitive().GetHyperparams()["primitive"].GetPrimitive().GetData())
	inputs = pipeline.GetSteps()[14].GetPrimitive().GetArguments()["inputs"].GetContainer().GetData()
	assert.Equal(t, "steps.12.produce", inputs)

	// next is the inference step, which doesn't have a primitive associated with it
	assert.NotNil(t, pipeline.GetSteps()[15].GetPlaceholder())
	inputs = pipeline.GetSteps()[15].GetPlaceholder().GetInputs()[0].GetData()
	assert.Equal(t, "steps.14.produce", inputs)

	assert.NoError(t, err)
}

func TestCreateUserDatasetPipelineMappingError(t *testing.T) {

	variables := []*model.Variable{
		{
			Key:          "test_var_0",
			OriginalType: "blordinal",
			Type:         "categorical",
			Index:        0,
			DistilRole:   []string{"data"},
		},
	}

	_, err := CreateUserDatasetPipeline("test_user_pipeline", "a test user pipeline",
		&UserDatasetDescription{
			AllFeatures:      variables,
			TargetFeature:    &model.Variable{Key: "test_target"},
			SelectedFeatures: []string{"test_var_0"},
			Filters:          nil,
		}, nil)
	assert.Error(t, err)
}

func TestCreateUserDatasetEmpty(t *testing.T) {

	variables := []*model.Variable{
		{
			Key:          "test_var_0",
			OriginalType: "categorical",
			Type:         "categorical",
			Index:        0,
			DistilRole:   []string{"data"},
		},
	}

	pipeline, err := CreateUserDatasetPipeline("test_user_pipeline", "a test user pipeline",
		&UserDatasetDescription{
			AllFeatures:      variables,
			TargetFeature:    &model.Variable{Key: "test_target"},
			SelectedFeatures: []string{"test_var_0"},
			Filters:          nil,
		}, nil)
	assert.Equal(t, 10, len(pipeline.GetSteps()))
	assert.Nil(t, err)

	pythonPath := pipeline.GetSteps()[0].GetPrimitive().GetPrimitive().GetPythonPath()
	assert.Equal(t, "d3m.primitives.data_transformation.extract_columns.Common", pythonPath)

	pythonPath = pipeline.GetSteps()[1].GetPrimitive().GetPrimitive().GetPythonPath()
	assert.Equal(t, "d3m.primitives.operator.dataset_map.DataFrameCommon", pythonPath)
	assert.Equal(t, int32(0), pipeline.GetSteps()[0].GetPrimitive().GetHyperparams()["primitive"].GetPrimitive().GetData())
	inputs := pipeline.GetSteps()[1].GetPrimitive().GetArguments()["inputs"].GetContainer().GetData()
	assert.Equal(t, "inputs.0", inputs)

	pythonPath = pipeline.GetSteps()[2].GetPrimitive().GetPrimitive().GetPythonPath()
	assert.Equal(t, "d3m.primitives.data_transformation.denormalize.Common", pythonPath)
	inputs = pipeline.GetSteps()[2].GetPrimitive().GetArguments()["inputs"].GetContainer().GetData()
	assert.Equal(t, "steps.1.produce", inputs)

	pythonPath = pipeline.GetSteps()[3].GetPrimitive().GetPrimitive().GetPythonPath()
	assert.Equal(t, "d3m.primitives.data_transformation.column_parser.DistilColumnParser", pythonPath)

	pythonPath = pipeline.GetSteps()[4].GetPrimitive().GetPrimitive().GetPythonPath()
	assert.Equal(t, "d3m.primitives.operator.dataset_map.DataFrameCommon", pythonPath)
	assert.Equal(t, int32(3), pipeline.GetSteps()[4].GetPrimitive().GetHyperparams()["primitive"].GetPrimitive().GetData())
	inputs = pipeline.GetSteps()[4].GetPrimitive().GetArguments()["inputs"].GetContainer().GetData()
	assert.Equal(t, "steps.2.produce", inputs)

	pythonPath = pipeline.GetSteps()[5].GetPrimitive().GetPrimitive().GetPythonPath()
	assert.Equal(t, "d3m.primitives.data_cleaning.data_cleaning.Datacleaning", pythonPath)

	pythonPath = pipeline.GetSteps()[6].GetPrimitive().GetPrimitive().GetPythonPath()
	assert.Equal(t, "d3m.primitives.operator.dataset_map.DataFrameCommon", pythonPath)
	assert.Equal(t, int32(5), pipeline.GetSteps()[6].GetPrimitive().GetHyperparams()["primitive"].GetPrimitive().GetData())
	inputs = pipeline.GetSteps()[6].GetPrimitive().GetArguments()["inputs"].GetContainer().GetData()
	assert.Equal(t, "steps.4.produce", inputs)

	// add attrobute semantic type
	hyperParams := pipeline.GetSteps()[7].GetPrimitive().GetHyperparams()
	assert.Equal(t, []int64{0}, ConvertToIntArray(hyperParams["columns"].GetValue().GetData().GetRaw().GetList()))
	assert.Equal(t, []string{"https://metadata.datadrivendiscovery.org/types/Attribute"}, ConvertToStringArray(hyperParams["semantic_types"].GetValue().GetData().GetRaw().GetList()))

	assert.Equal(t, int32(7), pipeline.GetSteps()[8].GetPrimitive().GetHyperparams()["primitive"].GetPrimitive().GetData())
	inputs = pipeline.GetSteps()[8].GetPrimitive().GetArguments()["inputs"].GetContainer().GetData()
	assert.Equal(t, "steps.6.produce", inputs)

	// next is the inference step, which doesn't have a primitive associated with it
	assert.NotNil(t, pipeline.GetSteps()[9].GetPlaceholder())
	inputs = pipeline.GetSteps()[9].GetPlaceholder().GetInputs()[0].GetData()
	assert.Equal(t, "steps.8.produce", inputs)
}

func TestCreatePCAFeaturesPipeline(t *testing.T) {
	pipeline, err := CreatePCAFeaturesPipeline("pca_features_test", "test pca feature ranking pipeline")
	assert.NoError(t, err)

	data, err := proto.Marshal(pipeline.Pipeline)
	assert.NoError(t, err)
	assert.NotNil(t, data)

	err = ioutil.WriteFile("/tmp/pca_features.pln", data, 0644)
	assert.NoError(t, err)
}

func TestCreateSimonPipeline(t *testing.T) {
	pipeline, err := CreateSimonPipeline("simon_test", "test simon classification pipeline")
	assert.NoError(t, err)

	data, err := proto.Marshal(pipeline.Pipeline)
	assert.NoError(t, err)
	assert.NotNil(t, data)

	err = ioutil.WriteFile("/tmp/simon.pln", data, 0644)
	assert.NoError(t, err)
}

func TestCreateDataCleaningPipeline(t *testing.T) {
	pipeline, err := CreateDataCleaningPipeline("data cleaning test", "test data cleaning pipeline", nil, true)
	assert.NoError(t, err)

	data, err := proto.Marshal(pipeline.Pipeline)
	assert.NoError(t, err)
	assert.NotNil(t, data)

	err = ioutil.WriteFile("/tmp/datacleaning.pln", data, 0644)
	assert.NoError(t, err)
}

func TestCreateSlothPipeline(t *testing.T) {
	timeSeriesVariables := []*model.Variable{
		{Key: "time", Type: "string", OriginalType: "unknown", Index: 0},
		{Key: "value", Type: "string", OriginalType: "unknown", Index: 1},
	}

	pipeline, err := CreateSlothPipeline("sloth_test", "test sloth object detection pipeline", "time", "value", nil, timeSeriesVariables)
	assert.NoError(t, err)

	data, err := proto.Marshal(pipeline.Pipeline)
	assert.NoError(t, err)
	assert.NotNil(t, data)

	err = ioutil.WriteFile("/tmp/sloth.pln", data, 0644)
	assert.NoError(t, err)
}

func TestCreateDukePipeline(t *testing.T) {
	pipeline, err := CreateDukePipeline("duke_test", "test duke data summary pipeline")
	assert.NoError(t, err)

	data, err := proto.Marshal(pipeline.Pipeline)
	assert.NoError(t, err)
	assert.NotNil(t, data)

	err = ioutil.WriteFile("/tmp/duke.pln", data, 0644)
	assert.NoError(t, err)
}

func TestCreateTargetRankingPipeline(t *testing.T) {
	vars := []*model.Variable{
		{
			Key:          "hall_of_fame",
			Index:        18,
			Type:         model.CategoricalType,
			OriginalType: model.CategoricalType,
		},
	}
	pipeline, err := CreateTargetRankingPipeline("target_ranking_test", "test target_ranking pipeline", vars[0], vars, map[string]bool{"hall_of_fame": true})
	assert.NoError(t, err)

	data, err := proto.Marshal(pipeline.Pipeline)
	assert.NoError(t, err)
	assert.NotNil(t, data)

	err = ioutil.WriteFile("/tmp/target_ranking.pln", data, 0644)
	assert.NoError(t, err)
}

func TestCreateGoatForwardPipeline(t *testing.T) {
	region := &model.Variable{
		Key:   "region",
		Index: 14,
	}
	pipeline, err := CreateGoatForwardPipeline("goat_forward_test", "test goat forward geocoding pipeline", region)
	assert.NoError(t, err)

	data, err := proto.Marshal(pipeline.Pipeline)
	assert.NoError(t, err)
	assert.NotNil(t, data)

	err = ioutil.WriteFile("/tmp/goat_forward.pln", data, 0644)
	assert.NoError(t, err)
}

func TestCreateGoatReversePipeline(t *testing.T) {
	lat := &model.Variable{
		Key:   "lat",
		Index: 20,
	}
	lon := &model.Variable{
		Key:   "lon",
		Index: 21,
	}
	pipeline, err := CreateGoatReversePipeline("goat_reverse_test", "test goat reverse geocoding pipeline", lat, lon)
	assert.NoError(t, err)

	data, err := proto.Marshal(pipeline.Pipeline)
	assert.NoError(t, err)
	assert.NotNil(t, data)

	err = ioutil.WriteFile("/tmp/goat_reverse.pln", data, 0644)
	assert.NoError(t, err)
}

func TestCreateJoinPipeline(t *testing.T) {
	joins := []*Join{{
		Left:     &model.Variable{HeaderName: "Doubles", Type: "real", OriginalType: "string"},
		Right:    &model.Variable{HeaderName: "horsepower", Type: "real", OriginalType: "string"},
		Accuracy: 0.8,
	},
	}
	joinInfo := &JoinDescription{
		Joins: joins,
		Type:  "left",
	}
	pipeline, err := CreateJoinPipeline("join_test", "test join pipeline", joinInfo)
	assert.NoError(t, err)

	data, err := proto.Marshal(pipeline.Pipeline)
	assert.NoError(t, err)
	assert.NotNil(t, data)

	err = ioutil.WriteFile("/tmp/join.pln", data, 0644)
	assert.NoError(t, err)
}

func TestCreateDSBoxJoinPipeline(t *testing.T) {
	pipeline, err := CreateDSBoxJoinPipeline("ds_join_test", "test ds box join pipeline", []string{"Doubles"}, []string{"horsepower"}, 0.8)
	assert.NoError(t, err)

	data, err := proto.Marshal(pipeline.Pipeline)
	assert.NoError(t, err)
	assert.NotNil(t, data)

	err = ioutil.WriteFile("/tmp/ds_join.pln", data, 0644)
	assert.NoError(t, err)
}

func TestCreateDenormalizePipeline(t *testing.T) {
	pipeline, err := CreateDenormalizePipeline("denorm_test", "test denorm pipeline")
	assert.NoError(t, err)

	data, err := proto.Marshal(pipeline.Pipeline)
	assert.NoError(t, err)
	assert.NotNil(t, data)

	err = ioutil.WriteFile("/tmp/denorm.pln", data, 0644)
	assert.NoError(t, err)
}

func TestCreateTimeseriesFormatterPipeline(t *testing.T) {
	pipeline, err := CreateTimeseriesFormatterPipeline("formatter_test", "test formatter pipeline", "0")
	assert.NoError(t, err)

	data, err := proto.Marshal(pipeline.Pipeline)
	assert.NoError(t, err)
	assert.NotNil(t, data)

	err = ioutil.WriteFile("/tmp/formatter.pln", data, 0644)
	assert.NoError(t, err)
}

func TestCreateDatamartDownloadPipeline(t *testing.T) {
	pipeline, err := CreateDatamartDownloadPipeline("download_test", "test download pipeline", searchResult, "NYU")
	assert.NoError(t, err)

	data, err := proto.Marshal(pipeline.Pipeline)
	assert.NoError(t, err)
	assert.NotNil(t, data)

	err = ioutil.WriteFile("/tmp/download.pln", data, 0644)
	assert.NoError(t, err)
}

func TestCreateDatamartAugmentPipeline(t *testing.T) {
	pipeline, err := CreateDatamartAugmentPipeline("augment_test", "test augment pipeline", searchResult, "NYU")
	assert.NoError(t, err)

	data, err := proto.Marshal(pipeline.Pipeline)
	assert.NoError(t, err)
	assert.NotNil(t, data)

	err = ioutil.WriteFile("/tmp/augment.pln", data, 0644)
	assert.NoError(t, err)
}
