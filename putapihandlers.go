package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dbHackathon2017/hackathon/common/primitives"
	"github.com/dbHackathon2017/hackathon/common/primitives/random"
)

type PutPension struct {
	Acctnum   string   `json:"acctnum"`
	Address   string   `json:"address"`
	Company   string   `json:"company"`
	Docs      []PutDoc `json:"docs"`
	Firstname string   `json:"firstname"`
	Lastname  string   `json:"lastname"`
	Phone     string   `json:"phone"`
	Ssn       string   `json:"ssn"`
}

type PutDoc struct {
	Hash      string `json:"hash"`
	Location  string `json:"location"`
	Path      string `json:"path"`
	Source    string `json:"source"`
	Timestamp string `json:"timestamp"`
}

func handleMakePension(w http.ResponseWriter, r *http.Request, data []byte) error {
	// penIDStr := r.FormValue("content")

	type POSTMakePen struct {
		Request string     `json:"request"`
		Params  PutPension `json:"params,omitempty"`
	}

	pr := new(POSTMakePen)

	err := json.Unmarshal(data, pr)
	if err != nil {
		fmt.Println(err)
		return err
	}

	docs := new(primitives.FileList)

	list := docsToFileList(pr.Params.Docs)
	docs = list
	docs.FixFiles()

	penID, err := MainCompany.CreatePension(pr.Params.Firstname,
		pr.Params.Lastname,
		pr.Params.Address,
		pr.Params.Phone,
		pr.Params.Ssn,
		pr.Params.Acctnum,
		*docs)
	if err != nil {
		return err
	}

	w.Write(jsonResp(penID.String()))
	return nil
}

func docsToFileList(pdocs []PutDoc) *primitives.FileList {
	docs := new(primitives.FileList)
	for _, d := range pdocs {
		sing := new(primitives.File)
		sing.SetFileName(d.Path)
		sing.Location = d.Location
		sing.Source = d.Source
		sing.Timestamp = random.RandomTimestamp()
		hash, err := primitives.HexToHash(d.Hash)
		if err != nil {
			sing.DocHash = *primitives.NewZeroHash()
		} else {
			sing.DocHash = *hash
		}

		docs.FileList = append(docs.FileList, *sing)
	}

	return docs
}

type ValChange struct {
	Penid     string   `json:"penid"`
	Person    string   `json:"person"`
	Valchange int      `json:"valchange"`
	Docs      []PutDoc `json:"docs"`
}

func handleAddValue(w http.ResponseWriter, r *http.Request, data []byte) error {
	type POSTAddVal struct {
		Request string    `json:"request"`
		Params  ValChange `json:"params,omitempty"`
	}

	pr := new(POSTAddVal)

	err := json.Unmarshal(data, pr)
	if err != nil {
		return err
	}

	docs := new(primitives.FileList)

	list := docsToFileList(pr.Params.Docs)
	docs = list

	pension := MainCompany.GetPensionByID(pr.Params.Penid)
	if pension == nil {
		return fmt.Errorf("Pension by id %s not found in company\n", pr.Params.Penid)
	}

	pp, err := primitives.NewPersonName(pr.Params.Person)
	if err != nil {
		return err
	}

	txid, err := pension.AddValue(pr.Params.Valchange,
		*pp,
		*docs,
		true)
	if err != nil {
		return err
	}

	w.Write(jsonResp(txid.String()))
	return nil
}

type TransferChain struct {
	FromPenid string   `json:"from-penid"`
	ToPenid   string   `json:"to-penid"`
	Person    string   `json:"person"`
	Docs      []PutDoc `json:"docs"`
}

func handleTransfer(w http.ResponseWriter, r *http.Request, data []byte) error {
	type POSTAddVal struct {
		Request string    `json:"request"`
		Params  ValChange `json:"params,omitempty"`
	}

	pr := new(POSTAddVal)

	err := json.Unmarshal(data, pr)
	if err != nil {
		return err
	}

	docs := new(primitives.FileList)

	list := docsToFileList(pr.Params.Docs)
	docs = list

	from := MainCompany.GetPensionByID(pr.Params.Penid)
	if from == nil {
		return fmt.Errorf("Pension by id %s not found in company\n", pr.Params.Penid)
	}

	to := MainCompany.GetPensionByID(pr.Params.Penid)
	if to == nil {
		return fmt.Errorf("Pension by id %s not found in company\n", pr.Params.Penid)
	}

	pp, err := primitives.NewPersonName(pr.Params.Person)
	if err != nil {
		return err
	}

	err = from.MoveChainTo(to,
		*pp,
		*docs)
	if err != nil {
		return err
	}

	w.Write(jsonResp("Success"))
	return nil
}
