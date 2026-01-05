package config

type ConfigString string

func (cs ConfigString) IsAssigned() bool {
	return cs != ""
}

type ConfigBool bool

func (cb ConfigBool) IsAssigned() bool {
	return bool(cb)
}
