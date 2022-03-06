package main

import (
	"flag"
	"fmt"
	"path/filepath"

	"github.com/030/p2iwd/internal/pull"
	"github.com/030/p2iwd/internal/push"
	homedir "github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
	jww "github.com/spf13/jwalterweatherman"
	"github.com/spf13/viper"
)

type Login struct {
	Debug, Dir, Host, Pass, Push, User string
}

type Project struct {
	Home, Host, Name, Pass, User string
}

const ext = "yml"

func (p *Project) viperBase(path, filename string) error {
	home, err := homedir.Dir()
	if err != nil {
		return err
	}

	viper.SetConfigName(filename)
	viper.SetConfigType(ext)
	viper.AddConfigPath(filepath.Join(home, "."+p.Name))

	if path != "" {
		viper.SetConfigFile(filepath.Join(path, filename+"."+ext))
	}

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("fatal error config file: %v", err)
	}

	return nil
}

func (p *Project) credsOrConfigValue(creds bool, key string) (string, error) {
	file := "config"
	if creds {
		file = "creds"
	}

	if err := p.viperBase(p.Home, file); err != nil {
		return "", err
	}
	value := viper.GetString(key)
	if value == "" {
		return "", fmt.Errorf("key '%s' not found. Check whether it resides and is populated populated in '%s'", key, viper.ConfigFileUsed())
	}
	return value, nil
}

func (p *Project) dockerRegistry() (Login, error) {
	debug, err := p.credsOrConfigValue(false, "debug")
	if err != nil {
		return Login{}, err
	}

	dir, err := p.credsOrConfigValue(false, "dir")
	if err != nil {
		return Login{}, err
	}

	host, err := p.credsOrConfigValue(false, "host")
	if err != nil {
		return Login{}, err
	}

	push, err := p.credsOrConfigValue(false, "push")
	if err != nil {
		return Login{}, err
	}

	pass, err := p.credsOrConfigValue(true, "pass")
	if err != nil {
		return Login{}, err
	}

	user, err := p.credsOrConfigValue(true, "user")
	if err != nil {
		return Login{}, err
	}

	return Login{Debug: debug, Dir: dir, Host: host, Pass: pass, Push: push, User: user}, err
}

func main() {
	home := flag.String("home", "", "the home folder that contains the config.yml and creds.yml file")

	flag.Parse()

	if *home != "" {
		viper.SetConfigFile(*home)
	}

	p := Project{Home: *home, Name: "p2iwd"}
	l, err := p.dockerRegistry()
	if err != nil {
		log.Fatal(err)
	}

	if l.Debug == "true" {
		enableDebug()
	}

	if l.Push == "true" {
		pdr := push.DockerRegistry{User: l.User, Pass: l.Pass, Dir: l.Dir, Host: l.Host}
		if err := pdr.All(); err != nil {
			log.Fatal(err)
		}
	} else {
		pdr := pull.DockerRegistry{User: l.User, Pass: l.Pass, Dir: l.Dir, Host: l.Host}
		if err := pdr.All(); err != nil {
			log.Fatal(err)
		}
	}
}

func enableDebug() {
	log.SetReportCaller(true)
	if true {
		log.SetLevel(log.DebugLevel)

		// Added to be able to debug viper (used to read the config file)
		// Viper is using a different logger
		jww.SetLogThreshold(jww.LevelTrace)
		jww.SetStdoutThreshold(jww.LevelDebug)
	}
}
