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
