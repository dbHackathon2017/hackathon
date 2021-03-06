package primitives_test

import (
	"testing"

	. "github.com/dbHackathon2017/hackathon/common/primitives"
)

func TestPrivateKey(t *testing.T) {
	p, err := GeneratePrivateKey()
	if err != nil {
		t.Error(err)
	}

	sec := p.Secret[:32]
	p2 := new(PrivateKey)
	err = p2.SetBytes(sec)
	if err != nil {
		t.Error(err)
	}

	if !p2.Public.IsSameAs(&p.Public) {
		t.Error("Should be same")
	}

	var _, _ = p, sec
}

func TestPrivateKeyMarshal(t *testing.T) {
	for i := 0; i < 1000; i++ {
		h, _ := RandomPrivateKey()
		data, _ := h.MarshalBinary()

		n := new(PrivateKey)
		newData, err := n.UnmarshalBinaryData(data)
		if err != nil {
			t.Error(err)
		}

		if !h.IsSameAs(n) {
			t.Error("Failed, should be same")
		}

		if len(newData) != 0 {
			t.Error("Failed, should have no bytes left")
		}

		if h.Empty() {
			t.Error("Should not be empty")
		}
	}
}

func TestPrivateKeyDiff(t *testing.T) {
	same := 0
	for i := 0; i < 1000; i++ {
		a, _ := RandomPrivateKey()
		b, _ := RandomPrivateKey()
		if a.IsSameAs(b) {
			same++
		}
	}
	if same > 15 {
		t.Error("More than 15 are the same, it is totally random, so it is likely the IsSameAs is broken.")
	}
}

func TestBadUnmarshalPrK(t *testing.T) {
	badData := []byte{}

	n := new(PrivateKey)

	_, err := n.UnmarshalBinaryData(badData)
	if err == nil {
		t.Error("Should panic or error out")
	}
}
