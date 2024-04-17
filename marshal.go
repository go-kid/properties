package properties

func Marshal(v any) ([]byte, error) {
	pm, err := NewFromAny(v)
	if err != nil {
		return nil, err
	}
	return pm.Marshal()
}
