package main

import (
	"fmt"
	"net/http"

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

		fpen, err := read.GetPensionFromFactom(pens[i].PensionID)
		if err == nil {
			sp.Lastint = fpen.LastInteraction()
		} else {
			sp.Lastint = "Unknown"
		}
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
	Transid      string        `json:"transid"`
	Authority    string        `json:"authority"`
	Value        string        `json:"value"`
	Transactions []Transaction `json:"transactions"`
}

type Transaction struct {
	Penid       string     `json:"penid"`
	Pension     string     `json:"pension"`
	ToPension   string     `json:"toPension"`
	Valchange   string     `json:"valchange"`
	Factomtype  string     `json:"factomtype"`
	Usertype    string     `json:"usertype"`
	Timestamp   string     `json:"timestamp"`
	Bctimestamp string     `json:"bctimestamp"`
	Actor       string     `json:"actor"`
	Docs        []Document `json:"docs"`
}

type Document struct {
	Path      string `json:"path"`
	Hash      string `json:"hash"`
	Timestamp string `json:"timestamp"`
	Source    string `json:"source"`
	Location  string `json:"location"`
}

func handlePension(w http.ResponseWriter, r *http.Request) error {
	penIDStr := r.FormValue("content")
	penID, err := primitives.HexToHash(penIDStr)
	if err != nil {
		return err
	}

	factomPen, err := read.GetPensionFromFactom(*penID)
	if err != nil {
		return err
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
	header.Acct = metaPen.AccountNumber

	penStruct := new(LongPension)
	penStruct.Penid = penID.String()

	var _, _ = factomPen, holder

	return nil

	/*
		pens := MainCompany.Pensions
		sPens := make([]ShortPensions, len(pens))
		for i, sp := range sPens {
			sp.Acct = pens[i].AccountNumber
			sp.Firstname = pens[i].FirstName
			sp.Lastname = pens[i].LastName
			sp.PenID = pens[i].PensionID.String()

			fpen, err := read.GetPensionFromFactom(pens[i].PensionID)
			if err == nil {
				sp.Lastint = fpen.LastInteraction()
			} else {
				sp.Lastint = "Unknown"
			}
		}

		container := new(ShortPensionsHolder)
		container.Holder = sPens

		w.Write(jsonResp(container))
		return nil*/
}
