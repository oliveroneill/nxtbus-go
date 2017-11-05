package nxtbus

import (
	"encoding/xml"
	"fmt"
	"testing"
)

func TestCreateStopMonitoringRequestBody(t *testing.T) {
	apiKey := "TEST_KEY_123"
	var stopID uint = 3412
	result, err := createStopMonitoringRequestBody(apiKey, stopID)
	if err != nil {
		t.Error("Request body failed:", err)
	}
	// ensure that the API key and Stop ID are added to the
	// correct fields
	parsed := &SiriSchema{}
	// parse the XML back into the struct
	err = xml.Unmarshal(result, parsed)
	if err != nil {
		t.Error("Unmarshal failed:", err)
	}
	parsedKey := parsed.ServiceRequest.RequestorRef
	if parsedKey != apiKey {
		t.Error("Expected", apiKey, "found", parsedKey)
	}
	parsedStopID := parsed.ServiceRequest.StopMonitoringRequest.MonitoringRef
	if parsedStopID != stopID {
		t.Error("Expected", apiKey, "found", parsedStopID)
	}
}

func TestParseStopMonitoringResponseBody(t *testing.T) {
	// check some fields are correctly marshalled into struct
	aimedArrival := "2013-12-30T17:38:00.000+01:00"
	expectedArrival := "2013-12-30T17:38:00.000+01:00"
	expectedDeparture := "2013-12-30T17:48:56.000+01:00"
	data := `
	<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
	<Siri version="2.0" xmlns:ns2="http://www.ifopt.org.uk/acsb"
		xmlns="http://www.sir i.org.uk/siri"
		xmlns:ns4="http://datex2.eu/schema/2_0RC1/2_0"
		xmlns:ns3="http://ww w.ifopt.org.uk/ifopt">
	<ServiceDelivery>
		<ResponseTimestamp>2013-12-30T17:44:03.000+01:00</ResponseTimestamp>
		<StopMonitoringDelivery version="2.0">
			<MonitoredStopVisit>
				<MonitoredVehicleJourney>
					<MonitoredCall>
						<VehicleAtStop>true</VehicleAtStop>
						<AimedArrivalTime>%s</AimedArrivalTime>
						<ExpectedArrivalTime>%s</ExpectedArrivalTime>
						<AimedDepartureTime>2013-12-30T17:18:56.000+01:00</AimedDepartureTime>
						<ExpectedDepartureTime>2013-12-30T17:28:56.000+01:00</ExpectedDepartureTime>
						<DeparturePlatformName>Stg 2</DeparturePlatformName>
				 </MonitoredCall>
			  </MonitoredVehicleJourney>
			</MonitoredStopVisit>
			<MonitoredStopVisit>
				<MonitoredVehicleJourney>
					<MonitoredCall>
						<VehicleAtStop>true</VehicleAtStop>
						<AimedArrivalTime>2013-12-31T17:18:56.000+01:00</AimedArrivalTime>
						<ExpectedArrivalTime>2013-11-30T17:18:56.000+01:00</ExpectedArrivalTime>
						<AimedDepartureTime>2013-12-30T17:18:56.000+01:00</AimedDepartureTime>
						<ExpectedDepartureTime>%s</ExpectedDepartureTime>
						<DeparturePlatformName>Stg 2</DeparturePlatformName>
				 </MonitoredCall>
			  </MonitoredVehicleJourney>
			</MonitoredStopVisit>
		 </StopMonitoringDelivery>
	   </ServiceDelivery>
	</Siri>`
	formattedData := fmt.Sprintf(data, aimedArrival, expectedArrival, expectedDeparture)
	parsed, err := parseStopMonitoringResponseBody([]byte(formattedData))
	if err != nil {
		t.Error("Parse failed:", err)
	}
	parsedAimedArrival := parsed.StopMonitoringDelivery.MonitoredStopVisits[0].AimedArrivalTime
	if parsedAimedArrival != aimedArrival {
		t.Error("Expected", aimedArrival, "found", parsedAimedArrival)
	}
	parsedExpectedArrival := parsed.StopMonitoringDelivery.MonitoredStopVisits[0].ExpectedArrivalTime
	if parsedExpectedArrival != expectedArrival {
		t.Error("Expected", expectedArrival, "found", parsedExpectedArrival)
	}
	// ensure it can read multiple stop visits
	parsedExpectedDeparture := parsed.StopMonitoringDelivery.MonitoredStopVisits[1].ExpectedDepartureTime
	if parsedExpectedDeparture != expectedDeparture {
		t.Error("Expected", expectedDeparture, "found", parsedExpectedDeparture)
	}
}
