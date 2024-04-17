package properties

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestProperties_Expand(t *testing.T) {
	type C struct {
		A string
	}
	tests := []struct {
		name string
		p    Properties
		want map[string]any
	}{
		{
			name: "1",
			p: Properties{
				"a.b.c":  123,
				"a.b.c2": "foo",
				"a.b.c3": &C{A: "abc"},
				"a.b2":   "bar",
				"a.c":    &C{A: "abc"},
				"c":      &C{A: "abc"},
			},
			want: map[string]any{
				"a": map[string]any{
					"b": map[string]any{
						"c":  123,
						"c2": "foo",
						"c3": &C{A: "abc"},
					},
					"b2": "bar",
					"c":  &C{A: "abc"},
				},
				"c": &C{A: "abc"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.p.Expand())
		})
	}
}

func TestProperties_Set_Get(t *testing.T) {
	type args struct {
		key string
		val any
	}
	type wants struct {
		key   string
		want  any
		want1 bool
	}
	tests := []struct {
		name  string
		p     Properties
		args  args
		wants []wants
	}{
		{
			name: "1",
			p:    New(),
			args: args{
				key: "a.b.c",
				val: 123,
			},
			wants: []wants{
				{
					key:   "a.b.c",
					want:  123,
					want1: true,
				},
			},
		},
		{
			name: "1",
			p:    New(),
			args: args{
				key: "a.b.c2",
				val: "foo",
			},
			wants: []wants{
				{
					key:   "a.b.c2",
					want:  "foo",
					want1: true,
				},
			},
		},
		{
			name: "1",
			p:    New(),
			args: args{
				key: "a.b2",
				val: "bar",
			},
			wants: []wants{
				{
					key:   "a.b2",
					want:  "bar",
					want1: true,
				},
			},
		},
		{
			name: "",
			p:    New(),
			args: args{
				key: "person",
				val: tPerson,
			},
			wants: []wants{
				{
					key:   "person.name",
					want:  tPerson.Name,
					want1: true,
				},
				{
					key:   "person.mobile.phone_code",
					want:  tPerson.Phone.Code,
					want1: true,
				},
				{
					key:   "person.mobile.phone_number",
					want:  tPerson.Phone.Number,
					want1: true,
				},
				{
					key:   "person.hobbies",
					want:  tPerson.Hobbies,
					want1: true,
				},
			},
		},
		{
			name: "",
			p:    New(),
			args: args{
				key: "person",
				val: tPerson,
			},
			wants: []wants{
				{
					key: "person",
					want: map[string]any{
						"person": tPersonMap,
					},
					want1: true,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.p.Set(tt.args.key, tt.args.val)
			for _, want := range tt.wants {
				got, got1 := tt.p.Get(want.key)
				assert.Equal(t, want.want, got)
				assert.Equal(t, want.want1, got1)
			}
		})
	}
}

func TestProperties_ValueSets(t *testing.T) {
	tests := []struct {
		name string
		p    Properties
		want []*ValueSet
	}{
		{
			name: "",
			p:    NewFromMap(tPersonMap),
			want: []*ValueSet{
				{
					Key:   "hobbies",
					Value: tPerson.Hobbies,
				},
				{
					Key:   "mobile.phone_code",
					Value: tPerson.Phone.Code,
				},
				{
					Key:   "mobile.phone_number",
					Value: tPerson.Phone.Number,
				},
				{
					Key:   "name",
					Value: tPerson.Name,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, tt.p.ValueSets(), "ValueSets()")
		})
	}
}
