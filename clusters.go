package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/borystomala/copier"
)

// ClustersService is an interface for interacting with Clusters endpoints of CloudThing API
// https://tenant-name.cloudthing.io/api/v1/clusters
type ClustersService interface {
	GetById(string, ...interface{}) (*Cluster, error)
	GetByLink(string, ...interface{}) (*Cluster, error)
	ListByLink(string, ...interface{}) ([]Cluster, *ListParams, error)
	ListByApplication(string, ...interface{}) ([]Cluster, *ListParams, error)
	ListByDevice(string, ...interface{}) ([]Cluster, *ListParams, error)
	CreateByLink(string, *ClusterRequestCreate) (*Cluster, error)
	CreateByApplication(string, *ClusterRequestCreate) (*Cluster, error)
	UpdateById(string, *ClusterRequestUpdate) (*Cluster, error)
	UpdateByLink(string, *ClusterRequestUpdate) (*Cluster, error)
	Delete(*Cluster) error
	DeleteByLink(string) error
	DeleteById(string) error

	get(*ClusterResponse) (*Cluster, error)
	getCollection(*ClustersResponse) ([]Cluster, *ListParams, error)
}

// ClustersServiceOp handles communication with Clusters related methods of API
type ClustersServiceOp struct {
	client *Client
}

// Cluster is a struct representing CloudThing Cluster
type Cluster struct {
	// Standard field for all resources
	ModelBase
	Name        string
	Description string
	Custom      map[string]interface{}

	// Points to Tenant if expansion was requested, otherwise nil
	Tenant *Tenant
	// Points to Applications if expansion was requested, otherwise nil
	Application *Application
	// Points to Tenant if expansion was requested, otherwise nil
	Groups []Group
	// Points to Applications if expansion was requested, otherwise nil
	Devices []Device

	Memberships []ClusterMembership

	// Links to related resources
	tenant      string
	application string
	groups      string
	devices     string
	memberships string
	resources   string

	// service for communication, internal use only
	service *ClustersServiceOp
}

// ClusterResponse is a struct representing item response from API
type ClusterResponse struct {
	ModelBase
	Name        string                 `json:"name,omitempty"`
	Description string                 `json:"description,omitempty"`
	Custom      map[string]interface{} `json:"custom,omitempty"`

	Tenant      map[string]interface{} `json:"tenant,omitempty"`
	Application map[string]interface{} `json:"application,omitempty"`
	Groups      map[string]interface{} `json:"groups,omitempty"`
	Devices     map[string]interface{} `json:"devices,omitempty"`
	Memberships map[string]interface{} `json:"memberships,omitempty"`
}

// ClusterResponse is a struct representing collection response from API
type ClustersResponse struct {
	ListParams
	Items []ClusterResponse `json:"items"`
}

// ClusterResponse is a struct representing item create request for API
type ClusterRequestCreate struct {
	Name        string                 `json:"name,omitempty"`
	Description string                 `json:"description,omitempty"`
	Custom      map[string]interface{} `json:"custom,omitempty"`
}

// ClusterResponse is a struct representing item update request for API
type ClusterRequestUpdate struct {
	Name        string                 `json:"name,omitempty"`
	Description string                 `json:"description,omitempty"`
	Custom      map[string]interface{} `json:"custom,omitempty"`
}

// TenantLink returns indicator of Tenant expansion and link to tenant.
// If expansion for Tenant was requested and resource is available via pointer
// it returns true, otherwise false. Link (href) is always returned.
func (d *Cluster) TenantLink() (bool, string) {
	return (d.Tenant != nil), d.tenant
}

// ApplicationsLink returns indicator of Cluster expansion and link to dierctory.
// If expansion for Cluster was requested and resource is available via pointer
// it returns true, otherwise false. Link (href) is always returned.
func (d *Cluster) ApplicationLink() (bool, string) {
	return (d.Application != nil), d.application
}

// ApplicationsLink returns indicator of Cluster expansion and link to dierctory.
// If expansion for Cluster was requested and resource is available via pointer
// it returns true, otherwise false. Link (href) is always returned.
func (d *Cluster) GroupsLink() (bool, string) {
	return (d.Groups != nil), d.groups
}

// ApplicationsLink returns indicator of Cluster expansion and link to dierctory.
// If expansion for Cluster was requested and resource is available via pointer
// it returns true, otherwise false. Link (href) is always returned.
func (d *Cluster) DevicesLink() (bool, string) {
	return (d.Devices != nil), d.devices
}

// ApplicationsLink returns indicator of Cluster expansion and link to dierctory.
// If expansion for Cluster was requested and resource is available via pointer
// it returns true, otherwise false. Link (href) is always returned.
func (d *Cluster) MembershipsLink() (bool, string) {
	return (d.Memberships != nil), d.memberships
}

func (d *Cluster) ResourcesLink() string {
	return d.Href + "/resources"
}

func (d *Cluster) ResourcesDataLink() string {
	return d.ResourcesLink() + "/data"
}

func (d *Cluster) ResourcesEventsLink() string {
	return d.ResourcesLink() + "/events"
}

func (d *Cluster) ResourcesCommandsLink() string {
	return d.ResourcesLink() + "/commands"
}

func (d *Cluster) ResourcesDataKeyLink(key string) string {
	return d.ResourcesDataLink() + "/" + key
}

func (d *Cluster) ResourcesEventsKeyLink(key string) string {
	return d.ResourcesEventsLink() + "/" + key
}

func (d *Cluster) ResourcesCommandsKeyLink(key string) string {
	return d.ResourcesCommandsLink() + "/" + key
}

// Save is a helper method for updating apikey.
// It calls UpdateByLink() on service under the hood.
func (t *Cluster) Save() error {
	tmp := &ClusterRequestUpdate{}
	copier.Copy(tmp, t)
	ten, err := t.service.UpdateByLink(t.Href, tmp)
	if err != nil {
		return err
	}

	tmpTenant := t.Tenant
	tmpApplication := t.Application
	tmpDevices := t.Devices
	tmpGroups := t.Groups
	tmpMemberships := t.Memberships

	*t = *ten
	t.Tenant = tmpTenant
	t.Application = tmpApplication
	t.Devices = tmpDevices
	t.Groups = tmpGroups
	t.Memberships = tmpMemberships

	return nil
}

// Save is a helper method for deleting apikey.
// It calls Delete() on service under the hood.
func (t *Cluster) Delete() error {
	err := t.service.Delete(t)
	if err != nil {
		return err
	}

	return nil
}

// GetById retrieves apikey by its ID
func (s *ClustersServiceOp) GetById(id string, args ...interface{}) (*Cluster, error) {
	endpoint := "clusters/"
	endpoint = fmt.Sprintf("%s%s", endpoint, id)

	return s.GetByLink(endpoint, args...)
}

// GetById retrieves apikey by its full link
func (s *ClustersServiceOp) GetByLink(endpoint string, args ...interface{}) (*Cluster, error) {
	resp, err := s.client.request("GET", endpoint, nil, args...)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Status code: %d", resp.StatusCode)
	}
	obj := &ClusterResponse{}
	dec := json.NewDecoder(resp.Body)
	dec.Decode(obj)

	return s.get(obj)
}

// get is internal method for transforming ClusterResponse into Cluster
func (s *ClustersServiceOp) get(r *ClusterResponse) (*Cluster, error) {
	obj := &Cluster{}
	copier.Copy(obj, r)
	if v, ok := r.Tenant["href"]; ok {
		obj.tenant = v.(string)
	}
	if v, ok := r.Application["href"]; ok {
		obj.application = v.(string)
	}
	if v, ok := r.Groups["href"]; ok {
		obj.groups = v.(string)
	}
	if v, ok := r.Devices["href"]; ok {
		obj.devices = v.(string)
	}
	if v, ok := r.Memberships["href"]; ok {
		obj.memberships = v.(string)
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
	if len(r.Application) > 1 {
		bytes, err := json.Marshal(r.Application)
		if err != nil {
			return nil, err
		}
		ten := &ApplicationResponse{}
		json.Unmarshal(bytes, ten)
		t, err := s.client.Applications.get(ten)
		if err != nil {
			return nil, err
		}
		obj.Application = t
	}
	if len(r.Devices) > 1 {
		bytes, err := json.Marshal(r.Devices)
		if err != nil {
			return nil, err
		}
		ten := &DevicesResponse{}
		json.Unmarshal(bytes, ten)
		t, _, err := s.client.Devices.getCollection(ten)
		if err != nil {
			return nil, err
		}
		obj.Devices = t
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
	if len(r.Memberships) > 1 {
		bytes, err := json.Marshal(r.Memberships)
		if err != nil {
			return nil, err
		}
		ten := &ClusterMembershipsResponse{}
		json.Unmarshal(bytes, ten)
		t, _, err := s.client.ClusterMemberships.getCollection(ten)
		if err != nil {
			return nil, err
		}
		obj.Memberships = t
	}
	obj.service = s
	return obj, nil
}

// get is internal method for transforming ApplicationResponse into Application
func (s *ClustersServiceOp) getCollection(r *ClustersResponse) ([]Cluster, *ListParams, error) {
	dst := make([]Cluster, len(r.Items))

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

// GetById retrieves collection of clusters of current tenant
func (s *ClustersServiceOp) ListByApplication(id string, args ...interface{}) ([]Cluster, *ListParams, error) {
	endpoint := fmt.Sprintf("applications/%s/clusters", id)
	return s.ListByLink(endpoint, args...)
}

// GetById retrieves collection of clusters of current tenant
func (s *ClustersServiceOp) ListByDevice(id string, args ...interface{}) ([]Cluster, *ListParams, error) {
	endpoint := fmt.Sprintf("devices/%s/clusters", id)
	return s.ListByLink(endpoint, args...)
}

// GetById retrieves collection of clusters by link
func (s *ClustersServiceOp) ListByLink(endpoint string, args ...interface{}) ([]Cluster, *ListParams, error) {
	resp, err := s.client.request("GET", endpoint, nil, args...)
	if err != nil {
		return nil, nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, nil, fmt.Errorf("Status code: %d", resp.StatusCode)
	}
	obj := &ClustersResponse{}
	dec := json.NewDecoder(resp.Body)
	dec.Decode(obj)

	return s.getCollection(obj)
}

// GetById updates apikey with specified ID
func (s *ClustersServiceOp) UpdateById(id string, t *ClusterRequestUpdate) (*Cluster, error) {
	endpoint := fmt.Sprintf("clusters/%s", id)
	return s.UpdateByLink(endpoint, t)
}

// GetById updates apikey specified by link
func (s *ClustersServiceOp) UpdateByLink(endpoint string, t *ClusterRequestUpdate) (*Cluster, error) {
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
	obj := &ClusterResponse{}
	dec := json.NewDecoder(resp.Body)
	dec.Decode(obj)
	return s.get(obj)
}

func (s *ClustersServiceOp) CreateByApplication(id string, dir *ClusterRequestCreate) (*Cluster, error) {
	endpoint := fmt.Sprintf("applications/%s/clusters", id)
	return s.CreateByLink(endpoint, dir)
}

// Create creates new apikey within tenant
func (s *ClustersServiceOp) CreateByLink(endpoint string, dir *ClusterRequestCreate) (*Cluster, error) {
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
	obj := &ClusterResponse{}
	dec := json.NewDecoder(resp.Body)
	dec.Decode(obj)
	return s.get(obj)
}

// Delete removes apikey
func (s *ClustersServiceOp) Delete(t *Cluster) error {
	return s.DeleteByLink(t.Href)
}

// Delete removes apikey by ID
func (s *ClustersServiceOp) DeleteById(id string) error {
	endpoint := fmt.Sprintf("clusters/%s", id)
	return s.DeleteByLink(endpoint)
}

// Delete removes apikey by link
func (s *ClustersServiceOp) DeleteByLink(endpoint string) error {
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
