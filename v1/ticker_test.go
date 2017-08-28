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

package okcoin_test

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/orijtech/okcoin/v1"
)

var blankTickerResponse = new(okcoin.TickerResponse)

func TestTicker(t *testing.T) {
	client, err := okcoin.NewDefaultClient()
	if err != nil {
		t.Fatalf("new client: %v", err)
	}

	tests := [...]struct {
		symbol  okcoin.Symbol
		wantErr bool
	}{
		0: {wantErr: true},
		1: {symbol: okcoin.LTCUSD},
		2: {symbol: "fugazi-coin", wantErr: true},
	}

	client.SetHTTPRoundTripper(&backend{route: tickerRoute})
	for i, tt := range tests {
		tres, err := client.Ticker(tt.symbol)
		if tt.wantErr {
			if err == nil {
				t.Errorf("#%d: want non-nil; got resp=%#v", i, tres)
			}
			continue
		}

		if err != nil {
			t.Errorf("#%d: got err: %v", i, err)
			continue
		}
		if reflect.DeepEqual(tres, blankTickerResponse) {
			t.Errorf("#%d: got blankTickerResponse", i)
		}
	}
}

type backend struct {
	route string
}

var _ http.RoundTripper = (*backend)(nil)

var (
	errUnimplemented = errors.New("unimplemented")
)

func (b *backend) RoundTrip(req *http.Request) (*http.Response, error) {
	switch b.route {
	case tickerRoute:
		return b.tickerRoundTrip(req)
	default:
		return nil, errUnimplemented
	}
}

func (b *backend) tickerRoundTrip(req *http.Request) (*http.Response, error) {
	if req.Method != "GET" {
		return makeResp(fmt.Sprintf(`got method %q want "GET"`, req.Method), http.StatusMethodNotAllowed, nil)
	}
	wantPathSuffix := "/api/v1/ticker.do"
	gotPath := req.URL.Path
	if !strings.HasSuffix(gotPath, wantPathSuffix) {
		return makeResp(fmt.Sprintf(`got suffix %q want %q`, gotPath, wantPathSuffix), http.StatusBadRequest, nil)
	}
	query := req.URL.Query()
	symbol := query.Get("symbol")
	outPath := fmt.Sprintf("./testdata/ticker-%s.json", symbol)
	f, err := os.Open(outPath)
	if err != nil {
		return makeResp(err.Error(), http.StatusInternalServerError, nil)
	}
	// This handle should be closed by the consumer.
	return makeResp("200 OK", http.StatusOK, f)
}

func makeResp(status string, statusCode int, body io.ReadCloser) (*http.Response, error) {
	resp := &http.Response{
		Status:     status,
		StatusCode: statusCode,
		Body:       body,
		Header:     make(http.Header),
	}
	return resp, nil
}

// Route declarations
const (
	tickerRoute = "/ticker"
)
