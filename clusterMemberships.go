package api

import (
    "bytes"
    "encoding/json" 
    "fmt"   
    "net/http"
)

// ClusterMembershipsService is an interafce for interacting with ClusterMemberships endpoints of CloudThing API
// https://tenant-name.cloudthing.io/api/v1/clusterMemberships

type ClusterMembershipsService interface {
    GetById(string) (*ClusterMembership, error)
    GetByHref(string) (*ClusterMembership, error)
    ListByHref(string, *ListOptions) ([]ClusterMembership, *ListParams, error)
    Create(*ClusterMembership) (*ClusterMembership, error)
    Delete(*ClusterMembership) (error)
    DeleteByHref(string) (error)
    DeleteById(string) (error)
}

// ClusterMembershipsServiceOp handles communication with Tenant related methods of API
type ClusterMembershipsServiceOp struct {
    client *Client
}

type ClusterMembership struct {
    ModelBase

    device          string          `json:"device,omitempty"`
    cluster         string          `json:"cluster,omitempty"`
    application     string          `json:"application,omitempty"`

    // service for communication, internal use only
    service         *ClusterMembershipsServiceOp `json:"-"` 
}

type ClusterMemberships struct{
    ListParams
    Items           []ClusterMembership     `json:"items"`
}

func (d *ClusterMembership) Application() (*Application, error) {
    return d.service.client.Applications.GetByHref(d.application)
}

func (d *ClusterMembership) Device() (*Device, error) {
    return d.service.client.Devices.GetByHref(d.device)
}

func (d *ClusterMembership) Cluster() (*Cluster, error) {
    return d.service.client.Clusters.GetByHref(d.cluster)
}

// GetById retrieves directory
func (s *ClusterMembershipsServiceOp) GetById(id string) (*ClusterMembership, error) {
    endpoint := "clusterMemberships/"
    endpoint = fmt.Sprintf("%s%s", endpoint, id)

    return s.GetByHref(endpoint)
}

func (s *ClusterMembershipsServiceOp) GetByHref(endpoint string) (*ClusterMembership, error) {
    resp, err := s.client.request("GET", endpoint, nil)
    if err != nil {
        return nil, err
    }

    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("Status code: %d", resp.StatusCode)
    }
    obj := &ClusterMembership{}
    dec := json.NewDecoder(resp.Body)
    dec.Decode(obj)
    obj.service = s
    return obj, nil
}
func (s *ClusterMembershipsServiceOp) ListByHref(endpoint string, lo *ListOptions) ([]ClusterMembership, *ListParams, error) {
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
    obj := &ClusterMemberships{}
    dec := json.NewDecoder(resp.Body)
    dec.Decode(obj)

    dst := make([]ClusterMembership, len(obj.Items))
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

func (s *ClusterMembershipsServiceOp) Create(dir *ClusterMembership) (*ClusterMembership, error) {
    endpoint := fmt.Sprintf("clusterMemberships")

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
    obj := &ClusterMembership{}
    dec := json.NewDecoder(resp.Body)
    dec.Decode(obj)
    obj.service = s
    return obj, nil
}

// Delete removes application
func (s *ClusterMembershipsServiceOp) Delete(t *ClusterMembership) (error) {
    return s.DeleteByHref(t.Href)
}

// Delete removes application by ID
func (s *ClusterMembershipsServiceOp) DeleteById(id string) (error) {
    endpoint := fmt.Sprintf("clusterMemberships/%s", id)
    return s.DeleteByHref(endpoint)
}

// Delete removes application by link
func (s *ClusterMembershipsServiceOp) DeleteByHref(endpoint string) (error) {
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
