package transactions

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/blocktree/go-owcrypt"
	"github.com/blocktree/lisk-adapter/crypto"
	"testing"
)

var (
	defaultPrivateKey, _ = hex.DecodeString("2bb80d537b1da3e38bd30361aa855686bde0eacd7162fef6a25fe97bf527a25b5d036a858ce89f844491762eb89e2bfbd50a4a0a0da658e4b2628b25b117ae09")
)

func TestTransaction_Sign(t *testing.T) {
	transaction := Transaction{SenderPublicKey: defaultSenderPublicKey}
	err := transaction.Sign(defaultPrivateKey)
	//hex,_ := transaction.Hash()
	//hex = []byte(string("test"))
	//hex,_ = transaction.Hash()

	//hex,_ = transaction.Hash()
	hex := "text"
	sig, _ := owcrypt.Signature(defaultPrivateKey, nil, 0, []byte(hex), uint16(len([]byte(hex))), owcrypt.ECC_CURVE_ED25519)
	sig2 := crypto.SignMessageWithPrivateKey(string(hex), defaultPrivateKey)
	fmt.Println("content:"+hex)
	fmt.Println("pub:"+"5d036a858ce89f844491762eb89e2bfbd50a4a0a0da658e4b2628b25b117ae09")
	fmt.Println("pri:"+ "2bb80d537b1da3e38bd30361aa855686bde0eacd7162fef6a25fe97bf527a25b5d036a858ce89f844491762eb89e2bfbd50a4a0a0da658e4b2628b25b117ae09")
	//fmt.Println("sig: "+ base64.StdEncoding.EncodeToString(sig))
	fmt.Println("sig2:"+ base64.StdEncoding.EncodeToString(sig2))

	res := owcrypt.Verify(defaultSenderPublicKey, nil, 0, []byte(hex),uint16(len([]byte(hex))), sig2, owcrypt.ECC_CURVE_ED25519)
	fmt.Println("res2: ",res)
	if base64.StdEncoding.EncodeToString(transaction.Signature) != "dPFiRMXaoW8zAooOemp6sGv9BR6obHnjHOdWmg28n5QzJzs85+sNIvJpLzOxOq3NnCFqbnEChsdzCgMyandbCQ==" || err != nil {
		t.Errorf("Transaction.Sign() generates wrong Signature: %v; error: %v", base64.StdEncoding.EncodeToString(transaction.Signature), err)
	}
	if base64.StdEncoding.EncodeToString(sig) != "dPFiRMXaoW8zAooOemp6sGv9BR6obHnjHOdWmg28n5QzJzs85+sNIvJpLzOxOq3NnCFqbnEChsdzCgMyandbCQ==" || err != nil {
		t.Errorf("Transaction.Sign() generates wrong Signature: %v; error: %v", base64.StdEncoding.EncodeToString(transaction.Signature), err)
	}
	if base64.StdEncoding.EncodeToString(sig2) != "dPFiRMXaoW8zAooOemp6sGv9BR6obHnjHOdWmg28n5QzJzs85+sNIvJpLzOxOq3NnCFqbnEChsdzCgMyandbCQ==" || err != nil {
		t.Errorf("Transaction.Sign() generates wrong Signature: %v; error: %v", base64.StdEncoding.EncodeToString(transaction.Signature), err)
	}
}

func TestTransaction_SecondSign(t *testing.T) {
	transaction := Transaction{SenderPublicKey: defaultSenderPublicKey}
	err := transaction.SecondSign(defaultPrivateKey)

	if base64.StdEncoding.EncodeToString(transaction.secondSignature) != "dPFiRMXaoW8zAooOemp6sGv9BR6obHnjHOdWmg28n5QzJzs85+sNIvJpLzOxOq3NnCFqbnEChsdzCgMyandbCQ==" || err != nil {
		t.Errorf("Transaction.SecondSign() generates wrong Signature: %v; error: %v", base64.StdEncoding.EncodeToString(transaction.Signature), err)
	}
}


