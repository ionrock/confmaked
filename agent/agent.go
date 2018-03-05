package agent

import (
	"time"

	log "github.com/Sirupsen/logrus"
)

type Makefile interface {
	GetName() string
	Download() ([]byte, error)
}

type Agent struct {
	Cadence   time.Duration
	Makefiles []Makefile
	Workspace *Workspace
}

func New(makefiles []Makefile, dir string) *Agent {
	return &Agent{
		Cadence:   time.Second * 5,
		Makefiles: makefiles,
		Workspace: &Workspace{RootDir: dir},
	}
}

func (a *Agent) Run() error {
	log.Info("starting run")
	for _, mf := range a.Makefiles {
		log.WithField("name", mf.GetName()).Info("executing make task")
		content, err := mf.Download()
		if err != nil {
			return err
		}

		a.Workspace.Add(content, mf.GetName())
	}

	return a.Workspace.Apply()
}
