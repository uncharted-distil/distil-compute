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

package metadata

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMetadataFromSchema(t *testing.T) {

	meta, err := LoadMetadataFromOriginalSchema("./testdata/datasetDoc.json", true)
	assert.NoError(t, err)

	assert.Equal(t, meta.Name, "test dataset")
	assert.Equal(t, meta.ID, "test_dataset")
	assert.Equal(t, meta.Description, "YOU ARE STANDING AT THE END OF A ROAD BEFORE A SMALL BRICK BUILDING.")
	assert.Equal(t, len(meta.DataResources[0].Variables), 4)
	assert.Equal(t, meta.DataResources[0].Variables[0].StorageName, "bravo")
	assert.Equal(t, meta.DataResources[0].Variables[0].Role, []string{"index"})
	assert.Equal(t, meta.DataResources[0].Variables[0].Type, "integer")
	assert.Equal(t, meta.DataResources[0].Variables[0].OriginalType, "integer")
	assert.Equal(t, meta.DataResources[0].Variables[1].StorageName, "alpha")
	assert.Equal(t, meta.DataResources[0].Variables[1].Role, []string{"attribute"})
	assert.Equal(t, meta.DataResources[0].Variables[1].Type, "string")
	assert.Equal(t, meta.DataResources[0].Variables[1].OriginalType, "string")
	assert.Equal(t, meta.DataResources[0].Variables[2].StorageName, "whiskey")
	assert.Equal(t, meta.DataResources[0].Variables[2].Role, []string{"suggestedTarget"})
	assert.Equal(t, meta.DataResources[0].Variables[2].Type, "integer")
	assert.Equal(t, meta.DataResources[0].Variables[2].OriginalType, "integer")
	assert.Equal(t, meta.DataResources[0].Variables[3].StorageName, "d3mIndex")
	assert.Equal(t, meta.DataResources[0].Variables[3].Role, []string{"index"})
	assert.Equal(t, meta.DataResources[0].Variables[3].Type, "integer")
	assert.Equal(t, meta.DataResources[0].Variables[3].OriginalType, "integer")
}
