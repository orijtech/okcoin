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

package okcoin_test

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"testing"

	"github.com/orijtech/okcoin/v1"
)

const (
	apiKey1    = "foo"
	apiSecret1 = "$B4r^"
	fundsRoute = "/funds"
)

func TestFunds(t *testing.T) {
	t.Parallel()

	client, err := okcoin.NewDefaultClient()
	if err != nil {
		t.Fatalf("new client: %v", err)
	}
	client.SetHTTPRoundTripper(&backend{route: fundsRoute})

	tests := []struct {
		creds   *okcoin.Credentials
		want    *okcoin.Funds
		wantErr string
	}{
		{&okcoin.Credentials{}, nil, `expecting "api_key"`},
		{&okcoin.Credentials{APIKey: apiKey1, Secret: apiSecret1}, fundsFromFile("funds1.json"), ""},
	}

	for i, tt := range tests {
		client.SetCredentials(tt.creds)
		funds, err := client.Funds()
		if tt.wantErr != "" {
			if err == nil {
				t.Errorf("#%d: want non-nil error", i)
			} else if got, want := err.Error(), tt.wantErr; !strings.Contains(got, want) {
				t.Errorf("#%d: got=%q want=%q", i, got, want)
			}
			continue
		}
		if err != nil {
			t.Errorf("#%d: got unexpected err: %v", i, err)
			continue
		}
		if g, w := funds, tt.want; !reflect.DeepEqual(g, w) {
			t.Errorf("#%d:\ngot=%#v\nwant=%#v\n", i, g, w)
		}
	}
}

var knownAPIKeyToSecrets = map[string]string{
	apiKey1: apiSecret1,
}

func (b *backend) fundsRoundTrip(req *http.Request) (*http.Response, error) {
	if req.Method != "POST" {
		return makeResp(fmt.Sprintf(`got method %q want "POST"`, req.Method), http.StatusMethodNotAllowed, nil)
	}
	qv := req.URL.Query()
	apiKey := qv.Get("api_key")
	if apiKey == "" {
		return makeResp(`expecting "api_key" in query string`, http.StatusBadRequest, nil)
	}
	gotSignature := qv.Get("sign")
	if gotSignature != "" && gotSignature != strings.ToUpper(gotSignature) {
		return makeResp(`"sign" must be an all uppercase string`, http.StatusBadRequest, nil)
	}
	knownSecret, ok := knownAPIKeyToSecrets[apiKey]
	if !ok || knownSecret == "" {
		return makeResp("unknown API secret", http.StatusUnauthorized, nil)
	}
	qv.Del("sign")
	h := md5.New()
	fmt.Fprintf(h, "%s&secret_key=%s", qv.Encode(), knownSecret)
	wantSignature := strings.ToUpper(fmt.Sprintf("%x", h.Sum(nil)))
	if wantSignature != gotSignature {
		return makeResp("signatures do not match", http.StatusBadRequest, nil)
	}
	return respFromFile("./testdata/funds1.json")
}

type fundsIntermediate struct {
	Funds map[string]*okcoin.Funds `json:"info"`

	Result bool `json:"result"`
}

var blankFunds = new(okcoin.Funds)

func fundsFromFile(basename string) (*okcoin.Funds) {
	blob, err := ioutil.ReadFile("./testdata/" + basename)
	if err != nil {
		return nil
	}
	fi := new(fundsIntermediate)
	if err := json.Unmarshal(blob, fi); err != nil {
		return nil
	}
	funds := fi.Funds["funds"]
	if funds == nil || reflect.DeepEqual(funds, blankFunds) {
		return nil
	}
	return funds
}
