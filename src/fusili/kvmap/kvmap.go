package kvmap

import (
	"fmt"
	"reflect"
)

func getSetting(config map[string]interface{}, setting string, kind reflect.Kind,
	mandatory bool, fallbackValue interface{}) (interface{}, error) {

	if _, ok := config[setting]; !ok {
		if mandatory {
			return fallbackValue, fmt.Errorf("missing mandatory setting `%s'", setting)
		}

		return fallbackValue, nil
	}

	configKind := reflect.ValueOf(config[setting]).Kind()
	if configKind != kind {
		return fallbackValue, fmt.Errorf(
			"setting %q value should be of type %s and not %s",
			setting,
			kind.String(),
			configKind.String())
	}

	return config[setting], nil
}

func GetString(config map[string]interface{}, setting string, mandatory bool) (string, error) {
	value, err := getSetting(config, setting, reflect.String, mandatory, "")

	return value.(string), err
}

func GetInt(config map[string]interface{}, setting string, mandatory bool) (int, error) {
	value, err := getSetting(config, setting, reflect.Float64, mandatory, 0.0)

	return int(value.(float64)), err
}

func GetFloat(config map[string]interface{}, setting string, mandatory bool) (float64, error) {
	value, err := getSetting(config, setting, reflect.Float64, mandatory, 0.0)

	return value.(float64), err
}

func GetBool(config map[string]interface{}, setting string, mandatory bool) (bool, error) {
	value, err := getSetting(config, setting, reflect.Bool, mandatory, false)

	return value.(bool), err
}

func GetStringSlice(config map[string]interface{}, setting string, mandatory bool) ([]string, error) {
	value, err := getSetting(config, setting, reflect.Slice, mandatory, nil)
	if err != nil || value == nil {
		return nil, err
	}

	out := make([]string, 0)
	for _, v := range value.([]interface{}) {
		if reflect.ValueOf(v).Kind() != reflect.String {
			return nil, fmt.Errorf("setting `%s' should be slice of strings and not %s", setting,
				reflect.ValueOf(v).Kind().String())
		} else {
			out = append(out, v.(string))
		}
	}

	return out, nil
}

func GetStringMap(config map[string]interface{}, setting string, mandatory bool) (map[string]interface{}, error) {
	value, err := getSetting(config, setting, reflect.Map, mandatory, nil)

	return value.(map[string]interface{}), err
}
