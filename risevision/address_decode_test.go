package risevision

import (
	"encoding/hex"
	"fmt"
	"testing"
)

func TestAddressDecoder_PublicKeyToAddress(t *testing.T) {
	pub, _ := hex.DecodeString("6e6490ba9ffa3ed276048e23c52f09a7622e02111124e9c770d1a6ac11a723c6")
	decoder := AddressDecoder{}
	addr, err := decoder.PublicKeyToAddress(pub, false)
	if err != nil {
		t.Errorf("PublicKeyToAddress error: %v", err)
		return
	}
	fmt.Println(addr)
}
