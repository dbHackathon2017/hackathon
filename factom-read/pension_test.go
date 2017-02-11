package read_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/FactomProject/factom"
	"github.com/dbHackathon2017/hackathon/common"
	"github.com/dbHackathon2017/hackathon/common/constants"
	"github.com/dbHackathon2017/hackathon/common/primitives"
	. "github.com/dbHackathon2017/hackathon/factom-read"
	"github.com/dbHackathon2017/hackathon/factom-write"
)

func TestReadPension(t *testing.T) {
	// Set to remote
	factom.SetFactomdServer(constants.REMOTE_HOST)

	// Get key for pension
	pk, _ := primitives.RandomPrivateKey()
	p := common.RandomPension()
	p.AuthKey = pk.Public

	// Get ECAddress
	ec := write.GetECAddress()

	c, err := write.SubmitPensionToFactom(p, ec)
	if err != nil {
		t.Error(err)
	}
	p.FixPids()

	// Wait for chain to enter factom
	// chain must be in first
	time.Sleep(4 * time.Second)

	trans := common.RandomValChangeTransaction(&p.PensionID)
	trans.PensionID = p.PensionID

	for i := 0; i < len(p.Transactions); i++ {
		_, err := write.SubmitValueChangeTransactionToPension(p.Transactions[i], ec, *pk)
		if err != nil {
			t.Error(err)
		}
		p.Value += p.Transactions[i].ValueChange
	}

	//p.Transactions = append(p.Transactions, *trans)

	fmt.Println(c.String() + "\n")
	//p.Value += trans.ValueChange

	time.Sleep(3 * time.Second)

	pen2, err := GetPensionFromFactom(p.PensionID)
	if err != nil {
		t.Error(err)
	}

	if !pen2.IsSameAs(p) {
		t.Error("Should be same")
	}
}
