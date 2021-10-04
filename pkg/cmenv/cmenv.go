package cmenv

type EnvType int

// Environment constants
const (
	EnvDev EnvType = iota + 1
	EnvStag
	EnvProd
)

var isDev = false
var notProd = true
var serviceName string
var env EnvType
var envHooks []func(env EnvType) error
var envNames = map[EnvType]string{
	EnvDev:  "dev",
	EnvStag: "stag",
	EnvProd: "prod",
}
var envValues = map[string]EnvType{
	"dev":  EnvDev,
	"stag": EnvStag,
	"prod": EnvProd,
}

func (e EnvType) String() string {
	return envNames[e]
}

func SetHook(fn func(env EnvType) error) {
	if env != 0 {
		panic("Already initialize environment")
	}
	envHooks = append(envHooks, fn)
}

func ServiceName() string {
	if serviceName == "" {
		return "unknown service"
	}
	return serviceName
}

func SetEnvironment(name, e string) EnvType {
	if env != 0 {
		panic("Already initialize environment")
	}
	if serviceName != "" {
		panic("already initialized")
	}
	serviceName = name

	env = envValues[e]
	switch env {
	case EnvDev:
		isDev = true

	case EnvStag:

	case EnvProd:
		notProd = false

	default:
		panic("invalid environment: " + e)
	}
	for _, hook := range envHooks {
		if err := hook(env); err != nil {
			panic(err)
		}
	}
	return env
}

func Env() EnvType {
	return env
}

func IsDev() bool {
	return isDev
}

func IsDevOrStag() bool {
	return notProd
}

func NotProd() bool {
	return notProd
}

func IsProd() bool {
	return !notProd
}
