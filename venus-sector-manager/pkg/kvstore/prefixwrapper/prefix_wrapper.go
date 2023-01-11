package prefixwrapper

import (
	"context"
	"fmt"

	"github.com/ipfs-force-community/venus-cluster/venus-sector-manager/pkg/kvstore"
	"github.com/ipfs-force-community/venus-cluster/venus-sector-manager/pkg/logging"
)

var log = logging.New("kv").With("driver", "prefix-wrapper")

var _ kvstore.KVStore = (*WrappedKVStore)(nil)

func NewWrappedKVStore(prefix []byte, inner kvstore.KVStore) (*WrappedKVStore, error) {
	prefixLen := len(prefix)
	if prefixLen == 0 {
		return nil, fmt.Errorf("empty prefix is not allowed")
	}

	if prefix[prefixLen-1] != '/' {
		prefixLen++
		p := make([]byte, prefixLen)

		// copied must be the size of prefix, thus it is also the last index of the p
		copied := copy(p, prefix)
		p[copied] = '/'
		prefix = p
	}

	log.Debugw("kv wrapped", "prefix", string(prefix), "prefix-len", prefixLen)

	return &WrappedKVStore{
		prefix:    prefix,
		prefixLen: prefixLen,
		inner:     inner,
	}, nil
}

type WrappedKVStore struct {
	prefix    []byte
	prefixLen int
	inner     kvstore.KVStore
}

func (w *WrappedKVStore) makeKey(raw kvstore.Key) kvstore.Key {
	key := make(kvstore.Key, w.prefixLen+len(raw))
	copy(key[:w.prefixLen], w.prefix)
	copy(key[w.prefixLen:], raw)
	return key
}

func (w *WrappedKVStore) Get(ctx context.Context, key kvstore.Key) (kvstore.Val, error) {
	return w.inner.Get(ctx, w.makeKey(key))
}

func (w *WrappedKVStore) Has(ctx context.Context, key kvstore.Key) (bool, error) {
	return w.inner.Has(ctx, w.makeKey(key))
}

func (w *WrappedKVStore) View(ctx context.Context, key kvstore.Key, cb kvstore.Callback) error {
	return w.inner.View(ctx, w.makeKey(key), cb)
}

func (w *WrappedKVStore) Put(ctx context.Context, key kvstore.Key, val kvstore.Val) error {
	return w.inner.Put(ctx, w.makeKey(key), val)
}

func (w *WrappedKVStore) Del(ctx context.Context, key kvstore.Key) error {
	return w.inner.Del(ctx, w.makeKey(key))
}

func (w *WrappedKVStore) Scan(ctx context.Context, prefix kvstore.Prefix) (kvstore.Iter, error) {
	iter, err := w.inner.Scan(ctx, w.makeKey(prefix))
	if err != nil {
		return nil, err
	}

	return &WrappedIter{
		prefixLen: w.prefixLen,
		inner:     iter,
	}, nil
}

func (w *WrappedKVStore) Run(ctx context.Context) error {
	return nil
}

func (w *WrappedKVStore) Close(ctx context.Context) error {
	return nil
}

type WrappedIter struct {
	prefixLen int
	inner     kvstore.Iter
}

func (wi *WrappedIter) Next() bool { return wi.inner.Next() }

func (wi *WrappedIter) View(ctx context.Context, cb kvstore.Callback) error {
	return wi.inner.View(ctx, cb)
}

func (wi *WrappedIter) Close() { wi.inner.Close() }

func (wi *WrappedIter) Key() kvstore.Key {
	key := wi.inner.Key()
	if len(key) > wi.prefixLen {
		key = key[wi.prefixLen:]
	}

	return key
}
