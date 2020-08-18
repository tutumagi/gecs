package gecs

import (
	"math"
	"strconv"
)

// EntityID entity runtime id
type EntityID uint32

func (e EntityID) String() string {
	if DefaultPlaceholder == e {
		return "^"
	}
	return strconv.FormatUint(uint64(e), 10)
}

// DefaultPlaceholder default entity id
const DefaultPlaceholder EntityID = math.MaxUint32

// ComponentID component runtime id
type ComponentID uint32

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
