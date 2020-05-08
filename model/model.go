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

package model

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/Jeffail/gabs/v2"
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

	variableNameSizeLimit = 50
	datasetIDSizeLimit    = 50

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
	// ResTypeRaw is a raw data resource type
	ResTypeRaw = "raw"

	// FeatureTypeTrain is the training feature type.
	FeatureTypeTrain = "train"
	// FeatureTypeTarget is the target feature type.
	FeatureTypeTarget = "target"
	// RoleIndex is the role used for index fields.
	RoleIndex = "index"
	// RoleMultiIndex is the role used for index fields which are not unique in the learning data.
	RoleMultiIndex = "multiIndex"
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
	// VarDescriptionField is the field name for the variable description.
	VarDescriptionField = "colDescription"
	// VarImportanceField is the field name for the variable importnace.
	VarImportanceField = "importance"
	// VarSuggestedTypesField is the field name for the suggested variable types.
	VarSuggestedTypesField = "suggestedTypes"
	// VarDistilRole is the variable role in distil.
	VarDistilRole = "distilRole"
	// VarDistilRoleIndex indicates a var has an index role in distil.
	VarDistilRoleIndex = "index"
	// VarDistilRoleData indicates a var has a data role in distil.
	VarDistilRoleData = "data"
	// VarDistilRoleGrouping indicates a var has a grouping role in distil.
	VarDistilRoleGrouping = "grouping"
	// VarDistilRoleMetadata is the distil role for metadata variables
	VarDistilRoleMetadata = "metadata"
	// VarDeleted flags whether the variable is deleted.
	VarDeleted = "deleted"
	// VarGroupingField is the field name for the variable grouping.
	VarGroupingField = "grouping"
	// VarMinField is the field name for the min value.
	VarMinField = "min"
	// VarMaxField is the field name for the max value.
	VarMaxField = "max"

	// TypeTypeField is the type field of a suggested type
	TypeTypeField = "type"
	// TypeProbabilityField is the probability field of a suggested type
	TypeProbabilityField = "probability"
	// TypeProvenanceField is the provenance field of a suggested type
	TypeProvenanceField = "provenance"

	// DatasetPrefix is the prefix used for a normalized dataset id
	DatasetPrefix = "d_"

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

// GroupingProperties represents a grouping properties.
type GroupingProperties struct {
	XCol       string `json:"xCol"`
	YCol       string `json:"yCol"`
	ClusterCol string `json:"clusterCol"`
}

// Grouping represents a variable grouping.
type Grouping struct {
	Dataset    string             `json:"dataset"`
	Type       string             `json:"type"`
	IDCol      string             `json:"idCol"`
	SubIDs     []string           `json:"subIds"`
	Hidden     []string           `json:"hidden"`
	Properties GroupingProperties `json:"properties"`
}

// Variable represents a single variable description.
type Variable struct {
	Name             string                 `json:"colName"`
	Type             string                 `json:"colType,omitempty"`
	Description      string                 `json:"colDescription,omitempty"`
	OriginalType     string                 `json:"colOriginalType,omitempty"`
	SelectedRole     string                 `json:"selectedRole,omitempty"`
	Role             []string               `json:"role,omitempty"`
	DistilRole       string                 `json:"distilRole,omitempty"`
	OriginalVariable string                 `json:"colOriginalName"`
	DisplayName      string                 `json:"colDisplayName,omitempty"`
	Importance       int                    `json:"importance"`
	Index            int                    `json:"colIndex"`
	SuggestedTypes   []*SuggestedType       `json:"suggestedTypes,omitempty"`
	RefersTo         map[string]interface{} `json:"refersTo,omitempty"`
	Deleted          bool                   `json:"deleted"`
	Grouping         *Grouping              `json:"grouping"`
	Min              float64                `json:"min"`
	Max              float64                `json:"max"`
}

// DataResource represents a set of variables found in a data asset.
type DataResource struct {
	ResID        string              `json:"resID"`
	ResType      string              `json:"resType"`
	ResPath      string              `json:"resPath"`
	IsCollection bool                `json:"isCollection"`
	Variables    []*Variable         `json:"columns,omitempty"`
	ResFormat    map[string][]string `json:"resFormat"`
}

// SuggestedType represents a classified variable type.
type SuggestedType struct {
	Type        string  `json:"type"`
	Probability float64 `json:"probability"`
	Provenance  string  `json:"provenance"`
}

// DatasetOrigin represents the originating information for a dataset
type DatasetOrigin struct {
	SearchResult  string `json:"searchResult"`
	Provenance    string `json:"provenance"`
	SourceDataset string `json:"sourceDataset"`
}

// Metadata represents a collection of dataset descriptions.
type Metadata struct {
	ID               string
	ParentDatasetIDs []string
	Name             string
	StorageName      string
	Description      string
	Summary          string
	SummaryMachine   string
	Raw              bool
	DataResources    []*DataResource
	Schema           *gabs.Container
	Classification   *gabs.Container
	NumRows          int64
	NumBytes         int64
	SchemaSource     string
	Redacted         bool
	DatasetFolder    string
	DatasetOrigins   []*DatasetOrigin
	SearchResult     string
	SearchProvenance string
	SourceDataset    string
	Type             string
}

// NewMetadata creates a new metadata instance.
func NewMetadata(id string, name string, description string, storageName string) *Metadata {
	return &Metadata{
		ID:            id,
		Name:          name,
		StorageName:   storageName,
		Description:   description,
		DataResources: make([]*DataResource, 0),
	}
}

// NewDataResource creates a new data resource instance.
func NewDataResource(id string, typ string, format map[string][]string) *DataResource {
	return &DataResource{
		ResID:     id,
		ResType:   typ,
		ResFormat: format,
		Variables: make([]*Variable, 0),
	}
}

// NormalizeDatasetID modifies a dataset ID to be compatible with postgres
// naming requirements.
func NormalizeDatasetID(id string) string {
	// datasets can't have '.' and should be lowercase.
	normalized := nameRegex.ReplaceAllString(id, "_")
	normalized = strings.ToLower(normalized)

	// add a prefix to handle cases where numbers are the first character.
	normalized = fmt.Sprintf("%s%s", DatasetPrefix, normalized)
	// truncate so that name is not longer than allowed table name limit - need to leave space
	// for name suffixes as well
	if len(normalized) > datasetIDSizeLimit {
		normalized = normalized[:datasetIDSizeLimit]
	}
	return normalized
}

// NormalizeVariableName normalizes a variable name.
func NormalizeVariableName(name string) string {
	nameNormalized := nameRegex.ReplaceAllString(name, "_")
	if len(nameNormalized) > variableNameSizeLimit {
		nameNormalized = nameNormalized[:variableNameSizeLimit]
	}
	return nameNormalized
}

func doesNameAlreadyExist(name string, existingVariables []*Variable) bool {
	for _, v := range existingVariables {
		if v != nil && v.Name == name {
			return true
		}
	}
	return false
}

func ensureUniqueNameRecursive(name string, existingVariables []*Variable, count int) string {
	newName := fmt.Sprintf("%s_%d", name, count)
	if doesNameAlreadyExist(newName, existingVariables) {
		return ensureUniqueNameRecursive(name, existingVariables, count+1)
	}
	return newName
}

func ensureUniqueName(name string, existingVariables []*Variable) string {
	if doesNameAlreadyExist(name, existingVariables) {
		return ensureUniqueNameRecursive(name, existingVariables, 0)
	}
	return name
}

// NewVariable creates a new variable.
func NewVariable(index int, name, displayName, originalName, typ, originalType, description string, role []string, distilRole string, refersTo map[string]interface{}, existingVariables []*Variable, normalizeName bool) *Variable {
	normalized := name
	if normalizeName {
		// normalize name
		normalized = NormalizeVariableName(name)

		// normalized name needs to be unique
		normalized = ensureUniqueName(normalized, existingVariables)
	}

	// select the first role by default.
	selectedRole := ""
	if len(role) > 0 {
		selectedRole = role[0]
	}
	if distilRole == "" {
		distilRole = VarDistilRoleData
	}
	if originalName == "" {
		originalName = normalized
	}
	if displayName == "" {
		displayName = name
	}
	if originalType == "" {
		originalType = typ
	}

	return &Variable{
		Name:             normalized,
		Index:            index,
		Type:             typ,
		Description:      description,
		OriginalType:     originalType,
		Role:             role,
		SelectedRole:     selectedRole,
		DistilRole:       distilRole,
		OriginalVariable: originalName,
		DisplayName:      displayName,
		RefersTo:         refersTo,
		SuggestedTypes:   make([]*SuggestedType, 0),
	}
}

// AddVariable creates and add a new variable to the data resource.
func (dr *DataResource) AddVariable(name string, originalName string, typ string, description string, role []string, distilRole string) {
	v := NewVariable(len(dr.Variables), name, "", originalName, typ, typ, description, role, distilRole, nil, dr.Variables, false)
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

	if len(m.DataResources) > 0 {
		return m.DataResources[0]
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
	case IntegerType, LongitudeType, LatitudeType, RealType, TimestampType:
		return dataTypeFloat
	case OrdinalType, CategoricalType, StringType:
		return dataTypeText
	case DateTimeType:
		return dataTypeDate
	case RealVectorType, RealListType:
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
	case LongitudeType, LatitudeType, RealType:
		return "'NaN'::double precision"
	case IntegerType, TimestampType:
		return int(0)
	case DateTimeType:
		return fmt.Sprintf("'%s'", time.Time{}.Format(dateFormat))
	case RealVectorType, RealListType:
		return "'{}'"
	default:
		return "''"
	}
}

// PostgresValueForFieldType generates the select field value for a given variable type.
func PostgresValueForFieldType(typ string, field string) string {
	fieldQuote := fmt.Sprintf("\"%s\"", field)
	switch typ {
	case RealListType:
		return fmt.Sprintf("string_to_array(%s, ',')", fieldQuote)
	case DateTimeType:
		// datetime may be only time so need to support both cases
		// times can have first value missing a 0 so want to first get a time value then add it to epoch time 0
		return fmt.Sprintf("CASE WHEN length(%[1]s) IN (4, 5) AND position(':' in %[1]s) > 0 THEN CONCAT('1970-01-01 ', to_char(to_timestamp(%[1]s, 'MI:SS'), 'HH24:MI:SS')) ELSE %[1]s END", fieldQuote)
	default:
		return fmt.Sprintf("\"%s\"", field)
	}
}

// IsTA2Field indicates whether or not a particular variable is recognized by a TA2.
func IsTA2Field(distilRole string) bool {
	return distilRole == VarDistilRoleData || distilRole == VarDistilRoleIndex
}

// IsIndexRole returns true if the d3m role is an index role.
func IsIndexRole(role string) bool {
	return role == RoleIndex || role == RoleMultiIndex
}
