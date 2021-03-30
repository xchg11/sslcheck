package main
import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type data struct {
	CheckSsl struct {
		Domains []string `yaml:"domains"`
		Debug   bool     `yaml:"debug"`
		Notify  struct {
			Mail struct {
				Server string   `yaml:"server"`
				Port   int      `yaml:"port"`
				Email  []string `yaml:"email"`
				Resent string   `yaml:"resent"`
				Auth   struct {
					Login    string `yaml:"login"`
					Password string `yaml:"password"`
				}
			}
			Telegram struct {
				Touser   []string `yaml:"touser"`
				Apitoken string   `yaml:"apitoken"`
				Debug    bool     `yaml:"debug"`
			}
		}
	}
}

func readConfig(filename string) error {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	config = &data{}
	err = yaml.Unmarshal(buf, config)
	if err != nil {
		return err
	}
	return nil
}
