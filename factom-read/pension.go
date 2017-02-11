package read

import (
	"bytes"
	"fmt"
	"log"
	"time"

	"github.com/FactomProject/factom"
	"github.com/dbHackathon2017/hackathon/common"
	"github.com/dbHackathon2017/hackathon/common/constants"
	"github.com/dbHackathon2017/hackathon/common/primitives"
)

var _ = constants.FAC_LIQUID_SEND

func GetTransactionFromTxID(id primitives.Hash) (*common.Transaction, error) {
	ent, err := factom.GetEntry(id.String())
	fmt.Println(ent, err)
	if err != nil {
		return nil, err
	}
	t := buildValChangeTransactionsFromFactomEntry(ent, nil)
	if t == nil {
		return nil, fmt.Errorf("Entry errored")
	}

	return t, nil
}

func GetPensionFromFactom(id primitives.Hash) (*common.Pension, error) {
	ents, err := factom.GetAllChainEntries(id.String())
	var _ = ents
	if err != nil {
		return nil, err
	}

	p := new(common.Pension)
	p.Active = true
	// Need to grab the Pension chain enty
	for _, e := range ents {
		if bytes.Compare(e.ExtIDs[0], []byte("Pension Chain")) == 0 {
			p = buildPensionFromFactomEntry(e)
			if p != nil {
				break // Exit loop
			}
		}
	}

	if p == nil {
		return nil, fmt.Errorf("Could not build pension from chain")
	}

	transactions := make([]*common.Transaction, 0)
	// Now build transactions
	for _, e := range ents {
		if bytes.Compare(e.ExtIDs[0], []byte("Transaction Value Change")) == 0 {
			t := buildValChangeTransactionsFromFactomEntry(e, p)
			if t == nil {
				fmt.Println("Bad transaction")
				continue // Bad transaction
			} else {
				transactions = append(transactions, t)
			}
		} else if bytes.Compare(e.ExtIDs[0], []byte("Transaction Move Chain")) == 0 {
			t := applyMoveChain(e, p)
			if t == nil {
				fmt.Println("Bad Trans move transaction")
				continue // Bad transaction
			} else {
				transactions = append(transactions, t)
			}
		} else if bytes.Compare(e.ExtIDs[0], []byte("Transaction Request Chain")) == 0 {
			t := applyMoveChain(e, p)
			if t == nil {
				fmt.Println("Bad Request Trans move transaction")
				continue // Bad transaction
			} else {
				transactions = append(transactions, t)
			}
		}
	}

	p.Transactions = transactions

	return p, nil
}

func applyMoveChain(e *factom.Entry, p *common.Pension) *common.Transaction {
	t := applyTransaction(e, p)
	if t == nil {
		return nil
	}
	p.Value += t.ValueChange
	p.Active = false
	return t
}

func applyRequestChain(e *factom.Entry, p *common.Pension) *common.Transaction {
	t := applyTransaction(e, p)
	if t == nil {
		return nil
	}
	return t
}

func buildValChangeTransactionsFromFactomEntry(e *factom.Entry, p *common.Pension) *common.Transaction {
	t := applyTransaction(e, p)
	if t == nil {
		return nil
	}
	if p != nil {
		p.Value += t.ValueChange
	}
	return t
}

func applyTransaction(e *factom.Entry, p *common.Pension) *common.Transaction {
	if len(e.ExtIDs) != 9 {
		return nil
	}

	buf := new(bytes.Buffer)

	t := new(common.Transaction)
	ut, err := primitives.BytesToUint32(e.ExtIDs[1])
	if err != nil {
		log.Println("Usetype fail")
		return nil
	}
	t.UserType = ut

	valC, err := primitives.BytesToUint32(e.ExtIDs[2])
	if err != nil {
		log.Println("Valchange fail")
		return nil
	}
	t.ValueChange = int(valC)

	pid, err := primitives.BytesToHash(e.ExtIDs[3])
	if err != nil {
		log.Println("PID fail")
		return nil
	}
	t.PensionID = *pid

	if p != nil && !pid.IsSameAs(&p.PensionID) {
		return nil
	}

	toPid, err := primitives.BytesToHash(e.ExtIDs[4])
	if err != nil {
		log.Println("PID fail")
		return nil
	}
	t.ToPensionID = *toPid

	per := new(primitives.PersonName)
	err = per.UnmarshalBinary(e.ExtIDs[5])
	if err != nil {
		log.Println("PersonName fail")
		return nil
	}
	t.Person = *per

	ts := new(time.Time)
	err = ts.UnmarshalBinary(e.ExtIDs[6])
	if err != nil {
		log.Println("Timestamp fail")
		return nil
	}
	t.Timestamp = *ts

	pk, err := primitives.PublicKeyFromBytes(e.ExtIDs[7])
	if err != nil {
		log.Println("Pubkey fail")
		return nil
	}
	if p != nil && !pk.IsSameAs(&p.AuthKey) {
		log.Println("PubKey not same")
		return nil
	}

	for i := 0; i < 7; i++ {
		buf.Write(e.ExtIDs[i])
	}

	msg := buf.Next(buf.Len())
	sig := e.ExtIDs[8]

	valid := pk.Verify(msg, sig)
	if !valid {
		log.Println("Not valid")
		return nil
	}

	docs := new(primitives.FileList)
	err = docs.UnmarshalBinary(e.Content)
	if err != nil {
		log.Println("Doclist fail")
		return nil
	}

	t.Docs = *docs
	ehash, err := primitives.BytesToHash(e.Hash())
	if err != nil {
		log.Println("Fail to get ehash")
	}
	t.TransactionID = *ehash
	return t
}

func buildPensionFromFactomEntry(e *factom.Entry) *common.Pension {
	if len(e.ExtIDs) != 5 {
		return nil
	}
	var err error
	p := new(common.Pension)

	comp := new(primitives.PersonName)
	err = comp.UnmarshalBinary(e.ExtIDs[1])
	if err != nil {
		log.Println("Comp fails")
		return nil
	}
	p.Company = *comp

	ak, err := primitives.PublicKeyFromBytes(e.ExtIDs[2])
	if err != nil {
		log.Println("Pub key fails")
		return nil
	}
	p.AuthKey = *ak

	uh, err := primitives.BytesToHash(e.ExtIDs[3])
	if err != nil {
		log.Println("Pub key fails")
		return nil
	}
	p.UniqueHash = *uh

	pid, err := primitives.HexToHash(e.ChainID)
	if err != nil {
		log.Println("Pid fails")
		return nil
	}

	p.PensionID = *pid

	err = p.Docs.UnmarshalBinary(e.Content)
	if err != nil {
		return nil
	}

	return p
}
