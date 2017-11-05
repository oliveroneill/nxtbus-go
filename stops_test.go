package nxtbus

import (
	"testing"
)

func TestStopNameToID(t *testing.T) {
	// Not a unit test. This will test to ensure that with the
	// given data, the function will identify the correct stop.
	// This will need to be updated every time `stops.txt` is
	// updated
	var expected uint = 1364
	id, err := StopNameToID("Langdon Av opp Wanniassa Hills PS")
	if err != nil {
		t.Error("Failed", err)
	}
	if id != expected {
		t.Error("Expected", expected, "found", id)
	}
}
