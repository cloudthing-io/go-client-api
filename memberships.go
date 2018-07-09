package api

import (
    "bytes"
    "encoding/json" 
    "fmt"   
    "net/http"
    _"strings"
    "github.com/borystomala/copier"
)

// MembershipsService is an interface for interacting with Memberships endpoints of CloudThing API
// https://tenant-name.cloudthing.io/api/v1/memberships
type MembershipsService interface {
    GetById(string, ...interface{}) (*Membership, error)
    GetByLink(string, ...interface{}) (*Membership, error)
    ListByLink(string, ...interface{}) ([]Membership, *ListParams, error)
    ListByUser(string, ...interface{}) ([]Membership, *ListParams, error)
    ListByUsergroup(string, ...interface{}) ([]Membership, *ListParams, error)
    CreateByLink(string, *MembershipRequestCreate) (*Membership, error)
    CreateByUser(string, *MembershipRequestCreate) (*Membership, error)
    CreateByUsergroup(string, *MembershipRequestCreate) (*Membership, error)
    Delete(*Membership) (error)
    DeleteByLink(string) (error)
    DeleteById(string) (error)

    get(*MembershipResponse) (*Membership, error)
    getCollection(*MembershipsResponse) ([]Membership, *ListParams, error)
}

// MembershipsServiceOp handles communication with Memberships related methods of API
type MembershipsServiceOp struct {
    client *Client
}

// Membership is a struct representing CloudThing Membership
type Membership struct {
    // Standard field for all resources
    ModelBase

    User            *User 
    Usergroup       *Usergroup

    // Links to related resources
    user            string
    usergroup       string

    // service for communication, internal use only
    service         *MembershipsServiceOp
}

// MembershipResponse is a struct representing item response from API
type MembershipResponse struct {
    ModelBase

    User            map[string]interface{}  `json:"user,omitempty"`
    Usergroup       map[string]interface{}  `json:"usergroup,omitempty"` 
}

// MembershipResponse is a struct representing collection response from API
type MembershipsResponse struct{
    ListParams
    Items           []MembershipResponse   `json:"items"`
}

// MembershipResponse is a struct representing item create request for API
type MembershipRequestCreate struct {
    User            *Link                   `json:"user,omitempty"`
    Usergroup       *Link                   `json:"usergroup,omitempty"` 
}

// TenantLink returns indicator of Tenant expansion and link to tenant.
// If expansion for Tenant was requested and resource is available via pointer
// it returns true, otherwise false. Link (href) is always returned. 
func (d *Membership) UserLink() (bool, string) {
    return (d.User != nil), d.user
}

// DevicesLink returns indicator of Devices expansion and link to list of devices.
// If expansion for Devices was requested and resource is available via pointer
// it returns true, otherwise false. Link (href) is always returned. 
func (d *Membership) UsergroupLink() (bool, string) {
    return (d.Usergroup != nil), d.usergroup
}

// Delete is a helper method for deleting product.
// It calls Delete() on service under the hood.
func (t *Membership) Delete() error {
    err := t.service.Delete(t)
    if err != nil {
        return err
    }

    return nil
}

// GetById retrieves product by its ID
func (s *MembershipsServiceOp) GetById(id string, args ...interface{}) (*Membership, error) {
    endpoint := "memberships/"
    endpoint = fmt.Sprintf("%s%s", endpoint, id)

    return s.GetByLink(endpoint, args...)
}

// GetById retrieves product by its full link
func (s *MembershipsServiceOp) GetByLink(endpoint string, args ...interface{}) (*Membership, error) {
    resp, err := s.client.request("GET", endpoint, nil, args...)
    if err != nil {
        return nil, err
    }

    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, ApiError{StatusCode: resp.StatusCode, Message: "non-ok status returned"}
    }
    obj := &MembershipResponse{}
    dec := json.NewDecoder(resp.Body)
    dec.Decode(obj)

    return s.get(obj)
}

// get is internal method for transforming MembershipResponse into Membership
func (s *MembershipsServiceOp) get(r *MembershipResponse) (*Membership, error) {
    obj := &Membership{}
    copier.Copy(obj, r)
    if v, ok :=  r.User["href"]; ok {
        obj.user = v.(string)
    }
    if v, ok :=  r.Usergroup["href"]; ok {
        obj.usergroup = v.(string)
    }
    if len(r.User) > 1 {        
        bytes, err := json.Marshal(r.User)
        if err != nil {
            return nil, err
        }
        ten := &UserResponse{}
        json.Unmarshal(bytes, ten)
        t, err := s.client.Users.get(ten)
        if err != nil {
            return nil, err
        }
        obj.User = t
    }
    if len(r.Usergroup) > 1 {        
        bytes, err := json.Marshal(r.Usergroup)
        if err != nil {
            return nil, err
        }
        ten := &UsergroupResponse{}
        json.Unmarshal(bytes, ten)
        t, err := s.client.Usergroups.get(ten)
        if err != nil {
            return nil, err
        }
        obj.Usergroup = t
    }
    obj.service = s
    return obj, nil
}

// get is internal method for transforming MembershipResponse into Membership
func (s *MembershipsServiceOp) getCollection(r *MembershipsResponse) ([]Membership, *ListParams, error) {
    dst := make([]Membership, len(r.Items))

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

// GetById retrieves collection of memberships of current tenant
func (s *MembershipsServiceOp) ListByUser(id string, args ...interface{}) ([]Membership, *ListParams, error) {
    endpoint := fmt.Sprintf("users/%s/memberships", id)
    return s.ListByLink(endpoint, args...)
}

// GetById retrieves collection of memberships of current tenant
func (s *MembershipsServiceOp) ListByUsergroup(id string, args ...interface{}) ([]Membership, *ListParams, error) {
    endpoint := fmt.Sprintf("usergroups/%s/memberships", id)
    return s.ListByLink(endpoint, args...)
}

// GetById retrieves collection of memberships by link
func (s *MembershipsServiceOp) ListByLink(endpoint string, args ...interface{}) ([]Membership, *ListParams, error) {
    resp, err := s.client.request("GET", endpoint, nil, args...)
    if err != nil {
        return nil, nil, err
    }

    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, nil, ApiError{StatusCode: resp.StatusCode, Message: "non-ok status returned"}
    }
    obj := &MembershipsResponse{}
    dec := json.NewDecoder(resp.Body)
    dec.Decode(obj)

    return s.getCollection(obj)
}

func (s *MembershipsServiceOp) CreateByUser(id string, dir *MembershipRequestCreate) (*Membership, error) {
    endpoint := fmt.Sprintf("users/%s/memberships", id)
    return s.CreateByLink(endpoint, dir)
}

func (s *MembershipsServiceOp) CreateByUsergroup(id string, dir *MembershipRequestCreate) (*Membership, error) {
    endpoint := fmt.Sprintf("usergroups/%s/memberships", id)
    return s.CreateByLink(endpoint, dir)
}

// Create creates new product within tenant
func (s *MembershipsServiceOp) CreateByLink(endpoint string, dir *MembershipRequestCreate) (*Membership, error) {
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
    obj := &MembershipResponse{}
    dec := json.NewDecoder(resp.Body)
    dec.Decode(obj)
    return s.get(obj)
}

// Delete removes product
func (s *MembershipsServiceOp) Delete(t *Membership) (error) {
    return s.DeleteByLink(t.Href)
}

// Delete removes product by ID
func (s *MembershipsServiceOp) DeleteById(id string) (error) {
    endpoint := fmt.Sprintf("memberships/%s", id)
    return s.DeleteByLink(endpoint)
}

// Delete removes product by link
func (s *MembershipsServiceOp) DeleteByLink(endpoint string) (error) {
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
