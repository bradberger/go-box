package box

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"golang.org/x/net/context"
)

type RefreshTokenResponse struct {
	AccessToken  string         `json:"access_token"`
	ExpiresIn    int64          `json:"expires_in"`
	RestrictedTo *[]interface{} `json:"restricted_to,omitempty"`
	TokenType    string         `json:"token_type"`
	RefreshToken string         `json:"refresh_token"`
}

func RefreshToken(ctx context.Context, clientID, clientSecret, refreshToken string) (*RefreshTokenResponse, *ErrorCodeResponse, error) {

	var tokenResp RefreshTokenResponse
	form := url.Values{
		"grant_type":    {"refresh_token"},
		"refresh_token": {refreshToken},
		"client_id":     {clientID},
		"client_secret": {clientSecret},
	}
	req, err := http.NewRequest("POST", "https://app.box.com/api/oauth2/token", bytes.NewBufferString(form.Encode()))
	if err != nil {
		return nil, nil, err
	}

	log.Printf("Sending: %+v\n", req)

	resp, err := getHTTPClient(ctx).Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	// Read the content
	var bodyBytes []byte
	if resp.Body != nil {
		bodyBytes, _ = ioutil.ReadAll(resp.Body)
	}

	// Restore the io.ReadCloser to its original state
	resp.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	log.Printf("Response: %+v", resp)
	log.Printf("Body: %s\n", string(bodyBytes))

	if resp.StatusCode >= 400 {
		var e ErrorCodeResponse
		if err := json.NewDecoder(resp.Body).Decode(&e); err != nil {
			e.Message = http.StatusText(resp.StatusCode)
			e.Status = resp.StatusCode
			e.Code = http.StatusText(resp.StatusCode)
		}
		return nil, &e, fmt.Errorf(e.Message)
	}
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, nil, err
	}

	return &tokenResp, nil, nil
}
