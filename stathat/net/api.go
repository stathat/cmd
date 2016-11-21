package net

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/stathat/cmd/stathat/config"
)

type ExportAPI interface {
	Get(path string, params url.Values, dest interface{}) error
	Delete(path string, params url.Values, dest interface{}) error
	Post(path string, params url.Values, dest interface{}) error
}

var DefaultAPI ExportAPI

func init() {
	DefaultAPI = NewAPI()
}

type API struct{}

func NewAPI() *API {
	return &API{}
}

func (a *API) do(r *http.Request, dest interface{}) error {
	r.SetBasicAuth(config.AccessKey(), "")
	if config.Debug("api") {
		fmt.Printf("http %s -> %s\n", r.Method, r.URL)
	}
	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if config.Debug("api") {
		fmt.Printf("body type: %T\n", resp.Body)
		fmt.Printf("http %s -> %s -> response: %+v\n", r.Method, r.URL, resp)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("http error: %d %s", resp.StatusCode, http.StatusText(resp.StatusCode))
	}

	if dest == nil {
		return nil
	}

	dec := json.NewDecoder(resp.Body)
	return dec.Decode(dest)
}

func (a *API) Get(path string, params url.Values, dest interface{}) error {
	req, err := getReq("GET", path, params)
	if err != nil {
		return err
	}
	return a.do(req, dest)
}

func (a *API) Delete(path string, params url.Values, dest interface{}) error {
	req, err := getReq("DELETE", path, params)
	if err != nil {
		return err
	}
	return a.do(req, dest)
}

func (a *API) Post(path string, params url.Values, dest interface{}) error {
	req, err := formReq("POST", path, params)
	if err != nil {
		return err
	}
	return a.do(req, dest)
}

func (a *API) Put(path string, params url.Values, dest interface{}) error {
	req, err := formReq("PUT", path, params)
	if err != nil {
		return err
	}
	return a.do(req, dest)
}

func getReq(method, path string, vals url.Values) (*http.Request, error) {
	r, err := http.NewRequest(method, "", nil)
	if err != nil {
		return nil, fmt.Errorf("http.NewRequest error: %s", err)
	}
	r.URL, err = newURL(path)
	if err != nil {
		return nil, fmt.Errorf("newURL error: %s", err)
	}
	r.URL.RawQuery = vals.Encode()
	return r, nil
}

func formReq(method, path string, params url.Values) (*http.Request, error) {
	pr := strings.NewReader(params.Encode())
	r, err := http.NewRequest(method, "", pr)
	if err != nil {
		return nil, err
	}
	r.URL, err = newURL(path)
	if err != nil {
		return nil, err
	}
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return r, nil
}

func newURL(path string) (*url.URL, error) {
	ur, err := url.Parse(config.Host())
	if err != nil {
		return nil, err
	}
	return &url.URL{
		Scheme: ur.Scheme,
		Host:   ur.Host,
		Path:   "/a1/" + path,
	}, nil
}
