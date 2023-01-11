package kvstore

import (
	"github.com/ipfs-force-community/venus-cluster/venus-sector-manager/pkg/kvstore"
	"github.com/ipfs-force-community/venus-cluster/venus-sector-manager/pkg/plugin"
	pluginkvstore "github.com/ipfs-force-community/venus-cluster/venus-sector-manager/pkg/plugin/kvstore"
)

var plog = kvstore.Log.With("driver", "plugin")

func OpenPluginDB(path string, meta map[string]string) (kvstore.DB, error) {

	dbPlugin, err := plugin.Load(path)
	if err != nil {
		// TODO(0x5459): error handing
		return nil, err
	}
	db := pluginkvstore.DeclareKVStoreManifest(dbPlugin.Manifest)
	plog.With("plugin_name", db.Name)
	return db.Constructor(meta)
}
