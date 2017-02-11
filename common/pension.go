package common

import (
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

	// The current amount of tokens in the pension.
	Value   int
	AuthKey primitives.PublicKey
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
	p.FixPids()
	return p
}

func (p *Pension) LastInteraction() string {
	if len(p.Transactions) >= 1 {
		sort.Sort(TransList(p.Transactions))
		return p.Transactions[0].Timestamp.Format(layout)
	}
	return "NA"

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
	return true
}
