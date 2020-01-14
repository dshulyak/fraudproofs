package fraud

import (
	"crypto/rand"
	"testing"

	"github.com/stretchr/testify/require"
)

func genChunks(n int) [][]byte {
	chunks := make([][]byte, n)
	for i := range chunks {
		chunk := make([]byte, 100)
		rand.Read(chunk)
		chunks[i] = chunk
	}
	return chunks
}

func TestMerkleProof(t *testing.T) {
	chunks := genChunks(8)

	root := Root(chunks)
	for i := range chunks {
		proof := Proof(chunks, i)
		require.True(t, Verify(root, chunks[i], proof, i), "proof for %dth chunk is not valid", i)
	}
}

func TestMerkleProofInvalid(t *testing.T) {
	chunks := genChunks(8)

	root := Root(chunks)
	for i := range chunks {
		proof := Proof(chunks, i)
		chunk := make([]byte, len(chunks[i]))
		copy(chunk, chunks[i])
		chunk[0] ^= 0xff
		require.False(t, Verify(root, chunk, proof, i), "proof for %dth corrupted chunk is valid", i)
	}
}
