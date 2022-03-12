package storage

import "github.com/arya-analytics/aryacore/pkg/util/config"

const (
	configLeaf       = "storage"
	configLeafObject = "object"
	configLeafMD     = "md"
	configLeafCache  = "cache"
)

type configTree struct {
	config.Tree
}

func ConfigTree() configTree {
	return configTree{Tree: config.Tree{Base: configLeaf}}
}

func (ct configTree) Object() config.Tree {
	return ct.SubTree(configLeafObject)
}

func (ct configTree) MetaData() config.Tree {
	return ct.SubTree(configLeafMD)
}

func (ct configTree) Cache() config.Tree {
	return ct.SubTree(configLeafCache)
}
