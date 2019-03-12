package config

//IConfigurable ...
type IConfigurable interface {
	Validate() error
}

//IRuntimeConfigurable notifies the user when it changes
type IRuntimeConfigurable interface {
	IConfigurable
	Loaded()
	//Released()
}

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

func init() {
}
