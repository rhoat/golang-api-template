package config

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

type OtelDestination int8

var (
	ErrInvalidDestination = errors.New("invalid OtelDestination value: Valid values are stdout,http,grpc")
)

const (
	STDOUT OtelDestination = iota + 1
	HTTP
	GRPC
)

var otelDestinationToString = map[OtelDestination]string{
	STDOUT: "stdout",
	HTTP:   "http",
	GRPC:   "grpc",
}

// stringToOtelDestination maps strings to their OtelDestination values.
var stringToOtelDestination = map[string]OtelDestination{
	"stdout": STDOUT,
	"http":   HTTP,
	"grpc":   GRPC,
}

// String implements the fmt.Stringer interface.
func (d *OtelDestination) String() string {
	return otelDestinationToString[*d]
}

func otelDestinationDecodeHook(from reflect.Type, to reflect.Type, data interface{}) (interface{}, error) {
	if from.Kind() == reflect.String && to == reflect.TypeOf(OtelDestination(0)) {
		if val, ok := stringToOtelDestination[strings.ToLower(data.(string))]; ok {
			return val, nil
		}
		return nil, fmt.Errorf("invalid OtelDestination: %s", data)
	}
	return data, nil
}
