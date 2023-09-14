package address

import (
	"testing"

	"github.com/btcsuite/btcd/btcutil/bech32"
)

func FuzzAddr(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		bech32Data, _ := bech32.EncodeM("hrp", data)
		DecodeAddress(bech32Data, &MainNetTap)
	})
}
