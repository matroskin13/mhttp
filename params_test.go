package mhttp

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func getParamsJSON(params []string) string {
	parser, _ := ParseParams(params)
	jsonBytes, _ := parser.ToJSON()

	return string(jsonBytes)
}

func TestParseParams(t *testing.T) {
	assert.Equal(t, "{\"x\":\"hey\"}", getParamsJSON([]string{"x=hey"}), "single param")
	assert.Equal(t, "{\"x\":\"hey\",\"y\":1}", getParamsJSON([]string{"x=hey", "y:=1"}), "string and number params")
	assert.Equal(t, "{\"x\":\"hey\",\"y\":true}", getParamsJSON([]string{"x=hey", "y:=true"}), "string and boolean params")

	assert.Equal(
		t,
		"{\"x\":\"cool\",\"y\":{\"one\":1,\"two\":\"hello\"}}",
		getParamsJSON([]string{"x=cool", "y.one:=1", "y.two=hello"}),
		"nested params",
	)

	assert.Equal(
		t,
		"{\"x\":{\"y\":{\"z\":{\"a\":1}}}}",
		getParamsJSON([]string{"x.y.z.a:=1"}),
		"super nested",
	)
}
