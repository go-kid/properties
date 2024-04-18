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
		switch value.(type) {
		case map[string]any, Properties:
			tmp := value.(map[string]any)
			buildMap(next, val, &tmp, mode)
		case []any:
			tmpArr := value.([]any)
			var subM map[string]any
			if len(tmpArr) > idx {
				if am, ok := tmpArr[idx].(map[string]any); ok {
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
			value = setSlice(tmpArr, idx, subM)
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
			switch value.(type) {
			case []any:
				value = setSlice(value.([]any), idx, val)
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
			switch value.(type) {
			case []any:
				value = append(value.([]any), val)
			default:
				value = []any{value, val}
			}
		}
	}
	rtmp[key] = value
}

func getMap(m map[string]any, path string) (any, bool) {
	arr := strings.SplitN(path, ".", 2)
	if len(arr) == 2 {
		key := arr[0]
		next := arr[1]
		if sub, ok := m[key]; ok {
			switch sub.(type) {
			case map[string]any, Properties:
				return getMap(sub.(map[string]any), next)
			default:
				return nil, false
			}
		}
	}
	a, ok := m[path]
	return a, ok
}

func decodeToMap(v any) (map[string]any, error) {
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
	return m, nil
}

func flatten(prePath string, m Properties) map[string]any {
	tmp := make(map[string]any)
	for k, v := range m {
		switch v.(type) {
		case map[string]any, Properties:
			subTmp := flatten(path(prePath, k), v.(map[string]any))
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
	var (
		ql, qr int
	)
	key = s
	for i, c := range s {
		switch c {
		case '[':
			ql = i + 1
		case ']':
			qr = i
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
