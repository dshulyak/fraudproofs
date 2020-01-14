package fraud

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"
)

// FIXME this is not a valid, very lazy test
func TestMatrixFromData(t *testing.T) {
	example := `d2ba8e 1826d2 16d057 edee48 787878 787878 787878 787878
1c13ce 3e579f 948058 47a1c7 787878 787878 787878 787878
5b5441 69d10c abafe3 f161f6 787878 787878 787878 787878
3ded87 f50f8c f07eac 79e5ed 787878 787878 787878 787878
787878 787878 787878 787878 787878 787878 787878 787878
787878 787878 787878 787878 787878 787878 787878 787878
787878 787878 787878 787878 787878 787878 787878 787878
787878 787878 787878 787878 787878 787878 787878 787878`

	data := make([]byte, 4<<10)
	rand.Seed(100)
	rand.Read(data)

	mx := MatrixFromData(data)
	require.Equal(t, example, mx.PrettyString(true))
}

func TestShareProof(t *testing.T) {
	data := make([]byte, 4<<10)
	rand.Read(data)

	mx := PrecomputeExtendedMatrix(data)
	share := Share{3, 4}
	rowRoot := Root(mx.GetRow(nil, share.Row))
	columnRoot := Root(mx.GetColumn(nil, share.Column))

	proof := CreateShareProof(mx, share)
	require.NoError(t, VerifyShareProof(rowRoot, columnRoot, proof))
}
