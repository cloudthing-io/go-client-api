package api

import (
    "bytes"
    "encoding/json" 
    "fmt"   
    "net/http"
)

// ProductsService is an interafce for interacting with Products endpoints of CloudThing API
// https://tenant-name.cloudthing.io/api/v1/products

type ProductsService interface {
    GetById(string) (*Product, error)
    GetByLink(string) (*Product, error)
    List(*ListOptions) ([]Product, *ListParams, error)
    Create(*Product) (*Product, error)
    Update(*Product) (*Product, error)
    Delete(*Product) (error)
    DeleteByLink(string) (error)
    DeleteById(string) (error)
}

// ProductsServiceOp handles communication with Tenant related methods of API
type ProductsServiceOp struct {
    client *Client
}

type ProductProperty struct { // oops, it's device property, change needed
    Key             string          `json:"key"`
    Value           interface{}     `json:"value"`
}

type ProductSimpleResource struct {
    Id              string          `json:"id"`
    Name            string          `json:"name"`
    Description     string          `json:"description"`
}

type ProductPayload struct {
    Name            string          `json:"name"`
    Serialization   string          `json:"serialization"`
    Value           string          `json:"value"`
}

type ProductCommandResource struct {
    ProductSimpleResource
    Payloads        []ProductPayload`json:"payloads"`
}

type ProductResources struct {
    Data            []ProductSimpleResource `json:"data"`
    Events          []ProductSimpleResource `json:"events"`
    Commands        []ProductCommandResource`json:"commands"`
}

type Product struct {
    ModelBase
    Name            string          `json:"name,omitempty"`
    Description     string          `json:"description,omitempty"`
    Custom          interface{}     `json:"custom,omitempty"`
    Properties      []ProductProperty `json:"properties,omitempty"`
    Resources       *ProductResources`json:"resources,omitempty"`

    tenant          string          `json:"tenant,omitempty"`
    devices         string          `json:"users,omitempty"`

    // service for communication, internal use only
    service         *ProductsServiceOp `json:"-"` 
}

type Products struct{
    ListParams
    Items           []Product     `json:"items"`
}

func (d *Product) Tenant() (*Tenant, error) {
    return d.service.client.Tenant.Get()
}

func (d *Product) Devices() ([]Device, *ListParams, error) {
    return d.service.client.Devices.ListByLink(d.devices, nil)
}

// Save updates tenant by calling Update() on service under the hood
func (t *Product) Save() error {
    tmp := *t

    ten, err := t.service.Update(&tmp)
    if err != nil {
        return err
    }

    *t = *ten
    return nil
}

// GetById retrieves directory
func (s *ProductsServiceOp) GetById(id string) (*Product, error) {
    endpoint := "products/"
    endpoint = fmt.Sprintf("%s%s", endpoint, id)

    return s.GetByLink(endpoint)
}

func (s *ProductsServiceOp) GetByLink(id string) (*Product, error) {
    resp, err := s.client.request("GET", id, nil)
    if err != nil {
        return nil, err
    }

    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("Status code: %d", resp.StatusCode)
    }
    obj := &Product{}
    dec := json.NewDecoder(resp.Body)
    dec.Decode(obj)
    obj.service = s
    return obj, nil
}

func (s *ProductsServiceOp) List(lo *ListOptions) ([]Product, *ListParams, error) {
    if lo == nil {
        lo = &ListOptions {
            Page: 1,
            Limit: DefaultLimit,
        }
    }
    endpoint := fmt.Sprintf("tenants/%s/products?limit=%d&page=%d", s.client.tenantId, lo.Limit, lo.Page)

    resp, err := s.client.request("GET", endpoint, nil)
    if err != nil {
        return nil, nil, err
    }

    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, nil, fmt.Errorf("Status code: %d", resp.StatusCode)
    }
    obj := &Products{}
    dec := json.NewDecoder(resp.Body)
    dec.Decode(obj)

    dst := make([]Product, len(obj.Items))
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
func (s *ProductsServiceOp) Update(t *Product) (*Product, error) {
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
    obj := &Product{}
    dec := json.NewDecoder(resp.Body)
    dec.Decode(obj)
    obj.service = s
    return obj, nil
}

func (s *ProductsServiceOp) Create(dir *Product) (*Product, error) {
    endpoint := fmt.Sprintf("products")

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
    obj := &Product{}
    dec := json.NewDecoder(resp.Body)
    dec.Decode(obj)
    obj.service = s
    return obj, nil
}

// Delete removes application
func (s *ProductsServiceOp) Delete(t *Product) (error) {
    return s.DeleteByLink(t.Href)
}

// Delete removes application by ID
func (s *ProductsServiceOp) DeleteById(id string) (error) {
    endpoint := fmt.Sprintf("products/%s", id)
    return s.DeleteByLink(endpoint)
}

// Delete removes application by link
func (s *ProductsServiceOp) DeleteByLink(endpoint string) (error) {
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
