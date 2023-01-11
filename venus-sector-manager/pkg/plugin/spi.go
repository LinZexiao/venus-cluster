package plugin

import (
	"context"
	"fmt"
	"reflect"
	"unsafe"
)

const (
	// ManifestSymbol defines VSM plugin's entrance symbol.
	// Plugin take manifest info from this symbol.
	ManifestSymbol = "PluginManifest"
)

var (
	ErrInvalidPluginManifest = fmt.Errorf("invalid plugin manifest")
)

// Kind presents the kind of plugin.
type Kind uint8

const (
	// KVStore indicates it is a KVStore plugin.
	KVStore Kind = 1 + iota
	// ObjStore indicates it is a ObjStore plugin.
	ObjStore
)

type Manifest struct {
	Name        string
	Description string
	BuildTime   string
	// OnInit defines the plugin init logic.
	// it will be called after domain init.
	// return error will stop load plugin process and VSM startup.
	OnInit func(ctx context.Context, manifest *Manifest) error
	// OnShutDown defines the plugin cleanup logic.
	// return error will write log and continue shutdown.
	OnShutdown func(ctx context.Context, manifest *Manifest) error

	Kind Kind
}

// ExportManifest exports a manifest to VSM as a known format.
// it just casts sub-manifest to manifest.
func ExportManifest(m interface{}) *Manifest {
	v := reflect.ValueOf(m)
	return (*Manifest)(unsafe.Pointer(v.Pointer()))
}
