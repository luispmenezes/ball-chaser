package parsers

import (
	"github.com/luispmenezes/ball-chaser/internal/bitreader"
	"github.com/luispmenezes/ball-chaser/pkg/ballchaser/model"
)

func ParseHeader(reader *bitreader.Reader) model.Header {
	var header model.Header

	header.Size = reader.ReadInt(32)
	header.CRC = reader.ReadInt(32)
	header.EngineVersion = reader.ReadInt(32)
	header.LicenseeVersion = reader.ReadInt(32)

	if header.EngineVersion >= 868 && header.LicenseeVersion >= 18 {
		header.NetVersion = reader.ReadInt(32)
	}

	header.Label = reader.ReadString(0)

	header.Properties = parseProperties(reader)

	return header
}
