// Copyright © 2019 Developer Network, LLC
//
// This file is subject to the terms and conditions defined in
// file 'LICENSE', which is part of this source code package.

package montecarlopi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"atomizer.io/engine"
	"devnw.com/alog"
	"github.com/google/uuid"
)

// Ensure compile time adherence to interface
var _ engine.Atom = (*MonteCarlo)(nil)

// MonteCarlo is the atom for estimating PI
type MonteCarlo struct {
	tossed    int
	tosses    chan int
	timeout   time.Duration
	conductor engine.Conductor
}

// Constants for helping determine a timeout for tosses
const denom int = 500
const modifier int = 30

// Process test method
func (mc *MonteCarlo) Process(
	ctx context.Context,
	conductor engine.Conductor,
	electron *engine.Electron,
) (result []byte, err error) {
	mc.conductor = conductor

	mc.tosses = make(chan int)

	var e = &mcelectron{}
	err = json.Unmarshal(electron.Payload, e)
	if err != nil {
		err = fmt.Errorf("error un-marshaling %s | err: %s", string(electron.Payload), err)
		return nil, err
	}

	if e.Tosses < 1 {
		return nil, errors.New("0 is not a valid toss")
	}

	// Setup the timeout with a minimum of 30 seconds
	mc.timeout = time.Second * (time.Duration(e.Tosses/denom) + time.Duration(modifier))
	mc.tossed = e.Tosses

	r := mc.estimate(ctx)

	for i := 0; i < e.Tosses; i++ {
		select {
		case <-ctx.Done():
		default:
			if err = mc.toss(ctx); err != nil {
				e.Tosses--
				alog.Warn(err, "error received on sending, attempting toss again")
			}
		}
	}

	// Get the results from the estimation function finishing processing
	select {
	case <-ctx.Done():
		return
	case calc, ok := <-r:
		if ok {
			result = calc
		}
	}

	return result, err
}

func (mc *MonteCarlo) toss(ctx context.Context) (err error) {
	e := &engine.Electron{
		ID:     uuid.New().String(),
		AtomID: engine.ID(&Toss{}),
	}

	response, err := mc.conductor.Send(ctx, e)
	if err != nil {
		return fmt.Errorf("error sending electron [%s] | %s", e.ID, err.Error())
	}

	go func(ctx context.Context, response <-chan *engine.Properties) {
		ctx, cancel := context.WithTimeout(ctx, mc.timeout)
		defer cancel()

		if response == nil {
			alog.Error(nil, "response channel from send on conductor is nil")
			return
		}

		select {
		case <-ctx.Done():
			mc.tosses <- -1
		case r, ok := <-response:
			if !ok {
				mc.tosses <- -1
				alog.Errorf(nil, "response closed prematurely for electron [%s]", e.ID)
			}

			if r.Error != nil {
				alog.Error(r.Error)
				return
			}

			t := Toss{}
			err = json.Unmarshal(r.Result, &t)
			if err != nil {
				alog.Errorf(err, "error while un-marshaling toss from %s\n", r.Result)
			}

			select {
			case <-ctx.Done():
			case mc.tosses <- t.Value:
			}
		}
	}(ctx, response)

	return err
}

func (mc *MonteCarlo) estimate(ctx context.Context) <-chan []byte {
	result := make(chan []byte)

	go func(result chan<- []byte) {
		defer close(result)

		// Until the channel closes calculate how many tosses are greater than 0
		in, tosses, errs := mc.readtosses(ctx)
		pi := mc.calculate(float64(in), float64(tosses))

		res := struct {
			In     int     `json:"in"`
			Tosses int     `json:"tosses"`
			Errors int     `json:"errors"`
			PI     float64 `json:"pi"`
		}{
			in,
			tosses,
			errs,
			pi,
		}

		if b, err := json.Marshal(res); err == nil {
			select {
			case <-ctx.Done():
				return
			case result <- b:
				alog.Printf("estimation finished, pushed result [%s] to conductor", string(b))
				return
			}
		} else {
			alog.Error(err, "error marshaling response")
		}
	}(result)

	return result
}

func (mc *MonteCarlo) readtosses(ctx context.Context) (in, tosses, errs int) {
	var count int

	// Execute the calculation
	// EstimatePI takes in the toss results and estimates PI
	for count < mc.tossed {
		count++
		select {
		case <-ctx.Done():
			return
		case v, ok := <-mc.tosses:
			if ok {
				tosses++
				if v > 0 {
					in++
				} else if v < 0 {
					errs++
					tosses--
				}
			} else {
				return
			}
		}
	}

	return in, tosses, errs
}

func (mc *MonteCarlo) calculate(in, tosses float64) float64 {
	return (4 * in) / tosses
}
