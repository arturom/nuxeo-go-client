package nuxeogoclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

func readBodyToBuffer(resp *http.Response) (*bytes.Buffer, error) {
	defer resp.Body.Close()
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(resp.Body)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

func readStringResponse(resp *http.Response) (string, error) {
	buf, err := readBodyToBuffer(resp)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func unmarshallJSONResponse(resp *http.Response, v any) error {
	buf, err := readBodyToBuffer(resp)
	if err != nil {
		return err
	}
	fmt.Println(buf.String())
	return json.Unmarshal(buf.Bytes(), v)
}

type NuxeoClient interface {
	RequestAuthenticationToken(appName, deviceId, deviceDescription, permission string) (string, error)
	Users() UsersClient
	Uploads() UploadsClient
}

type client struct {
	baseUrl string
	auth    AuthMethod
}

func NewNuxeoClient(baseUrl string, auth AuthMethod) NuxeoClient {
	c := new(client)
	c.baseUrl = baseUrl
	c.auth = auth
	return c
}

func (c client) initRequest(method, path string, body io.Reader) (*http.Request, error) {
	url, err := url.JoinPath(c.baseUrl, path)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	c.auth.applyAuth(req)
	return req, nil
}

func (c client) sendRequest(req *http.Request) (*http.Response, error) {
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 400 {
		return resp, newRequestError(resp.StatusCode)
	}
	return resp, nil
}

func (c client) Get(path string) (*http.Response, error) {
	req, err := c.initRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	return c.sendRequest(req)
}

func (c client) GetJson(path string, v any) (*http.Response, error) {
	resp, err := c.Get(path)
	if err != nil {
		return nil, err
	}

	if err = unmarshallJSONResponse(resp, v); err != nil {
		return nil, err
	}

	return resp, nil
}

func (c client) Post(path string, body io.Reader) (*http.Response, error) {
	req, err := c.initRequest(http.MethodPost, path, body)
	if err != nil {
		return nil, err
	}
	return c.sendRequest(req)
}

func (c client) PostJson(path string, body io.Reader, v any) (*http.Response, error) {
	resp, err := c.Post(path, body)
	if err != nil {
		return nil, err
	}

	if err = unmarshallJSONResponse(resp, v); err != nil {
		return nil, err
	}

	return resp, nil
}

func (c client) Delete(path string) (*http.Response, error) {
	req, err := c.initRequest(http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}
	return c.sendRequest(req)
}

func (c client) DeleteJson(path string, v any) (*http.Response, error) {
	resp, err := c.Delete(path)
	if err != nil {
		return nil, err
	}

	if err = unmarshallJSONResponse(resp, v); err != nil {
		return nil, err
	}

	return resp, nil
}

func (c client) RequestAuthenticationToken(appName, deviceId, deviceDescription, permission string) (string, error) {
	req, err := c.initRequest(http.MethodGet, "/authentication/token", nil)
	if err != nil {
		return "", err
	}
	q := req.URL.Query()
	q.Add("applicationName", appName)
	q.Add("deviceId", deviceId)
	q.Add("deviceDescription", deviceDescription)
	q.Add("permission", permission)
	req.URL.RawQuery = q.Encode()
	resp, err := c.sendRequest(req)
	if err != nil {
		return "", err
	}

	return readStringResponse(resp)
}

func (c client) Users() UsersClient {
	return newUserClient(c)
}

func (c client) Uploads() UploadsClient {
	return newUploadsClient(c)
}
