package omni

import "encoding/json"

type Tx struct {
	// The hex-encoded hash of the transaction
	Txid string
	// The transaction fee in bitcoins
	Fee string
	// The Bitcoin address of the sender
	Sendingaddress string
	// A Bitcoin address used as reference (if any)
	Referenceaddress string
	// Whether the transaction involes an address in the wallet
	Ismine bool
	// The transaction version
	Version uint32
	// The transaction type as number
	Type_int uint32
	// The transaction type as string
	Type string
	// The token property id
	Propertyid uint32
	// Can be divisible
	Divisible bool
	// Amount of the transaction
	Amount string
	// Whether the transaction is valid
	Valid bool
	// The hash of the block that contains the transaction
	Blockhash string
	// The timestamp of the block that contains the transaction
	Blocktime uint32
	// The position of the transaction in the block
	Positioninblock uint64
	// The height of the block that contains the transaction
	Block uint64
	// The number of transaction confirmations
	Confirmations uint64
}

type BalanceResult struct {
	Balance, Reserved string
}

type BitcoinTransaction struct {
	Txid              string                 `json:"txid"`
	Amount            json.Number            `json:"amount"`
	Fee               json.Number            `json:"fee"`
	Blockhash         string                 `json:"blockhash"`
	Blockindex        int64                  `json:"blockindex"`
	Blocktime         int64                  `json:"blocktime"`
	Time              int64                  `json:"time"`
	Timereceived      int64                  `json:"timereceived"`
	Bip125Replaceable string                 `json:"bip125-replaceable"` //"yes|no|unknown",  (string) Whether this transaction could be replaced due to BIP125 (replace-by-fee);may be unknown for unconfirmed transactions not in the mempool
	Hex               string                 `json:"hex"`
	Details           []BtcTransactionDetail `json:"details"`
}

type BtcTransactionDetail struct {
	Address  string      `bson:"address"`
	Category string      `bson:"category"`
	Amount   json.Number `bson:"amount"`
	Label    string      `bson:"label"`
	Vout     uint64      `bson:"vout"`
}

//type Unspent struct {
//	Txid          string      `json:"txid"`
//	Vout          int         `json:"vout"`
//	Address       string      `json:"address"`
//	Label         string      `json:"label"`
//	ScriptPubKey  string      `json:"scriptPubKey"`
//	Amount        json.Number `json:"amount"`
//	Confirmations json.Number `json:"confirmations"`
//	RedeemScript  string      `json:"redeemScript"`
//	WitnessScript string      `json:"witnessScript"`
//	Spendable     bool        `json:"spendable"`
//	Solvable      bool        `json:"solvable"`
//	Desc          string      `json:"desc"`
//	Safe          bool        `json:"safe"`
//}

type Input struct {
	Txid     string `json:"txid"`
	Vout     int    `json:"vout"`
	Sequence int    `json:"sequence,omitempty"`
}

type SignRawTransactionResult struct {
	Hex      string `json:"hex"`
	Complete bool   `json:"complete"`
}

type WalletAddressBalance struct {
	Address  string              `json:"address"`
	Balances []OmniWalletBalance `json:"balances"`
}

type OmniWalletBalance struct {
	Propertyid uint32 `json:"propertyid"`
	Name       string `json:"name"`
	Balance    string `json:"balance"`
	Reserved   string `json:"reserved"`
	Frozen     string `json:"frozen"`
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

type AddressBalance struct {
	Balance  json.Number `json:"balance"`
	Received json.Number `json:"received"`
	Immature json.Number `json:"immature"`
}

type AddressObject struct {
	Addresses []string `json:"addresses"`
}

type ListUnspentResult struct {
	Txid          string  `json:"txid"`
	Vout          uint32  `json:"vout"`
	Address       string  `json:"address"`
	Label         string  `json:"label"`
	ScriptPubKey  string  `json:"scriptPubKey"`
	Amount        float64 `json:"amount"`
	Confirmations int64   `json:"confirmations"`
	RedeemScript  string  `json:"redeemScript,omitempty"`
	WitnessScript string  `json:"witnessScript,omitempty"`
	Spendable     bool    `json:"spendable"`
	Solvable      bool    `json:"solvable"`
	Desc          string  `json:"desc"`
	Safe          bool    `json:"safe"`
}

type BlockHeader struct {
	Hash              string      `json:"hash"`
	Confirmations     json.Number `json:"confirmations"`
	Height            json.Number `json:"height"`
	Version           json.Number `json:"version"`
	VersionHex        string      `json:"versionHex"`
	MerkleRoot        string      `json:"merkleroot"`
	Time              json.Number `json:"time"`
	MedianTime        json.Number `json:"mediantime"`
	Nonce             json.Number `json:"nonce"`
	Bits              string      `json:"bits"`
	Difficulty        json.Number `json:"difficulty"`
	ChainWork         string      `json:"chainwork"`
	Ntx               json.Number `json:"nTx"`
	PreviousBlockHash string      `json:"previousblockhash"`
	NextBlockHash     string      `json:"nextblockhash"`
}

// GetBlockVerboseTxResult models the data from the getblock command when the
// verbose flag is set to 2.  When the verbose flag is set to 0, getblock returns a
// hex-encoded string. When the verbose flag is set to 1, getblock returns an object
// whose tx field is an array of transaction hashes. When the verbose flag is set to 2,
// getblock returns an object whose tx field is an array of raw transactions.
// Use GetBlockVerboseResult to unmarshal data received from passing verbose=1 to getblock.
type GetBlockVerboseTxResult struct {
	Hash          string              `json:"hash"`
	Confirmations int64               `json:"confirmations"`
	StrippedSize  int32               `json:"strippedsize"`
	Size          int32               `json:"size"`
	Weight        int32               `json:"weight"`
	Height        int64               `json:"height"`
	Version       int32               `json:"version"`
	VersionHex    string              `json:"versionHex"`
	MerkleRoot    string              `json:"merkleroot"`
	Tx            []BtcRawTransaction `json:"tx,omitempty"`
	Time          int64               `json:"time"`
	Nonce         uint32              `json:"nonce"`
	Bits          string              `json:"bits"`
	Difficulty    float64             `json:"difficulty"`
	PreviousHash  string              `json:"previousblockhash"`
	NextHash      string              `json:"nextblockhash,omitempty"`
}

// TxRawResult models the data from the getrawtransaction command.
type TxRawResult struct {
	Hex           string `json:"hex"`
	Txid          string `json:"txid"`
	Hash          string `json:"hash,omitempty"`
	Size          int32  `json:"size,omitempty"`
	Vsize         int32  `json:"vsize,omitempty"`
	Weight        int32  `json:"weight,omitempty"`
	Version       int32  `json:"version"`
	LockTime      uint32 `json:"locktime"`
	Vin           []Vin  `json:"vin"`
	Vout          []Vout `json:"vout"`
	BlockHash     string `json:"blockhash,omitempty"`
	Confirmations uint64 `json:"confirmations,omitempty"`
	Time          int64  `json:"time,omitempty"`
	Blocktime     int64  `json:"blocktime,omitempty"`
	BlockHeight   uint64 `json:"height" bson:"height"`
}

type AddressTxIdsParams struct {
	Addresses []string    `json:"addresses"`
	Start     json.Number `json:"start,omitempty"`
	End       json.Number `json:"end,omitempty"`
}