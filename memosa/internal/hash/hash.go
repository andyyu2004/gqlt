package hash

import (
	"crypto/md5"
	"encoding/binary"

	"github.com/andyyu2004/memosa/internal/lib"
)

const hashSize = 16

type Hashed [hashSize]byte

func Hash(value any) Hashed {
	hasher := md5.New()
	lib.Assert(hasher.Size() == hashSize)
	lib.Ensure(binary.Write(hasher, binary.NativeEndian, value))
	return Hashed(hasher.Sum(nil))
}
