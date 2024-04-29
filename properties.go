package properties

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
)

type Properties map[string]any

func (p Properties) Set(key string, val any) {
	p.SetWithMode(key, val, 0)
}

func (p Properties) Add(key string, val any) {
	p.SetWithMode(key, val, Append)
}

func (p Properties) SetWithMode(key string, val any, mode SetMode) {
	if val == nil {
		p.set(key, val, mode)
		return
	}
	switch t := reflect.TypeOf(val); t.Kind() {
	case reflect.Map, reflect.Struct:
		err := p.flatSet(key, val, mode)
		if err != nil {
			panic(err)
		}
	case reflect.Pointer:
		if eleKind := t.Elem().Kind(); eleKind == reflect.Map || eleKind == reflect.Struct {
			err := p.flatSet(key, val, mode)
			if err != nil {
				panic(err)
			}
			return
		}
		fallthrough
	default:
		p.set(key, val, mode)
	}
}

func (p Properties) Get(key string) (any, bool) {
	return p.get(key)
}

func (p Properties) flatSet(key string, val any, mode SetMode) error {
	subProp, err := decodeToMap(val)
	if err != nil {
		return err
	}
	p.set(key, subProp, mode)
	return nil
}

func (p Properties) set(key string, val any, mode SetMode) {
	buildMap(key, val, (*map[string]any)(&p), mode)
}

func (p Properties) get(key string) (any, bool) {
	return getMap(p, key)
}

type ValueSet struct {
	Key   string
	Value any
}

func (p Properties) ValueSets() []*ValueSet {
	properties := p.ToPropertiesMap()
	sets := make([]*ValueSet, len(properties))
	i := 0
	for s, a := range properties {
		sets[i] = &ValueSet{
			Key:   s,
			Value: a,
		}
		i++
	}
	sort.Slice(sets, func(i, j int) bool {
		return sets[i].Key < sets[j].Key
	})
	return sets
}

func (p Properties) ToPropertiesMap() map[string]any {
	return flatten("", p)
}

func (p Properties) Marshal() ([]byte, error) {
	var pairs = make(map[string]any)
	for key, a := range p.ToPropertiesMap() {
		f, err := formatPropertiesPair(key, a)
		if err != nil {
			return nil, err
		}
		assignMap(f, pairs)
	}

	var sb strings.Builder
	var latestGroup string
	var sortedKeys []string
	for s, _ := range pairs {
		sortedKeys = append(sortedKeys, s)
	}
	sort.Slice(sortedKeys, func(i int, j int) bool {
		return sortedKeys[i] < sortedKeys[j]
	})
	for _, key := range sortedKeys {
		v := pairs[key]
		groupIndex := strings.Index(key, ".")
		var group string
		if groupIndex > 0 {
			group = key[:groupIndex]
		} else {
			group = key
		}
		if latestGroup != "" && latestGroup != group {
			sb.WriteString("\n")
		}
		latestGroup = group
		sb.WriteString(fmt.Sprintf("%s=%v\n", key, v))
	}
	return []byte(sb.String()), nil
}

func (p Properties) Unmarshal(v any) error {
	return decode(p, v)
}

func (p Properties) UnmarshalKey(k string, v any) error {
	input, ok := p.Get(k)
	if !ok {
		return fmt.Errorf("unmarshal key error: key '%v' is not found", k)
	}
	return decode(input, v)
}
