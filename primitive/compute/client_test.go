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

package compute

/*
import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClient(t *testing.T) {

	client, err := NewClient("localhost:45042", "./datasets", true)
	assert.NoError(t, err)

	searchSolutionsRequest := &SearchSolutionsRequest{
		Problem: &ProblemDescription{
			Problem: &Problem{
				TaskType: TaskType_REGRESSION,
				PerformanceMetrics: []*ProblemPerformanceMetric{
					&ProblemPerformanceMetric{
						Metric: PerformanceMetric_MEAN_SQUARED_ERROR,
					},
				},
			},
			Inputs: []*ProblemInput{
				&ProblemInput{
					DatasetId: "196_autoMpg",
					Targets: []*ProblemTarget{
						&ProblemTarget{
							TargetIndex: 0,
							ResourceId:  "learningData",
							ColumnIndex: 8,
							ColumnName:  "class",
						},
					},
				},
			},
		},
	}

	searchID, err := client.StartSearch(context.Background(), searchSolutionsRequest)
	assert.NoError(t, err)

	solutions, err := client.SearchSolutions(context.Background(), searchID)
	assert.NoError(t, err)

	for _, solution := range solutions {

		assert.NotEmpty(t, solution.SolutionId)

		_, err := client.GenerateSolutionScores(context.Background(), solution.SolutionId)
		assert.NoError(t, err)

		_, err = client.GenerateSolutionFit(context.Background(), solution.SolutionId)
		assert.NoError(t, err)

		produceSolutionRequest := &ProduceSolutionRequest{
			SolutionId: solution.SolutionId,
			Inputs: []*Value{
				{
					Value: &Value_DatasetUri{
						DatasetUri: "testURI",
					},
				},
			},
		}

		_, err = client.GeneratePredictions(context.Background(), produceSolutionRequest)
		assert.NoError(t, err)
	}

	err = client.EndSearch(context.Background(), searchID)
	assert.NoError(t, err)
}
*/
