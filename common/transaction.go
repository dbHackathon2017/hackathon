package common

import (
	"bytes"
	"fmt"
	"time"

	"github.com/dbHackathon2017/hackathon/common/constants"
	"github.com/dbHackathon2017/hackathon/common/primitives"
	"github.com/dbHackathon2017/hackathon/common/primitives/random"
)

var _ = fmt.Sprintf("")

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
	PensionID     primitives.Hash
	ToPensionID   primitives.Hash
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
	t.ToPensionID = *pensionID
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

func (t *Transaction) UnmarshalBinary(data []byte) error {
	_, err := t.UnmarshalBinaryData(data)
	return err
}

func (t *Transaction) UnmarshalBinaryData(data []byte) (newData []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("[Pension] A panic has occurred while unmarshaling: %s", r)
			return
		}
	}()
	newData = data

	newData, err = t.PensionID.UnmarshalBinaryData(newData)
	if err != nil {
		return data, err
	}

	newData, err = t.ToPensionID.UnmarshalBinaryData(newData)
	if err != nil {
		return data, err
	}

	newData, err = t.TransactionID.UnmarshalBinaryData(newData)
	if err != nil {
		return data, err
	}

	u, err := primitives.BytesToUint32(newData[:4])
	if err != nil {
		return data, err
	}
	t.ValueChange = int(u)
	newData = newData[4:]

	u, err = primitives.BytesToUint32(newData[:4])
	if err != nil {
		return data, err
	}
	t.FactomType = u
	newData = newData[4:]

	u, err = primitives.BytesToUint32(newData[:4])
	if err != nil {
		return data, err
	}
	t.UserType = u
	newData = newData[4:]

	err = t.Timestamp.UnmarshalBinary(newData[:15])
	if err != nil {
		return data, err
	}
	newData = newData[15:]

	newData, err = t.Docs.UnmarshalBinaryData(newData)
	if err != nil {
		return data, err
	}

	newData, err = t.Person.UnmarshalBinaryData(newData)
	if err != nil {
		return data, err
	}

	return newData, nil
}

func (t *Transaction) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)

	data, err := t.PensionID.MarshalBinary()
	if err != nil {
		return nil, err
	}
	buf.Write(data)

	data, err = t.ToPensionID.MarshalBinary()
	if err != nil {
		return nil, err
	}
	buf.Write(data)

	data, err = t.TransactionID.MarshalBinary()
	if err != nil {
		return nil, err
	}
	buf.Write(data)

	data = primitives.Uint32ToBytes(uint32(t.ValueChange))
	buf.Write(data)

	data = primitives.Uint32ToBytes(t.FactomType)
	buf.Write(data)

	data = primitives.Uint32ToBytes(t.UserType)
	buf.Write(data)

	data, err = t.Timestamp.MarshalBinary()
	if err != nil {
		return nil, err
	}
	buf.Write(data)

	data, err = t.Docs.MarshalBinary()
	if err != nil {
		return nil, err
	}
	buf.Write(data)

	data, err = t.Person.MarshalBinary()
	if err != nil {
		return nil, err
	}
	buf.Write(data)

	return buf.Next(buf.Len()), nil
}

func (t *Transaction) GetTimeStampFormatted() string {
	return t.Timestamp.Format(layout)
}

func (t *Transaction) GetUserTypeString() string {
	if t.UserType == constants.USER_LIQUID_SEND {
		return "Chain Liquidation"
	} else if t.UserType == constants.USER_TRANS_DOC_CHANGE {
		return "Document Change"
	} else if t.UserType == constants.USER_TRANS_VAL_CHANGE {
		if t.ValueChange < 0 {
			return "Withdraw"
		}
		return "Deposit"
	} else if t.UserType == constants.USER_LIQUID_REQUEST {
		return "Merge Finalized"
	} else if t.UserType == constants.USER_LIQUID_CONFIRMED {
		return "Merge Finalized"
	}
	return "--"
}

func (t *Transaction) GetFactomTypeString() string {
	if t.UserType == constants.FAC_LIQUID_SEND {
		return "Chain Liquidation"
	} else if t.UserType == constants.FAC_TRANS_VAL_CHANGE {
		return "ValueChange"
	} else if t.UserType == constants.FAC_LIQUID_REQUEST {
		return "Retquest Merge In"
	} else if t.UserType == constants.FAC_LIQUID_CONFIRM {
		return "Merge Finalized"
	}
	return "--"
}

func (a *Transaction) IsSameAs(b *Transaction) bool {
	if !a.PensionID.IsSameAs(&b.PensionID) {
		return false
	}

	if !a.ToPensionID.IsSameAs(&b.ToPensionID) {
		return false
	}

	if !a.TransactionID.IsSameAs(&b.TransactionID) {
		return false
	}

	if a.ValueChange != b.ValueChange {
		return false
	}

	if a.FactomType != b.FactomType {
		return false
	}

	if a.UserType != b.UserType {
		return false
	}

	if !a.Docs.IsSameAs(&b.Docs) {
		return false
	}

	if !a.Person.IsSameAs(&b.Person) {
		return false
	}

	return true
}

type TransList []*Transaction

func (s TransList) Len() int {
	return len(s)
}
func (s TransList) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s TransList) Less(i, j int) bool {
	return s[i].Timestamp.Before(s[j].Timestamp)
}
