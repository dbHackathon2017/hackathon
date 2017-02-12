package common

import (
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"sort"

	"github.com/dbHackathon2017/hackathon/common/primitives"
	//"transaction"
)

const layout = "Jan 2, 2006"

// Pension structs contain all information about a pension relative to
// factom. This means it does not contain personal information (name/address/etc)
// instead all pensions are tied to a sha256 hash. This hash is it's location in
// factom.
type Pension struct {
	// ChainID in Factom. This is how we find a pension
	PensionID    primitives.Hash
	Transactions []*Transaction
	Company      primitives.PersonName // Company name, use same primitive
	UniqueHash   primitives.Hash
	Docs         primitives.FileList

	// The current amount of tokens in the pension.
	Value   int
	AuthKey primitives.PublicKey
	Active  bool
}

func RandomPension() *Pension {
	p := new(Pension)
	p.PensionID = *primitives.RandomHash()
	p.Transactions = make([]*Transaction, 0)
	for i := 0; i < rand.Intn(20); i++ {
		p.Transactions = append(p.Transactions, RandomValChangeTransaction(&p.PensionID))
	}
	p.Company = *primitives.RandomName()
	p.Value = 0
	p.AuthKey = *primitives.RandomPublicKey()
	p.UniqueHash = *primitives.RandomHash()
	p.FixPids()
	p.Active = true
	return p
}

func (t *Pension) UnmarshalBinary(data []byte) error {
	_, err := t.UnmarshalBinaryData(data)
	return err
}

func (p *Pension) UnmarshalBinaryData(data []byte) (newData []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("[Pension] A panic has occurred while unmarshaling: %s", r)
			return
		}
	}()

	newData = data

	newData, err = p.PensionID.UnmarshalBinaryData(newData)
	if err != nil {
		return data, err
	}

	u, err := primitives.BytesToUint32(newData[:4])
	if err != nil {
		return data, err
	}
	newData = newData[4:]

	p.Transactions = make([]*Transaction, u)
	var i uint32
	for i = 0; i < u; i++ {
		x := new(Transaction)
		newData, err = x.UnmarshalBinaryData(newData)
		if err != nil {
			return data, err
		}
		p.Transactions[i] = x
	}

	newData, err = p.Company.UnmarshalBinaryData(newData)
	if err != nil {
		return data, err
	}

	newData, err = p.UniqueHash.UnmarshalBinaryData(newData)
	if err != nil {
		return data, err
	}

	newData, err = p.Docs.UnmarshalBinaryData(newData)
	if err != nil {
		return data, err
	}

	u, err = primitives.BytesToUint32(newData[:4])
	if err != nil {
		return data, err
	}
	p.Value = int(u)
	newData = newData[4:]

	newData, err = p.AuthKey.UnmarshalBinaryData(newData)
	if err != nil {
		return data, err
	}

	if newData[0] == 0xFF {
		p.Active = true
	} else {
		p.Active = false
	}
	newData = newData[1:]
	return
}

func (p *Pension) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)

	data, err := p.PensionID.MarshalBinary()
	if err != nil {
		return nil, err
	}
	buf.Write(data)

	data = primitives.Uint32ToBytes(uint32(len(p.Transactions)))
	buf.Write(data)

	for _, t := range p.Transactions {
		data, err := t.MarshalBinary()
		if err != nil {
			return nil, err
		}
		buf.Write(data)
	}

	data, err = p.Company.MarshalBinary()
	if err != nil {
		return nil, err
	}
	buf.Write(data)

	data, err = p.UniqueHash.MarshalBinary()
	if err != nil {
		return nil, err
	}
	buf.Write(data)

	data, err = p.Docs.MarshalBinary()
	if err != nil {
		return nil, err
	}
	buf.Write(data)

	data = primitives.Uint32ToBytes(uint32(p.Value))
	buf.Write(data)

	data, err = p.AuthKey.MarshalBinary()
	if err != nil {
		return nil, err
	}
	buf.Write(data)

	switch p.Active {
	case true:
		buf.Write([]byte{0xFF})
	case false:
		buf.Write([]byte{0x00})
	}

	return buf.Next(buf.Len()), nil
}

func (p *Pension) LastInteraction() string {
	if len(p.Transactions) >= 1 {
		sort.Sort(TransList(p.Transactions))
		return p.Transactions[0].Timestamp.Format(layout)
	}
	return "--"
}

func (p *Pension) FixPids() {
	for i := 0; i < len(p.Transactions); i++ {
		p.Transactions[i].PensionID = p.PensionID
	}
}

func (a *Pension) IsSameAs(b *Pension) bool {
	if !a.PensionID.IsSameAs(&b.PensionID) {
		log.Printf("Not PID Same: Found %s, expect %s\n", a.PensionID.String(), b.PensionID.String())
		return false
	}

	if len(a.Transactions) != len(b.Transactions) {
		log.Printf("Not Same Len: Found %d, expect %d\n", len(a.Transactions), len(b.Transactions))
		return false
	}

	for i := range a.Transactions {
		if !a.Transactions[i].IsSameAs(b.Transactions[i]) {
			return false
		}
	}

	if !a.Company.IsSameAs(&b.Company) {
		log.Printf("Not Same: Found %s, expect %s\n", a.Company.String(), b.Company.String())
		return false
	}

	if a.Value != b.Value {
		log.Printf("Not Same: Found %d, expect %d\n", a.Value, b.Value)
		return false
	}

	if !a.AuthKey.IsSameAs(&b.AuthKey) {
		log.Printf("Not Same: Found %s, expect %s\n", a.AuthKey.String(), b.AuthKey.String())
		return false
	}

	if !a.UniqueHash.IsSameAs(&b.UniqueHash) {
		return false
	}
	return true
}
