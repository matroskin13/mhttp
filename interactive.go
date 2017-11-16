package mhttp

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/hokaccha/go-prettyjson"
)

const (
	CommandGet = "get"
)

var commandNotFound = errors.New("command not found")

func getInMap(props []string, b []byte) ([]byte, error) {
	current := b

	for _, prop := range props {
		data := map[string]json.RawMessage{}

		err := json.Unmarshal(current, &data)

		if err != nil {
			return current, err
		}

		if value, ok := data[prop]; ok {
			current = value
		} else {
			err = errors.New("property is not defined")

			return current, errors.New("property is not defined")
		}
	}

	return current, nil
}

func execute(in string, b []byte) error {
	values := strings.Split(strings.TrimSpace(in), " ")

	if len(values) == 0 {
		return commandNotFound
	}

	command := values[0]

	switch command {
	case CommandGet:
		if len(values) < 2 {
			return errors.New("property name is not defined")
		}

		properties := strings.Split(values[1], ".")

		b, err := getInMap(properties, b)

		if err != nil {
			return err
		}

		s, _ := prettyjson.Format(b)

		fmt.Println(string(s))
	}

	return nil
}

func InitInteractive(b []byte) error {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("Enter command: ")
		text, _ := reader.ReadString('\n')

		err := execute(text, b)

		if err != nil {
			return err
		}
	}

	return nil
}
