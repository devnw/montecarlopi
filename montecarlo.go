package montecarlopi

import (
	"context"

	"github.com/benjivesterby/atomizer"
)

// MonteCarlo is the atom for estimating PI
type MonteCarlo struct {
	payload []byte
}

// ID test method
func (mc *MonteCarlo) ID() string {
	return "montecarlo"
}

// Process test method
func (mc *MonteCarlo) Process(ctx context.Context, electron atomizer.Electron, outbound chan<- atomizer.Electron) (<-chan []byte, <-chan error) {
	var results = make(chan []byte)
	var errors = make(chan error)

	go func() {

		// Step 1: parse my electron - How many tosses?

		// Step 2: Push the toss electrons onto the outbound channel

		// Step 3: Wait for the callbacks on the outbound electrons to finish processing

		// Step 4: Calculate PI estimation using the returned values from each toss "atom"

		// Step 5: Return the results on the results channel to the atomizer

	}()

	return results, errors
}

func (mc *MonteCarlo) EstimatePI(tosses <-chan int) {
	// Until the channel closes calculate how many tosses are greater than 0

	// Execute the calculation
	// EstimatePI takes in the toss results and estimates PI
	// double EstimatePI(long long inCircle, int numTosses) {
	// 	return (4 * inCircle)/((double)numTosses);
	// }
}
