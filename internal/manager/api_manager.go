package manager

import (
	"bytes"
	"fmt"
	"net/http"
	"time"
)

// TODO Add error handling to the APIManager

// APIManager handles API requests with proper authentication.
const guestToken = "0fc237962e95129004c313015d220aef4c7ffddc465cf984d1e63130b6e180c8"

type APIManager struct {
	client     *http.Client
	guestToken string
	userToken  string
	isUserAuth bool
}

// NewAPIManager initializes the API manager with a guest token.
func NewAPIManager() *APIManager {
	return &APIManager{
		client:     &http.Client{Timeout: 10 * time.Second},
		guestToken: guestToken,
		isUserAuth: false,
	}
}

// NewAPIManagerWithToken initializes the API manager with a guest token and optionally a user token.
func NewAPIManagerWithToken(userToken string) *APIManager {
	manager := &APIManager{
		client:     &http.Client{Timeout: 10 * time.Second},
		guestToken: guestToken,
		isUserAuth: false,
	}

	if userToken != "" {
		manager.SetUserToken(userToken)
	}

	return manager
}

// SetUserToken sets the user token and marks the user as authenticated.
func (a *APIManager) SetUserToken(token string) {
	a.userToken = token
	a.isUserAuth = true
}

func (a *APIManager) Request(method, url string, body []byte) (*http.Response, error) {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	if a.isUserAuth {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", a.userToken))
	} else {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", a.guestToken))
	}

	req.Header.Set("Content-Type", "application/json")

	return a.client.Do(req)
}

func (a *APIManager) Get(url string) (*http.Response, error) {
	return a.Request("GET", url, nil)
}

// Post performs a POST request.
func (a *APIManager) Post(url string, body []byte) (*http.Response, error) {
	return a.Request("POST", url, body)
}

// Put performs a PUT request.
func (a *APIManager) Put(url string, body []byte) (*http.Response, error) {
	return a.Request("PUT", url, body)
}

// Delete performs a DELETE request.
func (a *APIManager) Delete(url string) (*http.Response, error) {
	return a.Request("DELETE", url, nil)
}

func (a *APIManager) GetGuestToken() string {
	return guestToken
}
