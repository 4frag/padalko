package utils

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/yarlson/tap"
)

type Number interface {
	~int | ~float32 | ~float64
}

func parseNumber[T Number](input string) (T, error) {
	var zero T

	switch any(zero).(type) {
	case int:
		v, err := strconv.Atoi(input)
		return T(v), err

	case float32:
		v, err := strconv.ParseFloat(input, 32)
		return T(v), err

	case float64:
		v, err := strconv.ParseFloat(input, 64)
		return T(v), err

	default:
		return zero, errors.New("unsupported type")
	}
}


func InputNumber[T Number](
	ctx context.Context,
	prompt string,
	placeholder string,
	defaultValue T,
	validator func(T) error,
) T {
	input := tap.Text(ctx, tap.TextOptions{
		Message:     prompt,
		Placeholder: placeholder,
		Validate: func(input string) error {
			val, err := parseNumber[T](input)
			if err != nil {
				return err
			}

			if validator != nil {
				if err := validator(val); err != nil {
					return err
				}
			}
			return nil
		},
	})

	if input == "" {
		return defaultValue
	}

	val, _ := parseNumber[T](input)
	return val
}

func InputMatrixRow[T Number](
	ctx context.Context,
	prompt string,
	cols int,
	validator func(T) error,
) ([]T, error) {

	input := tap.Text(ctx, tap.TextOptions{
		Message: prompt,
		Placeholder: "Enter values separated by spaces",
		Validate: func(input string) error {
			values := strings.Fields(input)
			if len(values) != cols {
				return fmt.Errorf("expected %d values", cols)
			}

			for _, v := range values {
				val, err := parseNumber[T](v)
				if err != nil {
					return err
				}
				if validator != nil {
					if err := validator(val); err != nil {
						return err
					}
				}
			}
			return nil
		},
	})

	parts := strings.Fields(input)
	row := make([]T, cols)
	for i, p := range parts {
		row[i], _ = parseNumber[T](p)
	}

	return row, nil
}

type MenuReturn interface {
	~int | ~float32 | ~float64 | ~string
}

type MenuItem[T MenuReturn] struct {
    ID          T
    Name        string
    Description string
    Handler     func(ctx context.Context) error
}

func CreateMenu[T MenuReturn](ctx context.Context, message string, items []MenuItem[T]) T {
	options_raw := make([]tap.SelectOption[T], len(items))
	for i := range len(items) {
		options_raw[i] = tap.SelectOption[T]{
			Value: items[i].ID,
			Label: items[i].Name,
			Hint: items[i].Description,
		}
	}
	options := tap.SelectOptions[T]{
		Message: message,
		Options: options_raw,
	}

	return tap.Select(ctx, options)
}