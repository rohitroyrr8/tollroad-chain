package types

import "encoding/binary"

var _ binary.ByteOrder

const (
	// UserVaultKeyPrefix is the prefix to retrieve all UserVault
	UserVaultKeyPrefix = "UserVault/value/"
)

// UserVaultKey returns the store key to retrieve a UserVault from the index fields
func UserVaultKey(
	index string,
) []byte {
	var key []byte

	indexBytes := []byte(index)
	key = append(key, indexBytes...)
	key = append(key, []byte("/")...)

	return key
}
