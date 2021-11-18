package ds

func NewConnPooler(configs Configs) *ConnPooler {
	return &ConnPooler{
		configs:  configs,
		Adapters: map[*ConnAdapter]bool{},
	}
}

type ConnPooler struct {
	configs  map[string] Config
	Adapters map[*ConnAdapter] bool
}

func (cp *ConnPooler) GetConfig(key string) Config {
	return cp.configs[key]
}

func (cp *ConnPooler) SetConfig(key string, config Config) {
	cp.configs[key] = config
}

func (cp *ConnPooler) DelConfig(key string) {
	delete(cp.configs, key)
}

func (cp *ConnPooler) GetOrCreate(key string) Conn {
	ca, notFound := cp.getAdapter(key)
	if notFound {
		newCa := cp.createAdapter(key)
		cp.addAdapter(newCa)
		return newCa.conn
	}
	return ca.conn
}

func (cp *ConnPooler) addAdapter(ca *ConnAdapter) {
	cp.Adapters[ca] = true
}

func (cp *ConnPooler) removeAdapter(ca *ConnAdapter) {
	delete(cp.Adapters, ca)
}

func (cp *ConnPooler) createAdapter(key string) *ConnAdapter {
	config := cp.GetConfig(key)
	connector := getConnector(config.Engine)
	ca := NewConnAdapter(key, config, connector)
	return &ca
}

func (cp *ConnPooler) getAdapter(key string) (ca *ConnAdapter, notFound bool) {
	for ca := range cp.Adapters {
		if ca.key == key {
			return ca, false
		}
	}
	return &ConnAdapter{}, true
}
