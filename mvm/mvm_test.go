package mvm

import (
	"context"
	"testing"
)

func TestFetchAddress(t *testing.T) {
	ctx := context.Background()
	userID := "e8e8cd79-cd40-4796-8c54-3a13cfe50115"
	addr, err := UserToAddress(ctx, userID)
	if err != nil {
		t.Fatal(err)
	}

	id, err := AddressToUser(ctx, addr)
	if err != nil {
		t.Fatal(err)
	}

	if id != userID {
		t.Fatalf("expected %s, got %s", userID, id)
	}
}
