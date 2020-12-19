package pkg

import (
	"os"
	"testing"

	"it/losangeles971/therocksdk/internal"
)

func getCredential() internal.Credential {
	cred := internal.Credential{
		RoKey: os.Getenv("RO_KEY"),
		RoSecret: os.Getenv("RO_SECRET"),
	}
	return cred
}

func TestFunds(t *testing.T) {
	funds, err := GetFunds(getCredential())
	if err != nil {
		t.Errorf("Error %v", err)
		t.Fail()
	}
	if len(funds.Funds) < 1 {
		t.Errorf("No funds")
		t.Fail()
	}
	fund, ok := funds.GetFund("BTCEUR")
	if !ok {
		t.Errorf("No BTCEUR")
		t.Fail()
	}
	if fund.Buy_Fee == 0 || fund.Base_Currency_Decimals == 0 || fund.Base_Currency != "EUR" {
		t.Errorf("BTCEUR is invalid")
		t.Fail()
	}
	s1, err2 := IsSupported("BTC", "EUR", getCredential())
	if err2 != nil || s1 != internal.PAIR_DIRECT_SUPPORT {
		t.Errorf("BTCEUR should be directly supported")
		t.Fail()
	}
	s2, err3 := IsSupported("EUR", "BTC", getCredential())
	if err3 != nil || s2 != internal.PAIR_REVERSE_SUPPORT {
		t.Errorf("BTCEUR should be reversely supported")
		t.Fail()
	}
}

func TestCurrencies(t *testing.T) {
	currencies, err := GetCurrencies(getCredential())
	if err != nil {
		t.Errorf("Error %v", err)
		t.Fail()
	}
	if len(currencies.Currencies) < 1 {
		t.Errorf("No funds")
		t.Fail()
	}
	btc, ok := currencies.GetCurrency("BTC")
	if !ok {
		t.Errorf("No BTC")
		t.Fail()
	}
	if btc.Common_Name != "Bitcoin" || btc.Decimals == 0 {
		t.Errorf("BTC currency is corrupted")
		t.Fail()
	}
}

func TestWallet(t *testing.T) {
	_, err := GetWallet(getCredential())
	if err != nil {
		t.Errorf("Error %v", err)
		t.Fail()
	}
}

func TestTickers(t *testing.T) {
	tickers, err := GetTickers(getCredential())
	if err != nil {
		t.Errorf("Error %v", err)
		t.Fail()
	}
	ticker, ok := tickers.GetTicker("BTCEUR")
	if !ok {
		t.Errorf("No BTCEUR ticker")
		t.Fail()
	}
	if ticker.Volume_Traded == 0 {
		t.Errorf("BTCEUR ticker is invalid")
		t.Fail()
	}
}