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
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"reflect"

	"github.com/orijtech/otils"
)

type TickerResponse struct {
	TimeAtEpoch float64 `json:"date,string"`

	Ticker *Ticker `json:"ticker"`
}

type Ticker struct {
	Buy    float64 `json:"buy,string,omitempty"`
	High   float64 `json:"high,string,omitempty"`
	Last   float64 `json:"last,string,omitempty"`
	Low    float64 `json:"low,string,omitempty"`
	Sell   float64 `json:"sell,string,omitempty"`
	Volume float64 `json:"vol,string,omitempty"`
}

var (
	errBlankSymbol         = errors.New("expecting a non-blank symbol")
	errBlankTickerResponse = errors.New("expecting a non-blank ticker response")
)

var (
	blankTickerResponse = new(TickerResponse)
)

func (c *Client) Ticker(sym Symbol) (*TickerResponse, error) {
	if sym == "" {
		return nil, errBlankSymbol
	}
	fullURL := fmt.Sprintf("%s/ticker.do?symbol=%s", baseURL, sym)
	res, err := c.httpClient().Get(fullURL)
	if err != nil {
		return nil, err
	}
	if res.Body != nil {
		defer res.Body.Close()
	}
	if !otils.StatusOK(res.StatusCode) {
		return nil, fmt.Errorf("%s", res.StatusCode)
	}

	blob, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	tRes := new(TickerResponse)
	if err := json.Unmarshal(blob, tRes); err != nil {
		return nil, err
	}
	if reflect.DeepEqual(tRes, blankTickerResponse) {
		return nil, errBlankTickerResponse
	}
	return tRes, nil
}
