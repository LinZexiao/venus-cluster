package plugin

import (
	goplugin "plugin"
)

type Plugin struct {
	*Manifest
	library *goplugin.Plugin
}

func Load(path string) (plugin *Plugin, err error) {
	plugin = &Plugin{}
	plugin.library, err = goplugin.Open(path)
	if err != nil {
		// TODO(0x5459): error handing
		return
	}
	manifestSym, err := plugin.library.Lookup(ManifestSymbol)
	if err != nil {
		// TODO(0x5459): error handing
		return
	}
	var ok bool
	manifestFn, ok := manifestSym.(func() *Manifest)
	if !ok {
		err = ErrInvalidPluginManifest
		return
	}
	plugin.Manifest = manifestFn()
	return
}
