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
	TotalTransactions int `json:"total-transaction"`
	TotalPensions     int `json:"total-pension"`
	TotalValue        int `json:"value"`
}

func handleCompanyStats(w http.ResponseWriter, r *http.Request, data []byte) error {
	cs := new(CompanyStats)
	cs.TotalPensions = len(MainCompany.Pensions)

	for _, p := range MainCompany.Pensions {
		fpen := GetFromPensionCache(p.PensionID.String())
		if fpen != nil {
			cs.TotalTransactions += len(fpen.Transactions)
			cs.TotalValue += fpen.Value
		}
	}

	w.Write(jsonResp(cs))
	return nil
}
