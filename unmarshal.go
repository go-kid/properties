package properties

import (
	"bytes"
	"fmt"
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
			err := parsePropertiesPair(pm, line)
			if err != nil {
				return fmt.Errorf("parse properties failed at line %d: %v", i, err)
			}
		}
		if err != nil {
			break
		}
	}
	return pm.Unmarshal(v)
}
