package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type wsHandler struct {
	conn       *websocket.Conn
	quote      Quote
	quoteBatch chan struct{}
	news       []News
	newsBatch  chan struct{}
}

func (h wsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	h.conn = conn
	defer h.conn.Close()
	for {
		message := Request{}
		err := h.conn.ReadJSON(&message)
		if err != nil {
			break
		}
		h.dispatch(message)
	}
}

type Request struct {
	Id        string                 `json:"id"`
	Namespace string                 `json:"namespace"`
	Name      string                 `json:"name"`
	Args      map[string]interface{} `json:"args,omitEmpty"`
}

type Response struct {
	Id   string      `json:"id"`
	Name string      `json:"name"`
	Data interface{} `json:"data,omitEmpty"`
}

type Update struct {
	Name string      `json:"name"`
	Data interface{} `json:"data"`
}

type Stock struct {
	Symbol string `json:"symbol"`
	Name   string `json:"name"`
}

type Quote struct {
	Symbol       string  `json:"symbol"`
	CompanyName  string  `json:"companyName"`
	Sector       string  `json:"sector"`
	LatestPrice  float32 `json:"latestPrice"`
	Open         float32 `json:"open"`
	Close        float32 `json:"close"`
	LatestUpdate int64   `json:"latestUpdate"`
}

type Chart struct {
	Date  string  `json:"date"`
	Label string  `json:"label"`
	Open  float32 `json:"open"`
	Close float32 `json:"close"`
}

type News struct {
	DateTime string `json:"datetime"`
	Headline string `json:"headline"`
	Source   string `json:"source"`
	Url      string `json:"url"`
	Related  string `json:"related"`
}

const (
	Get         = "get"
	Subscribe   = "subscribe"
	Unsubscribe = "unsubscribe"
)

const (
	StockList    = "stock.list"
	StockRandom  = "stock.random"
	StockQuote   = "stock.quote"
	StockChart   = "stock.chart"
	StockPeers   = "stock.peers"
	StockCurrent = "stock.current"
	StockNews    = "stock.news"
)

func (h wsHandler) getStocks() []Stock {
	resp, _ := http.Get("https://api.iextrading.com/1.0/ref-data/symbols")
	defer resp.Body.Close()
	stocks := []Stock{}
	decoder := json.NewDecoder(resp.Body)
	decoder.Decode(&stocks)
	return stocks
}

func (h wsHandler) getQuote(stock string) Quote {
	resp, _ := http.Get(fmt.Sprintf("https://api.iextrading.com/1.0/stock/%s/quote", stock))
	defer resp.Body.Close()
	quote := Quote{}
	decoder := json.NewDecoder(resp.Body)
	decoder.Decode(&quote)
	return quote
}

func (h wsHandler) getNews(stock string) []News {
	resp, _ := http.Get(fmt.Sprintf("https://api.iextrading.com/1.0/stock/%s/news", stock))
	defer resp.Body.Close()
	news := []News{}
	decoder := json.NewDecoder(resp.Body)
	decoder.Decode(&news)
	return news
}

func (h wsHandler) handleStockList(response *Response) {
	response.Data = h.getStocks()
}

func (h wsHandler) handleStockRandom(response *Response) {
	stocks := h.getStocks()
	response.Data = stocks[rand.Intn(len(stocks))]
}

func (h wsHandler) handleStockQuote(response *Response, stock string) {
	response.Data = h.getQuote(stock)
}

func (h wsHandler) handleStockChart(response *Response, stock string) {
	resp, _ := http.Get(fmt.Sprintf("https://api.iextrading.com/1.0/stock/%s/chart/ytd", stock))
	defer resp.Body.Close()
	charts := []Chart{}
	decoder := json.NewDecoder(resp.Body)
	decoder.Decode(&charts)
	response.Data = charts
}

func (h wsHandler) handleStockPeers(response *Response, stock string) {
	resp, _ := http.Get(fmt.Sprintf("https://api.iextrading.com/1.0/stock/%s/peers", stock))
	defer resp.Body.Close()
	peers := []string{}
	decoder := json.NewDecoder(resp.Body)
	decoder.Decode(&peers)
	response.Data = peers
}

func (h wsHandler) handleStockCurrent(response *Response, stock string) {
	resp, _ := http.Get(fmt.Sprintf("https://api.iextrading.com/1.0/stock/%s/price", stock))
	defer resp.Body.Close()
	current := 0.0
	decoder := json.NewDecoder(resp.Body)
	decoder.Decode(&current)
	response.Data = current
}

func (h wsHandler) handleStockNews(response *Response, stock string) {
	response.Data = h.getNews(stock)
}

func (h wsHandler) handleSubscribeQuote(stock string) chan struct{} {
	if h.quoteBatch != nil {
		close(h.quoteBatch)
		h.quoteBatch = nil
	}
	h.quote = h.getQuote(stock)
	ticker := time.NewTicker(10 * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				updatedQuote := h.getQuote(stock)
				if h.quote.LatestUpdate != updatedQuote.LatestUpdate {
					h.quote = updatedQuote
					h.conn.WriteJSON(Update{StockQuote, h.quote})
				}
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
	return quit
}

func (h wsHandler) handleSubscribeNews(stock string) chan struct{} {
	if h.newsBatch != nil {
		close(h.newsBatch)
		h.newsBatch = nil
	}
	h.news = h.getNews(stock)
	ticker := time.NewTicker(10 * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				updatedNews := h.getNews(stock)
				if h.news[0].DateTime != updatedNews[0].DateTime {
					h.news = updatedNews
					h.conn.WriteJSON(Update{StockNews, h.news})
				}
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
	return quit
}

func (h wsHandler) dispatch(message Request) {
	response := Response{message.Id, message.Name, nil}
	switch message.Namespace {
	case Get:
		switch message.Name {
		case StockList:
			h.handleStockList(&response)
		case StockRandom:
			h.handleStockRandom(&response)
		case StockQuote:
			h.handleStockQuote(&response, message.Args["stock"].(string))
		case StockChart:
			h.handleStockChart(&response, message.Args["stock"].(string))
		case StockPeers:
			h.handleStockPeers(&response, message.Args["stock"].(string))
		case StockCurrent:
			h.handleStockCurrent(&response, message.Args["stock"].(string))
		case StockNews:
			h.handleStockNews(&response, message.Args["stock"].(string))
		}
	case Subscribe:
		switch message.Name {
		case StockQuote:
			h.quoteBatch = h.handleSubscribeQuote(message.Args["stock"].(string))
		case StockNews:
			h.newsBatch = h.handleSubscribeNews(message.Args["stock"].(string))
		}
	case Unsubscribe:
		switch message.Name {
		case StockQuote:
			if h.quoteBatch != nil {
				close(h.quoteBatch)
				h.quoteBatch = nil
			}
		case StockNews:
			if h.newsBatch != nil {
				close(h.newsBatch)
				h.newsBatch = nil
			}
		}
	}
	h.conn.WriteJSON(response)
}

func getPort() string {
	var port = os.Getenv("PORT")
	// Set a default port if there is nothing in the environment
	if port == "" {
		port = "4747"
		fmt.Println("INFO: No PORT environment variable detected, defaulting to " + port)
	}
	return ":" + port
}

func main() {
	rand.Seed(time.Now().UnixNano())
	http.Handle("/", wsHandler{})
	log.Fatal(http.ListenAndServe(getPort(), nil))
}
