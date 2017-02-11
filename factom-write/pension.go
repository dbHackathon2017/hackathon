package write

import (
	"fmt"
	"github.com/FactomProject/factom"
	"github.com/dbHackathon2017/hackathon/common"
	"github.com/dbHackathon2017/hackathon/common/constants"
	"github.com/dbHackathon2017/hackathon/common/primitives"
	"time"
)

var _ = fmt.Sprintf("")
var _ = time.Second
var _ = constants.CHAIN_PREFIX

// SubmitPensionToFactom makes a pension chain
// Factom Chain
//		"Pension Chain"
// 		Company
//		PubKey
//		Nonce
func SubmitPensionToFactom(pen *common.Pension, ec *factom.ECAddress) (*primitives.Hash, error) {
	e := new(factom.Entry)

	comp, err := pen.Company.MarshalBinary()
	if err != nil {
		return nil, err
	}

	e.ExtIDs = append(e.ExtIDs, []byte("Pension Chain")) // 0
	e.ExtIDs = append(e.ExtIDs, comp)                    // 1
	e.ExtIDs = append(e.ExtIDs, pen.AuthKey.Bytes())     // 2
	e.ExtIDs = append(e.ExtIDs, pen.UniqueHash.Bytes())  // 3
	nonce := FindValidNonce(e)
	e.ExtIDs = append(e.ExtIDs, nonce) // 4

	c := factom.NewChain(e)

	_, err = factom.CommitChain(c, ec)
	if err != nil {
		return nil, err
	}

	// We can send this on a go routine
	//go func() {
	// fmt.Println("Revealed!")
	//	time.Sleep(constants.REVEAL_WAIT)
	_, _ = factom.RevealChain(c)
	// fmt.Println(str, err)
	//}()

	chainID, err := primitives.HexToHash(c.FirstEntry.ChainID)
	if err != nil {
		return nil, err
	}
	pen.PensionID = *chainID

	return chainID, nil
}
