package api

import (
    "bytes"
    "encoding/json" 
    "fmt"   
    "net/http"
)

// ApikeysService is an interafce for interacting with Apikeys endpoints of CloudThing API
// https://tenant-name.cloudthing.io/api/v1/directories

type ApikeysService interface {
    GetById(string) (*Apikey, error)
    GetByHref(string) (*Apikey, error)
    ListByHref(string, *ListOptions) ([]Apikey, *ListParams, error)
    Create(*Apikey) (*Apikey, error)
    Update(*Apikey) (*Apikey, error)
    Delete(*Apikey) (error)
    DeleteByHref(string) (error)
    DeleteById(string) (error)
}

// ApikeysServiceOp handles communication with Tenant related methods of API
type ApikeysServiceOp struct {
    client *Client
}

type Apikey struct {
    ModelBase
    Name            string          `json:"name,omitempty"`
    Description     string          `json:"description,omitempty"`
    Key             string          `json:"key,omitempty"`
    Secret          string          `json:"secret,omitempty"`
    Status          string          `json:"status,omitempty"`
    Custom          interface{}     `json:"custom,omitempty"`

    tenant          string          `json:"tenant,omitempty"`
    applications         string          `json:"applications,omitempty"`  

    // service for communication, internal use only
    service         *ApikeysServiceOp `json:"-"` 
}

type Apikeys struct{
    ListParams
    Items           []Apikey     `json:"items"`
}

func (d *Apikey) Tenant() (*Tenant, error) {
    return d.service.client.Tenant.Get()
}

func (d *Apikey) Applications() ([]Application, *ListParams, error) {
    return d.service.client.Applications.ListByHref(d.applications, nil)
}

// Save updates application by calling Update() on service under the hood
func (t *Apikey) Save() error {
    tmp := *t

    ten, err := t.service.Update(&tmp)
    if err != nil {
        return err
    }

    *t = *ten
    return nil
}

func (t *Apikey) Delete() error {
    err := t.service.Delete(t)
    if err != nil {
        return err
    }

    return nil
}

// GetById retrieves application
func (s *ApikeysServiceOp) GetById(id string) (*Apikey, error) {
    endpoint := "apikeys/"
    endpoint = fmt.Sprintf("%s%s", endpoint, id)

    return s.GetByHref(endpoint)
}

func (s *ApikeysServiceOp) GetByHref(endpoint string) (*Apikey, error) {
    resp, err := s.client.request("GET", endpoint, nil)
    if err != nil {
        return nil, err
    }

    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("Status code: %d", resp.StatusCode)
    }
    obj := &Apikey{}
    dec := json.NewDecoder(resp.Body)
    dec.Decode(obj)
    obj.service = s
    return obj, nil
}

func (s *ApikeysServiceOp) ListByHref(endpoint string, lo *ListOptions) ([]Apikey, *ListParams, error) {
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
    obj := &Apikeys{}
    dec := json.NewDecoder(resp.Body)
    dec.Decode(obj)

    dst := make([]Apikey, len(obj.Items))
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
func (s *ApikeysServiceOp) Update(t *Apikey) (*Apikey, error) {
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
    obj := &Apikey{}
    dec := json.NewDecoder(resp.Body)
    dec.Decode(obj)
    obj.service = s
    return obj, nil
}

func (s *ApikeysServiceOp) Create(dir *Apikey) (*Apikey, error) {
    endpoint := fmt.Sprintf("tenants/%s/apikeys", s.client.tenantId)

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
    obj := &Apikey{}
    dec := json.NewDecoder(resp.Body)
    dec.Decode(obj)
    obj.service = s
    return obj, nil
}

// Delete removes application
func (s *ApikeysServiceOp) Delete(t *Apikey) (error) {
    return s.DeleteByHref(t.Href)
}

// Delete removes application by ID
func (s *ApikeysServiceOp) DeleteById(id string) (error) {
    endpoint := fmt.Sprintf("apikeys/%s", id)
    return s.DeleteByHref(endpoint)
}

// Delete removes application by link
func (s *ApikeysServiceOp) DeleteByHref(endpoint string) (error) {
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
