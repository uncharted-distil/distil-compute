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
		"classification":         "CLASSIFICATION",
		"regression":             "REGRESSION",
		"clustering":             "CLUSTERING",
		"linkPrediction":         "LINK_PREDICTION",
		"vertexNomination":       "VERTEX_NOMINATION",
		"communityDetection":     "COMMUNITY_DETECTION",
		"graphClustering":        "GRAPH_CLUSTERING",
		"graphMatching":          "GRAPH_MATCHING",
		"timeSeriesForecasting":  "TIME_SERIES_FORECASTING",
		"collaborativeFiltering": "COLLABORATIVE_FILTERING",
		"objectDetection":        "OBJECT_DETECTION",
	}
	problemTaskSubMap = map[string]string{
		"none":           "NONE",
		"binary":         "BINARY",
		"multiClass":     "MULTICLASS",
		"multiLabel":     "MULTILABEL",
		"univariate":     "UNIVARIATE",
		"multivariate":   "MULTIVARIATE",
		"overlapping":    "OVERLAPPING",
		"nonOverlapping": "NONOVERLAPPING",
	}
	defaultProblemSubTaskMap = map[string]string{
		"classification":         "multiClass",
		"regression":             "univariate",
		"clustering":             "none",
		"vertexNomination":       "none",
		"communityDetection":     "none",
		"graphClustering":        "none",
		"graphMatching":          "none",
		"timeSeriesForecasting":  "none",
		"collaborativeFiltering": "none",
		"objectDetection":        "overlapping",
	}
	defaultTaskMetricMap = map[string]string{
		"classification":         "f1Macro",
		"regression":             "rSquared",
		"clustering":             "normalizedMutualInformation",
		"linkPrediction":         "accuracy",
		"vertexNomination":       "accuracy",
		"communityDetection":     "accuracy",
		"graphClustering":        "normalizedMutualInformation",
		"graphMatching":          "accuracy",
		"timeSeriesForecasting":  "rSquared",
		"collaborativeFiltering": "rSquared",
		"objectDetection":        "objectDetectionAP",
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

// ConvertProblemTaskSubToTA2 converts a problem schema metric to a TA2 task sub.
func ConvertProblemTaskSubToTA2(metric string) string {
	return problemTaskSubMap[metric]
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

// GetDefaultTaskSubTypeTA3 returns the default TA3 subtask type for a supplied
// TA3 task type.
func GetDefaultTaskSubTypeTA3(task string) string {
	return defaultProblemSubTaskMap[task]
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

// ConvertTaskTypeFromTA3ToTA2 converts a task from TA3 to TA2.
func ConvertTaskTypeFromTA3ToTA2(taskType string) pipeline.TaskType {
	ta2Task := ConvertProblemTaskToTA2(taskType)
	if ta2Task == "" {
		log.Warnf("unrecognized task type ('%s'), defaulting to undefined", taskType)
		return pipeline.TaskType_TASK_TYPE_UNDEFINED
	}
	task, ok := pipeline.TaskType_value[ta2Task]
	if !ok {
		log.Warnf("undefined task type found ('%s'), defaulting to undefined", ta2Task)
		return pipeline.TaskType_TASK_TYPE_UNDEFINED
	}
	return pipeline.TaskType(task)
}

// ConvertTaskSubTypeFromTA3ToTA2 converts a task sub type from TA3 to TA2.
func ConvertTaskSubTypeFromTA3ToTA2(taskSubType string) pipeline.TaskSubtype {
	ta2TaskSub := ConvertProblemTaskSubToTA2(taskSubType)
	if ta2TaskSub == "" {
		log.Warnf("unrecognized task sub type ('%s'), defaulting to undefined", taskSubType)
		return pipeline.TaskSubtype_TASK_SUBTYPE_UNDEFINED
	}
	task, ok := pipeline.TaskSubtype_value[ta2TaskSub]
	if !ok {
		log.Warnf("undefined task sub type found ('%s'), defaulting to undefined", ta2TaskSub)
		return pipeline.TaskSubtype_TASK_SUBTYPE_UNDEFINED
	}
	return pipeline.TaskSubtype(task)
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
