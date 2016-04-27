package api

import (
    "bytes"
    "encoding/json" 
    "fmt"   
    "net/http"
)

// GroupMembershipsService is an interafce for interacting with GroupMemberships endpoints of CloudThing API
// https://tenant-name.cloudthing.io/api/v1/groupMemberships

type GroupMembershipsService interface {
    GetById(string) (*GroupMembership, error)
    GetByHref(string) (*GroupMembership, error)
    ListByHref(string, *ListOptions) ([]GroupMembership, *ListParams, error)
    Create(*GroupMembership) (*GroupMembership, error)
    Delete(*GroupMembership) (error)
    DeleteByHref(string) (error)
    DeleteById(string) (error)
}

// GroupMembershipsServiceOpg handles communication with Tenant related methods of API
type GroupMembershipsServiceOp struct {
    client *Client
}

type GroupMembership struct {
    ModelBase

    device          string          `json:"device,omitempty"`
    group           string          `json:"group,omitempty"`

    // service for communication, internal use only
    service         *GroupMembershipsServiceOp `json:"-"` 
}

type GroupMemberships struct{
    ListParams
    Items           []GroupMembership     `json:"items"`
}
/*
func (d *GroupMembership) Application() (*Application, error) {
    return d.service.client.Applications.GetByHref(d.application)
}*/

func (d *GroupMembership) Device() (*Device, error) {
    return d.service.client.Devices.GetByHref(d.device)
}

func (d *GroupMembership) Group() (*Group, error) {
    return d.service.client.Groups.GetByHref(d.group)
}

// GetById retrieves directory
func (s *GroupMembershipsServiceOp) GetById(id string) (*GroupMembership, error) {
    endpoint := "groupMemberships/"
    endpoint = fmt.Sprintf("%s%s", endpoint, id)

    return s.GetByHref(endpoint)
}

func (s *GroupMembershipsServiceOp) GetByHref(endpoint string) (*GroupMembership, error) {
    resp, err := s.client.request("GET", endpoint, nil)
    if err != nil {
        return nil, err
    }

    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("Status code: %d", resp.StatusCode)
    }
    obj := &GroupMembership{}
    dec := json.NewDecoder(resp.Body)
    dec.Decode(obj)
    obj.service = s
    return obj, nil
}
func (s *GroupMembershipsServiceOp) ListByHref(endpoint string, lo *ListOptions) ([]GroupMembership, *ListParams, error) {
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
    obj := &GroupMemberships{}
    dec := json.NewDecoder(resp.Body)
    dec.Decode(obj)

    dst := make([]GroupMembership, len(obj.Items))
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


func (s *GroupMembershipsServiceOp) Create(dir *GroupMembership) (*GroupMembership, error) {
    endpoint := fmt.Sprintf("groupMemberships")

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
    obj := &GroupMembership{}
    dec := json.NewDecoder(resp.Body)
    dec.Decode(obj)
    obj.service = s
    return obj, nil
}

// Delete removes application
func (s *GroupMembershipsServiceOp) Delete(t *GroupMembership) (error) {
    return s.DeleteByHref(t.Href)
}

// Delete removes application by ID
func (s *GroupMembershipsServiceOp) DeleteById(id string) (error) {
    endpoint := fmt.Sprintf("groupMemberships/%s", id)
    return s.DeleteByHref(endpoint)
}

// Delete removes application by link
func (s *GroupMembershipsServiceOp) DeleteByHref(endpoint string) (error) {
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
