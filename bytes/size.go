package bytesext

// Bytes is a type alias to int64 in order to better express the desired data type.
type Bytes = int64

// Common byte unit sizes
const (
	BYTE = 1

	// Decimal (Powers of 10 for Humans)
	KB = 1000 * BYTE
	MB = 1000 * KB
	GB = 1000 * MB
	TB = 1000 * GB
	PB = 1000 * TB
	EB = 1000 * PB
	ZB = 1000 * EB
	YB = 1000 * ZB

	// Binary (Powers of 2 for Computers)
	KiB = 1024 * BYTE
	MiB = 1024 * KiB
	GiB = 1024 * MiB
	TiB = 1024 * GiB
	PiB = 1024 * TiB
	EiB = 1024 * PiB
	ZiB = 1024 * EiB
	YiB = 1024 * ZiB
)
