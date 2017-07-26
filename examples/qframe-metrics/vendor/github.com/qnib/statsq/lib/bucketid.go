package statsq

import (
	"crypto/sha1"
	"encoding/hex"
	"github.com/qnib/qframe-types"
	"log"
	"strings"
)

type BucketID struct {
	ID         string
	BucketName string
	Dimensions qtypes.Dimensions
}

func NewBucketID(name string, dims qtypes.Dimensions) BucketID {
	bid := BucketID{
		ID:         "",
		BucketName: name,
		Dimensions: dims,
	}
	bid.GenerateID()
	return bid
}

func (bid *BucketID) GetDims() map[string]string {
	return bid.Dimensions.Map
}

func (bid *BucketID) GenerateID() {
	if bid.ID != "" {
		log.Panicf("BucketID already has ID '%s'", bid.ID)
	}
	idRaw := []string{bid.BucketName}
	for key, val := range bid.Dimensions.Map {
		idRaw = append(idRaw, key+"="+val)
	}
	s := strings.Join(idRaw, "_")
	bid.ID = GenID(s)
}

func GenID(s string) string {
	h := sha1.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}
