package commands

type Command interface {
	Matches(pattern string) bool
	Parse(pattern string) error
	Execute() error
	String() string
}
