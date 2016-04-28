package api

import (
    "bytes"
    "encoding/json" 
    "fmt"   
    "net/http"
)

// MembershipsService is an interafce for interacting with Memberships endpoints of CloudThing API
// https://tenant-name.cloudthing.io/api/v1/memberships

type MembershipsService interface {
    GetById(string) (*Membership, error)
    GetByLink(string) (*Membership, error)
    ListByLink(string, *ListOptions) ([]Membership, *ListParams, error)
    Create(*Membership) (*Membership, error)
    Delete(*Membership) (error)
    DeleteByLink(string) (error)
    DeleteById(string) (error)
}

// MembershipsServiceOp handles communication with Tenant related methods of API
type MembershipsServiceOp struct {
    client *Client
}

type Membership struct {
    ModelBase

    user            string          `json:"user,omitempty"`
    usergroup       string          `json:"usergroup,omitempty"`

    // service for communication, internal use only
    service         *MembershipsServiceOp `json:"-"` 
}

type Memberships struct{
    ListParams
    Items           []Membership     `json:"items"`
}

func (d *Membership) User() (*User, error) {
    return d.service.client.Users.GetByLink(d.user)
}

func (d *Membership) Usergroup() (*Usergroup, error) {
    return d.service.client.Usergroups.GetByLink(d.usergroup)
}

// GetById retrieves directory
func (s *MembershipsServiceOp) GetById(id string) (*Membership, error) {
    endpoint := "memberships/"
    endpoint = fmt.Sprintf("%s%s", endpoint, id)

    return s.GetByLink(endpoint)
}

func (s *MembershipsServiceOp) GetByLink(endpoint string) (*Membership, error) {
    resp, err := s.client.request("GET", endpoint, nil)
    if err != nil {
        return nil, err
    }

    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("Status code: %d", resp.StatusCode)
    }
    obj := &Membership{}
    dec := json.NewDecoder(resp.Body)
    dec.Decode(obj)
    obj.service = s
    return obj, nil
}
func (s *MembershipsServiceOp) ListByLink(endpoint string, lo *ListOptions) ([]Membership, *ListParams, error) {
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
    obj := &Memberships{}
    dec := json.NewDecoder(resp.Body)
    dec.Decode(obj)

    dst := make([]Membership, len(obj.Items))
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

func (s *MembershipsServiceOp) Create(dir *Membership) (*Membership, error) {
    endpoint := fmt.Sprintf("memberships")

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
    obj := &Membership{}
    dec := json.NewDecoder(resp.Body)
    dec.Decode(obj)
    obj.service = s
    return obj, nil
}

// Delete removes application
func (s *MembershipsServiceOp) Delete(t *Membership) (error) {
    return s.DeleteByLink(t.Href)
}

// Delete removes application by ID
func (s *MembershipsServiceOp) DeleteById(id string) (error) {
    endpoint := fmt.Sprintf("memberships/%s", id)
    return s.DeleteByLink(endpoint)
}

// Delete removes application by link
func (s *MembershipsServiceOp) DeleteByLink(endpoint string) (error) {
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
