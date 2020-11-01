package sha1

import (
	"encoding/binary"
)

const (
	A = 0x67452301
	B = 0xEFCDAB89
	C = 0x98BADCFE
	D = 0x10325476
	E = 0xC3D2E1F0
)

type SHA1Values struct {
	a int
	b int
	c int
	d int
	e int
}

type SHA1Custom interface {
	GetHash(string) string
}

func breakStringIntoBlocks(text string) [][]byte {
	result := [][]byte{}
	byteString := []byte(text)
	size := len(byteString)
	rest := size % 512
	blocksCount := int((size - rest) / 512)
	for i := 0; i < blocksCount; i++ {
		result = append(result, byteString[i*512:(i+1)*512])
	}
	if rest != 0 {
		finalBlock := byteString[blocksCount*512:]
		if rest > 448 {
			finalBlock = append(finalBlock, 1)
			zeros := make([]byte, 511-rest)
			finalBlock = append(finalBlock, zeros...)
			result = append(result, finalBlock)
			finalBlock = make([]byte, 448)
		} else {
			if rest < 448 {
				finalBlock = append(finalBlock, 1)
			}
			if rest < 447 {
				zeros := make([]byte, 447-rest)
				finalBlock = append(finalBlock, zeros...)
			}
		}
		bytesSize := make([]byte, 64)
		binary.BigEndian.PutUint64(bytesSize, uint64(size))
		finalBlock = append(finalBlock, bytesSize...)
	}
	return result
}
