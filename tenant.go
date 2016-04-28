package api

import (
	"fmt"
	"encoding/json"
	"bytes"
	"net/http"
    "github.com/jinzhu/copier"
)

// TenantService is an interafce for interacting with Tenants endpoints of CloudThing API
// https://tenant-name.cloudthing.io/api/v1/tenants

type TenantService interface {
	Get() (*Tenant, error)
    Update(*Tenant) (*Tenant, error)

    get(*TenantResponse) (*Tenant, error)
}

// TenantServiceOp handles communication with Tenant related methods of API
type TenantServiceOp struct {
	client *Client
}

type TenantResponse struct {
    ModelBase
    ShortName       string          `json:"shortName,omitempty"`
    Name            string          `json:"name,omitempty"`
    Custom          map[string]interface{}     `json:"custom,omitempty"`

    Directories     map[string]interface{}           `json:"directories,omitempty"`
    Applications    map[string]interface{}`json:"applications,omitempty"`
    Products        map[string]interface{}`json:"products,omitempty"`    
}

// Tenant represents an Organization (tenant) within CloudThing
type Tenant struct {
    ModelBase
    ShortName       string          `json:"shortName,omitempty"`
    Name       		string          `json:"name,omitempty"`
    Custom          interface{}     `json:"custom,omitempty"`

    Directories     []Directory           `json:"directories,omitempty"`
    Applications    []Application           `json:"applications,omitempty"`
    Products        []Product           `json:"products,omitempty"`    

    // service for communication, internal use only
    service 		*TenantServiceOp `json:"-"` 
}

// Directories retrieves directories of current tenant
func (t *Tenant) GetDirectories() ([]Directory, *ListParams, error) {
    return t.service.client.Directories.List(nil)
}

// Directories retrieves directories of current tenant
func (t *Tenant) GetApplications() ([]Application, *ListParams, error) {
    return t.service.client.Applications.List(nil)
}

// Products retrieves directories of current tenant
func (t *Tenant) GetProducts() ([]Product, *ListParams, error) {
    return t.service.client.Products.List(nil)
}

// Save updates tenant by calling Update() on service under the hood
func (t *Tenant) Save() error {
	tmp := *t

    ten, err := t.service.Update(&tmp)
    if err != nil {
    	return err
    }

    *t = *ten
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
func (s *TenantServiceOp) Update(t *Tenant) (*Tenant, error) {
	endpoint := t.Href

    t.CreatedAt = nil
    t.UpdatedAt = nil
    t.Href = ""
    t.ShortName = ""
    t.Custom = nil

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