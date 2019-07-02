package risevision_addrdec

import (
	"crypto/sha256"
)

var (
	Default = AddressDecoderV2{}
)

//AddressDecoderV2
type AddressDecoderV2 struct {
	IsTestNet bool
}

// GetAddressFromPublicKey takes a Lisk public key and returns the associated address
func GetAddressFromPublicKey(publicKey []byte) string {
	publicKeyHash := sha256.Sum256(publicKey)

	return GetBigNumberStringFromBytes(GetFirstEightBytesReversed(publicKeyHash[:sha256.Size])) + "R"
}

//AddressEncode 地址编码
func (dec *AddressDecoderV2) AddressEncode(hash []byte, opts ...interface{}) (string, error) {
	address := GetAddressFromPublicKey(hash)
	return address, nil
}
