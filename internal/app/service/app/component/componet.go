package component

type Components struct {
	components []Component
}

func NewComponents() *Components {
	return &Components{
		components: make([]Component, 0),
	}
}

func (cs *Components) Add(c Component) {
	cs.components = append(cs.components, c)
}

func (cs *Components) Items() []Component {
	return cs.components
}
