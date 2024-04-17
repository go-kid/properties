package properties

import (
	"bytes"
	"fmt"
	"github.com/mitchellh/mapstructure"
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
			pm.Set(pairs[0], pairs[1]) //todo parse to any
		}
		if err != nil {
			break
		}
	}
	config := newDecodeConfig(v)
	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return err
	}
	return decoder.Decode(pm.Expand())
}
