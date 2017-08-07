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
	GetDataByClusterID(string, ...interface{}) ([]DataPoint, *ListParams, error)
	GetEventsByClusterID(string, ...interface{}) ([]EventPoint, *ListParams, error)
	GetCommandsByClusterID(string, ...interface{}) ([]CommandPoint, *ListParams, error)
	WriteDataForClusterID(string, []DataPoint) ([]DataPoint, error)
	WriteEventsForClusterID(string, []EventPoint) ([]EventPoint, error)
	WriteCommandsForClusterID(string, []CommandPoint) ([]CommandPoint, error)
	GetDataByLink(string, ...interface{}) ([]DataPoint, *ListParams, error)
	GetEventsByLink(string, ...interface{}) ([]EventPoint, *ListParams, error)
	GetCommandsByLink(string, ...interface{}) ([]CommandPoint, *ListParams, error)
	WriteDataForLink(string, []DataPoint) ([]DataPoint, error)
	WriteEventsForLink(string, []EventPoint) ([]EventPoint, error)
	WriteCommandsForLink(string, []CommandPoint) ([]CommandPoint, error)
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

	return s.GetDataByLink(endpoint, filters...)
}

// GetEventsByDeviceID requests from CloudThing device's events with set filters
func (s *ResourcesServiceOp) GetEventsByDeviceID(deviceID string, filters ...interface{}) ([]EventPoint, *ListParams, error) {
	endpoint := fmt.Sprintf("/api/v1/devices/%s/resources/events", deviceID)

	return s.GetEventsByLink(endpoint, filters...)
}

// GetCommandsByDeviceID requests from CloudThing device's commands with set filters
func (s *ResourcesServiceOp) GetCommandsByDeviceID(deviceID string, filters ...interface{}) ([]CommandPoint, *ListParams, error) {
	endpoint := fmt.Sprintf("/api/v1/devices/%s/resources/commands", deviceID)

	return s.GetCommandsByLink(endpoint, filters...)
}

// WriteDataForDeviceID sends data to CloudThing device's for saving
func (s *ResourcesServiceOp) WriteDataForDeviceID(deviceID string, points []DataPoint) ([]DataPoint, error) {
	endpoint := fmt.Sprintf("/api/v1/devices/%s/resources/data", deviceID)

	return s.WriteDataForLink(endpoint, points)
}

// WriteEventsForDeviceID sends events to CloudThing device's for saving
func (s *ResourcesServiceOp) WriteEventsForDeviceID(deviceID string, points []EventPoint) ([]EventPoint, error) {
	endpoint := fmt.Sprintf("/api/v1/devices/%s/resources/events", deviceID)

	return s.WriteEventsForLink(endpoint, points)
}

// WriteCommandsForDeviceID sends events to CloudThing device's for saving
func (s *ResourcesServiceOp) WriteCommandsForDeviceID(deviceID string, points []CommandPoint) ([]CommandPoint, error) {
	endpoint := fmt.Sprintf("/api/v1/devices/%s/resources/commands", deviceID)

	return s.WriteCommandsForLink(endpoint, points)
}

// GetDataByClusterID requests from CloudThing cluster's data with set filters
func (s *ResourcesServiceOp) GetDataByClusterID(clusterID string, filters ...interface{}) ([]DataPoint, *ListParams, error) {
	endpoint := fmt.Sprintf("/api/v1/clusters/%s/resources/data", clusterID)

	return s.GetDataByLink(endpoint, filters...)
}

// GetEventsByClusterID requests from CloudThing cluster's events with set filters
func (s *ResourcesServiceOp) GetEventsByClusterID(clusterID string, filters ...interface{}) ([]EventPoint, *ListParams, error) {
	endpoint := fmt.Sprintf("/api/v1/clusters/%s/resources/events", clusterID)

	return s.GetEventsByLink(endpoint, filters...)
}

// GetCommandsByClusterID requests from CloudThing cluster's commands with set filters
func (s *ResourcesServiceOp) GetCommandsByClusterID(clusterID string, filters ...interface{}) ([]CommandPoint, *ListParams, error) {
	endpoint := fmt.Sprintf("/api/v1/clusters/%s/resources/commands", clusterID)

	return s.GetCommandsByLink(endpoint, filters...)
}

// WriteDataForClusterID sends data to CloudThing cluster's for saving
func (s *ResourcesServiceOp) WriteDataForClusterID(clusterID string, points []DataPoint) ([]DataPoint, error) {
	endpoint := fmt.Sprintf("/api/v1/clusters/%s/resources/data", clusterID)

	return s.WriteDataForLink(endpoint, points)
}

// WriteEventsForClusterID sends events to CloudThing cluster's for saving
func (s *ResourcesServiceOp) WriteEventsForClusterID(clusterID string, points []EventPoint) ([]EventPoint, error) {
	endpoint := fmt.Sprintf("/api/v1/clusters/%s/resources/events", clusterID)

	return s.WriteEventsForLink(endpoint, points)
}

// WriteCommandsForClusterID sends events to CloudThing cluster's for saving
func (s *ResourcesServiceOp) WriteCommandsForClusterID(clusterID string, points []CommandPoint) ([]CommandPoint, error) {
	endpoint := fmt.Sprintf("/api/v1/clusters/%s/resources/commands", clusterID)

	return s.WriteCommandsForLink(endpoint, points)
}

// GetDataByLink requests from CloudThing data with set filters
func (s *ResourcesServiceOp) GetDataByLink(link string, filters ...interface{}) ([]DataPoint, *ListParams, error) {
	obj := &DataResponse{}
	err := s.getResourcesByEndpoint(obj, link, filters...)
	if err != nil {
		return nil, nil, err
	}
	return obj.Items, &obj.ListParams, nil
}

// GetEventsByLink requests from CloudThing events with set filters
func (s *ResourcesServiceOp) GetEventsByLink(link string, filters ...interface{}) ([]EventPoint, *ListParams, error) {
	obj := &EventsResponse{}
	err := s.getResourcesByEndpoint(obj, link, filters...)
	if err != nil {
		return nil, nil, err
	}
	return obj.Items, &obj.ListParams, nil
}

// GetCommandsByLink requests from CloudThing commands with set filters
func (s *ResourcesServiceOp) GetCommandsByLink(link string, filters ...interface{}) ([]CommandPoint, *ListParams, error) {
	obj := &CommandsResponse{}
	err := s.getResourcesByEndpoint(obj, link, filters...)
	if err != nil {
		return nil, nil, err
	}
	return obj.Items, &obj.ListParams, nil
}

// WriteDataForLink sends data to CloudThing for saving
func (s *ResourcesServiceOp) WriteDataForLink(link string, points []DataPoint) ([]DataPoint, error) {
	obj := make([]DataPoint, 0)
	err := s.writeResourcesByEndpoint(obj, points, link)
	if err != nil {
		return nil, err
	}
	return obj, nil
}

// WriteEventsForLink sends events to CloudThing for saving
func (s *ResourcesServiceOp) WriteEventsForLink(link string, points []EventPoint) ([]EventPoint, error) {
	obj := make([]EventPoint, 0)
	err := s.writeResourcesByEndpoint(obj, points, link)
	if err != nil {
		return nil, err
	}
	return obj, nil
}

// WriteCommandsForLink sends events to CloudThing for saving
func (s *ResourcesServiceOp) WriteCommandsForLink(link string, points []CommandPoint) ([]CommandPoint, error) {
	obj := make([]CommandPoint, 0)
	err := s.writeResourcesByEndpoint(obj, points, link)
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
		return fmt.Errorf("CloudThing API returned non-OK status code: %d", resp.StatusCode)
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
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("CloudThing API returned non-OK status code: %d", resp.StatusCode)
	}
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&responseBody)
	if err != nil {
		return err
	}
	return nil
}
