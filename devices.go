package api

import (
    "bytes"
    "encoding/json" 
    "fmt"   
    "net/http"
)

// DevicesService is an interafce for interacting with Devices endpoints of CloudThing API
// https://tenant-name.cloudthing.io/api/v1/devices

type DevicesService interface {
    GetById(string) (*Device, error)
    GetByHref(string) (*Device, error)
    List(*ListOptions) ([]Device, *ListParams, error)
    ListByHref(string, *ListOptions) ([]Device, *ListParams, error)
    Create(*Device) (*Device, error)
    Update(*Device) (*Device, error)
    Delete(*Device) (error)
    DeleteByHref(string) (error)
    DeleteById(string) (error)
}

// DevicesServiceOp handles communication with Tenant related methods of API
type DevicesServiceOp struct {
    client *Client
}

type DeviceProperty struct {
    Key             string          `json:"key"`
    Value           interface{}     `json:"value"`
}

type Device struct {
    ModelBase
    Token               string          `json:"token,omitempty"`
    Activated           *bool           `json:"activated"`
    Custom              interface{}     `json:"custom,omitempty"`
    Properties          []DeviceProperty`json:"properties,omitempty"`

    tenant              string          `json:"tenant,omitempty"`
    product             string          `json:"product,omitempty"`
    clusters            string          `json:"clusters,omitempty"`
    groups              string          `json:"groups,omitempty"`
    clusterMemberships  string          `json:"clusterMemberships,omitempty"`
    groupMemberships    string          `json:"groupMemberships,omitempty"`


    // service for communication, internal use only
    service         *DevicesServiceOp `json:"-"` 
}

type Devices struct{
    ListParams
    Items           []Device     `json:"items"`
}

func (d *Device) Tenant() (*Tenant, error) {
    return d.service.client.Tenant.Get()
}

func (d *Device) Product() (*Product, error) {
    return d.service.client.Products.GetByHref(d.product)
}

func (d *Device) Clusters() ([]Cluster, *ListParams, error) {
    return d.service.client.Clusters.ListByHref(d.clusters, nil)
}

func (d *Device) Groups() ([]Group, *ListParams, error) {
    return d.service.client.Groups.ListByHref(d.groups, nil)
}

func (d *Device) ClusterMemberships() ([]ClusterMembership, *ListParams, error) {
    return d.service.client.ClusterMemberships.ListByHref(d.clusterMemberships, nil)
}

func (d *Device) GroupMemberships() ([]GroupMembership, *ListParams, error) {
    return d.service.client.GroupMemberships.ListByHref(d.groupMemberships, nil)
}

// Save updates tenant by calling Update() on service under the hood
func (t *Device) Save() error {
    tmp := *t

    ten, err := t.service.Update(&tmp)
    if err != nil {
        return err
    }

    *t = *ten
    return nil
}

// GetById retrieves directory
func (s *DevicesServiceOp) GetById(id string) (*Device, error) {
    endpoint := "devices/"
    endpoint = fmt.Sprintf("%s%s", endpoint, id)

    return s.GetByHref(endpoint)
}

func (s *DevicesServiceOp) GetByHref(endpoint string) (*Device, error) {
    resp, err := s.client.request("GET", endpoint, nil)
    if err != nil {
        return nil, err
    }

    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("Status code: %d", resp.StatusCode)
    }
    obj := &Device{}
    dec := json.NewDecoder(resp.Body)
    dec.Decode(obj)
    obj.service = s
    return obj, nil
}

func (s *DevicesServiceOp) List(lo *ListOptions) ([]Device, *ListParams, error) {
    endpoint := fmt.Sprintf("tenants/%s/devices", s.client.tenantId)

    return s.ListByHref(endpoint, lo)
}

func (s *DevicesServiceOp) ListByHref(endpoint string, lo *ListOptions) ([]Device, *ListParams, error) {
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
    obj := &Devices{}
    dec := json.NewDecoder(resp.Body)
    dec.Decode(obj)

    dst := make([]Device, len(obj.Items))
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
func (s *DevicesServiceOp) Update(t *Device) (*Device, error) {
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
    obj := &Device{}
    dec := json.NewDecoder(resp.Body)
    dec.Decode(obj)
    obj.service = s
    return obj, nil
}

func (s *DevicesServiceOp) Create(dir *Device) (*Device, error) {
    endpoint := fmt.Sprintf("devices")

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
    obj := &Device{}
    dec := json.NewDecoder(resp.Body)
    dec.Decode(obj)
    obj.service = s
    return obj, nil
}

// Delete removes application
func (s *DevicesServiceOp) Delete(t *Device) (error) {
    return s.DeleteByHref(t.Href)
}

// Delete removes application by ID
func (s *DevicesServiceOp) DeleteById(id string) (error) {
    endpoint := fmt.Sprintf("devices/%s", id)
    return s.DeleteByHref(endpoint)
}

// Delete removes application by link
func (s *DevicesServiceOp) DeleteByHref(endpoint string) (error) {
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
