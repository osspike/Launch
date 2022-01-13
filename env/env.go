package env

import (
	"encoding/json"
	"os"
	"strings"
)

type Variables map[string]Value

type Value struct {
	V  string
	Vs []string
}

func (v *Value) UnmarshalJSON(b []byte) (err error) {
	err = json.Unmarshal(b, &v.V)
	if err == nil {
		return
	}

	err = json.Unmarshal(b, &v.Vs)
	return
}

func (v *Value) String() string {
	if v.V != "" {
		return v.V
	} else {
		return strings.Join(v.Vs, string(os.PathListSeparator))
	}
}

func (vs *Variables) Set() {
	for k, v := range *vs {
		os.Setenv(k, v.String())
	}
}
