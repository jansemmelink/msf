package config

import (
	"fmt"

	"github.com/jansemmelink/msf/lib/log"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

//Register config
func Register(name string, data IConfigurable, doc string) {
	if err := cs.Register(name, data, doc); err != nil {
		panic(errors.Wrapf(err, "failed to register config(%s)", name))
	}
}

//Get named config
func Get(name string) (IConfigurable, error) {
	return cs.Get(name)
}

//MustGet is same as Get, but panics on error
func MustGet(name string) IConfigurable {
	c, err := Get(name)
	if err != nil {
		panic(fmt.Sprintf("Failed to get config \"%s\": %v", name, err))
	}
	return c
}

var cs = newConfigSet()

func newConfigSet() *configSet {
	cs := &configSet{
		all:  make(map[string]*config),
		dirs: []string{"./conf"},
	}
	cs.Register("log", &log.Config{}, "Configuration for process log levels")
	return cs
}

type configSet struct {
	all  map[string]*config
	dirs []string
}

func (cs configSet) AddDir(dir string) {
	cs.dirs = append(cs.dirs, dir)
}

//Register a configurable item which may/may not be configured
//this just tells the program what types of config it supports,
//and if you call this in a module's init() function, this
//config will be part of the generated documentation
func (cs *configSet) Register(name string, data IConfigurable, doc string) error {
	if len(name) < 1 {
		return fmt.Errorf("config.Register(%s, %T) without a name", name, data)
	}

	if _, ok := cs.all[name]; ok {
		return fmt.Errorf("config.Register(%s, %T) with duplicate name", name, data)
	}

	cs.all[name] = &config{
		data:        data,
		doc:         doc,
		rt:          nil,
		viperConfig: nil,
		loaded:      false,
	}
	log.Debugf("Registered config(%s)=%T", name, data)
	return nil
}

//RegisterRt registers config that can change at run-time
func (cs *configSet) RegisterRt(name string, data IRuntimeConfigurable, doc string) error {
	if err := cs.Register(name, data, doc); err != nil {
		return errors.Wrapf(err, "failed to add as config")
	}
	config, _ := cs.all[name]
	config.rt = data
	log.Debugf("Registered runtime config(%s)=%T", name, data)
	return nil
}

//Has checks if this is configured, without parsing and validating the config
func (cs *configSet) Has(name string) bool {
	log.Errorf("NYI configSet.Has()")
	return false
}

func (cs *configSet) Get(name string) (IConfigurable, error) {
	c, ok := cs.all[name]
	if !ok {
		return nil, fmt.Errorf("config(%s) not registered", name)
	}

	if !c.loaded {
		//first time using this config after registration:
		//read this config using a new copy of our viper instance:
		c.viperConfig = viper.New()
		for _, dir := range cs.dirs {
			c.viperConfig.AddConfigPath(dir)
		}
		c.viperConfig.SetConfigName(name)

		if err := c.viperConfig.ReadInConfig(); err != nil {
			return nil, errors.Wrapf(err, "failed to read config [%s]", name)
		}

		if err := c.viperConfig.Unmarshal(c.data); err != nil {
			return nil, errors.Wrapf(err, "failed to unmarshal config(%s)", name)
		}

		if err := c.data.Validate(); err != nil {
			return nil, errors.Wrapf(err, "invalid config(%s)", name)
		}

		c.loaded = true
		log.Debugf("Loaded config(%s): %+v", name, c.data)
	}

	//todo: return a copy of the config that won't change ever and that
	//the user can change without affecting others
	//...
	return c.data, nil
}

// 	/*
// 	 * Notify listeners of the config changes
// 	 */
// 	for _, changed := range publicConfig.changed {

// 		if oldConfig != nil {
// 			changed.Released(oldConfig.config)
// 		}

// 		changed.Loaded(newConfig.config)

// 	} // for each load

// 	/*
// 	 * Done
// 	 */
// 	log.Debugf("Successfully loaded config [%s]",
// 		publicConfig.configName)
// 	return nil

// 	data := nil

// 	var publicConfig Config
// 	publicConfig.init(
// 		configName,
// 		configType,
// 		viperConfig)

// 	/*
// 	 * Load the configuration into the user struct
// 	 */
// 	if err := publicConfig.loadConfig(
// 		true); err != nil {

// 		return nil, errors.Wrapf(err,
// 			"failed to load configuration file")

// 	} // if failed to load configuration

// 	/*
// 	 * Enable viper to notify us of config file changes
// 	 */
// 	viperConfig.OnConfigChange(func(in fsnotify.Event) {

// 		defer log.Sync()

// 		if in.Op&fsnotify.Write == fsnotify.Write {

// 			log.Debugf("Config file changed. Event [%s]",
// 				in)

// 			if err := publicConfig.loadConfig(
// 				false); err != nil {

// 				log.Errorf("%+v", errors.Wrap(err,
// 					"failed to load config"))

// 			} // if failed to load config

// 		} // if write

// 	}) // OnConfigChange

// 	viperConfig.WatchConfig()

// 	/*
// 	 * Done
// 	 */
// 	log.Debugf("Successfully initialised config [%s]",
// 		configName)
// 	return &publicConfig, nil

// }

type config struct {
	data        IConfigurable
	doc         string
	rt          IRuntimeConfigurable
	viperConfig *viper.Viper
	loaded      bool
}
