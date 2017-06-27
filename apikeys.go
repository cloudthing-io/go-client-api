package api

import (
    "bytes"
    "encoding/json" 
    "fmt"   
    "net/http"
    _"strings"
    "github.com/borystomala/copier"
)

// ApikeysService is an interface for interacting with Apikeys endpoints of CloudThing API
// https://tenant-name.cloudthing.io/api/v1/apikeys
type ApikeysService interface {
    GetById(string, ...interface{}) (*Apikey, error)
    GetByLink(string, ...interface{}) (*Apikey, error)
    List(...interface{}) ([]Apikey, *ListParams, error)
    ListByLink(string, ...interface{}) ([]Apikey, *ListParams, error)
    Create(*ApikeyRequestCreate) (*Apikey, error)
    UpdateById(string, *ApikeyRequestUpdate) (*Apikey, error)
    UpdateByLink(string, *ApikeyRequestUpdate) (*Apikey, error)
    Delete(*Apikey) (error)
    DeleteByLink(string) (error)
    DeleteById(string) (error)

    get(*ApikeyResponse) (*Apikey, error)
    getCollection(*ApikeysResponse) ([]Apikey, *ListParams, error)
}

// ApikeysServiceOp handles communication with Apikeys related methods of API
type ApikeysServiceOp struct {
    client *Client
}

// Apikey is a struct representing CloudThing Apikey
type Apikey struct {
    // Standard field for all resources
    ModelBase
    // Apikey name
    Name            string
    // Description of apikey
    Description     string
    // Apikey status, may be ENABLED or DISABLED
    Status          string
    // Key
    Key             string 
    // Secret
    Secret          string
    // Field for tenant's custom data
    Custom          map[string]interface{}

    // Points to Tenant if expansion was requested, otherwise nil
    Tenant          *Tenant
    // Points to Applications if expansion was requested, otherwise nil
    Applications    []Application

    // Links to related resources
    tenant          string
    applications    string

    // service for communication, internal use only
    service         *ApikeysServiceOp
}

// ApikeyResponse is a struct representing item response from API
type ApikeyResponse struct {
    ModelBase
    Name            string                  `json:"name,omitempty"`
    Description     string                  `json:"description,omitempty"`
    Key             string                  `json:"key,omitempty"`
    Secret          string                  `json:"secret,omitempty"`
    Status          string                  `json:"status,omitempty"`
    Custom          map[string]interface{}  `json:"custom,omitempty"`

    Tenant          map[string]interface{}  `json:"tenant,omitempty"`
    Applications    map[string]interface{}  `json:"applications,omitempty"`
}

// ApikeyResponse is a struct representing collection response from API
type ApikeysResponse struct{
    ListParams
    Items           []ApikeyResponse   `json:"items"`
}

// ApikeyResponse is a struct representing item create request for API
type ApikeyRequestCreate struct {
    Name            string                  `json:"name,omitempty"`
    Description     string                  `json:"description,omitempty"`
    Status          string                  `json:"status,omitempty"`
    Custom          map[string]interface{}  `json:"custom,omitempty"`
}

// ApikeyResponse is a struct representing item update request for API
type ApikeyRequestUpdate struct {
    Name            string                  `json:"name,omitempty"`
    Description     string                  `json:"description,omitempty"`
    Status          string                  `json:"status,omitempty"`
    Custom          map[string]interface{}  `json:"custom,omitempty"`
}

// TenantLink returns indicator of Tenant expansion and link to tenant.
// If expansion for Tenant was requested and resource is available via pointer
// it returns true, otherwise false. Link (href) is always returned. 
func (d *Apikey) TenantLink() (bool, string) {
    return (d.Tenant != nil), d.tenant
}

// ApplicationsLink returns indicator of Directory expansion and link to dierctory.
// If expansion for Directory was requested and resource is available via pointer
// it returns true, otherwise false. Link (href) is always returned. 
func (d *Apikey) ApplicationsLink() (bool, string) {
    return (d.Applications != nil), d.applications
}

// Save is a helper method for updating apikey.
// It calls UpdateByLink() on service under the hood.
func (t *Apikey) Save() error {
    tmp := &ApikeyRequestUpdate{}
    copier.Copy(tmp, t)
    ten, err := t.service.UpdateByLink(t.Href, tmp)
    if err != nil {
        return err
    }

    tmpTenant := t.Tenant
    tmpApplications := t.Applications

    *t = *ten
    t.Tenant = tmpTenant
    t.Applications = tmpApplications

    return nil
}

// Save is a helper method for deleting apikey.
// It calls Delete() on service under the hood.
func (t *Apikey) Delete() error {
    err := t.service.Delete(t)
    if err != nil {
        return err
    }

    return nil
}

// GetById retrieves apikey by its ID
func (s *ApikeysServiceOp) GetById(id string, args ...interface{}) (*Apikey, error) {
    endpoint := "apikeys/"
    endpoint = fmt.Sprintf("%s%s", endpoint, id)

    return s.GetByLink(endpoint, args...)
}

// GetById retrieves apikey by its full link
func (s *ApikeysServiceOp) GetByLink(endpoint string, args ...interface{}) (*Apikey, error) {
    resp, err := s.client.request("GET", endpoint, nil, args...)
    if err != nil {
        return nil, err
    }

    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("Status code: %d", resp.StatusCode)
    }
    obj := &ApikeyResponse{}
    dec := json.NewDecoder(resp.Body)
    dec.Decode(obj)

    return s.get(obj)
}

// get is internal method for transforming ApikeyResponse into Apikey
func (s *ApikeysServiceOp) get(r *ApikeyResponse) (*Apikey, error) {
    obj := &Apikey{}
    copier.Copy(obj, r)
    if v, ok :=  r.Tenant["href"]; ok {
        obj.tenant = v.(string)
    }
    if v, ok :=  r.Applications["href"]; ok {
        obj.applications = v.(string)
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
   if len(r.Applications) > 1 {        
        bytes, err := json.Marshal(r.Applications)
        if err != nil {
            return nil, err
        }
        ten := &ApplicationsResponse{}
        json.Unmarshal(bytes, ten)
        t, _, err := s.client.Applications.getCollection(ten)
        if err != nil {
            return nil, err
        }
        obj.Applications = t
    }
    obj.service = s
    return obj, nil
}

// get is internal method for transforming ApplicationResponse into Application
func (s *ApikeysServiceOp) getCollection(r *ApikeysResponse) ([]Apikey, *ListParams, error) {
    dst := make([]Apikey, len(r.Items))

    for i, _ := range r.Items {
        t, err := s.get(&r.Items[i])
        if err == nil {
            dst[i] = *t
        }
    }

    lp := &ListParams {
        Href: r.Href,
        Prev: r.Prev,
        Next: r.Next,
        Limit: r.Limit,
        Size: r.Size,
        Page: r.Page,
    }
    return dst, lp, nil
}

// GetById retrieves collection of apikeys of current tenant
func (s *ApikeysServiceOp) List(args ...interface{}) ([]Apikey, *ListParams, error) {
    endpoint := fmt.Sprintf("tenants/%s/apikeys", s.client.tenantId)
    return s.ListByLink(endpoint, args...)
}

// GetById retrieves collection of apikeys by link
func (s *ApikeysServiceOp) ListByLink(endpoint string, args ...interface{}) ([]Apikey, *ListParams, error) {
    resp, err := s.client.request("GET", endpoint, nil, args...)
    if err != nil {
        return nil, nil, err
    }

    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, nil, fmt.Errorf("Status code: %d", resp.StatusCode)
    }
    obj := &ApikeysResponse{}
    dec := json.NewDecoder(resp.Body)
    dec.Decode(obj)

    return s.getCollection(obj)
}

// GetById updates apikey with specified ID
func (s *ApikeysServiceOp) UpdateById(id string, t *ApikeyRequestUpdate) (*Apikey, error) {
    endpoint := fmt.Sprintf("apikeys/%s", id)
    return s.UpdateByLink(endpoint, t)
}

// GetById updates apikey specified by link
func (s *ApikeysServiceOp) UpdateByLink(endpoint string, t *ApikeyRequestUpdate) (*Apikey, error) {
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
    obj := &ApikeyResponse{}
    dec := json.NewDecoder(resp.Body)
    dec.Decode(obj)
    return s.get(obj)
}

// Create creates new apikey within tenant
func (s *ApikeysServiceOp) Create(dir *ApikeyRequestCreate) (*Apikey, error) {
    endpoint := fmt.Sprintf("tenants/%s/apikeys", s.client.tenantId)

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
    obj := &ApikeyResponse{}
    dec := json.NewDecoder(resp.Body)
    dec.Decode(obj)
    return s.get(obj)
}

// Delete removes apikey
func (s *ApikeysServiceOp) Delete(t *Apikey) (error) {
    return s.DeleteByLink(t.Href)
}

// Delete removes apikey by ID
func (s *ApikeysServiceOp) DeleteById(id string) (error) {
    endpoint := fmt.Sprintf("apikeys/%s", id)
    return s.DeleteByLink(endpoint)
}

// Delete removes apikey by link
func (s *ApikeysServiceOp) DeleteByLink(endpoint string) (error) {
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
