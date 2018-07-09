package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/borystomala/copier"
)

// UsersService is an interface for interacting with Users endpoints of CloudThing API
// https://tenant-name.cloudthing.io/api/v1/users
type UsersService interface {
	GetCurrent(...interface{}) (*User, error)
	GetById(string, ...interface{}) (*User, error)
	GetByLink(string, ...interface{}) (*User, error)
	ListByLink(string, ...interface{}) ([]User, *ListParams, error)
	ListByDirectory(string, ...interface{}) ([]User, *ListParams, error)
	ListByUsergroup(string, ...interface{}) ([]User, *ListParams, error)
	CreateByLink(string, *UserRequestCreate) (*User, error)
	CreateByDirectory(string, *UserRequestCreate) (*User, error)
	UpdateById(string, *UserRequestUpdate) (*User, error)
	UpdateByLink(string, *UserRequestUpdate) (*User, error)
	Delete(*User) error
	DeleteByLink(string) error
	DeleteById(string) error

	get(*UserResponse) (*User, error)
	getCollection(*UsersResponse) ([]User, *ListParams, error)
}

// UsersServiceOp handles communication with Users related methods of API
type UsersServiceOp struct {
	client *Client
}

// User is a struct representing CloudThing User
type User struct {
	// Standard field for all resources
	ModelBase

	Username            string                 `json:"username,omitempty"`
	Email               string                 `json:"email,omitempty"`
	FirstName           string                 `json:"firstName,omitempty"`
	Surname             string                 `json:"surname,omitempty"`
	Password            string                 `json:"password,omitempty"`
	Activated           bool                   `json:"activated,omitempty"`
	LastSuccessfulLogin *time.Time             `json:"lastSuccessfulLogin,omitempty"`
	LastFailedLogin     *time.Time             `json:"lastFailedLogin,omitempty"`
	ActivationCode      string                 `json:"activationCode,omitempty"`
	Custom              map[string]interface{} `json:"custom,omitempty"`

	// Points to Tenant if expansion was requested, otherwise nil
	Tenant *Tenant `json:"tenant,omitempty"`
	// Points to Applications if expansion was requested, otherwise nil
	Directory *Directory `json:"directory,omitempty"`
	// Points to Tenant if expansion was requested, otherwise nil
	Applications []Application `json:"applications,omitempty"`
	Usergroups   []Usergroup   `json:"usergroups,omitempty"`

	Memberships []Membership `json:"memberships,omitempty"`

	tenant       string `json:"tenant,omitempty"`
	directory    string `json:"directory,omitempty"`
	applications string `json:"applications,omitempty"`
	usergroups   string `json:"usergroups,omitempty"`
	memberships  string `json:"memberships,omitempty"`

	// service for communication, internal use only
	service *UsersServiceOp
}

// UserResponse is a struct representing item response from API
type UserResponse struct {
	ModelBase
	Username            string                 `json:"username,omitempty"`
	Email               string                 `json:"email,omitempty"`
	FirstName           string                 `json:"firstName,omitempty"`
	Surname             string                 `json:"surname,omitempty"`
	LastSuccessfulLogin *time.Time             `json:"lastSuccessfulLogin,omitempty"`
	LastFailedLogin     *time.Time             `json:"lastFailedLogin,omitempty"`
	ActivationCode      string                 `json:"activationCode,omitempty"`
	Custom              map[string]interface{} `json:"custom,omitempty"`

	Tenant       map[string]interface{} `json:"tenant,omitempty"`
	Applications map[string]interface{} `json:"applications,omitempty"`
	Directory    map[string]interface{} `json:"directory,omitempty"`
	Usergroups   map[string]interface{} `json:"usergroups,omitempty"`
	Memberships  map[string]interface{} `json:"memberships,omitempty"`
}

// UserResponse is a struct representing collection response from API
type UsersResponse struct {
	ListParams
	Items []UserResponse `json:"items"`
}

// UserResponse is a struct representing item create request for API
type UserRequestCreate struct {
	Username  string                 `json:"username,omitempty"`
	Email     string                 `json:"email,omitempty"`
	FirstName string                 `json:"firstName,omitempty"`
	Surname   string                 `json:"surname,omitempty"`
	Password  string                 `json:"password,omitempty"`
	Activated bool                   `json:"activated,omitempty"`
	Custom    map[string]interface{} `json:"custom,omitempty"`
}

// UserResponse is a struct representing item update request for API
type UserRequestUpdate struct {
	Username  string                 `json:"username,omitempty"`
	Email     string                 `json:"email,omitempty"`
	FirstName string                 `json:"firstName,omitempty"`
	Surname   string                 `json:"surname,omitempty"`
	Password  string                 `json:"password,omitempty"`
	Activated bool                   `json:"activated,omitempty"`
	Custom    map[string]interface{} `json:"custom,omitempty"`
}

// TenantLink returns indicator of Tenant expansion and link to tenant.
// If expansion for Tenant was requested and resource is available via pointer
// it returns true, otherwise false. Link (href) is always returned.
func (d *User) TenantLink() (bool, string) {
	return (d.Tenant != nil), d.tenant
}

// ApplicationsLink returns indicator of User expansion and link to dierctory.
// If expansion for User was requested and resource is available via pointer
// it returns true, otherwise false. Link (href) is always returned.
func (d *User) ApplicationsLink() (bool, string) {
	return (d.Applications != nil), d.applications
}

// ApplicationsLink returns indicator of User expansion and link to dierctory.
// If expansion for User was requested and resource is available via pointer
// it returns true, otherwise false. Link (href) is always returned.
func (d *User) DirectoryLink() (bool, string) {
	return (d.Directory != nil), d.directory
}

// ApplicationsLink returns indicator of User expansion and link to dierctory.
// If expansion for User was requested and resource is available via pointer
// it returns true, otherwise false. Link (href) is always returned.
func (d *User) UsergroupsLink() (bool, string) {
	return (d.Usergroups != nil), d.usergroups
}

// MembershipsLink returns indicator of User expansion and link to dierctory.
// If expansion for User was requested and resource is available via pointer
// it returns true, otherwise false. Link (href) is always returned.
func (d *User) MembershipsLink() (bool, string) {
	return (d.Memberships != nil), d.memberships
}

// Save is a helper method for updating apikey.
// It calls UpdateByLink() on service under the hood.
func (t *User) Save() error {
	tmp := &UserRequestUpdate{}
	copier.Copy(tmp, t)
	ten, err := t.service.UpdateByLink(t.Href, tmp)
	if err != nil {
		return err
	}

	tmpTenant := t.Tenant
	tmpApplications := t.Applications
	tmpDirectory := t.Directory
	tmpUsergroups := t.Usergroups
	tmpMemberships := t.Memberships

	*t = *ten
	t.Tenant = tmpTenant
	t.Applications = tmpApplications
	t.Directory = tmpDirectory
	t.Usergroups = tmpUsergroups
	t.Memberships = tmpMemberships

	return nil
}

// Save is a helper method for deleting apikey.
// It calls Delete() on service under the hood.
func (t *User) Delete() error {
	err := t.service.Delete(t)
	if err != nil {
		return err
	}

	return nil
}

// Get retrieves current user
func (s *UsersServiceOp) GetCurrent(args ...interface{}) (*User, error) {
	endpoint := "users/current"

	resp, err := s.client.request("GET", endpoint, nil, args...)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusForbidden || resp.StatusCode == http.StatusUnauthorized {
		// this is probably due to redirect
		endpoint = resp.Request.URL.String()
		resp, err = s.client.request("GET", endpoint, nil, args...)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return nil, ApiError{StatusCode: resp.StatusCode, Message: "non-ok status returned"}
		}
	} else if resp.StatusCode != http.StatusOK {
		return nil, ApiError{StatusCode: resp.StatusCode, Message: "non-ok status returned"}
	}
	user := &UserResponse{}
	dec := json.NewDecoder(resp.Body)
	dec.Decode(user)
	return s.get(user)
}

// GetById retrieves apikey by its ID
func (s *UsersServiceOp) GetById(id string, args ...interface{}) (*User, error) {
	endpoint := "users/"
	endpoint = fmt.Sprintf("%s%s", endpoint, id)

	return s.GetByLink(endpoint, args...)
}

// GetById retrieves apikey by its full link
func (s *UsersServiceOp) GetByLink(endpoint string, args ...interface{}) (*User, error) {
	resp, err := s.client.request("GET", endpoint, nil, args...)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, ApiError{StatusCode: resp.StatusCode, Message: "non-ok status returned"}
	}
	obj := &UserResponse{}
	dec := json.NewDecoder(resp.Body)
	dec.Decode(obj)

	return s.get(obj)
}

// get is internal method for transforming UserResponse into User
func (s *UsersServiceOp) get(r *UserResponse) (*User, error) {
	obj := &User{}
	copier.Copy(obj, r)
	if v, ok := r.Tenant["href"]; ok {
		obj.tenant = v.(string)
	}
	if v, ok := r.Applications["href"]; ok {
		obj.applications = v.(string)
	}
	if v, ok := r.Directory["href"]; ok {
		obj.directory = v.(string)
	}
	if v, ok := r.Usergroups["href"]; ok {
		obj.usergroups = v.(string)
	}
	if v, ok := r.Memberships["href"]; ok {
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
	if len(r.Applications) > 1 {
		bytes, err := json.Marshal(r.Applications)
		if err != nil {
			return nil, err
		}
		ten := &ApplicationsResponse{}
		json.Unmarshal(bytes, ten)
		t, _, err := s.client.Applications.getCollection(ten)
		if err != nil {
			return nil, err
		}
		obj.Applications = t
	}
	if len(r.Directory) > 1 {
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
	}
	if len(r.Usergroups) > 1 {
		bytes, err := json.Marshal(r.Usergroups)
		if err != nil {
			return nil, err
		}
		ugs := &UsergroupsResponse{}
		json.Unmarshal(bytes, ugs)
		u, _, err := s.client.Usergroups.getCollection(ugs)
		if err != nil {
			return nil, err
		}
		obj.Usergroups = u
	}
	obj.service = s
	return obj, nil
}

// get is internal method for transforming ApplicationResponse into Application
func (s *UsersServiceOp) getCollection(r *UsersResponse) ([]User, *ListParams, error) {
	dst := make([]User, len(r.Items))

	for i, _ := range r.Items {
		t, err := s.get(&r.Items[i])
		if err == nil {
			dst[i] = *t
		}
	}

	lp := &ListParams{
		Href:  r.Href,
		Prev:  r.Prev,
		Next:  r.Next,
		Limit: r.Limit,
		Size:  r.Size,
		Page:  r.Page,
	}
	return dst, lp, nil
}

// GetById retrieves collection of users of current tenant
func (s *UsersServiceOp) ListByDirectory(id string, args ...interface{}) ([]User, *ListParams, error) {
	endpoint := fmt.Sprintf("directories/%s/users", id)
	return s.ListByLink(endpoint, args...)
}

// ListByUsergroup retrieves collection of users of current usergroup
func (s *UsersServiceOp) ListByUsergroup(id string, args ...interface{}) ([]User, *ListParams, error) {
	endpoint := fmt.Sprintf("usergroups/%s/users", id)
	return s.ListByLink(endpoint, args...)
}

// GetById retrieves collection of users by link
func (s *UsersServiceOp) ListByLink(endpoint string, args ...interface{}) ([]User, *ListParams, error) {
	resp, err := s.client.request("GET", endpoint, nil, args...)
	if err != nil {
		return nil, nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, nil, ApiError{StatusCode: resp.StatusCode, Message: "non-ok status returned"}
	}
	obj := &UsersResponse{}
	dec := json.NewDecoder(resp.Body)
	dec.Decode(obj)

	return s.getCollection(obj)
}

// GetById updates apikey with specified ID
func (s *UsersServiceOp) UpdateById(id string, t *UserRequestUpdate) (*User, error) {
	endpoint := fmt.Sprintf("users/%s", id)
	return s.UpdateByLink(endpoint, t)
}

// GetById updates apikey specified by link
func (s *UsersServiceOp) UpdateByLink(endpoint string, t *UserRequestUpdate) (*User, error) {
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
		return nil, ApiError{StatusCode: resp.StatusCode, Message: "non-ok status returned"}
	}
	obj := &UserResponse{}
	dec := json.NewDecoder(resp.Body)
	dec.Decode(obj)
	return s.get(obj)
}

func (s *UsersServiceOp) CreateByDirectory(id string, dir *UserRequestCreate) (*User, error) {
	endpoint := fmt.Sprintf("directories/%s/users", id)
	return s.CreateByLink(endpoint, dir)
}

// Create creates new apikey within tenant
func (s *UsersServiceOp) CreateByLink(endpoint string, dir *UserRequestCreate) (*User, error) {
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
	obj := &UserResponse{}
	dec := json.NewDecoder(resp.Body)
	dec.Decode(obj)
	return s.get(obj)
}

// Delete removes apikey
func (s *UsersServiceOp) Delete(t *User) error {
	return s.DeleteByLink(t.Href)
}

// Delete removes apikey by ID
func (s *UsersServiceOp) DeleteById(id string) error {
	endpoint := fmt.Sprintf("users/%s", id)
	return s.DeleteByLink(endpoint)
}

// Delete removes apikey by link
func (s *UsersServiceOp) DeleteByLink(endpoint string) error {
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
