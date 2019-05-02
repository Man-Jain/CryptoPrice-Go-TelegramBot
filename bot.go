package main

import (
	"time"
	tb "gopkg.in/tucnak/telebot.v2"
	"fmt"
	"net/http"
	//"strconv"
	"io/ioutil"
	"encoding/json"
	"os"
)

func main(){

	type PriceMap struct {
		Status string `json:"status"`
		Data   struct {
			Stats struct {
				Total  int    `json:"total"`
				Offset int    `json:"offset"`
				Base   string `json:"base"`
			} `json:"stats"`
			Base struct {
				Symbol string `json:"symbol"`
				Sign   string `json:"sign"`
			} `json:"base"`
			Coins []struct {
				ID          int    `json:"id"`
				Slug        string `json:"slug"`
				Symbol      string `json:"symbol"`
				Name        string `json:"name"`
				Description string `json:"description"`
				Color       string `json:"color"`
				IconType    string `json:"iconType"`
				IconURL     string `json:"iconUrl"`
				WebsiteURL  string `json:"websiteUrl"`
				Socials     []struct {
					Name string `json:"name"`
					URL  string `json:"url"`
					Type string `json:"type"`
				} `json:"socials"`
				ConfirmedSupply   bool     `json:"confirmedSupply"`
				NumberOfMarkets   int      `json:"numberOfMarkets"`
				NumberOfExchanges int      `json:"numberOfExchanges"`
				Type              string   `json:"type"`
				Volume            int64    `json:"volume"`
				MarketCap         int64    `json:"marketCap"`
				Price             string   `json:"price"`
				CirculatingSupply float64  `json:"circulatingSupply"`
				TotalSupply       float64  `json:"totalSupply"`
				ApprovedSupply    bool     `json:"approvedSupply"`
				FirstSeen         int64    `json:"firstSeen"`
				Change            float64  `json:"change"`
				Rank              int      `json:"rank"`
				History           []string `json:"history"`
				AllTimeHigh       struct {
					Price     string `json:"price"`
					Timestamp int64  `json:"timestamp"`
				} `json:"allTimeHigh"`
				Penalty bool `json:"penalty"`
			} `json:"coins"`
		} `json:"data"`
	}

	b, err := tb.NewBot(
		tb.Settings{
			Token: os.Getenv("key"),
			Poller: &tb.LongPoller{Timeout: 10 * time.Second},
		})

		if err != nil {
			fmt.Println("Error couldn't start Bot")
			return
		}

		fmt.Println("Crypto Bot started")

		b.Handle("/hello", func(m *tb.Message){
			fmt.Println(m.Text)
			b.Send(m.Sender, "Hello, I am a Crypto Bot made in Golang. Send me crypto symbols and I'll give you the prices in INR.\nSend /help for more help.")
		})

		b.Handle("/help", func(m *tb.Message){
			fmt.Println(m.Text)
			b.Send(m.Sender, "Send me Crypto symbol for Eg:- ETH for Ethereum, BTC for Bitcoin etc.")
		})

		b.Handle("/start", func(m *tb.Message){
			fmt.Println(m.Text)
			b.Send(m.Sender, "/hello - For bot introduction\n/help - To get help on bot")
		})

		b.Handle(tb.OnText, func(m *tb.Message){
			req, err := http.NewRequest("GET", "https://api.coinranking.com/v1/public/coins",nil)

				if err != nil {
					return
				}

				q := req.URL.Query()
				q.Add("base","INR")
				q.Add("symbols",m.Text)
				req.URL.RawQuery = q.Encode()
				client := &http.Client{}

				resp, err := client.Do(req)

				if err != nil {
					return
				}

				data, _ := ioutil.ReadAll(resp.Body)
				var parsed_data PriceMap
				json.Unmarshal(data, &parsed_data)
				fmt.Printf("Results: %v\n", parsed_data.Data.Coins[0].Price)
				//?base=INR&timePeriod=7d&symbols=ETH
				//b.Send(m.Sender, "Price is: " + strconv.Itoa(price_map[m.Text]))
				b.Send(m.Sender, "The price for " + m.Text + " is:- Rs. " + string(parsed_data.Data.Coins[0].Price))
			})

			b.Start()
		}
