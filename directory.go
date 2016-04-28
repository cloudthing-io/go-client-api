package api

import (
    "bytes"
    "encoding/json" 
    "fmt"   
    "net/http"
    _"strings"
    "github.com/borystomala/copier"
)

// DirectoriesService is an interface for interacting with Directories endpoints of CloudThing API
// https://tenant-name.cloudthing.io/api/v1/directories
type DirectoriesService interface {
    GetById(string, ...interface{}) (*Directory, error)
    GetByLink(string, ...interface{}) (*Directory, error)
    List(...interface{}) ([]Directory, *ListParams, error)
    ListByLink(string, ...interface{}) ([]Directory, *ListParams, error)
    Create(*DirectoryRequestCreate) (*Directory, error)
    UpdateById(string, *DirectoryRequestUpdate) (*Directory, error)
    UpdateByLink(string, *DirectoryRequestUpdate) (*Directory, error)
    Delete(*Directory) (error)
    DeleteByLink(string) (error)
    DeleteById(string) (error)

    get(*DirectoryResponse) (*Directory, error)
    getCollection(*DirectoriesResponse) ([]Directory, *ListParams, error)
}

// DirectoriesServiceOp handles communication with Directories related methods of API
type DirectoriesServiceOp struct {
    client *Client
}

// Directory is a struct representing CloudThing Directory
type Directory struct {
    // Standard field for all resources
    ModelBase
    // Directory name
    Name            string
    // Description of apikey
    Description     string
    // Indicates whether directory is official or not
    Official        *bool 
    // Field for tenant's custom data
    Custom          map[string]interface{}

    // Points to Tenant if expansion was requested, otherwise nil
    Tenant          *Tenant
    // Points to Applications if expansion was requested, otherwise nil
    Applications    []Application
    // Points to Tenant if expansion was requested, otherwise nil
    Users           *User
    // Points to Applications if expansion was requested, otherwise nil
    Usergroups      []Usergroup

    // Links to related resources
    tenant          string
    applications    string
    users           string
    usergroups      string

    // service for communication, internal use only
    service         *DirectoriesServiceOp
}

// DirectoryResponse is a struct representing item response from API
type DirectoryResponse struct {
    ModelBase
    Name            string                  `json:"name,omitempty"`
    Description     string                  `json:"description,omitempty"`
    Official        *bool                   `json:"official,omitempty"`
    Custom          map[string]interface{}  `json:"custom,omitempty"`

    Tenant          map[string]interface{}  `json:"tenant,omitempty"`
    Applications    map[string]interface{}  `json:"applications,omitempty"`
    Users           map[string]interface{}  `json:"users,omitempty"`
    Usergroups      map[string]interface{}  `json:"usergroups,omitempty"`
}

// DirectoryResponse is a struct representing collection response from API
type DirectoriesResponse struct{
    ListParams
    Items           []DirectoryResponse     `json:"items"`
}

// DirectoryResponse is a struct representing item create request for API
type DirectoryRequestCreate struct {
    Name            string                  `json:"name,omitempty"`
    Description     string                  `json:"description,omitempty"`
    Custom          map[string]interface{}  `json:"custom,omitempty"`
}

// DirectoryResponse is a struct representing item update request for API
type DirectoryRequestUpdate struct {
    Name            string                  `json:"name,omitempty"`
    Description     string                  `json:"description,omitempty"`
    Custom          map[string]interface{}  `json:"custom,omitempty"`
}

// TenantLink returns indicator of Tenant expansion and link to tenant.
// If expansion for Tenant was requested and resource is available via pointer
// it returns true, otherwise false. Link (href) is always returned. 
func (d *Directory) TenantLink() (bool, string) {
    return (d.Tenant != nil), d.tenant
}

// ApplicationsLink returns indicator of Directory expansion and link to dierctory.
// If expansion for Directory was requested and resource is available via pointer
// it returns true, otherwise false. Link (href) is always returned. 
func (d *Directory) ApplicationsLink() (bool, string) {
    return (d.Applications != nil), d.applications
}

// ApplicationsLink returns indicator of Directory expansion and link to dierctory.
// If expansion for Directory was requested and resource is available via pointer
// it returns true, otherwise false. Link (href) is always returned. 
func (d *Directory) UsersLink() (bool, string) {
    return (d.Users != nil), d.users
}

// ApplicationsLink returns indicator of Directory expansion and link to dierctory.
// If expansion for Directory was requested and resource is available via pointer
// it returns true, otherwise false. Link (href) is always returned. 
func (d *Directory) UsergroupsLink() (bool, string) {
    return (d.Usergroups != nil), d.usergroups
}

// Save is a helper method for updating apikey.
// It calls UpdateByLink() on service under the hood.
func (t *Directory) Save() error {
    tmp := &DirectoryRequestUpdate{}
    copier.Copy(tmp, t)
    ten, err := t.service.UpdateByLink(t.Href, tmp)
    if err != nil {
        return err
    }

    tmpTenant := t.Tenant
    tmpApplications := t.Applications
    tmpUsers := t.Users
    tmpUsergroups := t.Usergroups

    *t = *ten
    t.Tenant = tmpTenant
    t.Applications = tmpApplications
    t.Users = tmpUsers
    t.Usergroups = tmpUsergroups

    return nil
}

// Save is a helper method for deleting apikey.
// It calls Delete() on service under the hood.
func (t *Directory) Delete() error {
    err := t.service.Delete(t)
    if err != nil {
        return err
    }

    return nil
}

// GetById retrieves apikey by its ID
func (s *DirectoriesServiceOp) GetById(id string, args ...interface{}) (*Directory, error) {
    endpoint := "directories/"
    endpoint = fmt.Sprintf("%s%s", endpoint, id)

    return s.GetByLink(endpoint, args...)
}

// GetById retrieves apikey by its full link
func (s *DirectoriesServiceOp) GetByLink(endpoint string, args ...interface{}) (*Directory, error) {
    resp, err := s.client.request("GET", endpoint, nil, args...)
    if err != nil {
        return nil, err
    }

    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("Status code: %d", resp.StatusCode)
    }
    obj := &DirectoryResponse{}
    dec := json.NewDecoder(resp.Body)
    dec.Decode(obj)

    return s.get(obj)
}

// get is internal method for transforming DirectoryResponse into Directory
func (s *DirectoriesServiceOp) get(r *DirectoryResponse) (*Directory, error) {
    obj := &Directory{}
    copier.Copy(obj, r)
    if v, ok :=  r.Tenant["href"]; ok {
        obj.tenant = v.(string)
    }
    if v, ok :=  r.Applications["href"]; ok {
        obj.applications = v.(string)
    }
    if v, ok :=  r.Users["href"]; ok {
        obj.users = v.(string)
    }
    if v, ok :=  r.Usergroups["href"]; ok {
        obj.usergroups = v.(string)
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
func (s *DirectoriesServiceOp) getCollection(r *DirectoriesResponse) ([]Directory, *ListParams, error) {
    dst := make([]Directory, len(r.Items))

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

// GetById retrieves collection of directories of current tenant
func (s *DirectoriesServiceOp) List(args ...interface{}) ([]Directory, *ListParams, error) {
    endpoint := fmt.Sprintf("tenants/%s/directories", s.client.tenantId)
    return s.ListByLink(endpoint, args...)
}

// GetById retrieves collection of directories by link
func (s *DirectoriesServiceOp) ListByLink(endpoint string, args ...interface{}) ([]Directory, *ListParams, error) {
    resp, err := s.client.request("GET", endpoint, nil, args...)
    if err != nil {
        return nil, nil, err
    }

    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, nil, fmt.Errorf("Status code: %d", resp.StatusCode)
    }
    obj := &DirectoriesResponse{}
    dec := json.NewDecoder(resp.Body)
    dec.Decode(obj)

    return s.getCollection(obj)
}

// GetById updates apikey with specified ID
func (s *DirectoriesServiceOp) UpdateById(id string, t *DirectoryRequestUpdate) (*Directory, error) {
    endpoint := fmt.Sprintf("directories/%s", id)
    return s.UpdateByLink(endpoint, t)
}

// GetById updates apikey specified by link
func (s *DirectoriesServiceOp) UpdateByLink(endpoint string, t *DirectoryRequestUpdate) (*Directory, error) {
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
    obj := &DirectoryResponse{}
    dec := json.NewDecoder(resp.Body)
    dec.Decode(obj)
    return s.get(obj)
}

// Create creates new apikey within tenant
func (s *DirectoriesServiceOp) Create(dir *DirectoryRequestCreate) (*Directory, error) {
    endpoint := fmt.Sprintf("tenants/%s/directories", s.client.tenantId)

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
    obj := &DirectoryResponse{}
    dec := json.NewDecoder(resp.Body)
    dec.Decode(obj)
    return s.get(obj)
}

// Delete removes apikey
func (s *DirectoriesServiceOp) Delete(t *Directory) (error) {
    return s.DeleteByLink(t.Href)
}

// Delete removes apikey by ID
func (s *DirectoriesServiceOp) DeleteById(id string) (error) {
    endpoint := fmt.Sprintf("directories/%s", id)
    return s.DeleteByLink(endpoint)
}

// Delete removes apikey by link
func (s *DirectoriesServiceOp) DeleteByLink(endpoint string) (error) {
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
