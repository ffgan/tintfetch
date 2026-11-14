package configs

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// Fields returns configuration tags and their addressable string values.
func Fields(config *Config) map[string]*string {
	fields := make(map[string]*string)
	var walk func(reflect.Value)
	walk = func(v reflect.Value) {
		for i := 0; i < v.NumField(); i++ {
			f, meta := v.Field(i), v.Type().Field(i)
			if f.Kind() == reflect.Struct {
				walk(f)
				continue
			}
			if f.Kind() == reflect.String {
				fields[meta.Tag.Get("tintfetch")] = f.Addr().Interface().(*string)
			}
		}
	}
	walk(reflect.ValueOf(config).Elem())
	return fields
}

func Validate(config *Config) error {
	fields := Fields(config)
	for key, bounds := range map[string][2]int{"memory_precision": {0, 6}, "bar_length": {1, 200}, "block_width": {1, 100}, "block_height": {1, 100}, "gap": {0, 100}, "xoffset": {0, 10000}, "yoffset": {0, 10000}, "public_ip_timeout": {1, 60}} {
		value, ok := fields[key]
		if !ok {
			return fmt.Errorf("missing configuration field %s", key)
		}
		n, err := strconv.Atoi(*value)
		if err != nil || n < bounds[0] || n > bounds[1] {
			return fmt.Errorf("%s must be an integer in %d..%d", key, bounds[0], bounds[1])
		}
	}
	for _, row := range config.InfoList {
		if len(row) < 1 || len(row) > 3 {
			return fmt.Errorf("invalid info row")
		}
	}
	return nil
}

func shellQuote(s string) string { return "'" + strings.ReplaceAll(s, "'", "'\\''") + "'" }
