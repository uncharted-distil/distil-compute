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

import log "github.com/unchartedsoftware/plog"

const (
	// Distil Internal Type Keys.  These are the application's set of recognized types within the
	// server and in the client code, and any external types, such as that used by Lincoln Labs in
	// the D3M datasets, those used by TA2 systems and in the D3M runtime, or those produced by analytics like
	// Simon, should be mapped and translated to/from this set.
	//
	// NOTE: these are copied to `distil/public/util/types.ts` and
	// should be kept up to date in case of changes.

	// AddressType is the schema type for address values
	AddressType = "address"
	// IndexType is the schema type for index values
	IndexType = "index"
	// IntegerType is the schema type for int values
	IntegerType = "integer"
	// RealType is the schema type for real values, and is equivalent to FloatType
	RealType = "real"
	// RealVectorType is the schema type for a vector of real values
	RealVectorType = "realVector"
	// RealListType is the schema type for a list of real values
	RealListType = "realList"
	// BoolType is the schema type for bool values
	BoolType = "boolean"
	// DateTimeType is the schema type for date/time values
	DateTimeType = "dateTime"
	// OrdinalType is the schema type for ordinal values
	OrdinalType = "ordinal"
	// CategoricalType is the schema type for categorical values
	CategoricalType = "categorical"
	// NumericalType is the schema type for numerical values
	NumericalType = "numerical"
	// StringType is the schema type for string/text values
	StringType = "string"
	// CityType is the schema type for city values
	CityType = "city"
	// CountryType is the schema type for country values
	CountryType = "country"
	// EmailType is the schema type for email values
	EmailType = "email"
	// LatitudeType is the schema type for latitude values
	LatitudeType = "latitude"
	// LongitudeType is the schema type for longitude values
	LongitudeType = "longitude"
	// PhoneType is the schema type for phone values
	PhoneType = "phone"
	// PostalCodeType is the schema type for postal code values
	PostalCodeType = "postal_code"
	// StateType is the schema type for state values
	StateType = "state"
	// URIType is the schema type for URI values
	URIType = "uri"
	// ImageType is the schema type for Image values
	ImageType = "image"
	// MultiBandImageType is t he schema type for multi-band (satellite) images
	MultiBandImageType = "multiband_image"
	// TimeSeriesType is the schema type for timeseries values
	TimeSeriesType = "timeseries"
	// GeoCoordinateType is the schema type for geocoordinate values
	GeoCoordinateType = "geocoordinate"
	// RemoteSensingType is the schema type for remote sensing values
	RemoteSensingType = "remote_sensing"
	// TimestampType is the schema type for timestamp values
	TimestampType = "timestamp"
	// UnknownType is the schema type for unknown values
	UnknownType = "unknown"

	// Simon types - these are types produced by the Simon analytic which is run during ingest.  These types
	// should not be used/stored internally, but instead translated into the Distil types at the application
	// boundaries.

	// SimonAddressType is the Simon type representing a street address
	SimonAddressType = "address"
	// SimonBooleanType is the Simon type representing a boolean
	SimonBooleanType = "boolean"
	// SimonDateTimeType is the Simon representation of a date time value
	SimonDateTimeType = "datetime"
	// SimonEmailType is the Simon representation of an email address
	SimonEmailType = "email"
	// SimonFloatType is the Simon type representing a float
	SimonFloatType = "float"
	// SimonIntegerType is the Simon type representing an integer
	SimonIntegerType = "int"
	// SimonPhoneType is the Simon representation of a phone number
	SimonPhoneType = "phone"
	// SimonStringType is the Simon type representing a string
	SimonStringType = "text"
	// SimonURIType is the Simon type representing a URI
	SimonURIType = "uri"
	// SimonCategoricalType is the Simon type representing categorical values
	SimonCategoricalType = "categorical"
	// SimonOrdinalType is the Simon type representing ordinalvalues
	SimonOrdinalType = "ordinal"
	// SimonStateType is the Simon type representing state values
	SimonStateType = "state"
	// SimonCityType is the Simon type representing city values
	SimonCityType = "city"
	// SimonPostalCodeType is the Simon type representing postal code values
	SimonPostalCodeType = "postal_code"
	// SimonLatitudeType is the Simon type representing geographic latitude values
	SimonLatitudeType = "latitude"
	// SimonLongitudeType is the Simon type representing geographic longitude values
	SimonLongitudeType = "longitude"
	// SimonCountryType is the Simon type representing country values
	SimonCountryType = "country"
	// SimonCountryCodeType is the Simon type representing country code values
	SimonCountryCodeType = "country_code"

	// TA2 Semantic Type Keys - defined in
	// https://gitlab.com/datadrivendiscovery/d3m/blob/devel/d3m/metadata/schemas/v0/definitions.json
	// These are the agreed upond set of types that are consumable by a downstream TA2 system and the
	// D3M runtime.  These should not be used/stored internally, but instead translated into the Distil
	// types at the application boundaries.

	// TA2AttributeType is the semantic type representing an attribute
	TA2AttributeType = "https://metadata.datadrivendiscovery.org/types/Attribute"
	// TA2UnknownType is the semantic type representing the unknown field type
	TA2UnknownType = "https://metadata.datadrivendiscovery.org/types/UnknownType"
	// TA2StringType is the semantic type reprsenting a text/string
	TA2StringType = "http://schema.org/Text"
	// TA2IntegerType is the TA2 semantic type for an integer value
	TA2IntegerType = "http://schema.org/Integer"
	// TA2RealType is the TA2 semantic type for a real value
	TA2RealType = "http://schema.org/Float"
	// TA2BooleanType is the TA2 semantic type for a boolean value
	TA2BooleanType = "http://schema.org/Boolean"
	// TA2LocationType is the TA2 semantic type for a location value
	TA2LocationType = "https://metadata.datadrivendiscovery.org/types/Location"
	// TA2DateTimeType is the TA2 semantic type for a datetime value
	TA2DateTimeType = "http://schema.org/DateTime"
	// TA2TimeType is the TA2 semantic type for a time (role) value
	TA2TimeType = "https://metadata.datadrivendiscovery.org/types/Time"
	// TA2CategoricalType is the TA2 semantic type for categorical data
	TA2CategoricalType = "https://metadata.datadrivendiscovery.org/types/CategoricalData"
	// TA2OrdinalType is the TA2 semantic type for ordinal (ordered categorical) data
	TA2OrdinalType = "https://metadata.datadrivendiscovery.org/types/OrdinalData"
	// TA2ImageType is the TA2 semantic type for image data
	TA2ImageType = "http://schema.org/ImageObject"
	// TA2TimeSeriesType is the TA2 semantic type for timeseries data
	TA2TimeSeriesType = "https://metadata.datadrivendiscovery.org/types/Timeseries"
	// TA2RealVectorType is the TA2 semantic type for vector data
	TA2RealVectorType = "https://metadata.datadrivendiscovery.org/types/Vector"

	// TA2 Role keys

	// TA2TargetType is the semantic type indicating a prediction target
	TA2TargetType = "https://metadata.datadrivendiscovery.org/types/Target"
	// TA2GroupingKeyType is the semantic type indicating a grouping key
	TA2GroupingKeyType = "https://metadata.datadrivendiscovery.org/types/GroupingKey"

	// Lincoln Labs D3M Dataset Type Keys - these are the types used in the D3M dataset format.
	// These should not be used/stored internally, but instead translated into the Distil
	// types at the application boundaries.

	// BooleanSchemaType is the schema doc type for boolean data
	BooleanSchemaType = "boolean"
	// IntegerSchemaType is the schema doc type for integer data
	IntegerSchemaType = "integer"
	// RealSchemaType is the schema doc type for real data
	RealSchemaType = "real"
	// StringSchemaType is the schema doc type for string/text data
	StringSchemaType = "string"
	// TextSchemaTypeDeprecated is the old schema doc type for string/text data
	TextSchemaTypeDeprecated = "text"
	// CategoricalSchemaType is the schema doc type for categorical data
	CategoricalSchemaType = "categorical"
	// DatetimeSchemaType is the schema doc type for datetime data
	DatetimeSchemaType = "dateTime"
	// RealVectorSchemaType is the schema doc type for a vector of real data
	RealVectorSchemaType = "realVector"
	// JSONSchemaType is the schema doc type for json data
	JSONSchemaType = "json"
	// GeoJSONSchemaType is the schema doc type for geo json data
	GeoJSONSchemaType = "geojson"
	// ImageSchemaType is the schema doc type for image data
	ImageSchemaType = "image"
	// TimeSeriesSchemaType is the schema doc type for image data
	TimeSeriesSchemaType = "timeseries"
	// TimestampSchemaType is the schema doc type for image data
	TimestampSchemaType = "timestamp"
	// UnknownSchemaType is the scehma type representing the unknown field type
	UnknownSchemaType = "unknown"
)

var (
	categoricalTypes = map[string]bool{
		CategoricalType: true,
		OrdinalType:     true,
		BoolType:        true,
		AddressType:     true,
		CityType:        true,
		CountryType:     true,
		EmailType:       true,
		PhoneType:       true,
		PostalCodeType:  true,
		StateType:       true,
		URIType:         true,
		UnknownType:     true,
	}
	numericalTypes = map[string]bool{
		LongitudeType: true,
		LatitudeType:  true,
		RealType:      true,
		IntegerType:   true,
		IndexType:     true,
	}
	floatingPointTypes = map[string]bool{
		LongitudeType: true,
		LatitudeType:  true,
		RealType:      true,
	}

	// Maps from Distil internal type to TA2 supported type
	ta2TypeMap = map[string]string{
		AddressType:        TA2StringType,
		IndexType:          TA2IntegerType,
		IntegerType:        TA2IntegerType,
		RealType:           TA2RealType,
		BoolType:           TA2BooleanType,
		DateTimeType:       TA2TimeType,
		OrdinalType:        TA2CategoricalType,
		CategoricalType:    TA2CategoricalType,
		NumericalType:      TA2RealType,
		StringType:         TA2StringType,
		CityType:           TA2CategoricalType,
		CountryType:        TA2CategoricalType,
		EmailType:          TA2StringType,
		LatitudeType:       TA2RealType,
		LongitudeType:      TA2RealType,
		PhoneType:          TA2StringType,
		PostalCodeType:     TA2StringType,
		StateType:          TA2CategoricalType,
		URIType:            TA2StringType,
		ImageType:          TA2StringType,
		MultiBandImageType: TA2StringType,
		TimestampType:      TA2TimeType,
		TimeSeriesType:     TA2TimeSeriesType,
		UnknownType:        TA2UnknownType,
		RealVectorType:     TA2RealVectorType,
		RealListType:       TA2RealVectorType,
	}

	// Maps from Distil internal type to D3M dataset doc type
	schemaTypeMap = map[string]string{
		AddressType:        StringSchemaType,
		IndexType:          IntegerSchemaType,
		IntegerType:        IntegerSchemaType,
		RealType:           RealSchemaType,
		BoolType:           BooleanSchemaType,
		DateTimeType:       DatetimeSchemaType,
		OrdinalType:        CategoricalSchemaType,
		CategoricalType:    CategoricalSchemaType,
		NumericalType:      RealSchemaType,
		StringType:         StringSchemaType,
		CityType:           StringSchemaType,
		CountryType:        StringSchemaType,
		EmailType:          StringSchemaType,
		LatitudeType:       RealSchemaType,
		LongitudeType:      RealSchemaType,
		PhoneType:          StringSchemaType,
		PostalCodeType:     StringSchemaType,
		StateType:          StringSchemaType,
		URIType:            StringSchemaType,
		ImageType:          StringSchemaType,
		MultiBandImageType: StringSchemaType,
		TimeSeriesType:     TimeSeriesSchemaType,
		TimestampType:      TimestampSchemaType,
		UnknownType:        UnknownSchemaType,
		RealVectorType:     RealVectorSchemaType,
	}

	// Maps from Lincoln Labs D3M dataset doc type to Distil internal type
	llTypeMap = map[string]string{
		BooleanSchemaType:        BoolType,
		IntegerSchemaType:        IntegerType,
		RealSchemaType:           RealType,
		CategoricalSchemaType:    CategoricalType,
		DatetimeSchemaType:       DateTimeType,
		RealVectorSchemaType:     RealVectorType,
		JSONSchemaType:           StringType,
		GeoJSONSchemaType:        StringType,
		ImageSchemaType:          ImageType,
		TimeSeriesSchemaType:     TimeSeriesType,
		TimestampSchemaType:      TimestampType,
		TextSchemaTypeDeprecated: StringType,
		StringSchemaType:         StringType,
		UnknownSchemaType:        UnknownType,
	}

	// Maps from Simon type to Distil internal type
	simonTypeMap = map[string]string{
		SimonAddressType:     AddressType,
		SimonBooleanType:     BoolType,
		SimonDateTimeType:    DateTimeType,
		SimonEmailType:       EmailType,
		SimonFloatType:       RealType,
		SimonIntegerType:     IntegerType,
		SimonPhoneType:       PhoneType,
		SimonStringType:      StringType,
		SimonURIType:         URIType,
		SimonCategoricalType: CategoricalType,
		SimonOrdinalType:     OrdinalType,
		SimonStateType:       StateType,
		SimonCityType:        CityType,
		SimonPostalCodeType:  PostalCodeType,
		SimonLatitudeType:    LatitudeType,
		SimonLongitudeType:   LongitudeType,
		SimonCountryType:     CountryType,
		SimonCountryCodeType: CountryType,
	}

	simonBasicTypes = map[string]bool{
		IndexType:      true,
		IntegerType:    true,
		RealType:       true,
		RealVectorType: true,
		BoolType:       true,
		DateTimeType:   true,
		NumericalType:  true,
		StringType:     true,
	}
)

// IsNumerical indicates whether or not a schema type is numeric for the purposes
// of analysis.
func IsNumerical(typ string) bool {
	return numericalTypes[typ]
}

// IsFloatingPoint indicates whether or not a schema type is a floating point
// value.
func IsFloatingPoint(typ string) bool {
	return floatingPointTypes[typ]
}

// IsDatabaseFloatingPoint indicates whether or not a database type is a floating point
// value.
func IsDatabaseFloatingPoint(typ string) bool {
	return typ == dataTypeFloat
}

// IsCategorical indicates whether or not a schema type is categorical for the purposes
// of analysis.
func IsCategorical(typ string) bool {
	return categoricalTypes[typ]
}

// IsText indicates whether or not a schema type is text for the purposes
// of analysis.
func IsText(typ string) bool {
	return typ == StringType
}

// IsVector indicates whether or not a schema type is a vector for the purposes
// of analysis.
func IsVector(typ string) bool {
	return typ == RealVectorType
}

// IsList indicates whether or not a schema type is a list for the purposes
// of analysis.
func IsList(typ string) bool {
	return typ == RealListType
}

// IsImage indicates whether or not a schema type is an image for the purposes
// of analysis.
func IsImage(typ string) bool {
	return typ == ImageType
}

// IsMultiBandImage indicates whether or not a schema type is a multi-band (satellite) image for the
// purposes of analysis
func IsMultiBandImage(typ string) bool {
	return typ == MultiBandImageType
}

// IsTimeSeries indicates whether or not a schema type is a timeseries for the purposes
// of analysis.
func IsTimeSeries(typ string) bool {
	return typ == TimeSeriesType
}

// IsGeoCoordinate indicates whether or not a schema type is a geo coordinate
// for the purposes of analysis.
func IsGeoCoordinate(typ string) bool {
	return typ == GeoCoordinateType
}

// IsRemoteSensing indicates whether or not a schema type is a remote sensing
// for the purposes of analysis.
func IsRemoteSensing(typ string) bool {
	return typ == RemoteSensingType
}

// IsTimestamp indicates whether or not a schema type is a timestamp for the purposes
// of analysis.
func IsTimestamp(typ string) bool {
	return typ == TimestampType
}

// IsDateTime indicates whether or not a schema type is a date time for the purposes
// of analysis.
func IsDateTime(typ string) bool {
	return typ == DateTimeType
}

// HasFeatureVar indicates whether or not a schema type has a corresponding feature var.
func HasFeatureVar(typ string) bool {
	return IsImage(typ)
}

// HasClusterVar indicates whether or not a schema type has a corresponding cluster var.
func HasClusterVar(typ string) bool {
	return IsTimeSeries(typ)
}

// MapTA2Type maps an internal Distil type to a TA2 type.  Distil
// type should be recognized at this point and it is programmer error if not.
func MapTA2Type(typ string) string {
	return ta2TypeMap[typ]
}

// MapSchemaType maps an internal Distil type to an LL D3M dataset doc type.  Distil
// type should be recognized at this point and it is programmer error if not.
func MapSchemaType(typ string) string {
	return schemaTypeMap[typ]
}

// MapLLType maps a LL D3M dataset doc type to an internal type
func MapLLType(typ string) string {
	mapped := llTypeMap[typ]
	if mapped == "" {
		log.Warnf("D3M dataset doc type '%s' has no mappping to Distl type defined and will be passed through", typ)
		mapped = typ
	}

	return mapped
}

// MapSimonType maps a Simon type to an internal data type.
func MapSimonType(typ string) string {
	mapped := simonTypeMap[typ]
	if mapped == "" {
		log.Warnf("Simon type '%s' has no mappping to Distil type defined and will be passed through", typ)
		mapped = typ
	}

	return mapped
}

// IsBasicSimonType check whether or not a type is a basic simon type.
func IsBasicSimonType(typ string) bool {
	return simonBasicTypes[typ]
}
