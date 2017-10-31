// Copyright 2017 orijtech. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package okcoin

import (
	"crypto/md5"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"

	"github.com/orijtech/otils"
)

const (
	baseURL = "https://www.okcoin.com/api/v1"
)

type Symbol string

const (
	BCCUSD Symbol = "bcc_usd"
	BTCUSD Symbol = "btc_usd"
	ETCUSD Symbol = "etc_usd"
	ETHUSD Symbol = "eth_usd"
	LTCUSD Symbol = "ltc_usd"
)

type Client struct {
	rt http.RoundTripper
	mu sync.RWMutex

	_apiSecret string
	_apiKey    string
}

const (
	envAPIKeyKey    = "OKCOIN_API_KEY"
	envAPISecretKey = "OKCOIN_API_SECRET"
)

func NewClientFromEnv() (*Client, error) {
	var errsList []string
	apiKey := fromEnvOrAppendError(envAPIKeyKey, &errsList)
	apiSecret := fromEnvOrAppendError(envAPISecretKey, &errsList)
	if len(errsList) > 0 {
		return nil, errors.New(strings.Join(errsList, "\n"))
	}
	return &Client{_apiSecret: apiSecret, _apiKey: apiKey}, nil
}

func fromEnvOrAppendError(envKey string, errsList *[]string) string {
	if value := os.Getenv(envKey); value != "" {
		return value
	}
	*errsList = append(*errsList, fmt.Sprintf("%q was not set", envKey))
	return ""
}

func (c *Client) SetHTTPRoundTripper(rt http.RoundTripper) {
	c.mu.Lock()
	c.rt = rt
	c.mu.Unlock()
}

func (c *Client) httpClient() *http.Client {
	c.mu.RLock()
	rt := c.rt
	c.mu.RUnlock()

	return &http.Client{
		Transport: rt,
	}
}

func NewDefaultClient() (*Client, error) {
	return new(Client), nil
}

func (c *Client) doHTTPReq(req *http.Request) ([]byte, http.Header, error) {
	res, err := c.httpClient().Do(req)
	if err != nil {
		return nil, nil, err
	}
	if res.Body != nil {
		defer res.Body.Close()
	}
	if !otils.StatusOK(res.StatusCode) {
		return nil, res.Header, fmt.Errorf("%s %d", res.Status, res.StatusCode)
	}

	blob, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, res.Header, err
	}
	return blob, res.Header, nil
}

func (c *Client) prepareSignedAuthBody(qv url.Values) (url.Values, error) {
	// As per https://www.okcoin.com/intro_signParams.html
	// Get the query string + "secret_key"=$SECRET_KEY
	// * All parameters/queries except "sign" must be signed
	// * Must use an MD5 signature with apiSecret as the key and the query values
	// * Send the "sign" query value along
	// * Sign's value MUST be an upper-case string
	h := md5.New()
	fmt.Fprintf(h, "%s&secret_key=%s", qv.Encode(), c.apiSecret())
	signature := fmt.Sprintf("%x", h.Sum(nil))
	qv.Set("sign", strings.ToUpper(signature))
	return qv, nil
}

type Credentials struct {
	APIKey string `json:"api_key"`
	Secret string `json:"secret"`
}

func (c *Client) SetCredentials(creds *Credentials) {
	if creds == nil {
		return
	}
	c.mu.Lock()
	c._apiKey = creds.APIKey
	c._apiSecret = creds.Secret
	c.mu.Unlock()
}

func (c *Client) apiSecret() string {
	c.mu.Lock()
	secretKey := c._apiSecret
	c.mu.Unlock()
	return secretKey
}

func (c *Client) apiKey() string {
	c.mu.Lock()
	apiKey := c._apiKey
	c.mu.Unlock()
	return apiKey
}
