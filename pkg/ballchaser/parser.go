package ballchaser

import (
	"ball-chaser/internal/bitreader"
	"ball-chaser/pkg/ballchaser/model"
	"ball-chaser/pkg/ballchaser/parsers"
	"io"
)

type parser struct {
	bitReader bitreader.Reader
	Header    model.Header
}

func NewParser(replayStream io.Reader) parser {
	p := parser{
		bitReader: *bitreader.NewLargeBitReader(replayStream),
	}
	p.Header = parsers.ParseHeader(&p.bitReader)
	return p
}
