package restapi

func NewAPITarget(category, ttype string, fields []*InputField) *Target {
	return &Target{
		Category:    category,
		Type:        ttype,
		InputFields: fields,
	}
}

type InputFieldsBuilder struct {
	fields []*InputField
}

func NewInputFieldsBuilder() *InputFieldsBuilder {
	inputs := []*InputField{}
	return &InputFieldsBuilder{
		fields: inputs,
	}
}

func (b *InputFieldsBuilder) With(name, value string) *InputFieldsBuilder {
	field := &InputField{
		Name:  name,
		Value: value,
	}

	b.fields = append(b.fields, field)
	return b
}

func (b *InputFieldsBuilder) Create() []*InputField {
	return b.fields
}
