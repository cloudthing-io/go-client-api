package api

import (
    "bytes"
    "encoding/json" 
    "fmt"   
    "net/http"
    "strings"
)

// ApplicationsService is an interafce for interacting with Applications endpoints of CloudThing API
// https://tenant-name.cloudthing.io/api/v1/directories

type ApplicationsService interface {
    GetById(string) (*Application, error)
    GetByHref(string) (*Application, error)
    List(*ListOptions) ([]Application, *ListParams, error)
    ListByHref(string, *ListOptions) ([]Application, *ListParams, error)
    Create(*Application) (*Application, error)
    Update(*Application) (*Application, error)
    Delete(*Application) (error)
    DeleteByHref(string) (error)
    DeleteById(string) (error)
}

// ApplicationsServiceOp handles communication with Tenant related methods of API
type ApplicationsServiceOp struct {
    client *Client
}

type Application struct {
    ModelBase
    Name            string          `json:"name,omitempty"`
    Official        bool            `json:"official,omitempty"`
    Description     string          `json:"description,omitempty"`
    Status          string          `json:"status,omitempty"`
    Custom          interface{}     `json:"custom,omitempty"`

    tenant          string          `json:"tenant,omitempty"`
    directory       string          `json:"directory,omitempty"`
    devices         string          `json:"devices,omitempty"`   

    // service for communication, internal use only
    service         *ApplicationsServiceOp `json:"-"` 
}

type Applications struct{
    ListParams
    Items           []Application     `json:"items"`
}

func (d *Application) Tenant() (*Tenant, error) {
    return d.service.client.Tenant.Get()
}

func (d *Application) Directory() (*Directory, error) {
    id := strings.Split(d.directory, "/")
    return d.service.client.Directories.GetById(id[len(id)-1])
}

func (d *Application) Devices() ([]Device, *ListParams, error) {
    return d.service.client.Devices.ListByHref(d.devices, nil)
}

// Save updates application by calling Update() on service under the hood
func (t *Application) Save() error {
    tmp := *t

    ten, err := t.service.Update(&tmp)
    if err != nil {
        return err
    }

    *t = *ten
    return nil
}

func (t *Application) Delete() error {
    err := t.service.Delete(t)
    if err != nil {
        return err
    }

    return nil
}

// GetById retrieves application
func (s *ApplicationsServiceOp) GetById(id string) (*Application, error) {
    endpoint := "applications/"
    endpoint = fmt.Sprintf("%s%s", endpoint, id)

    return s.GetByHref(endpoint)
}

func (s *ApplicationsServiceOp) GetByHref(endpoint string) (*Application, error) {
    resp, err := s.client.request("GET", endpoint, nil)
    if err != nil {
        return nil, err
    }

    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("Status code: %d", resp.StatusCode)
    }
    obj := &Application{}
    dec := json.NewDecoder(resp.Body)
    dec.Decode(obj)
    obj.service = s
    return obj, nil
}

func (s *ApplicationsServiceOp) List(lo *ListOptions) ([]Application, *ListParams, error) {
    endpoint := fmt.Sprintf("tenants/%s/applications", s.client.tenantId)
    return s.ListByHref(endpoint, lo)
}

func (s *ApplicationsServiceOp) ListByHref(endpoint string, lo *ListOptions) ([]Application, *ListParams, error) {
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
    obj := &Applications{}
    dec := json.NewDecoder(resp.Body)
    dec.Decode(obj)

    dst := make([]Application, len(obj.Items))
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
func (s *ApplicationsServiceOp) Update(t *Application) (*Application, error) {
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
    obj := &Application{}
    dec := json.NewDecoder(resp.Body)
    dec.Decode(obj)
    obj.service = s
    return obj, nil
}

func (s *ApplicationsServiceOp) Create(dir *Application) (*Application, error) {
    endpoint := fmt.Sprintf("tenants/%s/applications", s.client.tenantId)

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
    obj := &Application{}
    dec := json.NewDecoder(resp.Body)
    dec.Decode(obj)
    obj.service = s
    return obj, nil
}

// Delete removes application
func (s *ApplicationsServiceOp) Delete(t *Application) (error) {
    return s.DeleteByHref(t.Href)
}

// Delete removes application by ID
func (s *ApplicationsServiceOp) DeleteById(id string) (error) {
    endpoint := fmt.Sprintf("applications/%s", id)
    return s.DeleteByHref(endpoint)
}

// Delete removes application by link
func (s *ApplicationsServiceOp) DeleteByHref(endpoint string) (error) {
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
