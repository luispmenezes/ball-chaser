package model

import (
	"github.com/pkg/errors"
)

type Property struct {
	Name   string
	Type   string
	Length uint
	Value  interface{}
}

func (p *Property) getIntValue() (uint, error) {
	if p.Type != "IntProperty" && p.Type != "QWordProperty" {
		return 0, errors.New(p.Name + " has unexpected type")
	}
	return p.Value.(uint), nil
}

func (p *Property) getStringValue() (string, error) {
	if p.Type != "StrProperty" && p.Type != "NameProperty" {
		return "", errors.New(p.Name + " has unexpected type")
	}
	return p.Value.(string), nil
}

func (p *Property) getFloatValue() (float32, error) {
	if p.Type != "FloatProperty" {
		return 0.0, errors.New(p.Name + " has unexpected type")
	}
	return p.Value.(float32), nil
}

func (p *Property) getByteValue() (map[string]string, error) {
	if p.Type != "ByteProperty" {
		return map[string]string{}, errors.New(p.Name + " has unexpected type")
	}
	return p.Value.(map[string]string), nil
}

func (p *Property) getBoolValue() (bool, error) {
	if p.Type != "BoolProperty" {
		return false, errors.New(p.Name + " has unexpected type")
	}
	return p.Value.([]uint8)[0] == 1, nil
}

func (p *Property) getArrayValue() ([]map[string]Property, error) {
	if p.Type != "BoolProperty" {
		return []map[string]Property{}, errors.New(p.Name + " has unexpected type")
	}
	return p.Value.([]map[string]Property), nil
}
