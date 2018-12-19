package description

import (
	"sort"
	"strings"

	"github.com/pkg/errors"
	"github.com/unchartedsoftware/distil-compute/model"
	"github.com/unchartedsoftware/distil-compute/pipeline"
)

const defaultResource = "0"

// CreateUserDatasetPipeline creates a pipeline description to capture user feature selection and
// semantic type information.
func CreateUserDatasetPipeline(name string, description string, allFeatures []*model.Variable,
	targetFeature string, selectedFeatures []string, filters []*model.Filter) (*pipeline.PipelineDescription, error) {

	// save the selected features in a set for quick lookup
	selectedSet := map[string]bool{}
	for _, v := range selectedFeatures {
		selectedSet[strings.ToLower(v)] = true
	}
	columnIndices := mapColumns(allFeatures, selectedSet)

	// create the semantic type update primitive
	updateSemanticTypes, err := createUpdateSemanticTypes(allFeatures, selectedSet)
	if err != nil {
		return nil, err
	}

	// create the feature selection primitive
	removeFeatures, err := createRemoveFeatures(allFeatures, selectedSet)
	if err != nil {
		return nil, err
	}

	// If neither have any content, we'll skip the template altogether.
	if len(updateSemanticTypes) == 0 && removeFeatures == nil {
		return nil, nil
	}

	filterData := createFilterData(filters, columnIndices)

	// instantiate the pipeline
	builder := NewBuilder(name, description)
	for _, v := range updateSemanticTypes {
		builder = builder.AddStep(v)
	}
	if removeFeatures != nil {
		builder = builder.AddStep(removeFeatures)
	}
	for _, f := range filterData {
		builder = builder.AddStep(f)
	}

	pip, err := builder.AddInferencePoint().Compile()
	if err != nil {
		return nil, err
	}

	// Input set to arbitrary string for now
	pip.Inputs = []*pipeline.PipelineDescriptionInput{{
		Name: "dataset",
	}}
	return pip, nil
}

func createRemoveFeatures(allFeatures []*model.Variable, selectedSet map[string]bool) (*StepData, error) {
	// create a list of features to remove
	removeFeatures := []int{}
	for _, v := range allFeatures {
		if !selectedSet[strings.ToLower(v.Name)] {
			removeFeatures = append(removeFeatures, v.Index)
		}
	}

	if len(removeFeatures) == 0 {
		return nil, nil
	}

	// instantiate the feature remove primitive
	featureSelect, err := NewRemoveColumnsStep(defaultResource, removeFeatures)
	if err != nil {
		return nil, err
	}
	return featureSelect, nil
}

type update struct {
	removeIndices []int
	addIndices    []int
}

func newUpdate() *update {
	return &update{
		addIndices:    []int{},
		removeIndices: []int{},
	}
}

func createUpdateSemanticTypes(allFeatures []*model.Variable, selectedSet map[string]bool) ([]*StepData, error) {
	// create maps of (semantic type, index list) - primitive allows for semantic types to be added to /
	// remove from multiple columns in a single operation
	updateMap := map[string]*update{}
	for _, v := range allFeatures {
		// empty selected set means all selected
		if len(selectedSet) == 0 || selectedSet[strings.ToLower(v.Name)] {
			addType := model.MapTA2Type(v.Type)
			if addType == "" {
				return nil, errors.Errorf("variable `%s` internal type `%s` can't be mapped to ta2", v.Name, v.Type)
			}
			removeType := model.MapTA2Type(v.OriginalType)
			if removeType == "" {
				return nil, errors.Errorf("remove variable `%s` internal type `%s` can't be mapped to ta2", v.Name, v.OriginalType)
			}

			// only apply change when types are different
			if addType != removeType {
				if _, ok := updateMap[addType]; !ok {
					updateMap[addType] = newUpdate()
				}
				updateMap[addType].addIndices = append(updateMap[addType].addIndices, v.Index)

				if _, ok := updateMap[removeType]; !ok {
					updateMap[removeType] = newUpdate()
				}
				updateMap[removeType].removeIndices = append(updateMap[removeType].removeIndices, v.Index)
			}
		}
	}

	// Copy the created maps into the column update structure used by the primitive.  Force
	// alpha ordering to make debugging / testing predictable
	keys := make([]string, 0, len(updateMap))
	for k := range updateMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	semanticTypeUpdates := []*StepData{}
	for _, k := range keys {
		v := updateMap[k]
		var addKey string
		if len(v.addIndices) > 0 {
			addKey = k
		}
		add := &ColumnUpdate{
			SemanticTypes: []string{addKey},
			Indices:       v.addIndices,
		}
		var removeKey string
		if len(v.removeIndices) > 0 {
			removeKey = k
		}
		remove := &ColumnUpdate{
			SemanticTypes: []string{removeKey},
			Indices:       v.removeIndices,
		}
		semanticTypeUpdate, err := NewUpdateSemanticTypeStep(defaultResource, add, remove)
		if err != nil {
			return nil, err
		}
		semanticTypeUpdates = append(semanticTypeUpdates, semanticTypeUpdate)
	}
	return semanticTypeUpdates, nil
}

func createFilterData(filters []*model.Filter, columnIndices map[string]int) []*StepData {

	// Map the fiters to pipeline primitives
	filterSteps := []*StepData{}
	for _, f := range filters {
		var filter *StepData
		inclusive := f.Mode == model.IncludeFilter
		colIndex := columnIndices[f.Key]

		switch f.Type {
		case model.NumericalFilter:
			filter = NewNumericRangeFilterStep(defaultResource, colIndex, inclusive, *f.Min, *f.Max, false)
		case model.CategoricalFilter:
			filter = NewTermFilterStep(defaultResource, colIndex, inclusive, f.Categories, true)
		case model.RowFilter:
			filter = NewTermFilterStep(defaultResource, colIndex, inclusive, f.D3mIndices, true)
		case model.FeatureFilter, model.TextFilter:
			filter = NewTermFilterStep(defaultResource, colIndex, inclusive, f.Categories, false)
		}

		filterSteps = append(filterSteps, filter)
	}
	return filterSteps
}

// CreateSlothPipeline creates a pipeline to peform timeseries clustering on a dataset.
func CreateSlothPipeline(name string, description string, timeColumn string, valueColumn string,
	timeSeriesFeatures []*model.Variable) (*pipeline.PipelineDescription, error) {

	timeIdx, err := getIndex(timeSeriesFeatures, timeColumn)
	if err != nil {
		return nil, err
	}

	valueIdx, err := getIndex(timeSeriesFeatures, valueColumn)
	if err != nil {
		return nil, err
	}

	// insantiate the pipeline
	pipeline, err := NewBuilder(name, description).
		AddStep(NewDenormalizeStep()).
		AddStep(NewDatasetToDataframeStep()).
		AddStep(NewTimeSeriesLoaderStep(-1, timeIdx, valueIdx)).
		AddStep(NewSlothStep()).
		Compile()

	if err != nil {
		return nil, err
	}

	return pipeline, nil
}

// CreateDukePipeline creates a pipeline to peform image featurization on a dataset.
func CreateDukePipeline(name string, description string) (*pipeline.PipelineDescription, error) {
	// insantiate the pipeline
	pipeline, err := NewBuilder(name, description).
		AddStep(NewDatasetToDataframeStep()).
		AddStep(NewDukeStep()).
		Compile()

	if err != nil {
		return nil, err
	}
	return pipeline, nil
}

// CreateSimonPipeline creates a pipeline to run semantic type inference on a dataset's
// columns.
func CreateSimonPipeline(name string, description string) (*pipeline.PipelineDescription, error) {
	// insantiate the pipeline
	pipeline, err := NewBuilder(name, description).
		AddStep(NewDatasetToDataframeStep()).
		AddStep(NewSimonStep()).
		Compile()

	if err != nil {
		return nil, err
	}
	return pipeline, nil
}

// CreateCrocPipeline creates a pipeline to run image featurization on a dataset.
func CreateCrocPipeline(name string, description string, targetColumns []string, outputLabels []string) (*pipeline.PipelineDescription, error) {
	// insantiate the pipeline
	pipeline, err := NewBuilder(name, description).
		AddStep(NewDenormalizeStep()).
		AddStep(NewDatasetToDataframeStep()).
		AddStep(NewCrocStep(targetColumns, outputLabels)).
		Compile()

	if err != nil {
		return nil, err
	}
	return pipeline, nil
}

// CreateUnicornPipeline creates a pipeline to run image clustering on a dataset.
func CreateUnicornPipeline(name string, description string, targetColumns []string, outputLabels []string) (*pipeline.PipelineDescription, error) {
	// insantiate the pipeline
	pipeline, err := NewBuilder(name, description).
		AddStep(NewDenormalizeStep()).
		AddStep(NewDatasetToDataframeStep()).
		AddStep(NewUnicornStep(targetColumns, outputLabels)).
		Compile()

	if err != nil {
		return nil, err
	}
	return pipeline, nil
}

// CreatePCAFeaturesPipeline creates a pipeline to run feature ranking on an input dataset.
func CreatePCAFeaturesPipeline(name string, description string) (*pipeline.PipelineDescription, error) {
	// insantiate the pipeline
	pipeline, err := NewBuilder(name, description).
		AddStep(NewDatasetToDataframeStep()).
		AddStep(NewPCAFeaturesStep()).
		Compile()

	if err != nil {
		return nil, err
	}
	return pipeline, nil
}

// CreateDenormalizePipeline creates a pipeline to run the denormalize primitive on an input dataset.
func CreateDenormalizePipeline(name string, description string) (*pipeline.PipelineDescription, error) {
	// insantiate the pipeline
	pipeline, err := NewBuilder(name, description).
		AddStep(NewDenormalizeStep()).
		AddStep(NewDatasetToDataframeStep()).
		Compile()

	if err != nil {
		return nil, err
	}
	return pipeline, nil
}

// CreateTargetRankingPipeline creates a pipeline to run feature ranking on an input dataset.
func CreateTargetRankingPipeline(name string, description string, target string, features []*model.Variable) (*pipeline.PipelineDescription, error) {
	// compute index associated with column name
	targetIdx := -1
	for _, f := range features {
		if strings.EqualFold(target, f.Name) {
			targetIdx = f.Index
			break
		}
	}
	if targetIdx < 0 {
		return nil, errors.Errorf("can't find var '%s'", name)
	}

	// ranking is dependent on user updated semantic types, so we need to make sure we apply
	// those to the original data
	updateSemanticTypeStep, err := createUpdateSemanticTypes(features, map[string]bool{})
	if err != nil {
		return nil, err
	}
	steps := []Step{}
	for _, s := range updateSemanticTypeStep {
		steps = append(steps, s)
	}

	// insantiate the pipeline
	pipeline, err := NewBuilder(name, description).
		AddStep(NewDenormalizeStep()).            // denormalize
		AddSteps(steps).                          // apply recorded semantic type changes
		AddStep(NewDatasetToDataframeStep()).     // extract main dataframe
		AddStep(NewColumnParserStep()).           // convert obj/str to python raw types
		AddStep(NewTargetRankingStep(targetIdx)). // compute feature ranks relative to target
		Compile()

	if err != nil {
		return nil, err
	}
	return pipeline, nil
}

// CreateGoatForwardPipeline creates a forward geocoding pipeline.
func CreateGoatForwardPipeline(name string, description string, source string) (*pipeline.PipelineDescription, error) {
	// insantiate the pipeline
	pipeline, err := NewBuilder(name, description).
		AddStep(NewDenormalizeStep()).        // denormalize
		AddStep(NewDatasetToDataframeStep()). // extract main dataframe
		AddStep(NewGoatForwardStep(source)).  // geocode
		Compile()

	if err != nil {
		return nil, err
	}
	return pipeline, nil
}

// CreateGoatReversePipeline creates a forward geocoding pipeline.
func CreateGoatReversePipeline(name string, description string, lonSource string, latSource string,
	features []*model.Variable) (*pipeline.PipelineDescription, error) {
	// map col names to indices
	indices := mapColumns(features, map[string]bool{lonSource: true, latSource: true})
	if len(indices) != 2 {
		return nil, errors.Errorf("can't find one of %s, %s", lonSource, latSource)
	}

	// insantiate the pipeline
	pipeline, err := NewBuilder(name, description).
		AddStep(NewDenormalizeStep()).                                       // denormalize
		AddStep(NewDatasetToDataframeStep()).                                // extract main dataframe
		AddStep(NewGoatReverseStep(indices[lonSource], indices[latSource])). // geocode
		Compile()

	if err != nil {
		return nil, err
	}
	return pipeline, nil
}

func mapColumns(allFeatures []*model.Variable, selectedSet map[string]bool) map[string]int {
	colIndices := make(map[string]int)
	index := 0
	for _, f := range allFeatures {
		if selectedSet[f.Name] {
			colIndices[f.Name] = index
			index = index + 1
		}
	}

	return colIndices
}

func getIndex(allFeatures []*model.Variable, name string) (int, error) {
	for _, f := range allFeatures {
		if strings.EqualFold(name, f.Name) {
			return f.Index, nil
		}
	}
	return -1, errors.Errorf("can't find var '%s'", name)
}
