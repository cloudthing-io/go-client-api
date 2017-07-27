package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"gitlab.com/cloudthing/structures"
)

// ResourcesService is an interface for interacting with Resources endpoints of CloudThing API
// https://tenant-name.cloudthing.io/api/v1/model/id/resources
type ResourcesService interface {
	GetDataByDeviceID(string, ...interface{}) ([]DataPoint, *ListParams, error)
	GetEventsByDeviceID(string, ...interface{}) ([]EventPoint, *ListParams, error)
	GetCommandsByDeviceID(string, ...interface{}) ([]CommandPoint, *ListParams, error)
	WriteDataForDeviceID(string, []DataPoint) ([]DataPoint, error)
	WriteEventsForDeviceID(string, []EventPoint) ([]EventPoint, error)
	WriteCommandsForDeviceID(string, []CommandPoint) ([]CommandPoint, error)
}

// ResourcesServiceOp handles communication with Resources related methods of API
type ResourcesServiceOp struct {
	client *Client
}

// DataPoint is a single key-value pair with associated time. Used in describing data series
type DataPoint struct {
	Time  string          `json:"time"`
	Value interface{}     `json:"value"`
	Key   string          `json:"key,omitempty"`
	Geo   *structures.Geo `json:"geo,omitempty"`
}

// EventPoint is a single timeseries point with associated time and added payload. Used in events series
type EventPoint struct {
	Time    string          `json:"time"`
	Payload interface{}     `json:"payload"`
	Key     string          `json:"key,omitempty"`
	Geo     *structures.Geo `json:"geo,omitempty"`
}

// CommandPoint is a single timeseries point with associated time and added payload. Used in commands
type CommandPoint EventPoint

// TimeParams allows to add 'start' and 'end' query parameters
type TimeParams struct {
	Start *time.Time
	End   *time.Time
}

// DataResponse is a response from CloudThing when calling for data resource
type DataResponse struct {
	ListParams
	Items []DataPoint `json:"items"`
}

// EventsResponse is a response from CloudThing when calling for events resource
type EventsResponse struct {
	ListParams
	Items []EventPoint `json:"items"`
}

// CommandsResponse is a response from CloudThing when calling for commands resource
type CommandsResponse struct {
	ListParams
	Items []CommandPoint `json:"items"`
}

func (t *TimeParams) String() string {
	return fmt.Sprintf("start=%s&end=%s", t.Start.UTC().Format(time.RFC3339), t.End.UTC().Format(time.RFC3339))
}

// GetDataByDeviceID requests from CloudThing device's data with set filters
func (s *ResourcesServiceOp) GetDataByDeviceID(deviceID string, filters ...interface{}) ([]DataPoint, *ListParams, error) {
	endpoint := fmt.Sprintf("/api/v1/devices/%s/resources/data", deviceID)

	obj := &DataResponse{}
	err := s.getResourcesByEndpoint(obj, endpoint, filters...)
	if err != nil {
		return nil, nil, err
	}
	return obj.Items, &obj.ListParams, nil
}

// GetEventsByDeviceID requests from CloudThing device's events with set filters
func (s *ResourcesServiceOp) GetEventsByDeviceID(deviceID string, filters ...interface{}) ([]EventPoint, *ListParams, error) {
	endpoint := fmt.Sprintf("/api/v1/devices/%s/resources/events", deviceID)

	obj := &EventsResponse{}
	err := s.getResourcesByEndpoint(obj, endpoint, filters...)
	if err != nil {
		return nil, nil, err
	}
	return obj.Items, &obj.ListParams, nil
}

// GetCommandsByDeviceID requests from CloudThing device's commands with set filters
func (s *ResourcesServiceOp) GetCommandsByDeviceID(deviceID string, filters ...interface{}) ([]CommandPoint, *ListParams, error) {
	endpoint := fmt.Sprintf("/api/v1/devices/%s/resources/commands", deviceID)

	obj := &CommandsResponse{}
	err := s.getResourcesByEndpoint(obj, endpoint, filters...)
	if err != nil {
		return nil, nil, err
	}
	return obj.Items, &obj.ListParams, nil
}

// WriteDataForDeviceID sends data to CloudThing device's for saving
func (s *ResourcesServiceOp) WriteDataForDeviceID(deviceID string, points []DataPoint) ([]DataPoint, error) {
	endpoint := fmt.Sprintf("/api/v1/devices/%s/resources/data", deviceID)

	obj := make([]DataPoint, 0)
	err := s.writeResourcesByEndpoint(obj, points, endpoint)
	if err != nil {
		return nil, err
	}
	return obj, nil
}

// WriteEventsForDeviceID sends events to CloudThing device's for saving
func (s *ResourcesServiceOp) WriteEventsForDeviceID(deviceID string, points []EventPoint) ([]EventPoint, error) {
	endpoint := fmt.Sprintf("/api/v1/devices/%s/resources/events", deviceID)

	obj := make([]EventPoint, 0)
	err := s.writeResourcesByEndpoint(obj, points, endpoint)
	if err != nil {
		return nil, err
	}
	return obj, nil
}

// WriteCommandsForDeviceID sends events to CloudThing device's for saving
func (s *ResourcesServiceOp) WriteCommandsForDeviceID(deviceID string, points []CommandPoint) ([]CommandPoint, error) {
	endpoint := fmt.Sprintf("/api/v1/devices/%s/resources/commands", deviceID)

	obj := make([]CommandPoint, 0)
	err := s.writeResourcesByEndpoint(obj, points, endpoint)
	if err != nil {
		return nil, err
	}
	return obj, nil
}

func (s *ResourcesServiceOp) getResourcesByEndpoint(responseBody interface{}, endpoint string, filters ...interface{}) error {
	resp, err := s.client.request("GET", endpoint, nil, filters...)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Status code: %d", resp.StatusCode)
	}

	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(responseBody)
	if err != nil {
		return err
	}
	return nil
}

func (s *ResourcesServiceOp) writeResourcesByEndpoint(responseBody interface{}, data interface{}, endpoint string) error {
	enc, err := json.Marshal(data)
	if err != nil {
		return err
	}
	buf := bytes.NewBuffer(enc)
	resp, err := s.client.request("POST", endpoint, buf)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&responseBody)
	if err != nil {
		return err
	}
	return nil
}
