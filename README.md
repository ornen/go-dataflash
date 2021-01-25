# go-dataflash

This is a Dataflash log parser for Go.

## Example
The following example parses a dataflash log file and prints the messages.
```go
package main

import (
	"io"
	"log"
	"os"

	"github.com/ornen/go-dataflash"
)

func main() {
	input, err := os.Open("test.bin")

	if err != nil {
		log.Fatal(err)
	}

	defer input.Close()

	var reader = dataflash.NewReader(input)

	for {
		var message, err = reader.Read()

		if err != nil {
			if err == io.EOF {
				break
			} else {
				log.Println(err)
				continue
			}
		}

		log.Println(message)
	}
}
```

## License

This code is licensed under the Apache License 2.0.
