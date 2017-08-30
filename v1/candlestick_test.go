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
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/orijtech/okcoin/v1"
)

var (
	blankCandleStickResponse = new(okcoin.CandleStickResponse)
	blankCandleStick         = new(okcoin.CandleStick)
)

func TestCandleStick(t *testing.T) {
	t.Parallel()

	client, err := okcoin.NewDefaultClient()
	if err != nil {
		t.Fatalf("new client: %v", err)
	}

	tests := [...]struct {
		req         *okcoin.CandleStickRequest
		wantAtLeast int
		wantErr     bool
	}{
		0: {req: nil, wantAtLeast: 2},
		1: {req: &okcoin.CandleStickRequest{Since: -1}, wantAtLeast: 2},
		2: {req: &okcoin.CandleStickRequest{Since: 1504068205.797958}, wantAtLeast: 2},
	}

	client.SetHTTPRoundTripper(&backend{route: candleStickRoute})
	for i, tt := range tests {
		cres, err := client.CandleStick(tt.req)
		if tt.wantErr {
			if err == nil {
				t.Errorf("#%d: want non-nil; got resp=%#v", i, cres)
			}
			continue
		}

		if err != nil {
			t.Errorf("#%d: got err: %v", i, err)
			continue
		}
		if reflect.DeepEqual(cres, blankCandleStickResponse) {
			t.Errorf("#%d: got blankCandleStickResponse", i)
			continue
		}
		validNonBlankCandleSticks := 0
		for _, cstick := range cres.CandleSticks {
			if !reflect.DeepEqual(cstick, blankCandleStick) {
				validNonBlankCandleSticks += 1
			}
		}
		if g, w := validNonBlankCandleSticks, tt.wantAtLeast; g < w {
			t.Errorf("#%d: got=%d wantAtLeast=%d", i, g, w)
		}
	}
}

func (b *backend) candleStickRoundTrip(req *http.Request) (*http.Response, error) {
	if req.Method != "GET" {
		return makeResp(fmt.Sprintf(`method: got=%q want="GET"`, req.Method), http.StatusMethodNotAllowed, nil)
	}
	wantPathSuffix := "/api/v1/kline.do"
	gotPath := req.URL.Path
	if !strings.HasSuffix(gotPath, wantPathSuffix) {
		return makeResp(fmt.Sprintf(`got suffix %q want %q`, gotPath, wantPathSuffix), http.StatusBadRequest, nil)
	}
	query := req.URL.Query()
	if sinceStr := query.Get("since"); sinceStr != "" {
		since, err := strconv.ParseFloat(sinceStr, 64)
		if err != nil {
			return makeResp(fmt.Sprintf(`"since": parse err %v`, err), http.StatusBadRequest, nil)
		}
		if since < 0.0 {
			return makeResp(fmt.Sprintf(`"since": got=%f wantAtLeast 0`, since), http.StatusBadRequest, nil)
		}
	}
	symbol := query.Get("symbol")
	outPath := fmt.Sprintf("./testdata/candlestick-%s.json", symbol)
	f, err := os.Open(outPath)
	if err != nil {
		return makeResp(err.Error(), http.StatusInternalServerError, nil)
	}
	// This handle should be closed by the consumer.
	return makeResp("200 OK", http.StatusOK, f)
}

const (
	candleStickRoute = "/candle-stick"
)
