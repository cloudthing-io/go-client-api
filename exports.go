package api

import (
    "bytes"
    "encoding/json" 
    "fmt"   
    "net/http"
    _"strings"
)

// ExportsService is an interafce for interacting with Exports endpoints of CloudThing API
// https://tenant-name.cloudthing.io/api/v1/directories

type ExportsService interface {
    GetById(string) (*Export, error)
    GetByLink(string) (*Export, error)
    List(*ListOptions) ([]Export, *ListParams, error)
    ListByLink(string, *ListOptions) ([]Export, *ListParams, error)
    ListByApplication(string, *ListOptions) ([]Export, *ListParams, error)
    CreateByLink(string,*Export) (*Export, error)
    CreateByApplication(string,*Export) (*Export, error)
    Update(*Export) (*Export, error)
    Delete(*Export) (error)
    DeleteByLink(string) (error)
    DeleteById(string) (error)
}

// ExportsServiceOp handles communication with Tenant related methods of API
type ExportsServiceOp struct {
    client *Client
}

type ExportEntry struct {
    Type            string          `json:"type"`
    Name            string          `json:"name"`
    Read            bool            `json:"read"`
    Write           bool            `json:"write"`
    GrantRead       bool            `json:"grantRead"`
    GrantWrite      bool            `json:"grantWrite"`
}

type Export struct {
    ModelBase
    ModelType       string          `json:"modelType,omitempty"`    
    LimitsType      string          `json:"limitsType,omitempty"`    
    Export          []ExportEntry   `json:"export,omitempty"`
    TenExpPerm      string          `json:"tenantExportingPermission,omitempty"`

    Limits          *Link           `json:"limits,omitempty"`
    Product         *Link           `json:"product,omitempty"`
    TenantExp       *Link           `json:"tenantExp,omitempty"`
    TenantImp       *Link           `json:"tenantImp,omitempty"`
    Application     *Link           `json:"application,omitempty"`
    TenExpPermExp   *Link           `json:"tenantExportingPermissionExport,omitempty"`   

    // service for communication, internal use only
    service         *ExportsServiceOp `json:"-"` 
}

type Exports struct{
    ListParams
    Items           []Export     `json:"items"`
}

func (d *Export) GetTenant() (*Tenant, error) {
    return d.service.client.Tenant.Get()
}

func (d *Export) GetProduct() (*Product, error) {
    if d.ModelType != "DEVICE" {
        return nil, fmt.Errorf("This export's model is not device")
    }
    return d.service.client.Products.GetByLink(d.Product.Href)
}

func (d *Export) SetProduct(p *Product) (error) {
    if d.ModelType != "DEVICE" {
        return fmt.Errorf("This export's model is not device")
    }
    d.Product = &Link{p.Href}
    fmt.Println(d)
    return nil
}
func (d *Export) GetApplication() (*Application, error) {
    return d.service.client.Applications.GetByLink(d.Application.Href)
}

func (d *Export) GetTenExpPermExp() (*Export, error) {
    return d.service.GetByLink(d.TenExpPermExp.Href)
}

func (d *Export) GetLimits() (string, interface{}, error) {
    if d.LimitsType == "" {
        return "", nil, nil
    } else if d.LimitsType == "DEVICE" {
         res, err := d.service.client.Devices.GetByLink(d.Limits.Href)
         return d.LimitsType, res, err
    } else if d.LimitsType == "GROUP" {
         res, err := d.service.client.Groups.GetByLink(d.Limits.Href)
         return d.LimitsType, res, err
    } else if d.LimitsType == "CLUSTER" {
         res, err := d.service.client.Clusters.GetByLink(d.Limits.Href)
         return d.LimitsType, res, err
    }
   return "", nil, nil
}

// Save updates export by calling Update() on service under the hood
func (t *Export) Save() error {
    tmp := *t

    ten, err := t.service.Update(&tmp)
    if err != nil {
        return err
    }

    *t = *ten
    return nil
}

func (t *Export) Delete() error {
    err := t.service.Delete(t)
    if err != nil {
        return err
    }

    return nil
}

// GetById retrieves export
func (s *ExportsServiceOp) GetById(id string) (*Export, error) {
    endpoint := "exports/"
    endpoint = fmt.Sprintf("%s%s", endpoint, id)

    return s.GetByLink(endpoint)
}

func (s *ExportsServiceOp) GetByLink(endpoint string) (*Export, error) {
    resp, err := s.client.request("GET", endpoint, nil)
    if err != nil {
        return nil, err
    }

    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("Status code: %d", resp.StatusCode)
    }
    obj := &Export{}
    dec := json.NewDecoder(resp.Body)
    dec.Decode(obj)
    obj.service = s
    return obj, nil
}


func (s *ExportsServiceOp) List(lo *ListOptions) ([]Export, *ListParams, error) {
    endpoint := fmt.Sprintf("/exports")
    return s.ListByLink(endpoint, lo)
}

func (s *ExportsServiceOp) ListByApplication(id string, lo *ListOptions) ([]Export, *ListParams, error) {
    endpoint := fmt.Sprintf("applications/%s/exports", id)
    return s.ListByLink(endpoint, lo)
}

func (s *ExportsServiceOp) ListByLink(endpoint string, lo *ListOptions) ([]Export, *ListParams, error) {
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
    obj := &Exports{}
    dec := json.NewDecoder(resp.Body)
    dec.Decode(obj)

    dst := make([]Export, len(obj.Items))
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

// Update updates tenant
func (s *ExportsServiceOp) Update(t *Export) (*Export, error) {
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
    obj := &Export{}
    dec := json.NewDecoder(resp.Body)
    dec.Decode(obj)
    obj.service = s
    return obj, nil
}

func (s *ExportsServiceOp) CreateByApplication(id string, dir *Export) (*Export, error) {
    endpoint := fmt.Sprintf("applications/%s/exports", id)
    return s.CreateByLink(endpoint, dir)
}

func (s *ExportsServiceOp) CreateByLink(endpoint string, dir *Export) (*Export, error) {
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
    obj := &Export{}
    dec := json.NewDecoder(resp.Body)
    dec.Decode(obj)
    obj.service = s
    return obj, nil
}

// Delete removes export
func (s *ExportsServiceOp) Delete(t *Export) (error) {
    return s.DeleteByLink(t.Href)
}

// Delete removes export by ID
func (s *ExportsServiceOp) DeleteById(id string) (error) {
    endpoint := fmt.Sprintf("exports/%s", id)
    return s.DeleteByLink(endpoint)
}

// Delete removes export by link
func (s *ExportsServiceOp) DeleteByLink(endpoint string) (error) {
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
