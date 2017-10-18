package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/borystomala/copier"
)

// DevicesService is an interface for interacting with Devices endpoints of CloudThing API
// https://tenant-name.cloudthing.io/api/v1/devices
type DevicesService interface {
	GetById(string, ...interface{}) (*Device, error)
	GetByLink(string, ...interface{}) (*Device, error)
	ListByLink(string, ...interface{}) ([]Device, *ListParams, error)
	ListByCluster(string, ...interface{}) ([]Device, *ListParams, error)
	ListByApplication(string, ...interface{}) ([]Device, *ListParams, error)
	ListByGroup(string, ...interface{}) ([]Device, *ListParams, error)
	ListByProduct(string, ...interface{}) ([]Device, *ListParams, error)
	CreateByLink(string, *DeviceRequestCreate) (*Device, error)
	CreateByProduct(string, *DeviceRequestCreate) (*Device, error)
	UpdateById(string, *DeviceRequestUpdate) (*Device, error)
	UpdateByLink(string, *DeviceRequestUpdate) (*Device, error)
	Delete(*Device) error
	DeleteByLink(string) error
	DeleteById(string) error

	get(*DeviceResponse) (*Device, error)
	getCollection(*DevicesResponse) ([]Device, *ListParams, error)
}

// DevicesServiceOp handles communication with Devices related methods of API
type DevicesServiceOp struct {
	client *Client
}

// Device is a struct representing CloudThing Device
type Device struct {
	// Standard field for all resources
	ModelBase
	Token      string
	Activated  *bool
	Custom     map[string]interface{}
	Properties []DeviceProperty

	// Points to Tenant if expansion was requested, otherwise nil
	Tenant *Tenant
	// Points to Applications if expansion was requested, otherwise nil
	Product *Product
	// Points to Tenant if expansion was requested, otherwise nil
	Clusters           []Cluster
	ClusterMemberships []ClusterMembership
	// Points to Applications if expansion was requested, otherwise nil
	Groups           []Group
	GroupMemberships []GroupMembership

	// Links to related resources
	tenant             string
	product            string
	clusters           string
	groups             string
	clusterMemberships string
	groupMemberships   string

	// service for communication, internal use only
	service *DevicesServiceOp
}

type DeviceProperty struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

// DeviceResponse is a struct representing item response from API
type DeviceResponse struct {
	ModelBase
	Token      string                 `json:"token,omitempty"`
	Activated  *bool                  `json:"activated"`
	Custom     map[string]interface{} `json:"custom,omitempty"`
	Properties []DeviceProperty       `json:"properties,omitempty"`

	Tenant             map[string]interface{} `json:"tenant,omitempty"`
	Product            map[string]interface{} `json:"product,omitempty"`
	Clusters           map[string]interface{} `json:"clusters,omitempty"`
	Groups             map[string]interface{} `json:"groups,omitempty"`
	ClusterMemberships map[string]interface{} `json:"clusterMemberships,omitempty"`
	GroupMemberships   map[string]interface{} `json:"groupMemberships,omitempty"`
}

// DeviceResponse is a struct representing collection response from API
type DevicesResponse struct {
	ListParams
	Items []DeviceResponse `json:"items"`
}

// DeviceResponse is a struct representing item create request for API
type DeviceRequestCreate struct {
	Custom     map[string]interface{} `json:"custom,omitempty"`
	Properties []DeviceProperty       `json:"properties,omitempty"`
}

// DeviceResponse is a struct representing item update request for API
type DeviceRequestUpdate struct {
	Custom     map[string]interface{} `json:"custom,omitempty"`
	Properties []DeviceProperty       `json:"properties,omitempty"`
}

// TenantLink returns indicator of Tenant expansion and link to tenant.
// If expansion for Tenant was requested and resource is available via pointer
// it returns true, otherwise false. Link (href) is always returned.
func (d *Device) TenantLink() (bool, string) {
	return (d.Tenant != nil), d.tenant
}

// ApplicationsLink returns indicator of Device expansion and link to dierctory.
// If expansion for Device was requested and resource is available via pointer
// it returns true, otherwise false. Link (href) is always returned.
func (d *Device) ProductLink() (bool, string) {
	return (d.Product != nil), d.product
}

// ApplicationsLink returns indicator of Device expansion and link to dierctory.
// If expansion for Device was requested and resource is available via pointer
// it returns true, otherwise false. Link (href) is always returned.
func (d *Device) ClustersLink() (bool, string) {
	return (d.Clusters != nil), d.clusters
}

// ApplicationsLink returns indicator of Device expansion and link to dierctory.
// If expansion for Device was requested and resource is available via pointer
// it returns true, otherwise false. Link (href) is always returned.
func (d *Device) GroupsLink() (bool, string) {
	return (d.Groups != nil), d.groups
}

// ApplicationsLink returns indicator of Device expansion and link to dierctory.
// If expansion for Device was requested and resource is available via pointer
// it returns true, otherwise false. Link (href) is always returned.
func (d *Device) ClusterMembershipsLink() (bool, string) {
	return (d.ClusterMemberships != nil), d.clusterMemberships
}

// ApplicationsLink returns indicator of Device expansion and link to dierctory.
// If expansion for Device was requested and resource is available via pointer
// it returns true, otherwise false. Link (href) is always returned.
func (d *Device) GroupMembershipsLink() (bool, string) {
	return (d.GroupMemberships != nil), d.groupMemberships
}

func (d *Device) ResourcesLink() string {
	return d.Href + "/resources"
}

func (d *Device) ResourcesDataLink() string {
	return d.ResourcesLink() + "/data"
}

func (d *Device) ResourcesEventsLink() string {
	return d.ResourcesLink() + "/events"
}

func (d *Device) ResourcesCommandsLink() string {
	return d.ResourcesLink() + "/commands"
}

func (d *Device) ResourcesDataKeyLink(key string) string {
	return d.ResourcesDataLink() + "/" + key
}

func (d *Device) ResourcesEventsKeyLink(key string) string {
	return d.ResourcesEventsLink() + "/" + key
}

func (d *Device) ResourcesCommandsKeyLink(key string) string {
	return d.ResourcesCommandsLink() + "/" + key
}

// Save is a helper method for updating apikey.
// It calls UpdateByLink() on service under the hood.
func (t *Device) Save() error {
	tmp := &DeviceRequestUpdate{}
	copier.Copy(tmp, t)
	ten, err := t.service.UpdateByLink(t.Href, tmp)
	if err != nil {
		return err
	}

	tmpTenant := t.Tenant
	tmpProduct := t.Product
	tmpGroups := t.Groups
	tmpClusters := t.Clusters
	tmpGroupMemberships := t.GroupMemberships
	tmpClusterMemberships := t.ClusterMemberships

	*t = *ten
	t.Tenant = tmpTenant
	t.Product = tmpProduct
	t.Groups = tmpGroups
	t.Clusters = tmpClusters
	t.GroupMemberships = tmpGroupMemberships
	t.ClusterMemberships = tmpClusterMemberships

	return nil
}

// Save is a helper method for deleting apikey.
// It calls Delete() on service under the hood.
func (t *Device) Delete() error {
	err := t.service.Delete(t)
	if err != nil {
		return err
	}

	return nil
}

// GetById retrieves apikey by its ID
func (s *DevicesServiceOp) GetById(id string, args ...interface{}) (*Device, error) {
	endpoint := "devices/"
	endpoint = fmt.Sprintf("%s%s", endpoint, id)

	return s.GetByLink(endpoint, args...)
}

// GetById retrieves apikey by its full link
func (s *DevicesServiceOp) GetByLink(endpoint string, args ...interface{}) (*Device, error) {
	resp, err := s.client.request("GET", endpoint, nil, args...)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Status code: %d", resp.StatusCode)
	}
	obj := &DeviceResponse{}
	dec := json.NewDecoder(resp.Body)
	dec.Decode(obj)

	return s.get(obj)
}

// get is internal method for transforming DeviceResponse into Device
func (s *DevicesServiceOp) get(r *DeviceResponse) (*Device, error) {
	obj := &Device{}
	copier.Copy(obj, r)
	if v, ok := r.Tenant["href"]; ok {
		obj.tenant = v.(string)
	}
	if v, ok := r.Product["href"]; ok {
		obj.product = v.(string)
	}
	if v, ok := r.Clusters["href"]; ok {
		obj.clusters = v.(string)
	}
	if v, ok := r.Groups["href"]; ok {
		obj.groups = v.(string)
	}
	if v, ok := r.ClusterMemberships["href"]; ok {
		obj.clusterMemberships = v.(string)
	}
	if v, ok := r.GroupMemberships["href"]; ok {
		obj.groupMemberships = v.(string)
	}
	if len(r.Tenant) > 1 {
		bytes, err := json.Marshal(r.Tenant)
		if err != nil {
			return nil, err
		}
		ten := &TenantResponse{}
		json.Unmarshal(bytes, ten)
		t, err := s.client.Tenant.get(ten)
		if err != nil {
			return nil, err
		}
		obj.Tenant = t
	}
	if len(r.Product) > 1 {
		bytes, err := json.Marshal(r.Product)
		if err != nil {
			return nil, err
		}
		ten := &ProductResponse{}
		json.Unmarshal(bytes, ten)
		t, err := s.client.Products.get(ten)
		if err != nil {
			return nil, err
		}
		obj.Product = t
	}
	if len(r.Clusters) > 1 {
		bytes, err := json.Marshal(r.Clusters)
		if err != nil {
			return nil, err
		}
		ten := &ClustersResponse{}
		json.Unmarshal(bytes, ten)
		t, _, err := s.client.Clusters.getCollection(ten)
		if err != nil {
			return nil, err
		}
		obj.Clusters = t
	}
	if len(r.Groups) > 1 {
		bytes, err := json.Marshal(r.Groups)
		if err != nil {
			return nil, err
		}
		ten := &GroupsResponse{}
		json.Unmarshal(bytes, ten)
		t, _, err := s.client.Groups.getCollection(ten)
		if err != nil {
			return nil, err
		}
		obj.Groups = t
	}
	if len(r.ClusterMemberships) > 1 {
		bytes, err := json.Marshal(r.ClusterMemberships)
		if err != nil {
			return nil, err
		}
		ten := &ClusterMembershipsResponse{}
		json.Unmarshal(bytes, ten)
		t, _, err := s.client.ClusterMemberships.getCollection(ten)
		if err != nil {
			return nil, err
		}
		obj.ClusterMemberships = t
	}
	if len(r.GroupMemberships) > 1 {
		bytes, err := json.Marshal(r.GroupMemberships)
		if err != nil {
			return nil, err
		}
		ten := &GroupMembershipsResponse{}
		json.Unmarshal(bytes, ten)
		t, _, err := s.client.GroupMemberships.getCollection(ten)
		if err != nil {
			return nil, err
		}
		obj.GroupMemberships = t
	}
	obj.service = s
	return obj, nil
}

// get is internal method for transforming ApplicationResponse into Application
func (s *DevicesServiceOp) getCollection(r *DevicesResponse) ([]Device, *ListParams, error) {
	dst := make([]Device, len(r.Items))

	for i, _ := range r.Items {
		t, err := s.get(&r.Items[i])
		if err == nil {
			dst[i] = *t
		}
	}

	lp := &ListParams{
		Href:  r.Href,
		Prev:  r.Prev,
		Next:  r.Next,
		Limit: r.Limit,
		Size:  r.Size,
		Page:  r.Page,
	}
	return dst, lp, nil
}

// GetById retrieves collection of devices of current tenant
func (s *DevicesServiceOp) ListByCluster(id string, args ...interface{}) ([]Device, *ListParams, error) {
	endpoint := fmt.Sprintf("clusters/%s/devices", id)
	return s.ListByLink(endpoint, args...)
}

// GetById retrieves collection of devices of current tenant
func (s *DevicesServiceOp) ListByApplication(id string, args ...interface{}) ([]Device, *ListParams, error) {
	endpoint := fmt.Sprintf("applications/%s/devices", id)
	return s.ListByLink(endpoint, args...)
}

// GetById retrieves collection of devices of current tenant
func (s *DevicesServiceOp) ListByGroup(id string, args ...interface{}) ([]Device, *ListParams, error) {
	endpoint := fmt.Sprintf("groups/%s/devices", id)
	return s.ListByLink(endpoint, args...)
}

// GetById retrieves collection of devices of current tenant
func (s *DevicesServiceOp) ListByProduct(id string, args ...interface{}) ([]Device, *ListParams, error) {
	endpoint := fmt.Sprintf("products/%s/devices", id)
	return s.ListByLink(endpoint, args...)
}

// GetById retrieves collection of devices by link
func (s *DevicesServiceOp) ListByLink(endpoint string, args ...interface{}) ([]Device, *ListParams, error) {
	resp, err := s.client.request("GET", endpoint, nil, args...)
	if err != nil {
		return nil, nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, nil, fmt.Errorf("Status code: %d", resp.StatusCode)
	}
	obj := &DevicesResponse{}
	dec := json.NewDecoder(resp.Body)
	dec.Decode(obj)

	return s.getCollection(obj)
}

// GetById updates apikey with specified ID
func (s *DevicesServiceOp) UpdateById(id string, t *DeviceRequestUpdate) (*Device, error) {
	endpoint := fmt.Sprintf("devices/%s", id)
	return s.UpdateByLink(endpoint, t)
}

// GetById updates apikey specified by link
func (s *DevicesServiceOp) UpdateByLink(endpoint string, t *DeviceRequestUpdate) (*Device, error) {
	enc, err := json.Marshal(t)
	if err != nil {
		return nil, err
	}

	buf := bytes.NewBuffer(enc)

	resp, err := s.client.request("POST", endpoint, buf)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Status code: %d", resp.StatusCode)
	}
	obj := &DeviceResponse{}
	dec := json.NewDecoder(resp.Body)
	dec.Decode(obj)
	return s.get(obj)
}

func (s *DevicesServiceOp) CreateByProduct(id string, dir *DeviceRequestCreate) (*Device, error) {
	endpoint := fmt.Sprintf("products/%s/devices", id)
	return s.CreateByLink(endpoint, dir)
}

// Create creates new apikey within tenant
func (s *DevicesServiceOp) CreateByLink(endpoint string, dir *DeviceRequestCreate) (*Device, error) {
	enc, err := json.Marshal(dir)
	if err != nil {
		return nil, err
	}

	buf := bytes.NewBuffer(enc)

	resp, err := s.client.request("POST", endpoint, buf)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("Status code: %d", resp.StatusCode)
	}
	obj := &DeviceResponse{}
	dec := json.NewDecoder(resp.Body)
	dec.Decode(obj)
	return s.get(obj)
}

// Delete removes apikey
func (s *DevicesServiceOp) Delete(t *Device) error {
	return s.DeleteByLink(t.Href)
}

// Delete removes apikey by ID
func (s *DevicesServiceOp) DeleteById(id string) error {
	endpoint := fmt.Sprintf("devices/%s", id)
	return s.DeleteByLink(endpoint)
}

// Delete removes apikey by link
func (s *DevicesServiceOp) DeleteByLink(endpoint string) error {
	resp, err := s.client.request("DELETE", endpoint, nil)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("Status code: %d", resp.StatusCode)
	}
	return nil
}
