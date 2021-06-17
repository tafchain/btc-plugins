package types

type BtcBlockSync struct {
	BlockHeight int64  `json:"block_height" bson:"block_height"`
	BlockHash   string `json:"block_hash" bson:"block_hash"`
	BlockTime   int64  `json:"block_time" bson:"block_time"`
	CreateTime  int64  `json:"create_time" bson:"create_time"`
	UpdateTime  int64  `json:"update_time" bson:"update_time"`
}
