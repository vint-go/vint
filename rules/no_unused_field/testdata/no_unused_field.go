package fixtures

import "sync"

// Invalid: unexported field never accessed anywhere.

type User struct {
	Name     string
	Age      int
	nickname string // MATCH /field nickname is unused/
}

func NewUser(name string, age int) User {
	return User{Name: name, Age: age}
}

func (u User) Display() string {
	return u.Name
}

// Invalid: unexported field never accessed.

type config struct {
	host    string
	port    int
	timeout int // MATCH /field timeout is unused/
}

func newConfig() config {
	return config{host: "localhost", port: 8080}
}

func (c config) Address() string {
	return c.host + ":" + string(rune(c.port))
}

// Invalid: unexported field never accessed.

type response struct {
	status  int
	body    string
	headers map[string]string // MATCH /field headers is unused/
}

func ok(body string) response {
	return response{status: 200, body: body}
}

func (r response) String() string {
	return r.body
}

// Valid: all fields used.

type userValid struct {
	Name     string
	Age      int
	nickname string
}

func NewUserValid(name string, age int, nick string) userValid {
	return userValid{Name: name, Age: age, nickname: nick}
}

func (u userValid) Nickname() string {
	return u.nickname
}

// Valid: exported fields are considered used by default.

type Config struct {
	Host    string
	Port    int
	Timeout int
}

// Valid: sync.Mutex embedded field is considered used.

type SafeMap struct {
	mu sync.Mutex
	m  map[string]string
}

func (s *SafeMap) Set(k, v string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.m[k] = v
}

func (s *SafeMap) Get(k string) string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.m[k]
}

// Valid: anonymous struct fields are always considered used (rule 11.1).

func process() {
	data := struct {
		Name  string
		Value int
	}{
		Name:  "test",
		Value: 42,
	}
	_ = data.Name
}
