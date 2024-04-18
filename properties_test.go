package properties

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestProperties_Build(t *testing.T) {
	type C struct {
		A string
	}
	type input struct {
		key string
		val any
	}
	t.Run("Set", func(t *testing.T) {
		tests := []struct {
			name   string
			inputs []input
			want   Properties
		}{
			{
				name: "1",
				inputs: []input{
					{key: "", val: nil},
					{key: "a.b.c", val: 123},
					{key: "a.b.c2", val: "foo"},
					{key: "a.b.c3", val: &C{A: "abc"}},
					{key: "a.b2", val: "bar"},
					{key: "a.c", val: &C{A: "abc"}},
					{key: "c", val: &C{A: "abc"}},
				},
				want: map[string]any{
					"": nil,
					"a": map[string]any{
						"b": map[string]any{
							"c":  123,
							"c2": "foo",
							"c3": map[string]any{
								"A": "abc",
							},
						},
						"b2": "bar",
						"c": map[string]any{
							"A": "abc",
						},
					},
					"c": map[string]any{
						"A": "abc",
					},
				},
			},
		}
		for _, tt := range tests {
			pm := New()
			for _, a := range tt.inputs {
				pm.Set(a.key, a.val)
			}
			t.Run(tt.name, func(t *testing.T) {
				assert.Equal(t, tt.want, pm)
			})
		}
	})

	t.Run("Add", func(t *testing.T) {
		tests := []struct {
			name   string
			inputs []input
			want   Properties
		}{
			{
				name: "1",
				inputs: []input{
					{key: "a", val: nil},
					{key: "a", val: 123},
					{key: "a", val: "foo"},
					{key: "a", val: &C{A: "abc"}},
					{key: "b.c", val: nil},
					{key: "b.c", val: 123},
					{key: "b.c", val: "foo"},
					{key: "b.c", val: &C{A: "abc"}},
				},
				want: map[string]any{
					"a": []any{nil, 123, "foo", map[string]any{"A": "abc"}},
					"b": map[string]any{
						"c": []any{nil, 123, "foo", map[string]any{"A": "abc"}},
					},
				},
			},
		}
		for _, tt := range tests {
			pm := New()
			for _, a := range tt.inputs {
				pm.Add(a.key, a.val)
			}
			t.Run(tt.name, func(t *testing.T) {
				assert.Equal(t, tt.want, pm)
			})
		}
	})
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
					key:   "person",
					want:  tPersonMap,
					want1: true,
				},
			},
		},
		{
			name: "",
			p:    New(),
			args: args{
				key: "a.d",
				val: 123,
			},
			wants: []wants{
				{
					key:   "a.b.c.d",
					want:  nil,
					want1: false,
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
			p:    tPersonMap,
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

func TestProperties_Marshal(t *testing.T) {
	tests := []struct {
		name    string
		p       Properties
		want    string
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "",
			p:    tPersonMap,
			want: `hobbies[0]=foo

hobbies[1]=bar

mobile.phone_code=11
mobile.phone_number=123321123321

name=foo
`,
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.p.Marshal()
			if !tt.wantErr(t, err, fmt.Sprintf("Marshal()")) {
				return
			}
			assert.Equalf(t, tt.want, string(got), string(got))
		})
	}
}

func TestProperties_Unmarshal(t *testing.T) {
	type args struct {
		v any
	}
	tests := []struct {
		name    string
		p       Properties
		args    args
		want    any
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "",
			p:    tPersonMap,
			args: args{
				&person{},
			},
			want:    tPerson,
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.wantErr(t, tt.p.Unmarshal(tt.args.v), fmt.Sprintf("Unmarshal(%v)", tt.args.v))
			assert.Equal(t, tt.want, tt.args.v)
		})
	}
}

func TestProperties_SetWithMode(t *testing.T) {
	type args struct {
		key  string
		val  any
		mode SetMode
	}
	tests := []struct {
		name  string
		p     Properties
		args  []args
		wants Properties
	}{
		{
			name: "default overwrite",
			p:    New(),
			args: []args{
				{
					key:  "a.b",
					val:  1,
					mode: 0,
				},
				{
					key:  "a.b",
					val:  2,
					mode: 0,
				},
			},
			wants: map[string]any{
				"a": map[string]any{
					"b": 2,
				},
			},
		},
		{
			name: "append",
			p:    New(),
			args: []args{
				{
					key:  "a.b",
					val:  1,
					mode: Append,
				},
				{
					key:  "a.b",
					val:  2,
					mode: Append,
				},
			},
			wants: map[string]any{
				"a": map[string]any{
					"b": []any{1, 2},
				},
			},
		},
		{
			name: "overwrite type error",
			p:    New(),
			args: []args{
				{
					key:  "a.b",
					val:  1,
					mode: 0,
				},
				{
					key:  "a.b.c",
					val:  2,
					mode: 0,
				},
			},
			wants: map[string]any{
				"a": map[string]any{
					"b": []any{1, 2},
				},
			},
		},
		{
			name: "overwrite type",
			p:    New(),
			args: []args{
				{
					key:  "a.b",
					val:  1,
					mode: OverwriteType,
				},
				{
					key:  "a.b.c",
					val:  2,
					mode: OverwriteType,
				},
			},
			wants: map[string]any{
				"a": map[string]any{
					"b": map[string]any{
						"c": 2,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if err := recover(); err != nil {
					assert.Error(t, err.(error))
				}
			}()
			for _, arg := range tt.args {
				tt.p.SetWithMode(arg.key, arg.val, arg.mode)
			}
			assert.Equal(t, tt.wants, tt.p)
		})
	}
}
