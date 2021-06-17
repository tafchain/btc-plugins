package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"vbh/btc-plugins/api/types"
	"vbh/btc-plugins/blockchain/tx"
)
import . "vbh/btc-plugins/api/routers/internal"

func (p *Router) PendingTransaction(c *gin.Context) {
	r := Resp{C: c}
	txId, ok := c.GetPostForm("tx_id")
	log.Println("---->pending txid=", txId)
	if !ok || txId == "" {
		r.Err(InvalidParam, fmt.Errorf("invalid txid = %s", txId))
		return
	}
	if len(txId) != 64 {
		r.Err(InvalidParam, fmt.Errorf("invalid txid = %s", txId))
		return
	}
	e := tx.PendingTransaction(types.PendingTransaction{TxId: txId})
	if e != nil {
		r.Err(PendingTransactionErr1, e)
	}
	r.Ok(txId)
}
