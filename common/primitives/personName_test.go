package primitives_test

import (
	"fmt"
	"testing"

	. "github.com/dbHackathon2017/hackathon/common/primitives"
)

var _ = fmt.Sprintf("")

func TestNames(t *testing.T) {
	for i := 0; i < 1000; i++ {
		l := RandomName()
		data, err := l.MarshalBinary()
		if err != nil {
			t.Error(err)
		}

		n := new(PersonName)
		newData, err := n.UnmarshalBinaryData(data)
		if err != nil {
			t.Error(err)
		}
		if !n.IsSameAs(l) {
			t.Error("Should match.")
		}

		if len(newData) != 0 {
			t.Error("Failed, should have no bytes left")
		}

	}
}

func TestDiffNames(t *testing.T) {
	a := RandomName()
	data, err := a.MarshalBinary()
	if err != nil {
		t.Error(err)
	}

	b := new(PersonName)
	_, err = b.UnmarshalBinaryData(data)
	if err != nil {
		t.Error(err)
	}

	a.SetString("One")
	b.SetString("Two")

	if a.IsSameAs(b) {
		t.Error("Should be different")
	}
}

func TestBadUnmarshalNames(t *testing.T) {
	badData := []byte{}

	n := new(PersonName)

	_, err := n.UnmarshalBinaryData(badData)
	if err == nil {
		t.Error("Should panic or error out")
	}

	s := new(PersonName)
	_, err = s.UnmarshalBinaryData(badData)
	if err == nil {
		t.Error("Should panic or error out")
	}
}
