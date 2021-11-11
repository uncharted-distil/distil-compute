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

package model

import (
	"fmt"
	"path"
	"regexp"
	"strings"

	"github.com/Jeffail/gabs/v2"
)

const (
	// ExplainValues is the variable that holds the col name for the explain values
	ExplainValues = "explain_values"
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
	// RoleAttribute is the role used for attribute fields.
	RoleAttribute = "attribute"

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
	// VarKeyField is the field name for the variable key.
	VarKeyField = "key"
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
	// VarDistilRoleSystemData indicates a var is important for the system but not intended for the user
	VarDistilRoleSystemData = "system-data"
	// VarDistilRoleFeaturized indicates a var has been featurized
	VarDistilRoleFeaturized = "featurized-data"
	// VarDistilRoleLabel indicates a var has been created through the label work flow
	VarDistilRoleLabel = "label"
	// VarDistilRoleIndex indicates a var has an index role in distil.
	VarDistilRoleIndex = "index"
	// VarDistilRoleData indicates a var has a data role in distil.
	VarDistilRoleData = "data"
	// VarDistilRoleGrouping indicates a var has a grouping role in distil.
	VarDistilRoleGrouping = "grouping"
	// VarDistilRoleMetadata is the distil role for metadata variables
	VarDistilRoleMetadata = "metadata"
	// VarDistilRoleAugmented is the distil role for variables created for UX, not for model creation. i.e. Outliers detection
	VarDistilRoleAugmented = "augmented"
	// VarDistilRolePadding is the distil role for variables that are empty or dummy variables (their only purpose is to add alignment for pipelines)
	VarDistilRolePadding = "padding"
	// VarDeleted flags whether the variable is deleted.
	VarDeleted = "deleted"
	// VarImmutableField is the name for the flag indicating whether the variable is immutable.
	VarImmutableField = "immutable"
	// VarGroupingField is the field name for the variable grouping.
	VarGroupingField = "grouping"
	// VarMinField is the field name for the min value.
	VarMinField = "min"
	// VarMaxField is the field name for the max value.
	VarMaxField = "max"
	// VarValuesField is the field name for the categorical values present in the data.
	VarValuesField = "values"

	// TypeTypeField is the type field of a suggested type
	TypeTypeField = "type"
	// TypeProbabilityField is the probability field of a suggested type
	TypeProbabilityField = "probability"
	// TypeProvenanceField is the provenance field of a suggested type
	TypeProvenanceField = "provenance"

	// DatasetPrefix is the prefix used for a normalized dataset id
	DatasetPrefix = "d_"
)

var (
	nameRegex     = regexp.MustCompile("[^a-zA-Z0-9]")
	truncateRegex = regexp.MustCompile(`(.*)(_\d+$)`)
)

// BaseGrouping provides access to the basic grouping information.
type BaseGrouping interface {
	GetDataset() string
	GetType() string
	GetIDCol() string
	GetSubIDs() []string
	GetHidden() []string
	IsNil() bool
}

// ClusteredGrouping provides access to grouping cluster information.
type ClusteredGrouping interface {
	GetClusterCol() string
}

// Grouping represents a variable grouping.
type Grouping struct {
	Dataset string   `json:"dataset"`
	Type    string   `json:"type"`
	IDCol   string   `json:"idCol"`
	SubIDs  []string `json:"subIds"`
	Hidden  []string `json:"hidden"`
}

// GeoCoordinateGrouping is used for geocoordinate grouping information.
type GeoCoordinateGrouping struct {
	Grouping
	XCol string `json:"xCol"`
	YCol string `json:"yCol"`
}

// TimeseriesGrouping is used for timeseries grouping information.
type TimeseriesGrouping struct {
	Grouping
	ClusterCol string `json:"clusterCol"`
	XCol       string `json:"xCol"`
	YCol       string `json:"yCol"`
}

// MultiBandImageGrouping is used for remote sensing grouping information.
type MultiBandImageGrouping struct {
	Grouping
	BandCol    string `json:"bandCol"`
	ImageCol   string `json:"imageCol"`
	ClusterCol string `json:"clusterCol"`
}

// GeoBoundsGrouping is used for geo bounds groups.
type GeoBoundsGrouping struct {
	Grouping
	CoordinatesCol string `json:"coordinatesCol"`
	PolygonCol     string `json:"polygonCol"`
}

// GetDataset returns the grouping dataset.
func (g *Grouping) GetDataset() string {
	return g.Dataset
}

// GetType returns the grouping type.
func (g *Grouping) GetType() string {
	return g.Type
}

// GetIDCol returns the grouping id column name.
func (g *Grouping) GetIDCol() string {
	return g.IDCol
}

// GetSubIDs returns the grouping sub id column names.
func (g *Grouping) GetSubIDs() []string {
	return g.SubIDs
}

// GetHidden returns the grouping hidden column names.
func (g *Grouping) GetHidden() []string {
	return g.Hidden
}

// IsNil checks if this is a typed nil.
func (g *Grouping) IsNil() bool {
	return g == nil
}

// GetClusterCol returns the cluster column name for a remote sensing group.
func (t *TimeseriesGrouping) GetClusterCol() string {
	return t.ClusterCol
}

// GetClusterCol returns the cluster column name for a remote sensing group.
func (t *MultiBandImageGrouping) GetClusterCol() string {
	return t.ClusterCol
}

// GetPolygonCol returns the polygon representation of the geo bounds group.
func (t *GeoBoundsGrouping) GetPolygonCol() string {
	return t.PolygonCol
}

// Variable represents a single variable description.
type Variable struct {
	Key              string                 `json:"key"`
	HeaderName       string                 `json:"colName"`
	Type             string                 `json:"colType,omitempty"`
	Description      string                 `json:"colDescription,omitempty"`
	OriginalType     string                 `json:"colOriginalType,omitempty"`
	SelectedRole     string                 `json:"selectedRole,omitempty"`
	Role             []string               `json:"role,omitempty"`
	DistilRole       []string               `json:"distilRole,omitempty"`
	OriginalVariable string                 `json:"colOriginalName"`
	DisplayName      string                 `json:"colDisplayName,omitempty"`
	Importance       float64                `json:"importance"`
	Index            int                    `json:"colIndex"`
	SuggestedTypes   []*SuggestedType       `json:"suggestedTypes,omitempty"`
	RefersTo         map[string]interface{} `json:"refersTo,omitempty"`
	Deleted          bool                   `json:"deleted"`
	Immutable        bool                   `json:"immutable"`
	Grouping         BaseGrouping           `json:"grouping"`
	Min              float64                `json:"min"`
	Max              float64                `json:"max"`
	Values           []string               `json:"values"`
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
	Classification   *ClassificationData
	NumRows          int64
	NumBytes         int64
	SchemaSource     string
	Redacted         bool
	DatasetFolder    string
	DatasetOrigins   []*DatasetOrigin
	SearchResult     string
	SearchProvenance string
	SourceDataset    string
	LearningDataset  string
	Type             string
	Digest           string
	Immutable        bool
	Clone            bool
}

// ClassificationData contains semantic type information by column index.
type ClassificationData struct {
	Labels        [][]string  `json:"labels"`
	Probabilities [][]float64 `json:"label_probabilities"`
	Path          string      `json:"path"`
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
		// Standard approach to deconflicting names in the system is to append `_N`.  We need to make
		// sure we don't truncate that portion.
		matches := truncateRegex.FindStringSubmatch(normalized)
		if len(matches) != 3 {
			return normalized[:datasetIDSizeLimit]
		}
		bodyLength := datasetIDSizeLimit - len(matches[2])
		normalized = fmt.Sprintf("%s%s", matches[1][:bodyLength], matches[2])
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

func doesKeyAlreadyExist(key string, existingVariables []*Variable) bool {
	for _, v := range existingVariables {
		if v != nil && v.Key == key {
			return true
		}
	}
	return false
}

func ensureUniqueKeyRecursive(key string, existingVariables []*Variable, count int) string {
	newKey := fmt.Sprintf("%s_%d", key, count)
	if doesKeyAlreadyExist(newKey, existingVariables) {
		return ensureUniqueKeyRecursive(key, existingVariables, count+1)
	}
	return newKey
}

func ensureUniqueKey(key string, existingVariables []*Variable) string {
	if doesKeyAlreadyExist(key, existingVariables) {
		return ensureUniqueKeyRecursive(key, existingVariables, 0)
	}
	return key
}

// NewVariable creates a new variable.
func NewVariable(index int, key, displayName, headerName, originalName, typ, originalType, description string, role []string, distilRole []string, refersTo map[string]interface{}, existingVariables []*Variable, normalizeName bool) *Variable {
	normalized := key
	if normalizeName {
		// normalize name
		normalized = NormalizeVariableName(key)

		// normalized name needs to be unique
		normalized = ensureUniqueKey(normalized, existingVariables)
	}

	// select the first role by default.
	selectedRole := ""
	if len(role) > 0 {
		selectedRole = role[0]
	}
	if len(distilRole) == 0 {
		distilRole = []string{VarDistilRoleData}
	}
	if originalName == "" {
		originalName = normalized
	}
	if displayName == "" {
		displayName = key
	}
	if headerName == "" {
		// CSV is flexible with header names, but other formats like parquet are not
		headerName = normalized
	}
	if originalType == "" {
		originalType = typ
	}

	return &Variable{
		Key:              normalized,
		Index:            index,
		Type:             typ,
		Description:      description,
		OriginalType:     originalType,
		Role:             role,
		SelectedRole:     selectedRole,
		DistilRole:       distilRole,
		OriginalVariable: originalName,
		DisplayName:      displayName,
		HeaderName:       headerName,
		RefersTo:         refersTo,
		SuggestedTypes:   make([]*SuggestedType, 0),
	}
}

// AddVariable creates and add a new variable to the data resource.
func (dr *DataResource) AddVariable(name string, originalName string, typ string, description string, role []string, distilRole []string) {
	v := NewVariable(len(dr.Variables), name, "", name, originalName, typ, typ, description, role, distilRole, nil, dr.Variables, false)
	dr.Variables = append(dr.Variables, v)
}

// GetMainDataResource returns the data resource that contains the D3M index.
func (m *Metadata) GetMainDataResource() *DataResource {
	// main data resource has d3m index variable
	for _, dr := range m.DataResources {
		for _, v := range dr.Variables {
			if v.Key == D3MIndexFieldName {
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
		header[hIndex] = field.HeaderName
	}

	return header
}

// GetResourcePath returns the absolute path of the data resource.
func GetResourcePath(schemaFile string, dataResource *DataResource) string {
	return GetResourcePathFromFolder(path.Dir(schemaFile), dataResource)
}

// GetResourcePathFromFolder returns the absolute path of the data resource.
func GetResourcePathFromFolder(datasetFolder string, dataResource *DataResource) string {
	// path can either be absolute or relative to the schema file
	drPath := dataResource.ResPath
	if len(drPath) > 0 && drPath[0] != '/' {
		drPath = path.Join(datasetFolder, drPath)
	}

	return drPath
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

// IsGrouping returns true if the variable is a grouping.
func (v *Variable) IsGrouping() bool {
	return v.Grouping != nil && !v.Grouping.IsNil()
}

// Clone clones a variable into a new variable.
func (v *Variable) Clone() *Variable {
	// HACK: RefersTo and Grouping are not doing deep copies
	clone := &Variable{}
	clone.Role = append([]string{}, v.Role...)
	clone.Values = append(clone.Values, v.Values...)
	for _, s := range v.SuggestedTypes {
		cs := *s
		clone.SuggestedTypes = append(clone.SuggestedTypes, &cs)
	}
	if v.RefersTo != nil {
		clone.RefersTo = map[string]interface{}{}
		for k, t := range v.RefersTo {
			clone.RefersTo[k] = t
		}
	}
	if v.Grouping != nil {
		clone.Grouping = v.Grouping
	}
	clone.Key = v.Key
	clone.HeaderName = v.HeaderName
	clone.Type = v.Type
	clone.Description = v.Description
	clone.OriginalType = v.OriginalType
	clone.SelectedRole = v.SelectedRole
	clone.DistilRole = v.DistilRole
	clone.OriginalVariable = v.OriginalVariable
	clone.DisplayName = v.DisplayName
	clone.Importance = v.Importance
	clone.Index = v.Index
	clone.Deleted = v.Deleted
	clone.Immutable = v.Immutable
	clone.Min = v.Min
	clone.Max = v.Max

	return clone
}

// IsTA2Field indicates whether or not a particular variable is recognized by a TA2.
func (v *Variable) IsTA2Field() bool {
	if v.HasAnyRole([]string{VarDistilRoleData, VarDistilRoleIndex, VarDistilRoleSystemData}) {
		return true
	}

	if v.HasRole(VarDistilRoleGrouping) && IsAttributeRole(v.SelectedRole) {
		return true
	}

	return false
}

// IsIndexRole returns true if the d3m role is an index role.
func IsIndexRole(role string) bool {
	return role == RoleIndex || role == RoleMultiIndex
}

// IsAttributeRole returns true if the d3m role is an attribute role.
func IsAttributeRole(role string) bool {
	return role == RoleAttribute
}

// HasRole checks to see if the supplied role exists within the variable's DistilRole property
func (v *Variable) HasRole(role string) bool {

	for _, element := range v.DistilRole {
		if element == role {
			return true
		}
	}

	return false
}

// HasAnyRole checks that at least 1 of the supplied roles exist within the variable's DistilRole property
// equivalent to OR operation
func (v *Variable) HasAnyRole(roles []string) bool {
	roleMap := map[string]bool{}

	for _, role := range roles {
		roleMap[role] = true
	}

	for _, role := range v.DistilRole {
		if _, ok := roleMap[role]; ok {
			return true
		}
	}

	return false
}

// HasAllRole checks that all supplied roles exist within the variable's DistilRole property
// equivalent to AND operation
func (v *Variable) HasAllRole(roles []string) bool {
	roleMap := map[string]bool{}

	for _, role := range roles {
		roleMap[role] = true
	}

	for _, role := range v.DistilRole {
		if _, ok := roleMap[role]; ok {
			return false
		}
	}

	return true
}
