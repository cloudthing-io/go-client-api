package api

import (
    "bytes"
    "encoding/json" 
    "fmt"   
    "net/http"
    "github.com/borystomala/copier"
)

// UsergroupsService is an interface for interacting with Usergroups endpoints of CloudThing API
// https://tenant-name.cloudthing.io/api/v1/usergroups
type UsergroupsService interface {
    GetById(string, ...interface{}) (*Usergroup, error)
    GetByLink(string, ...interface{}) (*Usergroup, error)
    ListByLink(string, ...interface{}) ([]Usergroup, *ListParams, error)
    ListByDirectory(string, ...interface{}) ([]Usergroup, *ListParams, error)
    CreateByLink(string, *UsergroupRequestCreate) (*Usergroup, error)
    CreateByDirectory(string, *UsergroupRequestCreate) (*Usergroup, error)
    UpdateById(string, *UsergroupRequestUpdate) (*Usergroup, error)
    UpdateByLink(string, *UsergroupRequestUpdate) (*Usergroup, error)
    Delete(*Usergroup) (error)
    DeleteByLink(string) (error)
    DeleteById(string) (error)

    get(*UsergroupResponse) (*Usergroup, error)
    getCollection(*UsergroupsResponse) ([]Usergroup, *ListParams, error)
}

// UsergroupsServiceOp handles communication with Usergroups related methods of API
type UsergroupsServiceOp struct {
    client *Client
}

// Usergroup is a struct representing CloudThing Usergroup
type Usergroup struct {
    // Standard field for all resources
    ModelBase
    Name            string
    Custom          map[string]interface{}

    // Points to Tenant if expansion was requested, otherwise nil
    Tenant          *Tenant
    // Points to Directory if expansion was requested, otherwise nil
    Directory       *Directory
    // Points to Users if expansion was requested, otherwise nil
    Users           []User
    // Points to Memberships if expansion was requested, otherwise nil
    Memberships     []Membership

    // Links to related resources
    tenant          string
    directory       string
    users           string
    memberships     string

    // service for communication, internal use only
    service         *UsergroupsServiceOp
}

// UsergroupResponse is a struct representing item response from API
type UsergroupResponse struct {
    ModelBase
    Name            string                  `json:"name,omitempty"`
    Custom          map[string]interface{} `json:"custom,omitempty"`

    Tenant          map[string]interface{}  `json:"tenant,omitempty"`
    Directory       map[string]interface{}  `json:"directory,omitempty"`
    Users           map[string]interface{}  `json:"users,omitempty"`
    Memberships     map[string]interface{}  `json:"memberships,omitempty"`
}

// UsergroupResponse is a struct representing collection response from API
type UsergroupsResponse struct{
    ListParams
    Items           []UsergroupResponse     `json:"items"`
}

// UsergroupResponse is a struct representing item create request for API
type UsergroupRequestCreate struct {
    Name            string                  `json:"name,omitempty"`
    Custom          map[string]interface{}  `json:"custom,omitempty"`
}

// UsergroupResponse is a struct representing item update request for API
type UsergroupRequestUpdate struct {
    Name            string                  `json:"name,omitempty"`
    Custom          map[string]interface{}  `json:"custom,omitempty"`
}

// TenantLink returns indicator of Tenant expansion and link to tenant.
// If expansion for Tenant was requested and resource is available via pointer
// it returns true, otherwise false. Link (href) is always returned. 
func (d *Usergroup) TenantLink() (bool, string) {
    return (d.Tenant != nil), d.tenant
}

// ApplicationsLink returns indicator of Usergroup expansion and link to dierctory.
// If expansion for Usergroup was requested and resource is available via pointer
// it returns true, otherwise false. Link (href) is always returned. 
func (d *Usergroup) DirectoryLink() (bool, string) {
    return (d.Directory != nil), d.directory
}

// ApplicationsLink returns indicator of Usergroup expansion and link to dierctory.
// If expansion for Usergroup was requested and resource is available via pointer
// it returns true, otherwise false. Link (href) is always returned. 
func (d *Usergroup) UsersLink() (bool, string) {
    return (d.Users != nil), d.users
}

// MembershipsLink returns indicator of Usergroup expansion and link to dierctory.
// If expansion for Usergroup was requested and resource is available via pointer
// it returns true, otherwise false. Link (href) is always returned. 
func (d *Usergroup) MembershipsLink() (bool, string) {
    return (d.Memberships != nil), d.memberships
}

// Save is a helper method for updating apikey.
// It calls UpdateByLink() on service under the hood.
func (t *Usergroup) Save() error {
    tmp := &UsergroupRequestUpdate{}
    copier.Copy(tmp, t)
    ten, err := t.service.UpdateByLink(t.Href, tmp)
    if err != nil {
        return err
    }

    tmpTenant := t.Tenant
    tmpDirectory := t.Directory
    tmpUsers := t.Users
    tmpMemberships := t.Memberships

    *t = *ten
    t.Tenant = tmpTenant
    t.Directory = tmpDirectory
    t.Users = tmpUsers
    t.Memberships = tmpMemberships

    return nil
}

// Save is a helper method for deleting apikey.
// It calls Delete() on service under the hood.
func (t *Usergroup) Delete() error {
    err := t.service.Delete(t)
    if err != nil {
        return err
    }

    return nil
}

// GetById retrieves apikey by its ID
func (s *UsergroupsServiceOp) GetById(id string, args ...interface{}) (*Usergroup, error) {
    endpoint := "usergroups/"
    endpoint = fmt.Sprintf("%s%s", endpoint, id)

    return s.GetByLink(endpoint, args...)
}

// GetById retrieves apikey by its full link
func (s *UsergroupsServiceOp) GetByLink(endpoint string, args ...interface{}) (*Usergroup, error) {
    resp, err := s.client.request("GET", endpoint, nil, args...)
    if err != nil {
        return nil, err
    }

    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("Status code: %d", resp.StatusCode)
    }
    obj := &UsergroupResponse{}
    dec := json.NewDecoder(resp.Body)
    dec.Decode(obj)

    return s.get(obj)
}

// get is internal method for transforming UsergroupResponse into Usergroup
func (s *UsergroupsServiceOp) get(r *UsergroupResponse) (*Usergroup, error) {
    obj := &Usergroup{}
    copier.Copy(obj, r)
    if v, ok :=  r.Tenant["href"]; ok {
        obj.tenant = v.(string)
    }
    if v, ok :=  r.Directory["href"]; ok {
        obj.directory = v.(string)
    }
    if v, ok :=  r.Users["href"]; ok {
        obj.users = v.(string)
    }
    if v, ok :=  r.Memberships["href"]; ok {
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
   if len(r.Users) > 1 {        
        bytes, err := json.Marshal(r.Users)
        if err != nil {
            return nil, err
        }
        ten := &UsersResponse{}
        json.Unmarshal(bytes, ten)
        t, _, err := s.client.Users.getCollection(ten)
        if err != nil {
            return nil, err
        }
        obj.Users = t
    }
    if len(r.Directory) > 1 {        
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
    }
    obj.service = s
    return obj, nil
}

// get is internal method for transforming ApplicationResponse into Application
func (s *UsergroupsServiceOp) getCollection(r *UsergroupsResponse) ([]Usergroup, *ListParams, error) {
    dst := make([]Usergroup, len(r.Items))

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

// GetById retrieves collection of usergroups of current tenant
func (s *UsergroupsServiceOp) ListByDirectory(id string, args ...interface{}) ([]Usergroup, *ListParams, error) {
    endpoint := fmt.Sprintf("directories/%s/usergroups", id)
    return s.ListByLink(endpoint, args...)
}

// GetById retrieves collection of usergroups by link
func (s *UsergroupsServiceOp) ListByLink(endpoint string, args ...interface{}) ([]Usergroup, *ListParams, error) {
    resp, err := s.client.request("GET", endpoint, nil, args...)
    if err != nil {
        return nil, nil, err
    }

    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, nil, fmt.Errorf("Status code: %d", resp.StatusCode)
    }
    obj := &UsergroupsResponse{}
    dec := json.NewDecoder(resp.Body)
    dec.Decode(obj)

    return s.getCollection(obj)
}

// GetById updates apikey with specified ID
func (s *UsergroupsServiceOp) UpdateById(id string, t *UsergroupRequestUpdate) (*Usergroup, error) {
    endpoint := fmt.Sprintf("usergroups/%s", id)
    return s.UpdateByLink(endpoint, t)
}

// GetById updates apikey specified by link
func (s *UsergroupsServiceOp) UpdateByLink(endpoint string, t *UsergroupRequestUpdate) (*Usergroup, error) {
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
    obj := &UsergroupResponse{}
    dec := json.NewDecoder(resp.Body)
    dec.Decode(obj)
    return s.get(obj)
}

func (s *UsergroupsServiceOp) CreateByDirectory(id string, dir *UsergroupRequestCreate) (*Usergroup, error) {
    endpoint := fmt.Sprintf("directories/%s/usergroups", id)
    return s.CreateByLink(endpoint, dir)
}

// Create creates new apikey within tenant
func (s *UsergroupsServiceOp) CreateByLink(endpoint string, dir *UsergroupRequestCreate) (*Usergroup, error) {
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
    obj := &UsergroupResponse{}
    dec := json.NewDecoder(resp.Body)
    dec.Decode(obj)
    return s.get(obj)
}

// Delete removes apikey
func (s *UsergroupsServiceOp) Delete(t *Usergroup) (error) {
    return s.DeleteByLink(t.Href)
}

// Delete removes apikey by ID
func (s *UsergroupsServiceOp) DeleteById(id string) (error) {
    endpoint := fmt.Sprintf("usergroups/%s", id)
    return s.DeleteByLink(endpoint)
}

// Delete removes apikey by link
func (s *UsergroupsServiceOp) DeleteByLink(endpoint string) (error) {
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
