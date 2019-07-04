package types

// Key describes the keys
type Key struct {
	ID       string `json:"id"    bson:"id"`
	Issued   bool   `json:"issued"   bson:"issued"`
	Canceled bool   `json:"canceled" bson:"canceled"`
}
