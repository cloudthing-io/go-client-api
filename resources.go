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
	GetDataByDevice(*Device, ...interface{}) ([]DataPoint, *ListParams, error)
	GetEventsByDevice(*Device, ...interface{}) ([]EventPoint, *ListParams, error)
	GetCommandsByDevice(*Device, ...interface{}) ([]CommandPoint, *ListParams, error)
	WriteDataForDevice(*Device, []DataPoint) ([]DataPoint, error)
	WriteEventsForDevice(*Device, []EventPoint) ([]EventPoint, error)
	WriteCommandsForDevice(*Device, []CommandPoint) ([]CommandPoint, error)
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

// GetDataByDevice requests from CloudThing device's data with set filters
func (s *ResourcesServiceOp) GetDataByDevice(device *Device, filters ...interface{}) ([]DataPoint, *ListParams, error) {
	endpoint := fmt.Sprintf("/api/v1/devices/%s/resources/data", device.GetId())

	obj := &DataResponse{}
	err := s.getResourcesByEndpoint(obj, endpoint, filters...)
	if err != nil {
		return nil, nil, err
	}
	return obj.Items, &obj.ListParams, nil
}

// GetEventsByDevice requests from CloudThing device's events with set filters
func (s *ResourcesServiceOp) GetEventsByDevice(device *Device, filters ...interface{}) ([]EventPoint, *ListParams, error) {
	endpoint := fmt.Sprintf("/api/v1/devices/%s/resources/events", device.GetId())

	obj := &EventsResponse{}
	err := s.getResourcesByEndpoint(obj, endpoint, filters...)
	if err != nil {
		return nil, nil, err
	}
	return obj.Items, &obj.ListParams, nil
}

// GetCommandsByDevice requests from CloudThing device's commands with set filters
func (s *ResourcesServiceOp) GetCommandsByDevice(device *Device, filters ...interface{}) ([]CommandPoint, *ListParams, error) {
	endpoint := fmt.Sprintf("/api/v1/devices/%s/resources/commands", device.GetId())

	obj := &CommandsResponse{}
	err := s.getResourcesByEndpoint(obj, endpoint, filters...)
	if err != nil {
		return nil, nil, err
	}
	return obj.Items, &obj.ListParams, nil
}

// WriteDataForDevice sends data to CloudThing device's for saving
func (s *ResourcesServiceOp) WriteDataForDevice(device *Device, points []DataPoint) ([]DataPoint, error) {
	endpoint := fmt.Sprintf("/api/v1/devices/%s/resources/data", device.GetId())

	obj := make([]DataPoint, 0)
	err := s.writeResourcesByEndpoint(obj, points, endpoint)
	if err != nil {
		return nil, err
	}
	return obj, nil
}

// WriteEventsForDevice sends events to CloudThing device's for saving
func (s *ResourcesServiceOp) WriteEventsForDevice(device *Device, points []EventPoint) ([]EventPoint, error) {
	endpoint := fmt.Sprintf("/api/v1/devices/%s/resources/events", device.GetId())

	obj := make([]EventPoint, 0)
	err := s.writeResourcesByEndpoint(obj, points, endpoint)
	if err != nil {
		return nil, err
	}
	return obj, nil
}

// WriteCommandsForDevice sends events to CloudThing device's for saving
func (s *ResourcesServiceOp) WriteCommandsForDevice(device *Device, points []CommandPoint) ([]CommandPoint, error) {
	endpoint := fmt.Sprintf("/api/v1/devices/%s/resources/commands", device.GetId())

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
