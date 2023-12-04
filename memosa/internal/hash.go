package internal

import (
	"crypto/md5"
	"encoding/binary"
)

const hashSize = 16

type hashed [hashSize]byte

func hash(value any) hashed {
	hasher := md5.New()
	assert(hasher.Size() == hashSize)
	ensure(binary.Write(hasher, binary.NativeEndian, value))
	return hashed(hasher.Sum(nil))
}
