package mhttp

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"
)

const StringSeparator = "="
const NotStringSeparator = ":="

type ParamsParser struct {
	data map[string]interface{}
}

func ParseParams(params []string) (ParamsParser, error) {
	parser := ParamsParser{}

	result, err := parser.parse(params)

	if err != nil {
		return parser, err
	}

	parser.data = result

	return parser, nil
}

func (p *ParamsParser) parse(params []string) (map[string]interface{}, error) {
	response := make(map[string]interface{})

	for _, param := range params {
		var separator string

		// hack for "x=blabla:=1"
		if strings.Contains(param[:strings.Index(param, "=")+1], NotStringSeparator) {
			separator = ":="
		} else {
			separator = "="
		}

		items := strings.Split(param, separator)
		key := items[0]
		value := items[1]

		var finalValue interface{} = value

		if len(items) < 1 {
			return response, errors.New("params is invalid")
		}

		properties := strings.Split(key, ".")

		current := response
		lastKey := ""

		for i, property := range properties {
			if _, ok := current[property]; i < len(properties)-1 {
				if !ok {
					empty := make(map[string]interface{})

					current[property] = empty

					current = empty
				} else {
					current = current[property].(map[string]interface{})
				}
			}

			lastKey = property
		}

		if separator == NotStringSeparator {
			switch value {
			case "true":
				finalValue = true
			case "false":
				finalValue = false
			default:
				number, err := strconv.Atoi(value)

				if err != nil {
					return response, err
				}

				finalValue = number
			}
		}

		current[lastKey] = finalValue
	}

	return response, nil
}

func (p ParamsParser) ToJSON() ([]byte, error) {
	return json.Marshal(p.data)
}
