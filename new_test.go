package properties

import (
	"reflect"
	"testing"
)

func TestToPropertiesMap(t *testing.T) {
	type args struct {
		m Properties
	}
	tests := []struct {
		name string
		args args
		want map[string]any
	}{
		{
			name: "1",
			args: args{
				m: Properties{
					"a": map[string]any{
						"b": map[string]any{
							"c":  123,
							"c2": "foo",
						},
						"b2": "bar",
						"c":  []string{"foo", "bar"},
					},
				},
			},
			want: map[string]any{
				"a.b.c":  123,
				"a.b.c2": "foo",
				"a.b2":   "bar",
				"a.c":    []string{"foo", "bar"},
			},
		},
		{
			name: "1",
			args: args{
				m: tPersonMap,
			},
			want: tPersonProp,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.args.m.ToPropertiesMap(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewFromMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

type person struct {
	Name  string `properties:"name"`
	Phone struct {
		Code   string `properties:"phone_code"`
		Number string `properties:"phone_number"`
	} `properties:"mobile"`
	Hobbies []string `properties:"hobbies"`
}

var (
	tPerson = &person{
		Name: "foo",
		Phone: struct {
			Code   string `properties:"phone_code"`
			Number string `properties:"phone_number"`
		}{
			Code:   "11",
			Number: "123321123321",
		},
		Hobbies: []string{"foo", "bar"},
	}
	tPersonProp = Properties{
		"name":                tPerson.Name,
		"mobile.phone_code":   tPerson.Phone.Code,
		"mobile.phone_number": tPerson.Phone.Number,
		"hobbies":             tPerson.Hobbies,
	}
	tPersonMap = map[string]any{
		"name": tPerson.Name,
		"mobile": map[string]any{
			"phone_code":   tPerson.Phone.Code,
			"phone_number": tPerson.Phone.Number,
		},
		"hobbies": tPerson.Hobbies,
	}
)

func TestNewFromAny(t *testing.T) {
	type args struct {
		v any
	}
	tests := []struct {
		name    string
		args    args
		want    Properties
		wantErr bool
	}{
		{
			name: "",
			args: args{
				v: tPerson,
			},
			want:    tPersonMap,
			wantErr: false,
		},
		{
			name: "",
			args: args{
				v: "abc",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "",
			args: args{
				v: 123,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewFromAny(tt.args.v)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewFromAny() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewFromAny() got = %v, want %v", got, tt.want)
			}
		})
	}
}
