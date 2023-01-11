package kvstore

import (
	"unsafe"

	"github.com/ipfs-force-community/venus-cluster/venus-sector-manager/pkg/kvstore"
	"github.com/ipfs-force-community/venus-cluster/venus-sector-manager/pkg/plugin"
)

// TODO(0x5459): docs
type KVStoreManifest struct {
	plugin.Manifest

	Constructor func(meta map[string]string) (kvstore.DB, error)
}

// DeclareKVStoreManifest declares manifest as KVStoreManifest.
func DeclareKVStoreManifest(m *plugin.Manifest) *KVStoreManifest {
	return (*KVStoreManifest)(unsafe.Pointer(m))
}
