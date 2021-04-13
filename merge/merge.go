package merge

import (
	"fmt"
	"reflect"

	"gopkg.in/yaml.v2"
)

// Merge the given YAML files into a single YAML file
func Yaml(first, second []byte, rest ...[]byte) ([]byte, error) {
	files := append([][]byte{first, second}, rest...)
	return innerYaml(files...)
}

func innerYaml(files ...[]byte) ([]byte, error) {
	var (
		result map[string]interface{}
		err    error
	)

	for _, file := range files {
		var override map[string]interface{}
		if err := yaml.Unmarshal(file, &override); err != nil {
			return nil, err
		}

		// this will only happen for the first loop
		if result == nil {
			result = override
			continue
		}

		for key, value := range override {
			result[key], err = mergeUnmarshalled(result[key], value)
			if err != nil {
				return nil, err
			}
		}

	}

	bytes, err := yaml.Marshal(result)
	if err != nil {
		return nil, fmt.Errorf("could not marshal final YAML: %w", err)
	}

	return bytes, nil
}

// Merge two marshalled YAML values
func mergeUnmarshalled(base, override interface{}) (interface{}, error) {
	// primitive values can't be merged, and override element takes precedence
	// this is the base case of our recursive decent into the YAML structure
	if isPrimitive(override) {
		return override, nil
	}

	baseValue, overrideValue := reflect.ValueOf(base), reflect.ValueOf(override)

	// if the base isn't valid/set, no alternative but to return the override
	if !baseValue.IsValid() {
		return override, nil
	}

	// likewise if the override isn't set, we have no alternative
	if !overrideValue.IsValid() {
		return base, nil
	}

	baseType, overrideType := baseValue.Type(), overrideValue.Type()

	// sanity check: only equivalent types of fields can be merged
	if baseType.Kind() != overrideType.Kind() {
		panic(fmt.Sprintf("mismatching kinds: %s and %s", baseType.Kind(), overrideType.Kind()))
	}

	// concatenate slices
	if baseType.Kind() == reflect.Slice {
		var out []interface{}
		for i := 0; i < baseValue.Len(); i++ {
			out = append(out, baseValue.Index(i).Interface())
		}
		for i := 0; i < overrideValue.Len(); i++ {
			out = append(out, overrideValue.Index(i).Interface())
		}
		return out, nil
	}

	// merge maps, everything else override
	if baseType.Kind() != reflect.Map {
		return override, nil
	}

	var err error
	out := make(map[string]interface{})

	// loop through all values of the map, merging individual items
	for _, key := range valueKeys(baseValue, overrideValue) {

		getInterfaceValue := func(value reflect.Value) interface{} {

			if !(value.IsValid() && !value.IsZero() && value.CanInterface()) {
				return nil
			}
			return value.Interface()
		}

		baseValueAtKey := baseValue.MapIndex(key)
		baseAtKey := getInterfaceValue(baseValueAtKey)

		overrideValueAtKey := overrideValue.MapIndex(key)
		overrideAtKey := getInterfaceValue(overrideValueAtKey)

		// recursively descend into each node
		out[key.Interface().(string)], err = mergeUnmarshalled(baseAtKey, overrideAtKey)
		if err != nil {
			return nil, err
		}
	}
	return out, nil
}

// extract all unique keys of the given map reflect.Values
func valueKeys(first, second reflect.Value) []reflect.Value {
	firstKeys := first.MapKeys()
	for _, key := range second.MapKeys() {
		var hasKey bool
		for _, firstKey := range firstKeys {
			if firstKey.Interface() == key.Interface() {
				hasKey = true
				break
			}
		}

		if !hasKey {
			firstKeys = append(firstKeys, key)
		}
	}
	return firstKeys
}

func isPrimitive(value interface{}) bool {
	switch value.(type) {
	case bool, string,
		int, int8, int16, int32, int64,
		uint, uint8, uint16, uint32, uint64,
		float32, float64:
		return true
	default:
		return false
	}
}
