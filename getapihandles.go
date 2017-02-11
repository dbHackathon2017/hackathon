package main

import (
	"encoding/json"
	"fmt"
	//"io/ioutil"
	"net/http"

	"github.com/dbHackathon2017/hackathon/common"
	"github.com/dbHackathon2017/hackathon/common/primitives"
	"github.com/dbHackathon2017/hackathon/company"
	"github.com/dbHackathon2017/hackathon/factom-read"
)

// Get all pensions
type ShortPensionsHolder struct {
	Holder []ShortPensions `json:"pensions"`
}

type ShortPensions struct {
	Acct      string `json:"acct"`
	Firstname string `json:"firstname"`
	PenID     string `json:"id"`
	Lastint   string `json:"lastint"`
	Lastname  string `json:"lastname"`
}

func handleAllPensions(w http.ResponseWriter, r *http.Request) error {
	pens := MainCompany.Pensions
	sPens := make([]ShortPensions, len(pens))
	for i, sp := range sPens {
		sp.Acct = pens[i].AccountNumber
		sp.Firstname = pens[i].FirstName
		sp.Lastname = pens[i].LastName
		sp.PenID = pens[i].PensionID.String()

		fpen := new(common.Pension)
		if pp := GetFromPensionCache(pens[i].PensionID.String()); pp != nil {
			fpen = pp
		} else {
			/*fpen, _ = read.GetPensionFromFactom(pens[i].PensionID)
			if fpen != nil {
				AddToPensionCache(fpen.PensionID.String(), *fpen)
			}*/
			fpen = nil
		}
		fmt.Println(fpen)
		if fpen != nil {
			sp.Lastint = fpen.LastInteraction()
		} else {
			sp.Lastint = "Unknown"
		}

		sPens[i] = sp
	}

	container := new(ShortPensionsHolder)
	container.Holder = sPens

	w.Write(jsonResp(container))
	return nil
}

// Get pension
type LongPensionHolder struct {
	Header  PensionHeader `json:"header"`
	Pension LongPension   `json:"pension"`
}

type PensionHeader struct {
	PensionID string `json:"id"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Acct      string `json:"acct"`
	Lastint   string `json:"lastint"`
	Addr      string `json:"addr"`
	Phone     string `json:"phone"`
	Company   string `json:"company"`
	Ssn       string `json:"ssn"`
}

type LongPension struct {
	Penid        string        `json:"penid"`
	Authority    string        `json:"authority"`
	Value        string        `json:"value"`
	Transactions []Transaction `json:"transactions"`
}

type Transaction struct {
	TransactionID string     `json:"txid"`
	Pension       string     `json:"pension"`
	ToPension     string     `json:"toPension"`
	Valchange     string     `json:"valchange"`
	Factomtype    string     `json:"factomtype"`
	Usertype      string     `json:"usertype"`
	Timestamp     string     `json:"timestamp"`
	Bctimestamp   string     `json:"bctimestamp"`
	Actor         string     `json:"actor"`
	Docs          []Document `json:"docs"`
}

type Document struct {
	Path      string `json:"path"`
	Hash      string `json:"hash"`
	Timestamp string `json:"timestamp"`
	Source    string `json:"source"`
	Location  string `json:"location"`
}

func handlePension(w http.ResponseWriter, r *http.Request, data []byte) error {
	// penIDStr := r.FormValue("content")
	type POSTPenRequest struct {
		Request string `json:"request"`
		Params  string `json:"params,omitempty"`
	}

	pr := new(POSTPenRequest)

	/*data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("Exit 1")
		return err
	}*/

	err := json.Unmarshal(data, pr)
	if err != nil {
		return err
	}

	penIDStr := pr.Params

	penID, err := primitives.HexToHash(penIDStr)
	if err != nil {
		return err
	}

	factomPen := GetFromPensionCache(penID.String())
	if factomPen == nil {
		factomPen, err = read.GetPensionFromFactom(*penID)
		if err != nil {
			return err
		}
	}

	metaPen := new(company.PensionAndMetadata)
	for i := range MainCompany.Pensions {
		if MainCompany.Pensions[i].PensionID.IsSameAs(penID) {
			metaPen = MainCompany.Pensions[i]
		}
	}

	if metaPen == nil {
		return fmt.Errorf("You don't own pension %s", penID.String())
	}

	holder := new(LongPensionHolder)

	header := new(PensionHeader)
	header.PensionID = penID.String()
	header.Firstname = metaPen.FirstName
	header.Lastname = metaPen.LastName
	header.Acct = metaPen.AccountNumber
	header.Lastint = factomPen.LastInteraction()
	header.Addr = metaPen.Address
	header.Phone = metaPen.PhoneNumber
	header.Ssn = metaPen.SSN
	holder.Header = *header

	penStruct := new(LongPension)
	penStruct.Penid = penID.String()
	penStruct.Authority = metaPen.CompanyName
	penStruct.Value = fmt.Sprintf("%d", factomPen.Value)

	transStruct := make([]Transaction, len(factomPen.Transactions))
	for i, t := range factomPen.Transactions {
		sing := new(Transaction)
		sing.TransactionID = t.TransactionID.String()
		sing.Actor = t.Person.String()
		sing.Factomtype = t.GetFactomTypeString()
		sing.Usertype = t.GetUserTypeString()
		sing.Pension = t.PensionID.String()
		sing.ToPension = t.ToPensionID.String()
		sing.Valchange = fmt.Sprintf("%d", t.ValueChange)
		sing.Timestamp = t.GetTimeStampFormatted()
		sing.Bctimestamp = t.GetTimeStampFormatted()

		docStruct := make([]Document, len(t.Docs.GetFiles()))
		for i, d := range t.Docs.GetFiles() {
			singDoc := new(Document)
			singDoc.Source = d.Source
			singDoc.Location = d.Location
			singDoc.Path = d.Name
			singDoc.Timestamp = d.GetTimeStampFormatted()
			singDoc.Hash = d.DocHash.String()

			docStruct[i] = *singDoc
		}

		sing.Docs = docStruct
		transStruct[i] = *sing
	}

	penStruct.Transactions = transStruct
	holder.Pension = *penStruct

	w.Write(jsonResp(holder))
	return nil
}

func handleTransaction(w http.ResponseWriter, r *http.Request, data []byte) error {
	// penIDStr := r.FormValue("content")
	type POSTTranRequest struct {
		Request string `json:"request"`
		Params  string `json:"params,omitempty"`
	}

	pr := new(POSTTranRequest)

	/*data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("Exit 1")
		return err
	}*/

	err := json.Unmarshal(data, pr)
	if err != nil {
		return err
	}

	transIDStr := pr.Params

	transID, err := primitives.HexToHash(transIDStr)
	if err != nil {
		return err
	}

	t, err := read.GetTransactionFromTxID(*transID)
	if err != nil {
		return err
	}

	sing := new(Transaction)
	sing.TransactionID = t.TransactionID.String()
	sing.Actor = t.Person.String()
	sing.Factomtype = t.GetFactomTypeString()
	sing.Usertype = t.GetUserTypeString()
	sing.Pension = t.PensionID.String()
	sing.ToPension = t.ToPensionID.String()
	sing.Valchange = fmt.Sprintf("%d", t.ValueChange)
	sing.Timestamp = t.GetTimeStampFormatted()
	sing.Bctimestamp = t.GetTimeStampFormatted()

	docStruct := make([]Document, len(t.Docs.GetFiles()))
	for i, d := range t.Docs.GetFiles() {
		singDoc := new(Document)
		singDoc.Source = d.Source
		singDoc.Location = d.Location
		singDoc.Path = d.Name
		singDoc.Timestamp = d.GetTimeStampFormatted()
		singDoc.Hash = d.DocHash.String()

		docStruct[i] = *singDoc
	}

	sing.Docs = docStruct

	w.Write(jsonResp(sing))
	return nil
}
