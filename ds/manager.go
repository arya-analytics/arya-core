package ds

import "fmt"

func NewConnManager(cfg Config) *ConnManager {
	return &ConnManager{
		cfg: cfg,
		Adapters: []ConnAdapter{},
	}
}

type ConnManager struct {
	cfg      Config
	Adapters []ConnAdapter
}

func (cm *ConnManager) getParams(configKey string) ConnParams {
	return cm.cfg[configKey]
}

func (cm *ConnManager) GetOrCreate(configKey string) Conn {
	ca, notFound := cm.get(configKey)
	if notFound == false {
		fmt.Println(cm.Adapters)
		return ca.conn
	}
	newCa := cm.create(configKey)
	cm.addAdapter(newCa)
	return newCa.conn
}

func (cm *ConnManager) addAdapter(ca ConnAdapter) {
	cm.Adapters = append(cm.Adapters, ca)
}

func (cm *ConnManager) create(configKey string) ConnAdapter {
	params := cm.getParams(configKey)
	connector := getConnector(params.Engine)
	ca := NewConnAdapter(configKey, params, connector)
	return ca
}

func (cm *ConnManager) get(configKey string) (ca ConnAdapter, notFound bool) {
	for i := range cm.Adapters {
		if cm.Adapters[i].configKey == configKey {
			return cm.Adapters[i], false
		}
	}
	return ConnAdapter{}, true
}
