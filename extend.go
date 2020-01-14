package fraud

import "github.com/klauspost/reedsolomon"

func NewExtender(chunks, chunkSize int) Extender {
	enc, err := reedsolomon.New(chunks, chunks)
	if err != nil {
		// reedsolomon.New can fail only if invalid options are provided
		panic(err.Error())
	}
	return Extender{
		chunkSize: chunkSize,
		enc:       enc,
	}
}

type Extender struct {
	chunkSize int
	enc       reedsolomon.Encoder
}

func (e Extender) Extend(m *Matrix) error {
	half := m.RowSize() / 2
	if err := e.ExtendHorizontally(m, 0, half); err != nil {
		return err
	}
	if err := e.ExtendVertically(m, 0, half); err != nil {
		return err
	}
	if err := e.ExtendHorizontally(m, half, half*2); err != nil {
		return err
	}
	return nil
}

func (e Extender) ExtendHorizontally(m *Matrix, lo, hi int) error {
	for i := lo; i < hi; i++ {
		row := m.GetRow(nil, i)
		for j := range row {
			if row[j] == nil {
				row[j] = make([]byte, e.chunkSize)
			}
		}
		if err := e.enc.Encode(row); err != nil {
			return err
		}
		m.AddRow(i, 0, row)
	}
	return nil
}

func (e Extender) ExtendVertically(m *Matrix, lo, hi int) error {
	for i := lo; i < hi; i++ {
		column := m.GetColumn(nil, i)
		for j := range column {
			if column[j] == nil {
				column[j] = make([]byte, e.chunkSize)
			}
		}
		if err := e.enc.Encode(column); err != nil {
			return err
		}
		m.AddColumn(i, 0, column)
	}
	return nil
}
