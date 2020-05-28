package header

import (
	"fmt"
	"github.com/luispmenezes/ball-chaser/internal/bitreader"
	"github.com/pkg/errors"
	"log"
)

func parseProperties(reader *bitreader.Reader) map[string]Property {
	propertyMap := make(map[string]Property)

	for {
		prop, err := parseSingleProp(reader)

		if err != nil {
			log.Fatal(err)
		}

		if prop.Name == "None" {
			break
		}

		propertyMap[prop.Name] = prop
	}

	return propertyMap
}

func parseSingleProp(reader *bitreader.Reader) (Property, error) {
	var prop Property

	prop.Name = reader.ReadString(0)

	if prop.Name == "None" {
		return prop, nil
	}

	prop.Type = reader.ReadString(0)

	prop.Length = reader.ReadInt32()
	reader.ReadInt32()

	if prop.Type == "IntProperty" {
		prop.Value = reader.ReadInt32()
	} else if prop.Type == "StrProperty" || prop.Type == "NameProperty" {
		prop.Value = reader.ReadString(0)
	} else if prop.Type == "FloatProperty" {
		prop.Value = reader.ReadFloat()
	} else if prop.Type == "ByteProperty" {
		prop.Value = map[string]string{reader.ReadString(0): reader.ReadString(0)}
	} else if prop.Type == "BoolProperty" {
		prop.Value = reader.ReadBytes(1)
	} else if prop.Type == "QWordProperty" {
		prop.Value = reader.ReadInt(64)
	} else if prop.Type == "ArrayProperty" {
		var array []map[string]Property
		arrayLength := reader.ReadInt32()
		for i := 0; i < int(arrayLength); i++ {
			arrayElement := make(map[string]Property)
			for {
				arrayElementProp, err := parseSingleProp(reader)
				if err != nil {
					return Property{}, errors.Wrap(err, fmt.Sprintf("failed parsing array property: %s", prop.Name))
				}
				if arrayElementProp.Name == "None" {
					break
				} else {
					arrayElement[arrayElementProp.Name] = arrayElementProp
				}
			}

			array = append(array, arrayElement)
		}
		prop.Value = array
		return prop, nil
	} else {
		return prop, errors.New(fmt.Sprintf("unknown prop: %s", prop.Type))
	}

	return prop, nil
}
