// Code generated by protoc-gen-go. DO NOT EDIT.
// source: problem.proto

package pipeline

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	_ "github.com/golang/protobuf/protoc-gen-go/descriptor"
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
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

// Top level classification of the problem.
type TaskType int32

const (
	// Default value. Not to be used.
	TaskType_TASK_TYPE_UNDEFINED     TaskType = 0
	TaskType_CLASSIFICATION          TaskType = 1
	TaskType_REGRESSION              TaskType = 2
	TaskType_CLUSTERING              TaskType = 3
	TaskType_LINK_PREDICTION         TaskType = 4
	TaskType_VERTEX_NOMINATION       TaskType = 5
	TaskType_COMMUNITY_DETECTION     TaskType = 6
	TaskType_GRAPH_CLUSTERING        TaskType = 7
	TaskType_GRAPH_MATCHING          TaskType = 8
	TaskType_TIME_SERIES_FORECASTING TaskType = 9
	TaskType_COLLABORATIVE_FILTERING TaskType = 10
	TaskType_OBJECT_DETECTION        TaskType = 11
)

var TaskType_name = map[int32]string{
	0:  "TASK_TYPE_UNDEFINED",
	1:  "CLASSIFICATION",
	2:  "REGRESSION",
	3:  "CLUSTERING",
	4:  "LINK_PREDICTION",
	5:  "VERTEX_NOMINATION",
	6:  "COMMUNITY_DETECTION",
	7:  "GRAPH_CLUSTERING",
	8:  "GRAPH_MATCHING",
	9:  "TIME_SERIES_FORECASTING",
	10: "COLLABORATIVE_FILTERING",
	11: "OBJECT_DETECTION",
}

var TaskType_value = map[string]int32{
	"TASK_TYPE_UNDEFINED":     0,
	"CLASSIFICATION":          1,
	"REGRESSION":              2,
	"CLUSTERING":              3,
	"LINK_PREDICTION":         4,
	"VERTEX_NOMINATION":       5,
	"COMMUNITY_DETECTION":     6,
	"GRAPH_CLUSTERING":        7,
	"GRAPH_MATCHING":          8,
	"TIME_SERIES_FORECASTING": 9,
	"COLLABORATIVE_FILTERING": 10,
	"OBJECT_DETECTION":        11,
}

func (x TaskType) String() string {
	return proto.EnumName(TaskType_name, int32(x))
}

func (TaskType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_b319862c9661813c, []int{0}
}

// Secondary classification of the problem.
type TaskSubtype int32

const (
	// Default value. Not to be used.
	TaskSubtype_TASK_SUBTYPE_UNDEFINED TaskSubtype = 0
	// No secondary task is applicable for this problem.
	TaskSubtype_NONE           TaskSubtype = 1
	TaskSubtype_BINARY         TaskSubtype = 2
	TaskSubtype_MULTICLASS     TaskSubtype = 3
	TaskSubtype_MULTILABEL     TaskSubtype = 4
	TaskSubtype_UNIVARIATE     TaskSubtype = 5
	TaskSubtype_MULTIVARIATE   TaskSubtype = 6
	TaskSubtype_OVERLAPPING    TaskSubtype = 7
	TaskSubtype_NONOVERLAPPING TaskSubtype = 8
)

var TaskSubtype_name = map[int32]string{
	0: "TASK_SUBTYPE_UNDEFINED",
	1: "NONE",
	2: "BINARY",
	3: "MULTICLASS",
	4: "MULTILABEL",
	5: "UNIVARIATE",
	6: "MULTIVARIATE",
	7: "OVERLAPPING",
	8: "NONOVERLAPPING",
}

var TaskSubtype_value = map[string]int32{
	"TASK_SUBTYPE_UNDEFINED": 0,
	"NONE":                   1,
	"BINARY":                 2,
	"MULTICLASS":             3,
	"MULTILABEL":             4,
	"UNIVARIATE":             5,
	"MULTIVARIATE":           6,
	"OVERLAPPING":            7,
	"NONOVERLAPPING":         8,
}

func (x TaskSubtype) String() string {
	return proto.EnumName(TaskSubtype_name, int32(x))
}

func (TaskSubtype) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_b319862c9661813c, []int{1}
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
	PerformanceMetric_ROOT_MEAN_SQUARED_ERROR_AVG        PerformanceMetric = 12
	PerformanceMetric_MEAN_ABSOLUTE_ERROR                PerformanceMetric = 13
	PerformanceMetric_R_SQUARED                          PerformanceMetric = 14
	PerformanceMetric_NORMALIZED_MUTUAL_INFORMATION      PerformanceMetric = 15
	PerformanceMetric_JACCARD_SIMILARITY_SCORE           PerformanceMetric = 16
	PerformanceMetric_PRECISION_AT_TOP_K                 PerformanceMetric = 17
	PerformanceMetric_OBJECT_DETECTION_AVERAGE_PRECISION PerformanceMetric = 18
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
	12:  "ROOT_MEAN_SQUARED_ERROR_AVG",
	13:  "MEAN_ABSOLUTE_ERROR",
	14:  "R_SQUARED",
	15:  "NORMALIZED_MUTUAL_INFORMATION",
	16:  "JACCARD_SIMILARITY_SCORE",
	17:  "PRECISION_AT_TOP_K",
	18:  "OBJECT_DETECTION_AVERAGE_PRECISION",
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
	"ROOT_MEAN_SQUARED_ERROR_AVG":        12,
	"MEAN_ABSOLUTE_ERROR":                13,
	"R_SQUARED":                          14,
	"NORMALIZED_MUTUAL_INFORMATION":      15,
	"JACCARD_SIMILARITY_SCORE":           16,
	"PRECISION_AT_TOP_K":                 17,
	"OBJECT_DETECTION_AVERAGE_PRECISION": 18,
	"LOSS":                               100,
}

func (x PerformanceMetric) String() string {
	return proto.EnumName(PerformanceMetric_name, int32(x))
}

func (PerformanceMetric) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_b319862c9661813c, []int{2}
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
	// TODO: Remove deprecated fields in a future version.
	Id                   string                      `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`                   // Deprecated: Do not use.
	Version              string                      `protobuf:"bytes,2,opt,name=version,proto3" json:"version,omitempty"`         // Deprecated: Do not use.
	Name                 string                      `protobuf:"bytes,3,opt,name=name,proto3" json:"name,omitempty"`               // Deprecated: Do not use.
	Description          string                      `protobuf:"bytes,4,opt,name=description,proto3" json:"description,omitempty"` // Deprecated: Do not use.
	TaskType             TaskType                    `protobuf:"varint,5,opt,name=task_type,json=taskType,proto3,enum=TaskType" json:"task_type,omitempty"`
	TaskSubtype          TaskSubtype                 `protobuf:"varint,6,opt,name=task_subtype,json=taskSubtype,proto3,enum=TaskSubtype" json:"task_subtype,omitempty"`
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

// Deprecated: Do not use.
func (m *Problem) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

// Deprecated: Do not use.
func (m *Problem) GetVersion() string {
	if m != nil {
		return m.Version
	}
	return ""
}

// Deprecated: Do not use.
func (m *Problem) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

// Deprecated: Do not use.
func (m *Problem) GetDescription() string {
	if m != nil {
		return m.Description
	}
	return ""
}

func (m *Problem) GetTaskType() TaskType {
	if m != nil {
		return m.TaskType
	}
	return TaskType_TASK_TYPE_UNDEFINED
}

func (m *Problem) GetTaskSubtype() TaskSubtype {
	if m != nil {
		return m.TaskSubtype
	}
	return TaskSubtype_TASK_SUBTYPE_UNDEFINED
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

type ProblemInput struct {
	// Should match one of input datasets given to the pipeline search.
	// Every "Dataset" object has an "id" associated with it and is available
	// in its metadata. That ID is then used here to reference those inputs.
	DatasetId string `protobuf:"bytes,1,opt,name=dataset_id,json=datasetId,proto3" json:"dataset_id,omitempty"`
	// Targets should resolve to columns in a given dataset.
	Targets              []*ProblemTarget `protobuf:"bytes,2,rep,name=targets,proto3" json:"targets,omitempty"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_unrecognized     []byte           `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *ProblemInput) Reset()         { *m = ProblemInput{} }
func (m *ProblemInput) String() string { return proto.CompactTextString(m) }
func (*ProblemInput) ProtoMessage()    {}
func (*ProblemInput) Descriptor() ([]byte, []int) {
	return fileDescriptor_b319862c9661813c, []int{3}
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
	return fileDescriptor_b319862c9661813c, []int{4}
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
	XXX_NoUnkeyedLiteral struct{}            `json:"-"`
	XXX_unrecognized     []byte              `json:"-"`
	XXX_sizecache        int32               `json:"-"`
}

func (m *ProblemDescription) Reset()         { *m = ProblemDescription{} }
func (m *ProblemDescription) String() string { return proto.CompactTextString(m) }
func (*ProblemDescription) ProtoMessage()    {}
func (*ProblemDescription) Descriptor() ([]byte, []int) {
	return fileDescriptor_b319862c9661813c, []int{5}
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

func init() {
	proto.RegisterEnum("TaskType", TaskType_name, TaskType_value)
	proto.RegisterEnum("TaskSubtype", TaskSubtype_name, TaskSubtype_value)
	proto.RegisterEnum("PerformanceMetric", PerformanceMetric_name, PerformanceMetric_value)
	proto.RegisterType((*ProblemPerformanceMetric)(nil), "ProblemPerformanceMetric")
	proto.RegisterType((*Problem)(nil), "Problem")
	proto.RegisterType((*ProblemTarget)(nil), "ProblemTarget")
	proto.RegisterType((*ProblemInput)(nil), "ProblemInput")
	proto.RegisterType((*DataAugmentation)(nil), "DataAugmentation")
	proto.RegisterType((*ProblemDescription)(nil), "ProblemDescription")
}

func init() { proto.RegisterFile("problem.proto", fileDescriptor_b319862c9661813c) }

var fileDescriptor_b319862c9661813c = []byte{
	// 1086 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x74, 0x55, 0xdd, 0x72, 0xdb, 0x36,
	0x13, 0xfd, 0xa8, 0x7f, 0xad, 0x64, 0x19, 0x46, 0xf2, 0x25, 0x6c, 0x7e, 0x26, 0x8a, 0xa6, 0x4d,
	0x3d, 0xbe, 0x50, 0x26, 0xe9, 0x7d, 0x67, 0x20, 0x0a, 0x72, 0x90, 0xf0, 0x47, 0x05, 0x49, 0xb7,
	0xce, 0x0d, 0x86, 0x36, 0x19, 0x0f, 0xc7, 0x92, 0xc8, 0x92, 0x54, 0xdb, 0xbc, 0x48, 0x7b, 0xd7,
	0xcb, 0x3e, 0x40, 0x9f, 0xa0, 0x8f, 0xd6, 0x01, 0x48, 0xca, 0xb2, 0xdd, 0xdc, 0x71, 0xcf, 0xd9,
	0x5d, 0xec, 0x1e, 0x1c, 0x48, 0x70, 0x90, 0x66, 0xc9, 0xc5, 0x2a, 0x5a, 0x4f, 0xd3, 0x2c, 0x29,
	0x92, 0x27, 0xe3, 0xab, 0x24, 0xb9, 0x5a, 0x45, 0xaf, 0x55, 0x74, 0xb1, 0xfd, 0xf4, 0x3a, 0x8c,
	0xf2, 0xcb, 0x2c, 0x4e, 0x8b, 0x24, 0x2b, 0x33, 0x26, 0x3f, 0x83, 0xbe, 0x2c, 0x4b, 0x96, 0x51,
	0xf6, 0x29, 0xc9, 0xd6, 0xc1, 0xe6, 0x32, 0xb2, 0xa2, 0x22, 0x8b, 0x2f, 0xf1, 0x09, 0x74, 0xd6,
	0xea, 0x4b, 0xd7, 0xc6, 0xda, 0xf1, 0xe8, 0x2d, 0x9e, 0xde, 0xcb, 0xe1, 0x55, 0x06, 0x1e, 0x82,
	0x76, 0xad, 0x37, 0xc6, 0xda, 0x71, 0x9b, 0x6b, 0xd7, 0xf8, 0x29, 0xf4, 0xd3, 0x24, 0x17, 0xab,
	0xe0, 0x22, 0x5a, 0xe9, 0xcd, 0xb1, 0x76, 0xdc, 0xe7, 0xbd, 0x34, 0xc9, 0x4d, 0x19, 0x4f, 0xfe,
	0x6c, 0x40, 0xb7, 0x3a, 0x13, 0x63, 0x68, 0xc4, 0xa1, 0x6a, 0xdf, 0x9f, 0x35, 0x74, 0x8d, 0x37,
	0xe2, 0x10, 0x3f, 0x83, 0xee, 0x2f, 0x51, 0x96, 0xc7, 0xc9, 0x46, 0x35, 0x2c, 0x89, 0x1a, 0xc2,
	0x8f, 0xa0, 0xb5, 0x09, 0xd6, 0x51, 0xd9, 0x55, 0x51, 0x2a, 0xc6, 0x5f, 0xc3, 0xa0, 0x5e, 0x4e,
	0x56, 0xb6, 0x76, 0xf4, 0x3e, 0x8c, 0x5f, 0x41, 0xbf, 0x08, 0xf2, 0x6b, 0x51, 0x7c, 0x4e, 0x23,
	0xbd, 0xad, 0xb6, 0xea, 0x4f, 0xbd, 0x20, 0xbf, 0xf6, 0x3e, 0xa7, 0x11, 0xef, 0x15, 0xd5, 0x17,
	0x7e, 0x0d, 0x43, 0x95, 0x97, 0x6f, 0x2f, 0x54, 0x6a, 0x47, 0xa5, 0x0e, 0x55, 0xaa, 0x5b, 0x62,
	0x7c, 0x50, 0xdc, 0x04, 0xf8, 0x3d, 0x3c, 0x48, 0x6f, 0xc4, 0x11, 0xa5, 0x2a, 0xb9, 0xde, 0x1d,
	0x37, 0x8f, 0x07, 0x6f, 0xbf, 0x9a, 0x7e, 0x49, 0x63, 0x8e, 0xd3, 0xbb, 0x50, 0x3e, 0xf9, 0x47,
	0x83, 0x83, 0xaa, 0xc0, 0x0b, 0xb2, 0xab, 0xa8, 0xc0, 0x2f, 0xe5, 0x38, 0xf2, 0x4b, 0xc4, 0x9b,
	0x30, 0xfa, 0x4d, 0x09, 0xd6, 0x96, 0x03, 0x48, 0x8c, 0x49, 0x08, 0xbf, 0x80, 0x41, 0x16, 0xe5,
	0xc9, 0x36, 0xbb, 0x8c, 0x44, 0x1c, 0x96, 0xca, 0x71, 0xa8, 0x21, 0x16, 0xca, 0x1e, 0x97, 0xc9,
	0x6a, 0xbb, 0xde, 0x54, 0x3d, 0x9a, 0x65, 0x8f, 0x12, 0xdb, 0xf5, 0xa8, 0x52, 0x94, 0xc4, 0xad,
	0xb2, 0x47, 0x09, 0xd9, 0x52, 0xe4, 0x6f, 0xe1, 0xf0, 0x72, 0xb5, 0xcd, 0x8b, 0x28, 0xcb, 0xc5,
	0x66, 0xbb, 0xbe, 0x88, 0x32, 0x25, 0x62, 0x9b, 0x8f, 0x6a, 0xd8, 0x56, 0xe8, 0xe4, 0x47, 0x18,
	0x56, 0x1b, 0xb0, 0x4d, 0xba, 0x2d, 0xf0, 0x73, 0x80, 0x30, 0x28, 0x82, 0x5c, 0x6e, 0x50, 0xdd,
	0x37, 0xef, 0x57, 0x08, 0x0b, 0xf1, 0x31, 0x74, 0xcb, 0x5d, 0x72, 0xbd, 0xa1, 0x14, 0x1b, 0x4d,
	0x6f, 0x09, 0xc0, 0x6b, 0x7a, 0xb2, 0x00, 0x34, 0x0f, 0x8a, 0x80, 0x6c, 0xaf, 0xd6, 0xd1, 0xa6,
	0x08, 0x8a, 0xd2, 0x12, 0x9d, 0x30, 0x59, 0x07, 0xf1, 0x46, 0xd7, 0xc6, 0xcd, 0xe3, 0x3e, 0xaf,
	0x22, 0xfc, 0x04, 0x7a, 0xd7, 0xd1, 0xe7, 0x5f, 0x93, 0x2c, 0x2c, 0xdb, 0xf6, 0xf9, 0x2e, 0x9e,
	0xfc, 0xd1, 0x00, 0x5c, 0x1d, 0x31, 0xdf, 0xf3, 0xc7, 0x04, 0xba, 0xd5, 0x0b, 0x52, 0x43, 0x0e,
	0xde, 0xf6, 0xea, 0x41, 0x78, 0x4d, 0xe0, 0x6f, 0xa0, 0x13, 0xcb, 0xa5, 0xea, 0x59, 0x0f, 0xa6,
	0xfb, 0xab, 0xf2, 0x8a, 0xc4, 0x23, 0x65, 0xed, 0xd2, 0xfc, 0xd2, 0xd6, 0xfa, 0x8d, 0xad, 0x4b,
	0x61, 0x77, 0x96, 0xc6, 0x95, 0xa5, 0xdb, 0x0a, 0x2e, 0xed, 0x3c, 0xbe, 0x6d, 0xe7, 0x8e, 0xa2,
	0x6e, 0x59, 0x59, 0x6e, 0x1d, 0x5f, 0x45, 0x79, 0xa1, 0x77, 0x15, 0x59, 0x45, 0xf8, 0x7b, 0x38,
	0x92, 0xc2, 0x8a, 0x60, 0x4f, 0x22, 0xbd, 0xa7, 0x26, 0x3d, 0x9a, 0xde, 0xd5, 0x8e, 0xa3, 0xf0,
	0x0e, 0x72, 0xf2, 0x7b, 0x03, 0x7a, 0xf5, 0x8b, 0xc0, 0x8f, 0xe1, 0x81, 0x47, 0xdc, 0x0f, 0xc2,
	0x3b, 0x5f, 0x52, 0xe1, 0xdb, 0x73, 0xba, 0x60, 0x36, 0x9d, 0xa3, 0xff, 0x61, 0x0c, 0x23, 0xc3,
	0x24, 0xae, 0xcb, 0x16, 0xcc, 0x20, 0x1e, 0x73, 0x6c, 0xa4, 0xe1, 0x11, 0x00, 0xa7, 0xa7, 0x9c,
	0xba, 0xae, 0x8c, 0x1b, 0x32, 0x36, 0x4c, 0xdf, 0xf5, 0x28, 0x67, 0xf6, 0x29, 0x6a, 0xe2, 0x07,
	0x70, 0x68, 0x32, 0xfb, 0x83, 0x58, 0x72, 0x3a, 0x67, 0x86, 0x2a, 0x6a, 0xe1, 0xff, 0xc3, 0xd1,
	0x19, 0xe5, 0x1e, 0xfd, 0x49, 0xd8, 0x8e, 0xc5, 0xec, 0xb2, 0x57, 0x5b, 0x1e, 0x6c, 0x38, 0x96,
	0xe5, 0xdb, 0xcc, 0x3b, 0x17, 0x73, 0xea, 0xd1, 0x32, 0xbf, 0x83, 0x1f, 0x02, 0x3a, 0xe5, 0x64,
	0xf9, 0x4e, 0xec, 0xb5, 0xee, 0xca, 0x71, 0x4a, 0xd4, 0x22, 0x9e, 0xf1, 0x4e, 0x62, 0x3d, 0xfc,
	0x14, 0x1e, 0x7b, 0xcc, 0xa2, 0xc2, 0xa5, 0x9c, 0x51, 0x57, 0x2c, 0x1c, 0x4e, 0x0d, 0xe2, 0x7a,
	0x92, 0xec, 0x4b, 0xd2, 0x70, 0x4c, 0x93, 0xcc, 0x1c, 0x4e, 0x3c, 0x76, 0x46, 0xc5, 0x82, 0x99,
	0x55, 0x37, 0x90, 0x67, 0x38, 0xb3, 0xf7, 0xd4, 0xf0, 0xf6, 0x4e, 0x1e, 0x9c, 0xfc, 0xa5, 0xc1,
	0x60, 0xef, 0xfd, 0xe3, 0x27, 0xf0, 0x48, 0x69, 0xe3, 0xfa, 0xb3, 0x7b, 0xf2, 0xf4, 0xa0, 0x65,
	0x3b, 0x36, 0x45, 0x1a, 0x06, 0xe8, 0xcc, 0x98, 0x4d, 0xf8, 0x79, 0x29, 0x88, 0xe5, 0x9b, 0x1e,
	0x53, 0xca, 0xa1, 0xe6, 0x2e, 0x36, 0xc9, 0x8c, 0x9a, 0xa8, 0x25, 0x63, 0xdf, 0x66, 0x67, 0x84,
	0x33, 0xe2, 0x51, 0xd4, 0xc6, 0x08, 0x86, 0x8a, 0xaf, 0x91, 0x0e, 0x3e, 0x84, 0x81, 0x73, 0x46,
	0xb9, 0x49, 0x96, 0xcb, 0xdd, 0xe2, 0xb6, 0x63, 0xef, 0x63, 0xbd, 0x93, 0xbf, 0x9b, 0x70, 0x74,
	0xff, 0xd7, 0xfc, 0x21, 0x20, 0x8b, 0x7a, 0x9c, 0x19, 0xb7, 0x06, 0x1d, 0x42, 0x8f, 0x18, 0x86,
	0xcf, 0x89, 0x71, 0x8e, 0x34, 0x7c, 0x00, 0xfd, 0x25, 0xa7, 0x06, 0xab, 0x2e, 0x10, 0xa0, 0x23,
	0x35, 0x33, 0x4d, 0xd4, 0xc4, 0x1d, 0x68, 0x2c, 0xde, 0xa0, 0x96, 0x2c, 0x58, 0xbc, 0x11, 0x16,
	0x33, 0xb8, 0x83, 0xda, 0x75, 0x44, 0x64, 0xd4, 0xc1, 0x03, 0xe8, 0x72, 0xc7, 0x10, 0xc4, 0x37,
	0x50, 0x17, 0x1f, 0xc1, 0x41, 0x15, 0x54, 0xd9, 0xbd, 0x5b, 0x90, 0x2a, 0xe9, 0xe3, 0x47, 0x80,
	0x2d, 0x4a, 0x6c, 0xe1, 0xfe, 0xe0, 0x13, 0x4e, 0xe7, 0x82, 0x72, 0xee, 0x70, 0x04, 0xf2, 0x7e,
	0xb8, 0xe3, 0x78, 0xe2, 0x3f, 0xc8, 0x01, 0x7e, 0x01, 0x4f, 0xbf, 0x40, 0x0a, 0x72, 0x76, 0x8a,
	0x86, 0xd2, 0x3d, 0x8a, 0x23, 0x33, 0xd7, 0x31, 0x7d, 0x8f, 0x56, 0x95, 0x07, 0x72, 0x41, 0x5e,
	0x57, 0xa0, 0x11, 0x7e, 0x09, 0xcf, 0x6d, 0x87, 0x5b, 0xc4, 0x64, 0x1f, 0xe9, 0x5c, 0x58, 0xbe,
	0xe7, 0x13, 0x53, 0x30, 0x7b, 0x21, 0x31, 0x75, 0xeb, 0x87, 0xf8, 0x19, 0xe8, 0xef, 0x89, 0x61,
	0x10, 0x3e, 0x17, 0x2e, 0xb3, 0x98, 0x49, 0xb8, 0x74, 0xa4, 0x6b, 0x38, 0x9c, 0x22, 0x24, 0xc7,
	0xdf, 0x09, 0x26, 0x88, 0x27, 0x3c, 0x67, 0x29, 0x3e, 0xa0, 0x23, 0xfc, 0x0a, 0x26, 0x77, 0x1d,
	0x24, 0xc8, 0x19, 0xe5, 0xe4, 0x94, 0x8a, 0x1b, 0x85, 0xb1, 0xf4, 0x89, 0xe9, 0xb8, 0x2e, 0x0a,
	0x67, 0xf0, 0xb1, 0x97, 0xc6, 0x69, 0xb4, 0x8a, 0x37, 0xd1, 0x45, 0x47, 0xfd, 0x37, 0x7f, 0xf7,
	0x6f, 0x00, 0x00, 0x00, 0xff, 0xff, 0x97, 0x96, 0xb0, 0x3c, 0xce, 0x07, 0x00, 0x00,
}
