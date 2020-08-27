package flag

import "strings"

// UintFlag is an uint flag implements of Flag interface.
type UintFlag struct {
	Name     string
	Usage    string
	Default  uint
	Variable *uint
	Action   func(string, *Set)
}

// Apply implements of Flag Apply function.
func (f *UintFlag) Apply(set *Set) {
	for _, field := range strings.Split(f.Name, ",") {
		field = strings.TrimSpace(field)
		if f.Variable != nil {
			set.FlagSet.UintVar(f.Variable, field, f.Default, f.Usage)
		}
		set.FlagSet.Uint(field, f.Default, f.Usage)
		set.actions[field] = f.Action
	}
}