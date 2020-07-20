package ini

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"reflect"
	"strconv"
	"strings"
)

// LoadIni .
func LoadIni(reader io.Reader, configData interface{}) error {
	scanner := bufio.NewScanner(reader)
	lineno := 0
	kvs := make(map[string]string)

	for scanner.Scan() {
		line := scanner.Text()
		lineno++
		if len(line) == 0 {
			continue
		}

		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}

		if line[0] == '#' {
			continue
		}

		pos := strings.Index(line, "=")
		if pos == -1 {
			text := fmt.Sprintf("invalid line %d", lineno)
			return errors.New(text)
		}

		key := strings.TrimSpace(line[:pos])
		value := strings.TrimSpace(line[pos+1:])

		kvs[key] = value
	}

	structVal := reflect.ValueOf(configData).Elem()
	rt := structVal.Type()
	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		fieldName := field.Tag.Get("ini")
		if len(fieldName) == 0 {
			continue
		}

		value, present := kvs[fieldName]
		if !present {
			continue
		}

		switch field.Type.Kind() {
		case reflect.Int:
			fallthrough
		case reflect.Int8:
			fallthrough
		case reflect.Int16:
			fallthrough
		case reflect.Int32:
			fallthrough
		case reflect.Int64:
			x, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				text := fmt.Sprintf("%s=%s invalid integer", fieldName, value)
				return errors.New(text)
			}
			structVal.Field(i).SetInt(x)

		case reflect.Uint8:
			fallthrough
		case reflect.Uint16:
			fallthrough
		case reflect.Uint32:
			fallthrough
		case reflect.Uint64:
			x, err := strconv.ParseUint(value, 10, 64)
			if err != nil {
				text := fmt.Sprintf("%s=%s invalid integer", fieldName, value)
				return errors.New(text)
			}
			structVal.Field(i).SetUint(x)

		case reflect.Float32:
			fallthrough
		case reflect.Float64:
			x, err := strconv.ParseFloat(value, 64)
			if err != nil {
				text := fmt.Sprintf("%s=%s invalid float", fieldName, value)
				return errors.New(text)
			}
			structVal.Field(i).SetFloat(x)

		case reflect.Bool:
			x, err := strconv.ParseBool(value)
			if err != nil {
				text := fmt.Sprintf("%s=%s invalid bool", fieldName, value)
				return errors.New(text)
			}
			structVal.Field(i).SetBool(x)

		case reflect.String:
			structVal.Field(i).SetString(value)

		default:
			text := fmt.Sprintf("%s=%s invalid type", fieldName, value)
			return errors.New(text)
		}
	}

	return nil
}

// LoadIniFromFile .
func LoadIniFromFile(confFile string, configData interface{}) error {
	f, err := os.Open(confFile)
	if err != nil {
		return err
	}
	defer f.Close()

	return LoadIni(f, configData)
}
