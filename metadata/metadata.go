//
//   Copyright © 2019 Uncharted Software Inc.
//
//   Licensed under the Apache License, Version 2.0 (the "License");
//   you may not use this file except in compliance with the License.
//   You may obtain a copy of the License at
//
//	   http://www.apache.org/licenses/LICENSE-2.0
//
//   Unless required by applicable law or agreed to in writing, software
//   distributed under the License is distributed on an "AS IS" BASIS,
//   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//   See the License for the specific language governing permissions and
//   limitations under the License.

package metadata

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/Jeffail/gabs/v2"
	"github.com/pkg/errors"
	log "github.com/unchartedsoftware/plog"

	"github.com/uncharted-distil/distil-compute/model"
)

// DatasetSource flags the type of ingest action that created a dataset
type DatasetSource string

const (
	// ProvenanceSimon identifies the type provenance as Simon
	ProvenanceSimon = "d3m.primitives.distil.simon"
	// ProvenanceSchema identifies the type provenance as schema
	ProvenanceSchema = "schema"

	schemaVersion = "4.0.0"
	license       = "Unknown"
	summaryLength = 256

	// Seed flags a dataset as ingested from seed data
	Seed DatasetSource = "seed"

	// Contrib flags a dataset as being ingested from contributed data
	Contrib DatasetSource = "contrib"

	// Augmented flags a dataset as being ingested from augmented data
	Augmented DatasetSource = "augmented"
)

var (
	typeProbabilityThreshold = 0.8
)

// SummaryResult captures the output of a summarization primitive.
type SummaryResult struct {
	Summary string `json:"summary"`
}

// SetTypeProbabilityThreshold below which a suggested type is not used as
// variable type
func SetTypeProbabilityThreshold(threshold float64) {
	typeProbabilityThreshold = threshold
}

// IsMetadataVariable indicates whether or not a variable is additional metadata
// added to the source.
func IsMetadataVariable(v *model.Variable) bool {
	return strings.HasPrefix(v.Name, "_")
}

// LoadMetadataFromOriginalSchema loads metadata from a schema file.
func LoadMetadataFromOriginalSchema(schemaPath string, augmentFromData bool) (*model.Metadata, error) {
	meta := &model.Metadata{
		SchemaSource: model.SchemaSourceOriginal,
	}
	err := loadSchema(meta, schemaPath)
	if err != nil {
		return nil, err
	}
	err = loadName(meta)
	if err != nil {
		return nil, err
	}
	err = loadID(meta)
	if err != nil {
		return nil, err
	}
	err = loadAbout(meta)
	if err != nil {
		return nil, err
	}
	err = loadOriginalSchemaVariables(meta, schemaPath)
	if err != nil {
		return nil, err
	}

	// read the header of every data resource and augment the variables since
	// the metadata may not specify every variable
	if augmentFromData {
		for _, dr := range meta.DataResources {
			dataPath := path.Join(path.Dir(schemaPath), dr.ResPath)

			// collection data resources need special care since datapath is a folder
			if dr.IsCollection {
				var ok bool
				dataPath, ok = getSampleFilename(dr, dataPath)
				if !ok {
					continue
				}
			}

			// read header from the raw datafile.
			csvFile, err := os.Open(dataPath)
			if err != nil {
				return nil, errors.Wrap(err, "failed to open raw data file")
			}
			defer csvFile.Close()

			reader := csv.NewReader(csvFile)
			header, err := reader.Read()
			if err != nil {
				return nil, errors.Wrap(err, "failed to read header line")
			}

			dr.Variables = AugmentVariablesFromHeader(dr, header)
		}
	}

	return meta, nil
}

func getSampleFilename(dr *model.DataResource, resourceFolder string) (string, bool) {
	// each resource type needs separate handling
	switch dr.ResType {
	case model.ResTypeTime:
		// take any file in the resource folder
		files, err := getFiles(resourceFolder)
		if err != nil || len(files) == 0 {
			return "", false
		}
		return files[0], true
	}
	return "", false
}

func getFiles(inputPath string) ([]string, error) {
	contents, err := ioutil.ReadDir(inputPath)
	if err != nil {
		return nil, errors.Wrap(err, "unable to list directory content")
	}

	files := make([]string, 0)
	for _, f := range contents {
		if !f.IsDir() {
			files = append(files, path.Join(inputPath, f.Name()))
		}
	}

	return files, nil
}

// LoadMetadataFromMergedSchema loads metadata from a merged schema file.
func LoadMetadataFromMergedSchema(schemaPath string) (*model.Metadata, error) {
	meta := &model.Metadata{
		SchemaSource: model.SchemaSourceMerged,
	}
	err := loadMergedSchema(meta, schemaPath)
	if err != nil {
		return nil, err
	}
	err = loadName(meta)
	if err != nil {
		return nil, err
	}
	err = loadID(meta)
	if err != nil {
		return nil, err
	}
	err = loadAbout(meta)
	if err != nil {
		return nil, err
	}
	err = loadMergedSchemaVariables(meta)
	if err != nil {
		return nil, err
	}
	return meta, nil
}

// LoadMetadataFromRawFile loads metadata from a raw file
// and a classification file.
func LoadMetadataFromRawFile(datasetPath string, classificationPath string) (*model.Metadata, error) {
	directory := filepath.Dir(datasetPath)
	directory = filepath.Base(directory)
	meta := &model.Metadata{
		ID:           directory,
		Name:         directory,
		StorageName:  model.NormalizeDatasetID(directory),
		SchemaSource: model.SchemaSourceRaw,
	}

	dr, err := loadRawVariables(datasetPath)
	if err != nil {
		return nil, err
	}
	meta.DataResources = []*model.DataResource{dr}

	if classificationPath != "" {
		classification, err := parseClassificationFile(classificationPath)
		if err != nil {
			return nil, err
		}
		meta.Classification = classification

		err = addClassificationTypes(meta, classificationPath)
		if err != nil {
			return nil, err
		}
	}

	return meta, nil
}

// LoadMetadataFromClassification loads metadata from a merged schema and
// classification file.
func LoadMetadataFromClassification(schemaPath string, classificationPath string, normalizeVariableNames bool, mergedFallback bool) (*model.Metadata, error) {
	meta := &model.Metadata{
		SchemaSource: model.SchemaSourceClassification,
	}

	// If classification can't be loaded, try to load from merged schema.
	classification, err := parseClassificationFile(classificationPath)
	if err != nil {
		log.Warnf("unable to load classification file: %v", err)
		if mergedFallback {
			log.Warnf("attempting to load from merged schema")
			return LoadMetadataFromMergedSchema(schemaPath)
		}

		log.Warnf("attempting to load from original schema")
		return LoadMetadataFromOriginalSchema(schemaPath, true)
	}
	meta.Classification = classification

	err = loadMergedSchema(meta, schemaPath)
	if err != nil {
		return nil, err
	}
	err = loadName(meta)
	if err != nil {
		return nil, err
	}
	err = loadID(meta)
	if err != nil {
		return nil, err
	}
	err = loadAbout(meta)
	if err != nil {
		return nil, err
	}
	err = loadClassificationVariables(meta, normalizeVariableNames)
	if err != nil {
		return nil, err
	}
	return meta, nil
}

func parseClassificationFile(classificationPath string) (*model.ClassificationData, error) {
	b, err := ioutil.ReadFile(classificationPath)
	if err != nil {
		return nil, errors.Wrap(err, "unable to read classification file")
	}

	classification := &model.ClassificationData{}
	err = json.Unmarshal(b, classification)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse classification file")
	}

	return classification, nil
}

func addClassificationTypes(m *model.Metadata, classificationPath string) error {
	classification, err := parseClassificationFile(classificationPath)
	if err != nil {
		return errors.Wrap(err, "failed to parse classification file")
	}

	for index, variable := range m.DataResources[0].Variables {
		// get suggested types
		suggestedTypes, err := parseSuggestedTypes(m, variable.Name, index, classification.Labels, classification.Probabilities)
		if err != nil {
			return err
		}
		variable.SuggestedTypes = append(variable.SuggestedTypes, suggestedTypes...)
		variable.Type = getHighestProbablySuggestedType(variable.SuggestedTypes)
	}

	return nil
}

func loadRawVariables(datasetPath string) (*model.DataResource, error) {
	// read header from the raw datafile.
	csvFile, err := os.Open(datasetPath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open raw data file")
	}
	defer csvFile.Close()

	reader := csv.NewReader(csvFile)
	fields, err := reader.Read()
	if err != nil {
		return nil, errors.Wrap(err, "failed to read header line")
	}

	// All variables now in a single dataset since it is merged
	dataResource := &model.DataResource{
		Variables: make([]*model.Variable, 0),
	}

	for index, v := range fields {
		variable := model.NewVariable(
			index,
			v,
			"",
			"",
			model.UnknownType,
			model.UnknownType,
			"",
			[]string{"attribute"},
			model.VarDistilRoleData,
			nil,
			dataResource.Variables,
			true)
		variable.Type = model.UnknownType
		dataResource.Variables = append(dataResource.Variables, variable)
	}
	return dataResource, nil
}

func loadSchema(m *model.Metadata, schemaPath string) error {
	schema, err := gabs.ParseJSONFile(schemaPath)
	if err != nil {
		return errors.Wrap(err, "failed to parse schema file")
	}
	m.Schema = schema
	return nil
}

func loadMergedSchema(m *model.Metadata, schemaPath string) error {
	schema, err := gabs.ParseJSONFile(schemaPath)
	if err != nil {
		return errors.Wrap(err, "failed to parse merged schema file")
	}
	// confirm merged schema
	if schema.Path("about.mergedSchema").Data() == nil {
		return fmt.Errorf("schema file provided is not the proper merged schema")
	}
	m.Schema = schema
	return nil
}

// LoadImportance wiull load the importance feature selection metric.
func LoadImportance(m *model.Metadata, importanceFile string) error {
	// unmarshall the schema file
	importance, err := gabs.ParseJSONFile(importanceFile)
	if err != nil {
		return errors.Wrap(err, "failed to parse importance file")
	}
	// if no numeric fields, features will be null
	// NOTE: Assume all variables in a single resource since that is
	// how we would submit to ranking.
	if importance.Path("features").Data() != nil {
		metric := importance.Path("features").Children()
		if metric == nil {
			return errors.New("features attribute missing from file")
		}
		for index, v := range m.DataResources[0].Variables {
			// geocoded variables added after ranking on ingest
			if index < len(metric) {
				v.Importance = int(metric[index].Data().(float64)) + 1
			}
		}
	}
	return nil
}

func writeSummaryFile(summaryFile string, summary string) error {
	return ioutil.WriteFile(summaryFile, []byte(summary), 0644)
}

func getSummaryFallback(str string) string {
	if len(str) < summaryLength {
		return str
	}
	return str[:summaryLength] + "..."
}

// LoadSummaryFromDescription loads a summary from the description.
func LoadSummaryFromDescription(m *model.Metadata, summaryFile string) {
	// request summary
	summary := getSummaryFallback(m.Description)
	// set summary
	m.Summary = summary
	// cache summary file
	writeSummaryFile(summaryFile, m.Summary)
}

// LoadSummary loads a description summary
func LoadSummary(m *model.Metadata, summaryFile string, useCache bool) {
	// use cache if available
	if useCache {
		b, err := ioutil.ReadFile(summaryFile)
		if err == nil {
			m.Summary = string(b)
			return
		}
	}
	LoadSummaryFromDescription(m, summaryFile)
}

// LoadSummaryMachine loads a machine-learned summary.
func LoadSummaryMachine(m *model.Metadata, summaryFile string) error {
	b, err := ioutil.ReadFile(summaryFile)
	if err != nil {
		return errors.Wrap(err, "unable to read machine-learned summary")
	}

	summary := &SummaryResult{}
	err = json.Unmarshal(b, summary)
	if err != nil {
		return errors.Wrap(err, "unable to parse machine-learned summary")
	}

	m.SummaryMachine = summary.Summary

	return nil
}

func numLines(r io.Reader) (int, error) {
	buf := make([]byte, 32*1024)
	count := 0
	lineSep := []byte{'\n'}

	for {
		c, err := r.Read(buf)
		count += bytes.Count(buf[:c], lineSep)

		switch {
		case err == io.EOF:
			return count, nil

		case err != nil:
			return count, err
		}
	}
}

// LoadDatasetStats loads the dataset and computes various stats.
func LoadDatasetStats(m *model.Metadata, datasetPath string) error {

	// open the left and outfiles for line-by-line by processing
	f, err := os.Open(datasetPath)
	if err != nil {
		return errors.Wrap(err, "failed to open dataset file")
	}

	fi, err := f.Stat()
	if err != nil {
		return errors.Wrap(err, "failed to acquire stats on dataset file")
	}

	m.NumBytes = fi.Size()

	lines, err := numLines(f)
	if err != nil {
		return errors.Wrap(err, "failed to count rows in file")
	}

	m.NumRows = int64(lines)
	return nil
}

func loadID(m *model.Metadata) error {
	id, ok := m.Schema.Path("about.datasetID").Data().(string)
	if !ok {
		return errors.Errorf("no `about.datasetID` key found in schema")
	}
	m.ID = id
	return nil
}

func loadName(m *model.Metadata) error {
	name, ok := m.Schema.Path("about.datasetName").Data().(string)
	if !ok {
		return nil //errors.Errorf("no `name` key found in schema")
	}
	m.Name = name
	return nil
}

func loadAbout(m *model.Metadata) error {
	if m.Schema.Path("about.description").Data() != nil {
		m.Description = m.Schema.Path("about.description").Data().(string)
	}

	// default to using the normalized id as storage name
	if m.Schema.Path("about.storageName").Data() != nil {
		m.StorageName = m.Schema.Path("about.storageName").Data().(string)
	} else {
		m.StorageName = model.NormalizeDatasetID(m.Schema.Path("about.datasetID").Data().(string))
	}

	if m.Schema.Path("about.redacted").Data() != nil {
		m.Redacted = m.Schema.Path("about.redacted").Data().(bool)
	}

	return nil
}

func parseSchemaVariable(v *gabs.Container, existingVariables []*model.Variable, normalizeName bool) (*model.Variable, error) {
	if v.Path("colName").Data() == nil {
		return nil, fmt.Errorf("unable to parse column name")
	}
	varName := v.Path("colName").Data().(string)

	varDisplayName := ""
	if v.Path("colDisplayName").Data() != nil {
		varDisplayName = v.Path("colDisplayName").Data().(string)
	}
	varType := ""
	if v.Path("colType").Data() != nil {
		varType = v.Path("colType").Data().(string)
		varType = model.MapLLType(varType)
	}

	varDescription := ""
	if v.Path("colDescription").Data() != nil {
		varDescription = v.Path("colDescription").Data().(string)
	}

	varIndex := 0
	if v.Path("colIndex").Data() != nil {
		varIndex = int(v.Path("colIndex").Data().(float64))
	}

	var varRoles []string
	if v.Path("role").Data() != nil {
		rolesRaw := v.Path("role").Children()
		if rolesRaw == nil {
			return nil, errors.New("unable to parse column role")
		}
		varRoles = make([]string, len(rolesRaw))
		for i, r := range rolesRaw {
			varRoles[i] = r.Data().(string)
		}
	}

	varDistilRole := ""
	if v.Path("distilRole").Data() != nil {
		varDistilRole = v.Path("distilRole").Data().(string)
	}

	varOriginalName := ""
	if v.Path("colOriginalName").Data() != nil {
		varOriginalName = v.Path("colOriginalName").Data().(string)
	}

	varOriginalType := ""
	if v.Path("colOriginalType").Data() != nil {
		varOriginalType = v.Path("colOriginalType").Data().(string)
		varOriginalType = model.MapLLType(varOriginalType)
	} else {
		varOriginalType = varType
	}

	// parse the refersTo fields to properly serialize it if necessary
	var refersTo map[string]interface{}
	if v.Path("refersTo").Data() != nil {
		refersTo = make(map[string]interface{})
		refersToData := v.Path("refersTo")
		resID := ""
		resObject := make(map[string]interface{})

		if refersToData.Path("resID").Data() != nil {
			resID = refersToData.Path("resID").Data().(string)
		}

		if refersToData.Path("resObject").Data() != nil {
			resObjectMap := refersToData.Path("resObject").ChildrenMap()
			if len(resObjectMap) == 0 {
				// see if it is maybe a string and if it is, ignore
				data, ok := refersToData.Path("resObject").Data().(string)
				if !ok {
					return nil, errors.New("unable to parse resObject")
				}
				refersTo["resObject"] = data
			} else {
				for k, v := range resObjectMap {
					resObject[k] = v.Data().(string)
				}
				refersTo["resObject"] = resObject
			}
		}

		refersTo["resID"] = resID
	}
	variable := model.NewVariable(
		varIndex,
		varName,
		varDisplayName,
		varOriginalName,
		varType,
		varOriginalType,
		varDescription,
		varRoles,
		varDistilRole,
		refersTo,
		existingVariables,
		normalizeName)

	// parse suggested types if present
	suggestedTypesJSON := v.Path("suggestedTypes").Children()
	suggestedTypes := make([]*model.SuggestedType, len(suggestedTypesJSON))
	for i, suggestion := range suggestedTypesJSON {
		stType := ""
		if suggestion.Path("type").Data() != nil {
			stType = suggestion.Path("type").Data().(string)
		}
		stProbability := -1.0
		if suggestion.Path("probability").Data() != nil {
			stProbability = suggestion.Path("probability").Data().(float64)
		}
		stProvenance := ""
		if suggestion.Path("provenance").Data() != nil {
			stProvenance = suggestion.Path("provenance").Data().(string)
		}

		suggestedTypes[i] = &model.SuggestedType{
			Type:        stType,
			Probability: stProbability,
			Provenance:  stProvenance,
		}
	}

	// no suggested type present so initialize it
	if len(suggestedTypes) == 0 {
		probability := 1.0
		if variable.Type == model.UnknownType {
			probability = 0
		} else if model.IsSchemaComplexType(variable.Type) {
			probability = 1.5
		}

		suggestedTypes = append(suggestedTypes, &model.SuggestedType{
			Type:        variable.Type,
			Probability: probability,
			Provenance:  ProvenanceSchema,
		})
	}
	variable.SuggestedTypes = suggestedTypes

	return variable, nil
}

func cleanVarType(m *model.Metadata, name string, typ string) string {
	// set the d3m index to int regardless of what gets returned
	if name == model.D3MIndexName {
		return "index"
	}
	// map types
	switch typ {
	case "int":
		return "integer"
	default:
		return typ
	}
}

func parseClassification(m *model.Metadata, index int, labels []*gabs.Container) (string, error) {
	// parse classification
	col := labels[index]
	varTypeLabels := col.Children()
	if varTypeLabels == nil {
		return "", errors.Errorf("failed to parse classification for column `%d`", col)
	}
	if len(varTypeLabels) > 0 {
		// TODO: fix so we don't always just use first classification
		return varTypeLabels[0].Data().(string), nil
	}
	return model.DefaultVarType, nil
}

func parseSuggestedTypes(m *model.Metadata, name string, index int, labels [][]string, probabilities [][]float64) ([]*model.SuggestedType, error) {
	// variables added after classification will not have suggested types
	if index >= len(labels) {
		return nil, nil
	}

	// parse probabilities
	labelsCol := labels[index]
	probabilitiesCol := probabilities[index]
	var suggested []*model.SuggestedType
	for index, typ := range labelsCol {
		probability := probabilitiesCol[index]

		// adjust the probability for complex suggested types
		if !model.IsBasicSimonType(typ) {
			probability = probability * 1.5
		}

		suggested = append(suggested, &model.SuggestedType{
			Type:        cleanVarType(m, name, typ),
			Probability: probability,
			Provenance:  ProvenanceSimon,
		})
	}
	// sort by probability
	sort.Slice(suggested, func(i, j int) bool {
		return suggested[i].Probability > suggested[j].Probability
	})
	return suggested, nil
}

// AugmentVariablesFromHeader augments the metadata variables with variables
// found in the header. All variables found in the header default to strings.
func AugmentVariablesFromHeader(dr *model.DataResource, header []string) []*model.Variable {
	// map variables by col index for quick lookup
	metaVars := make(map[int]*model.Variable)
	for _, v := range dr.Variables {
		metaVars[v.Index] = v
	}

	// add missing variables (default to string) and drop any vars not in header
	augmentedVars := make([]*model.Variable, len(header))
	for i, c := range header {
		if i < len(augmentedVars) {
			v := metaVars[i]
			if v == nil {
				v = model.NewVariable(i, c, c, c, model.UnknownType, model.UnknownType, "", []string{"attribute"}, model.VarDistilRoleData, nil, augmentedVars, true)
			}
			augmentedVars[i] = v
		}
	}

	return augmentedVars
}

func loadOriginalSchemaVariables(m *model.Metadata, schemaPath string) error {
	dataResources := m.Schema.Path("dataResources").Children()
	if dataResources == nil {
		return errors.New("failed to parse data resources")
	}

	// Parse the variables for every schema
	m.DataResources = make([]*model.DataResource, len(dataResources))
	for i, sv := range dataResources {
		if sv.Path("resType").Data() == nil {
			return fmt.Errorf("unable to parse resource type")
		}
		resType := sv.Path("resType").Data().(string)

		var parser DataResourceParser
		switch resType {
		case model.ResTypeAudio, model.ResTypeImage, model.ResTypeText:
			parser = NewMedia(resType)
		case model.ResTypeTable:
			parser = &Table{}
		case model.ResTypeTime:
			parser = &Timeseries{}
		case model.ResTypeRaw:
			parser = &Raw{
				rootPath: path.Dir(schemaPath),
			}
		default:
			return errors.Errorf("Unrecognized resource type '%s'", resType)
		}

		dr, err := parser.Parse(sv)
		if err != nil {
			return errors.Wrapf(err, "Unable to parse data resource of type '%s'", resType)
		}

		m.DataResources[i] = dr
	}
	return nil
}

func loadMergedSchemaVariables(m *model.Metadata) error {
	schemaResources := m.Schema.Path("dataResources").Children()
	if schemaResources == nil {
		return errors.New("failed to parse merged resource data")
	}

	schemaVariables := schemaResources[0].Path("columns").Children()
	if schemaVariables == nil {
		return errors.New("failed to parse merged variable data")
	}

	resPath := schemaResources[0].Path("resPath").Data().(string)

	// Merged schema has only one set of variables
	m.DataResources = make([]*model.DataResource, 1)
	m.DataResources[0] = &model.DataResource{
		Variables: make([]*model.Variable, 0),
		ResPath:   resPath,
	}

	for _, v := range schemaVariables {
		variable, err := parseSchemaVariable(v, m.DataResources[0].Variables, true)
		if err != nil {
			return errors.Wrap(err, "failed to parse merged schema variable")
		}
		m.DataResources[0].Variables = append(m.DataResources[0].Variables, variable)
	}
	return nil
}

func loadClassificationVariables(m *model.Metadata, normalizeVariableNames bool) error {
	schemaResources := m.Schema.Path("dataResources").Children()
	if schemaResources == nil {
		return errors.New("failed to parse merged resource data")
	}

	schemaVariables := schemaResources[0].Path("columns").Children()
	if schemaVariables == nil {
		return errors.New("failed to parse merged variable data")
	}

	resPath := schemaResources[0].Path("resPath").Data().(string)
	resID := schemaResources[0].Path("resID").Data().(string)

	resFormats, err := parseResFormats(schemaResources[0])
	if err != nil {
		return err
	}

	// All variables now in a single dataset since it is merged
	m.DataResources = make([]*model.DataResource, 1)
	m.DataResources[0] = &model.DataResource{
		Variables: make([]*model.Variable, 0),
		ResID:     resID,
		ResType:   model.ResTypeTable,
		ResFormat: resFormats,
		ResPath:   resPath,
	}

	for index, v := range schemaVariables {
		variable, err := parseSchemaVariable(v, m.DataResources[0].Variables, normalizeVariableNames)
		if err != nil {
			return err
		}

		// get suggested types if not loaded from schema
		loaded := false
		for _, suggestion := range variable.SuggestedTypes {
			if suggestion.Provenance == ProvenanceSimon {
				loaded = true
				break
			}
		}
		if !loaded {
			suggestedTypes, err := parseSuggestedTypes(m, variable.Name, index, m.Classification.Labels, m.Classification.Probabilities)
			if err != nil {
				return err
			}
			variable.SuggestedTypes = append(variable.SuggestedTypes, suggestedTypes...)
			variable.Type = getHighestProbablySuggestedType(variable.SuggestedTypes)
		}

		m.DataResources[0].Variables = append(m.DataResources[0].Variables, variable)
	}
	return nil
}

func getHighestProbablySuggestedType(suggestedTypes []*model.SuggestedType) string {
	typ := model.DefaultVarType
	maxProb := -math.MaxFloat64
	for _, suggestedType := range suggestedTypes {
		if suggestedType.Probability > maxProb {
			maxProb = suggestedType.Probability
			typ = suggestedType.Type
		}
	}
	return typ
}

func mergeVariables(m *model.Metadata, left []*gabs.Container, right []*gabs.Container) []*gabs.Container {
	var res []*gabs.Container
	added := make(map[string]bool)
	for _, val := range left {
		name := val.Path("varName").Data().(string)
		_, ok := added[name]
		if ok {
			continue
		}
		res = append(res, val)
		added[name] = true
	}
	for _, val := range right {
		name := val.Path("varName").Data().(string)
		_, ok := added[name]
		if ok {
			continue
		}
		res = append(res, val)
		added[name] = true
	}
	return res
}

// WriteMergedSchema exports the current meta data as a merged schema file.
func WriteMergedSchema(m *model.Metadata, path string, mergedDataResource *model.DataResource) error {
	// create output format
	output := map[string]interface{}{
		"about": map[string]interface{}{
			"datasetID":            m.ID,
			"datasetName":          m.Name,
			"parentDatasetIDs":     m.ParentDatasetIDs,
			"storageName":          m.StorageName,
			"description":          m.Description,
			"datasetSchemaVersion": schemaVersion,
			"license":              license,
			"rawData":              m.Raw,
			"redacted":             m.Redacted,
			"mergedSchema":         "true",
		},
		"dataResources": []*model.DataResource{mergedDataResource},
	}
	bytes, err := json.MarshalIndent(output, "", "	")
	if err != nil {
		return errors.Wrap(err, "failed to marshal merged schema file output")
	}
	// write copy to disk
	return ioutil.WriteFile(path, bytes, 0644)
}

// WriteClassification writes classification information to disk.
func WriteClassification(classification *model.ClassificationData, classificationPath string) error {
	b, err := json.Marshal(classification)
	if err != nil {
		return errors.Wrapf(err, "unable to marshal classification")
	}

	err = ioutil.WriteFile(classificationPath, b, os.ModePerm)
	if err != nil {
		return errors.Wrapf(err, "unable to write classification file")
	}

	return nil
}

// DatasetMatches determines if the metadata variables match.
func DatasetMatches(m *model.Metadata, variables []string) bool {
	// Assume metadata is for a merged schema, so only has 1 data resource.

	// Lengths need to be the same.
	if len(variables) != len(m.DataResources[0].Variables) {
		return false
	}

	// Build the variable lookup for matching.
	newVariable := make(map[string]bool)
	for _, v := range variables {
		newVariable[v] = true
	}

	// Make sure every existing variable is present.
	for _, v := range m.DataResources[0].Variables {
		if !newVariable[v.Name] {
			return false
		}
	}

	// Same amount of varibles, and all the names match.
	return true
}
