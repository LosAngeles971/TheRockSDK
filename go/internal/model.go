package internal

import "fmt"

const (
	API_BASEURL = "https://api.therocktrading.com/v1/"
	RATE_LIMIT  = 9
	PAIR_DIRECT_SUPPORT = 0
	PAIR_REVERSE_SUPPORT = 1
)

type Credential struct {
	RoKey    string
	RoSecret string
	RwKey    string
	RwSecret string
}

type Fund struct {
	ID                    string  `json"id"`
	Description           string  `json"description"`
	Type                  string  `json"type"`
	Base_Currency          string  `json"base_currency"`
	Trade_Currency         string  `json"trade_currency"`
	Buy_Fee                float64 `json"buy_fee"`
	Sell_Fee               float64 `json"sell_fee"`
	Minimum_Price_Offer     float64 `json"minimum_price_offer"`
	Minimum_Quantity_Offer  float64 `json"minimum_quantity_offer"`
	Base_Currency_Decimals  float64 `json"base_currency_decimals"`
	Trade_Currency_Decimals float64 `json"trade_currency_decimals"`
}

type Funds struct {
	Funds []Fund `json"funds"`
}

func (funds *Funds) GetFund(id string) (Fund, bool) {
	for _, fund := range funds.Funds {
		if fund.ID == id {
			return fund, true
		}
	}
	return Fund{}, false
}

type Balance struct {
	Currency       string  `json"currency"`
	Balance        float64 `json"balance"`
	Trading_Balance float64 `json"trading_balance"`
}

type Wallet struct {
	Balances []Balance `json"balances"`
}

func (wallet *Wallet) GetBalance(id string) (Balance, bool) {
	for _, balance := range wallet.Balances {
		if balance.Currency == id {
			return balance, true
		}
	}
	return Balance{}, false
}

type Ticker struct {
	Fund_ID       string  `json"fund_id"`
	Date         string  `json"date"`
	Bid          float64 `json"bid"`
	Ask          float64 `json"ask"`
	Last         float64 `json"last"`
	Volume       float64 `json"volume"`
	Volume_Traded float64 `json"volume_traded"`
	Open         float64 `json"open"`
	High         float64 `json"high"`
	Low          float64 `json"low"`
	Close        float64 `json"close"`
}

type Tickers struct {
	Tickers []Ticker `json"tickers"`
}

func (tickers *Tickers) GetTicker(id string) (Ticker, bool) {
	for _, ticker := range tickers.Tickers {
		fmt.Println(ticker.Volume_Traded)
		if ticker.Fund_ID == id {
			return ticker, true
		}
	}
	return Ticker{}, false
}

type Currency struct {
	Symbol   string `json"symbol"`
	Common_Name     string `json"common_name"`
	Decimals int    `json"decimals"`
}

type Currencies struct {
	Currencies []Currency `json"currencies"`
}

func (currencies *Currencies) GetCurrency(id string) (Currency, bool) {
	for _, currency := range currencies.Currencies {
		if currency.Symbol == id {
			return currency, true
		}
	}
	return Currency{}, false
}

type Trade struct {
	ID   string `json"id"`
	Fund_ID   string `json"fund_id"`
	Amount   float64 `json"amount"`
	Price   float64 `json"price"`
	Side   string `json"side"`
	Dark bool `json"dark"`
	Date string `json"date"`
}

type Order struct {
	ID   string `json"id"`
	Fund_ID   string `json"fund_id"`
	Side   string `json"side"`
	Type   string `json"type"`
	Status   string `json"status"`
	Price   float64 `json"price"`
	Amount   float64 `json"amount"`
	Amount_Unfilled   float64 `json"amount_unfilled"`
	Conditional_Type   string `json"conditional_type"`
	Conditional_Price float64 `json"conditional_price"`
	Date string `json"date"`
	Close_On string `json"close_on"`
	Leverage float64 `json"leverage"`
	Position_ID int `json"position_id"`
}

type PlacingOrder struct {
	ID   string `json"id"`
	Fund_ID   string `json"fund_id"`
	Side   string `json"side"`
	Price   float64 `json"price"`
	Amount   float64 `json"amount"`
	Conditional_Type   string `json"conditional_type"`
	Conditional_Price float64 `json"conditional_price"`
	Leverage float64 `json"leverage"`
	Position_ID int `json"position_id"`
	Position_Order_ID int `json"position_order_id"`
}