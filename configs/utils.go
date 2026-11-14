package configs

import "fmt"

func SetByTag(obj any, key, value string) error {
	config, ok := obj.(*Config)
	if !ok || config == nil {
		return fmt.Errorf("expected *Config")
	}
	field, ok := Fields(config)[key]
	if !ok {
		return fmt.Errorf("unknown configuration key %q", key)
	}
	*field = value
	return nil
}
