package p2iwd

import (
	"github.com/030/p2iwd/internal/app/p2iwd/pull"
	"github.com/030/p2iwd/internal/app/p2iwd/push"
)

type DockerRegistry struct {
	Dir, Host, Pass, User string
}

func (dr *DockerRegistry) Backup() error {
	pdr := pull.DockerRegistry{Dir: dr.Dir, Host: dr.Host, Pass: dr.Pass, User: dr.User}
	if err := pdr.All(""); err != nil {
		return err
	}
	return nil
}

func (dr *DockerRegistry) Upload() error {
	pdr := push.DockerRegistry{Dir: dr.Dir, Host: dr.Host, Pass: dr.Pass, User: dr.User}
	if err := pdr.All(); err != nil {
		return err
	}
	return nil
}
