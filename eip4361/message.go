package eip4361

import (
	"fmt"
	"net/url"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cast"
)

type Message struct {
	// Domain is the RFC 3986 authority that is requesting the signing.
	Domain string
	// Address is the Ethereum address performing the signing conformant to capitalization encoded checksum specified in EIP-55 where applicable.
	Address string
	// Statement (optional) is a human-readable ASCII assertion that the user will sign, and it must not contain '\n' (the byte 0x0a).
	Statement string
	// URI is an RFC 3986 URI referring to the resource that is the subject of the signing (as in the subject of a claim).
	URI string
	// Version is the current version of the message, which MUST be 1 for this specification.
	Version int
	// ChainID is the EIP-155 Chain ID to which the session is bound, and the network where Contract Accounts MUST be resolved.
	ChainID int
	// Nonce is a randomized token typically chosen by the relying party and used to prevent replay attacks, at least 8 alphanumeric characters.
	Nonce string
	// IssuedAt is the ISO 8601 datetime string of the current time.
	IssuedAt time.Time
	// ExpirationTime (optional) is the ISO 8601 datetime string that, if present, indicates when the signed authentication message is no longer valid.
	ExpirationTime time.Time
	// NotBefore (optional) is the ISO 8601 datetime string that, if present, indicates when the signed authentication message will become valid.
	NotBefore time.Time
	// RequestID (optional) is an system-specific identifier that may be used to uniquely refer to the sign-in request.
	RequestID string
	// resources (optional) is a list of information or references to information the user wishes to have resolved as part of authentication by the relying party. They are expressed as RFC 3986 URIs separated by "\n- " where \n is the byte 0x0a.
	Resources []string
}

func (m *Message) String() string {
	return formatMessage(m)
}

func (m *Message) Validate() error {
	if err := validateDomain(m.Domain); err != nil {
		return err
	}

	if err := validateAddress(m.Address); err != nil {
		return err
	}

	if err := validateURI(m.URI); err != nil {
		return err
	}

	if err := validateVersion(cast.ToString(m.Version)); err != nil {
		return err
	}

	if m.IssuedAt.IsZero() {
		return fmt.Errorf("issued at must be set")
	}

	for _, resource := range m.Resources {
		if err := validateURI(resource); err != nil {
			return err
		}
	}

	return nil
}

func validateDomain(domain string) error {
	u, err := url.ParseRequestURI("https://" + domain)
	if err != nil {
		return err
	}

	s := u.Host
	if user := u.User.String(); user != "" {
		s = user + "@" + s
	}

	if s != domain {
		return fmt.Errorf("domain %q is not a valid RFC 3986 authority", domain)
	}

	return nil
}

func validateAddress(address string) error {
	if !common.IsHexAddress(address) {
		return fmt.Errorf("address %q is not a valid Ethereum address", address)
	}

	return nil
}

func validateURI(uri string) error {
	_, err := url.Parse(uri)
	return err
}

func validateVersion(version string) error {
	if cast.ToInt(version) != 1 {
		return fmt.Errorf("version %s is not supported", version)
	}

	return nil
}