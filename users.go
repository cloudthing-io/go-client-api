package api

import (
    "bytes"
    "encoding/json" 
    "fmt"   
    "net/http"
    "time"
)

// UsersService is an interafce for interacting with Users endpoints of CloudThing API
// https://tenant-name.cloudthing.io/api/v1/directories

type UsersService interface {
    GetById(string) (*User, error)
    GetByHref(string) (*User, error)
    ListByHref(string, *ListOptions) ([]User, *ListParams, error)
    Create(*User) (*User, error)
    Update(*User) (*User, error)
    Delete(*User) (error)
    DeleteByHref(string) (error)
    DeleteById(string) (error)
}

// UsersServiceOp handles communication with Tenant related methods of API
type UsersServiceOp struct {
    client *Client
}

type User struct {
    ModelBase
    Username        string          `json:"username,omitempty"`
    Email           string          `json:"email,omitempty"`
    FirstName       string          `json:"firstName,omitempty"`
    Surname         string          `json:"surname,omitempty"`
    LastSuccessfulLogin       *time.Time       `json:"lastSuccessfulLogin,omitempty"`
    LastFailedLogin       *time.Time       `json:"lastFailedLogin,omitempty"`
    Custom          interface{}     `json:"custom,omitempty"`

    tenant          string          `json:"tenant,omitempty"`
    directory       string          `json:"directory,omitempty"`
    applications         string          `json:"applications,omitempty"`
    usergroups       string          `json:"usergroups,omitempty"`
    memberships         string          `json:"memberships,omitempty"`    

    // service for communication, internal use only
    service         *UsersServiceOp `json:"-"` 
}

type Users struct{
    ListParams
    Items           []User     `json:"items"`
}

func (d *User) Tenant() (*Tenant, error) {
    return d.service.client.Tenant.Get()
}

func (d *User) Directory() (*Directory, error) {
    return d.service.client.Directories.GetByHref(d.directory)
}

func (d *User) Applications() ([]Application, *ListParams, error) {
    return d.service.client.Applications.ListByHref(d.applications, nil)
}

func (d *User) Usergroups() ([]Usergroup, *ListParams, error) {
    return d.service.client.Usergroups.ListByHref(d.usergroups, nil)
}

func (d *User) Memberships() ([]Membership, *ListParams, error) {
    return d.service.client.Memberships.ListByHref(d.memberships, nil)
}

// Save updates application by calling Update() on service under the hood
func (t *User) Save() error {
    tmp := *t

    ten, err := t.service.Update(&tmp)
    if err != nil {
        return err
    }

    *t = *ten
    return nil
}

func (t *User) Delete() error {
    err := t.service.Delete(t)
    if err != nil {
        return err
    }

    return nil
}

// GetById retrieves application
func (s *UsersServiceOp) GetById(id string) (*User, error) {
    endpoint := "users/"
    endpoint = fmt.Sprintf("%s%s", endpoint, id)

    return s.GetByHref(endpoint)
}

func (s *UsersServiceOp) GetByHref(endpoint string) (*User, error) {
    resp, err := s.client.request("GET", endpoint, nil)
    if err != nil {
        return nil, err
    }

    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("Status code: %d", resp.StatusCode)
    }
    obj := &User{}
    dec := json.NewDecoder(resp.Body)
    dec.Decode(obj)
    obj.service = s
    return obj, nil
}

func (s *UsersServiceOp) ListByHref(endpoint string, lo *ListOptions) ([]User, *ListParams, error) {
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
    obj := &Users{}
    dec := json.NewDecoder(resp.Body)
    dec.Decode(obj)

    dst := make([]User, len(obj.Items))
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
func (s *UsersServiceOp) Update(t *User) (*User, error) {
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
    obj := &User{}
    dec := json.NewDecoder(resp.Body)
    dec.Decode(obj)
    obj.service = s
    return obj, nil
}

func (s *UsersServiceOp) Create(dir *User) (*User, error) {
    endpoint := fmt.Sprintf("tenants/%s/users", s.client.tenantId)

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
    obj := &User{}
    dec := json.NewDecoder(resp.Body)
    dec.Decode(obj)
    obj.service = s
    return obj, nil
}

// Delete removes application
func (s *UsersServiceOp) Delete(t *User) (error) {
    return s.DeleteByHref(t.Href)
}

// Delete removes application by ID
func (s *UsersServiceOp) DeleteById(id string) (error) {
    endpoint := fmt.Sprintf("users/%s", id)
    return s.DeleteByHref(endpoint)
}

// Delete removes application by link
func (s *UsersServiceOp) DeleteByHref(endpoint string) (error) {
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
