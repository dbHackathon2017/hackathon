package company

import (
	"time"

	"github.com/dbHackathon2017/hackathon/common"
	"github.com/dbHackathon2017/hackathon/common/primitives"
	"github.com/dbHackathon2017/hackathon/common/primitives/random"
	"github.com/dbHackathon2017/hackathon/factom-read"
	"github.com/dbHackathon2017/hackathon/factom-write"
)

// Controls this whole project. Wraps the factom stuff and stores personal info.
// This is what keeps a list of pension IDs.

type FakeCompany struct {
	CompanyName primitives.PersonName
	SigningKey  primitives.PrivateKey

	Pensions []*PensionAndMetadata
}

type PensionAndMetadata struct {
	PensionID  primitives.Hash
	SigningKey primitives.PrivateKey

	// Client
	FirstName   string
	LastName    string
	Address     string
	PhoneNumber string
	CompanyName string
	SSN         string

	AccountNumber string
}

func RandomFakeCompay() *FakeCompany {
	fc := new(FakeCompany)
	fc.CompanyName = *primitives.RandomName()
	sec, _ := primitives.RandomPrivateKey()
	fc.SigningKey = *sec
	fc.Pensions = make([]*PensionAndMetadata, 0)

	return fc
}

func (fc *FakeCompany) CreatePension() (primitives.Hash, error) {
	pm := new(PensionAndMetadata)
	pm.SigningKey = fc.SigningKey

	p := new(common.Pension)

	p.AuthKey = fc.SigningKey.Public
	p.Value = 0
	p.Company = fc.CompanyName
	p.UniqueHash = *primitives.RandomHash()

	ec := write.GetECAddress()
	_, err := write.SubmitPensionToFactom(p, ec)
	if err != nil {
		return *primitives.NewZeroHash(), err
	}

	pm.FirstName = random.RandStringOfSize(20)
	pm.LastName = random.RandStringOfSize(20)
	pm.Address = random.RandStringOfSize(20)
	pm.PhoneNumber = random.RandStringOfSize(14)
	pm.CompanyName = fc.CompanyName.String()
	pm.SSN = random.RandStringOfSize(8)
	pm.AccountNumber = random.RandStringOfSize(8)

	pm.PensionID = p.PensionID
	fc.Pensions = append(fc.Pensions, pm)

	return p.PensionID, nil
}

// AddValue returns transaction hash and error
func (p *PensionAndMetadata) AddValue(valueChange int, person primitives.PersonName, docs primitives.FileList) (*primitives.Hash, error) {
	trans := new(common.Transaction)
	trans.PensionID = p.PensionID
	trans.ToPensionID = p.PensionID
	trans.Timestamp = time.Now()
	trans.Person = person
	trans.Docs = docs
	trans.ValueChange = valueChange

	ec := write.GetECAddress()
	return write.SubmitValueChangeTransactionToPension(trans, ec, p.SigningKey)
}

// Will move all value from a into b
func (a *PensionAndMetadata) MoveChainTo(b *PensionAndMetadata, person primitives.PersonName, docs primitives.FileList) error {
	aPen, err := read.GetPensionFromFactom(a.PensionID)
	if err != nil {
		return err
	}
	ec := write.GetECAddress()

	// Put the send transaction on A
	send := new(common.Transaction)
	send.PensionID = a.PensionID
	send.ToPensionID = a.PensionID
	send.Timestamp = time.Now()
	send.Person = person
	send.Docs = docs
	send.ValueChange = aPen.Value
	_, err = write.SubmitChainMoveTransactionToPension(send, ec, a.SigningKey)
	if err != nil {
		return err
	}

	// Now put the request transaction in the recieving
	req := new(common.Transaction)
	req.PensionID = b.PensionID
	req.ToPensionID = b.PensionID
	req.Timestamp = time.Now()
	req.Person = person
	req.Docs = docs
	req.ValueChange = aPen.Value
	_, err = write.SubmitRequestMoveTransactionToPension(req, ec, b.SigningKey)
	if err != nil {
		return err
	}

	return nil
}
