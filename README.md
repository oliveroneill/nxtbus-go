# nxtbus-go
Golang implementation of Transport Canberra's NXTBUS API. This currently only
has support for StopMonitoringRequest/Response.

## Installation
```
go get github.com/oliveroneill/nxtbus-go
```

## Usage
Example usage
```go
package main

import (
    "fmt"
    "github.com/oliveroneill/nxtbus-go"
)

func main() {
    resp, err := nxtbus.MakeStopMonitoringRequest("<API_KEY>", <STOP_ID>)
    if err != nil {
        fmt.Println("ERROR:", err)
    }
    // get the expected arrival time
    fmt.Println(resp.StopMonitoringDelivery.MonitoredStopVisits[0].ExpectedArrivalTime)
}
```

## TODO

  * Parse dates automatically when parsing XML response
  * Add all fields available for StopMonitoringResponse
  * Add other endpoints of NXTBUS API
