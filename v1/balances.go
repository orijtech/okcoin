// Copyright 2017 orijtech, Inc. All Rights Reserved.
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
	"net/http"
	"reflect"
)

type Funds struct {
	Asset     *Fund `json:"asset,omitempty"`
	Borrow    *Fund `json:"borrow,omitempty"`
	Free      *Fund `json:"free,omitempty"`
	Frozen    *Fund `json:"freezed,omitempty"`
	UnionFund *Fund `json:"union_fund,omitempty"`
}

type Fund struct {
	Net   float64 `json:"net,string,omitempty"`
	Total float64 `json:"total,string,omitempty"`
	BTC   float64 `json:"btc,string,omitempty"`
	ETH   float64 `json:"eth,string,omitempty"`
	LTC   float64 `json:"ltc,string,omitempty"`
	USD   float64 `json:"usd,string,omitempty"`
}

type fundsIntermediate struct {
	Funds  map[string]*Funds `json:"info"`
	Result bool              `json:"result"`
}

var errNoFundsReturned = errors.New("no funds information returned")
var blankFunds = new(Funds)

func (c *Client) Funds() (*Funds, error) {
	fullURL := fmt.Sprintf("%s/userinfo.do", baseURL)
	req, err := http.NewRequest("POST", fullURL, nil)
	if err != nil {
		return nil, err
	}
	qv := req.URL.Query()
	qv.Set("api_key", c.apiKey())
	qv, err = c.prepareSignedAuthBody(qv)
	if err != nil {
		return nil, err
	}
	req.URL.RawQuery = qv.Encode()
	blob, _, err := c.doHTTPReq(req)
	if err != nil {
		return nil, err
	}
	fi := new(fundsIntermediate)
	if err := json.Unmarshal(blob, fi); err != nil {
		return nil, err
	}
	if !fi.Result || len(fi.Funds) == 0 {
		return nil, errNoFundsReturned
	}
	funds := fi.Funds["funds"]
	if reflect.DeepEqual(funds, blankFunds) {
		return nil, errNoFundsReturned
	}
	return funds, nil
}
