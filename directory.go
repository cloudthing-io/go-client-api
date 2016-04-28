package api

import (
    "bytes"
    "encoding/json" 
    "fmt"   
    "net/http"
)

// DirectoriesService is an interafce for interacting with Directories endpoints of CloudThing API
// https://tenant-name.cloudthing.io/api/v1/directories

type DirectoriesService interface {
    GetById(string) (*Directory, error)
    GetByLink(string) (*Directory, error)
    List(*ListOptions) ([]Directory, *ListParams, error)
    Create(*Directory) (*Directory, error)
    Update(*Directory) (*Directory, error)
    Delete(*Directory) (error)
    DeleteByLink(string) (error)
    DeleteById(string) (error)
}

// DirectoriesServiceOp handles communication with Tenant related methods of API
type DirectoriesServiceOp struct {
    client *Client
}

type Directory struct {
    ModelBase
    Name            string          `json:"name,omitempty"`
    Official        *bool            `json:"official,omitempty"`
    Description     string          `json:"description,omitempty"`
    Custom          interface{}     `json:"custom,omitempty"`

    tenant          string          `json:"tenant,omitempty"`
    users           string          `json:"users,omitempty"`
    usergroups      string          `json:"usergroups,omitempty"`    
    applications    string          `json:"applications,omitempty"`

    // service for communication, internal use only
    service         *DirectoriesServiceOp `json:"-"` 
}

type Directories struct{
    ListParams
    Items           []Directory     `json:"items"`
}

func (d *Directory) Tenant() (*Tenant, error) {
    return d.service.client.Tenant.Get()
}

func (d *Directory) Users() ([]User, *ListParams, error) {
    return d.service.client.Users.ListByLink(d.users, nil)
}

func (d *Directory) Usergroups() ([]Usergroup, *ListParams, error) {
    return d.service.client.Usergroups.ListByLink(d.usergroups, nil)
}

func (d *Directory) Applications() ([]Application, *ListParams, error) {
    return d.service.client.Applications.ListByLink(d.applications, nil)
}


// Save updates tenant by calling Update() on service under the hood
func (t *Directory) Save() error {
    tmp := *t

    ten, err := t.service.Update(&tmp)
    if err != nil {
        return err
    }

    *t = *ten
    return nil
}

// GetById retrieves directory
func (s *DirectoriesServiceOp) GetById(id string) (*Directory, error) {
    endpoint := "directories/"
    endpoint = fmt.Sprintf("%s%s", endpoint, id)

    return s.GetByLink(endpoint)
}


func (s *DirectoriesServiceOp) GetByLink(endpoint string) (*Directory, error) {
    resp, err := s.client.request("GET", endpoint, nil)
    if err != nil {
        return nil, err
    }

    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("Status code: %d", resp.StatusCode)
    }
    obj := &Directory{}
    dec := json.NewDecoder(resp.Body)
    dec.Decode(obj)
    obj.service = s
    return obj, nil
}
func (s *DirectoriesServiceOp) List(lo *ListOptions) ([]Directory, *ListParams, error) {
    if lo == nil {
        lo = &ListOptions {
            Page: 1,
            Limit: DefaultLimit,
        }
    }
    endpoint := fmt.Sprintf("tenants/%s/directories?limit=%d&page=%d", s.client.tenantId, lo.Limit, lo.Page)

    resp, err := s.client.request("GET", endpoint, nil)
    if err != nil {
        return nil, nil, err
    }

    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, nil, fmt.Errorf("Status code: %d", resp.StatusCode)
    }
    obj := &Directories{}
    dec := json.NewDecoder(resp.Body)
    dec.Decode(obj)

    dst := make([]Directory, len(obj.Items))
    copy(dst, obj.Items)
    for i, _ := range dst {
        dst[i].service = s
    }

    lp := &ListParams {
        Href: obj.Href,
        Prev: obj.Prev,
        Next: obj.Next,
        Limit: obj.Limit,
        Size: obj.Size,
        Page: obj.Page,
    }
    return dst, lp, nil
}

// Update updates tenant
func (s *DirectoriesServiceOp) Update(t *Directory) (*Directory, error) {
    endpoint := t.Href

    t.CreatedAt = nil
    t.UpdatedAt = nil
    t.Href = ""

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
    obj := &Directory{}
    dec := json.NewDecoder(resp.Body)
    dec.Decode(obj)
    obj.service = s
    return obj, nil
}

func (s *DirectoriesServiceOp) Create(dir *Directory) (*Directory, error) {
    endpoint := fmt.Sprintf("tenants/%s/directories", s.client.tenantId)

    dir.CreatedAt = nil
    dir.UpdatedAt = nil
    dir.Href = ""

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
    obj := &Directory{}
    dec := json.NewDecoder(resp.Body)
    dec.Decode(obj)
    obj.service = s
    return obj, nil
}

// Delete removes application
func (s *DirectoriesServiceOp) Delete(t *Directory) (error) {
    return s.DeleteByLink(t.Href)
}

// Delete removes application by ID
func (s *DirectoriesServiceOp) DeleteById(id string) (error) {
    endpoint := fmt.Sprintf("directories/%s", id)
    return s.DeleteByLink(endpoint)
}

// Delete removes application by link
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
