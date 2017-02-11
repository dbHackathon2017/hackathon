package main

import (
	"net/http"

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
type AutoGenerated struct {
	Header struct {
		ID        string `json:"id"`
		Firstname string `json:"firstname"`
		Lastname  string `json:"lastname"`
		Acct      string `json:"acct"`
		Lastint   string `json:"lastint"`
		Addr      string `json:"addr"`
		Phone     string `json:"phone"`
		Company   string `json:"company"`
		Ssn       string `json:"ssn"`
	} `json:"header"`
	Pensions struct {
		Penid        string `json:"penid"`
		Transid      string `json:"transid"`
		Authority    string `json:"authority"`
		Value        string `json:"value"`
		Transactions []struct {
			Penid       string `json:"penid"`
			Pension     string `json:"pension"`
			ToPension   string `json:"toPension"`
			Valchange   string `json:"valchange"`
			Factomtype  string `json:"factomtype"`
			Usertype    string `json:"usertype"`
			Timestamp   string `json:"timestamp"`
			Bctimestamp string `json:"bctimestamp"`
			Actor       string `json:"actor"`
			Docs        []struct {
				Path      string `json:"path"`
				Hash      string `json:"hash"`
				Timestamp string `json:"timestamp"`
				Source    string `json:"source"`
				Location  string `json:"location"`
			} `json:"docs"`
		} `json:"transactions"`
	} `json:"pensions"`
}
