package write

import (
	"bytes"
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
//			PensionID (from)
//			ToPensionID  (to)
//			PersonSubmit (who did it)
//			Timestamp
//			PubKey
//			Sig up to TS
//		Content
//			DocumentData
func SubmitValueChangeTransactionToPension(trans *common.Transaction, ec *factom.ECAddress, sigKey primitives.PrivateKey) (*primitives.Hash, error) {
	if trans.ValueChange == 0 {
		trans.UserType = constants.USER_TRANS_DOC_CHANGE
	} else {
		trans.UserType = constants.USER_TRANS_VAL_CHANGE
	}
	trans.FactomType = constants.FAC_TRANS_VAL_CHANGE
	return submitTransactionToFactom("Transaction Value Change", trans, ec, sigKey)

}

// Send A to B in A
func SubmitChainMoveTransactionToPension(trans *common.Transaction, ec *factom.ECAddress, sigKey primitives.PrivateKey) (*primitives.Hash, error) {
	trans.UserType = constants.USER_LIQUID_SEND
	trans.FactomType = constants.FAC_LIQUID_SEND
	return submitTransactionToFactom("Transaction Move Chain", trans, ec, sigKey)
}

// Request from A to B in B
func SubmitRequestMoveTransactionToPension(trans *common.Transaction, ec *factom.ECAddress, sigKey primitives.PrivateKey) (*primitives.Hash, error) {
	trans.UserType = constants.USER_LIQUID_REQUEST
	trans.FactomType = constants.FAC_LIQUID_REQUEST
	return submitTransactionToFactom("Transaction Request Move Chain", trans, ec, sigKey)
}

// Not used
func SubmitConfirmMoveTransactionToPension(trans *common.Transaction, ec *factom.ECAddress, sigKey primitives.PrivateKey) (*primitives.Hash, error) {
	trans.UserType = constants.USER_LIQUID_CONFIRMED
	trans.FactomType = constants.FAC_LIQUID_CONFIRM
	return submitTransactionToFactom("Transaction Confirm Move Chain", trans, ec, sigKey)
}

func submitTransactionToFactom(message string, trans *common.Transaction, ec *factom.ECAddress, sigKey primitives.PrivateKey) (*primitives.Hash, error) {
	e := new(factom.Entry)

	ut := primitives.Uint32ToBytes(trans.UserType)

	neg := trans.ValueChange < 0

	ft := primitives.Uint32ToBytes(trans.FactomType)
	person, err := trans.Person.MarshalBinary()
	if err != nil {
		return nil, err
	}

	var val uint32
	if trans.ValueChange < 0 {
		val = uint32(-1 * trans.ValueChange)
	} else {
		val = uint32(trans.ValueChange)
	}

	vb := primitives.Uint32ToBytes(val)
	if neg {
		vb = append([]byte{0x01}, vb...)
	} else {
		vb = append([]byte{0x00}, vb...)
	}

	ts, err := trans.Timestamp.MarshalBinary()
	if err != nil {
		return nil, err
	}

	e.ExtIDs = append(e.ExtIDs, []byte(message))           // 0
	e.ExtIDs = append(e.ExtIDs, ut)                        // 1
	e.ExtIDs = append(e.ExtIDs, ft)                        // 2
	e.ExtIDs = append(e.ExtIDs, vb)                        // 3
	e.ExtIDs = append(e.ExtIDs, trans.PensionID.Bytes())   // 4
	e.ExtIDs = append(e.ExtIDs, trans.ToPensionID.Bytes()) // 5
	e.ExtIDs = append(e.ExtIDs, person)                    // 6
	e.ExtIDs = append(e.ExtIDs, ts)                        // 7

	buf := new(bytes.Buffer)
	for i := 0; i < 8; i++ {
		buf.Write(e.ExtIDs[i])
	}

	msg := buf.Next(buf.Len())
	sig := sigKey.Sign(msg)

	e.ExtIDs = append(e.ExtIDs, sigKey.Public.Bytes()) // 8
	e.ExtIDs = append(e.ExtIDs, sig)                   // 9

	trans.Docs.FixFiles()

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

	trans.TransactionID = *ehash

	return ehash, nil
}
