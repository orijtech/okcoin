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
	"log"

	"github.com/orijtech/okcoin/v1"
)

func Example_client_Ticker() {
	client, err := okcoin.NewDefaultClient()
	if err != nil {
		log.Fatal(err)
	}

	ticker, err := client.Ticker(okcoin.BTCUSD)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Last ticker")
	fmt.Printf("TimeAtEpoch: %f\n", ticker.TimeAtEpoch)
	fmt.Printf("BuyPrice: %.4f\n", ticker.Ticker.Buy)
	fmt.Printf("SellPrice: %.4f\n", ticker.Ticker.Sell)
	fmt.Printf("\nLow: %.4f\n", ticker.Ticker.Low)
	fmt.Printf("High: %.4f\n\n", ticker.Ticker.High)
	fmt.Printf("Volume: %.4f\n", ticker.Ticker.Volume)
}
