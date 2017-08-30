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

func Example_client_LastTrades() {
	client, err := okcoin.NewDefaultClient()
	if err != nil {
		log.Fatal(err)
	}

	tres, err := client.LastTrades(&okcoin.LastTradesRequest{
		Symbol:      okcoin.BTCUSD,
		LastTradeID: 354425297,
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Last Trades for Symbol: %s\n\n", tres.Symbol)
	for i, trade := range tres.Trades {
		fmt.Printf("#%d ID: %d\n", i, trade.ID)
		fmt.Printf("\tAmount: %.3f\n", trade.Amount)
		fmt.Printf("\tPrice: %.3f\n", trade.Price)
		fmt.Printf("\tType: %s\n\n", trade.Type)
	}
}

func Example_client_CandleStick() {
	client, err := okcoin.NewDefaultClient()
	if err != nil {
		log.Fatal(err)
	}

	cres, err := client.CandleStick(&okcoin.CandleStickRequest{
		Symbol: okcoin.ETHUSD,
		Period: okcoin.P5Min,
		Since:  1504068023.163283,
		N:      20,
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("CandleSticks for %s Period: %v\n\n", cres.Symbol, cres.Period)
	for i, cstick := range cres.CandleSticks {
		fmt.Printf("#%d %+v\n", i, cstick)
	}
}
