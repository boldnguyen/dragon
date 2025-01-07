package profile

import (
	"testing"
)

func TestNewProfile(t *testing.T) {
	p := NewProfile("Alice", 1)
	if p.Name != "Alice" {
		t.Errorf("Expected name Alice, got %s", p.Name)
	}
	if p.Level != 1 {
		t.Errorf("Expected level 1, got %d", p.Level)
	}
	if p.PlayerID != "Player_Alice" {
		t.Errorf("Expected PlayerID Player_Alice, got %s", p.PlayerID)
	}
}

func TestSetLevel(t *testing.T) {
	p := NewProfile("Bob", 1)
	p.SetLevel(5)
	if p.Level != 5 {
		t.Errorf("Expected level 5, got %d", p.Level)
	}
}

func TestChangeName(t *testing.T) {
	p := NewProfile("Charlie", 1)
	p.Tokens = 100
	err := p.ChangeName("NewCharlie", 50)
	if err != nil {
		t.Errorf("Failed to change name: %v", err)
	}
	if p.Name != "NewCharlie" {
		t.Errorf("Expected name NewCharlie, got %s", p.Name)
	}
	if p.Tokens != 50 {
		t.Errorf("Expected 50 tokens left, got %v", p.Tokens)
	}

	// Test for insufficient tokens
	err = p.ChangeName("AnotherTry", 60)
	if err == nil {
		t.Error("Changing name should have failed due to insufficient tokens")
	}
}

func TestChangeAvatar(t *testing.T) {
	p := NewProfile("David", 1)
	p.Tokens = 100
	err := p.ChangeAvatar("new_avatar", 50)
	if err != nil {
		t.Errorf("Failed to change avatar: %v", err)
	}
	if p.Avatar != "new_avatar" {
		t.Errorf("Expected avatar new_avatar, got %s", p.Avatar)
	}
	if p.Tokens != 50 {
		t.Errorf("Expected 50 tokens left, got %v", p.Tokens)
	}
}
