package pluginapp

/*
func loadConfigJSON(filename string, configs map[string]Config) (err error) {
	var m map[string]json.RawMessage
	if err = config.JSON.LoadFromFile(filename, &m); err != nil {
		return err
	}
	for k, v := range m {
		if cfg, ok := configs[k]; ok {
			if err = json.Unmarshal(v, cfg); err != nil {
				return fmt.Errorf("load %s plugin config: %v", k, err)
			}
		}
	}
	return nil
}

func loadConfigTOML(filename string, configs map[string]Config) (err error) {
	var m map[string]toml.Primitive
	if err = tomlcfg.TOML.LoadFromFile(filename, &m); err != nil {
		return err
	}
	for k, v := range m {
		if cfg, ok := configs[k]; ok {
			if err = toml.PrimitiveDecode(v, cfg); err != nil {
				return fmt.Errorf("load %s plugin config: %v", k, err)
			}
		}
	}
	return nil
}
*/
