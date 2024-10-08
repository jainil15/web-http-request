package config

import (
	"log"
	"os"
	"reflect"
	"sync"
)

var lock = &sync.Mutex{}

type EnvConfig struct {
	Port string `env:"PORT"`
}

func MustGetEnv(key string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		log.Fatalln("Missing env variable:", key)
		return ""
	}
	return value
}

func getEnv() *EnvConfig {
	env := EnvConfig{}
	st := reflect.TypeOf(env)
	sv := reflect.ValueOf(&env).Elem()
	for i := 0; i < st.NumField(); i++ {
		field := st.Field(i)
		key := field.Tag.Get("env")
		value := MustGetEnv(key)
		sv.Field(i).SetString(value)
	}
	return &env
}

var env *EnvConfig

func initConfig() *EnvConfig {
	if env == nil {
		lock.Lock()
		defer lock.Unlock()
		if env == nil {
			log.Println("Init config")
			env = getEnv()
		} else {
			log.Println("Config already initialized")
		}
	} else {
		log.Println("Config already initialized")
	}
	return env
}

func GetConfig() *EnvConfig {
	return initConfig()
}

var Env = GetConfig()
