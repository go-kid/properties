package properties

import (
	"bytes"
	"fmt"
	"github.com/go-kid/strconv2"
	"strings"
)

func Unmarshal(data []byte, v any) error {
	data = bytes.Trim(data, "\n")
	buffer := bytes.NewBuffer(data)
	var pm = New()
	for i := 1; ; i++ {
		line, err := buffer.ReadString('\n')
		line = strings.TrimSuffix(line, "\n")
		if len(line) != 0 && line[0] != '#' {
			pairs := strings.SplitN(line, "=", 2)
			if len(pairs) != 2 {
				return fmt.Errorf("parse properties failed at line %d: no pairs found: %#v", i, line)
			}
			key := pairs[0]
			var val any = pairs[1]
			if parsedVal, err := strconv2.ParseAny(pairs[1]); err == nil {
				val = parsedVal
			}
			pm.Set(key, val)
		}
		if err != nil {
			break
		}
	}
	return pm.Unmarshal(v)
}
