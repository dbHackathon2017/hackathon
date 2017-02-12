package main

import (
	"encoding/json"
	"fmt"
	//"io/ioutil"
	"net/http"
	"strings"

	"github.com/dbHackathon2017/hackathon/common"
	"github.com/dbHackathon2017/hackathon/common/primitives"
	"github.com/dbHackathon2017/hackathon/company"
	"github.com/dbHackathon2017/hackathon/factom-read"
)

func valToString(val int) string {
	//return fmt.Sprintf("€%d", val)
	currency := "€"
	valStr := "0.00"
	pre := ""
	if val < 0 {
		val = -1 * val
		pre = "-"
	}
	if val < 10 {
		valStr = fmt.Sprintf("0.0%d", val)
	} else if val < 100 {
		valStr = fmt.Sprintf("0.%d", val)
	} else {
		tmp := fmt.Sprintf("%d", val)
		valStr = fmt.Sprintf("%s.%s", tmp[:len(tmp)-2], tmp[len(tmp)-2:])
	}
	return currency + pre + valStr
}

// Get all pensions
type ShortPensionsHolder struct {
	Holder []ShortPensions `json:"pensions"`
}

type ShortPensions struct {
	Acct              string `json:"acct"`
	Firstname         string `json:"firstname"`
	PenID             string `json:"id"`
	Lastint           string `json:"lastint"`
	Lastname          string `json:"lastname"`
	Active            bool   `json:"active"`
	TotalTransactions string `json:"totaltransactions"`
}

func handleAllPensionsUser(w http.ResponseWriter, r *http.Request) error {
	return handleAllPensions(w, r, true)
}

func handleAllPensionsCompany(w http.ResponseWriter, r *http.Request) error {
	return handleAllPensions(w, r, false)
}

func handleAllPensions(w http.ResponseWriter, r *http.Request, user bool) error {
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

		if fpen != nil {
			sp.Lastint = fpen.LastInteraction()
			sp.Active = fpen.Active
			sp.TotalTransactions = fmt.Sprintf("%d", len(fpen.Transactions))
		} else {
			sp.Lastint = "Unknown"
			sp.TotalTransactions = "..."
			sp.Active = true
		}

		if user {
			sp.Active = pens[i].Bucket
		}
		sPens[i] = sp
	}

	container := new(ShortPensionsHolder)

	for i, j := 0, len(sPens)-1; i < j; i, j = i+1, j-1 {
		sPens[i], sPens[j] = sPens[j], sPens[i]
	}
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
	Active    bool   `json:"active"`
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
	header.Active = factomPen.Active
	header.Company = "NestEgg"
	holder.Header = *header

	penStruct := new(LongPension)
	penStruct.Penid = penID.String()
	penStruct.Authority = "NestEgg"
	penStruct.Value = valToString(factomPen.Value)
	if !factomPen.Active {
		penStruct.Value = "€0.00"
	}

	transStruct := make([]Transaction, len(factomPen.Transactions))
	for i, t := range factomPen.Transactions {
		sing := new(Transaction)
		sing.TransactionID = t.TransactionID.String()
		sing.Actor = t.Person.String()
		sing.Factomtype = t.GetFactomTypeString()
		sing.Usertype = t.GetUserTypeString()
		sing.Pension = t.PensionID.String()
		sing.ToPension = t.ToPensionID.String()
		sing.Valchange = valToString(t.ValueChange)
		sing.Timestamp = t.GetTimeStampFormatted()
		sing.Bctimestamp = t.GetTimeStampFormatted()

		docStruct := make([]Document, len(t.Docs.GetFiles()))
		for i, d := range t.Docs.GetFiles() {
			singDoc := new(Document)
			singDoc.Source = d.Source
			singDoc.Location = d.Location
			singDoc.Path = d.GetFullPath()
			singDoc.Timestamp = d.GetTimeStampFormatted()
			singDoc.Hash = d.DocHash.String()

			docStruct[i] = *singDoc
		}

		sing.Docs = docStruct
		transStruct[i] = *sing
	}

	for i, j := 0, len(transStruct)-1; i < j; i, j = i+1, j-1 {
		transStruct[i], transStruct[j] = transStruct[j], transStruct[i]
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
	docIDStr := pr.Params

	findDoc := strings.Contains(pr.Params, "/")
	if findDoc {
		arr := strings.Split(pr.Params, "/")
		if len(arr) > 1 {
			transIDStr = arr[0]
			docIDStr = arr[1]
		} else if len(arr) == 1 {
			transIDStr = arr[0]
		}
	}

	transID, err := primitives.HexToHash(transIDStr)
	if err != nil {
		return err
	}

	docID, err := primitives.HexToHash(docIDStr)
	if err != nil {
		if findDoc {
			return err
		}
	}

	t, err := read.GetTransactionFromTxID(*transID)
	if err != nil {
		return err
	}

	// Only return doc
	if findDoc {
		for _, d := range t.Docs.GetFiles() {
			if d.DocHash.IsSameAs(docID) {
				retDoc := new(Document)
				retDoc.Hash = d.DocHash.String()
				retDoc.Location = d.Location
				retDoc.Source = d.Source
				retDoc.Timestamp = d.GetTimeStampFormatted()
				retDoc.Path = d.GetFullPath()
				w.Write(jsonResp(retDoc))
				return nil
			}
		}
		return fmt.Errorf("Document not found for that transaction")
	}

	sing := new(Transaction)
	sing.TransactionID = t.TransactionID.String()
	sing.Actor = t.Person.String()
	sing.Factomtype = t.GetFactomTypeString()
	sing.Usertype = t.GetUserTypeString()
	sing.Pension = t.PensionID.String()
	sing.ToPension = t.ToPensionID.String()
	sing.Valchange = valToString(t.ValueChange)
	sing.Timestamp = t.GetTimeStampFormatted()
	sing.Bctimestamp = t.GetTimeStampFormatted()

	docStruct := make([]Document, len(t.Docs.GetFiles()))
	for i, d := range t.Docs.GetFiles() {
		singDoc := new(Document)
		singDoc.Source = d.Source
		singDoc.Location = d.Location
		singDoc.Path = d.GetFullPath()
		singDoc.Timestamp = d.GetTimeStampFormatted()
		singDoc.Hash = d.DocHash.String()

		docStruct[i] = *singDoc
	}

	for i, j := 0, len(docStruct)-1; i < j; i, j = i+1, j-1 {
		docStruct[i], docStruct[j] = docStruct[j], docStruct[i]
	}

	sing.Docs = docStruct

	w.Write(jsonResp(sing))
	return nil
}
