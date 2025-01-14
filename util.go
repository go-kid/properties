package properties

import (
	"fmt"
	"github.com/go-kid/strconv2"
	"github.com/mitchellh/mapstructure"
	"reflect"
	"strings"
)

type SetMode uint32

func (m SetMode) Eq(mode SetMode) bool {
	return m&mode > 0
}

const (
	Append SetMode = 1 << iota
	OverwriteType
)

func buildMap(path string, val any, tmp *map[string]any, mode SetMode) {
	rtmp := *tmp
	arr := strings.SplitN(path, ".", 2)
	key, idx, hasIdx := arrSplit(arr[0])
	value, ok := rtmp[key]
	if len(arr) == 2 {
		next := arr[1]
		if !ok {
			if hasIdx {
				value = make([]any, idx)
			} else {
				value = make(map[string]any)
			}
		}
		switch tval := value.(type) {
		case map[string]any:
			buildMap(next, val, &tval, mode)
		case Properties:
			buildMap(next, val, (*map[string]any)(&tval), mode)
		case []any:
			var subM map[string]any
			if len(tval) > idx {
				if am, ok := tval[idx].(map[string]any); ok {
					subM = am
				} else if am != nil {
					if mode.Eq(OverwriteType) {
						subM = make(map[string]any)
					} else {
						panic(fmt.Errorf("%T(%v) is not map[string]interface{}, can't set sub values", am, am))
					}
				}
			}
			if subM == nil {
				subM = make(map[string]any)
			}
			buildMap(next, val, &subM, mode)
			value = setSlice(tval, idx, subM)
		default:
			tmp := make(map[string]any)
			buildMap(next, val, &tmp, mode)
			if mode.Eq(OverwriteType) {
				value = tmp
			} else {
				panic(fmt.Errorf("can't assign %T(%v) to %T(%v)", tmp, tmp, value, value))
			}
		}
	} else {
		if hasIdx {
			if !ok {
				rtmp[key] = setSlice([]any{}, idx, val)
				return
			}
			switch tval := value.(type) {
			case []any:
				value = setSlice(tval, idx, val)
			default:
				tmp := setSlice([]any{}, idx, val)
				if mode.Eq(OverwriteType) {
					value = tmp
				} else {
					panic(fmt.Errorf("can't assign %T(%v) to %T(%v)", tmp, tmp, value, value))
				}
			}
		} else {
			if !ok || !mode.Eq(Append) {
				rtmp[key] = val
				return
			}
			switch tval := value.(type) {
			case []any:
				switch rtval := val.(type) {
				case []any:
					value = append(tval, rtval...)
				default:
					value = append(tval, rtval)
				}
			default:
				switch rtval := val.(type) {
				case []any:
					value = append([]any{tval}, rtval...)
				default:
					value = []any{tval, rtval}
				}
			}
		}
	}
	rtmp[key] = value
}

func getMap(m map[string]any, path string) (any, bool) {
	arr := strings.SplitN(path, ".", 2)
	key, idx, hasIdx := arrSplit(arr[0])
	value, ok := m[key]
	if !ok {
		return nil, false
	}
	if len(arr) == 2 {
		if hasIdx {
			value, ok = getAt(value, idx)
			if !ok {
				return nil, false
			}
		}
		next := arr[1]
		switch value.(type) {
		case map[string]any:
			return getMap(value.(map[string]any), next)
		case Properties:
			return getMap(value.(Properties), next)
		default:
			return nil, false
		}
	} else {
		if hasIdx {
			return getAt(value, idx)
		}
		return value, ok
	}
}

func getAt(value any, idx int) (any, bool) {
	if arr, ok := value.([]any); ok && idx >= 0 && idx < len(arr) {
		return arr[idx], true
	}
	return nil, false
}

func decodeToMap(v any) (map[string]any, error) {
	m := make(map[string]any)
	err := decode(v, &m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func flatten(prePath string, m Properties) map[string]any {
	tmp := make(map[string]any)
	for k, v := range m {
		switch v.(type) {
		case map[string]any:
			subTmp := flatten(path(prePath, k), v.(map[string]any))
			assignMap(subTmp, tmp)
		case Properties:
			subTmp := flatten(path(prePath, k), v.(Properties))
			assignMap(subTmp, tmp)
		default:
			tmp[path(prePath, k)] = v
		}
	}
	return tmp
}

func decode(input any, result any) error {
	config := newDecodeConfig(result)
	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return err
	}
	return decoder.Decode(input)
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

func formatPropertiesPair(key string, a any) (map[string]any, error) {
	var result = make(map[string]any)
	if a == nil {
		result[key] = "<nil>"
		return result, nil
	}
	switch a.(type) {
	case string:
		result[key] = a
	default:
		switch p := reflect.TypeOf(a); p.Kind() {
		case reflect.Map, reflect.Struct:
			tmp, err := decodeToMap(a)
			if err != nil {
				return nil, err
			}
			for s, a2 := range flatten("", tmp) {
				f, err := formatPropertiesPair(path(key, s), a2)
				if err != nil {
					return nil, err
				}
				assignMap(f, result)
			}
		case reflect.Slice, reflect.Array:
			arrVal := reflect.ValueOf(a)
			for i := 0; i < arrVal.Len(); i++ {
				arrVal.Index(i)
				f, err := formatPropertiesPair(arrFormat(key, i), arrVal.Index(i).Interface())
				if err != nil {
					return nil, err
				}
				assignMap(f, result)
			}
		case reflect.Pointer:
			f, err := formatPropertiesPair(key, reflect.ValueOf(a).Elem().Interface())
			if err != nil {
				return nil, err
			}
			assignMap(f, result)
		default:
			result[key] = a
		}
	}
	return result, nil
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

func arrFormat(k string, i int) string {
	return fmt.Sprintf("%s[%d]", k, i)
}

func assignMap(f, t map[string]any) {
	for s, a := range f {
		t[s] = a
	}
}

func arrSplit(s string) (key string, idx int, hasIdx bool) {
	key = s
	if ls := len(s); ls < 3 || s[ls-1] != ']' {
		return
	}
	var (
		ql int
		qr = len(s) - 1
	)
	for i := qr; i >= 0; i-- {
		if s[i] == '[' {
			ql = i + 1
			break
		}
	}
	if ql >= qr {
		return
	}
	idxStr := s[ql:qr]
	var err error
	idx, err = strconv2.ParseInt(idxStr, 10)
	if err != nil {
		return
	}
	hasIdx = true
	key = s[:ql-1]
	return
}

func setSlice(anies []any, idx int, a any) []any {
	var arr = make([]any, maxI(len(anies), idx+1))
	if len(anies) < idx {
		copy(arr, anies)
		arr[idx] = a
	} else {
		copy(arr[:idx], anies[:idx])
		copy(arr[idx:], anies[idx:])
		arr[idx] = a
	}
	return arr
}

func maxI(i, j int) int {
	if i > j {
		return i
	}
	return j
}
