package types

type PendingTransaction struct {
	TxId string `bson:"txid"`
}

// BtcRawTransaction models the data from the getrawtransaction command.
type BtcRawTransaction struct {
	Hex           string `json:"hex" bson:"hex"`
	Txid          string `json:"txid" bson:"txid"`
	Hash          string `json:"hash,omitempty" bson:"hash"`
	Size          int32  `json:"size,omitempty" bson:"size"`
	Vsize         int32  `json:"vsize,omitempty" bson:"vsize"`
	Version       int32  `json:"version" bson:"version"`
	LockTime      uint32 `json:"locktime" bson:"locktime"`
	Vin           []Vin  `json:"vin" bson:"vin"`
	Vout          []Vout `json:"vout" bson:"vout"`
	BlockHash     string `json:"blockhash,omitempty" bson:"blockhash"`
	Confirmations uint64 `json:"confirmations,omitempty" bson:"confirmations"`
	Time          int64  `json:"time,omitempty" bson:"time"`
	Blocktime     int64  `json:"blocktime,omitempty" bson:"blocktime"`
	BlockHeight   uint64 `json:"blockheight" bson:"blockheight"`
}

// Vin models parts of the tx data.  It is defined separately since
// getrawtransaction, decoderawtransaction, and searchrawtransaction use the
// same structure.
type Vin struct {
	Coinbase  string     `json:"coinbase" bson:"coinbase"`
	Txid      string     `json:"txid" bson:"txid"`
	Vout      uint32     `json:"vout" bson:"vout"`
	ScriptSig *ScriptSig `json:"scriptSig" bson:"scriptSig"`
	Sequence  uint32     `json:"sequence" bson:"sequence"`
	Witness   []string   `json:"txinwitness" bson:"txinwitness"`
}

// Vout models parts of the tx data.  It is defined separately since both
// getrawtransaction and decoderawtransaction use the same structure.
type Vout struct {
	Value        float64            `json:"value" bson:"value"`
	N            uint32             `json:"n" bson:"n"`
	ScriptPubKey ScriptPubKeyResult `json:"scriptPubKey" bson:"scriptPubKey"`
}

// ScriptSig models a signature script.  It is defined separately since it only
// applies to non-coinbase.  Therefore the field in the Vin structure needs
// to be a pointer.
type ScriptSig struct {
	Asm string `json:"asm" bson:"asm"`
	Hex string `json:"hex" bson:"hex"`
}

// ScriptPubKeyResult models the scriptPubKey data of a tx script.  It is
// defined separately since it is used by multiple commands.
type ScriptPubKeyResult struct {
	Asm       string   `json:"asm" bson:"asm"`
	Hex       string   `json:"hex,omitempty" bson:"hex"`
	ReqSigs   int32    `json:"reqSigs,omitempty" bson:"reqSigs"`
	Type      string   `json:"type" bson:"type"`
	Addresses []string `json:"addresses,omitempty" bson:"addresses"`
}

type BtcTransaction struct {
	TxId              string                 `bson:"txid"`
	Amount            string                 `bson:"amount"`
	Fee               string                 `bson:"fee"`
	BlockHash         string                 `bson:"blockhash"`
	BlockIndex        int64                  `bson:"blockindex"`
	BlockTime         int64                  `bson:"blocktime"`
	Time              int64                  `bson:"time"`
	TimeReceived      int64                  `bson:"timereceived"`
	Bip125Replaceable string                 `bson:"bip125_replaceable"` //"yes|no|unknown",  (string) Whether this transaction could be replaced due to BIP125 (replace-by-fee);may be unknown for unconfirmed transactions not in the mempool
	Hex               string                 `bson:"hex"`
	Details           []BtcTransactionDetail `bson:"details"`
}

type BtcTransactionDetail struct {
	Address  string `bson:"address"`
	Category string `bson:"category"`
	Amount   string `bson:"amount"`
	Label    string `bson:"label"`
	Vout     uint64 `bson:"vout"`
}
