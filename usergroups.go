package api

import (
    "bytes"
    "encoding/json" 
    "fmt"   
    "net/http"
)

// UsergroupsService is an interafce for interacting with Userusergroups endpoints of CloudThing API
// https://tenant-name.cloudthing.io/api/v1/usergroups

type UsergroupsService interface {
    GetById(string) (*Usergroup, error)
    GetByLink(string) (*Usergroup, error)
    ListByLink(string, *ListOptions) ([]Usergroup, *ListParams, error)
    Create(*Usergroup) (*Usergroup, error)
    Update(*Usergroup) (*Usergroup, error)
    Delete(*Usergroup) (error)
    DeleteByLink(string) (error)
    DeleteById(string) (error)
}

// UsergroupsServiceOp handles communication with Tenant related methods of API
type UsergroupsServiceOp struct {
    client *Client
}

type Usergroup struct {
    ModelBase
    Name                string          `json:"name,omitempty"`
    Custom              map[string]interface{} `json:"custom,omitempty"`

    tenant              string          `json:"tenant,omitempty"`
    directory           string          `json:"directory,omitempty"`
    users               string          `json:"users,omitempty"`
    memberships         string          `json:"memberships,omitempty"`

    // service for communication, internal use only
    service         *UsergroupsServiceOp `json:"-"` 
}

type Usergroups struct{
    ListParams
    Items           []Usergroup     `json:"items"`
}

func (d *Usergroup) Tenant() (*Tenant, error) {
    return d.service.client.Tenant.Get()
}

func (d *Usergroup) Directory() (*Directory, error) {
    return d.service.client.Directories.GetByLink(d.directory)
}

func (d *Usergroup) Users() ([]User, *ListParams, error) {
    return d.service.client.Users.ListByLink(d.users, nil)
}

func (d *Usergroup) Memberships() ([]Membership, *ListParams, error) {
    return d.service.client.Memberships.ListByLink(d.memberships, nil)
}

// Save updates tenant by calling Update() on service under the hood
func (t *Usergroup) Save() error {
    tmp := *t

    ten, err := t.service.Update(&tmp)
    if err != nil {
        return err
    }

    *t = *ten
    return nil
}

func (t *Usergroup) Delete() error {
    err := t.service.Delete(t)
    if err != nil {
        return err
    }

    return nil
}

// GetById retrieves directory
func (s *UsergroupsServiceOp) GetById(id string) (*Usergroup, error) {
    endpoint := "usergroups/"
    endpoint = fmt.Sprintf("%s%s", endpoint, id)

    return s.GetByLink(endpoint)
}

func (s *UsergroupsServiceOp) GetByLink(endpoint string) (*Usergroup, error) {
    resp, err := s.client.request("GET", endpoint, nil)
    if err != nil {
        return nil, err
    }

    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("Status code: %d", resp.StatusCode)
    }
    obj := &Usergroup{}
    dec := json.NewDecoder(resp.Body)
    dec.Decode(obj)
    obj.service = s
    return obj, nil
}
func (s *UsergroupsServiceOp) ListByLink(endpoint string, lo *ListOptions) ([]Usergroup, *ListParams, error) {
    if lo == nil {
        lo = &ListOptions {
            Page: 1,
            Limit: DefaultLimit,
        }
    }
    endpoint = fmt.Sprintf("%s?limit=%d&page=%d", endpoint, lo.Limit, lo.Page)

    resp, err := s.client.request("GET", endpoint, nil)
    if err != nil {
        return nil, nil, err
    }

    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, nil, fmt.Errorf("Status code: %d", resp.StatusCode)
    }
    obj := &Usergroups{}
    dec := json.NewDecoder(resp.Body)
    dec.Decode(obj)

    dst := make([]Usergroup, len(obj.Items))
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

// Update updates product
func (s *UsergroupsServiceOp) Update(t *Usergroup) (*Usergroup, error) {
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
    obj := &Usergroup{}
    dec := json.NewDecoder(resp.Body)
    dec.Decode(obj)
    obj.service = s
    return obj, nil
}

func (s *UsergroupsServiceOp) Create(dir *Usergroup) (*Usergroup, error) {
    endpoint := fmt.Sprintf("usergroups")

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
    obj := &Usergroup{}
    dec := json.NewDecoder(resp.Body)
    dec.Decode(obj)
    obj.service = s
    return obj, nil
}

// Delete removes application
func (s *UsergroupsServiceOp) Delete(t *Usergroup) (error) {
    return s.DeleteByLink(t.Href)
}

// Delete removes application by ID
func (s *UsergroupsServiceOp) DeleteById(id string) (error) {
    endpoint := fmt.Sprintf("usergroups/%s", id)
    return s.DeleteByLink(endpoint)
}

// Delete removes application by link
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
