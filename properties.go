package properties

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
)

type Properties map[string]any

func (p Properties) Set(key string, val any) {
	p.setWithMode(key, val, false)
}

func (p Properties) Add(key string, val any) {
	p.setWithMode(key, val, true)
}

func (p Properties) setWithMode(key string, val any, appendMode bool) {
	switch t := reflect.TypeOf(val); t.Kind() {
	case reflect.Map, reflect.Struct:
		err := p.flatSet(key, val, appendMode)
		if err != nil {
			panic(err)
		}
	case reflect.Pointer:
		if eleKind := t.Elem().Kind(); eleKind == reflect.Map || eleKind == reflect.Struct {
			err := p.flatSet(key, val, appendMode)
			if err != nil {
				panic(err)
			}
			return
		}
		fallthrough
	default:
		p.set(key, val, appendMode)
	}
}

func (p Properties) Get(key string) (any, bool) {
	v, ok := p.get(key)
	if ok {
		return v, true
	}
	tmp := make(map[string]any)
	for s, a := range p {
		if i := strings.Index(s, key); i >= 0 && s[i+len(key)] == '.' {
			buildMap(s, a, &tmp)
		}
	}
	if len(tmp) == 0 {
		return nil, false
	}
	return tmp, true
}

func (p Properties) flatSet(key string, val any, appendMode bool) error {
	subProp, err := convertAny2Prop(val)
	if err != nil {
		return err
	}
	for s, a := range subProp {
		p.set(path(key, s), a, appendMode)
	}
	return nil
}

func (p Properties) set(key string, val any, appendMode bool) {
	a, ok := p.get(key)
	if !ok {
		p[key] = val
		return
	}
	if appendMode {
		switch a.(type) {
		case []any:
			p[key] = append(a.([]any), val)
		default:
			p[key] = []any{a, val}
		}
	} else {
		p[key] = val
	}
}

func (p Properties) get(key string) (any, bool) {
	a, ok := p[key]
	return a, ok
}

type ValueSet struct {
	Key   string
	Value any
}

func (p Properties) ValueSets() []*ValueSet {
	sets := make([]*ValueSet, len(p))
	i := 0
	for s, a := range p {
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

func (p Properties) Expand() map[string]any {
	tmp := make(map[string]any)
	for k, v := range p {
		buildMap(k, v, &tmp)
	}
	return tmp
}

func (p Properties) Marshal() ([]byte, error) {
	var pairs = make(map[string]any)
	for key, a := range p {
		f, err := formatPropertiesPair(a)
		if err != nil {
			return nil, err
		}
		for fk, fv := range f {
			pairs[path(key, fk)] = fv
		}
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
