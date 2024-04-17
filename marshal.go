package properties

func Marshal(v any) ([]byte, error) {
	pm, err := convertAny2Prop(v)
	if err != nil {
		return nil, err
	}
	return pm.Marshal()
}
