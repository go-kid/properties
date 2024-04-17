package properties

func New() Properties {
	return make(Properties)
}

func NewFromMap(m map[string]any) Properties {
	return buildProperties("", m)
}

func NewFromAny(v any) (Properties, error) {
	return convertAny2Prop(v)
}
