package flag

import (
	"os"
	"strings"
)

type StringFlag struct {
	Name     string
	Usage    string
	EnvVar   string
	Default  string
	Variable *string
	Action   func(string, *Set)
}

// Apply implements of Flag Apply function.
func (f *StringFlag) Apply(set *Set) {
	for _, field := range strings.Split(f.Name, ",") {
		field = strings.TrimSpace(field)

		if f.Variable != nil {
			set.FlagSet.StringVar(f.Variable, field, f.Default, f.Usage)
		}
		set.FlagSet.String(field, f.Default, f.Usage)
		set.actions[field] = f.Action

		set.environs[field] = os.Getenv(f.EnvVar)
	}
}
