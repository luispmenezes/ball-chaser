package header

import "github.com/luispmenezes/ball-chaser/internal/bitreader"

func ParseHeader(reader *bitreader.Reader) Header {
	var header Header

	header.Size = reader.ReadInt32()
	header.CRC = reader.ReadInt32()
	header.EngineVersion = reader.ReadInt32()
	header.LicenseeVersion = reader.ReadInt32()

	if header.EngineVersion >= 868 && header.LicenseeVersion >= 18 {
		header.NetVersion = reader.ReadInt32()
	}

	header.Label = reader.ReadString(0)

	header.Properties = parseProperties(reader)

	return header
}
