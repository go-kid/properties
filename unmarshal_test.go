package properties

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
	"testing"
)

func TestUnmarshal(t *testing.T) {
	//fmt.Println(arrSplit("a.b[123]"))
	//fmt.Println(arrSplit("[123]"))
	//fmt.Println(arrSplit("]123["))
	//fmt.Println(arrSplit("123"))
	//fmt.Println(arrSplit("123[ab]"))
	//fmt.Println(arrSplit("123[123ab]"))
	var pm = map[string]any{}
	err := Unmarshal([]byte(`
FileApiConfig[0].callback=/create/station
FileApiConfig[0].data_source_name=create_station
FileApiConfig[0].file_type[0]=csv
FileApiConfig[0].file_type[1]=xlsx
FileApiConfig[0].notify_url=
FileApiConfig[0].version=1

FileApiConfig[1].callback=/create/strike
FileApiConfig[1].data_source_name=create_strike
FileApiConfig[1].file_type[0]=csv
FileApiConfig[1].file_type[1]=xlsx
FileApiConfig[1].notify_url=
FileApiConfig[1].version=1

FileApiConfig[2].callback=/create/area_code
FileApiConfig[2].data_source_name=create_area_code
FileApiConfig[2].file_type[0]=csv
FileApiConfig[2].file_type[1]=xlsx
FileApiConfig[2].notify_url=
FileApiConfig[2].version=1
`), &pm)
	assert.NoError(t, err)
	//pm := New()
	//pm.Set("FileApiConfig[0].callback", "/create/station")
	//pm.Set("FileApiConfig[0].data_source_name", "create_station")
	//pm.Set("FileApiConfig[0].file_type[0]", "csv")
	//pm.Set("FileApiConfig[0].file_type[1]", "xlsx")
	//pm.Set("FileApiConfig[0].notify_url", "")
	//
	//pm.Set("FileApiConfig[1].callback", "/create/strike")
	//pm.Set("FileApiConfig[1].data_source_name", "create_strike")
	//pm.Set("FileApiConfig[1].file_type[0]", "csv")
	//pm.Set("FileApiConfig[1].file_type[1]", "xlsx")
	//pm.Set("FileApiConfig[1].notify_url", "")
	////fmt.Println(pm)
	marshal, err := yaml.Marshal(pm)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(marshal))
	p, err := NewFromAny(pm)
	if err != nil {
		panic(err)
	}

	fmt.Println(p.Get("FileApiConfig[-1]"))

	fmt.Println(p.Get("FileApiConfig[0]"))
	fmt.Println(p.Get("FileApiConfig[0].file_type[0]"))
	fmt.Println(p.Get("FileApiConfig[0].file_type[1]"))
	fmt.Println(p.Get("FileApiConfig[0].file_type[2]"))

	fmt.Println(p.Get("FileApiConfig[1]"))
	fmt.Println(p.Get("FileApiConfig[1].file_type[0]"))
	fmt.Println(p.Get("FileApiConfig[1].file_type[1]"))
	fmt.Println(p.Get("FileApiConfig[1].file_type[2]"))

	fmt.Println(p.Get("FileApiConfig[2]"))
	fmt.Println(p.Get("FileApiConfig[2].file_type[0]"))
	fmt.Println(p.Get("FileApiConfig[2].file_type[1]"))
	fmt.Println(p.Get("FileApiConfig[2].file_type[2]"))

	fmt.Println(p.Get("FileApiConfig[3]"))
	fmt.Println(p.Get("FileApiConfig[3].file_type[0]"))
	fmt.Println(p.Get("FileApiConfig[3].file_type[1]"))
	fmt.Println(p.Get("FileApiConfig[3].file_type[2]"))

	//pm.Set("a.b[0].c[1]", 1)
	//pm.Set("a.b[0].c[2].d", 1)
	//pm.Add("a.b[0].c[2].d", 2)
	//pm.Set("b[0]", 0)
	//pm.Set("b[1]", 1)
	//fmt.Println(pm)
	//m := map[string]any{
	//	"a": map[string]any{
	//		"b": []any{
	//			map[string]any{
	//				"c": []any{
	//					0,
	//					1,
	//				},
	//			},
	//		},
	//	},
	//}
	//fmt.Println(m)
}
