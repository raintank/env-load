package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/grafana/grafana/pkg/api/dtos"
	m "github.com/grafana/grafana/pkg/models"
)

type Client struct {
	key     string
	baseURL url.URL
	*http.Client
}

//New creates a new grafana client
//auth can be in user:pass format, or it can be an api key
func New(auth, baseURL string) (*Client, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}
	key := ""
	if strings.Contains(auth, ":") {
		split := strings.Split(auth, ":")
		u.User = url.UserPassword(split[0], split[1])
	} else {
		key = fmt.Sprintf("Bearer %s", auth)
	}
	return &Client{
		key,
		*u,
		&http.Client{},
	}, nil
}

func (c *Client) newRequest(method, path string, body io.Reader) (*http.Request, error) {
	url := c.baseURL
	url.Path = path
	req, err := http.NewRequest(method, url.String(), body)
	if err != nil {
		return req, err
	}
	if c.key != "" {
		req.Header.Add("Authorization", c.key)
	}
	if body == nil {
		fmt.Println("request to ", url.String(), "with no body data")
	} else {
		fmt.Println("request to ", url.String(), "with body data", body.(*bytes.Buffer).String())
	}
	req.Header.Add("Content-Type", "application/json")
	return req, err
}

func (c *Client) Orgs() ([]Org, error) {
	orgs := make([]Org, 0)

	req, err := c.newRequest("GET", "/api/orgs/", nil)
	if err != nil {
		return orgs, err
	}
	resp, err := c.Do(req)
	if err != nil {
		return orgs, err
	}
	if resp.StatusCode != 200 {
		return orgs, errors.New(resp.Status)
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return orgs, err
	}
	err = json.Unmarshal(data, &orgs)
	return orgs, err
}

func (c *Client) Collectors() ([]Collector, error) {
	collectors := make([]Collector, 0)

	req, err := c.newRequest("GET", "/api/collectors/", nil)
	if err != nil {
		return collectors, err
	}
	resp, err := c.Do(req)
	if err != nil {
		return collectors, err
	}
	if resp.StatusCode != 200 {
		return collectors, errors.New(resp.Status)
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return collectors, err
	}
	err = json.Unmarshal(data, &collectors)
	return collectors, err
}

func (c *Client) NewOrg(name string) error {
	settings := map[string]string{
		"name": name,
	}
	data, err := json.Marshal(settings)
	req, err := c.newRequest("POST", "/api/orgs", bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return errors.New(resp.Status)
	}
	return err
}

func (c *Client) NewEndpoint(settings m.AddEndpointCommand) error {
	data, err := json.Marshal(settings)
	req, err := c.newRequest("PUT", "/api/endpoints", bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	data, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	fmt.Println("response BODY", string(data))
	if resp.StatusCode != 200 {
		return errors.New(resp.Status)
	}
	return err
}

func (c *Client) CreateUserForm(settings dtos.AdminCreateUserForm) error {
	data, err := json.Marshal(settings)
	req, err := c.newRequest("POST", "/api/admin/users", bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	data, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	fmt.Println("response BODY", string(data))
	if resp.StatusCode != 200 {
		return errors.New(resp.Status)
	}
	return err
}
