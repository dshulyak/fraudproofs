package fraud

import (
	"bytes"
	"fmt"
	"strconv"
)

var (
	empty = []byte("xxx")
)

func NewMatrix(k int) *Matrix {
	return &Matrix{
		k:      k,
		chunks: make([][]byte, k*k),
	}
}

// Matrix is a square matrix where k must be a power of two.
// Also application will use chunks with equal size, but this is not in the scope of Matrix responsibility.
type Matrix struct {
	k int
	// in total there will be k*k chunks
	chunks [][]byte
}

// AddRow adds a row to matrix. Arguments can't be modified after passing to matrix.
func (m *Matrix) AddRow(ridx, start int, row [][]byte) {
	position := m.k * ridx
	if len(row)+start > m.k {
		panic("row cannot be longer than" + strconv.Itoa(m.k))
	}
	copy(m.chunks[position+start:], row)
}

// AddColumn adds a column to the matrix. Argumnets can't be modified after passing to matrix.
func (m *Matrix) AddColumn(cidx, start int, column [][]byte) {
	if len(column)+start > m.k {
		panic("column cannot be longer than" + strconv.Itoa(m.k))
	}
	for i := range column {
		position := m.k*i + cidx + start
		m.chunks[position] = column[i]
	}
}

// GetColumn appends cidx columns entries to dst.
func (m *Matrix) GetColumn(dst [][]byte, cidx int) [][]byte {
	if dst == nil {
		dst = make([][]byte, 0, m.k)
	}
	for i := 0; i < m.k; i++ {
		position := cidx + m.k*i
		dst = append(dst, m.chunks[position])
	}
	return dst
}

func (m *Matrix) GetRow(dst [][]byte, ridx int) [][]byte {
	if dst == nil {
		dst = make([][]byte, 0, m.k)
	}
	for i := 0; i < m.k; i++ {
		position := ridx*m.k + i
		dst = append(dst, m.chunks[position])
	}
	return dst
}

func (m *Matrix) PrettyString(hex bool) string {
	buf := bytes.NewBuffer(nil) // TODO estimate size
	j := 0
	format := "%s"
	if hex {
		format = "%x"
	}
	last := len(m.chunks) - 1
	for i := range m.chunks {
		toPrint := m.chunks[i]
		if toPrint == nil {
			toPrint = empty
		} else {
			// align all values in the matrix so that it looks better
			if len(toPrint) > len(empty) {
				toPrint = toPrint[:len(empty)]
			}
		}
		fmt.Fprintf(buf, format, toPrint)
		j++
		if j == m.k && i != last {
			buf.WriteString("\n")
			j = 0
		} else if j != m.k {
			buf.WriteString(" ")
		}
	}
	return buf.String()
}

func (m *Matrix) RowSize() int {
	return m.k
}

func (m *Matrix) ColumnSize() int {
	return m.RowSize()
}

func (m *Matrix) Size() int {
	return len(m.chunks)
}

func (m *Matrix) Get(i, j int) []byte {
	return m.chunks[i*m.k+j]
}
