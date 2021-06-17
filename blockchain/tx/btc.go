package tx

import (
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"strings"
	"vbh/btc-plugins/api/types"
	"vbh/btc-plugins/blockchain/omni"
	"vbh/btc-plugins/log"
	"vbh/btc-plugins/persistence"
)

var pendingTx = make(chan types.PendingTransaction)

func init() {
	go processTx()
}

func processTx() {
	for {
		select {
		case tx := <-pendingTx:
			SaveTransaction(tx)
		}
	}
}

func PendingTransaction(t types.PendingTransaction) error {
	e := db.C(types.PendingTransactionName).Insert(t)
	if e != nil {
		return e
	}
	pendingTx <- t

	var txIds []types.PendingTransaction
	e = db.C(types.PendingTransactionName).Find(nil).All(&txIds)
	if e != nil {
		return e
	}
	for i := range txIds {
		SaveTransaction(txIds[i])
	}
	fmt.Println("PendingTransaction end, txids=", txIds)
	return nil
}

func SaveTransaction(t types.PendingTransaction) {
	client := omni.GetClient()
	tx, e := client.GetBtcRawTransaction(t.TxId, true)
	if e != nil {
		log.Errorf("GetBtcRawTransaction failed ,txid = %s,err=%v", t.TxId, e)
		return
	}
	omniTx, e := client.GetTransaction(t.TxId)
	if e == nil {
		saveOmniTransaction(omniTx)
	} else if !strings.Contains(e.Error(), "No Omni Layer Protocol transaction") && !strings.Contains(e.Error(), "Generic transaction population failure") {
		log.Errorf("Get omni Transaction failed ,txid = %s,err=%v", t.TxId, e)
	}

	header, e := client.GetBlockHeader(tx.BlockHash)
	if e != nil {
		log.Errorf("GetBlockHeader failed ,txid = %s,err=%v", t.TxId, e)
		return
	}
	height, e := header.Height.Int64()
	if e != nil {
		log.Errorf("header Height failed ,txid = %s,err=%v", t.TxId, e)
		return
	}
	tx.BlockHeight = uint64(height)
	_, e = db.C(types.BtcRawTransactionName).Upsert(bson.M{"txid": tx.Txid}, bson.M{"$set": tx})
	if e != nil {
		log.Errorf("Upsert failed ,txid = %s,err=%v", t.TxId, e)
		return
	}
	FindVinTransactions(tx)
	_, e = db.C(types.PendingTransactionName).RemoveAll(bson.M{"txid": t.TxId})
	if e != nil {
		_ = db.C(types.BtcTransactionName).Remove(tx)
		log.Errorf("Remove BtcTransactionName failed ,txid = %s,err=%v", t.TxId, e)
		return
	}
}

func saveOmniTransaction(tx *omni.Tx) error {
	_, e := db.C(types.BtcOmniTransactionName).Upsert(bson.M{"txid": tx.Txid}, bson.M{"$set": tx})
	if e != nil {
		log.Errorf("Upsert failed ,txid = %s,err=%v", tx.Txid, e)
		return e
	}
	return nil
}

func FindVinTransactions(tx *omni.BtcRawTransaction) {
	for _, v := range tx.Vin {
		if v.Txid != "" {
			SaveTransaction(types.PendingTransaction{TxId: v.Txid})
		}
	}
}
