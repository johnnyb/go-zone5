package zone5

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

var RefreshTimeMargin = time.Duration(30) * time.Second

type AuthenticatedConnection struct {
	Connection   *Connection
	AuthToken    string
	RefreshToken string
	Expiration   *time.Time
	Details      map[string]interface{}
}

func (conn *AuthenticatedConnection) GetCurrentToken() string {
	if time.Now().After(*conn.Expiration) {
		conn.PerformTokenRefresh()
	}

	return conn.AuthToken
}

func (conn *AuthenticatedConnection) PerformTokenRefresh() error {
	req, err := conn.Connection.NewRequest("POST", "/rest/oauth/access_token", strings.NewReader(url.Values{
		"username":      {conn.Details["user"].(map[string]interface{})["email"].(string)},
		"client_id":     {conn.Connection.ApiKey},
		"client_secret": {conn.Connection.ApiKeySecret},
		"grant_type":    {"refresh_token"},
		"refresh_token": {conn.RefreshToken},
	}.Encode()))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := conn.Connection.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return NewErrorWithResponse(resp)
	}

	resultData, err := unmarshalJsonFromReader(resp.Body)
	if err != nil {
		return err
	}

	// Expiration time, but give 30 seconds for potential errors
	exp := time.Now().Add(time.Duration(resultData["expires_in"].(float64))*time.Second - RefreshTimeMargin)
	tok := resultData["access_token"].(string)

	conn.Expiration = &exp
	conn.AuthToken = tok
	return nil
}

func (conn *AuthenticatedConnection) NewRequest(method, path string, body io.Reader) (*http.Request, error) {
	var version string
	version = os.Getenv("TASKBOT_VERSION")
	if len(version) == 0 {
		version = "unknown"
	}
	var lambda string
	lambda = os.Getenv("AWS_LAMBDA_FUNCTION_NAME")
	if len(lambda) == 0 {
		version = "taskbot"
	}
	userAgent := fmt.Sprintf("%s/%s", lambda, version)

	req, err := conn.Connection.NewRequest(method, path, body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+conn.GetCurrentToken())
	req.Header.Set("User-Agent", userAgent)
	return req, nil
}
