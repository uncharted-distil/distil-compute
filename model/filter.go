//
//   Copyright © 2020 Uncharted Software Inc.
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
	"sort"
)

const (
	// DefaultFilterSize represents the default filter search size.
	DefaultFilterSize = 100
	// FilterSizeLimit represents the largest filter size.
	FilterSizeLimit = 1000
	// GeoBoundsFilter represents a geobound filter type.
	GeoBoundsFilter = "geobounds"
	// CategoricalFilter represents a categorical filter type.
	CategoricalFilter = "categorical"
	// ClusterFilter represents a cluster filter type.
	ClusterFilter = "cluster"
	// NumericalFilter represents a numerical filter type.
	NumericalFilter = "numerical"
	// BivariateFilter represents a numerical filter type.
	BivariateFilter = "bivariate"
	// DatetimeFilter represents a datetime filter type.
	DatetimeFilter = "datetime"
	// TextFilter represents a text filter type.
	TextFilter = "text"
	// VectorFilter represents a text filter type.
	VectorFilter = "vector"
	// RowFilter represents a numerical filter type.
	RowFilter = "row"
	// IncludeFilter represents an inclusive filter mode.
	IncludeFilter = "include"
	// ExcludeFilter represents an exclusive filter mode.
	ExcludeFilter = "exclude"
)

// StringSliceEqual compares 2 string slices to see if they are equal.
func StringSliceEqual(a, b []string) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// Bounds defines a bounding box
type Bounds struct {
	MinX float64 `json:"minX"`
	MaxX float64 `json:"maxX"`
	MinY float64 `json:"minY"`
	MaxY float64 `json:"maxY"`
}

// Filter defines a variable filter.
type Filter struct {
	Key        string   `json:"key"`
	Type       string   `json:"type"`
	NestedType string   `json:"nestedType"`
	Mode       string   `json:"mode"`
	Min        *float64 `json:"min"`
	Max        *float64 `json:"max"`
	Bounds     *Bounds  `json:"bounds"`
	Categories []string `json:"categories"`
	D3mIndices []string `json:"d3mIndices"`
}

// NewNumericalFilter instantiates a numerical filter.
func NewNumericalFilter(key string, mode string, min float64, max float64) *Filter {
	return &Filter{
		Key:  key,
		Type: NumericalFilter,
		Mode: mode,
		Min:  &min,
		Max:  &max,
	}
}

// NewVectorFilter instantiates a vector filter.
func NewVectorFilter(key string, nestedType string, mode string, min float64, max float64) *Filter {
	return &Filter{
		Key:        key,
		Type:       VectorFilter,
		NestedType: nestedType,
		Mode:       mode,
		Min:        &min,
		Max:        &max,
	}
}

// NewDatetimeFilter instantiates a datetime filter.
func NewDatetimeFilter(key string, mode string, min float64, max float64) *Filter {
	return &Filter{
		Key:  key,
		Type: DatetimeFilter,
		Mode: mode,
		Min:  &min,
		Max:  &max,
	}
}

// NewBivariateFilter instantiates a numerical filter.
func NewBivariateFilter(key string, mode string, minX float64, maxX float64, minY float64, maxY float64) *Filter {
	return &Filter{
		Key:  key,
		Type: BivariateFilter,
		Mode: mode,
		Bounds: &Bounds{
			MinX: minX,
			MaxX: maxX,
			MinY: minY,
			MaxY: maxY,
		},
	}
}

// NewGeoBoundsFilter instantiates a geobounds filter.
func NewGeoBoundsFilter(key string, mode string, minX float64, maxX float64, minY float64, maxY float64) *Filter {
	return &Filter{
		Key:  key,
		Type: GeoBoundsFilter,
		Mode: mode,
		Bounds: &Bounds{
			MinX: minX,
			MaxX: maxX,
			MinY: minY,
			MaxY: maxY,
		},
	}
}

// NewCategoricalFilter instantiates a categorical filter.
func NewCategoricalFilter(key string, mode string, categories []string) *Filter {
	sort.Strings(categories)
	return &Filter{
		Key:        key,
		Type:       CategoricalFilter,
		Mode:       mode,
		Categories: categories,
	}
}

// NewClusterFilter instantiates a cluster filter.
func NewClusterFilter(key string, mode string, categories []string) *Filter {
	sort.Strings(categories)
	return &Filter{
		Key:        key,
		Type:       ClusterFilter,
		Mode:       mode,
		Categories: categories,
	}
}

// NewTextFilter instantiates a text filter.
func NewTextFilter(key string, mode string, categories []string) *Filter {
	sort.Strings(categories)
	return &Filter{
		Key:        key,
		Type:       TextFilter,
		Mode:       mode,
		Categories: categories,
	}
}

// NewRowFilter instantiates a row filter.
func NewRowFilter(mode string, d3mIndices []string) *Filter {
	return &Filter{
		Key:        "d3mIndex",
		Type:       RowFilter,
		Mode:       mode,
		D3mIndices: d3mIndices,
	}
}
