package crypto

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/blocktree/go-owcrypt"
	"testing"
)

var (
	defaultMessage      = "Some default text."
	signPublicKey, _    = hex.DecodeString("7ef45cd525e95b7a86244bbd4eb4550914ad06301013958f4dd64d32ef7bc588")
	signPrivateKey, _   = hex.DecodeString("314852d7afb0d4c283692fef8a2cb40e30c7a5df2ed79994178c10ac168d6d977ef45cd525e95b7a86244bbd4eb4550914ad06301013958f4dd64d32ef7bc588")
	defaultSignature, _ = hex.DecodeString("974eeac2c7e7d9da42aa273c8caae8e6eb766fa29a31b37732f22e6d2e61e8402106849e61e3551ff70d7d359170a6198669e1061b6b4aa61997e26b87e3a704")
	wrongSignature, _   = hex.DecodeString("974f2ac2c7e7d9da42aa273c8caae8e6eb766fa29a31b37732f22e6d2e61e8402106849e61e3551ff70d7d359170a6198669e1061b6b4aa61997e26b87e3a704")
	myPub ,_ = hex.DecodeString("5d036a858ce89f844491762eb89e2bfbd50a4a0a0da658e4b2628b25b117ae09")
	sig, _ = owcrypt.Signature([]byte{0x68,0xB2,0x11,0xB2,0xC0,0x1C,0xC8,0x86,0x90,0xBA,0x76,0xA0,0x78,0x95,0xA5,0xB4,0x80,0x5E,0x1C,0x11,0xFD,0xD3,0xAF,0x4C,0x86,0x3E,0x6D,0x4E,0xFE,0xB1,0x43,0x78}, nil, 0, []byte(defaultMessage), uint16(len([]byte(defaultMessage))), owcrypt.ECC_CURVE_ED25519)

	)

func TestSignMessageWithPrivateKey(t *testing.T) {
	if val := SignMessageWithPrivateKey(defaultMessage, signPrivateKey); !bytes.Equal(val, defaultSignature) {
		t.Errorf("SignMessageWithPrivateKey(%v,%v)=%v; want %v", defaultMessage, signPrivateKey, val, defaultSignature)
	}
}

func TestSignDataWithPrivateKey(t *testing.T) {
	if val := SignDataWithPrivateKey([]byte(defaultMessage), signPrivateKey); !bytes.Equal(val, defaultSignature) {
		t.Errorf("TestSignDataWithPrivateKey(%v,%v)=%v; want %v", []byte(defaultMessage), signPrivateKey, val, defaultSignature)
	}
}

func TestVerifyMessageWithPublicKey(t *testing.T) {
	if val, err := VerifyMessageWithPublicKey(defaultMessage, defaultSignature, signPublicKey); !val || err != nil {
		t.Errorf("SignMessageWithPrivateKey(%v,%v,%v)=%v,%v; want %v,%v", defaultMessage, defaultSignature, signPrivateKey, val, err, true, nil)
	}

	val,err := VerifyMessageWithPublicKey(defaultMessage, sig, myPub)
	if val || err != nil {
		t.Errorf("SignMessageWithPrivateKey(%v,%v,%v)=%v,%v; want %v,%v", defaultMessage, wrongSignature, signPrivateKey, val, err, false, nil)
	}
	fmt.Println("test:",val)
}

func TestVerifyDataWithPublicKey(t *testing.T) {
	if val, err := VerifyDataWithPublicKey([]byte(defaultMessage), defaultSignature, signPublicKey); !val || err != nil {
		t.Errorf("SignMessageWithPrivateKey(%v,%v,%v)=%v,%v; want %v,%v", []byte(defaultMessage), defaultSignature, signPrivateKey, val, err, true, nil)
	}

	if val, err := VerifyDataWithPublicKey([]byte(defaultMessage), wrongSignature, signPublicKey); val || err != nil {
		t.Errorf("SignMessageWithPrivateKey(%v,%v,%v)=%v,%v; want %v,%v", []byte(defaultMessage), wrongSignature, signPrivateKey, val, err, false, nil)
	}
}
