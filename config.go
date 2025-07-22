package yourconfig

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/ettle/strcase"
)

func MustLoadContext[T any](ctx context.Context) T {
	output, err := LoadContext[T](ctx)
	if err != nil {
		panic(fmt.Sprintf("must load: %s", err.Error()))
	}

	return output
}

func MustLoad[T any]() T {
	output, err := LoadContext[T](context.Background())
	if err != nil {
		panic(fmt.Sprintf("must load: %s", err.Error()))
	}

	return output
}

func Load[T any]() (T, error) {
	return LoadContext[T](context.Background())
}

func LoadContext[T any](ctx context.Context) (T, error) {
	var cfg T

	v := reflect.ValueOf(&cfg).Elem()
	t := v.Type()

	errs := make([]error, 0)

OUTER:
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tagStr := field.Tag.Get("cfg")
		if tagStr == "" {
			continue
		}

		singleValue := make([]string, 0)
		options := make(map[string]string, 0)

		values := strings.SplitSeq(tagStr, ",")
		for value := range values {
			if value == "" {
				continue
			}

			key, val, ok := strings.Cut(strings.TrimSpace(value), ":")
			if ok {
				options[key] = val
			} else {
				singleValue = append(singleValue, value)
			}
		}

		var tag tag
		if len(singleValue) == 0 {
			tag.Env = strcase.ToSNAKE(field.Name)
		} else {
			tag.Env = singleValue[0] // We always count the first value as the name, if set
		}

		for _, option := range singleValue {
			switch option {
			case "required":
				tag.Required = true
			}
		}

		for key, val := range options {
			switch key {
			case "required":
				required, err := strconv.ParseBool(val)
				if err != nil {
					errs = append(errs, fmt.Errorf("field: %s (key: %s), err: %w", field.Name, key, err))
					continue OUTER
				}

				tag.Required = required
			}
		}

		valueStr, err := defaultLogger.Load().Get(ctx, tag.Env)
		if err != nil {
			errs = append(errs, fmt.Errorf("field: %s failed to load: %w", field.Name, err))
			continue OUTER
		}
		if valueStr == "" && tag.Required {
			errs = append(errs, fmt.Errorf("field: %s (env=%s) is not set and is required", field.Name, tag.Env))
			continue OUTER
		}

		fieldValue := v.Field(i)
		if !fieldValue.CanSet() {
			errs = append(errs, fmt.Errorf("field: %s is not settable", field.Name))
			continue OUTER
		}

		fieldValue.SetString(valueStr)
	}

	if err := errors.Join(errs...); err != nil {
		return cfg, fmt.Errorf("config failed: %w", err)
	}

	return cfg, nil
}

type tag struct {
	Env      string
	Required bool
}
