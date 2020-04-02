package crypto

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"time"
)

func NewUUID() string {
	return NewUUIDat(time.Now())
}

func NewUUIDat(at time.Time) string {
	buf := make([]byte, 11)
	_, err := rand.Read(buf)
	if err != nil {
		panic(err)
	}

	s := fmt.Sprintf("%010x%022x", at.Unix(), buf)
	return s[:8] + "-" + s[8:12] + "-" + s[12:16] + "-" + s[16:20] + "-" + s[20:32]
}

var ReadableChars = []byte("ACDEFGHJKLMNPQRTWXY3469")
var AlphanumericChars = []byte("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz")

func RandAlphaNumForHumans(N int) string {
	return randStr(N, ReadableChars)
}

func RandAlphaNum(N int) string {
	return randStr(N, AlphanumericChars)
}

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
