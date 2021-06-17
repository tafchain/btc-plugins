package types

type BtcUserInfo struct {
	Address    string `json:"address" bson:"address"`
	UpdateTime int64  `json:"update_time" bson:"update_time"`
	State      int    `json:"state" bson:"state"` // 0 - 未同步历史交易 1 - 已同步历史交易
}

const (
	BtcUserInfoStateUnsync = iota
	BtcUserInfoStateSynced
)
