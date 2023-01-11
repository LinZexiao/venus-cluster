package badger

import (
	"context"
	"fmt"
	"path/filepath"
	"sync"

	badgerdriver "github.com/dgraph-io/badger/v2"

	"github.com/ipfs-force-community/venus-cluster/venus-sector-manager/pkg/kvstore"
	"github.com/ipfs-force-community/venus-cluster/venus-sector-manager/pkg/logging"
)

var _ kvstore.KVStore = (*BadgerKVStore)(nil)

var blog = logging.New("kv").With("driver", "badger")

type blogger struct {
	*logging.ZapLogger
}

func (bl *blogger) Warningf(format string, args ...interface{}) {
	bl.ZapLogger.Warnf(format, args...)
}

func OpenBadger(basePath string) kvstore.DB {
	return &badgerDB{
		basePath: basePath,
		dbs:      make(map[string]*badgerdriver.DB),
		mu:       sync.Mutex{},
	}
}

type BadgerKVStore struct {
	db *badgerdriver.DB
}

func (b *BadgerKVStore) Get(ctx context.Context, key kvstore.Key) (kvstore.Val, error) {
	var val []byte
	err := b.db.View(func(txn *badgerdriver.Txn) error {
		switch item, err := txn.Get(key); err {
		case nil:
			val, err = item.ValueCopy(nil)
			return err

		case badgerdriver.ErrKeyNotFound:
			return kvstore.ErrKeyNotFound

		default:
			return fmt.Errorf("get value from badger: %w", err)
		}
	})

	if err != nil {
		return nil, err
	}

	return val, nil
}

func (b *BadgerKVStore) Has(ctx context.Context, key kvstore.Key) (bool, error) {
	err := b.db.View(func(txn *badgerdriver.Txn) error {
		_, err := txn.Get(key)
		return err
	})

	switch err {
	case badgerdriver.ErrKeyNotFound:
		return false, nil
	case nil:
		return true, nil
	default:
		return false, fmt.Errorf("failed to check if block exists in badger blockstore: %w", err)
	}
}

func (b *BadgerKVStore) View(ctx context.Context, key kvstore.Key, cb kvstore.Callback) error {
	return b.db.View(func(txn *badgerdriver.Txn) error {
		switch item, err := txn.Get(key); err {
		case nil:
			return item.Value(cb)

		case badgerdriver.ErrKeyNotFound:
			return kvstore.ErrKeyNotFound

		default:
			return fmt.Errorf("get value from badger: %w", err)
		}
	})

}

func (b *BadgerKVStore) Put(ctx context.Context, key kvstore.Key, val kvstore.Val) error {
	return b.db.Update(func(txn *badgerdriver.Txn) error {
		return txn.Set(key, val)
	})
}

func (b *BadgerKVStore) Del(ctx context.Context, key kvstore.Key) error {
	return b.db.Update(func(txn *badgerdriver.Txn) error {
		return txn.Delete(key)
	})
}

func (b *BadgerKVStore) Scan(ctx context.Context, prefix kvstore.Prefix) (kvstore.Iter, error) {
	txn := b.db.NewTransaction(false)
	iter := txn.NewIterator(badgerdriver.DefaultIteratorOptions)

	return &BadgerIter{
		txn:    txn,
		iter:   iter,
		seeked: false,
		valid:  false,
		prefix: prefix,
	}, nil
}

var _ kvstore.Iter = (*BadgerIter)(nil)

type BadgerIter struct {
	txn  *badgerdriver.Txn
	iter *badgerdriver.Iterator
	item *badgerdriver.Item

	seeked bool
	valid  bool
	prefix []byte
}

func (bi *BadgerIter) Next() bool {
	if bi.seeked {
		bi.iter.Next()
	} else {
		if len(bi.prefix) == 0 {
			bi.iter.Rewind()
		} else {
			bi.iter.Seek(bi.prefix)
		}
		bi.seeked = true
	}

	if len(bi.prefix) == 0 {
		bi.valid = bi.iter.Valid()
	} else {
		bi.valid = bi.iter.ValidForPrefix(bi.prefix)
	}

	if bi.valid {
		bi.item = bi.iter.Item()
	}

	return bi.valid
}

func (bi *BadgerIter) Key() kvstore.Key {
	if !bi.valid {
		return nil
	}

	return bi.item.Key()
}

func (bi *BadgerIter) View(ctx context.Context, cb kvstore.Callback) error {
	if !bi.valid {
		return kvstore.ErrIterItemNotValid
	}

	return bi.item.Value(cb)
}

func (bi *BadgerIter) Close() {
	bi.iter.Close()
	bi.txn.Discard()
}

var _ kvstore.DB = (*badgerDB)(nil)

type badgerDB struct {
	basePath string
	dbs      map[string]*badgerdriver.DB
	mu       sync.Mutex
}

func (db *badgerDB) Run(context.Context) error {
	return nil
}

func (db *badgerDB) Close(context.Context) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	var lastError error
	for k, innerDB := range db.dbs {
		if err := innerDB.Close(); err != nil {
			lastError = err
		}
		delete(db.dbs, k)
	}
	return lastError
}

func (db *badgerDB) OpenCollection(name string) (kvstore.KVStore, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	if innerDB, ok := db.dbs[name]; ok {
		return &BadgerKVStore{db: innerDB}, nil
	}
	path := filepath.Join(db.basePath, name)
	opts := badgerdriver.DefaultOptions(path).WithLogger(&blogger{blog.With("path", path)})
	innerDB, err := badgerdriver.Open(opts)
	if err != nil {
		return nil, fmt.Errorf("open sub badger %s, %w", name, err)
	}
	db.dbs[name] = innerDB
	return &BadgerKVStore{db: innerDB}, nil
}
