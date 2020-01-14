package fraud

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func sliceFromString(row string) [][]byte {
	parts := strings.Split(row, " ")
	rst := make([][]byte, len(parts))
	for i := range parts {
		rst[i] = []byte(parts[i])
	}
	return rst
}

func TestMatrixAddRow(t *testing.T) {
	expected := `foo bar
baz tet`
	mx := NewMatrix(2)
	mx.AddRow(0, 0, [][]byte{[]byte("foo"), []byte("bar")})
	mx.AddRow(1, 0, [][]byte{[]byte("baz"), []byte("tet")})

	require.Equal(t, expected, mx.PrettyString(false))
}

func TestMatrixAddColumn(t *testing.T) {
	expected := `foo bar
baz tet`
	mx := NewMatrix(2)
	mx.AddColumn(0, 0, [][]byte{[]byte("foo"), []byte("baz")})
	mx.AddColumn(1, 0, [][]byte{[]byte("bar"), []byte("tet")})

	require.Equal(t, expected, mx.PrettyString(false))
}

func TestMatrixGetRow(t *testing.T) {
	mx := NewMatrix(4)
	row := sliceFromString("das fod fbq vwq")
	mx.AddRow(0, 0, row)

	require.Equal(t, row, mx.GetRow(nil, 0))
}

func TestMatrixGetColumn(t *testing.T) {
	mx := NewMatrix(4)
	column := sliceFromString("das fod fbq vwq")
	mx.AddColumn(0, 0, column)

	require.Equal(t, column, mx.GetColumn(nil, 0))
}
