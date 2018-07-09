package api

import (
    "bytes"
    "encoding/json" 
    "fmt"   
    "net/http"
    _"strings"
    "github.com/borystomala/copier"
)

// GroupMembershipsService is an interface for interacting with GroupMemberships endpoints of CloudThing API
// https://tenant-name.cloudthing.io/api/v1/groupMemberships
type GroupMembershipsService interface {
    GetById(string, ...interface{}) (*GroupMembership, error)
    GetByLink(string, ...interface{}) (*GroupMembership, error)
    ListByLink(string, ...interface{}) ([]GroupMembership, *ListParams, error)
    ListByDevice(string, ...interface{}) ([]GroupMembership, *ListParams, error)
    ListByGroup(string, ...interface{}) ([]GroupMembership, *ListParams, error)
    CreateByLink(string, *GroupMembershipRequestCreate) (*GroupMembership, error)
    CreateByDevice(string, *GroupMembershipRequestCreate) (*GroupMembership, error)
    CreateByGroup(string, *GroupMembershipRequestCreate) (*GroupMembership, error)
    Delete(*GroupMembership) (error)
    DeleteByLink(string) (error)
    DeleteById(string) (error)

    get(*GroupMembershipResponse) (*GroupMembership, error)
    getCollection(*GroupMembershipsResponse) ([]GroupMembership, *ListParams, error)
}

// GroupMembershipsServiceOp handles communication with GroupMemberships related methods of API
type GroupMembershipsServiceOp struct {
    client *Client
}

// GroupMembership is a struct representing CloudThing GroupMembership
type GroupMembership struct {
    // Standard field for all resources
    ModelBase

    Device          *Device 
    Group           *Group

    // Links to related resources
    device          string
    group           string

    // service for communication, internal use only
    service         *GroupMembershipsServiceOp
}

// GroupMembershipResponse is a struct representing item response from API
type GroupMembershipResponse struct {
    ModelBase

    Device          map[string]interface{}      `json:"device,omitempty"`
    Group           map[string]interface{}      `json:"group,omitempty"` 
}

// GroupMembershipResponse is a struct representing collection response from API
type GroupMembershipsResponse struct{
    ListParams
    Items           []GroupMembershipResponse   `json:"items"`
}

// GroupMembershipResponse is a struct representing item create request for API
type GroupMembershipRequestCreate struct {
    Device          *Link                       `json:"device,omitempty"`
    Group           *Link                       `json:"group,omitempty"` 
}

// TenantLink returns indicator of Tenant expansion and link to tenant.
// If expansion for Tenant was requested and resource is available via pointer
// it returns true, otherwise false. Link (href) is always returned. 
func (d *GroupMembership) DeviceLink() (bool, string) {
    return (d.Device != nil), d.device
}

// DevicesLink returns indicator of Devices expansion and link to list of devices.
// If expansion for Devices was requested and resource is available via pointer
// it returns true, otherwise false. Link (href) is always returned. 
func (d *GroupMembership) GroupLink() (bool, string) {
    return (d.Group != nil), d.group
}

// Delete is a helper method for deleting product.
// It calls Delete() on service under the hood.
func (t *GroupMembership) Delete() error {
    err := t.service.Delete(t)
    if err != nil {
        return err
    }

    return nil
}

// GetById retrieves product by its ID
func (s *GroupMembershipsServiceOp) GetById(id string, args ...interface{}) (*GroupMembership, error) {
    endpoint := "groupMemberships/"
    endpoint = fmt.Sprintf("%s%s", endpoint, id)

    return s.GetByLink(endpoint, args...)
}

// GetById retrieves product by its full link
func (s *GroupMembershipsServiceOp) GetByLink(endpoint string, args ...interface{}) (*GroupMembership, error) {
    resp, err := s.client.request("GET", endpoint, nil, args...)
    if err != nil {
        return nil, err
    }

    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, ApiError{StatusCode: resp.StatusCode, Message: "non-ok status returned"}
    }
    obj := &GroupMembershipResponse{}
    dec := json.NewDecoder(resp.Body)
    dec.Decode(obj)

    return s.get(obj)
}

// get is internal method for transforming GroupMembershipResponse into GroupMembership
func (s *GroupMembershipsServiceOp) get(r *GroupMembershipResponse) (*GroupMembership, error) {
    obj := &GroupMembership{}
    copier.Copy(obj, r)
    if v, ok :=  r.Device["href"]; ok {
        obj.device = v.(string)
    }
    if v, ok :=  r.Group["href"]; ok {
        obj.group = v.(string)
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
    if len(r.Group) > 1 {        
        bytes, err := json.Marshal(r.Group)
        if err != nil {
            return nil, err
        }
        ten := &GroupResponse{}
        json.Unmarshal(bytes, ten)
        t, err := s.client.Groups.get(ten)
        if err != nil {
            return nil, err
        }
        obj.Group = t
    }
    obj.service = s
    return obj, nil
}

// get is internal method for transforming GroupMembershipResponse into GroupMembership
func (s *GroupMembershipsServiceOp) getCollection(r *GroupMembershipsResponse) ([]GroupMembership, *ListParams, error) {
    dst := make([]GroupMembership, len(r.Items))

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

// GetById retrieves collection of groupMemberships of current tenant
func (s *GroupMembershipsServiceOp) ListByDevice(id string, args ...interface{}) ([]GroupMembership, *ListParams, error) {
    endpoint := fmt.Sprintf("devices/%s/groupMemberships", id)
    return s.ListByLink(endpoint, args...)
}

// GetById retrieves collection of groupMemberships of current tenant
func (s *GroupMembershipsServiceOp) ListByGroup(id string, args ...interface{}) ([]GroupMembership, *ListParams, error) {
    endpoint := fmt.Sprintf("groups/%s/groupMemberships", id)
    return s.ListByLink(endpoint, args...)
}

// GetById retrieves collection of groupMemberships by link
func (s *GroupMembershipsServiceOp) ListByLink(endpoint string, args ...interface{}) ([]GroupMembership, *ListParams, error) {
    resp, err := s.client.request("GET", endpoint, nil, args...)
    if err != nil {
        return nil, nil, err
    }

    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, nil, ApiError{StatusCode: resp.StatusCode, Message: "non-ok status returned"}
    }
    obj := &GroupMembershipsResponse{}
    dec := json.NewDecoder(resp.Body)
    dec.Decode(obj)

    return s.getCollection(obj)
}

func (s *GroupMembershipsServiceOp) CreateByDevice(id string, dir *GroupMembershipRequestCreate) (*GroupMembership, error) {
    endpoint := fmt.Sprintf("devices/%s/groupMemberships", id)
    return s.CreateByLink(endpoint, dir)
}

func (s *GroupMembershipsServiceOp) CreateByGroup(id string, dir *GroupMembershipRequestCreate) (*GroupMembership, error) {
    endpoint := fmt.Sprintf("groups/%s/groupMemberships", id)
    return s.CreateByLink(endpoint, dir)
}

// Create creates new product within tenant
func (s *GroupMembershipsServiceOp) CreateByLink(endpoint string, dir *GroupMembershipRequestCreate) (*GroupMembership, error) {
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
        return nil, ApiError{StatusCode: resp.StatusCode, Message: "non-ok status returned"}
    }
    obj := &GroupMembershipResponse{}
    dec := json.NewDecoder(resp.Body)
    dec.Decode(obj)
    return s.get(obj)
}

// Delete removes product
func (s *GroupMembershipsServiceOp) Delete(t *GroupMembership) (error) {
    return s.DeleteByLink(t.Href)
}

// Delete removes product by ID
func (s *GroupMembershipsServiceOp) DeleteById(id string) (error) {
    endpoint := fmt.Sprintf("groupMemberships/%s", id)
    return s.DeleteByLink(endpoint)
}

// Delete removes product by link
func (s *GroupMembershipsServiceOp) DeleteByLink(endpoint string) (error) {
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
