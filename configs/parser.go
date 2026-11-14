package configs

import (
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"unicode"
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
				fields[meta.Tag.Get("config")] = f.Addr().Interface().(*string)
			}
		}
	}
	walk(reflect.ValueOf(config).Elem())
	return fields
}

// Words parses the declarative shell-like subset. It never evaluates shell code.
func Words(s string) ([]string, error) {
	var words []string
	var word strings.Builder
	var quote rune
	started, escape := false, false
	flush := func() {
		if started {
			words = append(words, word.String())
			word.Reset()
			started = false
		}
	}
	for _, r := range s {
		if escape {
			word.WriteRune(r)
			started = true
			escape = false
			continue
		}
		if r == '\\' && quote != '\'' {
			escape = true
			started = true
			continue
		}
		if quote != 0 {
			if r == quote {
				quote = 0
				continue
			}
			if quote == '"' && (r == '$' || r == '`') {
				return nil, fmt.Errorf("shell expansion requires --shell-config")
			}
			word.WriteRune(r)
			continue
		}
		switch {
		case r == '\'' || r == '"':
			quote = r
			started = true
		case r == '#' && !started:
			flush()
			return words, nil
		case unicode.IsSpace(r):
			flush()
		case strings.ContainsRune("$`;|&<>(){}", r):
			return nil, fmt.Errorf("shell syntax requires --shell-config")
		default:
			started = true
			word.WriteRune(r)
		}
	}
	if quote != 0 || escape {
		return nil, fmt.Errorf("unterminated quote or escape")
	}
	flush()
	return words, nil
}

func ParseConfig(data string, config *Config) error {
	fields := Fields(config)
	inInfo, seenInfo := false, false
	for i, line := range strings.Split(data, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		fail := func(err error) error { return fmt.Errorf("line %d: %w", i+1, err) }
		if line == "print_info() {" {
			if inInfo || seenInfo {
				return fail(fmt.Errorf("duplicate print_info block"))
			}
			config.InfoList = nil
			inInfo = true
			seenInfo = true
			continue
		}
		if inInfo {
			if line == "}" {
				inInfo = false
				continue
			}
			words, err := Words(line)
			if err != nil {
				return fail(err)
			}
			if len(words) < 2 || len(words) > 3 || words[0] != "info" {
				return fail(fmt.Errorf("expected info [label] module"))
			}
			config.InfoList = append(config.InfoList, words[1:])
			continue
		}
		key, value, ok := strings.Cut(line, "=")
		if !ok {
			return fail(fmt.Errorf("expected key=value"))
		}
		dest, ok := fields[strings.TrimSpace(key)]
		if !ok {
			return fail(fmt.Errorf("unknown configuration key %q", key))
		}
		value = strings.TrimSpace(value)
		array := strings.HasPrefix(value, "(")
		if array {
			if !strings.HasSuffix(value, ")") {
				return fail(fmt.Errorf("unterminated array"))
			}
			value = strings.TrimSpace(value[1 : len(value)-1])
		}
		words, err := Words(value)
		if err != nil {
			return fail(err)
		}
		if len(words) > 1 && !array {
			return fail(fmt.Errorf("scalar value must be quoted"))
		}
		if array && (strings.TrimSpace(key) == "disk_show" || strings.TrimSpace(key) == "local_ip_interface") {
			quoted := make([]string, len(words))
			for i, word := range words {
				quoted[i] = shellQuote(word)
			}
			*dest = "(" + strings.Join(quoted, " ") + ")"
		} else {
			*dest = strings.Join(words, " ")
		}
	}
	if inInfo {
		return fmt.Errorf("unterminated print_info block")
	}
	return Validate(config)
}

func Validate(config *Config) error {
	fields := Fields(config)
	for key, bounds := range map[string][2]int{"mem_precision": {0, 6}, "bar_length": {1, 200}, "block_width": {1, 100}, "block_height": {1, 100}, "gap": {0, 100}, "xoffset": {0, 10000}, "yoffset": {0, 10000}, "public_ip_timeout": {1, 60}} {
		n, err := strconv.Atoi(*fields[key])
		if err != nil || n < bounds[0] || n > bounds[1] {
			return fmt.Errorf("%s must be an integer in %d..%d", key, bounds[0], bounds[1])
		}
	}
	for _, row := range config.InfoList {
		if len(row) < 1 || len(row) > 2 {
			return fmt.Errorf("invalid info row")
		}
	}
	return nil
}

func Format(config *Config) string {
	fields := Fields(config)
	keys := make([]string, 0, len(fields))
	for key := range fields {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	var out strings.Builder
	out.WriteString("# Tintfetch declarative configuration; no shell code is evaluated.\n")
	for _, key := range keys {
		fmt.Fprintf(&out, "%s=%s\n", key, shellQuote(*fields[key]))
	}
	out.WriteString("\nprint_info() {\n")
	for _, row := range config.InfoList {
		out.WriteString("    info")
		for _, word := range row {
			out.WriteString(" " + shellQuote(word))
		}
		out.WriteByte('\n')
	}
	out.WriteString("}\n")
	return out.String()
}

func shellQuote(s string) string { return "'" + strings.ReplaceAll(s, "'", "'\\''") + "'" }
