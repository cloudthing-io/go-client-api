package api

import (
    "net/http"
    "net/url"
    "fmt"
    "encoding/json"
    "io"
    "time"
    "github.com/dgrijalva/jwt-go"
    "strings"
)

const (
    // Current version of client
    version     =       "0.1.0"
    // Default user-agent
    userAgent   =       "cloudthing-go-client/" + version
    // Default media-type used in Accept and Content-Type headers
    mediaType   =       "application/json"
    // API endpoint for v1
    apiEndpoint =       "/api/v1/"
)

// Link represents unexpanded link to child/parent resource
type Link struct{
    Href             *string         `json:"href,omitempty"`
}

// ModelBase is a base structure for all models in CloudThing API
type ModelBase struct {
    Href            string          `json:"href,omitempty"`
    CreatedAt       *time.Time      `json:"createdAt,omitempty"`
    UpdatedAt       *time.Time      `json:"updatedAt,omitempty"`
}

// Token is a JSON Web Token (JWT) used for Authorization
type Token struct {
    Token       string      `json:"token"`
    Type        string      `json:"type"`
    ExpiresIn   int64       `json:"expiresIn"`
}

// Client manages communication with CloudThing API
type Client struct {
    // HTTP client used to communicate with API
    client          *http.Client  

    // Base URL for API requests
    BaseURL         *url.URL 

    // User agent for client
    UserAgent       string

    // JWT token for authorization
    token           *Token

    // Tenant ID
    tenantId        string

    // Services used for communication
    Tenant          TenantService
    Directories     DirectoriesService
    Applications    ApplicationsService
    Products        ProductsService
    Devices         DevicesService
    Clusters        ClustersService
    Groups          GroupsService
    Users           UsersService
    Apikeys         ApikeysService
    ClusterMemberships ClusterMembershipsService 
    GroupMemberships GroupMembershipsService 
    Memberships     MembershipsService
    Usergroups      UsergroupsService
}

// ListOptions specifies the optional parameters for requests with pagination support
type ListOptions struct {
    // Page of results to retrieve
    Page        int         `url:"page,omitempty"`
    Limit       int         `url:"limit,omitempty"`
}

type ListParams struct {
    Href            string          `json:"href"`
    Size            int             `json:"size"`
    Limit           int             `json:"limit"`
    Page            int             `json:"page"`
    Prev            *Link            `json:"prev,omitempty"`
    Next            *Link            `json:"next,omitempty"`
}

const (
    DefaultLimit = 25
)

// NewClient returns a new CloudThing API client
func NewClient(httpClient *http.Client, baseURL string) (*Client, error) {
    if httpClient == nil {
        httpClient = http.DefaultClient
    }

    base, err := url.Parse(baseURL + apiEndpoint)
    if err != nil {
        return nil, err
    }

    c := &Client {
        client: httpClient,        
        BaseURL: base,
        UserAgent: userAgent,
    }

    c.Tenant = &TenantServiceOp{client: c}
    c.Directories = &DirectoriesServiceOp{client: c}
    c.Applications = &ApplicationsServiceOp{client: c}
    c.Products = &ProductsServiceOp{client: c}
    c.Devices = &DevicesServiceOp{client: c}
    c.Clusters = &ClustersServiceOp{client: c}
    c.Users = &UsersServiceOp{client: c}
    c.ClusterMemberships = &ClusterMembershipsServiceOp{client: c}
    c.GroupMemberships = &GroupMembershipsServiceOp{client: c} 
    c.Memberships = &MembershipsServiceOp{client: c}
    c.Usergroups = &UsergroupsServiceOp{client: c}
    c.Apikeys = &ApikeysServiceOp{client: c}

    return c, nil
}

// SetUserAgent sets required UserAgent string for further requests
func (c *Client) SetUserAgent(ua string) {
    c.UserAgent = ua
}

// SetBasicAuth uses provided basic authorization params for authenticating against 
// CloudThing API and retrieves and stores JWT token if succeeded for future requests.
func (c *Client) SetBasicAuth(username, password string) error {
    endpoint := "auth/token"
    endp, err := url.Parse(endpoint)
    if err != nil {
        return err
    }

    u := c.BaseURL.ResolveReference(endp)

    req, err := http.NewRequest("POST", u.String(), nil)
    if err != nil {
        return err
    }
    req.SetBasicAuth(username, password)
    resp, err := c.client.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("Failed to authenticate user: %d\n", resp.StatusCode)
    }

    token := &Token{}
    dec := json.NewDecoder(resp.Body)
    dec.Decode(token)

    c.setToken(token)
    return nil
}

// SetTokenAuth uses provided JWT token for authenticating against CloudThing API
// and stores it if succeeded for future requests.
func (c *Client) SetTokenAuth(token *Token) error {
    req, err := http.NewRequest("GET", c.BaseURL.String(), nil)
    if err != nil {
        return err
    }
    req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token.Token))
    resp, err := c.client.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("Failed to authenticate user: %d\n", resp.StatusCode)
    }

    c.setToken(token)
    return nil
}

// setToken parses JWT token, extracts tenant ID and sets token and in in client
func (c *Client) setToken(t *Token) {
    token, _ := jwt.Parse(t.Token, nil)
    claims := token.Claims.(jwt.MapClaims)
    iss := strings.Split(claims["iss"].(string), "/")
    c.tenantId = iss[len(iss)-1]
    c.token = t
}

// Creates new request or sending to API
func (c *Client) request(method, endpoint string, body io.Reader) (*http.Response, error) {
    if !c.IsAuthenticated() {
        return nil, fmt.Errorf("You need to authenticate first")
    }

    u, err := url.Parse(endpoint)
    if err != nil {
        return nil, err
    }

    if !u.IsAbs() {
        u = c.BaseURL.ResolveReference(u)
    }

    req, err := http.NewRequest(method, u.String(), body)
    if err != nil {
        return nil, err
    }

    req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.token.Token))
    req.Header.Add("Accept", mediaType)
    req.Header.Add("Content-Type", mediaType)
    req.Header.Add("User-Agent", c.UserAgent)

    resp, err := c.client.Do(req)
    if err != nil {
        return nil, err
    }
    return resp, nil
}

// Checkes whether Client is authentciated and able to create requests
func (c *Client) IsAuthenticated() bool {
    if c.token != nil {
        t, _ := jwt.Parse(c.token.Token, nil)
        var exp int64
        var vexp bool
        var err error
        now := time.Now().Unix()
        claims := t.Claims.(jwt.MapClaims)
        switch num := claims["exp"].(type) {
        case json.Number:
            if exp, err = num.Int64(); err == nil {
                vexp = true
            }
        case float64:
            vexp = true
            exp = int64(num)
        }

        if vexp && exp > now {
            return true
        }
    }
    return false
}

// GetToken retr=urn curent JWT token
func (c *Client) GetToken() *Token {
    return c.token
}