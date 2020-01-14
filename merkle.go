package fraud

import (
	"bytes"
	"math/bits"

	"golang.org/x/crypto/blake2s"
)

var (
	leafDomain  = []byte{0}
	innerDomain = []byte{1}

	modTwoMask = 0b1
)

// Root compute root for the complete merkle tree.
func Root(chunks [][]byte) []byte {
	leafs := computeLeafs(make([][]byte, 0, len(chunks)), chunks)
	return computeRoot(leafs)
}

// Proof generates a merkle proof, starting from a sibling leaf, excludes root from proof
// as it must be known.
func Proof(chunks [][]byte, position int) [][]byte {
	leafs := computeLeafs(make([][]byte, 0, len(chunks)), chunks)
	size := bits.Len(uint(len(leafs)))
	return computeProof(make([][]byte, 0, size-1), leafs, position)
}

// Verify that leaf at the specified position is a part of the complete merkle tree with a given root.
func Verify(root, chunk []byte, proof [][]byte, position int) bool {
	leaf := leafHash(chunk)
	return bytes.Compare(root, computeRootFromProof(leaf, proof, position)) == 0
}

func computeProof(dst [][]byte, leafs [][]byte, position int) [][]byte {
	digest, _ := blake2s.New256(nil)
	for len(leafs) > 1 {
		if position&modTwoMask == 0 {
			dst = append(dst, leafs[position+1])
		} else {
			dst = append(dst, leafs[position-1])
		}
		position = position >> 1

		half := len(leafs) / 2
		for i := 0; i < half; i++ {
			inner := make([]byte, 0, 32)
			digest.Write(innerDomain)
			digest.Write(leafs[i*2])
			digest.Write(leafs[i*2+1])
			leafs[i] = digest.Sum(inner)
			digest.Reset()
		}
		leafs = leafs[:half]
	}
	return dst
}

func computeRootFromProof(leaf []byte, proof [][]byte, position int) []byte {
	digest, _ := blake2s.New256(nil)
	for _, sibling := range proof {
		digest.Write(innerDomain)
		if position&modTwoMask == 0 {
			digest.Write(leaf)
			digest.Write(sibling)
		} else {
			digest.Write(sibling)
			digest.Write(leaf)
		}
		leaf = digest.Sum(leaf[:0])
		digest.Reset()
		position = position >> 1
	}
	return leaf
}

func computeLeafs(dst [][]byte, chunks [][]byte) [][]byte {
	for i := range chunks {
		dst = append(dst, leafHash(chunks[i]))
	}
	return dst
}

func computeRoot(leafs [][]byte) []byte {
	digest, _ := blake2s.New256(nil)
	for len(leafs) > 1 {
		half := len(leafs) / 2
		for i := 0; i < half; i++ {
			inner := make([]byte, 0, 32)
			digest.Write(innerDomain)
			digest.Write(leafs[i*2])
			digest.Write(leafs[i*2+1])
			leafs[i] = digest.Sum(inner)
			digest.Reset()
		}
		leafs = leafs[:half]
	}
	return leafs[0]
}

func leafHash(chunk []byte) []byte {
	digest, err := blake2s.New256(nil)
	if err != nil {
		panic(err.Error())
	}
	digest.Write(leafDomain)
	digest.Write(chunk)
	return digest.Sum(make([]byte, 0, 32))
}
