package api

import (
    "bytes"
    "encoding/json" 
    "fmt"   
    "net/http"
)

// GroupsService is an interafce for interacting with Groups endpoints of CloudThing API
// https://tenant-name.cloudthing.io/api/v1/groups

type GroupsService interface {
    GetById(string) (*Group, error)
    GetByLink(string) (*Group, error)
    ListByLink(string, *ListOptions) ([]Group, *ListParams, error)
    Create(*Group) (*Group, error)
    Update(*Group) (*Group, error)
    Delete(*Group) (error)
    DeleteByLink(string) (error)
    DeleteById(string) (error)
}

// GroupsServiceOp handles communication with Tenant related methods of API
type GroupsServiceOp struct {
    client *Client
}

type Group struct {
    ModelBase
    Name                string          `json:"name,omitempty"`
    Description         string          `json:"description,omitempty"`
    Custom              interface{}     `json:"custom,omitempty"`

    tenant              string          `json:"tenant,omitempty"`
    application         string          `json:"application,omitempty"`
    cluster             string          `json:"cluster,omitempty"`
    resources           string          `json:"resources,omitempty"`
    devices             string          `json:"devices,omitempty"`

    // service for communication, internal use only
    service         *GroupsServiceOp `json:"-"` 
}

type Groups struct{
    ListParams
    Items           []Group     `json:"items"`
}

func (d *Group) Tenant() (*Tenant, error) {
    return d.service.client.Tenant.Get()
}

func (d *Group) Application() (*Application, error) {
    return d.service.client.Applications.GetByLink(d.application)
}

func (d *Group) Cluster() (*Cluster, error) {
    return d.service.client.Clusters.GetByLink(d.cluster)
}

func (d *Group) Devices() ([]Device, *ListParams, error) {
    return d.service.client.Devices.ListByLink(d.devices, nil)
}

// Save updates tenant by calling Update() on service under the hood
func (t *Group) Save() error {
    tmp := *t

    ten, err := t.service.Update(&tmp)
    if err != nil {
        return err
    }

    *t = *ten
    return nil
}

// GetById retrieves directory
func (s *GroupsServiceOp) GetById(id string) (*Group, error) {
    endpoint := "groups/"
    endpoint = fmt.Sprintf("%s%s", endpoint, id)

    return s.GetByLink(endpoint)
}

func (s *GroupsServiceOp) GetByLink(endpoint string) (*Group, error) {
    resp, err := s.client.request("GET", endpoint, nil)
    if err != nil {
        return nil, err
    }

    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("Status code: %d", resp.StatusCode)
    }
    obj := &Group{}
    dec := json.NewDecoder(resp.Body)
    dec.Decode(obj)
    obj.service = s
    return obj, nil
}
func (s *GroupsServiceOp) ListByLink(endpoint string, lo *ListOptions) ([]Group, *ListParams, error) {
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
    obj := &Groups{}
    dec := json.NewDecoder(resp.Body)
    dec.Decode(obj)

    dst := make([]Group, len(obj.Items))
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
func (s *GroupsServiceOp) Update(t *Group) (*Group, error) {
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
    obj := &Group{}
    dec := json.NewDecoder(resp.Body)
    dec.Decode(obj)
    obj.service = s
    return obj, nil
}

func (s *GroupsServiceOp) Create(dir *Group) (*Group, error) {
    endpoint := fmt.Sprintf("groups")

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
    obj := &Group{}
    dec := json.NewDecoder(resp.Body)
    dec.Decode(obj)
    obj.service = s
    return obj, nil
}

// Delete removes application
func (s *GroupsServiceOp) Delete(t *Group) (error) {
    return s.DeleteByLink(t.Href)
}

// Delete removes application by ID
func (s *GroupsServiceOp) DeleteById(id string) (error) {
    endpoint := fmt.Sprintf("groups/%s", id)
    return s.DeleteByLink(endpoint)
}

// Delete removes application by link
func (s *GroupsServiceOp) DeleteByLink(endpoint string) (error) {
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
