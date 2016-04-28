package api

import (
    "bytes"
    "encoding/json" 
    "fmt"   
    "net/http"    
    "github.com/borystomala/copier"
)

// GroupsService is an interface for interacting with Groups endpoints of CloudThing API
// https://tenant-name.cloudthing.io/api/v1/groups
type GroupsService interface {
    GetById(string, ...interface{}) (*Group, error)
    GetByLink(string, ...interface{}) (*Group, error)
    ListByLink(string, ...interface{}) ([]Group, *ListParams, error)
    ListByCluster(string, ...interface{}) ([]Group, *ListParams, error)
    ListByDevice(string, ...interface{}) ([]Group, *ListParams, error)
    CreateByLink(string, *GroupRequestCreate) (*Group, error)
    CreateByCluster(string, *GroupRequestCreate) (*Group, error)
    UpdateById(string, *GroupRequestUpdate) (*Group, error)
    UpdateByLink(string, *GroupRequestUpdate) (*Group, error)
    Delete(*Group) (error)
    DeleteByLink(string) (error)
    DeleteById(string) (error)

    get(*GroupResponse) (*Group, error)
    getCollection(*GroupsResponse) ([]Group, *ListParams, error)
}

// GroupsServiceOp handles communication with Groups related methods of API
type GroupsServiceOp struct {
    client *Client
}

// Group is a struct representing CloudThing Group
type Group struct {
    // Standard field for all resources
    ModelBase
    Name                string
    Description         string 
    Custom              map[string]interface{}

    // Points to Tenant if expansion was requested, otherwise nil
    Tenant          *Tenant
    // Points to Applications if expansion was requested, otherwise nil
    Application     *Application
    // Points to Tenant if expansion was requested, otherwise nil
    Cluster         *Cluster
    // Points to Applications if expansion was requested, otherwise nil
    Devices         []Device
    Memberships     []GroupMembership

    // Links to related resources
    tenant          string
    application     string
    cluster         string
    devices         string
    memberships     string

    // service for communication, internal use only
    service         *GroupsServiceOp
}

// GroupResponse is a struct representing item response from API
type GroupResponse struct {
    ModelBase
    Name                string                  `json:"name,omitempty"`
    Description         string                  `json:"description,omitempty"`
    Custom              map[string]interface{}  `json:"custom,omitempty"`

    Tenant              map[string]interface{}  `json:"tenant,omitempty"`
    Application         map[string]interface{}  `json:"application,omitempty"`
    Cluster             map[string]interface{}  `json:"cluster,omitempty"`
    Devices             map[string]interface{}  `json:"devices,omitempty"`
    Memberships         map[string]interface{}  `json:"memberships,omitempty"`
}

// GroupResponse is a struct representing collection response from API
type GroupsResponse struct{
    ListParams
    Items           []GroupResponse              `json:"items"`
}

// GroupResponse is a struct representing item create request for API
type GroupRequestCreate struct {
    Name                string                  `json:"name,omitempty"`
    Description         string                  `json:"description,omitempty"`
    Custom              map[string]interface{}  `json:"custom,omitempty"`
}

// GroupResponse is a struct representing item update request for API
type GroupRequestUpdate struct {
    Name                string                  `json:"name,omitempty"`
    Description         string                  `json:"description,omitempty"`
    Custom              map[string]interface{}  `json:"custom,omitempty"`
}

// TenantLink returns indicator of Tenant expansion and link to tenant.
// If expansion for Tenant was requested and resource is available via pointer
// it returns true, otherwise false. Link (href) is always returned. 
func (d *Group) TenantLink() (bool, string) {
    return (d.Tenant != nil), d.tenant
}

// ApplicationsLink returns indicator of Group expansion and link to dierctory.
// If expansion for Group was requested and resource is available via pointer
// it returns true, otherwise false. Link (href) is always returned. 
func (d *Group) ApplicationLink() (bool, string) {
    return (d.Application != nil), d.application
}

// ApplicationsLink returns indicator of Group expansion and link to dierctory.
// If expansion for Group was requested and resource is available via pointer
// it returns true, otherwise false. Link (href) is always returned. 
func (d *Group) ClusterLink() (bool, string) {
    return (d.Cluster != nil), d.cluster
}

// ApplicationsLink returns indicator of Group expansion and link to dierctory.
// If expansion for Group was requested and resource is available via pointer
// it returns true, otherwise false. Link (href) is always returned. 
func (d *Group) DevicesLink() (bool, string) {
    return (d.Devices != nil), d.devices
}

// ApplicationsLink returns indicator of Group expansion and link to dierctory.
// If expansion for Group was requested and resource is available via pointer
// it returns true, otherwise false. Link (href) is always returned. 
func (d *Group) MembershipsLink() (bool, string) {
    return (d.Memberships != nil), d.memberships
}


// Save is a helper method for updating apikey.
// It calls UpdateByLink() on service under the hood.
func (t *Group) Save() error {
    tmp := &GroupRequestUpdate{}
    copier.Copy(tmp, t)
    ten, err := t.service.UpdateByLink(t.Href, tmp)
    if err != nil {
        return err
    }

    tmpTenant := t.Tenant
    tmpApplication := t.Application
    tmpDevices := t.Devices
    tmpCluster := t.Cluster

    *t = *ten
    t.Tenant = tmpTenant
    t.Application = tmpApplication
    t.Devices = tmpDevices
    t.Cluster = tmpCluster

    return nil
}

// Save is a helper method for deleting apikey.
// It calls Delete() on service under the hood.
func (t *Group) Delete() error {
    err := t.service.Delete(t)
    if err != nil {
        return err
    }

    return nil
}

// GetById retrieves apikey by its ID
func (s *GroupsServiceOp) GetById(id string, args ...interface{}) (*Group, error) {
    endpoint := "groups/"
    endpoint = fmt.Sprintf("%s%s", endpoint, id)

    return s.GetByLink(endpoint, args...)
}

// GetById retrieves apikey by its full link
func (s *GroupsServiceOp) GetByLink(endpoint string, args ...interface{}) (*Group, error) {
    resp, err := s.client.request("GET", endpoint, nil, args...)
    if err != nil {
        return nil, err
    }

    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("Status code: %d", resp.StatusCode)
    }
    obj := &GroupResponse{}
    dec := json.NewDecoder(resp.Body)
    dec.Decode(obj)

    return s.get(obj)
}

// get is internal method for transforming GroupResponse into Group
func (s *GroupsServiceOp) get(r *GroupResponse) (*Group, error) {
    obj := &Group{}
    copier.Copy(obj, r)
    if v, ok :=  r.Tenant["href"]; ok {
        obj.tenant = v.(string)
    }
    if v, ok :=  r.Application["href"]; ok {
        obj.application = v.(string)
    }
    if v, ok :=  r.Cluster["href"]; ok {
        obj.cluster = v.(string)
    }
    if v, ok :=  r.Devices["href"]; ok {
        obj.devices = v.(string)
    }
    if v, ok :=  r.Memberships["href"]; ok {
        obj.memberships = v.(string)
    }
    if len(r.Tenant) > 1 {        
        bytes, err := json.Marshal(r.Tenant)
        if err != nil {
            return nil, err
        }
        ten := &TenantResponse{}
        json.Unmarshal(bytes, ten)
        t, err := s.client.Tenant.get(ten)
        if err != nil {
            return nil, err
        }
        obj.Tenant = t
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
    if len(r.Devices) > 1 {        
        bytes, err := json.Marshal(r.Devices)
        if err != nil {
            return nil, err
        }
        ten := &DevicesResponse{}
        json.Unmarshal(bytes, ten)
        t, _, err := s.client.Devices.getCollection(ten)
        if err != nil {
            return nil, err
        }
        obj.Devices = t
    }
    if len(r.Memberships) > 1 {        
        bytes, err := json.Marshal(r.Memberships)
        if err != nil {
            return nil, err
        }
        ten := &GroupMembershipsResponse{}
        json.Unmarshal(bytes, ten)
        t, _, err := s.client.GroupMemberships.getCollection(ten)
        if err != nil {
            return nil, err
        }
        obj.Memberships = t
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
    obj.service = s
    return obj, nil
}

// get is internal method for transforming ApplicationResponse into Application
func (s *GroupsServiceOp) getCollection(r *GroupsResponse) ([]Group, *ListParams, error) {
    dst := make([]Group, len(r.Items))

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

// GetById retrieves collection of groups of current tenant
func (s *GroupsServiceOp) ListByCluster(id string, args ...interface{}) ([]Group, *ListParams, error) {
    endpoint := fmt.Sprintf("clusters/%s/groups", id)
    return s.ListByLink(endpoint, args...)
}

// GetById retrieves collection of groups of current tenant
func (s *GroupsServiceOp) ListByDevice(id string, args ...interface{}) ([]Group, *ListParams, error) {
    endpoint := fmt.Sprintf("devices/%s/groups", id)
    return s.ListByLink(endpoint, args...)
}

// GetById retrieves collection of groups by link
func (s *GroupsServiceOp) ListByLink(endpoint string, args ...interface{}) ([]Group, *ListParams, error) {
    resp, err := s.client.request("GET", endpoint, nil, args...)
    if err != nil {
        return nil, nil, err
    }

    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, nil, fmt.Errorf("Status code: %d", resp.StatusCode)
    }
    obj := &GroupsResponse{}
    dec := json.NewDecoder(resp.Body)
    dec.Decode(obj)

    return s.getCollection(obj)
}

// GetById updates apikey with specified ID
func (s *GroupsServiceOp) UpdateById(id string, t *GroupRequestUpdate) (*Group, error) {
    endpoint := fmt.Sprintf("groups/%s", id)
    return s.UpdateByLink(endpoint, t)
}

// GetById updates apikey specified by link
func (s *GroupsServiceOp) UpdateByLink(endpoint string, t *GroupRequestUpdate) (*Group, error) {
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
    obj := &GroupResponse{}
    dec := json.NewDecoder(resp.Body)
    dec.Decode(obj)
    return s.get(obj)
}

func (s *GroupsServiceOp) CreateByCluster(id string, dir *GroupRequestCreate) (*Group, error) {
    endpoint := fmt.Sprintf("clusters/%s/groups", id)
    return s.CreateByLink(endpoint, dir)
}

// Create creates new apikey within tenant
func (s *GroupsServiceOp) CreateByLink(endpoint string, dir *GroupRequestCreate) (*Group, error) {
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
    obj := &GroupResponse{}
    dec := json.NewDecoder(resp.Body)
    dec.Decode(obj)
    return s.get(obj)
}

// Delete removes apikey
func (s *GroupsServiceOp) Delete(t *Group) (error) {
    return s.DeleteByLink(t.Href)
}

// Delete removes apikey by ID
func (s *GroupsServiceOp) DeleteById(id string) (error) {
    endpoint := fmt.Sprintf("groups/%s", id)
    return s.DeleteByLink(endpoint)
}

// Delete removes apikey by link
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
