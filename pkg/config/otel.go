package config

import (
	"errors"
	"fmt"
	"reflect"

	gotel "github.com/rhoat/go-exercise/pkg/otel"
)

var (
	ErrInvalidDestination = errors.New("invalid OtelDestination value: Valid values are stdout,http,grpc")
)

func otelDestinationDecodeHook(from reflect.Type, to reflect.Type, data interface{}) (interface{}, error) {
	if from.Kind() == reflect.String && to == reflect.TypeOf(gotel.Destination(0)) {
		stringData, ok := data.(string)
		if !ok {
			return nil, fmt.Errorf("invalid OtelDestination: unable to cast to a string: %s", data)
		}
		destination, err := gotel.StringToDestination(stringData)
		if err != nil {
			return nil, fmt.Errorf("invalid OtelDestination: %s, error: %w", data, err)
		}
		return destination, nil
	}
	return data, nil
}
