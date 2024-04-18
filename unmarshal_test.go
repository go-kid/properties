package properties

import (
	"fmt"
	"testing"
)

func TestUnmarshal(t *testing.T) {
	//fmt.Println(arrSplit("a.b[123]"))
	//fmt.Println(arrSplit("[123]"))
	//fmt.Println(arrSplit("]123["))
	//fmt.Println(arrSplit("123"))
	//fmt.Println(arrSplit("123[ab]"))
	//fmt.Println(arrSplit("123[123ab]"))
	//	var m = map[string]any{}
	//	err := Unmarshal([]byte(`
	//#a.b=1
	//#a.c=2
	//#b.c=3
	//#b.d=4
	//#c.d[10]=0
	//#c.d[1]=1
	//d.e[0].f[0]=1
	//`), &m)
	//	assert.NoError(t, err)
	//	fmt.Printf("%#v", m)
	pm := New()
	pm.Set("a.b[0]", 1)
	fmt.Println(pm)
}
