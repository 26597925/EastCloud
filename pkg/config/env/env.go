package env

import (
	"os"
	"sapi/pkg/config/api"
	"sapi/pkg/util/stringext"
	"strconv"
	"strings"

	"github.com/imdario/mergo"
)

func Parse(prefixes []string) ([]byte, error) {
	var changes map[string]interface{}

	for _, env := range os.Environ() {

		if len(prefixes) > 0 && !stringext.ExistPrefix(prefixes, env) {
			continue
		}

		pair := strings.SplitN(env, "=", 2)
		value := pair[1]
		keys := strings.Split(strings.ToLower(pair[0]), "_")
		stringext.Reverse(keys)

		tmp := make(map[string]interface{})
		for i, k := range keys {
			if i == 0 {
				if intValue, err := strconv.Atoi(value); err == nil {
					tmp[k] = intValue
				} else if boolValue, err := strconv.ParseBool(value); err == nil {
					tmp[k] = boolValue
				} else {
					tmp[k] = value
				}
				continue
			}

			tmp = map[string]interface{}{k: tmp}
		}

		if err := mergo.Map(&changes, tmp); err != nil {
			return nil, err
		}
	}

	b, err := api.Encoders["json"].Encode(changes)
	if err != nil {
		return nil, err
	}

	return b, nil
}