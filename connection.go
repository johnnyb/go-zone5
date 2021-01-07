package zone5

import (
	"net/http"
	"time"
	"io"
)

type Connection struct {
	HTTPClient *http.Client
	BaseURL string
	ApiKey string
	ApiKeySecret string
}

func NewConnection(key, secret string) *Connection {
	conn := Connection{
		HTTPClient: DefaultHTTPClient,
		BaseURL: "https://api-sp.zone5cloud.com",
		ApiKey: key,
		ApiKeySecret: secret,
	}

	return &conn
}

func (conn *Connection) NewRequest(method, path string, body io.Reader) (*http.Request, error) {
	url := conn.URLForPath(path)
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("tp-nodecorate", "true")
	req.Header.Add("Api-Key", conn.ApiKey)
	req.Header.Add("Api-Key-Secret", conn.ApiKeySecret)
	return req, nil
}

func (conn *Connection) Login(username, password string) (*AuthenticatedConnection, error) {
	data := map[string]interface{}{
		"username": username,
		"password": password,
		"clientId": conn.ApiKey,
		"clientSecret": conn.ApiKeySecret,
	}
	req, err := conn.NewRequest("POST", "/rest/auth/login", mustJsonBody(data))
	if err != nil {
		return nil, err
	}
	req.Header.Add("tp-nodecorate", "true")
	req.Header.Add("Api-Key", conn.ApiKey)
	req.Header.Add("Api-Key-Secret", conn.ApiKeySecret)
	req.Header.Add("Content-Type", "application/json")

	resp, err := conn.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, NewErrorWithResponse(resp)
	}

	resultData, err := unmarshalJsonFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	exp := time.Unix(int64(resultData["tokenExp"].(float64) / 1000.0), 0)
	authConn := AuthenticatedConnection{
		Connection: conn,
		AuthToken: resultData["token"].(string),
		RefreshToken: resultData["refresh"].(string),
		Expiration: &exp,
		Details: resultData,
	}

	return &authConn, nil
}

func (conn *Connection) URLForPath(path string) string {
	return conn.BaseURL + path
}

