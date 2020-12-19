package pkg

import (
	"encoding/json"
	"errors"

	"it/losangeles971/therocksdk/internal"
)

var funds = internal.Funds{}
var currencies = internal.Currencies{}

// Get the list of supported funds (aka trading pairs), this method caches the result
func GetFunds(cred internal.Credential) (internal.Funds, error) {
	if len(funds.Funds) <1 {
		body, err := internal.Get(internal.API_BASEURL + "funds", cred.RoKey, cred.RoSecret)
		if err != nil {
			return funds, err
		}
		err = json.Unmarshal(body, &funds)
		if err != nil {
			return funds, err
		}
	}
	return funds, nil
}

// Get the list of supported currencies, this method caches the result
func GetCurrencies(cred internal.Credential) (internal.Currencies, error) {
	if len(currencies.Currencies) <1 {
		body, err := internal.Get(internal.API_BASEURL + "currencies", cred.RoKey, cred.RoSecret)
		if err != nil {
			return currencies, err
		}
		err = json.Unmarshal(body, &currencies)
		if err != nil {
			return currencies, err
		}
	}
	return currencies, nil
}

// It returns the ticker for the given fund (aka trading pair)
func GetTicker(fund string, cred internal.Credential) (internal.Ticker, error) {
	ticker := internal.Ticker{}
	body, err := internal.Get(internal.API_BASEURL + "funds/" + fund + "/ticker", cred.RoKey, cred.RoSecret)
	if err != nil {
		return ticker, err
	}
	err = json.Unmarshal(body, &ticker)
	if err != nil {
		return ticker, err
	}
	return ticker, nil
}

// It returns all tickers for the support trading pairs
func GetTickers(cred internal.Credential) (internal.Tickers, error) {
	tickers := internal.Tickers{}
	body, err := internal.Get(internal.API_BASEURL + "funds/tickers", cred.RoKey, cred.RoSecret)
	if err != nil {
		return tickers, err
	}
	err = json.Unmarshal(body, &tickers)
	if err != nil {
		return tickers, err
	}
	return tickers, nil
}

// It returns a balance object for each "owned" currency
func GetWallet(cred internal.Credential) (internal.Wallet, error) {
	wallet := internal.Wallet{}
	body, err := internal.Get(internal.API_BASEURL + "balances", cred.RoKey, cred.RoSecret)
	if err != nil {
		return wallet, err
	}
	err = json.Unmarshal(body, &wallet)
	if err != nil {
		return wallet, err
	}
	return wallet, nil
}

// It returns a running order
func GetOrder(fund string, id string, cred internal.Credential) (internal.Order, error) {
	order := internal.Order{}
	body, err := internal.Get(internal.API_BASEURL + "funds/" + fund + "/orders/" + id, cred.RoKey, cred.RoSecret)
	if err != nil {
		return order, err
	}
	err = json.Unmarshal(body, &order)
	if err != nil {
		return order, err
	}
	return order, nil
}

// It deletes a running order
func DeleteOrder(fund string, id string, cred internal.Credential) (internal.Order, error) {
	order := internal.Order{}
	body, err := internal.Delete(internal.API_BASEURL + "funds/" + fund + "/orders/" + id, cred.RoKey, cred.RoSecret)
	if err != nil {
		return order, err
	}
	err = json.Unmarshal(body, &order)
	if err != nil {
		return order, err
	}
	return order, nil
}

// It submits an order
func CreateOrder(placing internal.PlacingOrder, cred internal.Credential) (internal.Order, error) {
	order := internal.Order{}
	data, err := json.Marshal(placing)
	if err != nil {
		return order, err
	}
	body, err2 := internal.Post(internal.API_BASEURL + "funds/" + placing.Fund_ID + "/orders/", data, cred.RoKey, cred.RoSecret)
	if err2 != nil {
		return order, err
	}
	err = json.Unmarshal(body, &order)
	if err != nil {
		return order, err
	}
	return order, nil
}

// Check if a trading pair is supported by the Exchange
func IsSupported(c1Symbol string, c2Symbol string, cred internal.Credential) (int, error) {
	direct := c1Symbol + c2Symbol
	reverse := c2Symbol + c1Symbol
	funds, err := GetFunds(cred)
	if err != nil {
		return 0, err
	}
	for _, fund := range funds.Funds {
		if fund.ID == direct {
			return internal.PAIR_DIRECT_SUPPORT, nil
		}
		if fund.ID == reverse {
			return internal.PAIR_REVERSE_SUPPORT, nil
		}
	}
	return 0, errors.New("Unsupported pair " + direct + "/" + reverse)
}