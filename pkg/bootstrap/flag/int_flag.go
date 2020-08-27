package flag

import "strings"

type IntFlag struct {
	Name     string
	Usage    string
	Default  int
	Variable *int
	Action   func(string, *Set)
}

// Apply implements of Flag Apply function.
func (f *IntFlag) Apply(set *Set) {
	for _, field := range strings.Split(f.Name, ",") {
		field = strings.TrimSpace(field)
		if f.Variable != nil {
			set.FlagSet.IntVar(f.Variable, field, f.Default, f.Usage)
		}
		set.FlagSet.Int(field, f.Default, f.Usage)
		set.actions[field] = f.Action
	}
}
