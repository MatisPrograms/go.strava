package strava

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

// An OAuthAuthenticator holds state about how OAuth requests should be authenticated.
type OAuthAuthenticator struct {
	CallbackURL string // used to help generate the AuthorizationURL

	// The RequestClientGenerator builds the http.Client that will be used
	// to complete the token exchange. If nil, http.DefaultClient will be used.
	// On Google's App Engine http.DefaultClient is not available and this generator
	// can be used to create a client using the incoming request, for Example:
	//    func(r *http.Request) { return urlfetch.Client(appengine.NewContext(r)) }
	RequestClientGenerator func(r *http.Request) *http.Client
}

// Permission represents the access of an access_token.
// The permission type is requested during the token exchange.
type Permission string

// Permissions defines the available permissions
var Permissions = struct {
	Public           Permission
	ViewPrivate      Permission
	Write            Permission
	WriteViewPrivate Permission
	All              Permission
}{
	"public",
	"view_private",
	"write",
	"write,view_private",
	"read,read_all,profile:read_all,profile:write,activity:read,activity:read_all,activity:write",
}

// AuthorizationResponse is returned as a result of the token exchange
type AuthorizationResponse struct {
	AccessToken  string         `json:"access_token"`
	RefreshToken string         `json:"refresh_token"`
	ExpiresAt    int64          `json:"expires_at"`
	ExpiresIn    int64          `json:"expires_in"`
	Athlete      AthleteSummary `json:"athlete"`
	State        string         `json:"State"`
}

// StravaError represents an error response from the Strava API.
type StravaError struct {
	Code     string `json:"code"`
	Field    string `json:"field"`
	Resource string `json:"resource"`
}

// StravaErrorResponse represents an error response from the Strava API.
type StravaErrorResponse struct {
	Errors  []StravaError `json:"errors"`
	Message string        `json:"message"`
}

// CallbackPath returns the path portion of the CallbackURL.
// Useful when setting a http path handler, for example:
//
//	http.HandleFunc(stravaOAuth.CallbackURL(), stravaOAuth.HandlerFunc(successCallback, failureCallback))
func (auth OAuthAuthenticator) CallbackPath() (string, error) {
	if auth.CallbackURL == "" {
		return "", errors.New("callbackURL is empty")
	}
	url, err := url.Parse(auth.CallbackURL)
	if err != nil {
		return "", err
	}
	return url.Path, nil
}

// ExchangeToken handles the common logic for token exchange with the Strava API
func ExchangeToken(values url.Values) (*AuthorizationResponse, *http.Response, error) {
	// Append client_id and client_secret to the request
	values.Set("client_id", fmt.Sprintf("%d", ClientId))
	values.Set("client_secret", ClientSecret)

	// Send the request to the Strava API
	resp, err := http.DefaultClient.PostForm(basePath+"/oauth/token", values)

	// this was a poor request, maybe strava servers down?
	if err != nil {
		return nil, resp, err
	}
	defer resp.Body.Close()

	// Read the response body
	contents, _ := io.ReadAll(resp.Body)

	// if status code is not 200, then something went wrong
	if resp.StatusCode/100 != 2 {
		var stravaErr StravaErrorResponse
		json.Unmarshal(contents, &stravaErr)
		return nil, resp, errors.New("Strava API error")
	}

	// Parse the response body
	var response AuthorizationResponse
	err = json.Unmarshal(contents, &response)
	if err != nil {
		return nil, resp, err
	}

	// Return the data, response, and error
	return &response, resp, nil
}

// Authorize performs the second part of the OAuth exchange. The client has already been redirected to the
// Strava authorization page, has granted authorization to the application and has been redirected back to the
// defined URL. The code param was returned as a query string param in to the redirect_url.
func (auth OAuthAuthenticator) Authorize(code string, client *http.Client) (*AuthorizationResponse, error) {
	// make sure a code was passed
	if code == "" {
		return nil, OAuthInvalidCodeErr
	}

	// if a client wasn't passed use the default client
	if client == nil {
		client = http.DefaultClient
	}

	resp, err := client.PostForm(basePath+"/oauth/token",
		url.Values{"client_id": {fmt.Sprintf("%d", ClientId)}, "client_secret": {ClientSecret}, "code": {code}})

	// this was a poor request, maybe strava servers down?
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// check status code, could be 500, or most likely the client_secret is incorrect
	if resp.StatusCode/100 == 5 {
		return nil, OAuthServerErr
	}

	if resp.StatusCode/100 != 2 {
		var response Error
		contents, _ := ioutil.ReadAll(resp.Body)
		json.Unmarshal(contents, &response)

		if len(response.Errors) == 0 {
			return nil, OAuthServerErr
		}

		if response.Errors[0].Resource == "Application" {
			return nil, OAuthInvalidCredentialsErr
		}

		if response.Errors[0].Resource == "RequestToken" {
			return nil, OAuthInvalidCodeErr
		}

		return nil, &response
	}

	var response AuthorizationResponse
	contents, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(contents, &response)

	if err != nil {
		return nil, err
	}

	return &response, nil
}

// HandlerFunc builds a http.HandlerFunc that will complete the token exchange
// after a user authorizes an application on strava.com.
// This method handles the exchange and calls success or failure after it completes.
func (auth OAuthAuthenticator) HandlerFunc(
	success func(auth *AuthorizationResponse, w http.ResponseWriter, r *http.Request),
	failure func(err error, w http.ResponseWriter, r *http.Request)) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		// user denied authorization
		if r.FormValue("error") == "access_denied" {
			failure(OAuthAuthorizationDeniedErr, w, r)
			return
		}

		// use the client generator if provided.
		client := http.DefaultClient
		if auth.RequestClientGenerator != nil {
			client = auth.RequestClientGenerator(r)
		}

		resp, err := auth.Authorize(r.FormValue("code"), client)

		if err != nil {
			failure(err, w, r)
			return
		}

		resp.State = r.FormValue("state")

		success(resp, w, r)
	}
}

// AuthorizationURL constructs the url a user should use to authorize this specific application.
func (auth OAuthAuthenticator) AuthorizationURL(state string, scope Permission, force bool) string {
	path := fmt.Sprintf("%s/oauth/authorize?client_id=%d&response_type=code&redirect_uri=%s&scope=%v", basePath, ClientId, auth.CallbackURL, scope)

	if state != "" {
		path += "&state=" + state
	}

	if force {
		path += "&approval_prompt=force"
	}

	return path
}

// AuthorizationURL returns the url of the authorization endpoint.
func AuthorizationURL(callbackURL string, scope Permission) (string, error) {
	if ClientId == 0 || ClientId < 0 {
		return "", errors.New("clientId is empty")
	}
	if callbackURL == "" {
		return "", errors.New("callbackURL is empty")
	}
	if scope == "" {
		return "", errors.New("scope is empty")
	}
	return fmt.Sprintf("%s/oauth/authorize?client_id=%d&response_type=code&redirect_uri=%s&scope=%v", basePath, ClientId, callbackURL, scope), nil
}

/*********************************************************/

type OAuthService struct {
	client *Client
}

func NewOAuthService(client *Client) *OAuthService {
	return &OAuthService{client}
}

type OAuthDeauthorizeCall struct {
	service *OAuthService
}

func (s *OAuthService) Deauthorize() *OAuthDeauthorizeCall {
	return &OAuthDeauthorizeCall{
		service: s,
	}
}

func (c *OAuthDeauthorizeCall) Do() error {
	_, err := c.service.client.run("POST", "/oauth/deauthorize", nil)
	return err
}
