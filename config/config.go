package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"launch/color"
	"launch/env"
	"strings"
)

type Config struct {
	Vars      map[string]string    `json:"var"`
	Env       env.Variables        `json:"env"`
	Processes map[ProcName]Process `json:"proc"`
}

type ProcName string

type Process struct {
	Color color.Color `json:"color"`
	Path  string      `json:"path"`
	Args  []string    `json:"args"`
	Cwd   string      `json:"cwd"`
}

func (n *ProcName) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	if len(s) != 1 {
		return fmt.Errorf("command name must have lenght equal to 1")
	}

	*n = ProcName(s)
	return nil
}

func Load() (*Config, error) {
	f, err := ioutil.ReadFile("./launch.json")
	if err != nil {
		return nil, err
	}

	c := &Config{}
	err = json.Unmarshal(f, c)
	if err != nil {
		return nil, err
	}

	for varKey, varValue := range c.Vars {
		for envKey, envValue := range c.Env {
			newVarValue := strings.ReplaceAll(envValue.String(), varKey, varValue)
			c.Env[envKey] = env.Value{V: newVarValue}
		}

		for procName, procValue := range c.Processes {
			args := make([]string, 0, len(procValue.Args))
			for _, arg := range procValue.Args {
				args = append(args, strings.ReplaceAll(arg, varKey, varValue))
			}

			procValue.Path = strings.ReplaceAll(procValue.Path, varKey, varValue)
			procValue.Cwd = strings.ReplaceAll(procValue.Cwd, varKey, varValue)
			procValue.Args = args

			c.Processes[procName] = procValue
		}
	}

	return c, nil
}
