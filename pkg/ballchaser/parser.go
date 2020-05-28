package ballchaser

import (
	"github.com/luispmenezes/ball-chaser/internal/bitreader"
	"github.com/luispmenezes/ball-chaser/pkg/ballchaser/content"
	"github.com/luispmenezes/ball-chaser/pkg/ballchaser/header"
	"io"
)

type parser struct {
	bitReader bitreader.Reader
	Header    header.Header
	Content   content.Content
}

func NewParser(replayStream io.Reader) parser {
	p := parser{
		bitReader: *bitreader.NewLargeBitReader(replayStream),
	}
	p.Header = header.ParseHeader(&p.bitReader)
	return p
}

func (p *parser) ParseContent() {
	p.Content = content.ParseContent(&p.bitReader, p.Header.NetVersion)
}
