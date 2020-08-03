package env

import (
	"fmt"
	"os"
	"strconv"
)

// Bool func
func Bool(name string, req bool, def interface{}) bool {
	if v, ok := os.LookupEnv(name); ok {
		b, err := strconv.ParseBool(v)
		if err != nil {
			panic(err)
		}
		return b
	} else if req {
		panic(fmt.Sprintf("ENV %s is required", name))
	}

	if def != nil {
		if d, ok := def.(bool); ok {
			return d
		}
		panic(fmt.Sprintf(fmt.Sprintf("%s must be boolean value", name)))
	}
	return false
}

// Int func
func Int(name string, req bool, def interface{}) int {
	if v, ok := os.LookupEnv(name); ok {
		i, err := strconv.Atoi(v)
		if err != nil {
			panic(err)
		}
		return i
	} else if req {
		panic(fmt.Sprintf("ENV %s is required", name))
	}

	if def != nil {
		if d, ok := def.(int); ok {
			return d
		}
		panic(fmt.Sprintf(fmt.Sprintf("%s must be int value", name)))
	}
	return 0
}

// Str func
func Str(name string, req bool, def interface{}) string {
	if v, ok := os.LookupEnv(name); ok {
		return v
	} else if req {
		panic(fmt.Sprintf("ENV %s is required", name))
	}

	if def != nil {
		if d, ok := def.(string); ok {
			return d
		}
		panic(fmt.Sprintf(fmt.Sprintf("%s must be string value", name)))
	}
	return ""
}
