package eip4361

import (
	"bytes"
	"io"
	"unicode/utf8"

	"github.com/fkgi/abnf"
	"github.com/spf13/cast"
)

func Parse(msg string) (*Message, error) {
	r := bytes.NewReader([]byte(msg))

	domain, err := parseDomain(r)
	if err != nil {
		return nil, err
	}

	address, err := parseAddress(r)
	if err != nil {
		return nil, err
	}

	statement, err := parseStatement(r)
	if err != nil {
		return nil, err
	}

	uri, err := parseURI(r)
	if err != nil {
		return nil, err
	}

	ver, err := parseVersion(r)
	if err != nil {
		return nil, err
	}

	chainID, err := parseChainID(r)
	if err != nil {
		return nil, err
	}

	nonce, err := parseNonce(r)
	if err != nil {
		return nil, err
	}

	issuedAt, err := parseIssuedAt(r)
	if err != nil {
		return nil, err
	}

	expirationTime, err := parseExpirationTime(r)
	if err != nil {
		return nil, err
	}

	notBefore, err := parseNotBefore(r)
	if err != nil {
		return nil, err
	}

	requestID, err := parseRequestID(r)
	if err != nil {
		return nil, err
	}

	resources, err := parseResources(r)
	if err != nil {
		return nil, err
	}

	return &Message{
		Domain:         domain,
		Address:        address,
		Statement:      statement,
		URI:            uri,
		Version:        cast.ToInt(ver),
		ChainID:        cast.ToInt(chainID),
		Nonce:          nonce,
		IssuedAt:       cast.ToTime(issuedAt),
		ExpirationTime: cast.ToTime(expirationTime),
		NotBefore:      cast.ToTime(notBefore),
		RequestID:      requestID,
		Resources:      resources,
	}, nil
}

func isOptional(key Key) bool {
	switch key {
	case Statement, ExpirationTime, NotBefore, RequestID, Resources:
		return true
	default:
		return false
	}
}

func sizeOfRunes(runes []rune) int {
	var size int
	for _, r := range runes {
		size += utf8.RuneLen(r)
	}

	return size
}

func parseValues(r *bytes.Reader, key Key, f abnf.Rule) ([]string, error) {
	size := r.Len()
	t := abnf.ParseReader(r, size, size, f)
	if t == nil {
		return nil, key.invalidFormat()
	}

	subs := t.Children(int(key))
	values := make([]string, len(subs))

	for idx, sub := range subs {
		v := string(sub.V)
		values[idx] = v
	}

	if len(values) == 0 && !isOptional(key) {
		return nil, key.invalidFormat()
	}

	if x := size - r.Len() - sizeOfRunes(t.V); x != 0 {
		r.Seek(int64(-x), io.SeekCurrent)
	}

	return values, nil
}

func parseValue(r *bytes.Reader, key Key, f abnf.Rule) (string, error) {
	values, err := parseValues(r, key, f)
	if err != nil {
		return "", err
	}

	if len(values) > 0 {
		return values[0], nil
	}

	return "", nil
}

func parseDomain(r *bytes.Reader) (string, error) {
	key := Domain
	f := abnf.C(
		abnf.K(abnf.R1(abnf.VCHAR()), int(key)),
		abnf.SP(),
		abnf.VS("wants you to sign in with your Ethereum account:"),
		abnf.LF(),
	)

	return parseValue(r, key, f)
}

func parseAddress(r *bytes.Reader) (string, error) {
	key := Address

	hex := abnf.A(
		abnf.HEXDIG(),
		abnf.VR(0x61, 0x66),
	)

	f := abnf.C(
		abnf.K(
			abnf.C(
				abnf.VS("0x"),
				abnf.RN(40, hex),
			),
			int(key),
		),
		abnf.LF(),
	)

	return parseValue(r, key, f)
}

func parseStatement(r *bytes.Reader) (string, error) {
	key := Statement

	f := abnf.C(
		abnf.LF(),
		abnf.O(
			abnf.C(
				abnf.K(abnf.R1(abnf.A(abnf.VCHAR(), abnf.SP())), int(key)),
				abnf.LF(),
			),
		),
		abnf.LF(),
	)

	return parseValue(r, key, f)
}

func parseURI(r *bytes.Reader) (string, error) {
	key := URI

	f := abnf.C(
		abnf.VS("URI: "),
		abnf.K(abnf.R1(abnf.VCHAR()), int(key)),
		abnf.LF(),
	)

	return parseValue(r, key, f)
}

func parseVersion(r *bytes.Reader) (string, error) {
	key := Version

	f := abnf.C(
		abnf.VS("Version: "),
		abnf.K(abnf.R1(abnf.DIGIT()), int(key)),
		abnf.LF(),
	)

	return parseValue(r, key, f)
}

func parseChainID(r *bytes.Reader) (string, error) {
	key := ChainID

	f := abnf.C(
		abnf.VS("Chain ID: "),
		abnf.K(abnf.R1(abnf.DIGIT()), int(key)),
		abnf.LF(),
	)

	return parseValue(r, key, f)
}

func parseNonce(r *bytes.Reader) (string, error) {
	key := Nonce

	f := abnf.C(
		abnf.VS("Nonce: "),
		abnf.K(abnf.RV(8, -1, abnf.ALPHANUM()), int(key)),
		abnf.LF(),
	)

	return parseValue(r, key, f)
}

func parseIssuedAt(r *bytes.Reader) (string, error) {
	key := IssuedAt

	f := abnf.C(
		abnf.VS("Issued At: "),
		abnf.K(abnf.R1(abnf.VCHAR()), int(key)),
	)

	return parseValue(r, key, f)
}

func parseExpirationTime(r *bytes.Reader) (string, error) {
	key := ExpirationTime

	f := abnf.O(abnf.C(
		abnf.LF(),
		abnf.VS("Expiration Time: "),
		abnf.K(abnf.R1(abnf.VCHAR()), int(key)),
	))

	return parseValue(r, key, f)
}

func parseNotBefore(r *bytes.Reader) (string, error) {
	key := NotBefore

	f := abnf.O(abnf.C(
		abnf.LF(),
		abnf.VS("Not Before: "),
		abnf.K(abnf.R1(abnf.VCHAR()), int(key)),
	))

	return parseValue(r, key, f)
}

func parseRequestID(r *bytes.Reader) (string, error) {
	key := RequestID

	f := abnf.O(abnf.C(
		abnf.LF(),
		abnf.VS("Request ID: "),
		abnf.K(abnf.R1(abnf.VCHAR()), int(key)),
	))

	return parseValue(r, key, f)
}

func parseResources(r *bytes.Reader) ([]string, error) {
	key := Resources

	resources := abnf.O(abnf.C(
		abnf.LF(),
		abnf.VS("Resources:"),
		abnf.R0(abnf.C(
			abnf.LF(),
			abnf.VS("- "),
			abnf.K(abnf.R1(abnf.VCHAR()), int(key)),
		)),
	))

	f := abnf.C(resources, abnf.ETX())
	return parseValues(r, key, f)
}
