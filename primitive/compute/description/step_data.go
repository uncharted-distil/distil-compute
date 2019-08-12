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

package description

import (
	"fmt"
	"reflect"

	"github.com/pkg/errors"
	"github.com/uncharted-distil/distil-compute/pipeline"
)

const (
	stepInputsKey       = "inputs"
	stepKey             = "steps"
	pipelineOutputsName = "outputs"
)

// Step provides data for a pipeline description step and an operation
// to create a protobuf PipelineDescriptionStep from that data.
type Step interface {
	BuildDescriptionStep() (*pipeline.PipelineDescriptionStep, error)
	GetPrimitive() *pipeline.Primitive
	GetArguments() map[string]DataRef
	GetHyperparameters() map[string]interface{}
	GetOutputMethods() []string
}

// StepData contains the minimum amount of data used to describe a pipeline step
type StepData struct {
	Primitive       *pipeline.Primitive
	Arguments       map[string]DataRef
	Hyperparameters map[string]interface{}
	OutputMethods   []string
}

// PrimitiveReference is a marker struct for a an argument that is an integer, but should be interpreted
// as a reference to another primitive.
type PrimitiveReference struct {
	index int
}

// DataRef defines an interface for creating input reference strings that are used to connect primitive inputs
// to outputs.
type DataRef interface {
	RefString() string
}

// StepDataRef provides an input reference that points a step in the pipline, and its associated
// output method name.
type StepDataRef struct {
	stepNum int
	output  string
}

// RefString creates a string representation of a StepDataRef.
func (i *StepDataRef) RefString() string {
	return fmt.Sprintf("%s.%d.%s", stepKey, i.stepNum, i.output)
}

// PipelineDataRef points to an input of the pipeline itself (typically a dataset) instead of the output
// of another primitive.
type PipelineDataRef struct {
	inputNum int
}

// RefString creates a string representation of a PipelineDataRef.
func (p *PipelineDataRef) RefString() string {
	return fmt.Sprintf("%s.%d", stepInputsKey, p.inputNum)
}

// NewStepData creates a pipeline step instance from the required field subset.  Hyperparameters, Arguments and
// OutputMethods are all opional in the d3m runtime so so nil is a valid value.  Valid types fror hyper parameters are intXX, string, bool,
// or PrimitiveRef, which is the index of another primitive in the pipeline.
func NewStepData(
	primitive *pipeline.Primitive,
	outputMethods []string,
	hyperparameters map[string]interface{},
	arguments map[string]DataRef) *StepData {

	return &StepData{
		Primitive:       primitive,
		Hyperparameters: hyperparameters, // optional, nil is valid
		Arguments:       arguments,
		OutputMethods:   outputMethods,
	}
}

// GetPrimitive returns a primitive definition for a pipeline step.
func (s *StepData) GetPrimitive() *pipeline.Primitive {
	return s.Primitive
}

// GetArguments returns a map of arguments that will be passed to the methods
// of the primitive step.
func (s *StepData) GetArguments() map[string]DataRef {
	return s.Arguments
}

// GetHyperparameters returns a map of arguments that will be passed to the primitive methods
// of the primitive step.  Types are currently restricted to intXX, bool, string
func (s *StepData) GetHyperparameters() map[string]interface{} {
	return s.Hyperparameters
}

// GetOutputMethods returns a list of methods that will be called to generate
// primitive output.  These feed into downstream primitives.
func (s *StepData) GetOutputMethods() []string {
	return s.OutputMethods
}

// BuildDescriptionStep creates protobuf structs from step data.
func (s *StepData) BuildDescriptionStep() (*pipeline.PipelineDescriptionStep, error) {
	// generate arguments entries
	arguments := map[string]*pipeline.PrimitiveStepArgument{}
	for argName, argDataRef := range s.Arguments {
		arguments[argName] = &pipeline.PrimitiveStepArgument{
			// only handle container args rights now - extend to others if required
			Argument: &pipeline.PrimitiveStepArgument_Container{
				Container: &pipeline.ContainerArgument{
					Data: argDataRef.RefString(),
				},
			},
		}
	}

	// Generate hyper parameter entries - accepted go-natives types are currently intXX, string, bool, as well as list, map[string]
	// of those types.  Primitives are also accepted.
	hyperparameters := map[string]*pipeline.PrimitiveStepHyperparameter{}
	for k, v := range s.Hyperparameters {
		// We can handle hyperparameter args that are values, or primitive references
		var primitiveRef *PrimitiveReference
		switch t := v.(type) {
		case PrimitiveReference:
			primitiveRef = &t
		case *PrimitiveReference:
			primitiveRef = t
		}

		if primitiveRef != nil {
			// Primitive reference.  This is an int value that is the index of the primitive in protobuf / d3m pipeline structure.
			// The actual value will get filled in as a post process since we add all the primitive args in at the end.
			hyperparameters[k] = &pipeline.PrimitiveStepHyperparameter{
				Argument: &pipeline.PrimitiveStepHyperparameter_Primitive{
					Primitive: &pipeline.PrimitiveArgument{
						Data: int32(primitiveRef.index),
					},
				},
			}
		} else {
			// Values
			rawValue, err := parseValue(v)
			if err != nil {
				return nil, errors.Errorf("compile failed: hyperparameter `%s` unsupported - %s", k, err.Error())
			}
			hyperparameters[k] = &pipeline.PrimitiveStepHyperparameter{
				Argument: &pipeline.PrimitiveStepHyperparameter_Value{
					Value: &pipeline.ValueArgument{
						Data: &pipeline.Value{
							Value: &pipeline.Value_Raw{
								Raw: rawValue,
							},
						},
					},
				},
			}
		}
	}

	// list of methods that will generate output - order matters because the steps are
	// numbered
	outputMethods := []*pipeline.StepOutput{}
	for _, outputMethod := range s.OutputMethods {
		outputMethods = append(outputMethods,
			&pipeline.StepOutput{
				Id: outputMethod,
			})
	}

	// create the pipeline description structure
	descriptionStep := &pipeline.PipelineDescriptionStep{
		Step: &pipeline.PipelineDescriptionStep_Primitive{
			Primitive: &pipeline.PrimitivePipelineDescriptionStep{
				Primitive:   s.Primitive,
				Arguments:   arguments,
				Hyperparams: hyperparameters,
				Outputs:     outputMethods,
			},
		},
	}

	return descriptionStep, nil
}

func parseList(v interface{}) (*pipeline.ValueRaw, error) {
	// parse list contents as a list, map, or value
	valueList := []*pipeline.ValueRaw{}
	var value *pipeline.ValueRaw
	var err error

	// type switches to work well with generic arrays/maps so we have to revert to using reflection
	refValue := reflect.ValueOf(v)
	if refValue.Kind() != reflect.Slice {
		return nil, errors.Errorf("unexpected parameter %s", refValue.Kind())
	}
	for i := 0; i < refValue.Len(); i++ {
		refElement := refValue.Index(i)
		switch refElement.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.String, reflect.Bool, reflect.Float32, reflect.Float64:
			value, err = parseValue(refElement.Interface())
		case reflect.Slice:
			value, err = parseList(refElement.Interface())
		case reflect.Map:
			value, err = parseMap(refElement.Interface())
		default:
			err = errors.Errorf("unhandled list arg type %s", refElement.Kind())
		}

		if err != nil {
			return nil, err
		}

		valueList = append(valueList, value)
	}
	rawValue := &pipeline.ValueRaw{
		Raw: &pipeline.ValueRaw_List{
			List: &pipeline.ValueList{
				Items: valueList,
			},
		},
	}
	return rawValue, nil
}

func parseMap(vmap interface{}) (*pipeline.ValueRaw, error) {
	// parse map contents as list, map or value
	valueMap := map[string]*pipeline.ValueRaw{}
	var value *pipeline.ValueRaw
	var err error

	// type switches to work well with generic arrays/maps so we have to revert to using reflection
	refValue := reflect.ValueOf(vmap)
	if refValue.Kind() != reflect.Map {
		return nil, errors.Errorf("unexpected parameter %s", refValue.Kind())
	}
	keys := refValue.MapKeys()
	for _, key := range keys {

		if key.Kind() != reflect.String {
			return nil, errors.Errorf("non-string map key type %s", refValue.Kind())
		}

		refElement := refValue.MapIndex(key)
		switch refElement.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.String, reflect.Bool, reflect.Float32, reflect.Float64:
			value, err = parseValue(refElement.Interface())
		case reflect.Slice:
			value, err = parseList(refElement.Interface())
		case reflect.Map:
			value, err = parseMap(refElement.Interface())
		default:
			err = errors.Errorf("unhandled map arg type %s", refElement.Kind())
		}

		if err != nil {
			return nil, err
		}
		valueMap[key.String()] = value
	}

	v := &pipeline.ValueRaw{
		Raw: &pipeline.ValueRaw_Dict{
			Dict: &pipeline.ValueDict{
				Items: valueMap,
			},
		},
	}
	return v, nil
}

func parseValue(v interface{}) (*pipeline.ValueRaw, error) {
	refValue := reflect.ValueOf(v)
	switch refValue.Kind() {
	// parse a numeric, string or boolean value
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return &pipeline.ValueRaw{
			Raw: &pipeline.ValueRaw_Int64{
				Int64: refValue.Int(),
			},
		}, nil
	case reflect.Float32, reflect.Float64:
		return &pipeline.ValueRaw{
			Raw: &pipeline.ValueRaw_Double{
				Double: refValue.Float(),
			},
		}, nil
	case reflect.String:
		return &pipeline.ValueRaw{
			Raw: &pipeline.ValueRaw_String_{
				String_: refValue.String(),
			},
		}, nil
	case reflect.Bool:
		return &pipeline.ValueRaw{
			Raw: &pipeline.ValueRaw_Bool{
				Bool: refValue.Bool(),
			},
		}, nil
	case reflect.Slice:
		return parseList(v)
	case reflect.Map:
		return parseMap(v)
	default:
		return nil, errors.Errorf("unhandled value arg type %s", refValue.Kind())
	}
}
