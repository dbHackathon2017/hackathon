package write

import (
	"fmt"
	"time"

	"github.com/FactomProject/factom"
	"github.com/dbHackathon2017/hackathon/common"
	"github.com/dbHackathon2017/hackathon/common/constants"
	"github.com/dbHackathon2017/hackathon/common/primitives"
)

var _ = fmt.Sprintf("")
var _ = time.Second
var _ = constants.CHAIN_PREFIX

// Can make all transaction types and add them to a pension chain

// Factom Chain
//			"Transaction Value Change"
//			UserType
//			ValueChange
//			PensionID (Stop cross-chain-replay)
//			PersonSubmit (who did it)
//			Timestamp
//			PubKey
//			Sig up to TS
//		Content
//			DocumentData
func SubmitValueChangeTransactionToPension(trans common.Transaction, ec *factom.ECAddress, sigKey primitives.PrivateKey) (*primitives.Hash, error) {
	e := new(factom.Entry)

	ut := primitives.Uint32ToBytes(trans.UserType)
	vc := primitives.Uint32ToBytes(uint32(trans.ValueChange))
	person, err := trans.Person.MarshalBinary()
	if err != nil {
		return nil, err
	}

	ts, err := trans.Timestamp.MarshalBinary()
	if err != nil {
		return nil, err
	}

	e.ExtIDs = append(e.ExtIDs, []byte("Transaction Value Change")) // 0
	e.ExtIDs = append(e.ExtIDs, ut)                                 // 1
	e.ExtIDs = append(e.ExtIDs, vc)                                 // 2
	e.ExtIDs = append(e.ExtIDs, trans.PensionID.Bytes())            // 3
	e.ExtIDs = append(e.ExtIDs, person)                             // 4
	e.ExtIDs = append(e.ExtIDs, ts)                                 // 5

	msg := upToNonce(e.ExtIDs)
	sig := sigKey.Sign(msg)

	e.ExtIDs = append(e.ExtIDs, sigKey.Public.Bytes()) // 6
	e.ExtIDs = append(e.ExtIDs, sig)                   // 7

	docs, err := trans.Docs.MarshalBinary()
	if err != nil {
		return nil, err
	}

	e.Content = docs
	e.ChainID = trans.PensionID.String()

	_, err = factom.CommitEntry(e, ec)
	if err != nil {
		return nil, err
	}

	// We can send this on a go routine
	//go func() {
	// fmt.Println("Revealed!")
	//	time.Sleep(constants.REVEAL_WAIT)
	ehashStr, err := factom.RevealEntry(e)
	if err != nil {
		return nil, err
	}
	// fmt.Println(str, err)
	//}()

	ehash, err := primitives.HexToHash(ehashStr)
	if err != nil {
		return nil, err
	}

	return ehash, nil
}
