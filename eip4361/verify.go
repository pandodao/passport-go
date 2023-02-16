package eip4361

import (
	"fmt"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

func Verify(msg *Message, signature string) error {
	sig, err := hexutil.Decode(signature)
	if err != nil {
		return fmt.Errorf("decode signature: %w", err)
	}

	if len(sig) != crypto.SignatureLength {
		return fmt.Errorf("signature must be %d bytes long", crypto.SignatureLength)
	}

	// comment(storyicon): fix ledger wallet
	// https://ethereum.stackexchange.com/questions/103307/cannot-verifiy-a-signature-produced-by-ledger-in-solidity-using-ecrecover
	if sig[crypto.RecoveryIDOffset] == 0 || sig[crypto.RecoveryIDOffset] == 1 {
		sig[crypto.RecoveryIDOffset] += 27
	}

	if sig[crypto.RecoveryIDOffset] != 27 && sig[crypto.RecoveryIDOffset] != 28 {
		return fmt.Errorf("invalid Ethereum signature (V is not 27 or 28)")
	}
	sig[crypto.RecoveryIDOffset] -= 27 // Transform yellow paper V from 27/28 to 0/1

	pk, err := crypto.SigToPub(accounts.TextHash(formatMessage(msg)), sig)
	if err != nil {
		return err
	}

	if addr := crypto.PubkeyToAddress(*pk); addr.String() != msg.Address {
		return fmt.Errorf("recoved address not match")
	}

	return nil
}
