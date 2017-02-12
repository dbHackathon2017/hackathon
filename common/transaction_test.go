package common_test

import (
	"testing"

	. "github.com/dbHackathon2017/hackathon/common"
	"github.com/dbHackathon2017/hackathon/common/primitives"
)

func TestTransMarshal(t *testing.T) {
	tr := RandomValChangeTransaction(primitives.RandomHash())
	data, err := tr.MarshalBinary()
	if err != nil {
		t.Error(err)
	}

	ta := new(Transaction)
	nd, err := ta.UnmarshalBinaryData(data)
	if err != nil {
		t.Error(err)
	}

	if len(nd) != 0 {
		t.Errorf("Should be 0, found %d", len(nd))
	}

	if !tr.IsSameAs(ta) {
		t.Error("Not same")
	}
}
