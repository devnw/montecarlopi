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
	tossed int
	tosses chan int
}

// ID test method
func (mc *MonteCarlo) ID() string {
	return "montecarlo"
}

// Process test method
func (mc *MonteCarlo) Process(ctx context.Context, electron atomizer.Electron, outbound chan<- atomizer.Electron) <-chan []byte {
	var results = make(chan []byte)

	go func(results chan<- []byte) {

		mc.tosses = make(chan int)

		var e = &mcelectron{}
		if err := json.Unmarshal(electron.Payload(), e); err == nil {

			if e.Tosses > 0 {
				mc.tossed = e.Tosses

				go mc.estimate(ctx, results)

				for i := 0; i < e.Tosses; i++ {
					select {
					case <-ctx.Done():
						return
					case outbound <- mc.toss(ctx):
					}
				}

			} else {
				// TODO:
			}
		} else {
			alog.Errorf(err, "error un-marshalling %s", string(electron.Payload()))
		}

		// Step 1: parse my electron - How many tosses?

		// Step 2: Push the toss electrons onto the outbound channel

		// Step 3: Wait for the callbacks on the outbound electrons to finish processing

		// Step 4: Calculate PI estimation using the returned values from each toss "atom"

		// Step 5: Return the results on the results channel to the atomizer

	}(results)

	return results
}

func (mc *MonteCarlo) toss(ctx context.Context) atomizer.Electron {
	response := make(chan *atomizer.Properties)

	e := &atomizer.ElectronBase{
		ElectronID: uuid.New().String(),
		AtomID:     "toss",
		Resp:       response,
	}

	go func(ctx context.Context, response chan *atomizer.Properties) {
		defer close(response)

		select {
		case <-ctx.Done():
		case r, ok := <-response:
			if ok && r.Error == nil {

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
				mc.tosses <- -1
			}
		}
	}(ctx, response)

	return e
}

func (mc *MonteCarlo) estimate(ctx context.Context, result chan<- []byte) {
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
}

func (mc *MonteCarlo) readtosses(ctx context.Context) (in, tosses, errors int) {

	count := 0

	// Execute the calculation
	// EstimatePI takes in the toss results and estimates PI
	for {
		count++
		select {
		case <-ctx.Done():
			return
		case v, ok := <-mc.tosses:
			if ok && count < mc.tossed {
				alog.Printf("toss: %v\n", tosses)
				if v > 0 {
					in++
					tosses++
				} else if v == 0 {
					tosses++
				} else {
					errors++
				}
			} else {
				return
			}
		}
	}
}

func (mc *MonteCarlo) calculate(in, tosses float64) float64 {
	return (4 * in) / tosses
}

// double EstimatePI(long long inCircle, int numTosses) {
// 	return (4 * inCircle)/((double)numTosses);
// }
