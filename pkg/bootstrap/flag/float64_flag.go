package flag

import "strings"

type Float64Flag struct {
	Name     string
	Usage    string
	Default  float64
	Variable *float64
	Action   func(string, *Set)
}

// Apply implements of Flag Apply function.
func (f *Float64Flag) Apply(set *Set) {
	for _, field := range strings.Split(f.Name, ",") {
		field = strings.TrimSpace(field)
		if f.Variable != nil {
			set.FlagSet.Float64Var(f.Variable, field, f.Default, f.Usage)
		}
		set.FlagSet.Float64(field, f.Default, f.Usage)
		set.actions[field] = f.Action
	}
}