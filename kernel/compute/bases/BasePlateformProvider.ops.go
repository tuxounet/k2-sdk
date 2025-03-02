package bases

func (p *BasePlateformProvider[D]) Nuke() error {

	p.definitions = make([]D, 0)
	return nil
}

func (p *BasePlateformProvider[D]) RegisterDefinition(definition D) error {

	p.definitions = append(p.definitions, definition)

	return nil
}
