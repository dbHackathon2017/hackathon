package read

import (
	//"fmt"
	//"time"

	"github.com/FactomProject/factom"
	"github.com/dbHackathon2017/hackathon/common"
	"github.com/dbHackathon2017/hackathon/common/constants"
	"github.com/dbHackathon2017/hackathon/common/primitives"
)

var _ = constants.FAC_LIQUID_SEND

func GetPension(id primitives.Hash) (*common.Pension, error) {
	ents, err := factom.GetAllChainEntries(id.String())
	var _ = ents
	if err != nil {
		return nil, err
	}
	return common.RandomPension(), nil
}
