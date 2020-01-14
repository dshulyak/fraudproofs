package fraud

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestExtendMatrixAllStages(t *testing.T) {
	stages := []string{
		`baz buz bIz b]z
pit pot pet pct
xxx xxx xxx xxx
xxx xxx xxx xxx`,
		`baz buz bIz b]z
pit pot pet pct
Fqf FAf xxx xxx
Tyh T[h xxx xxx`,
		`baz buz bIz b]z
pit pot pet pct
Fqf FAf Ff F!f
Tyh T[h T=h Th`,
	}

	mx := NewMatrix(4) // size of extended matrix
	row1 := sliceFromString("baz buz")
	row2 := sliceFromString("pit pot")
	mx.AddRow(0, 0, row1)
	mx.AddRow(1, 0, row2)

	ext := NewExtender(2, 3)
	stage := 0
	require.NoError(t, ext.ExtendHorizontally(mx, 0, 2))
	require.Equal(t, stages[stage], mx.PrettyString(false))
	stage++
	require.NoError(t, ext.ExtendVertically(mx, 0, 2))
	require.Equal(t, stages[stage], mx.PrettyString(false))
	stage++
	require.NoError(t, ext.ExtendHorizontally(mx, 2, 4))
	require.Equal(t, stages[stage], mx.PrettyString(false))
}
