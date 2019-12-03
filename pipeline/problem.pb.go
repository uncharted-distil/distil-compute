// Code generated by protoc-gen-go. DO NOT EDIT.
// source: problem.proto

package pipeline

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

// Task keyword of the problem.
type TaskKeyword int32

const (
	// Default value. Not to be used.
	TaskKeyword_TASK_KEYWORD_UNDEFINED  TaskKeyword = 0
	TaskKeyword_CLASSIFICATION          TaskKeyword = 1
	TaskKeyword_REGRESSION              TaskKeyword = 2
	TaskKeyword_CLUSTERING              TaskKeyword = 3
	TaskKeyword_LINK_PREDICTION         TaskKeyword = 4
	TaskKeyword_VERTEX_NOMINATION       TaskKeyword = 5
	TaskKeyword_VERTEX_CLASSIFICATION   TaskKeyword = 6
	TaskKeyword_COMMUNITY_DETECTION     TaskKeyword = 7
	TaskKeyword_GRAPH_MATCHING          TaskKeyword = 8
	TaskKeyword_FORECASTING             TaskKeyword = 9
	TaskKeyword_COLLABORATIVE_FILTERING TaskKeyword = 10
	TaskKeyword_OBJECT_DETECTION        TaskKeyword = 11
	TaskKeyword_SEMISUPERVISED          TaskKeyword = 12
	TaskKeyword_BINARY                  TaskKeyword = 13
	TaskKeyword_MULTICLASS              TaskKeyword = 14
	TaskKeyword_MULTILABEL              TaskKeyword = 15
	TaskKeyword_UNIVARIATE              TaskKeyword = 16
	TaskKeyword_MULTIVARIATE            TaskKeyword = 17
	TaskKeyword_OVERLAPPING             TaskKeyword = 18
	TaskKeyword_NONOVERLAPPING          TaskKeyword = 19
	TaskKeyword_TABULAR                 TaskKeyword = 20
	TaskKeyword_RELATIONAL              TaskKeyword = 21
	TaskKeyword_IMAGE                   TaskKeyword = 22
	TaskKeyword_AUDIO                   TaskKeyword = 23
	TaskKeyword_VIDEO                   TaskKeyword = 24
	TaskKeyword_SPEECH                  TaskKeyword = 25
	TaskKeyword_TEXT                    TaskKeyword = 26
	TaskKeyword_GRAPH                   TaskKeyword = 27
	TaskKeyword_MULTIGRAPH              TaskKeyword = 28
	TaskKeyword_TIME_SERIES             TaskKeyword = 29
	TaskKeyword_GROUPED                 TaskKeyword = 30
	TaskKeyword_GEOSPATIAL              TaskKeyword = 31
	TaskKeyword_REMOTE_SENSING          TaskKeyword = 32
	TaskKeyword_LUPI                    TaskKeyword = 33
	TaskKeyword_MISSING_METADATA        TaskKeyword = 34
)

var TaskKeyword_name = map[int32]string{
	0:  "TASK_KEYWORD_UNDEFINED",
	1:  "CLASSIFICATION",
	2:  "REGRESSION",
	3:  "CLUSTERING",
	4:  "LINK_PREDICTION",
	5:  "VERTEX_NOMINATION",
	6:  "VERTEX_CLASSIFICATION",
	7:  "COMMUNITY_DETECTION",
	8:  "GRAPH_MATCHING",
	9:  "FORECASTING",
	10: "COLLABORATIVE_FILTERING",
	11: "OBJECT_DETECTION",
	12: "SEMISUPERVISED",
	13: "BINARY",
	14: "MULTICLASS",
	15: "MULTILABEL",
	16: "UNIVARIATE",
	17: "MULTIVARIATE",
	18: "OVERLAPPING",
	19: "NONOVERLAPPING",
	20: "TABULAR",
	21: "RELATIONAL",
	22: "IMAGE",
	23: "AUDIO",
	24: "VIDEO",
	25: "SPEECH",
	26: "TEXT",
	27: "GRAPH",
	28: "MULTIGRAPH",
	29: "TIME_SERIES",
	30: "GROUPED",
	31: "GEOSPATIAL",
	32: "REMOTE_SENSING",
	33: "LUPI",
	34: "MISSING_METADATA",
}

var TaskKeyword_value = map[string]int32{
	"TASK_KEYWORD_UNDEFINED":  0,
	"CLASSIFICATION":          1,
	"REGRESSION":              2,
	"CLUSTERING":              3,
	"LINK_PREDICTION":         4,
	"VERTEX_NOMINATION":       5,
	"VERTEX_CLASSIFICATION":   6,
	"COMMUNITY_DETECTION":     7,
	"GRAPH_MATCHING":          8,
	"FORECASTING":             9,
	"COLLABORATIVE_FILTERING": 10,
	"OBJECT_DETECTION":        11,
	"SEMISUPERVISED":          12,
	"BINARY":                  13,
	"MULTICLASS":              14,
	"MULTILABEL":              15,
	"UNIVARIATE":              16,
	"MULTIVARIATE":            17,
	"OVERLAPPING":             18,
	"NONOVERLAPPING":          19,
	"TABULAR":                 20,
	"RELATIONAL":              21,
	"IMAGE":                   22,
	"AUDIO":                   23,
	"VIDEO":                   24,
	"SPEECH":                  25,
	"TEXT":                    26,
	"GRAPH":                   27,
	"MULTIGRAPH":              28,
	"TIME_SERIES":             29,
	"GROUPED":                 30,
	"GEOSPATIAL":              31,
	"REMOTE_SENSING":          32,
	"LUPI":                    33,
	"MISSING_METADATA":        34,
}

func (x TaskKeyword) String() string {
	return proto.EnumName(TaskKeyword_name, int32(x))
}

func (TaskKeyword) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_b319862c9661813c, []int{0}
}

// The evaluation metric for any potential solution.
type PerformanceMetric int32

const (
	// Default value. Not to be used.
	PerformanceMetric_METRIC_UNDEFINED PerformanceMetric = 0
	// The following are the only evaluation methods required
	// to be supported for the ScoreSolution call.
	PerformanceMetric_ACCURACY                           PerformanceMetric = 1
	PerformanceMetric_PRECISION                          PerformanceMetric = 2
	PerformanceMetric_RECALL                             PerformanceMetric = 3
	PerformanceMetric_F1                                 PerformanceMetric = 4
	PerformanceMetric_F1_MICRO                           PerformanceMetric = 5
	PerformanceMetric_F1_MACRO                           PerformanceMetric = 6
	PerformanceMetric_ROC_AUC                            PerformanceMetric = 7
	PerformanceMetric_ROC_AUC_MICRO                      PerformanceMetric = 8
	PerformanceMetric_ROC_AUC_MACRO                      PerformanceMetric = 9
	PerformanceMetric_MEAN_SQUARED_ERROR                 PerformanceMetric = 10
	PerformanceMetric_ROOT_MEAN_SQUARED_ERROR            PerformanceMetric = 11
	PerformanceMetric_MEAN_ABSOLUTE_ERROR                PerformanceMetric = 12
	PerformanceMetric_R_SQUARED                          PerformanceMetric = 13
	PerformanceMetric_NORMALIZED_MUTUAL_INFORMATION      PerformanceMetric = 14
	PerformanceMetric_JACCARD_SIMILARITY_SCORE           PerformanceMetric = 15
	PerformanceMetric_PRECISION_AT_TOP_K                 PerformanceMetric = 17
	PerformanceMetric_OBJECT_DETECTION_AVERAGE_PRECISION PerformanceMetric = 18
	PerformanceMetric_HAMMING_LOSS                       PerformanceMetric = 19
	// This metric can be used to ask TA2 to rank a solution as part of
	// all found solutions of a given "SearchSolutions" call. Rank is a
	// floating-point number. Lower numbers represent better solutions.
	// Presently evaluation requirements are that ranks should be non-negative
	// and that each ranked pipeline have a different rank (for all
	// solutions of a given SearchSolutions call). Only possible with
	// "RANKING" evaluation method.
	PerformanceMetric_RANK PerformanceMetric = 99
	// The rest are defined to allow expressing internal evaluation
	// scores used by TA2 during pipeline search. If any you are using
	// is missing, feel free to request it to be added.
	// Average loss of an unspecified loss function.
	PerformanceMetric_LOSS PerformanceMetric = 100
)

var PerformanceMetric_name = map[int32]string{
	0:   "METRIC_UNDEFINED",
	1:   "ACCURACY",
	2:   "PRECISION",
	3:   "RECALL",
	4:   "F1",
	5:   "F1_MICRO",
	6:   "F1_MACRO",
	7:   "ROC_AUC",
	8:   "ROC_AUC_MICRO",
	9:   "ROC_AUC_MACRO",
	10:  "MEAN_SQUARED_ERROR",
	11:  "ROOT_MEAN_SQUARED_ERROR",
	12:  "MEAN_ABSOLUTE_ERROR",
	13:  "R_SQUARED",
	14:  "NORMALIZED_MUTUAL_INFORMATION",
	15:  "JACCARD_SIMILARITY_SCORE",
	17:  "PRECISION_AT_TOP_K",
	18:  "OBJECT_DETECTION_AVERAGE_PRECISION",
	19:  "HAMMING_LOSS",
	99:  "RANK",
	100: "LOSS",
}

var PerformanceMetric_value = map[string]int32{
	"METRIC_UNDEFINED":                   0,
	"ACCURACY":                           1,
	"PRECISION":                          2,
	"RECALL":                             3,
	"F1":                                 4,
	"F1_MICRO":                           5,
	"F1_MACRO":                           6,
	"ROC_AUC":                            7,
	"ROC_AUC_MICRO":                      8,
	"ROC_AUC_MACRO":                      9,
	"MEAN_SQUARED_ERROR":                 10,
	"ROOT_MEAN_SQUARED_ERROR":            11,
	"MEAN_ABSOLUTE_ERROR":                12,
	"R_SQUARED":                          13,
	"NORMALIZED_MUTUAL_INFORMATION":      14,
	"JACCARD_SIMILARITY_SCORE":           15,
	"PRECISION_AT_TOP_K":                 17,
	"OBJECT_DETECTION_AVERAGE_PRECISION": 18,
	"HAMMING_LOSS":                       19,
	"RANK":                               99,
	"LOSS":                               100,
}

func (x PerformanceMetric) String() string {
	return proto.EnumName(PerformanceMetric_name, int32(x))
}

func (PerformanceMetric) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_b319862c9661813c, []int{1}
}

type ProblemPerformanceMetric struct {
	Metric PerformanceMetric `protobuf:"varint,1,opt,name=metric,proto3,enum=PerformanceMetric" json:"metric,omitempty"`
	// Additional params used by some metrics.
	K                    int32    `protobuf:"varint,2,opt,name=k,proto3" json:"k,omitempty"`
	PosLabel             string   `protobuf:"bytes,3,opt,name=pos_label,json=posLabel,proto3" json:"pos_label,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ProblemPerformanceMetric) Reset()         { *m = ProblemPerformanceMetric{} }
func (m *ProblemPerformanceMetric) String() string { return proto.CompactTextString(m) }
func (*ProblemPerformanceMetric) ProtoMessage()    {}
func (*ProblemPerformanceMetric) Descriptor() ([]byte, []int) {
	return fileDescriptor_b319862c9661813c, []int{0}
}

func (m *ProblemPerformanceMetric) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ProblemPerformanceMetric.Unmarshal(m, b)
}
func (m *ProblemPerformanceMetric) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ProblemPerformanceMetric.Marshal(b, m, deterministic)
}
func (m *ProblemPerformanceMetric) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ProblemPerformanceMetric.Merge(m, src)
}
func (m *ProblemPerformanceMetric) XXX_Size() int {
	return xxx_messageInfo_ProblemPerformanceMetric.Size(m)
}
func (m *ProblemPerformanceMetric) XXX_DiscardUnknown() {
	xxx_messageInfo_ProblemPerformanceMetric.DiscardUnknown(m)
}

var xxx_messageInfo_ProblemPerformanceMetric proto.InternalMessageInfo

func (m *ProblemPerformanceMetric) GetMetric() PerformanceMetric {
	if m != nil {
		return m.Metric
	}
	return PerformanceMetric_METRIC_UNDEFINED
}

func (m *ProblemPerformanceMetric) GetK() int32 {
	if m != nil {
		return m.K
	}
	return 0
}

func (m *ProblemPerformanceMetric) GetPosLabel() string {
	if m != nil {
		return m.PosLabel
	}
	return ""
}

type Problem struct {
	TaskKeywords         []TaskKeyword               `protobuf:"varint,8,rep,packed,name=task_keywords,json=taskKeywords,proto3,enum=TaskKeyword" json:"task_keywords,omitempty"`
	PerformanceMetrics   []*ProblemPerformanceMetric `protobuf:"bytes,7,rep,name=performance_metrics,json=performanceMetrics,proto3" json:"performance_metrics,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                    `json:"-"`
	XXX_unrecognized     []byte                      `json:"-"`
	XXX_sizecache        int32                       `json:"-"`
}

func (m *Problem) Reset()         { *m = Problem{} }
func (m *Problem) String() string { return proto.CompactTextString(m) }
func (*Problem) ProtoMessage()    {}
func (*Problem) Descriptor() ([]byte, []int) {
	return fileDescriptor_b319862c9661813c, []int{1}
}

func (m *Problem) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Problem.Unmarshal(m, b)
}
func (m *Problem) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Problem.Marshal(b, m, deterministic)
}
func (m *Problem) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Problem.Merge(m, src)
}
func (m *Problem) XXX_Size() int {
	return xxx_messageInfo_Problem.Size(m)
}
func (m *Problem) XXX_DiscardUnknown() {
	xxx_messageInfo_Problem.DiscardUnknown(m)
}

var xxx_messageInfo_Problem proto.InternalMessageInfo

func (m *Problem) GetTaskKeywords() []TaskKeyword {
	if m != nil {
		return m.TaskKeywords
	}
	return nil
}

func (m *Problem) GetPerformanceMetrics() []*ProblemPerformanceMetric {
	if m != nil {
		return m.PerformanceMetrics
	}
	return nil
}

type ProblemTarget struct {
	TargetIndex          int32    `protobuf:"varint,1,opt,name=target_index,json=targetIndex,proto3" json:"target_index,omitempty"`
	ResourceId           string   `protobuf:"bytes,2,opt,name=resource_id,json=resourceId,proto3" json:"resource_id,omitempty"`
	ColumnIndex          int32    `protobuf:"varint,3,opt,name=column_index,json=columnIndex,proto3" json:"column_index,omitempty"`
	ColumnName           string   `protobuf:"bytes,4,opt,name=column_name,json=columnName,proto3" json:"column_name,omitempty"`
	ClustersNumber       int32    `protobuf:"varint,5,opt,name=clusters_number,json=clustersNumber,proto3" json:"clusters_number,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ProblemTarget) Reset()         { *m = ProblemTarget{} }
func (m *ProblemTarget) String() string { return proto.CompactTextString(m) }
func (*ProblemTarget) ProtoMessage()    {}
func (*ProblemTarget) Descriptor() ([]byte, []int) {
	return fileDescriptor_b319862c9661813c, []int{2}
}

func (m *ProblemTarget) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ProblemTarget.Unmarshal(m, b)
}
func (m *ProblemTarget) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ProblemTarget.Marshal(b, m, deterministic)
}
func (m *ProblemTarget) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ProblemTarget.Merge(m, src)
}
func (m *ProblemTarget) XXX_Size() int {
	return xxx_messageInfo_ProblemTarget.Size(m)
}
func (m *ProblemTarget) XXX_DiscardUnknown() {
	xxx_messageInfo_ProblemTarget.DiscardUnknown(m)
}

var xxx_messageInfo_ProblemTarget proto.InternalMessageInfo

func (m *ProblemTarget) GetTargetIndex() int32 {
	if m != nil {
		return m.TargetIndex
	}
	return 0
}

func (m *ProblemTarget) GetResourceId() string {
	if m != nil {
		return m.ResourceId
	}
	return ""
}

func (m *ProblemTarget) GetColumnIndex() int32 {
	if m != nil {
		return m.ColumnIndex
	}
	return 0
}

func (m *ProblemTarget) GetColumnName() string {
	if m != nil {
		return m.ColumnName
	}
	return ""
}

func (m *ProblemTarget) GetClustersNumber() int32 {
	if m != nil {
		return m.ClustersNumber
	}
	return 0
}

type ProblemPrivilegedData struct {
	PrivilegedDataIndex  int32    `protobuf:"varint,1,opt,name=privileged_data_index,json=privilegedDataIndex,proto3" json:"privileged_data_index,omitempty"`
	ResourceId           string   `protobuf:"bytes,2,opt,name=resource_id,json=resourceId,proto3" json:"resource_id,omitempty"`
	ColumnIndex          int32    `protobuf:"varint,3,opt,name=column_index,json=columnIndex,proto3" json:"column_index,omitempty"`
	ColumnName           string   `protobuf:"bytes,4,opt,name=column_name,json=columnName,proto3" json:"column_name,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ProblemPrivilegedData) Reset()         { *m = ProblemPrivilegedData{} }
func (m *ProblemPrivilegedData) String() string { return proto.CompactTextString(m) }
func (*ProblemPrivilegedData) ProtoMessage()    {}
func (*ProblemPrivilegedData) Descriptor() ([]byte, []int) {
	return fileDescriptor_b319862c9661813c, []int{3}
}

func (m *ProblemPrivilegedData) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ProblemPrivilegedData.Unmarshal(m, b)
}
func (m *ProblemPrivilegedData) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ProblemPrivilegedData.Marshal(b, m, deterministic)
}
func (m *ProblemPrivilegedData) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ProblemPrivilegedData.Merge(m, src)
}
func (m *ProblemPrivilegedData) XXX_Size() int {
	return xxx_messageInfo_ProblemPrivilegedData.Size(m)
}
func (m *ProblemPrivilegedData) XXX_DiscardUnknown() {
	xxx_messageInfo_ProblemPrivilegedData.DiscardUnknown(m)
}

var xxx_messageInfo_ProblemPrivilegedData proto.InternalMessageInfo

func (m *ProblemPrivilegedData) GetPrivilegedDataIndex() int32 {
	if m != nil {
		return m.PrivilegedDataIndex
	}
	return 0
}

func (m *ProblemPrivilegedData) GetResourceId() string {
	if m != nil {
		return m.ResourceId
	}
	return ""
}

func (m *ProblemPrivilegedData) GetColumnIndex() int32 {
	if m != nil {
		return m.ColumnIndex
	}
	return 0
}

func (m *ProblemPrivilegedData) GetColumnName() string {
	if m != nil {
		return m.ColumnName
	}
	return ""
}

type ForecastingHorizon struct {
	ResourceId           string   `protobuf:"bytes,1,opt,name=resource_id,json=resourceId,proto3" json:"resource_id,omitempty"`
	ColumnIndex          int32    `protobuf:"varint,2,opt,name=column_index,json=columnIndex,proto3" json:"column_index,omitempty"`
	ColumnName           string   `protobuf:"bytes,3,opt,name=column_name,json=columnName,proto3" json:"column_name,omitempty"`
	HorizonValue         float64  `protobuf:"fixed64,4,opt,name=horizon_value,json=horizonValue,proto3" json:"horizon_value,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ForecastingHorizon) Reset()         { *m = ForecastingHorizon{} }
func (m *ForecastingHorizon) String() string { return proto.CompactTextString(m) }
func (*ForecastingHorizon) ProtoMessage()    {}
func (*ForecastingHorizon) Descriptor() ([]byte, []int) {
	return fileDescriptor_b319862c9661813c, []int{4}
}

func (m *ForecastingHorizon) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ForecastingHorizon.Unmarshal(m, b)
}
func (m *ForecastingHorizon) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ForecastingHorizon.Marshal(b, m, deterministic)
}
func (m *ForecastingHorizon) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ForecastingHorizon.Merge(m, src)
}
func (m *ForecastingHorizon) XXX_Size() int {
	return xxx_messageInfo_ForecastingHorizon.Size(m)
}
func (m *ForecastingHorizon) XXX_DiscardUnknown() {
	xxx_messageInfo_ForecastingHorizon.DiscardUnknown(m)
}

var xxx_messageInfo_ForecastingHorizon proto.InternalMessageInfo

func (m *ForecastingHorizon) GetResourceId() string {
	if m != nil {
		return m.ResourceId
	}
	return ""
}

func (m *ForecastingHorizon) GetColumnIndex() int32 {
	if m != nil {
		return m.ColumnIndex
	}
	return 0
}

func (m *ForecastingHorizon) GetColumnName() string {
	if m != nil {
		return m.ColumnName
	}
	return ""
}

func (m *ForecastingHorizon) GetHorizonValue() float64 {
	if m != nil {
		return m.HorizonValue
	}
	return 0
}

type ProblemInput struct {
	// Should match one of input datasets given to the pipeline search.
	// Every "Dataset" object has an "id" associated with it and is available
	// in its metadata. That ID is then used here to reference those inputs.
	DatasetId string `protobuf:"bytes,1,opt,name=dataset_id,json=datasetId,proto3" json:"dataset_id,omitempty"`
	// Targets should resolve to columns in a given dataset.
	Targets []*ProblemTarget `protobuf:"bytes,2,rep,name=targets,proto3" json:"targets,omitempty"`
	// A list of privileged data columns related to unavailable attributes during testing.
	// These columns do not have data available in the test split of a dataset.
	PrivilegedData       []*ProblemPrivilegedData `protobuf:"bytes,3,rep,name=privileged_data,json=privilegedData,proto3" json:"privileged_data,omitempty"`
	ForecastingHorizon   *ForecastingHorizon      `protobuf:"bytes,4,opt,name=forecasting_horizon,json=forecastingHorizon,proto3" json:"forecasting_horizon,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                 `json:"-"`
	XXX_unrecognized     []byte                   `json:"-"`
	XXX_sizecache        int32                    `json:"-"`
}

func (m *ProblemInput) Reset()         { *m = ProblemInput{} }
func (m *ProblemInput) String() string { return proto.CompactTextString(m) }
func (*ProblemInput) ProtoMessage()    {}
func (*ProblemInput) Descriptor() ([]byte, []int) {
	return fileDescriptor_b319862c9661813c, []int{5}
}

func (m *ProblemInput) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ProblemInput.Unmarshal(m, b)
}
func (m *ProblemInput) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ProblemInput.Marshal(b, m, deterministic)
}
func (m *ProblemInput) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ProblemInput.Merge(m, src)
}
func (m *ProblemInput) XXX_Size() int {
	return xxx_messageInfo_ProblemInput.Size(m)
}
func (m *ProblemInput) XXX_DiscardUnknown() {
	xxx_messageInfo_ProblemInput.DiscardUnknown(m)
}

var xxx_messageInfo_ProblemInput proto.InternalMessageInfo

func (m *ProblemInput) GetDatasetId() string {
	if m != nil {
		return m.DatasetId
	}
	return ""
}

func (m *ProblemInput) GetTargets() []*ProblemTarget {
	if m != nil {
		return m.Targets
	}
	return nil
}

func (m *ProblemInput) GetPrivilegedData() []*ProblemPrivilegedData {
	if m != nil {
		return m.PrivilegedData
	}
	return nil
}

func (m *ProblemInput) GetForecastingHorizon() *ForecastingHorizon {
	if m != nil {
		return m.ForecastingHorizon
	}
	return nil
}

type DataAugmentation struct {
	Domain               []string `protobuf:"bytes,1,rep,name=domain,proto3" json:"domain,omitempty"`
	Keywords             []string `protobuf:"bytes,2,rep,name=keywords,proto3" json:"keywords,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *DataAugmentation) Reset()         { *m = DataAugmentation{} }
func (m *DataAugmentation) String() string { return proto.CompactTextString(m) }
func (*DataAugmentation) ProtoMessage()    {}
func (*DataAugmentation) Descriptor() ([]byte, []int) {
	return fileDescriptor_b319862c9661813c, []int{6}
}

func (m *DataAugmentation) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DataAugmentation.Unmarshal(m, b)
}
func (m *DataAugmentation) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DataAugmentation.Marshal(b, m, deterministic)
}
func (m *DataAugmentation) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DataAugmentation.Merge(m, src)
}
func (m *DataAugmentation) XXX_Size() int {
	return xxx_messageInfo_DataAugmentation.Size(m)
}
func (m *DataAugmentation) XXX_DiscardUnknown() {
	xxx_messageInfo_DataAugmentation.DiscardUnknown(m)
}

var xxx_messageInfo_DataAugmentation proto.InternalMessageInfo

func (m *DataAugmentation) GetDomain() []string {
	if m != nil {
		return m.Domain
	}
	return nil
}

func (m *DataAugmentation) GetKeywords() []string {
	if m != nil {
		return m.Keywords
	}
	return nil
}

// Problem description matches the parsed problem description by
// the d3m_metadata.problem.parse_problem_description Python method.
// Problem outputs are not necessary for the purpose of this API
// and are needed only when executing an exported pipeline, but then
// TA2 gets full problem description anyway directly.
type ProblemDescription struct {
	Problem *Problem        `protobuf:"bytes,1,opt,name=problem,proto3" json:"problem,omitempty"`
	Inputs  []*ProblemInput `protobuf:"bytes,2,rep,name=inputs,proto3" json:"inputs,omitempty"`
	// ID of this problem. Required.
	Id string `protobuf:"bytes,3,opt,name=id,proto3" json:"id,omitempty"`
	// Version of this problem.
	Version              string              `protobuf:"bytes,4,opt,name=version,proto3" json:"version,omitempty"`
	Name                 string              `protobuf:"bytes,5,opt,name=name,proto3" json:"name,omitempty"`
	Description          string              `protobuf:"bytes,6,opt,name=description,proto3" json:"description,omitempty"`
	Digest               string              `protobuf:"bytes,7,opt,name=digest,proto3" json:"digest,omitempty"`
	DataAugmentation     []*DataAugmentation `protobuf:"bytes,8,rep,name=data_augmentation,json=dataAugmentation,proto3" json:"data_augmentation,omitempty"`
	OtherNames           []string            `protobuf:"bytes,9,rep,name=other_names,json=otherNames,proto3" json:"other_names,omitempty"`
	XXX_NoUnkeyedLiteral struct{}            `json:"-"`
	XXX_unrecognized     []byte              `json:"-"`
	XXX_sizecache        int32               `json:"-"`
}

func (m *ProblemDescription) Reset()         { *m = ProblemDescription{} }
func (m *ProblemDescription) String() string { return proto.CompactTextString(m) }
func (*ProblemDescription) ProtoMessage()    {}
func (*ProblemDescription) Descriptor() ([]byte, []int) {
	return fileDescriptor_b319862c9661813c, []int{7}
}

func (m *ProblemDescription) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ProblemDescription.Unmarshal(m, b)
}
func (m *ProblemDescription) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ProblemDescription.Marshal(b, m, deterministic)
}
func (m *ProblemDescription) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ProblemDescription.Merge(m, src)
}
func (m *ProblemDescription) XXX_Size() int {
	return xxx_messageInfo_ProblemDescription.Size(m)
}
func (m *ProblemDescription) XXX_DiscardUnknown() {
	xxx_messageInfo_ProblemDescription.DiscardUnknown(m)
}

var xxx_messageInfo_ProblemDescription proto.InternalMessageInfo

func (m *ProblemDescription) GetProblem() *Problem {
	if m != nil {
		return m.Problem
	}
	return nil
}

func (m *ProblemDescription) GetInputs() []*ProblemInput {
	if m != nil {
		return m.Inputs
	}
	return nil
}

func (m *ProblemDescription) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *ProblemDescription) GetVersion() string {
	if m != nil {
		return m.Version
	}
	return ""
}

func (m *ProblemDescription) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *ProblemDescription) GetDescription() string {
	if m != nil {
		return m.Description
	}
	return ""
}

func (m *ProblemDescription) GetDigest() string {
	if m != nil {
		return m.Digest
	}
	return ""
}

func (m *ProblemDescription) GetDataAugmentation() []*DataAugmentation {
	if m != nil {
		return m.DataAugmentation
	}
	return nil
}

func (m *ProblemDescription) GetOtherNames() []string {
	if m != nil {
		return m.OtherNames
	}
	return nil
}

func init() {
	proto.RegisterEnum("TaskKeyword", TaskKeyword_name, TaskKeyword_value)
	proto.RegisterEnum("PerformanceMetric", PerformanceMetric_name, PerformanceMetric_value)
	proto.RegisterType((*ProblemPerformanceMetric)(nil), "ProblemPerformanceMetric")
	proto.RegisterType((*Problem)(nil), "Problem")
	proto.RegisterType((*ProblemTarget)(nil), "ProblemTarget")
	proto.RegisterType((*ProblemPrivilegedData)(nil), "ProblemPrivilegedData")
	proto.RegisterType((*ForecastingHorizon)(nil), "ForecastingHorizon")
	proto.RegisterType((*ProblemInput)(nil), "ProblemInput")
	proto.RegisterType((*DataAugmentation)(nil), "DataAugmentation")
	proto.RegisterType((*ProblemDescription)(nil), "ProblemDescription")
}

func init() { proto.RegisterFile("problem.proto", fileDescriptor_b319862c9661813c) }

var fileDescriptor_b319862c9661813c = []byte{
	// 1287 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xbc, 0x56, 0x4b, 0x72, 0xdb, 0x46,
	0x13, 0xfe, 0x41, 0x4a, 0x7c, 0x34, 0x1f, 0x1a, 0x0d, 0x2d, 0x19, 0x7e, 0xfd, 0xa6, 0x99, 0x4a,
	0xa2, 0xf2, 0x82, 0x55, 0x56, 0xf6, 0x49, 0x8d, 0x80, 0xa1, 0x34, 0x16, 0x1e, 0xcc, 0x00, 0x60,
	0x2c, 0x6f, 0xa6, 0x20, 0x12, 0x96, 0x51, 0x22, 0x09, 0x06, 0x80, 0x9c, 0xc7, 0x09, 0x72, 0x8a,
	0xdc, 0x20, 0x55, 0x59, 0x66, 0x99, 0x93, 0xe4, 0x06, 0xb9, 0x43, 0x6a, 0x06, 0xa0, 0x1e, 0x94,
	0x53, 0xde, 0x65, 0x37, 0xfd, 0x75, 0x4f, 0xf7, 0xd7, 0x5f, 0x37, 0x86, 0x84, 0xce, 0x2a, 0x4d,
	0xce, 0xe7, 0xd1, 0x62, 0xb8, 0x4a, 0x93, 0x3c, 0x19, 0x7c, 0x0f, 0xfa, 0xb8, 0x00, 0xc6, 0x51,
	0xfa, 0x2e, 0x49, 0x17, 0xe1, 0x72, 0x1a, 0xd9, 0x51, 0x9e, 0xc6, 0x53, 0xfc, 0x12, 0x6a, 0x0b,
	0x75, 0xd2, 0xb5, 0xbe, 0x76, 0xd0, 0x3d, 0xc4, 0xc3, 0x7b, 0x31, 0xbc, 0x8c, 0xc0, 0x6d, 0xd0,
	0x2e, 0xf5, 0x4a, 0x5f, 0x3b, 0xd8, 0xe6, 0xda, 0x25, 0x7e, 0x02, 0xcd, 0x55, 0x92, 0x89, 0x79,
	0x78, 0x1e, 0xcd, 0xf5, 0x6a, 0x5f, 0x3b, 0x68, 0xf2, 0xc6, 0x2a, 0xc9, 0x2c, 0x69, 0x0f, 0x7e,
	0xd1, 0xa0, 0x5e, 0xd6, 0xc4, 0xaf, 0xa0, 0x93, 0x87, 0xd9, 0xa5, 0xb8, 0x8c, 0x7e, 0xfa, 0x21,
	0x49, 0x67, 0x99, 0xde, 0xe8, 0x57, 0x0f, 0xba, 0x87, 0xed, 0xa1, 0x1f, 0x66, 0x97, 0xa7, 0x05,
	0xc8, 0xdb, 0xf9, 0x8d, 0x91, 0xe1, 0xd7, 0xd0, 0x5b, 0xdd, 0xd0, 0x10, 0x45, 0xfd, 0x4c, 0xaf,
	0xf7, 0xab, 0x07, 0xad, 0xc3, 0x47, 0xc3, 0x7f, 0xeb, 0x86, 0xe3, 0xd5, 0x26, 0x94, 0x0d, 0xfe,
	0xd4, 0xa0, 0x53, 0x5e, 0xf0, 0xc3, 0xf4, 0x22, 0xca, 0xf1, 0x0b, 0x68, 0xe7, 0xea, 0x24, 0xe2,
	0xe5, 0x2c, 0xfa, 0x51, 0x75, 0xbe, 0xcd, 0x5b, 0x05, 0xc6, 0x24, 0x84, 0x9f, 0x43, 0x2b, 0x8d,
	0xb2, 0xe4, 0x2a, 0x9d, 0x46, 0x22, 0x9e, 0xa9, 0xa6, 0x9b, 0x1c, 0xd6, 0x10, 0x9b, 0xc9, 0x1c,
	0xd3, 0x64, 0x7e, 0xb5, 0x58, 0x96, 0x39, 0xaa, 0x45, 0x8e, 0x02, 0xbb, 0xce, 0x51, 0x86, 0x2c,
	0xc3, 0x45, 0xa4, 0x6f, 0x15, 0x39, 0x0a, 0xc8, 0x09, 0x17, 0x11, 0xfe, 0x12, 0x76, 0xa6, 0xf3,
	0xab, 0x2c, 0x8f, 0xd2, 0x4c, 0x2c, 0xaf, 0x16, 0xe7, 0x51, 0xaa, 0x6f, 0xab, 0x34, 0xdd, 0x35,
	0xec, 0x28, 0x74, 0xf0, 0xbb, 0x06, 0x7b, 0xeb, 0x9e, 0xd3, 0xf8, 0x43, 0x3c, 0x8f, 0x2e, 0xa2,
	0x99, 0x19, 0xe6, 0x21, 0x3e, 0x84, 0xbd, 0xd5, 0x35, 0x22, 0x66, 0x61, 0x1e, 0xde, 0xe9, 0xa9,
	0xb7, 0xba, 0x13, 0xfe, 0xdf, 0xf5, 0x36, 0xf8, 0x55, 0x03, 0x3c, 0x4a, 0xd2, 0x68, 0x1a, 0x66,
	0x79, 0xbc, 0xbc, 0x38, 0x49, 0xd2, 0xf8, 0xe7, 0x64, 0xb9, 0x59, 0x5b, 0xfb, 0x64, 0xed, 0xca,
	0x27, 0x6b, 0x57, 0xef, 0xe9, 0xfa, 0x19, 0x74, 0xde, 0x17, 0xf5, 0xc4, 0x87, 0x70, 0x7e, 0x55,
	0xd0, 0xd3, 0x78, 0xbb, 0x04, 0x27, 0x12, 0x1b, 0xfc, 0xa5, 0x41, 0xbb, 0xd4, 0x94, 0x2d, 0x57,
	0x57, 0x39, 0x7e, 0x06, 0x20, 0xf5, 0xcb, 0xe4, 0x5a, 0xac, 0x99, 0x35, 0x4b, 0x84, 0xcd, 0xf0,
	0x01, 0xd4, 0x8b, 0x05, 0xc9, 0xf4, 0x8a, 0x5a, 0xc3, 0xee, 0xf0, 0xce, 0x56, 0xf1, 0xb5, 0x1b,
	0x7f, 0x03, 0x3b, 0x1b, 0x33, 0xd1, 0xab, 0xea, 0xc6, 0xfe, 0xf0, 0xa3, 0x43, 0xe4, 0xdd, 0xbb,
	0x53, 0xc2, 0x26, 0xf4, 0xde, 0xdd, 0x48, 0x27, 0x4a, 0xda, 0xaa, 0x8b, 0xd6, 0x61, 0x6f, 0x78,
	0x5f, 0x56, 0x8e, 0xdf, 0xdd, 0xc3, 0x06, 0x23, 0x40, 0x32, 0x1b, 0xb9, 0xba, 0x58, 0x44, 0xcb,
	0x3c, 0xcc, 0xe3, 0x64, 0x89, 0xf7, 0xa1, 0x36, 0x4b, 0x16, 0x61, 0xbc, 0xd4, 0xb5, 0x7e, 0xf5,
	0xa0, 0xc9, 0x4b, 0x0b, 0x3f, 0x86, 0xc6, 0xf5, 0xd7, 0x59, 0x51, 0x9e, 0x6b, 0x7b, 0xf0, 0x47,
	0x05, 0x70, 0xc9, 0xdb, 0x8c, 0xb2, 0x69, 0x1a, 0xaf, 0x54, 0xaa, 0x01, 0xd4, 0xcb, 0x57, 0x46,
	0x69, 0xd5, 0x3a, 0x6c, 0xac, 0xbb, 0xe3, 0x6b, 0x07, 0xfe, 0x1c, 0x6a, 0xb1, 0xd4, 0x76, 0x2d,
	0x59, 0x67, 0x78, 0x5b, 0x71, 0x5e, 0x3a, 0x71, 0x17, 0x2a, 0xf1, 0xac, 0x9c, 0x63, 0x25, 0x9e,
	0x61, 0x1d, 0xea, 0x1f, 0xa2, 0x34, 0x8b, 0xcb, 0x9e, 0x9b, 0x7c, 0x6d, 0x62, 0x0c, 0x5b, 0x6a,
	0xe6, 0xdb, 0x0a, 0x56, 0x67, 0xdc, 0x87, 0xd6, 0xec, 0x86, 0x97, 0x5e, 0x53, 0xae, 0xdb, 0x90,
	0xea, 0x3a, 0xbe, 0x88, 0xb2, 0x5c, 0xaf, 0x2b, 0x67, 0x69, 0xe1, 0xaf, 0x61, 0x57, 0x7d, 0x31,
	0xe1, 0x2d, 0x89, 0xd4, 0xe3, 0xd4, 0x3a, 0xdc, 0x1d, 0x6e, 0x6a, 0xc7, 0xd1, 0x6c, 0x53, 0xcd,
	0xe7, 0xd0, 0x4a, 0xf2, 0xf7, 0x51, 0xaa, 0xf6, 0x30, 0xd3, 0x9b, 0x4a, 0x38, 0x50, 0x90, 0xdc,
	0xc3, 0xec, 0xe5, 0xdf, 0x5b, 0xd0, 0xba, 0xf5, 0xc8, 0xe1, 0xc7, 0xb0, 0xef, 0x13, 0xef, 0x54,
	0x9c, 0xd2, 0xb3, 0xef, 0x5c, 0x6e, 0x8a, 0xc0, 0x31, 0xe9, 0x88, 0x39, 0xd4, 0x44, 0xff, 0xc3,
	0x18, 0xba, 0x86, 0x45, 0x3c, 0x8f, 0x8d, 0x98, 0x41, 0x7c, 0xe6, 0x3a, 0x48, 0xc3, 0x5d, 0x00,
	0x4e, 0x8f, 0x39, 0xf5, 0x3c, 0x69, 0x57, 0xa4, 0x6d, 0x58, 0x81, 0xe7, 0x53, 0xce, 0x9c, 0x63,
	0x54, 0xc5, 0x3d, 0xd8, 0xb1, 0x98, 0x73, 0x2a, 0xc6, 0x9c, 0x9a, 0xcc, 0x50, 0x97, 0xb6, 0xf0,
	0x1e, 0xec, 0x4e, 0x28, 0xf7, 0xe9, 0x1b, 0xe1, 0xb8, 0x36, 0x73, 0x8a, 0x5c, 0xdb, 0xf8, 0x11,
	0xec, 0x95, 0xf0, 0x46, 0x99, 0x1a, 0x7e, 0x08, 0x3d, 0xc3, 0xb5, 0xed, 0xc0, 0x61, 0xfe, 0x99,
	0x30, 0xa9, 0x4f, 0x8b, 0x54, 0x75, 0xc9, 0xe9, 0x98, 0x93, 0xf1, 0x89, 0xb0, 0x89, 0x6f, 0x9c,
	0xc8, 0x9a, 0x0d, 0xbc, 0x03, 0xad, 0x91, 0xcb, 0xa9, 0x41, 0x3c, 0x5f, 0x02, 0x4d, 0xfc, 0x04,
	0x1e, 0x1a, 0xae, 0x65, 0x91, 0x23, 0x97, 0x13, 0x9f, 0x4d, 0xa8, 0x18, 0x31, 0xab, 0x64, 0x08,
	0xf8, 0x01, 0x20, 0xf7, 0xe8, 0x35, 0x35, 0xfc, 0x5b, 0x79, 0x5b, 0x32, 0xaf, 0x47, 0x6d, 0xe6,
	0x05, 0x63, 0xca, 0x27, 0xcc, 0xa3, 0x26, 0x6a, 0x63, 0x80, 0xda, 0x11, 0x73, 0x08, 0x3f, 0x43,
	0x1d, 0xd9, 0xa7, 0x1d, 0x58, 0x3e, 0x53, 0x4c, 0x51, 0xf7, 0xda, 0xb6, 0xc8, 0x11, 0xb5, 0xd0,
	0x8e, 0xb4, 0x03, 0x87, 0x4d, 0x08, 0x67, 0xc4, 0xa7, 0x08, 0x61, 0x04, 0x6d, 0xe5, 0x5f, 0x23,
	0xbb, 0x92, 0xa5, 0x3b, 0xa1, 0xdc, 0x22, 0xe3, 0xb1, 0x24, 0x82, 0x65, 0x49, 0xc7, 0x75, 0x6e,
	0x63, 0x3d, 0xdc, 0x82, 0xba, 0x4f, 0x8e, 0x02, 0x8b, 0x70, 0xf4, 0xa0, 0xd0, 0xda, 0x52, 0x92,
	0x10, 0x0b, 0xed, 0xe1, 0x26, 0x6c, 0x33, 0x9b, 0x1c, 0x53, 0xb4, 0x2f, 0x8f, 0x24, 0x30, 0x99,
	0x8b, 0x1e, 0xca, 0xe3, 0x84, 0x99, 0xd4, 0x45, 0xba, 0x24, 0xec, 0x8d, 0x29, 0x35, 0x4e, 0xd0,
	0x23, 0xdc, 0x80, 0x2d, 0x9f, 0xbe, 0xf1, 0xd1, 0x63, 0x19, 0xa0, 0x24, 0x43, 0x4f, 0xae, 0x59,
	0x17, 0xf6, 0x53, 0xc9, 0xc9, 0x67, 0x36, 0x15, 0x1e, 0xe5, 0x8c, 0x7a, 0xe8, 0x99, 0xac, 0x7f,
	0xcc, 0xdd, 0x60, 0x4c, 0x4d, 0xf4, 0x7f, 0x19, 0x7d, 0x4c, 0x5d, 0x6f, 0x4c, 0x7c, 0x46, 0x2c,
	0xf4, 0x5c, 0x12, 0xe6, 0xd4, 0x76, 0x7d, 0x19, 0xef, 0x78, 0x92, 0x70, 0x5f, 0x96, 0xb1, 0x82,
	0x31, 0x43, 0x2f, 0xa4, 0xae, 0x36, 0xf3, 0x24, 0x2c, 0x6c, 0xea, 0x13, 0x93, 0xf8, 0x04, 0x0d,
	0x5e, 0xfe, 0x56, 0x85, 0xdd, 0xfb, 0x3f, 0xf1, 0x32, 0x96, 0xfa, 0x9c, 0x19, 0x77, 0xf6, 0xad,
	0x0d, 0x0d, 0x62, 0x18, 0x01, 0x27, 0xc6, 0x19, 0xd2, 0x70, 0x07, 0x9a, 0x63, 0x4e, 0x0d, 0x56,
	0x2e, 0x1a, 0x40, 0x4d, 0x8e, 0xd8, 0xb2, 0x50, 0x15, 0xd7, 0xa0, 0x32, 0x7a, 0x85, 0xb6, 0xe4,
	0x85, 0xd1, 0x2b, 0x61, 0x33, 0x83, 0xbb, 0x68, 0x7b, 0x6d, 0x11, 0x69, 0xd5, 0x64, 0x27, 0xdc,
	0x35, 0x04, 0x09, 0x0c, 0x54, 0xc7, 0xbb, 0xd0, 0x29, 0x8d, 0x32, 0xba, 0x71, 0x07, 0x52, 0x57,
	0x9a, 0x78, 0x1f, 0xb0, 0x4d, 0x89, 0x23, 0xbc, 0x6f, 0x03, 0xc2, 0xa9, 0x29, 0x28, 0xe7, 0x2e,
	0x47, 0x20, 0xd7, 0x89, 0xbb, 0xae, 0x2f, 0x3e, 0xe2, 0x6c, 0xc9, 0x4d, 0x55, 0x38, 0x39, 0xf2,
	0x5c, 0x2b, 0xf0, 0x69, 0xe9, 0x68, 0x4b, 0xfe, 0x7c, 0x1d, 0x8d, 0x3a, 0xf8, 0x05, 0x3c, 0x73,
	0x5c, 0x6e, 0x13, 0x8b, 0xbd, 0xa5, 0xa6, 0xb0, 0x03, 0x3f, 0x20, 0x96, 0x60, 0xce, 0x48, 0x62,
	0x6a, 0x07, 0xbb, 0xf8, 0x29, 0xe8, 0xaf, 0x89, 0x61, 0x10, 0x6e, 0x0a, 0x8f, 0xd9, 0xcc, 0x22,
	0x5c, 0x6e, 0xbf, 0x67, 0xb8, 0x9c, 0xa2, 0x1d, 0xc9, 0xee, 0x5a, 0x0f, 0x41, 0x7c, 0xe1, 0xbb,
	0x63, 0x71, 0x8a, 0x76, 0xf1, 0x17, 0x30, 0xd8, 0xdc, 0x67, 0x41, 0x26, 0x94, 0x93, 0x63, 0x2a,
	0x6e, 0x04, 0xc4, 0x72, 0x23, 0x4f, 0x88, 0x6d, 0xcb, 0xf9, 0x58, 0xae, 0xe7, 0xa1, 0x9e, 0x9c,
	0x1d, 0x27, 0xce, 0x29, 0x9a, 0xaa, 0x29, 0x4a, 0x6c, 0x76, 0x04, 0x6f, 0x1b, 0xab, 0x78, 0x15,
	0xcd, 0xe3, 0x65, 0x74, 0x5e, 0x53, 0xff, 0xd5, 0xbe, 0xfa, 0x27, 0x00, 0x00, 0xff, 0xff, 0x60,
	0xc8, 0xe7, 0x5d, 0xbc, 0x09, 0x00, 0x00,
}
