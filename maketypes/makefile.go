package maketypes

type Makefile interface {
	GetName()
	Download() error
}
