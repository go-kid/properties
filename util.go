package properties

import (
	"encoding/json"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"reflect"
	"strings"
)

func buildMap(path string, val any, tmp map[string]any) map[string]any {
	arr := strings.SplitN(path, ".", 2)
	if len(arr) > 1 {
		key := arr[0]
		next := arr[1]
		if tmp[key] == nil {
			tmp[key] = make(map[string]any)
		}
		switch sub := tmp[key]; sub.(type) {
		case map[string]any, Properties:
			tmp[key] = buildMap(next, val, sub.(map[string]any))
		default:
			tmp[key] = sub
		}
	} else {
		tmp[path] = val
	}
	return tmp
}

func convertAny2Prop(v any) (Properties, error) {
	if pm, ok := v.(Properties); ok {
		return pm, nil
	}
	if m, ok := v.(map[string]any); ok {
		return buildProperties("", m), nil
	}
	m := make(map[string]any)
	config := newDecodeConfig(&m)
	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return nil, err
	}
	err = decoder.Decode(v)
	if err != nil {
		return nil, err
	}
	return buildProperties("", m), nil
}

func buildProperties(prePath string, m map[string]any) Properties {
	tmp := make(Properties)
	for k, v := range m {
		switch v.(type) {
		case map[string]any:
			subTmp := buildProperties(path(prePath, k), v.(map[string]any))
			for subP, subV := range subTmp {
				tmp[subP] = subV
			}
		default:
			tmp[path(prePath, k)] = v
		}
	}
	return tmp
}

func newDecodeConfig(v any) *mapstructure.DecoderConfig {
	return &mapstructure.DecoderConfig{
		DecodeHook:           nil,
		ErrorUnused:          false,
		ErrorUnset:           false,
		ZeroFields:           false,
		WeaklyTypedInput:     true,
		Squash:               false,
		Metadata:             nil,
		Result:               v,
		TagName:              "properties",
		IgnoreUntaggedFields: false,
		MatchName:            nil,
	}
}

func path(first, second string) string {
	if first != "" {
		if second != "" {
			return first + "." + second
		} else {
			return first
		}
	} else {
		return second
	}
}

func formatPropertiesPair(a any) (map[string]any, error) {
	var result = make(map[string]any)
	if a == nil {
		result[""] = "<nil>"
		return result, nil
	}
	switch a.(type) {
	case string:
		result[""] = a
	default:
		switch p := reflect.TypeOf(a); p.Kind() {
		case reflect.Map, reflect.Struct:
			bytes, err := json.Marshal(a)
			if err != nil {
				return nil, err
			}
			var tmp = make(map[string]any)
			err = json.Unmarshal(bytes, &tmp)
			if err != nil {
				return nil, err
			}
			for s, a2 := range buildProperties("", tmp) {
				props, err := formatPropertiesPair(a2)
				if err != nil {
					return nil, err
				}
				for s2, a3 := range props {
					result[fmt.Sprintf("%s%s", s, s2)] = a3
				}
			}
		case reflect.Slice, reflect.Array:
			arrVal := reflect.ValueOf(a)
			for i := 0; i < arrVal.Len(); i++ {
				arrVal.Index(i)
				f, err := formatPropertiesPair(arrVal.Index(i).Interface())
				if err != nil {
					return nil, err
				}
				for s, a2 := range f {
					result[path(arrFormat(i), s)] = a2
				}
			}
		case reflect.Pointer:
			props, err := formatPropertiesPair(reflect.ValueOf(a).Elem().Interface())
			if err != nil {
				return nil, err
			}
			for s, a2 := range props {
				result[s] = a2
			}
		default:
			result[""] = a
		}
	}
	return result, nil
}

func arrFormat(i int) string {
	return fmt.Sprintf("[%d]", i)
}
