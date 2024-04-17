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
		if len(line) != 0 {
			pairs := strings.SplitN(line, "=", 2)
			if len(pairs) != 2 {
				return fmt.Errorf("parse properties failed at line %d, no pairs found", i)
			}
			parsedVal, err := strconv2.ParseAny(pairs[1])
			if err != nil {
				parsedVal = pairs[1]
			}
			pm.Set(pairs[0], parsedVal)
		}
		if err != nil {
			break
		}
	}
	return pm.Unmarshal(v)
}
