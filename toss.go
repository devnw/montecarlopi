package montecarlopi

import (
	"context"
	"encoding/json"
	"math/rand"

	"github.com/benjivesterby/atomizer"
)

// Toss randomly tosses darts
type Toss struct {
	Value int `json:"value"`
}

// ID test method
func (t *Toss) ID() string {
	return "toss"
}

// Process test method
func (t *Toss) Process(ctx context.Context, electron atomizer.Electron, outbound chan<- atomizer.Electron) <-chan []byte {
	var results = make(chan []byte)

	go func() {
		defer close(results)
		// Step 1: Generate my Random X/Y Coordinates
		x := rand.Float64()
		y := rand.Float64()

		// Step 2: Return my Random X/Y Coordinates
		t.Value = t.dsquared(x, y)

		if b, err := json.Marshal(t); err == nil {

			// Push result to conductor
			select {
			case <-ctx.Done():
				return
			case results <- b:
			}
		}
	}()

	return results
}

func (t *Toss) dsquared(x, y float64) int {
	sq := (x * x) + (y * y)

	in := 0
	if sq <= 1 {
		in = 1
	}

	return in
}

// TODO: Implement for project
// Toss takes in the randomly generated x and y coordinates
// determines if they are less than 1 (or inside the circle)
// and returns a 1 if they are and 0 if they aren't so that
// the results can be added without needing explicit logic
// for incrementation
// int Toss(double x, double y) {
// 	double dsquared = 0;

// 	dsquared = (x*x) + (y*y);

// 	return dsquared <= 1 ? 1 : 0;
// }
