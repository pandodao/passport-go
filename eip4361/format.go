package eip4361

import (
	"fmt"
	"strings"
	"time"
)

func formatMessage(msg *Message) string {
	var b strings.Builder

	// domain
	fmt.Fprintf(&b, "%s wants you to sign in with your Ethereum account:\n", msg.Domain)
	// b.WriteByte(0x0a)

	// address
	b.WriteString(msg.Address)
	b.WriteByte('\n')

	// statement
	b.WriteByte('\n')
	if msg.Statement != "" {
		b.WriteString(msg.Statement)
		b.WriteByte('\n')
	}
	b.WriteByte('\n')

	// uri
	fmt.Fprintf(&b, "URI: %s\n", msg.URI)

	// version
	fmt.Fprintf(&b, "Version: %d\n", msg.Version)

	// chain id
	fmt.Fprintf(&b, "Chain ID: %d\n", msg.ChainID)

	// nonce
	fmt.Fprintf(&b, "Nonce: %s\n", msg.Nonce)

	// issued at
	fmt.Fprintf(&b, "Issued At: %s\n", msg.IssuedAt.Format(time.RFC3339))

	// expiration time
	if !msg.ExpirationTime.IsZero() {
		fmt.Fprintf(&b, "Expiration Time: %s\n", msg.ExpirationTime.Format(time.RFC3339))
	}

	// not before
	if !msg.NotBefore.IsZero() {
		fmt.Fprintf(&b, "Not Before: %s\n", msg.NotBefore.Format(time.RFC3339))
	}

	// request id
	if msg.RequestID != "" {
		fmt.Fprintf(&b, "Request ID: %s\n", msg.RequestID)
	}

	// resources
	if len(msg.Resources) > 0 {
		b.WriteString("Resources:")
		for _, resource := range msg.Resources {
			fmt.Fprintf(&b, "\n- %s", resource)
		}
	}

	return b.String()
}
