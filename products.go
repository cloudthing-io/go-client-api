package api

import (
    "bytes"
    "encoding/json" 
    "fmt"   
    "net/http"
    _"strings"
    "github.com/borystomala/copier"
)

// ProductsService is an interface for interacting with Products endpoints of CloudThing API
// https://tenant-name.cloudthing.io/api/v1/products
type ProductsService interface {
    GetById(string, ...interface{}) (*Product, error)
    GetByLink(string, ...interface{}) (*Product, error)
    List(...interface{}) ([]Product, *ListParams, error)
    Create(*ProductRequestCreate) (*Product, error)
    UpdateById(string, *ProductRequestUpdate) (*Product, error)
    UpdateByLink(string, *ProductRequestUpdate) (*Product, error)
    Delete(*Product) (error)
    DeleteByLink(string) (error)
    DeleteById(string) (error)

    get(*ProductResponse) (*Product, error)
    getCollection(*ProductsResponse) ([]Product, *ListParams, error)
}

// ProductsServiceOp handles communication with Products related methods of API
type ProductsServiceOp struct {
    client *Client
}

// Product is a struct representing CloudThing Product
type Product struct {
    // Standard field for all resources
    ModelBase
    // Product name
    Name            string
    // Description of product
    Description     string
    // Field for tenant's custom data
    Custom          map[string]interface{}
    Properties      []ProductProperty
    Resources       *ProductResources

    // Points to Tenant if expansion was requested, otherwise nil
    Tenant          *Tenant
    // Points to Devices if expansion was requested, otherwise nil
    Devices         []Device

    // Links to related resources
    tenant          string
    devices         string

    // service for communication, internal use only
    service         *ProductsServiceOp
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

// ProductResponse is a struct representing item response from API
type ProductResponse struct {
    ModelBase
    Name            string                  `json:"name,omitempty"`
    Description     string                  `json:"description,omitempty"`
    Custom          map[string]interface{}  `json:"custom,omitempty"`
    Properties      []ProductProperty       `json:"properties,omitempty"`
    Resources       *ProductResources       `json:"resources,omitempty"`

    Tenant          map[string]interface{}  `json:"tenant,omitempty"`
    Devices         map[string]interface{}  `json:"devices,omitempty"` 
}

// ProductResponse is a struct representing collection response from API
type ProductsResponse struct{
    ListParams
    Items           []ProductResponse   `json:"items"`
}

// ProductResponse is a struct representing item create request for API
type ProductRequestCreate struct {
    Name            string                  `json:"name,omitempty"`
    Description     string                  `json:"description,omitempty"`
    Custom          map[string]interface{}  `json:"custom,omitempty"`
    Properties      []ProductProperty       `json:"properties,omitempty"`
    Resources       *ProductResources       `json:"resources,omitempty"`
}

// ProductResponse is a struct representing item update request for API
type ProductRequestUpdate struct {
    Name            string                  `json:"name,omitempty"`
    Description     string                  `json:"description,omitempty"`
    Custom          map[string]interface{}  `json:"custom,omitempty"`
    Properties      []ProductProperty       `json:"properties,omitempty"`
    Resources       *ProductResources       `json:"resources,omitempty"`
}

// TenantLink returns indicator of Tenant expansion and link to tenant.
// If expansion for Tenant was requested and resource is available via pointer
// it returns true, otherwise false. Link (href) is always returned. 
func (d *Product) TenantLink() (bool, string) {
    return (d.Tenant != nil), d.tenant
}

// DevicesLink returns indicator of Devices expansion and link to list of devices.
// If expansion for Devices was requested and resource is available via pointer
// it returns true, otherwise false. Link (href) is always returned. 
func (d *Product) DevicesLink() (bool, string) {
    return (d.Devices != nil), d.devices
}

// Save is a helper method for updating product.
// It calls UpdateByLink() on service under the hood.
func (t *Product) Save() error {
    tmp := &ProductRequestUpdate{}
    copier.Copy(tmp, t)
    ten, err := t.service.UpdateByLink(t.Href, tmp)
    if err != nil {
        return err
    }

    tmpTenant := t.Tenant
    tmpDevices := t.Devices

    *t = *ten
    t.Tenant = tmpTenant
    t.Devices = tmpDevices

    return nil
}

// Save is a helper method for deleting product.
// It calls Delete() on service under the hood.
func (t *Product) Delete() error {
    err := t.service.Delete(t)
    if err != nil {
        return err
    }

    return nil
}

// GetById retrieves product by its ID
func (s *ProductsServiceOp) GetById(id string, args ...interface{}) (*Product, error) {
    endpoint := "products/"
    endpoint = fmt.Sprintf("%s%s", endpoint, id)

    return s.GetByLink(endpoint, args...)
}

// GetById retrieves product by its full link
func (s *ProductsServiceOp) GetByLink(endpoint string, args ...interface{}) (*Product, error) {
    resp, err := s.client.request("GET", endpoint, nil, args...)
    if err != nil {
        return nil, err
    }

    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("Status code: %d", resp.StatusCode)
    }
    obj := &ProductResponse{}
    dec := json.NewDecoder(resp.Body)
    dec.Decode(obj)

    return s.get(obj)
}

// get is internal method for transforming ProductResponse into Product
func (s *ProductsServiceOp) get(r *ProductResponse) (*Product, error) {
    obj := &Product{}
    copier.Copy(obj, r)
    if v, ok :=  r.Tenant["href"]; ok {
        obj.tenant = v.(string)
    }
    if v, ok :=  r.Devices["href"]; ok {
        obj.devices = v.(string)
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
 /*   if len(r.Directory) > 1 {        
        bytes, err := json.Marshal(r.Directory)
        if err != nil {
            return nil, err
        }
        ten := &DirectoryResponse{}
        json.Unmarshal(bytes, ten)
        t, err := s.client.Directories.get(ten)
        if err != nil {
            return nil, err
        }
        obj.Directory = t
    }*/
    obj.service = s
    return obj, nil
}

// get is internal method for transforming ProductResponse into Product
func (s *ProductsServiceOp) getCollection(r *ProductsResponse) ([]Product, *ListParams, error) {
    dst := make([]Product, len(r.Items))

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

// GetById retrieves collection of products of current tenant
func (s *ProductsServiceOp) List(args ...interface{}) ([]Product, *ListParams, error) {
    endpoint := fmt.Sprintf("tenants/%s/products", s.client.tenantId)
    return s.ListByLink(endpoint, args...)
}

// GetById retrieves collection of products by link
func (s *ProductsServiceOp) ListByLink(endpoint string, args ...interface{}) ([]Product, *ListParams, error) {
    resp, err := s.client.request("GET", endpoint, nil, args...)
    if err != nil {
        return nil, nil, err
    }

    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, nil, fmt.Errorf("Status code: %d", resp.StatusCode)
    }
    obj := &ProductsResponse{}
    dec := json.NewDecoder(resp.Body)
    dec.Decode(obj)

    return s.getCollection(obj)
}

// GetById updates product with specified ID
func (s *ProductsServiceOp) UpdateById(id string, t *ProductRequestUpdate) (*Product, error) {
    endpoint := fmt.Sprintf("products/%s", id)
    return s.UpdateByLink(endpoint, t)
}

// GetById updates product specified by link
func (s *ProductsServiceOp) UpdateByLink(endpoint string, t *ProductRequestUpdate) (*Product, error) {
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
    obj := &ProductResponse{}
    dec := json.NewDecoder(resp.Body)
    dec.Decode(obj)
    return s.get(obj)
}

// Create creates new product within tenant
func (s *ProductsServiceOp) Create(dir *ProductRequestCreate) (*Product, error) {
    endpoint := fmt.Sprintf("tenants/%s/products", s.client.tenantId)

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
    obj := &ProductResponse{}
    dec := json.NewDecoder(resp.Body)
    dec.Decode(obj)
    return s.get(obj)
}

// Delete removes product
func (s *ProductsServiceOp) Delete(t *Product) (error) {
    return s.DeleteByLink(t.Href)
}

// Delete removes product by ID
func (s *ProductsServiceOp) DeleteById(id string) (error) {
    endpoint := fmt.Sprintf("products/%s", id)
    return s.DeleteByLink(endpoint)
}

// Delete removes product by link
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
