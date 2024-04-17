package properties

func New() Properties {
	return make(Properties)
}

func NewFromAny(v any) (Properties, error) {
	return decodeToMap(v)
}
