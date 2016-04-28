package api

import (
    "bytes"
    "encoding/json" 
    "fmt"   
    "net/http"
)

// ClustersService is an interafce for interacting with Clusters endpoints of CloudThing API
// https://tenant-name.cloudthing.io/api/v1/clusters

type ClustersService interface {
    GetById(string) (*Cluster, error)
    GetByLink(string) (*Cluster, error)
    ListByLink(string, *ListOptions) ([]Cluster, *ListParams, error)
    Create(*Cluster) (*Cluster, error)
    Update(*Cluster) (*Cluster, error)
    Delete(*Cluster) (error)
    DeleteByLink(string) (error)
    DeleteById(string) (error)
}

// ClustersServiceOp handles communication with Tenant related methods of API
type ClustersServiceOp struct {
    client *Client
}

type Cluster struct {
    ModelBase
    Name                string          `json:"name,omitempty"`
    Description         string          `json:"description,omitempty"`
    Custom              interface{}     `json:"custom,omitempty"`

    tenant              string          `json:"tenant,omitempty"`
    application         string          `json:"application,omitempty"`
    groups              string          `json:"groups,omitempty"`
    memberships         string          `json:"memberships,omitempty"`
    users               string          `json:"users,omitempty"`
    resources           string          `json:"resources,omitempty"`
    devices             string          `json:"devices,omitempty"`

    // service for communication, internal use only
    service         *ClustersServiceOp `json:"-"` 
}

type Clusters struct{
    ListParams
    Items           []Cluster     `json:"items"`
}

func (d *Cluster) Tenant() (*Tenant, error) {
    return d.service.client.Tenant.Get()
}

func (d *Cluster) Application() (*Application, error) {
    return d.service.client.Applications.GetByLink(d.application)
}

func (d *Cluster) Groups() ([]Group, *ListParams, error) {
    return d.service.client.Groups.ListByLink(d.groups, nil)
}

func (d *Cluster) Memberships() ([]Membership, *ListParams, error) {
    return d.service.client.Memberships.ListByLink(d.memberships, nil)
}

func (d *Cluster) Users() ([]User, *ListParams, error) {
    return d.service.client.Users.ListByLink(d.users, nil)
}

func (d *Cluster) Devices() ([]Device, *ListParams, error) {
    return d.service.client.Devices.ListByLink(d.devices, nil)
}

// Save updates tenant by calling Update() on service under the hood
func (t *Cluster) Save() error {
    tmp := *t

    ten, err := t.service.Update(&tmp)
    if err != nil {
        return err
    }

    *t = *ten
    return nil
}

func (t *Cluster) Delete() error {
    err := t.service.Delete(t)
    if err != nil {
        return err
    }

    return nil
}

// GetById retrieves directory
func (s *ClustersServiceOp) GetById(id string) (*Cluster, error) {
    endpoint := "clusters/"
    endpoint = fmt.Sprintf("%s%s", endpoint, id)

    return s.GetByLink(endpoint)
}

func (s *ClustersServiceOp) GetByLink(endpoint string) (*Cluster, error) {
    resp, err := s.client.request("GET", endpoint, nil)
    if err != nil {
        return nil, err
    }

    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("Status code: %d", resp.StatusCode)
    }
    obj := &Cluster{}
    dec := json.NewDecoder(resp.Body)
    dec.Decode(obj)
    obj.service = s
    return obj, nil
}
func (s *ClustersServiceOp) ListByLink(endpoint string, lo *ListOptions) ([]Cluster, *ListParams, error) {
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
    obj := &Clusters{}
    dec := json.NewDecoder(resp.Body)
    dec.Decode(obj)

    dst := make([]Cluster, len(obj.Items))
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
func (s *ClustersServiceOp) Update(t *Cluster) (*Cluster, error) {
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
    obj := &Cluster{}
    dec := json.NewDecoder(resp.Body)
    dec.Decode(obj)
    obj.service = s
    return obj, nil
}

func (s *ClustersServiceOp) Create(dir *Cluster) (*Cluster, error) {
    endpoint := fmt.Sprintf("clusters")

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
    obj := &Cluster{}
    dec := json.NewDecoder(resp.Body)
    dec.Decode(obj)
    obj.service = s
    return obj, nil
}

// Delete removes application
func (s *ClustersServiceOp) Delete(t *Cluster) (error) {
    return s.DeleteByLink(t.Href)
}

// Delete removes application by ID
func (s *ClustersServiceOp) DeleteById(id string) (error) {
    endpoint := fmt.Sprintf("clusters/%s", id)
    return s.DeleteByLink(endpoint)
}

// Delete removes application by link
func (s *ClustersServiceOp) DeleteByLink(endpoint string) (error) {
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
