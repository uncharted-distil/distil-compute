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

package result

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCSVResultParser(t *testing.T) {

	start := time.Now()
	result, err := ParseResultCSV("./testdata/test.csv")
	elapsed := time.Since(start)

	t.Logf("Elapsed: %s", elapsed)

	assert.NoError(t, err)
	assert.NotEmpty(t, result)

	t.Logf("%v", result)

	assert.Equal(t, []interface{}{"idx", "col a", "col b"}, result[0])
	assert.Equal(t, []interface{}{"0", []interface{}{"alpha", "bravo"}, "foxtrot"}, result[1])
	assert.Equal(t, []interface{}{"1", []interface{}{"charlie", "delta's oscar"}, "hotel"}, result[2])
	assert.Equal(t, []interface{}{"2", []interface{}{"a", "[", "b"}, []interface{}{"c", "\"", "e"}}, result[3])
	assert.Equal(t, []interface{}{"3", []interface{}{"a", "['\"", "b"}, []interface{}{"c", "\"", "e"}}, result[4])
	assert.Equal(t, []interface{}{"4", []interface{}{"-10.001", "20.1"}, []interface{}{"30", "40"}}, result[5])
	assert.Equal(t, []interface{}{"5", []interface{}{"int"}, []interface{}{"0.989599347114563"}}, result[6])
	assert.Equal(t, []interface{}{"7", []interface{}{"int", "categorical"}, []interface{}{"0.9885959029197693", "1"}}, result[8])
	assert.Equal(t, []interface{}{"10", "( ibid )", "hotel"}, result[11])
	assert.Equal(t, []interface{}{"11", []interface{}{"int"}, []interface{}{[]interface{}{"1", "2", "3"}, []interface{}{"4", "5", "6"}, []interface{}{"7", "8", "9"}}}, result[12])
	assert.Equal(t, []interface{}{"12", []interface{}{"int"}, []interface{}{[]interface{}{"1", "2", "3"}, []interface{}{"4", "5", "6"}, []interface{}{"7", "8", "9"}}}, result[13])
}

func TestCSVResultParserShallow(t *testing.T) {

	start := time.Now()
	result, err := ParseResultCSVShallow("./testdata/test.csv")
	elapsed := time.Since(start)

	t.Logf("Elapsed: %s", elapsed)

	assert.NoError(t, err)
	assert.NotEmpty(t, result)

	t.Logf("%v", result)

	assert.Equal(t, []interface{}{"idx", "col a", "col b"}, result[0])
	assert.Equal(t, []interface{}{"0", "['alpha', 'bravo']", "foxtrot"}, result[1])
	assert.Equal(t, []interface{}{"1", "('charlie', \"delta's oscar\")", "hotel"}, result[2])
}
