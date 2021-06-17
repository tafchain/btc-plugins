package tx

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"strings"
	"time"
	"vbh/btc-plugins/api/types"
	"vbh/btc-plugins/blockchain/omni"
	"vbh/btc-plugins/log"
	"vbh/btc-plugins/persistence"
)

func Init(period int, startHeight int64) {
	fmt.Println("period=", period, ",startHeight=", startHeight)
	go func() {
		client := omni.GetClient()
		ticker := time.NewTicker(100 * time.Duration(period) * time.Millisecond)
	NewLoop:
		for {
			var txIds []types.PendingTransaction
			e := db.C(types.PendingTransactionName).Find(nil).All(&txIds)
			if e != nil {
				log.Error(0, "Find All txIds failed,er=", e)
				continue
			}
			for i := range txIds {
				SaveTransaction(txIds[i])
			}
			if len(txIds) > 0 {
				log.Info("PendingTransaction end, txids=", txIds)
			}

			var addresses []types.BtcUserInfo
			e = db.C(types.BtcUserInfoName).Find(nil).All(&addresses)
			if e != nil {
				log.Error(1, "get addresses failed,err=", e)
				continue
			}
			mAddr := make(map[string]interface{})
			for i := range addresses {
				mAddr[addresses[i].Address] = struct{}{}
				if addresses[i].State == 0 {
					ids, e := client.GetAddressTxIds(&omni.AddressTxIdsParams{Addresses: []string{addresses[i].Address}, Start: "0"})
					if e != nil {
						log.Error(7, " get GetAddressTxIds address="+addresses[i].Address+", failed,err=", e)
						continue
					}
					log.Info("GetAddressTxIds ids=", ids)
					for i := range ids {
						e = loopTx(ids[i], client, mAddr)
						if e != nil {
							log.Error(8, " get loopTx failed,err=", e)
							continue
						}
					}
					ch, e := db.C(types.BtcUserInfoName).UpdateAll(bson.M{"address": addresses[i].Address}, bson.M{"$set": bson.M{"state": types.BtcUserInfoStateSynced, "update_time": time.Now().Unix()}})
					if e != nil {
						log.Error(8, " Update btc_user_info failed err=", e)
						continue
					}
					log.Infof("update address=%s"+" for change %+v", addresses[i].Address, ch)
				}
			}
			//log.Info("len mAddr=", len(mAddr))

			if len(mAddr) == 0 {
				continue
			}

			blockSync := types.BtcBlockSync{CreateTime: time.Now().Unix(), UpdateTime: time.Now().Unix()}

			select {
			case <-ticker.C:
				log.Debug("syncing...")
				curHeight, e := getBlock(client)
				if e != nil {
					log.Error(2, "get block height failed,err=", e)
					continue
				}
				log.Debug("last synced block height = ", curHeight)
				if curHeight == 0 {
					curHeight = startHeight
				}
				newHeight, e := client.GetBlockCount()
				if e != nil {
					log.Error(3, "GetBlockCount failed,err=", e)
					continue
				}
				log.Debug("block newHeight = ", newHeight)
				if newHeight == curHeight {
					continue
				}
				log.Infof("last synced block height = %d,block newHeight = %d", curHeight, newHeight)
				blockHash, e := client.GetBlockHash(curHeight + 1)
				if e != nil {
					log.Error(4, "get GetBlockHash failed,err=", e)
					continue
				}
				block, e := client.GetBlock(blockHash)
				if e != nil {
					log.Error(5, "get GetBlock failed,err=", e)
					continue
				}

				for _, tran := range block.Tx {
					e = loopTx(tran.Txid, client, mAddr)
					if e != nil {
						log.Error(6, "get loopTx failed,err=", e)
						goto NewLoop
					}
				}
				blockSync.BlockHeight = block.Height
				blockSync.BlockHash = block.Hash
				blockSync.BlockTime = block.Time
				e = db.C(types.BtcBlockSyncName).Insert(&blockSync)
				if e != nil {
					log.Error(e)
					break
				}
				log.Info("synced block ", blockSync)
				_, e = db.C(types.BtcBlockSyncName).RemoveAll(bson.M{"block_height": bson.M{"$lte": blockSync.BlockHeight - 7}})
				if e != nil {
					log.Error(e)
				}
			}
			log.Debug("syncing end...")
		}
	}()
}

func loopTx(outTxId string, client *omni.OmniClient, mAddr map[string]interface{}) error {
	//if level == 0 {
	//	return nil
	//}
	rawTx, e := getTransaction(outTxId, client)
	if e != nil {
		return e
	}
	e = checkAndSaveOmniTransaction(client, outTxId, mAddr)
	if e != nil {
		return e
	}
	//log.Errorf("%+v22222222222222", rawTx)
	//if level == 1 {
	//	return save(rawTx)
	//} else {
	exist, e := checkAndInsert(rawTx, mAddr)
	if e != nil {
		log.Error(e)
		return e
	}
	if !exist {
		for _, vin := range rawTx.Vin {
			if vin.Coinbase != "" {
				continue
			}
			if vin.Txid != "" {
				rawTx4, e := getTransaction(vin.Txid, client)
				if e != nil {
					return e
				}
				for _, vout := range rawTx4.Vout {
					if vout.N == vin.Vout {
						b, e := checkAndInsert(rawTx4, mAddr)
						if e != nil {
							log.Error(e)
							return e
						}
						if b {
							return save(rawTx)
						}
					}
				}
			}
		}
		return nil
	}

	for n := range rawTx.Vin {
		if rawTx.Vin[n].Coinbase != "" {
			continue
		}
		if rawTx.Vin[n].Txid != "" {
			//rawTx2, e := client.GetBtcRawTransaction(out.Vin[n].Txid, true)
			//if e != nil {
			//	log.Error(1, e)
			//	break
			//}
			//header, e := client.GetBlockHeader(rawTx2.BlockHash)
			//if e != nil {
			//	log.Errorf("GetBlockHeader failed ,txid = %s,err=%v", out.Vin[n].Txid, e)
			//	break
			//}
			//height, e := header.Height.Int64()
			//if e != nil {
			//	log.Errorf("header Height failed ,txid = %s,err=%v", out.Vin[n].Txid, e)
			//	break
			//}
			//rawTx2.BlockHeight = uint64(height)
			//if e != nil {
			//	log.Errorf("GetTransaction failed ,%+v", e)
			//	break
			//}

			rawTx2 := omni.BtcRawTransaction{Txid: rawTx.Vin[n].Txid}
			err := db.C(types.BtcRawTransactionName).Find(bson.M{"txid": rawTx2.Txid}).One(&rawTx2)
			if err == mgo.ErrNotFound {
				log.Info("Find prev raw, txid=", rawTx2.Txid)
				rawTx3, e := getTransaction(rawTx2.Txid, client)
				if e != nil {
					return e
				}
				e = save(rawTx3)
				if e != nil {
					log.Error(e)
					return e
				}

				//e = loopTx(rawTx2.Txid, client, nil, level-1)
				//if e != nil {
				//	return e
				//}

				//e := checkAndInsert(&rawTx2, mAddr)
				//if e != nil {
				//	log.Error(e)
				//	break
				//}
			} else if err != nil {
				log.Error("find prev raw err=", err)
				return err
			}
		}
		//}
	}

	return nil
}

func checkAndSaveOmniTransaction(client *omni.OmniClient, outTxId string, mAddr map[string]interface{}) error {
	omniTx, e := client.GetTransaction(outTxId)
	if e == nil {
		var exist bool
		if _, ok := mAddr[omniTx.Sendingaddress]; ok {
			exist = true
		}
		if _, ok := mAddr[omniTx.Referenceaddress]; ok {
			exist = true
		}
		if !exist {
			return nil
		}
		log.Infof("save omni transaction, txid=%s", omniTx.Txid)
		e = saveOmniTransaction(omniTx)
		if e != nil {
			return e
		}
	} else if !strings.Contains(e.Error(), "No Omni Layer Protocol transaction") && !strings.Contains(e.Error(), "Generic transaction population failure") {
		log.Errorf("Get omni Transaction failed ,txid = %s,err=%v", outTxId, e)
	}
	return nil
}

func getTransaction(outTxId string, client *omni.OmniClient) (*omni.BtcRawTransaction, error) {
	rawTx, e := client.GetBtcRawTransaction(outTxId, true)
	if e != nil {
		log.Error(e)
		return nil, e
	}

	header, e := client.GetBlockHeader(rawTx.BlockHash)
	if e != nil {
		log.Errorf("GetBlockHeader failed ,txid = %s,err=%v", rawTx.Txid, e)
		return nil, e
	}
	height, e := header.Height.Int64()
	if e != nil {
		log.Errorf("header Height failed ,height = %d,rawTx=%v,err=%v", height, rawTx, e)
		return nil, e
	}
	rawTx.BlockHeight = uint64(height)
	return rawTx, nil
}

func checkAndInsert(tx *omni.BtcRawTransaction, mAddr map[string]interface{}) (exist bool, e error) {
	for n := range tx.Vout {
		addresses := tx.Vout[n].ScriptPubKey.Addresses
		for i := range addresses {
			if _, ok := mAddr[addresses[i]]; ok {
				log.Infof("save tx for address=%+v", addresses[i])
				exist = true
				e = save(tx)
				return
			}
		}
	}
	return
}

func save(tx *omni.BtcRawTransaction) error {
	info, e := db.C(types.BtcRawTransactionName).Upsert(bson.M{"txid": tx.Txid}, bson.M{"$set": tx})
	log.Infof("insert tx info=%+v,txid=%+v,height=%v", info, tx.Txid, tx.BlockHeight)
	if e != nil {
		log.Errorf("insert failed ,%+v", tx)
		return e
	}
	return nil
}

func getBlock(client *omni.OmniClient) (int64, error) {
	lastSync := types.BtcBlockSync{}
	cond := db.C(types.BtcBlockSyncName).Find(nil).Sort("-block_height").Limit(1)
	n, e := cond.Count()
	if e != nil {
		return lastSync.BlockHeight, e
	}
	if n == 0 {
		return 0, nil
	}
	e = cond.One(&lastSync)
	if e != nil {
		return lastSync.BlockHeight, e
	}
	blockCount, e := client.GetBlockCount()
	if e != nil {
		return lastSync.BlockHeight, e
	}
	if lastSync.BlockHeight == blockCount {
		blockHash, err := client.GetBlockHash(blockCount)
		if err != nil {
			return lastSync.BlockHeight, e
		}
		if lastSync.BlockHash != blockHash {
			err = db.C(types.BtcBlockSyncName).Remove(&lastSync)
			if err != nil {
				return lastSync.BlockHeight, err
			}
			return getBlock(client)
		}
		return lastSync.BlockHeight, nil
	}
	return minHeight(lastSync.BlockHeight, blockCount), nil
}

func minHeight(h1 int64, h2 int64) int64 {
	if h1 > h2 {
		return h2
	}
	return h1
}
