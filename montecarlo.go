package montecarlopi

import (
	"context"
	"encoding/json"

	"github.com/benjivesterby/alog"
	"github.com/benjivesterby/atomizer"
	"github.com/google/uuid"
)

// MonteCarlo is the atom for estimating PI
type MonteCarlo struct {
	tossed    int
	tosses    chan int
	conductor atomizer.Conductor
}

// ID test method
func (mc *MonteCarlo) ID() string {
	return "montecarlo"
}

// Process test method
func (mc *MonteCarlo) Process(ctx context.Context, conductor atomizer.Conductor, electron *atomizer.Electron) (result []byte, err error) {
	mc.conductor = conductor

	mc.tosses = make(chan int)

	var e = &mcelectron{}
	if err = json.Unmarshal(electron.Payload, e); err == nil {

		if e.Tosses > 0 {
			mc.tossed = e.Tosses

			r := mc.estimate(ctx)

			for i := 0; i < e.Tosses; i++ {
				if err = mc.toss(ctx); err != nil {
					e.Tosses--
					alog.Warn(err, "error received on sending, attempting toss again")
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
		} else {
			// TODO:
		}
	} else {
		alog.Errorf(err, "error un-marshalling %s", string(electron.Payload))
	}

	return result, err

	// Step 1: parse my electron - How many tosses?

	// Step 2: Push the toss electrons onto the outbound channel

	// Step 3: Wait for the callbacks on the outbound electrons to finish processing

	// Step 4: Calculate PI estimation using the returned values from each toss "atom"

	// Step 5: Return the results on the results channel to the atomizer
}

func (mc *MonteCarlo) toss(ctx context.Context) (err error) {
	//duration := time.Second * 60

	e := &atomizer.Electron{
		ID:     uuid.New().String(),
		AtomID: "toss",
		//Timeout:    &duration,
	}

	var response <-chan *atomizer.Properties
	if response, err = mc.conductor.Send(ctx, e); err == nil {

		go func(ctx context.Context, response <-chan *atomizer.Properties) {
			// ctx, cancel := context.WithTimeout(ctx, duration)
			// defer cancel()

			if response != nil {
				select {
				case <-ctx.Done():
					select {
					case mc.tosses <- -1:
						alog.Error(nil, "context closed, returning toss error")
					}
				case r, ok := <-response:
					if ok {
						if r.Error == nil {
							t := &Toss{}
							if err := json.Unmarshal(r.Result, t); err == nil {

								select {
								case <-ctx.Done():
								case mc.tosses <- t.Value:
								}
							} else {
								alog.Errorf(err, "error while un-marshalling toss from %s\n", r.Result)
							}
						} else {
							alog.Error(r.Error, "error in toss response, sending error code")
						}
					} else {
						select {
						case mc.tosses <- -1:
							alog.Errorf(nil, "response closed prematurely for electron [%s]", e.ID)
						}
					}
				}
			} else {
				alog.Error(nil, "response channel from send on conductor is nil")
			}
		}(ctx, response)
	} else {
		alog.Errorf(err, "error sending electron [%s]", e.ID)
	}

	return err
}

func (mc *MonteCarlo) estimate(ctx context.Context) <-chan []byte {
	result := make(chan []byte)

	go func(result chan<- []byte) {
		defer close(result)

		// Until the channel closes calculate how many tosses are greater than 0
		in, tosses, errors := mc.readtosses(ctx)
		pi := mc.calculate(float64(in), float64(tosses))

		res := struct {
			In     int     `json:"in"`
			Tosses int     `json:"tosses"`
			Errors int     `json:"errors"`
			PI     float64 `json:"pi"`
		}{
			in,
			tosses,
			errors,
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
			// TODO:
		}
	}(result)

	return result
}

func (mc *MonteCarlo) readtosses(ctx context.Context) (in, tosses, errors int) {

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
				alog.Printf("toss: %v\n", tosses)
				if v > 0 {
					in++
				} else if v < 0 {
					errors++
					tosses--
				}
			} else {
				return
			}
		}
	}

	return in, tosses, errors
}

func (mc *MonteCarlo) calculate(in, tosses float64) float64 {
	return (4 * in) / tosses
}

// double EstimatePI(long long inCircle, int numTosses) {
// 	return (4 * inCircle)/((double)numTosses);
// }
