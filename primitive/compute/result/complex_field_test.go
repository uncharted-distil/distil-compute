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

package result

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParserSingleQuoted(t *testing.T) {
	field := &ComplexField{Buffer: "  ['c ar'  , '\\'plane', 'b* oat']"} // single quote can be escaped in python
	field.Init()

	err := field.Parse()
	assert.NoError(t, err)

	field.Execute()
	assert.Equal(t, []interface{}{"c ar", "'plane", "b* oat"}, field.arrayElements.elements)
}

func TestParserDoubleQuoted(t *testing.T) {
	field := &ComplexField{Buffer: "[\"&car\"  , \"\\plane\", \"boat's\"]"}
	field.Init()

	err := field.Parse()
	assert.NoError(t, err)

	field.Execute()
	assert.Equal(t, []interface{}{"&car", "\\plane", "boat's"}, field.arrayElements.elements)
}

func TestParserValues(t *testing.T) {
	field := &ComplexField{Buffer: "[10, 20, 30, \"forty  &*\", 4.9e-05, 4.9e05, 4.9e+05]"}
	field.Init()

	err := field.Parse()
	field.PrintSyntaxTree()
	assert.NoError(t, err)

	field.Execute()
	assert.Equal(t, []interface{}{"10", "20", "30", "forty  &*", "4.9e-05", "4.9e05", "4.9e+05"}, field.arrayElements.elements)
}

func TestParserFail(t *testing.T) {
	field := &ComplexField{Buffer: "[&*&, \"car\"  , \"plane\", \"boat's\"]"}
	field.Init()

	err := field.Parse()
	assert.Error(t, err)
}

func TestParserNested(t *testing.T) {
	field := &ComplexField{Buffer: "[[10, 20, 30, [alpha, bravo]], [40, 50, 60]]"}
	field.Init()

	err := field.Parse()
	field.PrintSyntaxTree()
	assert.NoError(t, err)

	field.Execute()

	assert.Equal(t, []interface{}{"alpha", "bravo"}, field.arrayElements.elements[0].([]interface{})[3].([]interface{}))
	assert.Equal(t, []interface{}{"40", "50", "60"}, field.arrayElements.elements[1].([]interface{}))
}

func TestParserTuple(t *testing.T) {
	field := &ComplexField{Buffer: "(10, 20, 30,)"}
	field.Init()

	err := field.Parse()
	field.PrintSyntaxTree()
	assert.NoError(t, err)

	field.Execute()

	assert.Equal(t, []interface{}{"10", "20", "30"}, field.arrayElements.elements)
}

func TestParserTupleFail(t *testing.T) {
	field := &ComplexField{Buffer: "(10)"}
	field.Init()

	err := field.Parse()
	field.PrintSyntaxTree()
	assert.Error(t, err)
}

func TestParserSingleTuple(t *testing.T) {
	field := &ComplexField{Buffer: "(10, )"}
	field.Init()

	err := field.Parse()
	field.PrintSyntaxTree()
	assert.NoError(t, err)

	field.Execute()

	assert.Equal(t, []interface{}{"10"}, field.arrayElements.elements)
}

func TestParserNestedTuple(t *testing.T) {
	field := &ComplexField{Buffer: "((10, 20, 30, (alpha, bravo)), (40, 50, 60))"}
	field.Init()

	err := field.Parse()
	field.PrintSyntaxTree()
	assert.NoError(t, err)

	field.Execute()

	assert.Equal(t, []interface{}{"alpha", "bravo"}, field.arrayElements.elements[0].([]interface{})[3].([]interface{}))
	assert.Equal(t, []interface{}{"40", "50", "60"}, field.arrayElements.elements[1].([]interface{}))
}

func TestParserNestedMixed(t *testing.T) {
	field := &ComplexField{Buffer: "([10, 20, 30, (alpha, bravo)], [40, 50, 60])"}
	field.Init()

	err := field.Parse()
	field.PrintSyntaxTree()
	assert.NoError(t, err)

	field.Execute()

	assert.Equal(t, []interface{}{"alpha", "bravo"}, field.arrayElements.elements[0].([]interface{})[3].([]interface{}))
	assert.Equal(t, []interface{}{"40", "50", "60"}, field.arrayElements.elements[1].([]interface{}))
}

func TestParserReset(t *testing.T) {
	start := time.Now()
	field := &ComplexField{Buffer: "  ['c ar'  , '\\'plane', 'b* oat']"} // single quote can be escaped in python
	field.Init()

	err := field.Parse()
	assert.NoError(t, err)

	field.Execute()

	elapsed := time.Since(start)
	t.Logf("With init: %v", elapsed)

	assert.Equal(t, []interface{}{"c ar", "'plane", "b* oat"}, field.arrayElements.elements)

	start = time.Now()

	field.Buffer = "[\"&car\"  , \"\\plane\", \"boat's\"]"
	field.Reset()

	err = field.Parse()
	assert.NoError(t, err)

	field.Execute()

	elapsed = time.Since(start)
	t.Logf("With reset: %v", elapsed)

	assert.Equal(t, []interface{}{"&car", "\\plane", "boat's"}, field.arrayElements.elements)
}
