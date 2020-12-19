package internal

import (
	"fmt"
	"os"
	"testing"
)

func getCredential() Credential {
	cred := Credential{
		RoKey: os.Getenv("RO_KEY"),
		RoSecret: os.Getenv("RO_SECRET"),
	}
	return cred
}

func TestSign(t *testing.T) {
	s := sign("test message for sign function", "test secret for sign function")
	fmt.Println(s)
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
	_, ok := funds.GetFund("BTCEUR")
	if !ok {
		t.Errorf("No BTCEUR")
		t.Fail()
	}
}