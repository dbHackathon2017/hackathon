package common

import (
	"github.com/dbHackathon2017/hackathon/common/primitives"
	"transaction"
	"math/rand"
)

// Pension structs contain all information about a pension relative to
// factom. This means it does not contain personal information (name/address/etc)
// instead all pensions are tied to a sha256 hash. This hash is it's location in
// factom.
type Pension struct {
	// ChainID in Factom. This is how we find a pension
	PensionID primitives.Hash
	var Transactions []Transaction
	Company primitives.PersonName // Company name, use same primitive

	// The current amount of tokens in the pension.
	Value   int
	AuthKey primitives.PublicKey
}

func RandomPension() *Pension {
	p := new(Pension)
	p.PensionID = *primitives.RandomHash()
	for i := 0; i < rand.Intn(20); i++ {
		Transactions = append(s, RandomValChangeTransaction(p.PensionID))
	}
	p.Company = *primitives.RandomName()
	p.Value = 0
	p.AuthKey = *primitives.RandomPublicKey()
	return p
}
