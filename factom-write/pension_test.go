package write_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/FactomProject/factom"
	"github.com/dbHackathon2017/hackathon/common"
	"github.com/dbHackathon2017/hackathon/common/constants"
	"github.com/dbHackathon2017/hackathon/common/primitives"
	. "github.com/dbHackathon2017/hackathon/factom-write"
)

func TestSubmitPension(t *testing.T) {
	// Set to remote
	factom.SetFactomdServer(constants.REMOTE_HOST)

	// Get key for pension
	pk, _ := primitives.RandomPrivateKey()
	p := common.RandomPension()
	p.AuthKey = pk.Public

	// Get ECAddress
	ec := GetECAddress()

	c, err := SubmitPensionToFactom(p, ec)
	if err != nil {
		t.Error(err)
	}

	// Wait for chain to enter factom
	// chain must be in first
	time.Sleep(4 * time.Second)

	trans := common.RandomValChangeTransaction(primitives.RandomHash())
	trans.PensionID = p.PensionID

	ehash, err := SubmitValueChangeTransactionToPension(trans, ec, *pk)
	if err != nil {
		t.Error(err)
	}

	fmt.Println(c.String() + "\n" + ehash.String())
}
