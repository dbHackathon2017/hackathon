package write

import (
	"bytes"
	"crypto/sha256"

	"github.com/FactomProject/factom"
	"github.com/dbHackathon2017/hackathon/common/constants"
	"github.com/dbHackathon2017/hackathon/common/primitives"
)

func FindValidNonce(e *factom.Entry) []byte {
	upToNonce := upToNonce(e.ExtIDs)
	var count uint64
	count = 0
	exit := false
	for exit == false {
		count++
		exit = checkNonce(upToNonce, count)

	}

	data, _ := primitives.Uint64ToBytes(count)
	return data
}

func checkNonce(upToNonce []byte, nonceInt uint64) bool {
	buf := new(bytes.Buffer)
	buf.Write(upToNonce)

	nonce, _ := primitives.Uint64ToBytes(nonceInt)
	//nonce := []byte(strconv.Itoa(nonceInt))
	result := sha256.Sum256(nonce)
	buf.Write(result[:])

	result = sha256.Sum256(buf.Bytes())

	chainFront := result[:constants.CHAIN_PREFIX_LENGTH_CHECK]

	if bytes.Compare(chainFront[:constants.CHAIN_PREFIX_LENGTH_CHECK],
		constants.CHAIN_PREFIX[:constants.CHAIN_PREFIX_LENGTH_CHECK]) == 0 {
		return true
	}
	return false
}

// upToNonce is exclusive
func upToNonce(extIDs [][]byte) []byte {
	buf := new(bytes.Buffer)
	for _, e := range extIDs {
		result := sha256.Sum256(e)
		buf.Write(result[:])
	}

	return buf.Next(buf.Len())
}

func GetECAddress() *factom.ECAddress {
	ec, _ := factom.GetECAddress("Es3cpDrGJRZpJBqZ3PwdohDpmMcXqmr8PuN2yyzBdB2rZ2McEtu1")
	return ec
}
