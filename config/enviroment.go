package config

import (
	"flag"
	"os"
)

type Variable struct {
	env   string
	name  string
	value *string
}

type EnvironmentVariables map[Variable]func(string) Option

func newVariable(env, name string) Variable {
	value := flag.String(name, "", env)
	return Variable{
		env:   env,
		name:  name,
		value: value,
	}
}

func newVariables() EnvironmentVariables {
	return EnvironmentVariables{
		newVariable("ADDRESS", "a"):          WithAddress,
		newVariable("REPORT_INTERVAL", "ri"): WithReportInterval,
		newVariable("POLL_INTERVAL", "p"):    WithPollInterval,
		newVariable("CLIENT_TIMEOUT", "c"):   WithClientTimeout,
		newVariable("STORE_INTERVAL", "i"):   WithStoreInterval,
		newVariable("RESTORE", "r"):          WithRestore,
		newVariable("KEY", "k"):              WithKey,
		newVariable("STORE_FILE", "f"):       WithStoreFile,
		newVariable("DATABASE_DSN", "d"):     WithDatabase,
	}
}

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
