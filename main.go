package main

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"io"
	"os"
)

func main() {
	help := hasArg("-h") || hasArg("--help")
	d := hasArg("-d")

	if help {
		printHelp()
		return
	}

	if d {
		err := decode()
		if err != nil {
			fmt.Printf("ERROR: %s\r\n", err)
			os.Exit(1)
		}
	} else {
		err := encode()
		if err != nil {
			fmt.Printf("ERROR: %s\r\n", err)
			os.Exit(1)
		}
	}
}

func hasArg(arg string) bool {
	for _, a := range os.Args {
		if a == arg {
			return true
		}
	}

	return false
}

func decode() error {
	decoder := base64.NewDecoder(base64.StdEncoding, os.Stdin)

	buffer := make([]byte, 1000)
	for {
		count, err := decoder.Read(buffer)
		if err != nil {
			if count == 0 && err == io.EOF {
				// end
				return nil
			}

			if count == 0 && err != nil {
				return err
			}
		}

		_, err = os.Stdout.Write(buffer[0:count])
		if err != nil {
			return err
		}
	}
}

func encode() error {
	encoder := base64.NewEncoder(base64.StdEncoding, os.Stdout)
	defer encoder.Close()
	reader := bufio.NewReader(os.Stdin)

	buffer := make([]byte, 1000)
	for {
		count, err := reader.Read(buffer)
		if err != nil {
			if count == 0 && err == io.EOF {
				// end
				return nil
			}

			if count == 0 && err != nil {
				return err
			}
		}

		written, err := encoder.Write(buffer[0:count])
		if written != count {
			return fmt.Errorf("Could not write all bytes")
		}
		if err != nil {
			return err
		}
	}
}

func printHelp() {
	fmt.Printf("Read text from stdin and output base64 to stdout.\r\n")
	fmt.Printf("Usage: gobase64 [-d]\r\n")
	fmt.Printf("\r\n")
	fmt.Printf("  -d Reverse base 64 encoding, read base64 from stdin print text to stdout.\r\n")
}
