package transactions

import (
	"fmt"
	"testing"
)

func TestNewTransaction(t *testing.T) {
	 ts, err := NewTransaction("7061235840388139206L", 1000, "secret", "", 0)
	 if err != nil {
		t.Errorf("NewTransaction() returns error: %v, nil; expected transaction", err)
	}
	str,_ := ts.MarshalJSON()
	fmt.Println( string(str))

}

func TestNewTransactionWithData(t *testing.T) {
	if _, err := NewTransactionWithData("", 0, "", "c", 0, "abc"); err != nil {
		t.Errorf("NewTransactionWithData() returns error: %v, nil; expected transaction", err)
	}

	if _, err := NewTransactionWithData("", 0, "", "c", 0, []byte("abc")); err != nil {
		t.Errorf("NewTransactionWithData() returns error: %v, nil; expected transaction", err)
	}

	if val, err := NewTransactionWithData("", 0, "", "c", 0, 0); err == nil {
		t.Errorf("NewTransactionWithData() returns wrong data: %v, nil; expected error", val)
	}
}

func TestNewSecondSignatureTransaction(t *testing.T) {
	if _, err := NewSecondSignatureTransaction("", "", "abc", 0); err != nil {
		t.Errorf("NewSecondSignatureTransaction() returns error: %v, nil; expected transaction", err)
	}
}

func TestNewVoteTransaction(t *testing.T) {
	if _, err := NewVoteTransaction("", "", "c", 0, [][]byte{defaultSenderPublicKey}, [][]byte{}); err != nil {
		t.Errorf("NewVoteTransaction() returns error: %v, nil; expected transaction", err)
	}

	if val, err := NewVoteTransaction("", "", "c", 0, [][]byte{[]byte("abc")}, [][]byte{}); err == nil {
		t.Errorf("NewVoteTransaction() returns wrong data: %v, nil; expected error", val)
	}
}

func TestNewMultisignatureRegistrationTransaction(t *testing.T) {
	if _, err := NewMultisignatureRegistrationTransaction("", "", "c", 0, [][]byte{defaultSenderPublicKey}, [][]byte{}, 0, 0); err != nil {
		t.Errorf("NewMultisignatureRegistrationTransaction() returns error: %v, nil; expected transaction", err)
	}

	if val, err := NewMultisignatureRegistrationTransaction("", "", "c", 0, [][]byte{[]byte("abc")}, [][]byte{}, 0, 0); err == nil {
		t.Errorf("NewMultisignatureRegistrationTransaction() returns wrong data: %v, nil; expected error", val)
	}
}
