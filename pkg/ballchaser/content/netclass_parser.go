package content

import "github.com/luispmenezes/ball-chaser/internal/bitreader"

func parseNetCache(reader *bitreader.Reader, classIndex map[uint]string) map[string]NetClass {
	classNetCacheLength := reader.ReadInt32()
	classes := make(map[string]NetClass)

	for i := 0; i < int(classNetCacheLength); i++ {
		var netClass NetClass
		netClass.Index = reader.ReadInt32()
		netClass.ParentId = reader.ReadInt32()
		netClass.Id = reader.ReadInt32()

		propertiesLength := reader.ReadInt32()
		netClass.Properties = make(map[uint]uint)

		for j := 0; j < int(propertiesLength); j++ {
			propIndex := reader.ReadInt32()
			propId := reader.ReadInt32()
			netClass.Properties[propId] = propIndex
		}

		classes[classIndex[netClass.Index]] = netClass
	}

	return classes
}
