package identity

// UserInfo interface
type UserInfo interface {
	ID() string
	Name() string
	Email() string
	Role() int
	Status() int
}
