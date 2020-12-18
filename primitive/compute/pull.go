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
	"io"
	"time"

	"github.com/pkg/errors"
)

// PullFunc is a function that pulls from an api.
type PullFunc func() error

// PullFromAPI pulls from the api.
func PullFromAPI(maxPulls int, timeout time.Duration, pull PullFunc) error {

	recvChan := make(chan error)

	count := 0
	for {

		// pull
		go func() {
			err := pull()
			recvChan <- err
		}()

		// set timeout timer
		timer := time.NewTimer(timeout)

		select {
		case err := <-recvChan:
			timer.Stop()
			if err == io.EOF {
				return nil
			} else if err != nil {
				return errors.Wrap(err, "rpc error")
			}
			count++
			if count > maxPulls {
				return nil
			}

		case <-timer.C:
			// timeout
			return errors.Errorf("solution request has timed out")
		}

	}
}
