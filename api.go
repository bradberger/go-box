package box

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/net/context"
)

// NewAPI returns a new API instance with the given authToken
func NewAPI(authToken string) *API {
	return &API{authToken: authToken, endpoint: "https://api.box.com/2.0"}
}

// API handles sending/receiving data to/from the MemberClicks API
type API struct {
	authToken, endpoint string
}

// SetEndpoint allows updating of the BOX api endpoint if needed for testing, etc.
func (a *API) SetEndpoint(endpointURL string) {
	a.endpoint = endpointURL
}

// Post sends a HTTP Post request to the given url with the reqData encoded as XML in the body, and returns the result
// decoded into respData.
func (a *API) PostJSON(ctx context.Context, uri string, reqData interface{}, respData interface{}) (*ErrorCodeResponse, error) {
	var body bytes.Buffer
	if err := json.NewEncoder(&body).Encode(reqData); err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/%s", a.endpoint, strings.TrimPrefix(uri, "/")), &body)
	if err != nil {
		return nil, err
	}
	return a.do(ctx, req, respData)
}

type Token struct {
	AccessToken  string        `json:"access_token"`
	ExpiresIn    int           `json:"expires_in"`
	TokenType    string        `json:"string"`
	RestrictedTo []interface{} `json:"restricted_to"`
}

// JWT gets a JWT token
func (a *API) JWT(ctx context.Context) (*Token, error) {

	var t Token
	form := url.Values{}
	form.Add("grant_type", "urn:ietf:params:oauth:grant-type:jwt-bearer")
	form.Add("assertion", JWTToken())
	form.Add("client_id", ClientID)
	form.Add("client_secret", ClientSecret)
	req, err := http.NewRequest("POST", "https://api.box.com/oauth2/token", bytes.NewBufferString(form.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	if _, err := a.do(ctx, req, &t); err != nil {
		log.Println(form)
		return nil, err
	}

	return &t, nil
}

// SetToken sets the token to be used on future requests
func (a *API) SetToken(accessToken string) {
	a.authToken = accessToken
}

func (a *API) Put(ctx context.Context, uri string, reqData interface{}, respData interface{}) (*ErrorCodeResponse, error) {
	var body bytes.Buffer
	if err := json.NewEncoder(&body).Encode(reqData); err != nil {
		return nil, err
	}
	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/%s", a.endpoint, strings.TrimPrefix(uri, "/")), &body)
	if err != nil {
		return nil, err
	}
	return a.do(ctx, req, respData)
}

// Get sends a HTTP GET request to the given uri and returns the XML encoded respone in respData
func (a *API) Get(ctx context.Context, uri string, respData interface{}) (*ErrorCodeResponse, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s", a.endpoint, strings.TrimPrefix(uri, "/")), nil)
	if err != nil {
		return nil, err
	}
	return a.do(ctx, req, respData)
}

func (a *API) do(ctx context.Context, req *http.Request, respData interface{}) (*ErrorCodeResponse, error) {

	//<base64URLencoded header>.<base64URLencoded claims>.<base64URLencoded signature>
	req.Header.Add("Authorization", "Bearer "+a.authToken)
	resp, err := getHTTPClient(ctx).Do(req)
	if err != nil {
		return nil, err
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
		if err := json.NewDecoder(resp.Body).Decode(&e); err == nil {
			e.Message = http.StatusText(resp.StatusCode)
			e.Status = resp.StatusCode
			e.Code = http.StatusText(resp.StatusCode)
		}
		return &e, fmt.Errorf("error: %v", e)
	}
	if err := json.NewDecoder(resp.Body).Decode(respData); err != nil {
		return nil, err
	}
	return nil, nil
}
