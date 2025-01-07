package wallet

import (
	"errors"
	"fmt"
)

type VirtualWallet struct {
	PublicKey string
	Balance   float64
	Assets    map[string]int
}

func NewVirtualWallet() *VirtualWallet {
	return &VirtualWallet{
		PublicKey: "0xVirtualPublicKey",
		Balance:   1000.0,
		Assets:    make(map[string]int),
	}
}

func (w *VirtualWallet) Connect() error {
	w.Assets["token"] = 1000
	w.Assets["dragons"] = 5
	w.Assets["items"] = 10
	fmt.Println("Virtual wallet connected with default assets")
	return nil
}

func (w *VirtualWallet) SyncAssets() {
	fmt.Println("Assets synced")
}

func (w *VirtualWallet) Deposit(amount float64) error {
	if amount <= 0 {
		return errors.New("deposit amount must be positive")
	}
	w.Balance += amount
	fmt.Printf("Deposited %.2f, new balance: %.2f\n", amount, w.Balance)
	return nil
}

func (w *VirtualWallet) Withdraw(amount float64) error {
	if amount > w.Balance || amount <= 0 {
		return errors.New("invalid withdrawal amount")
	}
	w.Balance -= amount
	fmt.Printf("Withdrawn %.2f, new balance: %.2f\n", amount, w.Balance)
	return nil
}

func (w *VirtualWallet) GetAssetCount(assetType string) (int, error) {
	count, exists := w.Assets[assetType]
	if !exists {
		return 0, errors.New("asset type not found")
	}
	return count, nil
}
