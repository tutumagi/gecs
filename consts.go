package gecs

import (
	"strconv"
)

// EntityID entity runtime id
type EntityID uint32

// DefaultPlaceholder default entity id
const DefaultPlaceholder EntityID = entityMask

func (e EntityID) toInt() int {
	return int(e)
}
func (e EntityID) String() string {
	if DefaultPlaceholder == e {
		return "^"
	}
	return strconv.FormatUint(uint64(e), 10)
}

// ComponentID component runtime id
type ComponentID uint32

// SingletonID singleton runtime id
type SingletonID uint32

// _EnttPageSize must be a power of two
const _EnttPageSize int = 1 << 15

// 4 because of entityID is uint32 which is place 4 bytes
const _EnttPerPage int = _EnttPageSize / 4

// EntityID for uint32 entityid

type _VersionType uint16
type _DifferenceType int64

const entityMask = 0xFFFFF // 20bits for entity number
const versionMask = 0xFFF  // 12bits for the version(reset in (0-4095))
const entityShift = 20     // 20bits for entity number

type _KVComIDInt map[ComponentID]int
type _KVStrInt map[string]int

type _Vars struct {
	id   int
	data interface{}
}
