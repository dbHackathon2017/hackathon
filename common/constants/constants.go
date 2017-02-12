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

var (
	FILE_NAMES []string = []string{"application.pdf",
		"notice.pdf",
		"information.pdf",
		"applicationaward.pdf",
		"legalenforcement.pdf",
		"depositform.pdf",
		"withdrawform.pdf",
		"pensionfund.pdf",
		"agreement.pdf",
		"ammendment.pdf",
		"companynotice.pdf",
		"form.pdf",
		"initialapp.pdf",
		"taxform.pdf",
		"govntapp.pdf",
		"incomeform.pdf",
		"Indictments_IdentityFraud.pdf",
		"termination.pdf"}

	FILE_HASHES []string = []string{"08ab8ebe9f1ed7a92e53128a2008642e49ee9be6164f9f916248f7d2e4e6cf29",
		"6f9a1108807b762d00b3e619b9100f8d493fb20f922ef76d8d9eaa9e1b542a8b",
		"55a548593fb0e1be5e5f199fd59cd94642e75636f37cad86ae4ed18ac1b336c6",
		"21f407390b49f6d73b9e7a5ca9276abe3000f6e374774e0f7aef927c5f049715",
		"57907ecf93f3038142e631491f7e33374a1624ddafd1c8011aa9a39e4870b84b",
		"ad07ecc9ae6a0541ebf5847565932ff66b16375e87aba7a1f9998de3fd7f1a9d",
		"871c56c2dc460c23769506ca6281f68d5ca617bcc5f43e0739b69bb1ad7c488b",
		"47cb2c59e03828cf9d7d879371c6da922cbd46646ebf1949cae94790c53e98ea",
		"f677d2e52dc9eb52fe9381ad562bd38007e6843cf653fa304c5cdf04e0920bb2",
		"2ebb1e25ff0e5bcb669d682b34868734e1666be47237decbbb6c4c4fefc0c3a1",
		"e0c73552e5fcea436aa2be0dda0397f8859e127181cf96054f4b415b9bbdb4fa",
		"2209fb05220829c7dedf596897e721fd23cef17d00ee3a3375dcec636eccc545",
		"2fc0d8bfa2775b4b693122a6571f236499c9fc98230bdce04164fa816bc30322",
		"21c6589e1be4313ed33daf17c83985cbc8381a37ff19682d475c37bced6bbe5c",
		"297697f35eb9ce39536af4beb84bcae9a801a405e49798e223f98e9895f8fa9b",
		"09fc354fb3811a897f980b09cfa02284f77461e8730fb81066caca8a8ae9983a",
		"7bd8adf584843e588b7f844f5e9d2ce8d1be17b23275da473ddb870c7dc604eb",
		"683b4871f45c46f428e8ea3af5d2cadb41287f15a31c82f2d0715cff659bc62d"}
)
