package dataflash

import "strings"

type Format struct {
	name    string
	length  uint8
	types   []string
	columns []string
}

func NewFormat(name string, length uint8, types string, columns string) Format {
	return Format{
		name:    name,
		length:  length,
		types:   strings.Split(types, ""),
		columns: strings.Split(columns, ","),
	}
}
