package omni

import (
	"encoding/json"
)

type OmniClient struct {
	ConnCfg    *ConnConfig
	PropertyId uint32
}

func NewOmniClient(connCfg *ConnConfig, propertyId uint32) *OmniClient {
	return &OmniClient{ConnCfg: connCfg, PropertyId: propertyId}
}

func (o *OmniClient) GetTransaction(txid string) (*Tx, error) {
	c := NewClient(o.ConnCfg)
	result := Tx{}
	if err := c.Call(&result, "omni_gettransaction", txid); err != nil {
		return &result, err
	}
	return &result, nil
}

func (o *OmniClient) GetBalance(address string) (balance, reserved string, err error) {
	c := NewClient(o.ConnCfg)
	var result BalanceResult

	if err = c.Call(&result, "omni_getbalance", address, o.PropertyId); err != nil {
		return
	}
	return result.Balance, result.Reserved, nil
}

func (o *OmniClient) FundedSend(fromaddress, toaddress string, propertyid uint32, amount, feeaddress string) (string, error) {
	c := NewClient(o.ConnCfg)
	var hash string
	if err := c.Call(&hash, "omni_funded_send", fromaddress, toaddress, propertyid, amount, feeaddress); err != nil {
		return hash, err
	}
	return hash, nil
}

func (o *OmniClient) SendAll(fromaddress, toaddress string, ecosystem uint32, redeemaddress, referenceamount string) (string, error) {
	c := NewClient(o.ConnCfg)
	var hash string
	if err := c.Call(&hash, "omni_sendall", fromaddress, toaddress, ecosystem, redeemaddress, referenceamount); err != nil {
		return hash, err
	}
	return hash, nil
}

func (o *OmniClient) GetBlockCount() (int64, error) {
	c := NewClient(o.ConnCfg)
	var count int64 = 0
	if err := c.Call(&count, "getblockcount"); err != nil {
		return count, err
	}
	return count, nil
}

func (o *OmniClient) GetBlockHash(height int64) (string, error) {
	c := NewClient(o.ConnCfg)
	var hash string
	if err := c.Call(&hash, "getblockhash", height); err != nil {
		return hash, err
	}
	return hash, nil
}

func (o *OmniClient) GetBlock(blockHash string) (*GetBlockVerboseTxResult, error) {
	c := NewClient(o.ConnCfg)
	var block GetBlockVerboseTxResult
	if err := c.Call(&block, "getblock", blockHash, 2); err != nil {
		return &block, err
	}
	return &block, nil
}

func (o *OmniClient) GetBtcBalance() (balance string, err error) {
	c := NewClient(o.ConnCfg)
	var b json.Number
	if err = c.Call(&b, "getbalance"); err != nil {
		return
	}
	balance = b.String()
	return
}

func (o *OmniClient) GetWalletBalance() (balance string, err error) {
	c := NewClient(o.ConnCfg)
	var b json.Number
	if err = c.Call(&b, "omni_getwalletbalances"); err != nil {
		return
	}
	balance = b.String()
	return
}

func (o *OmniClient) GetWalletAddressBalances() ([]WalletAddressBalance, error) {
	c := NewClient(o.ConnCfg)
	var b []WalletAddressBalance
	if err := c.Call(&b, "omni_getwalletaddressbalances"); err != nil {
		return nil, err
	}
	return b, nil
}

func (o *OmniClient) CreateBtcRawTransaction(inputs []Input, outputs []map[string]string) (hex string, err error) {
	c := NewClient(o.ConnCfg)
	err = c.Call(&hex, "createrawtransaction", inputs, outputs)
	return
}

func (o *OmniClient) SignBtcRawTransaction(hexstring string) (r *SignRawTransactionResult, err error) {
	c := NewClient(o.ConnCfg)
	err = c.Call(&r, "signrawtransaction", hexstring)
	return
}

func (o *OmniClient) SendBtcRawTransaction(hexstring string) (hex string, err error) {
	c := NewClient(o.ConnCfg)
	err = c.Call(&hex, "sendrawtransaction", hexstring)
	return
}

func (o *OmniClient) Send(fromaddress, toaddress string, amount string, redeemaddress string, referenceamount string) (string, error) {
	c := NewClient(o.ConnCfg)
	var hash string
	if err := c.Call(&hash, "omni_send", fromaddress, toaddress, o.PropertyId, amount); err != nil {
		return hash, err
	}
	return hash, nil
}

func (o *OmniClient) ListBtcUnspent(minconf uint64, maxconf ...uint64) (unspents []ListUnspentResult, err error) {
	c := NewClient(o.ConnCfg)
	if len(maxconf) == 0 {
		err = c.Call(&unspents, "listunspent", minconf)
	} else {
		err = c.Call(&unspents, "listunspent", minconf, maxconf[0])
	}
	return unspents, err
}

func (o *OmniClient) WalletPassphrase(passphrase string, timeout int) error {
	c := NewClient(o.ConnCfg)
	err := c.CallWithoutResponse("walletpassphrase", passphrase, timeout)
	if err != nil {
		return err
	}
	return nil
}

func (o *OmniClient) WalletLock() error {
	c := NewClient(o.ConnCfg)
	err := c.CallWithoutResponse("walletlock")
	return err
}

func (o *OmniClient) ListTransactions(address string, count, skip, startBlock, endBlock uint64) ([]Tx, error) {
	c := NewClient(o.ConnCfg)
	var result []Tx
	if err := c.Call(&result, "omni_listtransactions", address, count, skip, startBlock, endBlock); err != nil {
		return result, err
	}
	return result, nil
}

/**
Arguments:
1. minconf                            (numeric, optional, default=1) The minimum confirmations to filter
2. maxconf                            (numeric, optional, default=9999999) The maximum confirmations to filter
3. addresses                          (json array, optional, default=empty array) A json array of bitcoin addresses to filter
     [
       "address",                     (string) bitcoin address
       ...
     ]
4. include_unsafe                     (boolean, optional, default=true) Include outputs that are not safe to spend
                                      See description of "safe" attribute below.
5. query_options                      (json object, optional) JSON with query options
     {
       "minimumAmount": amount,       (numeric or string, optional, default=0) Minimum value of each UTXO in BTC
       "maximumAmount": amount,       (numeric or string, optional, default=unlimited) Maximum value of each UTXO in BTC
       "maximumCount": n,             (numeric, optional, default=unlimited) Maximum number of UTXOs
       "minimumSumAmount": amount,    (numeric or string, optional, default=unlimited) Minimum sum value of all UTXOs in BTC
     }

Result:
[                   (array of json object)
  {
    "txid" : "txid",          (string) the transaction id
    "vout" : n,               (numeric) the vout value
    "address" : "address",    (string) the bitcoin address
    "label" : "label",        (string) The associated label, or "" for the default label
    "scriptPubKey" : "key",   (string) the script key
    "amount" : x.xxx,         (numeric) the transaction output amount in BTC
    "confirmations" : n,      (numeric) The number of confirmations
    "redeemScript" : "script" (string) The redeemScript if scriptPubKey is P2SH
    "witnessScript" : "script" (string) witnessScript if the scriptPubKey is P2WSH or P2SH-P2WSH
    "spendable" : xxx,        (bool) Whether we have the private keys to spend this output
    "solvable" : xxx,         (bool) Whether we know how to spend this output, ignoring the lack of keys
    "desc" : xxx,             (string, only when solvable) A descriptor for spending this output
    "safe" : xxx              (bool) Whether this output is considered safe to spend. Unconfirmed transactions
                              from outside keys and unconfirmed replacement transactions are considered unsafe
                              and are not eligible for spending by fundrawtransaction and sendtoaddress.
  }
 */
func (o *OmniClient) ListBtcUnspentFilterByAddress(address []string) (unspents []ListUnspentResult, err error) {
	c := NewClient(o.ConnCfg)
	err = c.Call(&unspents, "listunspent", 0, 9999999, address) //min 0 max 9999999
	return unspents, err
}

func (o *OmniClient) GetReceivedBtcByAddress(address string, minconf ...uint64) (string, error) {
	c := NewClient(o.ConnCfg)
	var amount json.Number
	var err error
	if len(minconf) == 0 {
		err = c.Call(&amount, "getreceivedbyaddress", address)
	} else {
		err = c.Call(&amount, "getreceivedbyaddress", address, minconf[0])
	}
	return amount.String(), err
}

func (o *OmniClient) SetTxFee(feeAmount string) (result bool, err error) {
	c := NewClient(o.ConnCfg)
	err = c.Call(&result, "settxfee", feeAmount)
	return
}

func (o *OmniClient) GetBtcTransaction(txid string, includeWatchOnly bool) (*BitcoinTransaction, error) {
	c := NewClient(o.ConnCfg)
	tx := BitcoinTransaction{}
	err := c.Call(&tx, "gettransaction", txid, includeWatchOnly)
	return &tx, err
}

func (o *OmniClient) GetAddressBalance(addresses *AddressObject) (*AddressBalance, error) {
	c := NewClient(o.ConnCfg)
	b := AddressBalance{}
	err := c.Call(&b, "getaddressbalance", addresses)
	return &b, err
}

func (o *OmniClient) ImportAddress(address, label string, rescan, p2sh bool) error {
	c := NewClient(o.ConnCfg)
	return c.CallWithoutResponse("importaddress", address, label, rescan, p2sh)
}

/*
getrawtransaction "txid" ( verbose "blockhash" )

Return the raw transaction data.

By default this function only works for mempool transactions. When called with a blockhash
argument, getrawtransaction will return the transaction if the specified block is available and
the transaction is found in that block. When called without a blockhash argument, getrawtransaction
will return the transaction if it is in the mempool, or if -txindex is enabled and the transaction
is in a block in the blockchain.

Hint: Use gettransaction for wallet transactions.

If verbose is 'true', returns an Object with information about 'txid'.
If verbose is 'false' or omitted, returns a string that is serialized, hex-encoded data for 'txid'.

Arguments:
1. txid         (string, required) The transaction id
2. verbose      (boolean, optional, default=false) If false, return a string, otherwise return a json object
3. blockhash    (string, optional) The block in which to look for the transaction

Result (if verbose is not set or set to false):
"data"      (string) The serialized, hex-encoded data for 'txid'

Result (if verbose is set to true):
{
  "in_active_chain": b, (bool) Whether specified block is in the active chain or not (only present with explicit "blockhash" argument)
  "hex" : "data",       (string) The serialized, hex-encoded data for 'txid'
  "txid" : "id",        (string) The transaction id (same as provided)
  "hash" : "id",        (string) The transaction hash (differs from txid for witness transactions)
  "size" : n,             (numeric) The serialized transaction size
  "vsize" : n,            (numeric) The virtual transaction size (differs from size for witness transactions)
  "weight" : n,           (numeric) The transaction's weight (between vsize*4-3 and vsize*4)
  "version" : n,          (numeric) The version
  "locktime" : ttt,       (numeric) The lock time
  "vin" : [               (array of json objects)
     {
       "txid": "id",    (string) The transaction id
       "vout": n,         (numeric)
       "scriptSig": {     (json object) The script
         "asm": "asm",  (string) asm
         "hex": "hex"   (string) hex
       },
       "sequence": n      (numeric) The script sequence number
       "txinwitness": ["hex", ...] (array of string) hex-encoded witness data (if any)
     }
     ,...
  ],
  "vout" : [              (array of json objects)
     {
       "value" : x.xxx,            (numeric) The value in BTC
       "n" : n,                    (numeric) index
       "scriptPubKey" : {          (json object)
         "asm" : "asm",          (string) the asm
         "hex" : "hex",          (string) the hex
         "reqSigs" : n,            (numeric) The required sigs
         "type" : "pubkeyhash",  (string) The type, eg 'pubkeyhash'
         "addresses" : [           (json array of string)
           "address"        (string) bitcoin address
           ,...
         ]
       }
     }
     ,...
  ],
  "blockhash" : "hash",   (string) the block hash
  "confirmations" : n,      (numeric) The confirmations
  "blocktime" : ttt         (numeric) The block time in seconds since epoch (Jan 1 1970 GMT)
  "time" : ttt,             (numeric) Same as "blocktime"
}
 */
func (o *OmniClient) GetBtcRawTransaction(args ...interface{}) (*BtcRawTransaction, error) {
	c := NewClient(o.ConnCfg)
	rtx := BtcRawTransaction{}
	err := c.Call(&rtx, "getrawtransaction", args...)
	return &rtx, err
}

func (o *OmniClient) GetBlockHeader(blockhash string) (*BlockHeader, error) {
	c := NewClient(o.ConnCfg)
	h := BlockHeader{}
	e := c.Call(&h, "getblockheader", blockhash)
	return &h, e
}

func (o *OmniClient) GetAddressTxIds(p *AddressTxIdsParams) ([]string, error) {
	c := NewClient(o.ConnCfg)
	var txIds []string
	if p.End == json.Number("0") {
		p.End = json.Number("9999999")
	}
	e := c.Call(&txIds, "getaddresstxids", p)
	return txIds, e
}
