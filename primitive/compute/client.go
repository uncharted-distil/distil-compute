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

package compute

import (
	"fmt"
	"io"
	"os"
	"path"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/uncharted-distil/distil-compute/middleware"
	"github.com/uncharted-distil/distil-compute/pipeline"
	log "github.com/unchartedsoftware/plog"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const (
	// DefaultResourceID is the default name of the main resource id in a schema.
	DefaultResourceID = "learningData"
	// DefaultExposedOutputKey is the standard output key applied to pipelines.
	DefaultExposedOutputKey = "outputs.0"
	// SolutionPendingStatus represents that the solution request has been acknoledged by not yet sent to the API
	SolutionPendingStatus = "SOLUTION_PENDING"
	// SolutionFittingStatus represents that the solution request has been sent to the API.
	SolutionFittingStatus = "SOLUTION_FITTING"
	// SolutionScoringStatus represents that the solution request has been sent to the API.
	SolutionScoringStatus = "SOLUTION_SCORING"
	// SolutionProducingStatus represents that the solution request has been sent to the API.
	SolutionProducingStatus = "SOLUTION_PRODUCING"
	// SolutionErroredStatus represents that the solution request has terminated with an error.
	SolutionErroredStatus = "SOLUTION_ERRORED"
	// SolutionCompletedStatus represents that the solution request has completed successfully.
	SolutionCompletedStatus = "SOLUTION_COMPLETED"
	// SolutionCancelledStatus represents that the solution request was intentionally cancelled.
	SolutionCancelledStatus = "SOLUTION_CANCELLED"
	// RequestPendingStatus represents that the solution request has been acknoledged by not yet sent to the API
	RequestPendingStatus = "REQUEST_PENDING"
	// RequestRunningStatus represents that the solution request has been sent to the API.
	RequestRunningStatus = "REQUEST_RUNNING"
	// RequestErroredStatus represents that the solution request has terminated with an error.
	RequestErroredStatus = "REQUEST_ERRORED"
	// RequestCompletedStatus represents that the solution request has completed successfully.
	RequestCompletedStatus = "REQUEST_COMPLETED"
	// D3MLearningData provides the name of the training csv file as defined in the D3M schema
	D3MLearningData = "learningData.csv"
	// DistilParquetLearningData provides the name of the training parquet file introduced by Distil
	DistilParquetLearningData = "learningData.parquet"
	// D3MDataFolder provides the name of the directory containing the dataset
	D3MDataFolder = "tables"
	// D3MDataSchema provides the name of the D3M data schema file
	D3MDataSchema = "datasetDoc.json"
	// D3MDatasetSchemaVersion is the current version supported when persisting
	D3MDatasetSchemaVersion = "3.0"
	// D3MResourceType is the resource type of persisted datasets
	D3MResourceType = "table"
	// D3MResourceFormat is the resource format of persisted dataset
	D3MResourceFormat = "text/csv"
	// DistilParquesetResourceFormat is the parquet resource format of persisted dataset introduced by Distil
	DistilParquetResourceFormat = "application/parquet"

	defaultTrainTestRatio = 0.8
)

// Client provides facilities for managing GPRC solution requests. Requests are
// isssued and a context object containing rx channels is returned to the caller for consumption
// of results. The context for running requests can also be fetched, along with their buffered
// results. Spawning a grpc.ClientConn per RPC call is not considered good practice - the system
// is designed such that multiple go routines make RPC calls to a single shared client, and synch
// is managed internally.
type Client struct {
	client            pipeline.CoreClient
	conn              *grpc.ClientConn
	runner            *grpc.ClientConn
	UserAgent         string
	PullTimeout       time.Duration
	PullMax           int
	SkipPreprocessing bool
}

// SearchSolutionHandler is executed when a new search solution is returned.
type SearchSolutionHandler func(*pipeline.GetSearchSolutionsResultsResponse)

// NewClient creates a new pipline request dispatcher instance. This will establish
// the connection to the solution server or return an error on fail
func NewClient(serverAddr string, trace bool, userAgent string, label string,
	pullTimeout time.Duration, pullMax int, skipPreprocessing bool, logger middleware.MethodLogger) (*Client, error) {
	log.Infof("connecting to ta2 at %s", serverAddr)

	conn, err := grpc.Dial(
		serverAddr,
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithUnaryInterceptor(middleware.GenerateUnaryClientInterceptor(label, trace, logger)),
		grpc.WithStreamInterceptor(middleware.GenerateStreamClientInterceptor(trace, logger)),
	)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to connect to %s", serverAddr)
	}

	client := Client{
		client:            pipeline.NewCoreClient(conn),
		conn:              conn,
		UserAgent:         userAgent,
		PullTimeout:       pullTimeout,
		PullMax:           pullMax,
		SkipPreprocessing: skipPreprocessing,
	}

	return &client, nil
}

// Hello does the hello message exchange to check for basic connectivity.
func (c *Client) Hello() (string, error) {
	// check for basic ta2 connectivity
	helloResponse, err := c.client.Hello(context.Background(), &pipeline.HelloRequest{})
	if err != nil {
		return "", err
	}

	// Get the API version
	version := helloResponse.GetVersion()

	log.Infof("ta2 user agent: %s", helloResponse.GetUserAgent())
	log.Infof("ta2 API version: %s", version)
	log.Infof("ta2 Allowed value types: %+v", helloResponse.GetAllowedValueTypes())
	log.Infof("ta2 extensions: %+v", helloResponse.GetSupportedExtensions())

	if !strings.EqualFold(GetAPIVersion(), version) {
		log.Warnf("ta2 API version '%s' does not match expected version '%s", version, GetAPIVersion())
	}

	log.Infof("connected to %s", c.conn.Target())

	return version, nil
}

// Close the connection to the solution service
func (c *Client) Close() {
	log.Infof("client connection closed")
	c.conn.Close()
}

// StartSearch starts a solution search session.
func (c *Client) StartSearch(ctx context.Context, request *pipeline.SearchSolutionsRequest) (string, error) {

	searchSolutionResponse, err := c.client.SearchSolutions(ctx, request)
	if err != nil {
		return "", errors.Wrap(err, "failed to start search")
	}

	return searchSolutionResponse.SearchId, nil
}

// SearchSolutions generates candidate pipelines and executes a provided handler
// for each result. While handlers are executing asynchronously, this method
// will not return until all handlers have finished.
func (c *Client) SearchSolutions(ctx context.Context, searchID string, solutionHandler SearchSolutionHandler) error {

	searchPiplinesResultsRequest := &pipeline.GetSearchSolutionsResultsRequest{
		SearchId: searchID,
	}

	searchSolutionsResultsResponse, err := c.client.GetSearchSolutionsResults(ctx, searchPiplinesResultsRequest)
	if err != nil {
		return errors.Wrap(err, "failed to open search results stream")
	}

	// track handlers to ensure they all finish before returning
	wg := &sync.WaitGroup{}

	err = PullFromAPI(c.PullMax, c.PullTimeout, func() error {
		solutionResultResponse, err := searchSolutionsResultsResponse.Recv()
		if err == io.EOF {
			return nil
		}

		if err != nil {
			return errors.Wrap(err, "failed to get search result")
		}
		// ignore empty responses
		if solutionResultResponse.SolutionId != "" {
			wg.Add(1)
			go func() {
				solutionHandler(solutionResultResponse)
				wg.Done()
			}()
		}
		return nil
	})
	if err != nil {
		return err
	}

	// don't return until all handlers have finished executing
	wg.Wait()
	return nil
}

// GenerateSolutionScores generates scrores for candidate solutions.
func (c *Client) GenerateSolutionScores(ctx context.Context, solutionID string, datasetURI string, metrics []string) ([]*pipeline.GetScoreSolutionResultsResponse, error) {

	scoreSolutionRequest := &pipeline.ScoreSolutionRequest{
		SolutionId: solutionID,
		Inputs: []*pipeline.Value{
			{
				Value: &pipeline.Value_DatasetUri{
					DatasetUri: datasetURI,
				},
			},
		},
		PerformanceMetrics: ConvertMetricsFromTA3ToTA2(metrics),
		Configuration: &pipeline.ScoringConfiguration{
			Method:         HoldoutEvaluationMethod,
			TrainTestRatio: defaultTrainTestRatio,
		},
	}

	scoreSolutionResponse, err := c.client.ScoreSolution(ctx, scoreSolutionRequest)
	if err != nil {
		return nil, errors.Wrap(err, "failed to start solution scoring")
	}

	searchPiplinesResultsRequest := &pipeline.GetScoreSolutionResultsRequest{
		RequestId: scoreSolutionResponse.RequestId,
	}

	scoreSolutionResultsResponse, err := c.client.GetScoreSolutionResults(ctx, searchPiplinesResultsRequest)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open solution scoring results stream")
	}

	var solutionResultResponses []*pipeline.GetScoreSolutionResultsResponse

	err = PullFromAPI(c.PullMax, c.PullTimeout, func() error {
		solutionResultResponse, err := scoreSolutionResultsResponse.Recv()
		if err == io.EOF {
			return nil
		}

		if err != nil {
			return errors.Wrap(err, "failed to receive solution scoring result")
		}
		solutionResultResponses = append(solutionResultResponses, solutionResultResponse)
		return nil
	})
	if err != nil {
		return nil, err
	}

	return solutionResultResponses, nil
}

// GenerateSolutionFit generates fit for candidate solutions.
func (c *Client) GenerateSolutionFit(ctx context.Context, request *pipeline.FitSolutionRequest) ([]*pipeline.GetFitSolutionResultsResponse, error) {
	fitSolutionResponse, err := c.client.FitSolution(ctx, request)
	if err != nil {
		return nil, errors.Wrap(err, "failed to start solution fitting")
	}

	fitSolutionResultsRequest := &pipeline.GetFitSolutionResultsRequest{
		RequestId: fitSolutionResponse.RequestId,
	}

	fitSolutionResultsResponse, err := c.client.GetFitSolutionResults(ctx, fitSolutionResultsRequest)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open solution fitting result stream")
	}

	var solutionResultResponses []*pipeline.GetFitSolutionResultsResponse

	err = PullFromAPI(c.PullMax, c.PullTimeout, func() error {
		solutionResultResponse, err := fitSolutionResultsResponse.Recv()
		if err == io.EOF {
			return nil
		}

		if err != nil {
			return errors.Wrap(err, "failed to receive solution fitting result")
		}
		solutionResultResponses = append(solutionResultResponses, solutionResultResponse)
		return nil
	})
	if err != nil {
		return nil, err
	}

	return solutionResultResponses, nil
}

// GetSolutionDescription gets the solution description via API call.
func (c *Client) GetSolutionDescription(ctx context.Context, solutionID string) (*pipeline.DescribeSolutionResponse, error) {
	req := &pipeline.DescribeSolutionRequest{
		SolutionId: solutionID,
	}
	desc, err := c.client.DescribeSolution(ctx, req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open solution produce result stream")
	}

	return desc, nil
}

// GeneratePredictions generates predictions.
func (c *Client) GeneratePredictions(ctx context.Context, request *pipeline.ProduceSolutionRequest) (string, []*pipeline.GetProduceSolutionResultsResponse, error) {

	produceSolutionResponse, err := c.client.ProduceSolution(ctx, request)
	if err != nil {
		return "", nil, errors.Wrap(err, "failed to start solution produce")
	}

	produceSolutionResultsRequest := &pipeline.GetProduceSolutionResultsRequest{
		RequestId: produceSolutionResponse.RequestId,
	}

	produceSolutionResultsResponse, err := c.client.GetProduceSolutionResults(ctx, produceSolutionResultsRequest)
	if err != nil {
		return "", nil, errors.Wrap(err, "failed to open solution produce result stream")
	}

	var solutionResultResponses []*pipeline.GetProduceSolutionResultsResponse

	err = PullFromAPI(c.PullMax, c.PullTimeout, func() error {
		solutionResultResponse, err := produceSolutionResultsResponse.Recv()
		if err == io.EOF {
			return nil
		}

		if err != nil {
			return errors.Wrap(err, "failed to receive solution produce result")
		}
		solutionResultResponses = append(solutionResultResponses, solutionResultResponse)
		return nil
	})
	if err != nil {
		return "", nil, err
	}

	return produceSolutionResponse.RequestId, solutionResultResponses, nil
}

// StopSearch stop the solution search session.
func (c *Client) StopSearch(ctx context.Context, searchID string) error {

	stopSearchSolutions := &pipeline.StopSearchSolutionsRequest{
		SearchId: searchID,
	}

	_, err := c.client.StopSearchSolutions(ctx, stopSearchSolutions)
	return errors.Wrap(err, "failed to stop solution search")
}

// EndSearch ends the solution search session.
func (c *Client) EndSearch(ctx context.Context, searchID string) error {

	endSearchSolutions := &pipeline.EndSearchSolutionsRequest{
		SearchId: searchID,
	}

	_, err := c.client.EndSearchSolutions(ctx, endSearchSolutions)
	return errors.Wrap(err, "failed to end solution search")
}

// ExportSolution exports the solution.
func (c *Client) ExportSolution(ctx context.Context, solutionID string) error {
	exportSolution := &pipeline.SolutionExportRequest{
		Rank:       1,
		SolutionId: solutionID,
	}
	_, err := c.client.SolutionExport(ctx, exportSolution)
	return errors.Wrap(err, "failed to export solution")
}

// SaveSolution saves the solution.
func (c *Client) SaveSolution(ctx context.Context, solutionID string) (string, error) {
	saveSolution := &pipeline.SaveSolutionRequest{
		SolutionId: solutionID,
	}
	resp, err := c.client.SaveSolution(ctx, saveSolution)
	if err != nil {
		return "", errors.Wrap(err, "failed to save solution")
	}

	return resp.GetSolutionUri(), nil
}

// SaveFittedSolution saves the fitted solution.
func (c *Client) SaveFittedSolution(ctx context.Context, fittedSolutionID string) (string, error) {
	saveSolution := &pipeline.SaveFittedSolutionRequest{
		FittedSolutionId: fittedSolutionID,
	}
	resp, err := c.client.SaveFittedSolution(ctx, saveSolution)
	if err != nil {
		return "", errors.Wrap(err, "failed to save fitted solution")
	}

	return resp.GetFittedSolutionUri(), nil
}

// LoadSolution will have the server load a solution by URI.
func (c *Client) LoadSolution(ctx context.Context, solutionURI string) (string, error) {
	loadSolution := &pipeline.LoadSolutionRequest{
		SolutionUri: solutionURI,
	}
	resp, err := c.client.LoadSolution(ctx, loadSolution)
	if err != nil {
		return "", errors.Wrap(err, "failed to load solution")
	}

	return resp.GetSolutionId(), nil
}

// LoadFittedSolution will have the server load a fitted solution by URI.
func (c *Client) LoadFittedSolution(ctx context.Context, fittedSolutionURI string) (string, error) {
	loadSolution := &pipeline.LoadFittedSolutionRequest{
		FittedSolutionUri: fittedSolutionURI,
	}
	resp, err := c.client.LoadFittedSolution(ctx, loadSolution)
	if err != nil {
		return "", errors.Wrap(err, "failed to load fitted solution")
	}

	return resp.GetFittedSolutionId(), nil
}

// ExecutePipeline executes a pre-specified pipeline.
func (c *Client) ExecutePipeline(ctx context.Context, datasetURIs []string, pipelineDesc *pipeline.PipelineDescription) (*pipeline.PipelineExecuteResponse, error) {
	in := &pipeline.PipelineExecuteRequest{
		PipelineDescription: pipelineDesc,
		Inputs:              createInputValues(datasetURIs),
	}
	out := new(pipeline.PipelineExecuteResponse)
	err := c.runner.Invoke(ctx, "/Executor/ExecutePipeline", in, out)
	if err != nil {
		return nil, err
	}

	return out, nil
}

// createInputValues creates a protobuf value structure for each supplied dataset URI.
func createInputValues(datasetURIs []string) []*pipeline.Value {
	inputs := []*pipeline.Value{}
	for _, uri := range datasetURIs {
		datasetURI := BuildSchemaFileURI(uri)
		value := &pipeline.Value{
			Value: &pipeline.Value_DatasetUri{
				DatasetUri: datasetURI,
			},
		}
		inputs = append(inputs, value)
	}
	return inputs
}

// BuildSchemaFileURI updates a URI to match the expected structure of a
// the TA2TA3 API.
func BuildSchemaFileURI(uri string) string {
	formattedURI := uri
	// make sure it ends with the standard schema doc name
	if isDirectory(formattedURI) {
		log.Debugf("adding %s to %s since it is a directory", D3MDataSchema, formattedURI)
		formattedURI = path.Join(formattedURI, D3MDataSchema)
	} else if !strings.HasSuffix(formattedURI, D3MDataSchema) {
		log.Debugf("replacing file name in %s because it isnt the expected %s", formattedURI, D3MDataSchema)
		formattedURI = path.Join(path.Dir(formattedURI), D3MDataSchema)
	}

	// check for protocol
	if !strings.HasPrefix(formattedURI, "file://") {
		log.Debugf("adding file protocol to %s because it does not have it", formattedURI)
		formattedURI = fmt.Sprintf("file://%s", formattedURI)
	}

	log.Debugf("schema file uri: %s", formattedURI)

	return formattedURI
}

func isDirectory(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false
	}
	return fileInfo.IsDir()
}
