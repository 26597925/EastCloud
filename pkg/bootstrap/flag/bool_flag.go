package flag

import (
	"os"
	"strings"
)

type BoolFlag struct {
	Name     string
	Usage    string
	EnvVar   string
	Default  bool
	Variable *bool
	Action   func(string, *Set)
}

// Apply implements of Flag Apply function.
func (f *BoolFlag) Apply(set *Set) {
	for _, field := range strings.Split(f.Name, ",") {
		field = strings.TrimSpace(field)
		if f.Variable != nil {
			set.FlagSet.BoolVar(f.Variable, field, f.Default, f.Usage)
		}

		set.FlagSet.Bool(field, f.Default, f.Usage)
		set.actions[field] = f.Action
		set.environs[field] = os.Getenv(f.EnvVar)
	}
}
