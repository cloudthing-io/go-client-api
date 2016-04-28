package api

import (
    "bytes"
    "encoding/json" 
    "fmt"   
    "net/http"    
    "github.com/borystomala/copier"
)

// ExportsService is an interface for interacting with Exports endpoints of CloudThing API
// https://tenant-name.cloudthing.io/api/v1/exports
type ExportsService interface {
    GetById(string, ...interface{}) (*Export, error)
    GetByLink(string, ...interface{}) (*Export, error)
    ListByLink(string, ...interface{}) ([]Export, *ListParams, error)
    ListByApplication(string, ...interface{}) ([]Export, *ListParams, error)
    List(...interface{}) ([]Export, *ListParams, error)
    CreateByLink(string, *ExportRequestCreate) (*Export, error)
    CreateByApplication(string, *ExportRequestCreate) (*Export, error)
    UpdateById(string, *ExportRequestUpdate) (*Export, error)
    UpdateByLink(string, *ExportRequestUpdate) (*Export, error)
    Delete(*Export) (error)
    DeleteByLink(string) (error)
    DeleteById(string) (error)

    get(*ExportResponse) (*Export, error)
    getCollection(*ExportsResponse) ([]Export, *ListParams, error)
}

// ExportsServiceOp handles communication with Exports related methods of API
type ExportsServiceOp struct {
    client *Client
}

// Export is a struct representing CloudThing Export
type Export struct {
    // Standard field for all resources
    ModelBase
    ModelType       string
    LimitsType      string 
    Export          []ExportEntry
    TenExpPerm      string

    // Points to Tenant if expansion was requested, otherwise nil
    TenantExp           *Tenant
    // Points to Tenant if expansion was requested, otherwise nil
    TenantImp           *Tenant
    // Points to Applications if expansion was requested, otherwise nil
    Product             *Product
    // Points to Tenant if expansion was requested, otherwise nil
    Limits              interface{}
    Application         *Application
    TenExpPermExp       *Export

    // Links to related resources
    tenantExp           string
    tenantImp           string
    product             string
    limits              string
    application         string
    tenExpPermExp       string

    // service for communication, internal use only
    service         *ExportsServiceOp
}

type ExportEntry struct {
    Type            string                  `json:"type"`
    Name            string                  `json:"name"`
    Read            bool                    `json:"read"`
    Write           bool                    `json:"write"`
    GrantRead       bool                    `json:"grantRead"`
    GrantWrite      bool                    `json:"grantWrite"`
}

// ExportResponse is a struct representing item response from API
type ExportResponse struct {
    ModelBase
    ModelType       string                  `json:"modelType,omitempty"`    
    LimitsType      string                  `json:"limitsType,omitempty"`    
    Export          []ExportEntry           `json:"export,omitempty"`
    TenExpPerm      string                  `json:"tenantExportingPermission,omitempty"`

    Limits          interface{}             `json:"limits,omitempty"`
    Product         interface{}             `json:"product,omitempty"`
    TenantExp       map[string]interface{}  `json:"tenantExp,omitempty"`
    TenantImp       map[string]interface{}  `json:"tenantImp,omitempty"`
    Application     map[string]interface{}  `json:"application,omitempty"`
    TenExpPermExp   interface{}             `json:"tenantExportingPermissionExport,omitempty"` 
}


// ExportResponse is a struct representing collection response from API
type ExportsResponse struct{
    ListParams
    Items           []ExportResponse              `json:"items"`
}

// ExportResponse is a struct representing item create request for API
type ExportRequestCreate struct {
    ModelType       string                  `json:"modelType,omitempty"`    
    LimitsType      string                  `json:"limitsType,omitempty"`    
    Export          []ExportEntry           `json:"export,omitempty"`
    TenExpPerm      string                  `json:"tenantExportingPermission,omitempty"`
    Limits          *Link                   `json:"limits,omitempty"`
    Product         *Link                   `json:"product,omitempty"`
}

// ExportResponse is a struct representing item update request for API
type ExportRequestUpdate struct {
    ModelType       string                  `json:"modelType,omitempty"`    
    LimitsType      string                  `json:"limitsType,omitempty"`    
    Export          []ExportEntry           `json:"export,omitempty"`
    TenExpPerm      string                  `json:"tenantExportingPermission,omitempty"`
    Limits          *Link                   `json:"limits,omitempty"`
    Product         *Link                   `json:"product,omitempty"`
}

// TenantLink returns indicator of Tenant expansion and link to tenant.
// If expansion for Tenant was requested and resource is available via pointer
// it returns true, otherwise false. Link (href) is always returned. 
func (d *Export) TenantExpLink() (bool, string) {
    return (d.TenantExp != nil), d.tenantExp
}

// TenantLink returns indicator of Tenant expansion and link to tenant.
// If expansion for Tenant was requested and resource is available via pointer
// it returns true, otherwise false. Link (href) is always returned. 
func (d *Export) TenantImpLink() (bool, string) {
    return (d.TenantImp != nil), d.tenantImp
}

// ApplicationsLink returns indicator of Export expansion and link to dierctory.
// If expansion for Export was requested and resource is available via pointer
// it returns true, otherwise false. Link (href) is always returned. 
func (d *Export) ProductLink() (bool, string) {
    if d.ModelType != "DEVICE" {
        return false, ""
    }
    return (d.Product != nil), d.product
}

// ApplicationsLink returns indicator of Export expansion and link to dierctory.
// If expansion for Export was requested and resource is available via pointer
// it returns true, otherwise false. Link (href) is always returned. 
func (d *Export) LimitsLink() (bool, string) {
    if d.LimitsType == "" {
        return false, ""
    }
    return (d.Limits != nil), d.limits
}

// ApplicationsLink returns indicator of Export expansion and link to dierctory.
// If expansion for Export was requested and resource is available via pointer
// it returns true, otherwise false. Link (href) is always returned. 
func (d *Export) ApplicationLink() (bool, string) {
    return (d.Application != nil), d.application
}

// ApplicationsLink returns indicator of Export expansion and link to dierctory.
// If expansion for Export was requested and resource is available via pointer
// it returns true, otherwise false. Link (href) is always returned. 
func (d *Export) TenExpPermExpLink() (bool, string) {
    if d.TenExpPerm == "PRIMARY" {
        return false, ""
    }
    return (d.TenExpPermExp != nil), d.tenExpPermExp
}


// Save is a helper method for updating apikey.
// It calls UpdateByLink() on service under the hood.
func (t *Export) Save() error {
    tmp := &ExportRequestUpdate{}
    copier.Copy(tmp, t)
    ten, err := t.service.UpdateByLink(t.Href, tmp)
    if err != nil {
        return err
    }

    tmpTenantExp := t.TenantExp
    tmpTenantImp := t.TenantImp
    tmpProduct := t.Product
    tmpLimits := t.Limits
    tmpApplication := t.Application
    tmpTenExpPermExp := t.TenExpPermExp

    *t = *ten
    t.TenantExp = tmpTenantExp
    t.TenantImp = tmpTenantImp
    t.Product = tmpProduct
    t.Limits = tmpLimits
    t.Application = tmpApplication
    t.TenExpPermExp = tmpTenExpPermExp

    return nil
}

// Save is a helper method for deleting apikey.
// It calls Delete() on service under the hood.
func (t *Export) Delete() error {
    err := t.service.Delete(t)
    if err != nil {
        return err
    }

    return nil
}

// GetById retrieves apikey by its ID
func (s *ExportsServiceOp) GetById(id string, args ...interface{}) (*Export, error) {
    endpoint := "exports/"
    endpoint = fmt.Sprintf("%s%s", endpoint, id)

    return s.GetByLink(endpoint, args...)
}

// GetById retrieves apikey by its full link
func (s *ExportsServiceOp) GetByLink(endpoint string, args ...interface{}) (*Export, error) {
    resp, err := s.client.request("GET", endpoint, nil, args...)
    if err != nil {
        return nil, err
    }

    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("Status code: %d", resp.StatusCode)
    }
    obj := &ExportResponse{}
    dec := json.NewDecoder(resp.Body)
    dec.Decode(obj)

    return s.get(obj)
}

// get is internal method for transforming ExportResponse into Export
func (s *ExportsServiceOp) get(r *ExportResponse) (*Export, error) {
    obj := &Export{}
    copier.Copy(obj, r)
    if v, ok :=  r.TenantExp["href"]; ok {
        obj.tenantExp = v.(string)
    }
    if v, ok :=  r.TenantImp["href"]; ok {
        obj.tenantImp = v.(string)
    }
    if v, ok :=  r.Product.(map[string]interface{}); ok {
        if w, ok := v["href"]; ok {
            obj.product = w.(string)
        }
    }
    if v, ok :=  r.Limits.(map[string]interface{}); ok {
        if w, ok := v["href"]; ok {
            obj.limits = w.(string)
        }
    }
    if v, ok :=  r.Application["href"]; ok {
        obj.application = v.(string)
    }
    if v, ok :=  r.TenExpPermExp.(map[string]interface{}); ok {
        if w, ok := v["href"]; ok {
            obj.tenExpPermExp = w.(string)
        }
    }
    if len(r.TenantExp) > 1 {        
        bytes, err := json.Marshal(r.TenantExp)
        if err != nil {
            return nil, err
        }
        ten := &TenantResponse{}
        json.Unmarshal(bytes, ten)
        t, err := s.client.Tenant.get(ten)
        if err != nil {
            return nil, err
        }
        obj.TenantExp = t
    }
    if v, ok :=  r.Product.(map[string]interface{}); ok {        
        if len(v) > 1 {        
            bytes, err := json.Marshal(v)
            if err != nil {
                return nil, err
            }
            ten := &ProductResponse{}
            json.Unmarshal(bytes, ten)
            t, err := s.client.Products.get(ten)
            if err != nil {
                return nil, err
            }
            obj.Product = t
        }
    }
  /*  if len(r.Limit) > 1 {        
        bytes, err := json.Marshal(r.Product)
        if err != nil {
            return nil, err
        }
        ten := &ProductResponse{}
        json.Unmarshal(bytes, ten)
        t, err := s.client.Products.get(ten)
        if err != nil {
            return nil, err
        }
        obj.Product = t
    }*/
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

// get is internal method for transforming ApplicationResponse into Application
func (s *ExportsServiceOp) getCollection(r *ExportsResponse) ([]Export, *ListParams, error) {
    dst := make([]Export, len(r.Items))

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

// GetById retrieves collection of exports of current tenant
func (s *ExportsServiceOp) ListByApplication(id string, args ...interface{}) ([]Export, *ListParams, error) {
    endpoint := fmt.Sprintf("applications/%s/exports", id)
    return s.ListByLink(endpoint, args...)
}

// GetById retrieves collection of exports of current tenant
func (s *ExportsServiceOp) List(args ...interface{}) ([]Export, *ListParams, error) {
    endpoint := fmt.Sprintf("tenants/%s/exports", s.client.tenantId)
    return s.ListByLink(endpoint, args...)
}



// GetById retrieves collection of exports by link
func (s *ExportsServiceOp) ListByLink(endpoint string, args ...interface{}) ([]Export, *ListParams, error) {
    resp, err := s.client.request("GET", endpoint, nil, args...)
    if err != nil {
        return nil, nil, err
    }

    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, nil, fmt.Errorf("Status code: %d", resp.StatusCode)
    }
    obj := &ExportsResponse{}
    dec := json.NewDecoder(resp.Body)
    dec.Decode(obj)

    return s.getCollection(obj)
}

// GetById updates apikey with specified ID
func (s *ExportsServiceOp) UpdateById(id string, t *ExportRequestUpdate) (*Export, error) {
    endpoint := fmt.Sprintf("exports/%s", id)
    return s.UpdateByLink(endpoint, t)
}

// GetById updates apikey specified by link
func (s *ExportsServiceOp) UpdateByLink(endpoint string, t *ExportRequestUpdate) (*Export, error) {
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
    obj := &ExportResponse{}
    dec := json.NewDecoder(resp.Body)
    dec.Decode(obj)
    return s.get(obj)
}

func (s *ExportsServiceOp) CreateByApplication(id string, dir *ExportRequestCreate) (*Export, error) {
    endpoint := fmt.Sprintf("applications/%s/exports", id)
    return s.CreateByLink(endpoint, dir)
}

// Create creates new apikey within tenant
func (s *ExportsServiceOp) CreateByLink(endpoint string, dir *ExportRequestCreate) (*Export, error) {
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
    obj := &ExportResponse{}
    dec := json.NewDecoder(resp.Body)
    dec.Decode(obj)
    return s.get(obj)
}

// Delete removes apikey
func (s *ExportsServiceOp) Delete(t *Export) (error) {
    return s.DeleteByLink(t.Href)
}

// Delete removes apikey by ID
func (s *ExportsServiceOp) DeleteById(id string) (error) {
    endpoint := fmt.Sprintf("exports/%s", id)
    return s.DeleteByLink(endpoint)
}

// Delete removes apikey by link
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
