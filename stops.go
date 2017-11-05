package nxtbus

import (
	"encoding/csv"
	"errors"
	"io"
	"fmt"
	"strconv"
)

// StopsFile contains information for every stop
// Retrieved from: https://www.transport.act.gov.au/getting-around/bus-services/mobile-apps
const StopsFile = "stops.txt"

// Stops contains all Transport Canberra stops
var Stops = getStops()

// Stop contains information about a bus stop
type Stop struct {
	ID   uint
	Name string
}

// This will be run on import and will read through `stops.txt` to
// retrieve all stop information
func getStops() []Stop {
	// read through csv
	// Open is defined in `static.go` which is created by `staticfiles`
	csvFile, err := Open(StopsFile)
	if err != nil {
		fmt.Println(err)
		return []Stop{}
	}
	reader := csv.NewReader(csvFile)
	var stops []Stop
	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			continue
		}
		// convert strings to types as necessary
		id, err := strconv.ParseUint(line[0], 10, 0)
		// silently fail on parsing errors and continue
		if err != nil {
			continue
		}
		// populate slice
		stops = append(stops, Stop{
			ID:   uint(id),
			Name: line[2],
		})
	}
	return stops
}

// StopNameToID will find the first matching stop name and return the
// corresponding stop ID. This is useful when using Google Maps since it
// does not give you stop IDs
func StopNameToID(name string) (uint, error) {
	for _, stop := range Stops {
		if stop.Name == name {
			return stop.ID, nil
		}
	}
	return 0, errors.New("Couldn't find stop")
}
