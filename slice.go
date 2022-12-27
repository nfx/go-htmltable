package htmltable

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"
)

// NewSlice returns slice of annotated struct types from io.Reader
func NewSlice[T any](ctx context.Context, r io.Reader) ([]T, error) {
	f := &feeder[T]{
		Page: Page{ctx: ctx},
	}
	f.init(r)
	return f.slice()
}

// NewSliceFromPage finds a table matching the slice and returns the slice
func NewSliceFromPage[T any](p *Page) ([]T, error) {
	return (&feeder[T]{
		Page: *p,
	}).slice()
}

// NewSliceFromString is same as NewSlice(context.Context, io.Reader),
// but takes just a string.
func NewSliceFromString[T any](in string) ([]T, error) {
	return NewSlice[T](context.Background(), strings.NewReader(in))
}

// NewSliceFromString is same as NewSlice(context.Context, io.Reader),
// but takes just an http.Response
func NewSliceFromResponse[T any](resp *http.Response) ([]T, error) {
	return NewSlice[T](resp.Request.Context(), resp.Body)
}

// NewSliceFromString is same as NewSlice(context.Context, io.Reader),
// but takes just an URL.
func NewSliceFromURL[T any](url string) ([]T, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.Body != nil {
		defer resp.Body.Close()
	}
	return NewSliceFromResponse[T](resp)
}

type feeder[T any] struct {
	Page

	dummy T
}

func (f *feeder[T]) headers() ([]string, map[string]int, error) {
	dt := reflect.ValueOf(f.dummy)
	elem := dt.Type()
	headers := []string{}
	fields := map[string]int{}
	for i := 0; i < elem.NumField(); i++ {
		field := elem.Field(i)
		header := field.Tag.Get("header")
		if header == "" {
			continue
		}
		err := f.isTypeSupported(field)
		if err != nil {
			return nil, nil, err
		}
		fields[header] = i
		headers = append(headers, header)
	}
	return headers, fields, nil
}

func (f *feeder[T]) isTypeSupported(field reflect.StructField) error {
	k := field.Type.Kind()
	if k == reflect.String {
		return nil
	}
	if k == reflect.Int {
		return nil
	}
	if k == reflect.Bool {
		return nil
	}
	return fmt.Errorf("setting field is not supported, %s is %v",
		field.Name, field.Type.Name())
}

func (f *feeder[T]) table() (*Table, map[int]int, error) {
	headers, fields, err := f.headers()
	if err != nil {
		return nil, nil, err
	}
	table, err := f.FindWithColumns(headers...)
	if err != nil {
		return nil, nil, err
	}
	mapping := map[int]int{}
	for idx, header := range table.Header {
		field, ok := fields[header]
		if !ok {
			continue
		}
		mapping[idx] = field
	}
	return table, mapping, nil
}

func (f *feeder[T]) slice() ([]T, error) {
	table, mapping, err := f.table()
	if err != nil {
		return nil, err
	}
	dummy := reflect.ValueOf(f.dummy)
	dt := dummy.Type()
	sliceValue := reflect.MakeSlice(reflect.SliceOf(dt),
		len(table.Rows), len(table.Rows))
	for rowIdx, row := range table.Rows {
		item := sliceValue.Index(rowIdx)
		for idx, field := range mapping {
			if len(row) < len(mapping) && idx == len(row) {
				// either corrupt row or something like that
				continue
			}
			switch item.Field(field).Kind() {
			case reflect.String:
				item.Field(field).SetString(row[idx])
			case reflect.Bool:
				var v bool
				lower := strings.ToLower(row[idx])
				if lower == "yes" ||
					lower == "y" ||
					lower == "true" ||
					lower == "t" {
					v = true
				}
				item.Field(field).SetBool(v)
			case reflect.Int:
				var v int64
				_, err := fmt.Sscan(row[idx], &v)
				if err != nil {
					column := table.Header[idx]
					return nil, fmt.Errorf("row %d: %s: %w", rowIdx, column, err)
				}
				item.Field(field).SetInt(v)
			default: // noop
			}
		}
	}
	return sliceValue.Interface().([]T), nil
}
