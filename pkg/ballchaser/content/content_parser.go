package content

import "github.com/luispmenezes/ball-chaser/internal/bitreader"

func ParseContent(reader *bitreader.Reader, netVersion uint) Content {
	var content Content

	content.Length = reader.ReadInt32()
	content.CRC = reader.ReadInt32()

	levelLength := reader.ReadInt32()
	for i := 0; i < int(levelLength); i++ {
		content.Levels = append(content.Levels, reader.ReadString(0))
	}

	keyFrameLength := reader.ReadInt32()
	for i := 0; i < int(keyFrameLength); i++ {
		keyFrame := KeyFrame{Time: reader.ReadFloat(), Frame: reader.ReadInt32(), FilePosition: reader.ReadInt32()}
		content.KeyFrames = append(content.KeyFrames, keyFrame)
	}

	networkStreamLength := reader.ReadInt32()
	content.NetworkStream = reader.ReadBytes(int(networkStreamLength))

	debugStringLength := reader.ReadInt32()
	for i := 0; i < int(debugStringLength); i++ {
		debugString := DebugString{Frame: reader.ReadInt32(), Username: reader.ReadString(0), Text: reader.ReadString(0)}
		content.DebugStrings = append(content.DebugStrings, debugString)
	}

	tickMarkLength := reader.ReadInt32()
	for i := 0; i < int(tickMarkLength); i++ {
		tickMark := Tick{Type: reader.ReadString(0), Frame: reader.ReadInt32()}
		content.Ticks = append(content.Ticks, tickMark)
	}

	packagesLength := reader.ReadInt32()
	for i := 0; i < int(packagesLength); i++ {
		content.Packages = append(content.Packages, reader.ReadString(0))
	}

	objectsLength := reader.ReadInt32()
	for i := 0; i < int(objectsLength); i++ {
		content.Objects = append(content.Objects, reader.ReadString(0))
	}

	namesLength := reader.ReadInt32()
	for i := 0; i < int(namesLength); i++ {
		content.Names = append(content.Names, reader.ReadString(0))
	}

	classIndexLength := reader.ReadInt32()
	content.ClassIndex = make(map[string]uint)
	reverseIndex := make(map[uint]string)
	for i := 0; i < int(classIndexLength); i++ {
		class := reader.ReadString(0)
		index := reader.ReadInt32()
		content.ClassIndex[class] = index
		reverseIndex[index] = class
	}

	content.NetClass = parseNetCache(reader, reverseIndex)

	if netVersion >= 10 {
		reader.ReadInt32()
	}

	return content
}
