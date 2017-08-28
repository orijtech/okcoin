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
	"net/http"
	"sync"
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
