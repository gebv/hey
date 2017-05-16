package hey

// Chronograph represent main api
type Chronograph struct {
	Manager
}

// New create chronograph with default manager (tarantool)
func New() (*Chronograph, error) {
	manager, err := NewTarantoolManager()
	if err != nil {
		return nil, err
	}
	return NewWithManager(manager)
}

// NewWithManager create chronograph with passed store
func NewWithManager(manager Manager) (*Chronograph, error) {
	return &Chronograph{Manager: manager}, nil
}
