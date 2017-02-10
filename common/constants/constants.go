package constants

import (
	"time"
)

// Transaction Types for user readability
const (
	USER_TRANS_VAL_CHANGE uint32 = iota // Withdraw, deposit
	USER_TRANS_DOC_CHANGE               // Change pension metadata outside factom
	USER_LIQUID_REQUEST                 // Request a pension be merged into this one
	USER_LIQUID_SEND                    // Merge into another chain
	USER_LIQUID_CONFIRMED               // Confirm a pension, works both ways
)

// Transaction Types for Factom
const (
	FAC_TRANS_VAL_CHANGE uint32 = iota
	FAC_LIQUID_REQUEST
	FAC_LIQUID_SEND
	FAC_LIQUID_CONFIRM
)

// Supporting Doc types
const (
	DOC_PDF uint32 = iota
	DOC_IMAGE
	DOC_TXT
)

// Factom Chain info
var (
	CHAIN_PREFIX_LENGTH_CHECK int    = 1
	CHAIN_PREFIX              []byte = []byte{0x9E, 0x00, 0x00}
)

// Factom Settings
const (
	REVEAL_WAIT time.Duration = 5 * time.Second
	REMOTE_HOST string        = "192.210.237.146:8088"
	LOCAL_HOST  string        = "localhost:8088"
)

// Type sizes
const (
	HASH_BYTES_LENGTH int = 32
	MAX_FILENAME_LEN  int = 100
	MAX_NAME_LENGTH   int = 50
)
