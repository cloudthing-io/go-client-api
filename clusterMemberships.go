package api

import (
    "bytes"
    "encoding/json" 
    "fmt"   
    "net/http"
    _"strings"
    "github.com/borystomala/copier"
)

// ClusterMembershipsService is an interface for interacting with ClusterMemberships endpoints of CloudThing API
// https://tenant-name.cloudthing.io/api/v1/clusterMemberships
type ClusterMembershipsService interface {
    GetById(string, ...interface{}) (*ClusterMembership, error)
    GetByLink(string, ...interface{}) (*ClusterMembership, error)
    ListByLink(string, ...interface{}) ([]ClusterMembership, *ListParams, error)
    ListByDevice(string, ...interface{}) ([]ClusterMembership, *ListParams, error)
    ListByCluster(string, ...interface{}) ([]ClusterMembership, *ListParams, error)
    CreateByLink(string, *ClusterMembershipRequestCreate) (*ClusterMembership, error)
    CreateByDevice(string, *ClusterMembershipRequestCreate) (*ClusterMembership, error)
    CreateByCluster(string, *ClusterMembershipRequestCreate) (*ClusterMembership, error)
    Delete(*ClusterMembership) (error)
    DeleteByLink(string) (error)
    DeleteById(string) (error)

    get(*ClusterMembershipResponse) (*ClusterMembership, error)
    getCollection(*ClusterMembershipsResponse) ([]ClusterMembership, *ListParams, error)
}

// ClusterMembershipsServiceOp handles communication with ClusterMemberships related methods of API
type ClusterMembershipsServiceOp struct {
    client *Client
}

// ClusterMembership is a struct representing CloudThing ClusterMembership
type ClusterMembership struct {
    // Standard field for all resources
    ModelBase

    Device          *Device 
    Cluster         *Cluster
    Application     *Application

    // Links to related resources
    device          string
    cluster         string
    application     string

    // service for communication, internal use only
    service         *ClusterMembershipsServiceOp
}

// ClusterMembershipResponse is a struct representing item response from API
type ClusterMembershipResponse struct {
    ModelBase

    Device          map[string]interface{}      `json:"device,omitempty"`
    Cluster         map[string]interface{}      `json:"cluster,omitempty"` 
    Application         map[string]interface{}  `json:"application,omitempty"` 
}

// ClusterMembershipResponse is a struct representing collection response from API
type ClusterMembershipsResponse struct{
    ListParams
    Items           []ClusterMembershipResponse   `json:"items"`
}

// ClusterMembershipResponse is a struct representing item create request for API
type ClusterMembershipRequestCreate struct {
    Device          *Link                       `json:"device,omitempty"`
    Cluster         *Link                       `json:"cluster,omitempty"` 
}

// TenantLink returns indicator of Tenant expansion and link to tenant.
// If expansion for Tenant was requested and resource is available via pointer
// it returns true, otherwise false. Link (href) is always returned. 
func (d *ClusterMembership) DeviceLink() (bool, string) {
    return (d.Device != nil), d.device
}

// DevicesLink returns indicator of Devices expansion and link to list of devices.
// If expansion for Devices was requested and resource is available via pointer
// it returns true, otherwise false. Link (href) is always returned. 
func (d *ClusterMembership) ClusterLink() (bool, string) {
    return (d.Cluster != nil), d.cluster
}

// DevicesLink returns indicator of Devices expansion and link to list of devices.
// If expansion for Devices was requested and resource is available via pointer
// it returns true, otherwise false. Link (href) is always returned. 
func (d *ClusterMembership) ApplicationLink() (bool, string) {
    return (d.Application != nil), d.application
}

// Delete is a helper method for deleting product.
// It calls Delete() on service under the hood.
func (t *ClusterMembership) Delete() error {
    err := t.service.Delete(t)
    if err != nil {
        return err
    }

    return nil
}

// GetById retrieves product by its ID
func (s *ClusterMembershipsServiceOp) GetById(id string, args ...interface{}) (*ClusterMembership, error) {
    endpoint := "clusterMemberships/"
    endpoint = fmt.Sprintf("%s%s", endpoint, id)

    return s.GetByLink(endpoint, args...)
}

// GetById retrieves product by its full link
func (s *ClusterMembershipsServiceOp) GetByLink(endpoint string, args ...interface{}) (*ClusterMembership, error) {
    resp, err := s.client.request("GET", endpoint, nil, args...)
    if err != nil {
        return nil, err
    }

    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, ApiError{StatusCode: resp.StatusCode, Message: "non-ok status returned"}
    }
    obj := &ClusterMembershipResponse{}
    dec := json.NewDecoder(resp.Body)
    dec.Decode(obj)

    return s.get(obj)
}

// get is internal method for transforming ClusterMembershipResponse into ClusterMembership
func (s *ClusterMembershipsServiceOp) get(r *ClusterMembershipResponse) (*ClusterMembership, error) {
    obj := &ClusterMembership{}
    copier.Copy(obj, r)
    if v, ok :=  r.Device["href"]; ok {
        obj.device = v.(string)
    }
    if v, ok :=  r.Application["href"]; ok {
        obj.application = v.(string)
    }
    if v, ok :=  r.Cluster["href"]; ok {
        obj.cluster = v.(string)
    }
    if len(r.Device) > 1 {        
        bytes, err := json.Marshal(r.Device)
        if err != nil {
            return nil, err
        }
        ten := &DeviceResponse{}
        json.Unmarshal(bytes, ten)
        t, err := s.client.Devices.get(ten)
        if err != nil {
            return nil, err
        }
        obj.Device = t
    }
    if len(r.Cluster) > 1 {        
        bytes, err := json.Marshal(r.Cluster)
        if err != nil {
            return nil, err
        }
        ten := &ClusterResponse{}
        json.Unmarshal(bytes, ten)
        t, err := s.client.Clusters.get(ten)
        if err != nil {
            return nil, err
        }
        obj.Cluster = t
    }
    if len(r.Application) > 1 {        
        bytes, err := json.Marshal(r.Application)
        if err != nil {
            return nil, err
        }
        ten := &ApplicationResponse{}
        json.Unmarshal(bytes, ten)
        t, err := s.client.Applications.get(ten)
        if err != nil {
            return nil, err
        }
        obj.Application = t
    }
    obj.service = s
    return obj, nil
}

// get is internal method for transforming ClusterMembershipResponse into ClusterMembership
func (s *ClusterMembershipsServiceOp) getCollection(r *ClusterMembershipsResponse) ([]ClusterMembership, *ListParams, error) {
    dst := make([]ClusterMembership, len(r.Items))

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

// GetById retrieves collection of clusterMemberships of current tenant
func (s *ClusterMembershipsServiceOp) ListByDevice(id string, args ...interface{}) ([]ClusterMembership, *ListParams, error) {
    endpoint := fmt.Sprintf("devices/%s/clusterMemberships", id)
    return s.ListByLink(endpoint, args...)
}

// GetById retrieves collection of clusterMemberships of current tenant
func (s *ClusterMembershipsServiceOp) ListByCluster(id string, args ...interface{}) ([]ClusterMembership, *ListParams, error) {
    endpoint := fmt.Sprintf("clusters/%s/memberships", id)
    return s.ListByLink(endpoint, args...)
}

// GetById retrieves collection of clusterMemberships by link
func (s *ClusterMembershipsServiceOp) ListByLink(endpoint string, args ...interface{}) ([]ClusterMembership, *ListParams, error) {
    resp, err := s.client.request("GET", endpoint, nil, args...)
    if err != nil {
        return nil, nil, err
    }

    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, nil, ApiError{StatusCode: resp.StatusCode, Message: "non-ok status returned"}
    }
    obj := &ClusterMembershipsResponse{}
    dec := json.NewDecoder(resp.Body)
    dec.Decode(obj)

    return s.getCollection(obj)
}

func (s *ClusterMembershipsServiceOp) CreateByDevice(id string, dir *ClusterMembershipRequestCreate) (*ClusterMembership, error) {
    endpoint := fmt.Sprintf("devices/%s/clusterMemberships", id)
    return s.CreateByLink(endpoint, dir)
}

func (s *ClusterMembershipsServiceOp) CreateByCluster(id string, dir *ClusterMembershipRequestCreate) (*ClusterMembership, error) {
    endpoint := fmt.Sprintf("clusters/%s/memberships", id)
    return s.CreateByLink(endpoint, dir)
}

// Create creates new product within tenant
func (s *ClusterMembershipsServiceOp) CreateByLink(endpoint string, dir *ClusterMembershipRequestCreate) (*ClusterMembership, error) {
    enc, err := json.Marshal(dir)
    if err != nil {
        return nil, err
    }

    buf := bytes.NewBuffer(enc)
    fmt.Println(buf.String())
    resp, err := s.client.request("POST", endpoint, buf)
    if err != nil {
        return nil, err
    }

    defer resp.Body.Close()

    if resp.StatusCode != http.StatusCreated {
        return nil, ApiError{StatusCode: resp.StatusCode, Message: "non-ok status returned"}
    }
    obj := &ClusterMembershipResponse{}
    dec := json.NewDecoder(resp.Body)
    dec.Decode(obj)
    return s.get(obj)
}

// Delete removes product
func (s *ClusterMembershipsServiceOp) Delete(t *ClusterMembership) (error) {
    return s.DeleteByLink(t.Href)
}

// Delete removes product by ID
func (s *ClusterMembershipsServiceOp) DeleteById(id string) (error) {
    endpoint := fmt.Sprintf("clusterMemberships/%s", id)
    return s.DeleteByLink(endpoint)
}

// Delete removes product by link
func (s *ClusterMembershipsServiceOp) DeleteByLink(endpoint string) (error) {
    resp, err := s.client.request("DELETE", endpoint, nil)
    if err != nil {
        return err
    }

    defer resp.Body.Close()

    if resp.StatusCode != http.StatusNoContent {
        return ApiError{StatusCode: resp.StatusCode, Message: "non-ok status returned"}
    }
    return nil
}
