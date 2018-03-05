package maketypes

import "io/ioutil"

type LocalMakefile struct {
	Name string
	Path string
}

func (lm LocalMakefile) GetName() string {
	return lm.Name
}

func (lm LocalMakefile) Download() ([]byte, error) {
	return ioutil.ReadFile(lm.Path)
}
