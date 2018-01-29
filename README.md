# Two Pore Guys UI Engineer test - Stock browser

## Goal
This exercise is meant to demonstrate candidate's ability to build a sample application, using a modern JS framework (ReactJS is preferred).
The design is not the main concern of the exercise, still the application should be good looking enough,  so that user's eyes doesn't start bleeding while using the application.

## Objective
The application will make use of a websocket endpoint to list NASDAQ stocks, get data about a stock and subscribe to updates of those data.
The websocket is `wss://stock-browser.herokuapp.com/`.
It must be possible to:
- list stocks
- show quote for a stock
- subscribe / unsubscribe to quote changes and to reflect those changes
- show the evolution of stock's quote per day since the beginning of the year
- show the 10 last news related to a stock
- subscribe / unsubscribe to news about a stock and display the updated news
- navigate through stocks using peers relationship between them

**
*Data is provided by [IEX Developer Platform](https://iextrading.com/developer/) and based on NASDAQ realtime data.*

*NASDAQ being open from 4am to 8pm EST, there may not be updates on stocks outside of those hours.*
**

It is also recommended that you write at least one test as an exemplar of your approach to testing.

## Backend API

### Messages format
Messages can have three types:
- Request: From client to server
- Response: From server to client, as a request answer
- Update: From server to client, notify that data have changed on server side

Request messages have the form
```
{
  "id": <ANY_UNIQUE_ID>,
  "namespace": <"get"|"subscribe"|"unsubscribe"|"control">,
  "name": <NAME_OF_THE_COMMAND>,
  "args": { // OPTIONAL PROPERTY
    <KEY>: <VALUE>
  }
}
```

Response messages have the form
```
{
  "id": <ID_MATCHING_THE_REQUEST_ID>,
  "namespace": <"get"|"subscribe"|"unsubscribe"|"control">,
  "name": <NAME_OF_THE_COMMAND>,
  "data": <DATA_RETURNED_BY_THE_CALL>
}
```

Update messages have the form
```
{
  "name": <"stock.quote"|"stock.news">,
  "data": <UPDATED_DATA>
}
```

### Available namespaces

#### `get`: Query data from the server.
##### Commands

- `stock.list`: Returns the list of known stocks.
  ###### Return type: `Stock[]`
  ###### No argument required.

- `stock.random`: Returns one random stock.
  ###### Return type: `Stock`
  ###### No argument required.

- `stock.quote`: Returns current quote of a given stock.
  ###### Return type: `Quote`
  ###### Arguments:
  - `stock`: The symbol of the stock.

- `stock.chart`: Returns the evolution of quote per day since the beginning of the year for a given stock.
  ###### Return type: `Chart[]`
  ###### Arguments:
  - `stock`: The symbol of the stock.

- `stock.peers`: Returns a list of stocks related to a given stock.
  ###### Return type: `string[]`
  ###### Arguments:
  - `stock`: The symbol of the stock.

- `stock.current`: Returns the current value of a given stock.
  ###### Return type: `float`
  ###### Arguments:
  - `stock`: The symbol of the stock.

- `stock.news`: Returns the 10 last news related to a given stock.
  ###### Return type: `News[]`
  ###### Arguments:
  - `stock`: The symbol of the stock.

#### `subscribe`: Subscribe to updates
**Note:** *There can be only one subscription active at a time for a given subscription type, subscribing to a new stock updates will automatically cancel any previous subscription of this type.*
##### Commands

- `stock.quote`: Subscribes to updates on a quote for a given stock.
  ###### Arguments:
  - `stock`: The symbol of the stock.

- `stock.news`: Subscribes to updates on news for a given stock.
  ###### Arguments:
  - `stock`: The symbol of the stock.

#### `unsubscribe`: Cancel a subscription
**Note:** *There can be only one subscription active at a time for a given subscription type.*
##### Commands

- `stock.quote`: Cancel subscription to updates on a quote for a given stock.
  ###### No argument required.

- `stock.news`: Cancel subscription to updates on news for a given stock.
  ###### No argument required.

#### `control`: Control messages
##### Commands

- `ping`: Send a keep-alive signal to the other end.
  ###### No argument required.

- `pong`: Acknowledge to a ping message.
  ###### No argument required.

### Data model
#### `Stock`
```json
{
  "symbol": "AAPL",
  "title":  "Apple Inc."
}
```
#### `Quote`
```json
{
  "symbol":        "AAPL",
  "companyName":   "Apple Inc.",
  "sector":        "Technology",
  "latestPrice":   179.24,
  "open":          179.47,
  "close":         179.1,
  "latestUpdate":  1516309199625
}
```
#### `Chart`
```json
{
  "date": "2018-01-02",
  "label": "Jan 2",
  "open": 170.16,
  "close": 172.26
}
```
#### `News`
```json
{
  "datetime": "2018-01-18T15:26:00-05:00",
  "headline": "Tim Cook tells Cramer: New tax law helped pave the way for Apple's massive investment plan",
  "source": "CNBC",
  "url": "https://api.iextrading.com/1.0/stock/aapl/article/9002290080019164",
  "related": "AAPL"
}
```
