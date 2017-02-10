package common

import (
	"time"
	"pension"
	"github.com/dbHackathon2017/hackathon/common/constants"
	"github.com/dbHackathon2017/hackathon/common/primitives"
	"github.com/dbHackathon2017/hackathon/common/primitives/random"
)

// Transaction is an individual transaction in a pension chain. Each transaction
// has a type. This type is only relevent for readability.
// All transactions are processed the same way, they affect the total value of a chain.
// All supporting documents, text, and other metadata is only used by users.
//
// For Example: If a pension wants to change the address of a pension. They would
// 		make a transaction of type DOC_CHANGE and add hashes of the supporting documents.
//		The transaction would add $0.00 value. We would process it the same was as a $0.00
//		value one.
type Transaction struct {
	PensionID   primitives.Hash
	TransactionID primitives.Hash
	ValueChange int    // In every transaction
	FactomType  uint32 // Transaction type relative to factom
	UserType    uint32 // Transaction type relative to user
	Timestamp   time.Time

	Docs   primitives.FileList   // List of docs and some metadata
	Person primitives.PersonName // Name of person who made change to factom.
}

func RandomValChangeTransaction(pensionID *primitives.Hash) *Transaction {
	t := new(Transaction)
	t.PensionID = *pensionID
	t.TransactionID = *primitives.RandomHash()
	t.ValueChange = 100
	t.FactomType = constants.FAC_TRANS_VAL_CHANGE
	t.UserType = constants.USER_TRANS_VAL_CHANGE
	t.Docs = *primitives.RandomFileList(10)
	t.Person = *primitives.RandomName()

	t.Timestamp = time.Now()
	sec := random.RandomInt64Between(0, 100000)
	day := random.RandomInt64Between(0, 1000)

	t.Timestamp.Add(time.Duration(sec) * time.Second)
	t.Timestamp.Add(time.Duration(day) * time.Hour)

	return t
}
