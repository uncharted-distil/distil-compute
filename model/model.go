package model

import (
	"fmt"
	"regexp"
	"time"

	"github.com/jeffail/gabs"
)

const (
	// D3MIndexName is the variable name for the d3m index column
	D3MIndexName = "d3mIndex"
	// SchemaSourceClassification was loaded via classification
	SchemaSourceClassification = "classification"
	// SchemaSourceMerged was loaded via merged output
	SchemaSourceMerged = "merged"
	// SchemaSourceOriginal was loaded via original schema
	SchemaSourceOriginal = "original"
	// SchemaSourceRaw was loaded via raw data file
	SchemaSourceRaw = "raw"
	// VarRoleData is the distil role for data variables
	VarRoleData = "data"
	// VarRoleMetadata is the distil role for metadata variables
	VarRoleMetadata = "metadata"

	variableNameSizeLimit = 50

	// DefaultVarType is the variable type used by default
	DefaultVarType = "unknown"
	// ResTypeAudio is an audio data resource type
	ResTypeAudio = "audio"
	// ResTypeImage is an image data resource type
	ResTypeImage = "image"
	// ResTypeTable is a table data resource type
	ResTypeTable = "table"
	// ResTypeText is a text data resource type
	ResTypeText = "text"
	// ResTypeTime is a time series data resource type
	ResTypeTime = "timeseries"
	// ResTypeTime is a raw data resource type
	ResTypeRaw = "raw"

	// FeatureTypeTrain is the training feature type.
	FeatureTypeTrain = "train"
	// FeatureTypeTarget is the target feature type.
	FeatureTypeTarget = "target"
	// RoleIndex is the role used for index fields.
	RoleIndex = "index"
	// D3MIndexFieldName denotes the name of the index field.
	D3MIndexFieldName = "d3mIndex"
	// FeatureVarPrefix is the prefix of a metadata var name.
	FeatureVarPrefix = "_feature_"
	// ClusterVarPrefix is the prefix of a metadata var name.
	ClusterVarPrefix = "_cluster_"

	// Variables is the field name which stores the variables in elasticsearch.
	Variables = "variables"
	// VarNameField is the field name for the variable name.
	VarNameField = "colName"
	// VarIndexField is the field name for the variable index.
	VarIndexField = "colIndex"
	// VarRoleField is the field name for the variable role.
	VarRoleField = "role"
	// VarSelectedRoleField is the field name for the selected variable role.
	VarSelectedRoleField = "selectedRole"
	// VarDisplayVariableField is the field name for the display variable.
	VarDisplayVariableField = "colDisplayName"
	// VarOriginalVariableField is the field name for the original variable.
	VarOriginalVariableField = "colOriginalName"
	// VarTypeField is the field name for the variable type.
	VarTypeField = "colType"
	// VarOriginalTypeField is the field name for the orginal variable type.
	VarOriginalTypeField = "colOriginalType"
	// VarImportanceField is the field name for the variable importnace.
	VarImportanceField = "importance"
	// VarSuggestedTypesField is the field name for the suggested variable types.
	VarSuggestedTypesField = "suggestedTypes"
	// VarRoleIndex is the variable role of an index field.
	VarRoleIndex = "index"
	// VarDistilRole is the variable role in distil.
	VarDistilRole = "distilRole"
	// VarDeleted flags whether the variable is deleted.
	VarDeleted = "deleted"

	// TypeTypeField is the type field of a suggested type
	TypeTypeField = "type"
	// TypeProbabilityField is the probability field of a suggested type
	TypeProbabilityField = "probability"
	// TypeProvenanceField is the provenance field of a suggested type
	TypeProvenanceField = "provenance"

	// Database data types
	dataTypeText    = "TEXT"
	dataTypeDouble  = "double precision"
	dataTypeFloat   = "FLOAT8"
	dataTypeVector  = "FLOAT8[]"
	dataTypeInteger = "INTEGER"
	dataTypeDate    = "TIMESTAMP"
	dateFormat      = "2006-01-02T15:04:05Z"
)

var (
	nameRegex = regexp.MustCompile("[^a-zA-Z0-9]")
)

// Variable represents a single variable description.
type Variable struct {
	Name             string                 `json:"colName"`
	Type             string                 `json:"colType,omitempty"`
	OriginalType     string                 `json:"colOriginalType,omitempty"`
	SelectedRole     string                 `json:"selectedRole,omitempty"`
	Role             []string               `json:"role,omitempty"`
	DistilRole       string                 `json:"distilRole,omitempty"`
	OriginalVariable string                 `json:"colOriginalVariable"`
	OriginalName     string                 `json:"colOriginalName,omitempty"`
	DisplayName      string                 `json:"colDisplayName,omitempty"`
	Importance       int                    `json:"importance"`
	Index            int                    `json:"colIndex"`
	SuggestedTypes   []*SuggestedType       `json:"suggestedTypes,omitempty"`
	RefersTo         map[string]interface{} `json:"refersTo,omitempty"`
	Deleted          bool                   `json:"deleted"`
}

// DataResource represents a set of variables found in a data asset.
type DataResource struct {
	ResID        string      `json:"resID"`
	ResType      string      `json:"resType"`
	ResPath      string      `json:"resPath"`
	IsCollection bool        `json:"isCollection"`
	Variables    []*Variable `json:"columns,omitempty"`
	ResFormat    []string    `json:"resFormat"`
}

// SuggestedType represents a classified variable type.
type SuggestedType struct {
	Type        string  `json:"type"`
	Probability float64 `json:"probability"`
	Provenance  string  `json:"provenance"`
}

// Metadata represents a collection of dataset descriptions.
type Metadata struct {
	ID             string
	Name           string
	Description    string
	Summary        string
	SummaryMachine string
	Raw            bool
	DataResources  []*DataResource
	Schema         *gabs.Container
	Classification *gabs.Container
	NumRows        int64
	NumBytes       int64
	SchemaSource   string
	Redacted       bool
	DatasetFolder  string
}

// NewMetadata creates a new metadata instance.
func NewMetadata(id string, name string, description string) *Metadata {
	return &Metadata{
		ID:            id,
		Name:          name,
		Description:   description,
		DataResources: make([]*DataResource, 0),
	}
}

// NewDataResource creates a new data resource instance.
func NewDataResource(id string, typ string, format []string) *DataResource {
	return &DataResource{
		ResID:     id,
		ResType:   typ,
		ResFormat: format,
		Variables: make([]*Variable, 0),
	}
}

// NormalizeVariableName normalizes a variable name.
func NormalizeVariableName(name string) string {
	nameNormalized := nameRegex.ReplaceAllString(name, "_")
	if len(nameNormalized) > variableNameSizeLimit {
		nameNormalized = nameNormalized[:variableNameSizeLimit]
	}

	return nameNormalized
}

// NewVariable creates a new variable.
func NewVariable(index int, name, displayName, originalName, typ, originalType string, role []string, distilRole string, refersTo map[string]interface{}, existingVariables []*Variable, normalizeName bool) *Variable {
	normed := name
	if normalizeName {
		// normalize name
		normed = NormalizeVariableName(name)

		// normed name needs to be unique
		count := 0
		for _, v := range existingVariables {
			if v.Name == normed {
				count = count + 1
			}
		}
		if count > 0 {
			normed = fmt.Sprintf("%s_%d", normed, count)
		}
	}

	// select the first role by default.
	selectedRole := ""
	if len(role) > 0 {
		selectedRole = role[0]
	}
	if distilRole == "" {
		distilRole = VarRoleData
	}
	if originalName == "" {
		originalName = normed
	}

	if displayName == "" {
		displayName = name
	}

	return &Variable{
		Name:             normed,
		Index:            index,
		Type:             typ,
		OriginalType:     originalType,
		Role:             role,
		SelectedRole:     selectedRole,
		DistilRole:       distilRole,
		OriginalVariable: originalName,
		OriginalName:     normed,
		DisplayName:      displayName,
		RefersTo:         refersTo,
		SuggestedTypes:   make([]*SuggestedType, 0),
	}
}

// CanBeFeaturized determines if a data resource can be featurized.
func (dr *DataResource) CanBeFeaturized() bool {
	return dr.ResType == ResTypeImage
}

// AddVariable creates and add a new variable to the data resource.
func (dr *DataResource) AddVariable(name string, originalName string, typ string, role []string, distilRole string) {
	v := NewVariable(len(dr.Variables), name, "", originalName, typ, typ, role, distilRole, nil, dr.Variables, false)
	dr.Variables = append(dr.Variables, v)
}

// GetMainDataResource returns the data resource that contains the D3M index.
func (m *Metadata) GetMainDataResource() *DataResource {
	// main data resource has d3m index variable
	for _, dr := range m.DataResources {
		for _, v := range dr.Variables {
			if v.Name == D3MIndexName {
				return dr
			}
		}
	}

	return nil
}

// GenerateHeaders generates csv headers for the data resources.
func (m *Metadata) GenerateHeaders() ([][]string, error) {
	// each data resource needs a separate header
	headers := make([][]string, len(m.DataResources))

	for index, dr := range m.DataResources {
		header := dr.GenerateHeader()
		headers[index] = header
	}

	return headers, nil
}

// GenerateHeader generates csv headers for the data resource.
func (dr *DataResource) GenerateHeader() []string {
	header := make([]string, len(dr.Variables))

	// iterate over the fields
	for hIndex, field := range dr.Variables {
		header[hIndex] = field.Name
	}

	return header
}

// IsMediaReference returns true if a variable is a reference to a media resource.
func (v *Variable) IsMediaReference() bool {
	// if refers to has a res object of string, assume media reference`
	mediaReference := false
	if v.RefersTo != nil {
		if v.RefersTo["resObject"] != nil {
			_, ok := v.RefersTo["resObject"].(string)
			if ok {
				mediaReference = true
			}
		}
	}
	return mediaReference
}

// MapD3MTypeToPostgresType generates a postgres type from a d3m type.
func MapD3MTypeToPostgresType(typ string) string {
	// Integer types can be returned as floats.
	switch typ {
	case IndexType:
		return dataTypeInteger
	case IntegerType, FloatType, LongitudeType, LatitudeType, RealType:
		return dataTypeFloat
	case OrdinalType, CategoricalType, TextType:
		return dataTypeText
	case DateTimeType:
		return dataTypeDate
	case RealVectorType:
		return dataTypeVector
	default:
		return dataTypeText
	}
}

// DefaultPostgresValueFromD3MType generates a default postgres value from a d3m type.
func DefaultPostgresValueFromD3MType(typ string) interface{} {
	switch typ {
	case IndexType:
		return float64(0)
	case FloatType, LongitudeType, LatitudeType, RealType:
		return float64(0)
	case IntegerType:
		return int(0)
	case DateTimeType:
		return fmt.Sprintf("'%s'", time.Time{}.Format(dateFormat))
	case RealVectorType:
		return "'{}'"
	default:
		return "''"
	}
}
