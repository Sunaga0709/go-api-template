package environment

import (
	"fmt"
	"strings"
)

const (
	EnvLocal = "local"
	EnvDev   = "dev"
	EnvStg   = "stg"
	EnvPrd   = "prd"
)

type Environment struct {
	string
}

func NewEnvironment(env string) (Environment, error) {
	env = strings.ToLower(strings.TrimSpace(env))

	var validEnv string
	switch env {
	case "local":
		validEnv = EnvLocal
	case "dev", "develop", "development":
		validEnv = EnvDev
	case "stg", "stage", "staging":
		validEnv = EnvStg
	case "prd", "prod", "production":
		validEnv = EnvPrd
	}

	if validEnv == "" {
		return Environment{}, newError(fmt.Errorf("invalid environment: got = %s", env))
	}

	return Environment{validEnv}, nil
}

func (e Environment) String() string {
	return e.string
}

func (e Environment) IsLocal() bool {
	return e.string == EnvLocal
}

func (e Environment) IsDev() bool {
	return e.string == EnvDev
}

func (e Environment) IsStg() bool {
	return e.string == EnvStg
}

func (e Environment) IsPrd() bool {
	return e.string == EnvPrd
}
