package dataflash

import (
	"encoding/binary"
	"errors"
	"io"
)

var (
	// ErrMissingHeader means that the magic header was not found when trying to read a message.
	ErrMissingHeader = errors.New("missing magic message header")
	// ErrUnknownMessageType means that a message with an unknown type definition was found.
	ErrUnknownMessageType = errors.New("unknown message type")
)

var (
	endian = binary.LittleEndian
)

// NewReader instantiates and returns a dataflash message reader.
func NewReader(input io.Reader) *Reader {
	return &Reader{
		input:   input,
		formats: map[uint8]Format{0x80: NewFormat("FMT", 89, "BBnNZ", "Type,Length,Name,Format,Columns")},
	}
}

type Reader struct {
	input   io.Reader
	formats map[uint8]Format
}

// https://docs.python.org/3/library/struct.html
// https://github.com/ArduPilot/pymavlink/blob/master/DFReader.py

func (r *Reader) Read() (*Message, error) {
	header := make([]byte, 3)

	n, err := r.input.Read(header)

	if err != nil {
		return nil, err
	}

	if header[0] != Head1 || header[1] != Head2 || n < 3 {
		return nil, ErrMissingHeader
	}

	messageType := header[2]

	format, known := r.formats[messageType]

	//log.Println(format)
	if !known {
		return nil, ErrUnknownMessageType
	}

	var buffer []byte

	elements := map[string]interface{}{}

	for i, t := range format.types {
		column := format.columns[i]

		switch t {
		case "a":
			buffer = make([]byte, 64)
			r.input.Read(buffer)
			elements[column] = nullts(buffer)
		case "b":
			var value int8
			binary.Read(r.input, endian, &value)
			elements[column] = value
		case "B":
			var value uint8
			binary.Read(r.input, endian, &value)
			elements[column] = value
		case "c":
			var value int16
			binary.Read(r.input, endian, &value)
			elements[column] = float32(value) * 0.01
		case "C":
			var value uint16
			binary.Read(r.input, endian, &value)
			elements[column] = float32(value) * 0.01
		case "e":
			var value int32
			binary.Read(r.input, endian, &value)
			elements[column] = float32(value) * 0.01
		case "d":
			var value int32
			binary.Read(r.input, endian, &value)
			elements[column] = float32(value) * 0.01
		case "E":
			var value uint32
			binary.Read(r.input, endian, &value)
			elements[column] = float32(value) * 0.01
		case "f":
			var value float32
			binary.Read(r.input, endian, &value)
			elements[column] = value
		case "h":
			var value int16
			binary.Read(r.input, endian, &value)
			elements[column] = value
		case "H":
			var value uint16
			binary.Read(r.input, endian, &value)
			elements[column] = value
		case "i":
			var value int32
			binary.Read(r.input, endian, &value)
			elements[column] = value
		case "I":
			var value uint32
			binary.Read(r.input, endian, &value)
			elements[column] = value
		case "L":
			var value int32
			binary.Read(r.input, endian, &value)
			elements[column] = float32(value) * 1.0e-7
		case "M":
			var value int8
			binary.Read(r.input, endian, &value)
			elements[column] = value
		case "n":
			buffer = make([]byte, 4)
			r.input.Read(buffer)
			elements[column] = nullts(buffer)
		case "N":
			buffer = make([]byte, 16)
			r.input.Read(buffer)
			elements[column] = nullts(buffer)
		case "Z":
			buffer = make([]byte, 64)
			r.input.Read(buffer)
			elements[column] = nullts(buffer)
		case "q":
			var value int64
			binary.Read(r.input, endian, &value)
			elements[column] = value
		case "Q":
			var value uint64
			binary.Read(r.input, endian, &value)
			elements[column] = value
		}
	}

	if messageType == 128 {
		r.formats[elements["Type"].(uint8)] = NewFormat(elements["Name"].(string), elements["Length"].(uint8), elements["Format"].(string), elements["Columns"].(string))
	}

	return nil, nil
}

func nullts(buffer []byte) string {
	return string(buffer[:clen(buffer)])
}

func clen(n []byte) int {
	for i := 0; i < len(n); i++ {
		if n[i] == 0 {
			return i
		}
	}
	return len(n)
}
