package company

import (
	"github.com/dbHackathon2017/hackathon/common/primitives"
)

// Controls this whole project. Wraps the factom stuff and stores personal info.
// This is what keeps a list of pension IDs.

type FakeCompany struct {
	CompanyName primitives.PersonName
	SigningKey  primitives.PrivateKey

	Pensions []PensionAndMetadata
}

type PensionAndMetadata struct {
	PensionID primitives.Hash
	// Other personal info and metadata
}
