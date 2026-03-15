package fixtures

type User struct {
	Name  string
	Email string
	Age   int
}

type Option struct{}

type Service struct{}

func (s *Service) process(ctx any, orderID string, options ...Option) error {
	return nil
}

// Invalid: parameter list starts on a new line after the function name.
func createUser( // MATCH /multi-line function signature should have the first parameter on the same line as the func keyword/
	name string,
	email string,
	age int,
) (*User, error) {
	return &User{Name: name, Email: email, Age: age}, nil
}

// Invalid: parameter list starts on a new line after the method name.
func (s *Service) ProcessOrder( // MATCH /multi-line function signature should have the first parameter on the same line as the func keyword/
	ctx any,
	orderID string,
	options ...Option,
) error {
	return s.process(ctx, orderID, options...)
}

// Valid: first parameter is on the same line as the function name.
func createUserValid(name string,
	email string,
	age int,
) (*User, error) {
	return &User{Name: name, Email: email, Age: age}, nil
}

// Valid: first parameter is on the same line as the method name.
func (s *Service) ProcessOrderValid(ctx any,
	orderID string,
	options ...Option,
) error {
	return s.process(ctx, orderID, options...)
}

// Valid: single-line function signature.
func add(a, b int) int {
	return a + b
}
