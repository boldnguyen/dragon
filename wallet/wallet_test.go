package wallet

import (
	"testing"
)

func TestConnect(t *testing.T) {
	w := NewVirtualWallet()
	err := w.Connect()
	if err != nil {
		t.Errorf("Connection failed: %v", err)
	}
	if w.Balance != 1000.0 {
		t.Errorf("Expected balance of 1000.0, got %v", w.Balance)
	}
}

func TestDeposit(t *testing.T) {
	w := NewVirtualWallet()
	w.Connect() // Assuming this sets up the wallet
	err := w.Deposit(500.0)
	if err != nil {
		t.Errorf("Deposit failed: %v", err)
	}
	if w.Balance != 1500.0 {
		t.Errorf("Expected balance of 1500.0 after deposit, got %v", w.Balance)
	}
}

func TestWithdraw(t *testing.T) {
	w := NewVirtualWallet()
	w.Connect() // Assuming this sets up the wallet
	err := w.Withdraw(500.0)
	if err != nil {
		t.Errorf("Withdraw failed: %v", err)
	}
	if w.Balance != 500.0 {
		t.Errorf("Expected balance of 500.0 after withdrawal, got %v", w.Balance)
	}

	// Test for insufficient funds
	err = w.Withdraw(600.0)
	if err == nil {
		t.Error("Withdrawal should have failed due to insufficient funds")
	}
}

func TestGetAssetCount(t *testing.T) {
	w := NewVirtualWallet()
	w.Connect() // Assuming this sets up the wallet
	count, err := w.GetAssetCount("dragons")
	if err != nil {
		t.Errorf("Failed to get asset count: %v", err)
	}
	if count != 5 {
		t.Errorf("Expected 5 dragons, got %v", count)
	}

	_, err = w.GetAssetCount("non_existent_asset")
	if err == nil {
		t.Error("Expected error for non-existent asset")
	}
}
