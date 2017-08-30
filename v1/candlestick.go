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

type CandleStickRequest struct {
	N int `json:"n,omitempty"`

	Since  float64 `json:"since,omitempty"`
	Symbol Symbol  `json:"sym,omitempty"`

	Period Period `json:"period,omitempty"`
}

type Period string

const (
	P1Min   Period = "1min"
	P3Min   Period = "3min"
	P5Min   Period = "5min"
	P15Min  Period = "15min"
	P30Min  Period = "30min"
	P1Day   Period = "1day"
	P3Day   Period = "3day"
	P1Week  Period = "1week"
	P1Hour  Period = "1hour"
	P2Hour  Period = "2hour"
	P4Hour  Period = "4hour"
	P6Hour  Period = "6hour"
	P12Hour Period = "12hour"
)

type candleStickRequest struct {
	Symbol Symbol  `json:"symbol,omitempty"`
	Period Period  `json:"type,omitempty"`
	Since  float64 `json:"since,omitempty"`
	N      int     `json:"size,omitempty"`
}

const defaultPeriod = P1Hour

func (c *Client) CandleStick(csr *CandleStickRequest) (*CandleStickResponse, error) {
	if csr == nil {
		csr = new(CandleStickRequest)
	}
	symbol := csr.Symbol
	if symbol == "" {
		symbol = defaultSymbol
	}
	period := csr.Period
	if period == "" {
		period = defaultPeriod
	}

	since := float64(0)
	if csr.Since <= 0 {
		since = 0
	}

	creq := &candleStickRequest{
		Symbol: symbol,
		N:      csr.N,
		Since:  since,
		Period: period,
	}
	qv, err := otils.ToURLValues(creq)
	if err != nil {
		return nil, err
	}
	fullURL := fmt.Sprintf("%s/kline.do", baseURL)
	if len(qv) > 0 {
		fullURL = fmt.Sprintf("%s?%s", fullURL, qv.Encode())
	}
	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return nil, err
	}
	blob, _, err := c.doHTTPReq(req)
	if err != nil {
		return nil, err
	}
	var recv []*CandleStick
	if err := json.Unmarshal(blob, &recv); err != nil {
		return nil, err
	}
	cres := &CandleStickResponse{
		CandleSticks: recv,
		Symbol:       symbol,
		Since:        since,
		Period:       period,
	}

	return cres, nil
}

type CandleStickResponse struct {
	Symbol Symbol  `json:"symbol,omitempty"`
	Since  float64 `json:"since,omitempty"`
	Period Period  `json:"period,omitempty"`

	CandleSticks []*CandleStick `json:"candle_sticks,omitempty"`
}

type CandleStick struct {
	TimeStampMs float64 `json:"timestamp,omitempty"`
	Open        float64 `json:"open,omitempty"`
	High        float64 `json:"high,omitempty"`
	Low         float64 `json:"low,omitempty"`
	Close       float64 `json:"close,omitempty"`
	Volume      float64 `json:"volume,omitempty"`
}

const (
	rawCandleStickFieldCount = 6
)

func (cs *CandleStick) UnmarshalJSON(b []byte) error {
	// Expecting a datum of the form:
	// [
	//  1417564800000,  timestamp
	//  384.47,   open
	//  387.13,   high
	//  383.5,    low
	//  387.13,   close
	//  1062.04,  volume
	// ]
	recv := make([]float64, rawCandleStickFieldCount)
	if err := json.Unmarshal(b, &recv); err != nil {
		return err
	}

	if g, w := len(recv), rawCandleStickFieldCount; g < w {
		return fmt.Errorf("fields: got %d want %d; data=%s", g, w, b)
	}

	cs.TimeStampMs = recv[0]
	cs.Open = recv[1]
	cs.High = recv[2]
	cs.Low = recv[3]
	cs.Close = recv[4]
	cs.Volume = recv[5]
	return nil
}
