package common_test

import (
	"testing"

	. "github.com/dbHackathon2017/hackathon/common"
)

func TestPenMarshal(t *testing.T) {
	tr := RandomPension()
	data, err := tr.MarshalBinary()
	if err != nil {
		t.Error(err)
	}

	ta := new(Pension)
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
