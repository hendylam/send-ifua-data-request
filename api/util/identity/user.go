package identity

// User struct
type User struct {
	id     string
	name   string
	email  string
	role   int
	status int
}

// NewUser func
func NewUser(id string, name string, email string, role int, status int) *User {
	return &User{
		id:     id,
		name:   name,
		email:  email,
		role:   role,
		status: status,
	}
}

// ID func
func (u *User) ID() string {
	return u.id
}

// Name func
func (u *User) Name() string {
	return u.name
}

// Email func
func (u *User) Email() string {
	return u.email
}

// Role func
func (u *User) Role() int {
	return u.role
}

// Status func
func (u *User) Status() int {
	return u.status
}
