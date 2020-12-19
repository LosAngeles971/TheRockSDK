package internal

const (
	API_BASEURL = "https://api.therocktrading.com/v1/"
	RATE_LIMIT = 9
)

type Credential struct {
	RoKey string
	RoSecret string
	RwKey string
	RwSecret string
}

type Fund struct 
{
	ID string `json"id"`
	Description string `json"description"`
	Type string `json"type"`
	BaseCurrency string `json"base_currency"`
	TradeCurrency string `json"trade_currency"`
	BuyFee float64 `json"buy_fee"`
	SellFee float64 `json"sell_fee"`
	MinimumPriceOffer float64 `json"minimum_price_offer"`
	MinimumQuantityOffer float64 `json"minimum_quantity_offer"`
	BaseCurrencyDecimals float64 `json"base_currency_decimals"`
	TradeCurrencyDecimals float64 `json"trade_currency_decimals"`
}

type Funds struct {
	Funds []Fund `json"funds"`
}

func (funds *Funds) GetFund(id string) (Fund, bool) {
	for _, fund := range funds.Funds {
		if fund.ID == "BTCEUR" {
			return fund, true
		}
	}
	return Fund{}, false
}

type Balance struct 
{
	Currency string `json"currency"`
	Balance float64 `json"balance"`
	TradingBalance float64 `json"trading_balance"`
}

type Wallet struct
{
	Balances []Balance `json"balances"`
}

type Ticker struct 
{
	FundID string `json"fund_id"`
	Date string `json"date"`
	Bid float64 `json"bid"`
	Ask float64 `json"ask"`
	Last float64 `json"last"`
	Volume float64 `json"volume"`
	VolumeTraded float64 `json"volume_traded"`
	Open float64 `json"open"`
	High float64 `json"high"`
	Low float64 `json"low"`
	Close float64 `json"close"`
}

type Tickers struct 
{
	Tickers []Ticker `json"tickers"`
}