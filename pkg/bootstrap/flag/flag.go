package flag

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
)

type Flag interface {
	Apply(*Set)
}

type Set struct {
	*flag.FlagSet
	flags    []Flag
	actions  map[string]func(string, *Set)
	environs map[string]string
}

func NewFlagSet() *Set {
	procName := filepath.Base(os.Args[0])

	return &Set{
		FlagSet:  flag.NewFlagSet(procName, flag.ExitOnError),
		flags:    defaultFlags(),
		actions:  make(map[string]func(string, *Set)),
		environs: make(map[string]string),
	}
}

func (fs *Set) Register(flags ...Flag) {
	fs.flags = append(fs.flags, flags...)
}

func (fs *Set) Parse() error {
	if fs.Parsed() {
		return nil
	}
	for _, f := range fs.flags {
		f.Apply(fs)
	}

	if err := fs.FlagSet.Parse(os.Args[1:]); err != nil {
		return err
	}
	fs.FlagSet.Visit(func(f *flag.Flag) {
		if action, ok := fs.actions[f.Name]; ok && action != nil {
			action(f.Name, fs)
		}
		if env, ok := fs.environs[f.Name]; ok {
			fs.environs[f.Name] = env
		}
	})

	return nil
}

func (fs *Set) Lookup(name string) *flag.Flag {
	flag := fs.FlagSet.Lookup(name)

	if flag != nil {
		if env, ok := fs.environs[name]; ok {
			if env != "" {
				flag.Value.Set(env)
			}
		}
	}
	return flag
}

func (fs *Set) BoolE(name string) (bool, error) {
	flag := fs.Lookup(name)
	if flag != nil {
		return strconv.ParseBool(flag.Value.String())
	}

	return false, fmt.Errorf("undefined flag name: %s", name)
}

func (fs *Set) Bool(name string) bool {
	ret, _ := fs.BoolE(name)
	return ret
}

func (fs *Set) StringE(name string) (string, error) {
	flag := fs.Lookup(name)
	if flag != nil {
		return flag.Value.String(), nil
	}

	return "", fmt.Errorf("undefined flag name: %s", name)
}

func (fs *Set) String(name string) string {
	ret, _ := fs.StringE(name)
	return ret
}

func (fs *Set) IntE(name string) (int64, error) {
	flag := fs.Lookup(name)
	if flag != nil {
		return strconv.ParseInt(flag.Value.String(), 10, 64)
	}

	return 0, fmt.Errorf("undefined flag name: %s", name)
}

func (fs *Set) Int(name string) int64 {
	ret, _ := fs.IntE(name)
	return ret
}

func (fs *Set) UintE(name string) (uint64, error) {
	flag := fs.Lookup(name)
	if flag != nil {
		return strconv.ParseUint(flag.Value.String(), 10, 64)
	}

	return 0, fmt.Errorf("undefined flag name: %s", name)
}

func (fs *Set) Uint(name string) uint64 {
	ret, _ := fs.UintE(name)
	return ret
}

func (fs *Set) Float64E(name string) (float64, error) {
	flag := fs.Lookup(name)
	if flag != nil {
		return strconv.ParseFloat(flag.Value.String(), 64)
	}

	return 0.0, fmt.Errorf("undefined flag name: %s", name)
}

func (fs *Set) Float64(name string) float64 {
	ret, _ := fs.Float64E(name)
	return ret
}

func defaultFlags() []Flag {
	return []Flag{
		// HelpFlag prints usage of application.
		&BoolFlag{
			Name:  "help",
			Usage: "--help, show help information",
			Action: func(name string, fs *Set) {
				fs.PrintDefaults()
				os.Exit(0)
			},
		},
	}
}