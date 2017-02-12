package company

import (
	"bytes"
	"fmt"
	"log"
	"sync"

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

	DBPath string
	DB     *BoltDB
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

func RandomPenstionAndMetaData() *PensionAndMetadata {
	p := new(PensionAndMetadata)
	p.PensionID = *primitives.RandomHash()
	sec, _ := primitives.RandomPrivateKey()
	p.SigningKey = *sec

	p.FirstName = random.RandStringOfSize(random.RandomIntBetween(0, 100))
	p.LastName = random.RandStringOfSize(random.RandomIntBetween(0, 100))
	p.Address = random.RandStringOfSize(random.RandomIntBetween(0, 100))
	p.PhoneNumber = random.RandStringOfSize(random.RandomIntBetween(0, 100))
	p.CompanyName = random.RandStringOfSize(random.RandomIntBetween(0, 100))
	p.SSN = random.RandStringOfSize(random.RandomIntBetween(0, 100))
	p.AccountNumber = random.RandStringOfSize(random.RandomIntBetween(0, 100))
	return p
}

var PENSION_BUCKET []byte = []byte("Pensions") // Pension metadaa
var PRIV_KEY []byte = []byte("Secret")
var COMPANY_NAME []byte = []byte("CompName")
var PENSION_CACHE []byte = []byte("CompName") // Pension data in factom

var lock sync.RWMutex

func (fc *FakeCompany) Save(penCache []common.Pension, full bool) {
	if fc.DB == nil {
		return
	}

	lock.Lock()
	defer lock.Unlock()
	log.Println("Saving to db....")

	for _, p := range fc.Pensions {
		data, err := p.MarshalBinary()
		if err != nil {
			log.Printf("Failed to save pension %s\n", p.PensionID.String())
		}
		err = fc.DB.Put(PENSION_BUCKET, p.PensionID.Bytes(), data)
		if err != nil {
			log.Printf("Failed to save pension %s\n", p.PensionID.String())
		}
	}

	sec, err := fc.SigningKey.MarshalBinary()
	if err != nil {
		log.Printf("Failed to save Secret key %x\n", fc.SigningKey.Secret[:])
	}

	err = fc.DB.Put(PRIV_KEY, PRIV_KEY, sec)
	if err != nil {
		log.Printf("Failed to save Secret key %x\n", fc.SigningKey.Secret[:])
	}

	compData, err := fc.CompanyName.MarshalBinary()
	if err != nil {
		log.Printf("Failed to save Company Name \n")
	}

	err = fc.DB.Put(COMPANY_NAME, COMPANY_NAME, compData)
	if err != nil {
		log.Printf("Failed to save Company Name \n")
	}

	//if full {
	recs := make([]Record, 0)
	for _, p := range penCache {
		data, err := p.MarshalBinary()
		if err != nil {
			log.Printf("Failed to save pension %s\n", p.PensionID.String())
		}

		if len(data) < 100 {
			continue
		}
		rc := new(Record)
		rc.Bucket = PENSION_CACHE
		rc.Key = p.PensionID.Bytes()
		rc.Data = data
		recs = append(recs, *rc)
		/*err = fc.DB.Put(PENSION_CACHE, p.PensionID.Bytes(), data)
		if err != nil {
			log.Printf("Failed to save pension %s\n", p.PensionID.String())
		}*/
	}
	err = fc.DB.PutInBatch(recs)
	if err != nil {
		log.Printf("Failed to save factom-pensions \n")
	}
	//}

}

func (fc *FakeCompany) LoadFromDB() {
	fmt.Print("Loading pensions from DB...")
	if fc.DB == nil {
		fc.DB = NewBoltDB("company_cache.db")
	}

	lock.RLock()
	defer lock.RUnlock()

	/*keys, err := fc.DB.ListAllKeys(PENSION_BUCKET)
	if err != nil {
		log.Printf("Failed to load pensions\n")
	}*/

	fc.Pensions = make([]*PensionAndMetadata, 0)

	fmt.Printf("Loading Pensions...\n")
	dataSet, _, err := fc.DB.GetAll(PENSION_BUCKET)
	for _, data := range dataSet {
		t := new(PensionAndMetadata)
		err = t.UnmarshalBinary(data)
		if err != nil {
			log.Printf("Failed to unmarshal pensionmetadata")
			continue
		}

		fc.Pensions = append(fc.Pensions, t)
	}

	data, err := fc.DB.Get(PRIV_KEY, PRIV_KEY)
	if data != nil && len(data) > 0 {
		if err != nil {
			log.Printf("Failed to load Secret key\n")
		} else {
			sec := new(primitives.PrivateKey)
			err := sec.UnmarshalBinary(data)
			if err != nil {
				log.Printf("Failed to unmarshal Secret key\n")
			} else {
				fc.SigningKey = *sec
			}
		}
	}

	data, err = fc.DB.Get(COMPANY_NAME, COMPANY_NAME)
	if data != nil && len(data) > 0 {
		if err != nil {
			log.Printf("Failed to load Secret key\n")
		} else {
			comp := new(primitives.PersonName)
			err := comp.UnmarshalBinary(data)
			if err != nil {
				log.Printf("Failed to unmarshal Secret key\n")
			} else {
				fc.CompanyName = *comp
			}
		}
	}

	return
}

func (fc *FakeCompany) LoadPenCacheFromDB() []common.Pension {
	list := make([]common.Pension, 0)
	fmt.Printf("Loading factom-pensions from cache...")
	dataSet, _, err := fc.DB.GetAll(PENSION_CACHE)
	if err != nil {
		log.Println("Encountered error", err.Error())
		return nil
	}
	for _, data := range dataSet {
		p := new(common.Pension)
		err = p.UnmarshalBinary(data)
		if err != nil {
			//log.Printf("Failed to unmarshal pension %s\n", err.Error())
			continue
		}
		list = append(list, *p)
	}
	return list
}

func NewCompany(path string) *FakeCompany {
	f := new(FakeCompany)
	f.Pensions = make([]*PensionAndMetadata, 0)
	if path != "none" {
		f.DB = NewBoltDB(path)
	}

	return f
}

func RandomFakeCompay() *FakeCompany {
	fc := NewCompany("none")
	fc.CompanyName = *primitives.RandomName()
	sec, _ := primitives.RandomPrivateKey()
	fc.SigningKey = *sec
	fc.Pensions = make([]*PensionAndMetadata, 0)

	return fc
}

// CreatePension secures a pension into factom.
func (fc *FakeCompany) CreatePension(fn, ln, addr, pn, ssn, acct string, docs primitives.FileList) (primitives.Hash, error) {
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

	// Metadata
	pm.FirstName = fn
	pm.LastName = ln
	pm.Address = addr
	pm.PhoneNumber = pn
	pm.CompanyName = fc.CompanyName.String()
	pm.SSN = ssn
	pm.AccountNumber = acct

	pm.PensionID = p.PensionID
	fc.Pensions = append(fc.Pensions, pm)

	return p.PensionID, nil
}

// Return pension from our local list (NOT from factom)
func (fc *FakeCompany) GetPensionByID(id string) *PensionAndMetadata {
	for _, p := range fc.Pensions {
		if p.PensionID.String() == id {
			return p
		}
	}
	return nil
}

func (fc *FakeCompany) CreateRandomPension() (primitives.Hash, error) {
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
// Addvalue will add a token amount to pension
// Goes to Factom
func (p *PensionAndMetadata) AddValue(valueChange int, person primitives.PersonName, docs primitives.FileList, randTime bool) (*primitives.Hash, error) {
	trans := new(common.Transaction)
	trans.PensionID = p.PensionID
	trans.ToPensionID = p.PensionID
	//if randTime {
	trans.Timestamp = random.RandomTimestamp()
	//} else {
	//	trans.Timestamp = time.Now()
	//}
	trans.Person = person
	trans.Docs = docs
	trans.ValueChange = valueChange

	ec := write.GetECAddress()
	return write.SubmitValueChangeTransactionToPension(trans, ec, p.SigningKey)
}

// Will move all value from a into b
// Adds to factom
func (a *PensionAndMetadata) MoveChainTo(b *PensionAndMetadata, person primitives.PersonName, docs primitives.FileList) error {
	aPen, err := read.GetPensionFromFactom(a.PensionID)
	if err != nil {
		return err
	}
	ec := write.GetECAddress()

	// Put the send transaction on A
	send := new(common.Transaction)
	send.PensionID = a.PensionID
	send.ToPensionID = b.PensionID
	send.Timestamp = random.RandomTimestamp()
	send.Person = person
	send.Docs = docs
	send.ValueChange = aPen.Value
	e1, err := write.SubmitChainMoveTransactionToPension(send, ec, a.SigningKey)
	if err != nil {
		return err
	}

	// Now put the request transaction in the recieving
	req := new(common.Transaction)
	req.PensionID = b.PensionID
	req.ToPensionID = a.PensionID
	req.Timestamp = send.Timestamp
	req.Person = person
	req.Docs = docs
	req.ValueChange = aPen.Value
	e, err := write.SubmitRequestMoveTransactionToPension(req, ec, b.SigningKey)
	if err != nil {
		return err
	}
	fmt.Printf(e.String() + "\n" + e1.String() + "\n")

	return nil
}

func (a *PensionAndMetadata) IsSameAs(b *PensionAndMetadata) bool {
	if !a.PensionID.IsSameAs(&b.PensionID) {
		return false
	}

	if !a.SigningKey.IsSameAs(&b.SigningKey) {
		return false
	}

	if a.FirstName != b.FirstName {
		return false
	}

	if a.LastName != b.LastName {
		return false
	}

	if a.Address != b.Address {
		return false
	}

	if a.PhoneNumber != b.PhoneNumber {
		return false
	}

	if a.CompanyName != b.CompanyName {
		return false
	}

	if a.SSN != b.SSN {
		return false
	}

	if a.AccountNumber != b.AccountNumber {
		return false
	}

	return true
}

const MAX_LENGTH = 200

func (pm *PensionAndMetadata) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)

	data, err := pm.PensionID.MarshalBinary()
	if err != nil {
		return nil, err
	}
	buf.Write(data)

	data, err = pm.SigningKey.MarshalBinary()
	if err != nil {
		return nil, err
	}
	buf.Write(data)

	data, err = primitives.MarshalStringToBytes(pm.FirstName, MAX_LENGTH)
	if err != nil {
		return nil, err
	}
	buf.Write(data)

	data, err = primitives.MarshalStringToBytes(pm.LastName, MAX_LENGTH)
	if err != nil {
		return nil, err
	}
	buf.Write(data)

	data, err = primitives.MarshalStringToBytes(pm.Address, MAX_LENGTH)
	if err != nil {
		return nil, err
	}
	buf.Write(data)

	data, err = primitives.MarshalStringToBytes(pm.PhoneNumber, MAX_LENGTH)
	if err != nil {
		return nil, err
	}
	buf.Write(data)

	data, err = primitives.MarshalStringToBytes(pm.CompanyName, MAX_LENGTH)
	if err != nil {
		return nil, err
	}
	buf.Write(data)

	data, err = primitives.MarshalStringToBytes(pm.SSN, MAX_LENGTH)
	if err != nil {
		return nil, err
	}
	buf.Write(data)

	data, err = primitives.MarshalStringToBytes(pm.AccountNumber, MAX_LENGTH)
	if err != nil {
		return nil, err
	}
	buf.Write(data)

	return buf.Next(buf.Len()), nil
}

func (pm *PensionAndMetadata) UnmarshalBinary(data []byte) (err error) {
	_, err = pm.UnmarshalBinaryData(data)
	return
}

func (pm *PensionAndMetadata) UnmarshalBinaryData(data []byte) (newData []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("[PensionMetaData] A panic has occurred while unmarshaling: %s", r)
			return
		}
	}()

	newData = data

	newData, err = pm.PensionID.UnmarshalBinaryData(newData)
	if err != nil {
		return data, err
	}

	newData, err = pm.SigningKey.UnmarshalBinaryData(newData)
	if err != nil {
		return data, err
	}

	var str string
	str, newData, err = primitives.UnmarshalStringFromBytesData(newData, MAX_LENGTH)
	if err != nil {
		return data, err
	}
	pm.FirstName = str

	str, newData, err = primitives.UnmarshalStringFromBytesData(newData, MAX_LENGTH)
	if err != nil {
		return data, err
	}
	pm.LastName = str

	str, newData, err = primitives.UnmarshalStringFromBytesData(newData, MAX_LENGTH)
	if err != nil {
		return data, err
	}
	pm.Address = str

	str, newData, err = primitives.UnmarshalStringFromBytesData(newData, MAX_LENGTH)
	if err != nil {
		return data, err
	}
	pm.PhoneNumber = str

	str, newData, err = primitives.UnmarshalStringFromBytesData(newData, MAX_LENGTH)
	if err != nil {
		return data, err
	}
	pm.CompanyName = str

	str, newData, err = primitives.UnmarshalStringFromBytesData(newData, MAX_LENGTH)
	if err != nil {
		return data, err
	}
	pm.SSN = str

	str, newData, err = primitives.UnmarshalStringFromBytesData(newData, MAX_LENGTH)
	if err != nil {
		return data, err
	}
	pm.AccountNumber = str

	return

}
