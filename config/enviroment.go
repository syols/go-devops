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
		newVariable("ADDRESS", "a"):             withAddress,
		newVariable("REPORT_INTERVAL", "ri"):    withReportInterval,
		newVariable("POLL_INTERVAL", "p"):       withPollInterval,
		newVariable("CLIENT_TIMEOUT", "c"):      withClientTimeout,
		newVariable("STORE_INTERVAL", "i"):      withStoreInterval,
		newVariable("RESTORE", "r"):             withRestore,
		newVariable("KEY", "k"):                 withKey,
		newVariable("STORE_FILE", "f"):          withStoreFile,
		newVariable("DATABASE_DSN", "d"):        withDatabase,
		newVariable("CRYPTO_KEY", "crypto-key"): withCryptoKey,
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
