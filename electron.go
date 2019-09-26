package montecarlopi

import (
	"time"

	"github.com/benjivesterby/atomizer"
)

type mcelectron struct {
	// ID is the electron ID
	ID string `json:id`

	// Tosses is the number of tosses to generate
	Tosses int `json:tosses`
}

func (e *mcelectron) Atom() (ID string) {
	return ""
}

func (e *mcelectron) Timeout() (timeout *time.Duration) {
	return nil
}

func (e *mcelectron) Validate() (valid bool) {
	return valid
}

// Callback allows the system to return results of a spawned electron to the caller after it's been bonded
// TODO: Need to set it up so that an atom can communicate with the original source by sending messages through a channel which takes electrons
//  When the electron is sent back to another node a channel is opened by the send method of the source and blocking will occur on reading from that channel
//  rather than relying on a callback with a waitgroup which is less reliable
func (e *mcelectron) Callback(properties atomizer.Properties) (err error) {
	return err
}
