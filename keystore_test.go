package encrepo

import (
	secrand "crypto/rand"
	"path/filepath"
	"sort"
	"testing"

	"github.com/ipfs/go-datastore"
	ci "github.com/libp2p/go-libp2p-core/crypto"
	"github.com/stretchr/testify/require"
)

func TestKeystoreFromSQLiteDatastore(t *testing.T) {
	// Instantiate datastore
	ds, err := NewSQLiteDatastore("sqlite3", filepath.Join(t.TempDir(), "db.sqlite"), "keys")
	require.NoError(t, err)

	// Generate keys
	keysIDs := []string{"a", "b", "c"}
	keys := map[string]ci.PrivKey{}
	for _, id := range keysIDs {
		sk, _, err := ci.GenerateEd25519Key(secrand.Reader)
		require.NoError(t, err)
		keys[id] = sk
	}

	// Create keystore
	const prefix = "keys"
	ks := KeystoreFromDatastore(ds, prefix)

	// Put keys
	for id, val := range keys {
		require.NoError(t, ks.Put(id, val))
	}

	// Put data with same prefix in ds
	require.NoError(t, ds.Put(datastore.NewKey(prefix+"_foo"), []byte("42")))

	// Check that the key list contains the correct keys and not the data with same prefix
	l, err := ks.List()
	require.NoError(t, err)
	l2 := make([]string, len(l))
	for i, k := range l {
		l2[i] = datastore.NewKey(k).Name()
	}
	sort.Strings(l2)
	require.Equal(t, keysIDs, l2)

	// Check that key data matches
	for id, val := range keys {
		v, err := ks.Get(id)
		require.NoError(t, err)
		require.Equal(t, val, v)
	}
}
