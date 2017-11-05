package nxtbus

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// StopMonitoringRequestURL is the NXTBUS API url
const StopMonitoringRequestURL = "http://siri.nxtbus.act.gov.au:11000/%s/%s/service.xml"
// StopMonitoringServiceType is the indicator for monitoring stops
const StopMonitoringServiceType = "sm"

// SiriVersion is the version of Siri being used for the requests
const SiriVersion = "2.0"
// Xmlns is the Siri namespace
const Xmlns = "http://www.siri.org.uk/siri"
// Xsi is the schema instance
const Xsi = "http://www.w3.org/2001/XMLSchema-instance"
//Xsd is the schema
const Xsd = "http://www.w3.org/2001/XMLSchema"
// PreviewInterval is the time window that we want to monitor stops for
const PreviewInterval = 90

// SiriSchema is the request document for monitoring stops
type SiriSchema struct {
	ServiceRequest *ServiceRequest `xml:"ServiceRequest"`
	Xmlns          string          `xml:"xmlns,attr"`
	XmlnsXsi       string          `xml:"xmlns xsi,attr"`
	XmlnsXsd       string          `xml:"xmlns xsd,attr"`
	Version        string          `xml:"version,attr"`
	XMLName        xml.Name        `xml:"Siri"`
}

// ServiceRequest is the request information including API key
type ServiceRequest struct {
	RequestTimestamp      string
	StopMonitoringRequest *StopMonitoringRequest
	RequestorRef          string
}

// StopMonitoringRequest is the information regarding the stop and timings
type StopMonitoringRequest struct {
	Version          string `xml: "-version"`
	RequestTimestamp string
	PreviewInterval  string
	MonitoringRef    uint
}

// StopMonitoringResponse is the response document
type StopMonitoringResponse struct {
	StopMonitoringDelivery *StopMonitoringDelivery `xml:"ServiceDelivery>StopMonitoringDelivery"`
	ResponseTimestamp      string
}

// StopMonitoringDelivery contains a list of routes available at this stop
type StopMonitoringDelivery struct {
	MonitoredStopVisits []MonitoredStopVisit `xml:"MonitoredStopVisit"`
}

// MonitoredStopVisit contains information regarding each route through this
// stop
type MonitoredStopVisit struct {
	ExpectedDepartureTime string `xml:"MonitoredVehicleJourney>MonitoredCall>ExpectedDepartureTime"`
	ExpectedArrivalTime   string `xml:"MonitoredVehicleJourney>MonitoredCall>ExpectedArrivalTime"`
	StopPointRef          string `xml:"MonitoredVehicleJourney>MonitoredCall>StopPointRef"`
	AimedArrivalTime      string `xml:"MonitoredVehicleJourney>MonitoredCall>AimedArrivalTime"`
	AimedDepartureTime    string `xml:"MonitoredVehicleJourney>MonitoredCall>AimedDepartureTime"`
	DeparturePlatformName string `xml:"MonitoredVehicleJourney>MonitoredCall>DeparturePlatformName"`
}

// createStopMonitoringRequestBody will create a request body for the specified
// API key and Stop ID, this document can then be used to query NXTBUS
func createStopMonitoringRequestBody(apiKey string, stopID uint) ([]byte, error) {
	now := time.Now().Format("2006-01-02T15:04:05.000000")
	tmp := SiriSchema{
		Xmlns:    Xmlns,
		XmlnsXsi: Xsi,
		XmlnsXsd: Xsd,
		Version:  SiriVersion,
		ServiceRequest: &ServiceRequest{
			RequestTimestamp: now,
			StopMonitoringRequest: &StopMonitoringRequest{
				Version:          SiriVersion,
				RequestTimestamp: now,
				PreviewInterval:  fmt.Sprintf("PT%dM", PreviewInterval),
				MonitoringRef:    stopID,
			},
			RequestorRef: apiKey,
		},
	}
	byteXML, err := xml.Marshal(tmp)
	if err != nil {
		return nil, err
	}
	// add standard XML header
	parsed := xml.Header + string(byteXML)
	return []byte(parsed), nil
}

// createStopMonitoringRequestURL will create the URL to query, using the
// API key
func createStopMonitoringRequestURL(apiKey string) string {
	return fmt.Sprintf(StopMonitoringRequestURL, apiKey, StopMonitoringServiceType)
}

// parseStopMonitoringResponseBody will parse the HTTP response into the
// structs defined above
func parseStopMonitoringResponseBody(body []byte) (*StopMonitoringResponse, error) {
	parsed := &StopMonitoringResponse{}
	err := xml.Unmarshal(body, parsed)
	return parsed, err
}

// MakeStopMonitoringRequest will make the request and return the response
func MakeStopMonitoringRequest(apiKey string, stopID uint) (*StopMonitoringResponse, error) {
	url := createStopMonitoringRequestURL(apiKey)
	// generate request document
	body, err := createStopMonitoringRequestBody(apiKey, stopID)
	if err != nil {
		return nil, err
	}
	// make request
	resp, err := http.Post(url, "application/xml", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	// read body into byte array
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return parseStopMonitoringResponseBody(content)
}
