package config

import "github.com/spf13/viper"

func AccessKey() string {
	return viper.GetString("accesskey")
}

// XXX rename all this crap

func Host() string {
	return viper.GetString("host")
}

func PostHost() string {
	return viper.GetString("posthost")
}

var debug map[string]bool

// SetDebug sets the value of the specified boolean debugging flag.
// (from robpike.io/ivy/config/config.go)
func SetDebug(flag string, state bool) {
	if debug == nil {
		debug = make(map[string]bool)
	}
	debug[flag] = state
}

func Debug(flag string) bool {
	return debug[flag]
}
