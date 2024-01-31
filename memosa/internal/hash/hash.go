package hash

import (
	"crypto/md5"
	"encoding/json"

	"github.com/andyyu2004/gqlt/memosa/lib"
)

const hashSize = 16

type Hashed [hashSize]byte

func Hash(value any) Hashed {
	hasher := md5.New()
	lib.Assert(hasher.Size() == hashSize)

	lib.Ensure(json.NewEncoder(hasher).Encode(value))

	return Hashed(hasher.Sum(nil))
}
