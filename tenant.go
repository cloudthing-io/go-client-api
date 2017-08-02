package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/borystomala/copier"
)

// TenantService is an interafce for interacting with Tenants endpoints of CloudThing API
// https://tenant-name.cloudthing.io/api/v1/tenants

type TenantService interface {
	Get() (*Tenant, error)
	UpdateByLink(string, *TenantRequestUpdate) (*Tenant, error)

	get(*TenantResponse) (*Tenant, error)
}

// TenantServiceOp handles communication with Tenant related methods of API
type TenantServiceOp struct {
	client *Client
}

// Tenant represents an Organization (tenant) within CloudThing
type Tenant struct {
	ModelBase
	ShortName string
	Name      string
	Custom    map[string]interface{}

	Directories  []Directory
	Applications []Application
	Products     []Product

	// Links to related resources
	directories  string
	applications string
	products     string

	// service for communication, internal use only
	service *TenantServiceOp
}

type TenantResponse struct {
	ModelBase
	ShortName string                 `json:"shortName,omitempty"`
	Name      string                 `json:"name,omitempty"`
	Custom    map[string]interface{} `json:"custom,omitempty"`

	Directories  map[string]interface{} `json:"directories,omitempty"`
	Applications map[string]interface{} `json:"applications,omitempty"`
	Products     map[string]interface{} `json:"products,omitempty"`
}

type TenantRequestUpdate struct {
	Name   string                 `json:"name,omitempty"`
	Custom map[string]interface{} `json:"custom,omitempty"`
}

type TenantRequestCreate struct {
	Name   string                 `json:"name,omitempty"`
	Custom map[string]interface{} `json:"custom,omitempty"`
}

// Directories retrieves directories of current tenant
func (t *Tenant) DirectoriesLink() (bool, string) {
	return (t.Directories != nil), t.directories
}

// Directories retrieves directories of current tenant
func (t *Tenant) ApplicationsLink() (bool, string) {
	return (t.Applications != nil), t.applications
}

// Products retrieves directories of current tenant
func (t *Tenant) ProductsLink() (bool, string) {
	return (t.Products != nil), t.products
}

// Save is a helper method for updating tenant.
// It calls UpdateByLink() on service under the hood.
func (t *Tenant) Save() error {
	tmp := &TenantRequestUpdate{}
	copier.Copy(tmp, t)
	ten, err := t.service.UpdateByLink(t.Href, tmp)
	if err != nil {
		return err
	}

	tmpApplications := t.Applications
	tmpDirectories := t.Directories
	tmpProducts := t.Products

	*t = *ten
	t.Applications = tmpApplications
	t.Directories = tmpDirectories
	t.Products = tmpProducts

	return nil
}

// Get retrieves current tenant
func (s *TenantServiceOp) Get() (*Tenant, error) {
	endpoint := "tenants/"
	if s.client.tenantId == "" {
		endpoint = fmt.Sprintf("%s%s", endpoint, "current")
	} else {
		endpoint = fmt.Sprintf("%s%s", endpoint, s.client.tenantId)
	}

	resp, err := s.client.request("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusForbidden || resp.StatusCode == http.StatusUnauthorized {
		// this is probably due to redirect
		endpoint = resp.Request.URL.String()
		resp, err = s.client.request("GET", endpoint, nil)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("Status code: %d", resp.StatusCode)
		}
	} else if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Status code: %d", resp.StatusCode)
	}
	tenant := &TenantResponse{}
	dec := json.NewDecoder(resp.Body)
	dec.Decode(tenant)
	return s.get(tenant)
}

func (s *TenantServiceOp) get(t *TenantResponse) (*Tenant, error) {
	obj := &Tenant{}
	copier.Copy(obj, t)
	obj.service = s
	return obj, nil
}

// Update updates tenant
func (s *TenantServiceOp) UpdateByLink(endpoint string, t *TenantRequestUpdate) (*Tenant, error) {
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
	tenant := &Tenant{}
	dec := json.NewDecoder(resp.Body)
	dec.Decode(tenant)
	tenant.service = s
	return tenant, nil
}
