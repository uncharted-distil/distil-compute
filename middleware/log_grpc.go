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

package middleware

import (
	"fmt"
	"time"

	"github.com/golang/protobuf/proto" //nolint need to update to new protobuf api
	"github.com/pkg/errors"
	"github.com/uncharted-distil/distil-compute/pipeline"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// MethodLogger logs GRPC method calls.
type MethodLogger interface {
	LogAPIAction(method string, params map[string]string)
}

// SolutionMessage represents API messages that have a solution ID.
type SolutionMessage interface {
	GetSolutionId() string
}

// FittedSolutionMessage represents API messages that have a fitted solution ID.
type FittedSolutionMessage interface {
	GetFittedSolutionId() string
}

// GenerateUnaryClientInterceptor creates an interceptor function that will log unary grpc calls.
func GenerateUnaryClientInterceptor(label string, trace bool, logger MethodLogger) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		if logger != nil {
			params := getParamsFromMessage(req)
			logger.LogAPIAction(method, params)
		}

		startTime := time.Now()
		newRequestLogger().
			requestType(fmt.Sprintf("%s GRPC.UNARY [SEND]", label)).
			request(method).
			message(req.(proto.Message)).
			log(true)
		err := invoker(ctx, method, req, reply, cc, opts...)
		if err != nil {
			err = errors.Wrap(err, "invoker call failed")
		}
		dt := time.Since(startTime)
		newRequestLogger().
			requestType(fmt.Sprintf("%s GRPC.UNARY [RECV]", label)).
			request(method).
			message(reply.(proto.Message)).
			duration(dt).
			log(true)
		return err
	}
}

// GenerateStreamClientInterceptor creates an interceptor function that will log grpc streaming calls.
func GenerateStreamClientInterceptor(trace bool, logger MethodLogger) grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn,
		method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		clientStream, err := streamer(ctx, desc, cc, method, opts...)
		loggingClientStream := newLoggingClientStream(&clientStream, "GRPC.STREAM_CLIENT", method, trace, logger)
		if err != nil {
			err = errors.Wrap(err, "stream create call failed")
		}
		return loggingClientStream, err
	}
}

// LoggingClientStream implements a GRPC client stream that logs output
type LoggingClientStream struct {
	grpc.ClientStream
	requestType string
	method      string
	trace       bool
	logger      MethodLogger
}

func newLoggingClientStream(c *grpc.ClientStream, requestType string, request string, trace bool, logger MethodLogger) *LoggingClientStream {
	return &LoggingClientStream{*c, requestType, request, trace, logger}
}

// RecvMsg logs messages recieved over a GRPC stream
func (c *LoggingClientStream) RecvMsg(m interface{}) error {
	err := c.ClientStream.RecvMsg(m)
	if err != nil {
		return err
	}

	request := fmt.Sprintf("%s [RECV]", c.requestType)
	if c.trace {
		newRequestLogger().
			requestType(request).
			request(c.method).
			message(m.(proto.Message)).
			log(true)
	} else {
		newRequestLogger().
			requestType(request).
			request(c.method).
			log(true)
	}
	return err
}

// SendMsg logs messages sent out over a GRPC stream
func (c *LoggingClientStream) SendMsg(m interface{}) error {
	request := fmt.Sprintf("%s [SEND]", c.requestType)
	if c.logger != nil {
		params := getParamsFromMessage(m)
		c.logger.LogAPIAction(c.method, params)
	}

	if c.trace {
		newRequestLogger().
			requestType(request).
			request(c.method).
			message(m.(proto.Message)).
			log(true)
	} else {
		newRequestLogger().
			requestType(request).
			request(c.method).
			log(true)
	}
	return c.ClientStream.SendMsg(m)
}

func getParamsFromMessage(message interface{}) map[string]string {
	params := make(map[string]string)
	solutionID, ok := extractSolutionID(message)
	if ok {
		params["solution-id"] = solutionID
	}
	fittedSolutionID, ok := extractFittedSolutionID(message)
	if ok {
		params["fitted-solution-id"] = fittedSolutionID
	}
	name, ok := extractPipelineName(message)
	if ok {
		params["name"] = name
	}

	return params
}

func extractSolutionID(message interface{}) (string, bool) {
	solutionID := ""
	solutionMessage, ok := message.(SolutionMessage)
	if ok {
		solutionID = solutionMessage.GetSolutionId()
	}

	return solutionID, ok
}

func extractPipelineName(message interface{}) (string, bool) {
	searchMsg, ok := message.(*pipeline.SearchSolutionsRequest)
	if !ok {
		return "", false
	}

	if searchMsg.Template == nil {
		return "", false
	}

	return searchMsg.Template.Name, true
}

func extractFittedSolutionID(message interface{}) (string, bool) {
	fittedSolutionID := ""
	fittedSolutionMessage, ok := message.(FittedSolutionMessage)
	if ok {
		fittedSolutionID = fittedSolutionMessage.GetFittedSolutionId()
	}

	return fittedSolutionID, ok
}
