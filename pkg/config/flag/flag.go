package flag

import (
	"flag"
	"github.com/imdario/mergo"
	fg "sapi/pkg/bootstrap/flag"
	"sapi/pkg/config/api"
	"sapi/pkg/util/stringext"
	"strings"
)

func Parse(fs *fg.Set, prefixes []string) ([]byte, error) {

	var changes map[string]interface{}

	visitFn := func(f *flag.Flag) {
		n := strings.ToLower(f.Name)
		if !stringext.ExistPrefix(prefixes, n) {
			return
		}

		keys := strings.FieldsFunc(n, stringext.Split)
		stringext.Reverse(keys)

		tmp := make(map[string]interface{})
		for i, k := range keys {
			if i == 0 {
				tmp[k] = f.Value
				continue
			}

			tmp = map[string]interface{}{k: tmp}
		}

		mergo.Map(&changes, tmp) // need to sort errors handling
		return
	}

	fs.Visit(visitFn)

	b, err := api.Encoders["json"].Encode(changes)
	if err != nil {
		return nil, err
	}

	return b, nil
}