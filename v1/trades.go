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
	"fmt"
	"net/http"

	"github.com/orijtech/otils"
)

type LastTradesRequest struct {
	LastTradeID int    `json:"id,omitempty"`
	Symbol      Symbol `json:"symbol,omitempty"`
}

func (ltr *LastTradesRequest) Validate() error {
	if ltr == nil || ltr.Symbol == "" {
		return errBlankSymbol
	}
	return nil
}

const (
	defaultSymbol = BTCUSD
)

type lastTradesReq struct {
	LastTradeID int    `json:"since,omitempty"`
	Symbol      Symbol `json:"symbol,omitempty"`
}

type LastTradesResponse struct {
	Trades []*Trade `json:"trades,omitempty"`
	Symbol Symbol   `json:"symbol,omitempty"`
}

type Trade struct {
	Amount float64 `json:"amount,string,omitempty"`
	Type   string  `json:"type,omitempty"`
	ID     int64   `json:"tid,omitempty"`
	Price  float64 `json:"price,string,omitempty"`
	DateMs float64 `json:"date_ms,omitempty"`
	Date   float64 `json:"date,omitempty"`
}

func (c *Client) LastTrades(ltr *LastTradesRequest) (*LastTradesResponse, error) {
	if ltr == nil {
		ltr = new(LastTradesRequest)
	}
	lastTradeID := ltr.LastTradeID
	if lastTradeID <= 0 {
		lastTradeID = 0
	}
	symbol := ltr.Symbol
	if symbol == "" {
		symbol = defaultSymbol
	}
	qv, err := otils.ToURLValues(&lastTradesReq{
		LastTradeID: lastTradeID,
		Symbol:      symbol,
	})
	if err != nil {
		return nil, err
	}
	fullURL := fmt.Sprintf("%s/trades.do?%s", baseURL, qv.Encode())
	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return nil, err
	}
	blob, _, err := c.doHTTPReq(req)
	if err != nil {
		return nil, err
	}
	var trades []*Trade
	if err := json.Unmarshal(blob, &trades); err != nil {
		return nil, err
	}
	ltres := &LastTradesResponse{Trades: trades, Symbol: symbol}
	return ltres, nil
}
