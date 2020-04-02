package tables

import (
	"crypto/rand"
	"math/big"
)

func randStr(N int, charSet []byte) string {
	b := make([]byte, N)

	for i := range b {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(charSet))))
		if err != nil {
			panic(err)
		}

		b[i] = charSet[n.Int64()]
	}

	return string(b)
}

func randEmail() string {
	return randSlug() + "." + randSlug() + "@suburbia.io"
}

var slugChars = []byte("abcdefghijklmnopqrstuvwxyz")

func randSlug() string {
	return randStr(8, slugChars)
}

var countryChars = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randCountryCode() string {
	return randStr(2, slugChars)
}
