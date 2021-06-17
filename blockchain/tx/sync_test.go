package tx

import (
	"fmt"
	"testing"
	"vbh/btc-plugins/blockchain/omni"
	"vbh/btc-plugins/log"
)

func TestGetAddressTxIds(t *testing.T) {
	log.InitLogger("logs", "btc-plugins", "debug")
	omni.Init("127.0.0.1:18332", "rpc", "rpc", 1)
	client := omni.GetClient()
	ids, e := client.GetAddressTxIds(&omni.AddressTxIdsParams{Addresses: []string{}, Start: "0"})
	if e != nil {
		panic(e)
	}
	fmt.Println(ids)
}
