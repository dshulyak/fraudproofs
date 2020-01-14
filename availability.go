package fraud

import "fmt"

const (
	// matrixWidth must be computed based on the received data
	// in research paper: matrixWidth = sqrt(1/2*dataLength)

	initialK  = 4
	extendedK = initialK * 2
	chunks    = initialK * initialK
	extended  = extendedK * extendedK
)

func NewConfig(initialK int) Config {
	return Config{
		InitialK:     initialK,
		ExtendedK:    initialK * 2,
		InitialSize:  initialK * initialK,
		ExtendedSize: initialK * initialK * 4,
	}
}

type Config struct {
	InitialK     int
	ExtendedK    int
	InitialSize  int
	ExtendedSize int
}

// DataRoot computes merkle root from extended matrix of data shares.
// - split data into multiple equal sized chunks (number of chunks harcoded to 64*64)
// - create extended matrix of size 128*128 (x2 64)
// - add chunks to that matrix
// - initialize reed-solomon encoder with 64 shards and 64 parity shards
// - extend 64 row horizontally, then 64 columns vertically, then extend horizontally 64 rows
// created after previous extension
// - compute merkle root for each row and column in the extended matrix
// - compute merkle root from rows/columns roots
// - thats a block data root
// data expected to be encoded in a way to allow caller to recover original byte slices from it
func DataRoot(data []byte) []byte {
	mx := PrecomputeExtendedMatrix(data)

	roots := make([][]byte, 0, mx.RowSize()+mx.ColumnSize())
	for i := 0; i < mx.RowSize(); i++ {
		row := mx.GetRow(nil, i)
		roots = append(roots, Root(row))
	}
	for i := 0; i < mx.ColumnSize(); i++ {
		column := mx.GetColumn(nil, i)
		roots = append(roots, Root(column))
	}
	return Root(roots)
}

func PrecomputeExtendedMatrix(data []byte) *Matrix {
	mx := MatrixFromData(data)

	shareSize := len(data) / chunks
	ext := NewExtender(initialK, shareSize)
	if err := ext.Extend(mx); err != nil {
		panic("matrix is invalid " + err.Error())
	}
	return mx
}

func MatrixFromData(data []byte) *Matrix {
	// TODO use config
	lth := len(data)
	shareSize := lth / chunks
	offset := 0

	m := NewMatrix(extendedK)

	for i := 0; i < initialK; i++ {
		row := make([][]byte, initialK)
		for i := range row {
			share := make([]byte, shareSize)
			if offset < lth {
				offset += copy(share, data[offset:])
			}
			row[i] = share
		}
		m.AddRow(i, 0, row)
	}

	return m
}

type Share struct {
	Row, Column int
}

type ShareProof struct {
	Share
	Data        []byte
	RowProof    [][]byte
	ColumnProof [][]byte
}

// CreateShareProof returns share data chunk with column and row merkle proofs.
func CreateShareProof(mx *Matrix, share Share) *ShareProof {
	row := mx.GetRow(nil, share.Row)
	column := mx.GetColumn(nil, share.Column)
	// TODO is it enough to return one of them? or both are required?
	return &ShareProof{
		Share:       share,
		Data:        mx.Get(share.Row, share.Column),
		RowProof:    Proof(row, share.Column),
		ColumnProof: Proof(column, share.Row),
	}
}

// VerifyShareProof verifies both row and column proofs for a data chunk.
// Returns error if either one of them is not valid.
func VerifyShareProof(rowRoot, columnRoot []byte, proof *ShareProof) error {
	if !Verify(rowRoot, proof.Data, proof.RowProof, proof.Share.Column) {
		return fmt.Errorf("row merkle proof is invalid")
	}
	if !Verify(columnRoot, proof.Data, proof.ColumnProof, proof.Share.Row) {
		return fmt.Errorf("column merkle proof is invalid")
	}
	return nil
}
