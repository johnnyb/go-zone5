package zone5

import (
	"net/url"
)

func (conn *AuthenticatedConnection) LookupUserInformationByDetails(data url.Values) (map[string]interface{}, error) {
	req, err := conn.NewRequest("GET", "/rest/users?" + data.Encode(), nil)
	if err != nil {
		return nil, err
	}
	resp, err := conn.Connection.HTTPClient.Do(req)
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
	return resultData, nil
}

func (conn *AuthenticatedConnection) LookupUserInformationByEmail(email string)  (map[string]interface{}, error) {
	return conn.LookupUserInformationByDetails(url.Values{
		"email": {email},
	})
}

func (conn *AuthenticatedConnection) LookupUserInformationByIdentifier(ident string) (map[string]interface{}, error) {
	return conn.LookupUserInformationByDetails(url.Values{
		"uid": {ident},
	})

}

func (conn *AuthenticatedConnection) GetTokenForUserIdentifier(ident string) (string, error) {
	req, err := conn.NewRequest("GET", "/rest/users/jwt/" + conn.Connection.ApiKey + "/" + ident + "/7200", nil)
	if err != nil {
		return "", err
	}
	resp, err := conn.Connection.HTTPClient.Do(req)
	if err != nil {
		return "", err
	}

	resultData, err := unmarshalJsonFromReader(resp.Body)
	if err != nil {
		return "", err
	}

	return resultData["token"].(string), nil
}
