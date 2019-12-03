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

package compute

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"

	"github.com/golang/protobuf/proto"
	protobuf "github.com/golang/protobuf/protoc-gen-go/descriptor"
	log "github.com/unchartedsoftware/plog"

	"github.com/uncharted-distil/distil-compute/pipeline"
)

const (
	unknownAPIVersion = "unknown"

	// Task and sub-task types taken from the D3M LL problem schema - these are what is used internally within the
	// application at the server and client level to capture the model building task type.
	// For communication with TA2, these need to be translated via the `ConvertXXXFromTA3ToTA2` methods below.

	// ForecastingTask represents timeseries forcasting
	ForecastingTask = "forecasting"
	// ClassificationTask represents a classification task on image, timeseries or basic tabular data
	ClassificationTask = "classification"
	// RegressionTask represents a regression task on image, timeseries or basic tabular data
	RegressionTask = "regression"
	// ClusteringTask represents an unsupervised clustering task on image, timeseries or basic tabular data
	ClusteringTask = "clustering"
	// LinkPredictionTask represents a link prediction task on graph data
	LinkPredictionTask = "linkPrediction"
	// VertexClassificationTask represents a vertex nomination task on graph data
	VertexClassificationTask = "vertexClassification"
	// VertexNominationTask represents a vertex nomination task on graph data
	VertexNominationTask = "vertexNomination"
	// CommunityDetectionTask represents an unsupervised community detectiontask on on graph data
	CommunityDetectionTask = "communityDetection"
	// GraphMatchingTask represents an unsupervised matching task on graph data
	GraphMatchingTask = "graphMatching"
	// CollaborativeFilteringTask represents a collaborative filtering recommendation task on basic tabular data
	CollaborativeFilteringTask = "collaborativeFiltering"
	// ObjectDetectionTask represents an object detection task on image data
	ObjectDetectionTask = "objectDetection"
	// SemiSupervisedTask represents a semi-supervised classification task on tabular data
	SemiSupervisedTask = "semiSupervised"
	// BinaryTask represents task involving a single binary value for each prediction
	BinaryTask = "binary"
	// MultiClassTask represents a task involving a multi class value for each prediction
	MultiClassTask = "multiClass"
	// MultiLabelTask represents a task involving multiple lables for each each prediction
	MultiLabelTask = "multiLabel"
	// UnivariateTask represents a task involving predictions on a single variable
	UnivariateTask = "univariate"
	// MultivariateTask represents a task involving predictions on multiple variables
	MultivariateTask = "multivariate"
	// OverlappingTask represents a task involving overlapping predictions
	OverlappingTask = "overlapping"
	// NonOverlappingTask represents a task involving non-overlapping predictions
	NonOverlappingTask = "nonOverlapping"
	// TabularTask represents a task involving tabular data
	TabularTask = "tabular"
	// RelationalTask represents a task involving relational data
	RelationalTask = "relational"
	// ImageTask represents a task involving image data
	ImageTask = "image"
	// AudioTask represents a task involving audio data
	AudioTask = "audio"
	// VideoTask represents a task involving video data
	VideoTask = "video"
	// SpeechTask represents a task involving speech data
	SpeechTask = "speech"
	// TextTask represents a task involving text data
	TextTask = "text"
	// GraphTask represents a task involving graph data
	GraphTask = "graph"
	// MultiGraphTask represents a task involving multiple graph data
	MultiGraphTask = "multigraph"
	// TimeSeriesTask represents a task involving timeseries data
	TimeSeriesTask = "timeseries"
	// GroupedTask represents a task involving grouped data
	GroupedTask = "grouped"
	// GeospatialTask represents a task involving geospatial data
	GeospatialTask = "geospatial"
	// RemoteSensingTask represents a task involving remote sensing data
	RemoteSensingTask = "remoteSensing"
	// LupiTask represents a task involving LUPI (Learning Using Priveleged Information) data
	LupiTask = "lupi"
)

var (
	// cached ta3ta2 API version
	apiVersion       string
	problemMetricMap = map[string]string{
		"accuracy":                    "ACCURACY",
		"precision":                   "PRECISION",
		"recall":                      "RECALL",
		"f1":                          "F1",
		"f1Micro":                     "F1_MICRO",
		"f1Macro":                     "F1_MACRO",
		"rocAuc":                      "ROC_AUC",
		"rocAucMicro":                 "ROC_AUC_MICRO",
		"rocAucMacro":                 "ROC_AUC_MACRO",
		"meanSquaredError":            "MEAN_SQUARED_ERROR",
		"rootMeanSquaredError":        "ROOT_MEAN_SQUARED_ERROR",
		"rootMeanSquaredErrorAvg":     "ROOT_MEAN_SQUARED_ERROR_AVG",
		"meanAbsoluteError":           "MEAN_ABSOLUTE_ERROR",
		"rSquared":                    "R_SQUARED",
		"normalizedMutualInformation": "NORMALIZED_MUTUAL_INFORMATION",
		"jaccardSimilarityScore":      "JACCARD_SIMILARITY_SCORE",
		"precisionAtTopK":             "PRECISION_AT_TOP_K",
		"objectDetectionAP":           "OBJECT_DETECTION_AVERAGE_PRECISION",
	}
	problemTaskMap = map[string]string{
		ClassificationTask:         "CLASSIFICATION",
		RegressionTask:             "REGRESSION",
		ClusteringTask:             "CLUSTERING",
		LinkPredictionTask:         "LINK_PREDICTION",
		VertexNominationTask:       "VERTEX_NOMINATION",
		VertexClassificationTask:   "VERTEX_CLASSIFICATION",
		CommunityDetectionTask:     "COMMUNITY_DETECTION",
		GraphMatchingTask:          "GRAPH_MATCHING",
		ForecastingTask:            "FORECASTING",
		CollaborativeFilteringTask: "COLLABORATIVE_FILTERING",
		ObjectDetectionTask:        "OBJECT_DETECTION",
		SemiSupervisedTask:         "SEMISUPERVISED",
		BinaryTask:                 "BINARY",
		MultiClassTask:             "MULTICLASS",
		MultiLabelTask:             "MULTILABEL",
		UnivariateTask:             "UNIVARIATE",
		MultivariateTask:           "MULTIVARIATE",
		OverlappingTask:            "OVERLAPPING",
		NonOverlappingTask:         "NONOVERLAPPING",
		TabularTask:                "TABULAR",
		RelationalTask:             "RELATIONAL",
		ImageTask:                  "IMAGE",
		AudioTask:                  "AUDIO",
		VideoTask:                  "VIDEO",
		SpeechTask:                 "SPEECH",
		TextTask:                   "TEXT",
		GraphTask:                  "GRAPH",
		MultiGraphTask:             "MULTIGRAPH",
		TimeSeriesTask:             "TIME_SERIES",
		GroupedTask:                "GROUPED",
		GeospatialTask:             "GEOSPATIAL",
		RemoteSensingTask:          "REMOTE_SENSING",
		LupiTask:                   "LUPI",
	}
	defaultTaskMetricMap = map[string]string{
		ClassificationTask:         "f1Macro",
		RegressionTask:             "meanAbsoluteError",
		ClusteringTask:             "normalizedMutualInformation",
		LinkPredictionTask:         "accuracy",
		VertexNominationTask:       "accuracy",
		CommunityDetectionTask:     "accuracy",
		GraphMatchingTask:          "accuracy",
		ForecastingTask:            "rSquared",
		CollaborativeFilteringTask: "rSquared",
		ObjectDetectionTask:        "objectDetectionAP",
	}
	metricScoreMultiplier = map[string]float64{
		"ACCURACY":                           1,
		"PRECISION":                          1,
		"RECALL":                             1,
		"F1":                                 1,
		"F1_MICRO":                           1,
		"F1_MACRO":                           1,
		"ROC_AUC":                            1,
		"ROC_AUC_MICRO":                      1,
		"ROC_AUC_MACRO":                      1,
		"MEAN_SQUARED_ERROR":                 -1,
		"ROOT_MEAN_SQUARED_ERROR":            -1,
		"ROOT_MEAN_SQUARED_ERROR_AVG":        -1,
		"MEAN_ABSOLUTE_ERROR":                -1,
		"R_SQUARED":                          1,
		"NORMALIZED_MUTUAL_INFORMATION":      1,
		"JACCARD_SIMILARITY_SCORE":           1,
		"PRECISION_AT_TOP_K":                 1,
		"OBJECT_DETECTION_AVERAGE_PRECISION": 1,
	}
	metricLabel = map[string]string{
		"ACCURACY":                           "Accuracy",
		"PRECISION":                          "Precision",
		"RECALL":                             "Recall",
		"F1":                                 "F1",
		"F1_MICRO":                           "F1 Micro",
		"F1_MACRO":                           "F1 Macro",
		"ROC_AUC":                            "ROC AUC",
		"ROC_AUC_MICRO":                      "ROC AUC Micro",
		"ROC_AUC_MACRO":                      "ROC AUC Macro",
		"MEAN_SQUARED_ERROR":                 "MSE",
		"ROOT_MEAN_SQUARED_ERROR":            "RMSE",
		"ROOT_MEAN_SQUARED_ERROR_AVG":        "RMSE Avg",
		"MEAN_ABSOLUTE_ERROR":                "MAE",
		"R_SQUARED":                          "R Squared",
		"NORMALIZED_MUTUAL_INFORMATION":      "Normalized MI",
		"JACCARD_SIMILARITY_SCORE":           "Jaccard Similarity",
		"PRECISION_AT_TOP_K":                 "Precision Top K",
		"OBJECT_DETECTION_AVERAGE_PRECISION": "Avg Precision",
	}
)

// ConvertProblemMetricToTA2 converts a problem schema metric to a TA2 metric.
func ConvertProblemMetricToTA2(metric string) string {
	return problemMetricMap[metric]
}

// ConvertProblemTaskToTA2 converts a problem schema metric to a TA2 task.
func ConvertProblemTaskToTA2(metric string) string {
	return problemTaskMap[metric]
}

// GetMetricScoreMultiplier returns a weight to determine whether a higher or
// lower score is `better`.
func GetMetricScoreMultiplier(metric string) float64 {
	return metricScoreMultiplier[metric]
}

// GetMetricLabel returns a label string for a metric.
func GetMetricLabel(metric string) string {
	return metricLabel[metric]
}

// GetDefaultTaskMetricTA3 returns the default TA3 metric type for a supplied
// TA3 task type.
func GetDefaultTaskMetricTA3(task string) string {
	return defaultTaskMetricMap[task]
}

// ConvertMetricsFromTA3ToTA2 converts metrics from TA3 to TA2 values.
func ConvertMetricsFromTA3ToTA2(metrics []string) []*pipeline.ProblemPerformanceMetric {
	var res []*pipeline.ProblemPerformanceMetric
	for _, metric := range metrics {
		ta2Metric := ConvertProblemMetricToTA2(metric)
		var metricSet pipeline.PerformanceMetric
		if ta2Metric == "" {
			log.Warnf("unrecognized metric ('%s'), defaulting to undefined", metric)
			metricSet = pipeline.PerformanceMetric_METRIC_UNDEFINED
		} else {
			metricAdjusted, ok := pipeline.PerformanceMetric_value[ta2Metric]
			if !ok {
				log.Warnf("undefined metric found ('%s'), defaulting to undefined", ta2Metric)
				metricSet = pipeline.PerformanceMetric_METRIC_UNDEFINED
			} else {
				metricSet = pipeline.PerformanceMetric(metricAdjusted)
			}
		}
		res = append(res, &pipeline.ProblemPerformanceMetric{
			Metric: metricSet,
		})
	}
	return res
}

// ConvertTaskKeywordFromTA3ToTA2 converts a task from TA3 to TA2.
func ConvertTaskKeywordFromTA3ToTA2(taskKeyword string) pipeline.TaskKeyword {
	ta2Task := ConvertProblemTaskToTA2(taskKeyword)
	if ta2Task == "" {
		log.Warnf("unrecognized task type ('%s'), defaulting to undefined", taskKeyword)
		return pipeline.TaskKeyword_TASK_KEYWORD_UNDEFINED
	}
	task, ok := pipeline.TaskKeyword_value[ta2Task]
	if !ok {
		log.Warnf("undefined task type found ('%s'), defaulting to undefined", ta2Task)
		return pipeline.TaskKeyword_TASK_KEYWORD_UNDEFINED
	}
	return pipeline.TaskKeyword(task)
}

// ConvertTargetFeaturesTA3ToTA2 creates a problem target from a target name.
func ConvertTargetFeaturesTA3ToTA2(target string, columnIndex int) []*pipeline.ProblemTarget {
	return []*pipeline.ProblemTarget{
		{
			ColumnName:  target,
			ResourceId:  defaultResourceID,
			TargetIndex: 0,
			ColumnIndex: int32(columnIndex),
		},
	}
}

// ConvertDatasetTA3ToTA2 converts a dataset name from TA3 to TA2.
func ConvertDatasetTA3ToTA2(dataset string) string {
	return dataset
}

// GetAPIVersion retrieves the ta3-ta2 API version embedded in the pipeline_core.proto file.  This is
// a non-trivial operation, so the value is cached for quick access.
func GetAPIVersion() string {
	if apiVersion != "" {
		return apiVersion
	}

	// Get the raw file descriptor bytes
	fileDesc := proto.FileDescriptor(pipeline.E_ProtocolVersion.Filename)
	if fileDesc == nil {
		log.Errorf("failed to find file descriptor for %v", pipeline.E_ProtocolVersion.Filename)
		return unknownAPIVersion
	}

	// Open a gzip reader and decompress
	r, err := gzip.NewReader(bytes.NewReader(fileDesc))
	if err != nil {
		log.Errorf("failed to open gzip reader: %v", err)
		return unknownAPIVersion
	}
	defer r.Close()

	b, err := ioutil.ReadAll(r)
	if err != nil {
		log.Errorf("failed to decompress descriptor: %v", err)
		return unknownAPIVersion
	}

	// Unmarshall the bytes from the proto format
	fd := &protobuf.FileDescriptorProto{}
	if err := proto.Unmarshal(b, fd); err != nil {
		log.Errorf("malformed FileDescriptorProto: %v", err)
		return unknownAPIVersion
	}

	// Fetch the extension from the FileDescriptorOptions message
	ex, err := proto.GetExtension(fd.GetOptions(), pipeline.E_ProtocolVersion)
	if err != nil {
		log.Errorf("failed to fetch extension: %v", err)
		return unknownAPIVersion
	}

	apiVersion = *ex.(*string)

	return apiVersion
}
