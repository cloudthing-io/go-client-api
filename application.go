package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	_ "strings"

	"github.com/borystomala/copier"
)

// ApplicationsService is an interface for interacting with Applications endpoints of CloudThing API
// https://tenant-name.cloudthing.io/api/v1/applications
type ApplicationsService interface {
	GetById(string, ...interface{}) (*Application, error)
	GetByLink(string, ...interface{}) (*Application, error)
	List(...interface{}) ([]Application, *ListParams, error)
	ListByLink(string, ...interface{}) ([]Application, *ListParams, error)
	Create(*ApplicationRequestCreate) (*Application, error)
	UpdateById(string, *ApplicationRequestUpdate) (*Application, error)
	UpdateByLink(string, *ApplicationRequestUpdate) (*Application, error)
	Delete(*Application) error
	DeleteByLink(string) error
	DeleteById(string) error

	get(*ApplicationResponse) (*Application, error)
	getCollection(*ApplicationsResponse) ([]Application, *ListParams, error)
}

// ApplicationsServiceOp handles communication with Applications related methods of API
type ApplicationsServiceOp struct {
	client *Client
}

// Application is a struct representing CloudThing Application
type Application struct {
	// Standard field for all resources
	ModelBase
	// Application name
	Name string
	// Indicates whether application is official or not
	Official *bool
	// Description of application
	Description string
	// Application status, may be ENABLED or DISABLED
	Status string
	// Field for tenant's custom data
	Custom map[string]interface{}

	// Points to Tenant if expansion was requested, otherwise nil
	Tenant *Tenant
	// Points to Directory if expansion was requested, otherwise nil
	Directory *Directory
	// Points to Devices if expansion was requested, otherwise nil
	Devices []Device
	// Points to Clusters if expansion was requested, otherwise nil
	Clusters []Cluster

	// Links to related resources
	tenant    string
	directory string
	devices   string
	clusters  string

	// service for communication, internal use only
	service *ApplicationsServiceOp
}

// ApplicationResponse is a struct representing item response from API
type ApplicationResponse struct {
	ModelBase
	Name        string                 `json:"name,omitempty"`
	Official    *bool                  `json:"official,omitempty"`
	Description string                 `json:"description,omitempty"`
	Status      string                 `json:"status,omitempty"`
	Custom      map[string]interface{} `json:"custom,omitempty"`

	Tenant    map[string]interface{} `json:"tenant,omitempty"`
	Directory map[string]interface{} `json:"directory,omitempty"`
	Devices   map[string]interface{} `json:"devices,omitempty"`
	Clusters  map[string]interface{} `json:"clusters,omitempty"`
}

// ApplicationResponse is a struct representing collection response from API
type ApplicationsResponse struct {
	ListParams
	Items []ApplicationResponse `json:"items"`
}

// ApplicationResponse is a struct representing item create request for API
type ApplicationRequestCreate struct {
	Name        string                 `json:"name,omitempty"`
	Description string                 `json:"description,omitempty"`
	Status      string                 `json:"status,omitempty"`
	Custom      map[string]interface{} `json:"custom,omitempty"`
	Directory   *Link                  `json:"directory,omitempty"`
}

// ApplicationResponse is a struct representing item update request for API
type ApplicationRequestUpdate struct {
	Name        string                 `json:"name,omitempty"`
	Description string                 `json:"description,omitempty"`
	Status      string                 `json:"status,omitempty"`
	Custom      map[string]interface{} `json:"custom,omitempty"`
	//Directory       *Link                   `json:"directory,omitempty"`
}

// TenantLink returns indicator of Tenant expansion and link to tenant.
// If expansion for Tenant was requested and resource is available via pointer
// it returns true, otherwise false. Link (href) is always returned.
func (d *Application) TenantLink() (bool, string) {
	return (d.Tenant != nil), d.tenant
}

// DirectoryLink returns indicator of Directory expansion and link to dierctory.
// If expansion for Directory was requested and resource is available via pointer
// it returns true, otherwise false. Link (href) is always returned.
func (d *Application) DirectoryLink() (bool, string) {
	return (d.Directory != nil), d.directory
}

// DevicesLink returns indicator of Devices expansion and link to list of devices.
// If expansion for Devices was requested and resource is available via pointer
// it returns true, otherwise false. Link (href) is always returned.
func (d *Application) DevicesLink() (bool, string) {
	return (d.Devices != nil), d.devices
}

// ClustersLink returns indicator of Clusters expansion and link to list of clusters.
// If expansion for Clusters was requested and resource is available via pointer
// it returns true, otherwise false. Link (href) is always returned.
func (d *Application) ClustersLink() (bool, string) {
	return (d.Clusters != nil), d.clusters
}

// Save is a helper method for updating application.
// It calls UpdateByLink() on service under the hood.
func (t *Application) Save() error {
	tmp := &ApplicationRequestUpdate{}
	copier.Copy(tmp, t)
	ten, err := t.service.UpdateByLink(t.Href, tmp)
	if err != nil {
		return err
	}

	tmpTenant := t.Tenant
	tmpDirectory := t.Directory
	tmpDevices := t.Devices
	tmpClusters := t.Clusters

	*t = *ten
	t.Tenant = tmpTenant
	t.Directory = tmpDirectory
	t.Devices = tmpDevices
	t.Clusters = tmpClusters

	return nil
}

// Save is a helper method for deleting application.
// It calls Delete() on service under the hood.
func (t *Application) Delete() error {
	err := t.service.Delete(t)
	if err != nil {
		return err
	}

	return nil
}

// GetById retrieves application by its ID
func (s *ApplicationsServiceOp) GetById(id string, args ...interface{}) (*Application, error) {
	endpoint := "applications/"
	endpoint = fmt.Sprintf("%s%s", endpoint, id)

	return s.GetByLink(endpoint, args...)
}

// GetById retrieves application by its full link
func (s *ApplicationsServiceOp) GetByLink(endpoint string, args ...interface{}) (*Application, error) {
	resp, err := s.client.request("GET", endpoint, nil, args...)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Status code: %d", resp.StatusCode)
	}
	obj := &ApplicationResponse{}
	dec := json.NewDecoder(resp.Body)
	dec.Decode(obj)

	return s.get(obj)
}

// get is internal method for transforming ApplicationResponse into Application
func (s *ApplicationsServiceOp) get(r *ApplicationResponse) (*Application, error) {
	obj := &Application{}
	copier.Copy(obj, r)
	if v, ok := r.Tenant["href"]; ok {
		obj.tenant = v.(string)
	}
	if v, ok := r.Directory["href"]; ok {
		obj.directory = v.(string)
	}
	if v, ok := r.Devices["href"]; ok {
		obj.devices = v.(string)
	}
	if v, ok := r.Clusters["href"]; ok {
		obj.clusters = v.(string)
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
	/*   if len(r.Directory) > 1 {
	     bytes, err := json.Marshal(r.Directory)
	     if err != nil {
	         return nil, err
	     }
	     ten := &DirectoryResponse{}
	     json.Unmarshal(bytes, ten)
	     t, err := s.client.Directories.get(ten)
	     if err != nil {
	         return nil, err
	     }
	     obj.Directory = t
	 }*/
	obj.service = s
	return obj, nil
}

// get is internal method for transforming ApplicationResponse into Application
func (s *ApplicationsServiceOp) getCollection(r *ApplicationsResponse) ([]Application, *ListParams, error) {
	dst := make([]Application, len(r.Items))

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

// GetById retrieves collection of applications of current tenant
func (s *ApplicationsServiceOp) List(args ...interface{}) ([]Application, *ListParams, error) {
	endpoint := fmt.Sprintf("tenants/%s/applications", s.client.tenantId)
	return s.ListByLink(endpoint, args...)
}

// GetById retrieves collection of applications by link
func (s *ApplicationsServiceOp) ListByLink(endpoint string, args ...interface{}) ([]Application, *ListParams, error) {
	resp, err := s.client.request("GET", endpoint, nil, args...)
	if err != nil {
		return nil, nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, nil, fmt.Errorf("Status code: %d", resp.StatusCode)
	}
	obj := &ApplicationsResponse{}
	dec := json.NewDecoder(resp.Body)
	dec.Decode(obj)

	return s.getCollection(obj)
}

// GetById updates application with specified ID
func (s *ApplicationsServiceOp) UpdateById(id string, t *ApplicationRequestUpdate) (*Application, error) {
	endpoint := fmt.Sprintf("applications/%s", id)
	return s.UpdateByLink(endpoint, t)
}

// GetById updates application specified by link
func (s *ApplicationsServiceOp) UpdateByLink(endpoint string, t *ApplicationRequestUpdate) (*Application, error) {
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
	obj := &ApplicationResponse{}
	dec := json.NewDecoder(resp.Body)
	dec.Decode(obj)
	return s.get(obj)
}

// Create creates new application within tenant
func (s *ApplicationsServiceOp) Create(dir *ApplicationRequestCreate) (*Application, error) {
	endpoint := fmt.Sprintf("tenants/%s/applications", s.client.tenantId)

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
	obj := &ApplicationResponse{}
	dec := json.NewDecoder(resp.Body)
	dec.Decode(obj)
	return s.get(obj)
}

// Delete removes application
func (s *ApplicationsServiceOp) Delete(t *Application) error {
	return s.DeleteByLink(t.Href)
}

// Delete removes application by ID
func (s *ApplicationsServiceOp) DeleteById(id string) error {
	endpoint := fmt.Sprintf("applications/%s", id)
	return s.DeleteByLink(endpoint)
}

// Delete removes application by link
func (s *ApplicationsServiceOp) DeleteByLink(endpoint string) error {
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
