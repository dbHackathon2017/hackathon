package main

// Feeder feeds random transactions into random pensions

import (
	"log"
	"time"

	"github.com/dbHackathon2017/hackathon/common/primitives"
	"github.com/dbHackathon2017/hackathon/common/primitives/random"
	"github.com/dbHackathon2017/hackathon/company"
)

var minTPS int = 2
var maxTPS int = 15

func startFeeder() {
	for {
		dead := len(MainCompany.Pensions)
		if dead == 0 {
			continue
		}
		for _, p := range MainCompany.Pensions {
			pp := GetFromPensionCache(p.PensionID.String())
			if pp != nil {
				if pp.Active {
					dead--
				}
			}
		}

		amt := random.RandomIntBetween(minTPS*10, maxTPS*10)
		log.Printf("[FEEDER] Adding %d transactions from feeder...\n", amt)
		for i := 0; i < amt; i++ {
			action := random.RandomIntBetween(0, 100)
			i, j := random.RandomIntBetween(0, len(MainCompany.Pensions)), random.RandomIntBetween(0, len(MainCompany.Pensions))
			switch {
			case action >= 85:
				if dead > (len(MainCompany.Pensions) / 5) {
					log.Printf("[FEEDER] Refusing to move a chain, have too many dead. Dead %d\n", dead)
					continue // Stop killing everyone dammit
				}
				log.Println("[FEEDER] Liquidating a chain...")
				dead++
				moveTo(MainCompany.Pensions[i], MainCompany.Pensions[j])
			default:
				if pen := GetFromPensionCache(MainCompany.Pensions[i].PensionID.String()); pen != nil {
					if pen.Active {
						addRandVal(MainCompany.Pensions[i])
						continue
					}
				}

				// If got here, I is not active
				if pen := GetFromPensionCache(MainCompany.Pensions[j].PensionID.String()); pen != nil {
					if pen.Active {
						addRandVal(MainCompany.Pensions[j])
					}
				}

				// If both fail, oh well.
			}
		}

		time.Sleep(10 * time.Second)
	}
}

func addRandVal(p *company.PensionAndMetadata) {
	t := random.RandomIntBetween(0, 100)
	val := 0
	if t < 25 {
		val = random.RandomIntBetween(-5000, 0)
	} else if t > 35 {
		val = random.RandomIntBetween(0, 10000)
	} else {
		val = 0
	}
	p.AddValue(val, "NestEgg", *primitives.RandomFileList(10), true)
}

func moveTo(a *company.PensionAndMetadata, b *company.PensionAndMetadata) {
	a.MoveChainTo(b, "NestEgg", *primitives.RandomFileList(10))
}
