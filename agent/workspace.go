package agent

import (
	"bufio"
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"sync"

	log "github.com/Sirupsen/logrus"
)

type Workspace struct {
	RootDir   string
	Makefiles []string
	Names     []string
}

func (ws *Workspace) Add(content []byte, name string) error {
	logCtx := log.WithFields(log.Fields{
		"workspace": ws.RootDir,
		"makefile":  name,
	})

	dirname := filepath.Join(ws.RootDir, name)
	makefileName := filepath.Join(dirname, "Makefile")

	err := os.MkdirAll(dirname, 0744)
	if err != nil {
		return err
	}
	logCtx.WithField("makefile_path", makefileName).Infof("wrote %s", makefileName)
	err = ioutil.WriteFile(makefileName, content, 0644)
	if err != nil {
		return err
	}

	ws.Makefiles = append(ws.Makefiles, makefileName)
	ws.Names = append(ws.Names, name)
	return nil
}

func (ws *Workspace) Apply() error {
	for _, name := range ws.Names {
		logCtx := log.WithField("name", name)

		cmd := exec.Command("make")
		cmd.Dir = filepath.Join(ws.RootDir, name)

		watcher, err := watchCmd(cmd, logCtx)
		if err != nil {
			return err
		}

		cmd.Start()

		watcher.Wait()
		err = cmd.Wait()
		if err != nil {
			return err
		}
	}

	return nil
}

func watchCmd(cmd *exec.Cmd, logCtx *log.Entry) (*sync.WaitGroup, error) {
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		logCtx.WithError(err).Warn("error creating stdout pipe")
		return nil, err
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		logCtx.WithError(err).Warn("error creating stderr pipe")
		return nil, err
	}

	wg := new(sync.WaitGroup)
	wg.Add(2)

	go outReader(wg, stdout, logHandler(logCtx))
	go outReader(wg, stderr, logHandler(logCtx))

	return wg, nil
}

type outHandler func(string) string

func logHandler(logCtx *log.Entry) outHandler {
	return func(line string) string {
		logCtx.Info(line)
		return ""
	}
}

func outReader(wg *sync.WaitGroup, r io.Reader, handler outHandler) {
	defer wg.Done()

	reader := bufio.NewReader(r)
	var buffer bytes.Buffer

	for {
		buf := make([]byte, 1024)

		n, err := reader.Read(buf)
		if err != nil {
			return
		}

		buf = buf[:n]

		for {
			i := bytes.IndexByte(buf, '\n')
			if i < 0 {
				break
			}

			buffer.Write(buf[0:i])
			handler(buffer.String())
			buffer.Reset()
			buf = buf[i+1:]
		}
		buffer.Write(buf)
	}
}
