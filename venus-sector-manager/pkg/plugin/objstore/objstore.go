package objstore

import (
	"unsafe"

	"github.com/ipfs-force-community/venus-cluster/venus-sector-manager/pkg/objstore"
	"github.com/ipfs-force-community/venus-cluster/venus-sector-manager/pkg/plugin"
)

type ObjStoreManifest struct {
	plugin.Manifest

	Constructor func(cfg objstore.Config) (objstore.Store, error)
}

// DeclareObjStoreManifest declares manifest as ObjStoreManifest.
func DeclareObjStoreManifest(m *plugin.Manifest) *ObjStoreManifest {
	return (*ObjStoreManifest)(unsafe.Pointer(m))
}
