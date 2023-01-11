package main

import (
	"bytes"
	"context"
	"sync"

	"github.com/tidwall/btree"

	"github.com/ipfs-force-community/venus-cluster/venus-sector-manager/pkg/kvstore"
	"github.com/ipfs-force-community/venus-cluster/venus-sector-manager/pkg/plugin"
)

func OnInit(ctx context.Context, manifest *plugin.Manifest) error { return nil }

func Open(meta map[string]string) (kvstore.DB, error) {
	return &memdb{
		collections: make(map[string]*collection),
		lock:        sync.Mutex{},
	}, nil
}

type memdb struct {
	collections map[string]*collection
	lock        sync.Mutex
}

func (*memdb) Close(context.Context) error {
	return nil
}

func (m *memdb) OpenCollection(name string) (kvstore.KVStore, error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	kv, ok := m.collections[name]
	if !ok {
		kv = &collection{
			inner: btree.NewMap[string, kvstore.Val](32),
			lock:  sync.RWMutex{},
		}
		m.collections[name] = kv
	}
	return kv, nil
}

func (*memdb) Run(context.Context) error {
	return nil
}

type collection struct {
	inner *btree.Map[string, kvstore.Val]
	lock  sync.RWMutex
}

func (c *collection) Get(_ context.Context, k kvstore.Key) (kvstore.Val, error) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	v, ok := c.inner.Get(string(k))
	if !ok {
		return nil, kvstore.ErrKeyNotFound
	}
	return v, nil
}

func (c *collection) Has(ctx context.Context, k kvstore.Key) (has bool, err error) {
	_, err = c.Get(ctx, k)
	if err == kvstore.ErrKeyNotFound {
		has = false
		err = nil
	} else if err == nil {
		has = true
	}
	return
}

func (c *collection) View(ctx context.Context, k kvstore.Key, fn kvstore.Callback) error {
	v, err := c.Get(ctx, k)
	if err != nil {
		return err
	}
	return fn(v)
}

func (c *collection) Put(_ context.Context, k kvstore.Key, v kvstore.Val) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.inner.Set(string(k), v)
	c.inner.Iter()
	return nil
}

func (c *collection) Del(_ context.Context, k kvstore.Key) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.inner.Delete(string(k))
	return nil
}

func (c *collection) Scan(_ context.Context, prefix kvstore.Prefix) (kvstore.Iter, error) {
	innerIt := c.inner.Iter()
	if !innerIt.Seek(string(prefix)) {
		innerIt.Last()
	}
	return &iter{
		inner:  innerIt,
		prefix: prefix,
	}, nil
}

type iter struct {
	inner  btree.MapIter[string, kvstore.Val]
	prefix kvstore.Prefix
}

func (it *iter) Next() bool {
	if !it.inner.Next() {
		return false
	}
	return len(it.prefix) == 0 || bytes.HasPrefix(it.Key(), it.prefix)
}

func (it *iter) Key() kvstore.Key {
	return kvstore.Key(it.inner.Key())
}

func (it *iter) View(_ context.Context, fn kvstore.Callback) error {
	return fn(it.inner.Value())
}

func (it *iter) Close() {}
