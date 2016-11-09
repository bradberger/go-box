package box

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
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
func (a *API) Post(ctx context.Context, uri string, reqData interface{}, respData interface{}) (*ErrorCodeResponse, error) {
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
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", a.authToken))
	log.Printf("Request: %+v", req)
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
		if err := json.NewDecoder(resp.Body).Decode(&e); err != nil {
			e.Message = http.StatusText(resp.StatusCode)
			e.Status = resp.StatusCode
			e.Code = http.StatusText(resp.StatusCode)
		}
		return &e, fmt.Errorf(e.Message)
	}
	if err := json.NewDecoder(resp.Body).Decode(respData); err != nil {
		return nil, err
	}
	return nil, nil
}
