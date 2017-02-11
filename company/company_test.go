package company_test

import (
	"testing"

	. "github.com/dbHackathon2017/hackathon/company"
)

func TestCompMarsahl(t *testing.T) {
	fc := RandomPenstionAndMetaData()
	data, err := fc.MarshalBinary()
	if err != nil {
		t.Error(err)
	}

	fc2 := new(PensionAndMetadata)
	nd, err := fc2.UnmarshalBinaryData(data)
	if err != nil {
		t.Error(err)
	}

	if !fc.IsSameAs(fc2) {
		t.Error("Should be same")
	}

	if len(nd) > 0 {
		t.Error("Should be 0")
	}
}
