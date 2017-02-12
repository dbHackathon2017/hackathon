package main

import (
	"net/http"
)

// Company stats

// How many factom chains
// how many factom entries
// how many pensions
// Last record date

type CompanyStats struct {
	CompanyName       string `json:"company"`
	TotalTransactions int    `json:"totaltransaction"`
	TotalPensions     int    `json:"totalpension"`
	TotalValue        int    `json:"value"`
}

func handleCompanyStats(w http.ResponseWriter, r *http.Request, data []byte) error {
	cs := new(CompanyStats)
	cs.CompanyName = MainCompany.CompanyName.String()
	cs.TotalPensions = len(MainCompany.Pensions)

	for _, p := range MainCompany.Pensions {
		fpen := GetFromPensionCache(p.PensionID.String())
		if fpen != nil {
			cs.TotalTransactions += len(fpen.Transactions)
			cs.TotalValue += valToString(fpen.Value)
		}
	}

	w.Write(jsonResp(cs))
	return nil
}
