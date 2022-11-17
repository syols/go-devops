package config

import (
	"flag"
	"os"
)

// Variable struct
type Variable struct {
	value *string
	env   string
	name  string
}

// EnvironmentVariables function of a certain type
type EnvironmentVariables map[Variable]func(string) Option

// NewEnvironmentVariables creates EnvironmentVariables struct
func NewEnvironmentVariables() EnvironmentVariables {
	return EnvironmentVariables{
		newVariable("ADDRESS", "a"):          withAddress,
		newVariable("REPORT_INTERVAL", "ri"): withReportInterval,
		newVariable("POLL_INTERVAL", "p"):    withPollInterval,
		newVariable("CLIENT_TIMEOUT", "c"):   withClientTimeout,
		newVariable("STORE_INTERVAL", "i"):   withStoreInterval,
		newVariable("RESTORE", "r"):          withRestore,
		newVariable("KEY", "k"):              withKey,
		newVariable("STORE_FILE", "f"):       withStoreFile,
		newVariable("DATABASE_DSN", "d"):     withDatabase,
	}
}

// Options returns []Options  from env
func (e EnvironmentVariables) Options() (options []Option) {
	flag.Parse()
	for k, v := range e {
		variable, isOk := os.LookupEnv(k.env)
		if !isOk {
			variable = *k.value
		}

		if variable != "" {
			options = append(options, v(variable))
		}

	}
	return
}

func newVariable(env, name string) Variable {
	value := flag.String(name, "", env)
	return Variable{
		env:   env,
		name:  name,
		value: value,
	}
}
