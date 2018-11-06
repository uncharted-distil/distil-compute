package model

import (
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
)

// Variable represents a single variable description.
type Variable struct {
	Name             string                 `json:"colName"`
	Type             string                 `json:"colType,omitempty"`
	OriginalType     string                 `json:"colOriginalType,omitempty"`
	SelectedRole     string                 `json:"selectedRole,omitempty"`
	Role             []string               `json:"role,omitempty"`
	DistilRole       string                 `json:"distilRole,omitempty"`
	OriginalVariable string                 `json:"varOriginalName"`
	OriginalName     string                 `json:"colOriginalName,omitempty"`
	DisplayName      string                 `json:"colDisplayName,omitempty"`
	Importance       int                    `json:"importance,omitempty"`
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
	schema         *gabs.Container
	classification *gabs.Container
	NumRows        int64
	NumBytes       int64
	SchemaSource   string
	Redacted       bool
	DatasetFolder  string
}
