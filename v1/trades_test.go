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
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/orijtech/okcoin/v1"
)

func TestLastTrades(t *testing.T) {
	t.Parallel()

	client, err := okcoin.NewDefaultClient()
	if err != nil {
		t.Fatalf("new client: %v", err)
	}

	tests := [...]struct {
		req         *okcoin.LastTradesRequest
		wantAtLeast int
		wantSymbol  okcoin.Symbol
		wantErr     bool
	}{
		0: {wantAtLeast: 20, wantSymbol: okcoin.BTCUSD}, // a nil request should return the defaults
		1: {
			req: &okcoin.LastTradesRequest{
				Symbol: okcoin.LTCUSD,
				N:      60,
			},
			wantAtLeast: 20,
		},
	}

	client.SetHTTPRoundTripper(&backend{route: lastNTradesRoute})
	for i, tt := range tests {
		tres, err := client.LastTrades(tt.req)
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
		if tres == nil {
			t.Errorf("#%d: expected a non-nil trades response", i)
			continue
		}
		if g, w := len(tres.Trades), tt.wantAtLeast; g < w {
			t.Errorf("#%d: got=%d, wantAtLeast=%d", i, g, w)
		}
	}
}

func (b *backend) lastNTradesRoundTrip(req *http.Request) (*http.Response, error) {
	if req.Method != "GET" {
		return makeResp(fmt.Sprintf(`got method %q want "GET"`, req.Method), http.StatusMethodNotAllowed, nil)
	}
	wantPathSuffix := "/api/v1/trades.do"
	gotPath := req.URL.Path
	if !strings.HasSuffix(gotPath, wantPathSuffix) {
		return makeResp(fmt.Sprintf(`got suffix %q want %q`, gotPath, wantPathSuffix), http.StatusBadRequest, nil)
	}

	query := req.URL.Query()
	symbol := query.Get("symbol")
	if lastTradeIDStr := query.Get("since"); lastTradeIDStr != "" {
		// Ensure that the last requested tradeID is numeric.
		if _, err := strconv.ParseUint(lastTradeIDStr, 10, 64); err != nil {
			return makeResp(err.Error(), http.StatusBadRequest, nil)
		}
	}

	outPath := fmt.Sprintf("./testdata/trades-%s.json", symbol)
	f, err := os.Open(outPath)
	if err != nil {
		return makeResp(err.Error(), http.StatusInternalServerError, nil)
	}
	// This handle should be closed by the consumer.
	return makeResp("200 OK", http.StatusOK, f)
}

// Route declarations
const (
	lastNTradesRoute = "/last-trades"
)
