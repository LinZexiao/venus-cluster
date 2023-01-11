package main

import (
	objstore "github.com/ipfs-force-community/venus-objstore"

	"github.com/ipfs-force-community/venus-cluster/venus-sector-manager/pkg/objstore/filestore"
)

func Open(cfg objstore.Config) (objstore.Store, error) { // nolint: deadcode
	return filestore.Open(cfg, true)
}
