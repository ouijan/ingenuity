package utils

import "hash/fnv"

type HashId = uint64

func Hash(s string) HashId {
	h := fnv.New64a()
	h.Write([]byte(s))
	return h.Sum64()
}
