package barcode_snowflake

// 12 length character id, 39bit length
// most code from https://raw.githubusercontent.com/bwmarrin/snowflake/master/snowflake.go
import (
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/ezaurum/cthulthu/generators"
	"strconv"
	"sync"
	"time"
)

const (
	// Epoch is set to the twitter snowflake epoch of Nov 04 2010 01:42:54 UTC
	// You may customize this to set a different epoch for your application.
	Epoch int64 = 1288834974657

	// Number of bits to use for barcodeStringGenerator
	// Remember, you have a total 12 bits to share between Node/Step
	NodeBits uint8 = 4

	// Number of bits to use for Step
	// Remember, you have a total 12 bits to share between barcodeStringGenerator/Step
	StepBits uint8 = 8

	// use only last 39bit
	lastMask int64 = -1 ^ (-1 << TotalBits)

	TotalBits uint8 = 39
)

var (
	nodeMax   int64 = -1 ^ (-1 << NodeBits)
	nodeMask  int64 = nodeMax << StepBits
	stepMask  int64 = -1 ^ (-1 << StepBits)
	timeShift uint8 = NodeBits + StepBits
)

// A JSONSyntaxError is returned from UnmarshalJSON if an invalid ID is provided.
type JSONSyntaxError struct{ original []byte }

func (j JSONSyntaxError) Error() string {
	return fmt.Sprintf("invalid snowflake ID %q", string(j.original))
}

// A barcodeStringGenerator struct holds the basic information needed for a snowflake generator
// node

// An ID is a custom type used for a snowflake ID.  This is used so we can
// attach methods onto the ID.
type ID int64

// NewbarcodeStringGenerator returns a new snowflake node that can be used to generate snowflake
// IDs

func New(node int64) (*barcodeStringGenerator, error) {

	// re-calc in case custom NodeBits or StepBits were set
	nodeMax = -1 ^ (-1 << NodeBits)
	nodeMask = nodeMax << StepBits
	stepMask = -1 ^ (-1 << StepBits)
	timeShift = NodeBits + StepBits

	if node < 0 || node > nodeMax {
		return nil, errors.New("barcodeStringGenerator number must be between 0 and " + strconv.FormatInt(nodeMax, 10))
	}

	return &barcodeStringGenerator{
		time: 0,
		node: node,
		step: 0,
	}, nil
}

// Generate creates and returns a unique snowflake ID
func (n *barcodeStringGenerator) GenerateID() ID {

	n.mu.Lock()

	now := time.Now().UnixNano() / 1000000

	if n.time == now {
		n.step = (n.step + 1) & stepMask

		if n.step == 0 {
			for now <= n.time {
				now = time.Now().UnixNano() / 1000000
			}
		}
	} else {
		n.step = 0
	}

	n.time = now

	r := ID(lastMask & ((now-Epoch)<<timeShift |
		(n.node << StepBits) | (n.step)))

	/*r := ID((now-Epoch)<<timeShift |
			(n.node << StepBits) | (n.step))*/

	n.mu.Unlock()
	return r
}

// Int64 returns an int64 of the snowflake ID
func (f ID) Int64() int64 {
	return int64(f)
}

// String returns a string of the snowflake ID
func (f ID) String() string {
	return strconv.FormatInt(int64(f), 10)
}

// Base2 returns a string base2 of the snowflake ID
func (f ID) Base2() string {
	return strconv.FormatInt(int64(f), 2)
}

// Base36 returns a base36 string of the snowflake ID
func (f ID) Base36() string {
	return strconv.FormatInt(int64(f), 36)
}

// Base64 returns a base64 string of the snowflake ID
func (f ID) Base64() string {
	return base64.StdEncoding.EncodeToString(f.Bytes())
}

// Bytes returns a byte slice of the snowflake ID
func (f ID) Bytes() []byte {
	return []byte(f.String())
}

// IntBytes returns an array of bytes of the snowflake ID, encoded as a
// big endian integer.
func (f ID) IntBytes() [8]byte {
	var b [8]byte
	binary.BigEndian.PutUint64(b[:], uint64(f))
	return b
}

// Time returns an int64 unix timestamp of the snowflake ID time
func (f ID) Time() int64 {
	return (int64(f) >> timeShift) + Epoch
}

// barcodeStringGenerator returns an int64 of the snowflake ID node number
func (f ID) Node() int64 {
	return int64(f) & nodeMask >> StepBits
}

// Step returns an int64 of the snowflake step (or sequence) number
func (f ID) Step() int64 {
	return int64(f) & stepMask
}

// MarshalJSON returns a json byte array string of the snowflake ID.
func (f ID) MarshalJSON() ([]byte, error) {
	buff := make([]byte, 0, 22)
	buff = append(buff, '"')
	buff = strconv.AppendInt(buff, int64(f), 10)
	buff = append(buff, '"')
	return buff, nil
}

// UnmarshalJSON converts a json byte array of a snowflake ID into an ID type.
func (f *ID) UnmarshalJSON(b []byte) error {
	if len(b) < 3 || b[0] != '"' || b[len(b)-1] != '"' {
		return JSONSyntaxError{b}
	}

	i, err := strconv.ParseInt(string(b[1:len(b)-1]), 10, 64)
	if err != nil {
		return err
	}

	*f = ID(i)
	return nil
}

var _ generators.IDGenerator = &barcodeStringGenerator{}

type barcodeStringGenerator struct {
	mu   sync.Mutex
	time int64
	node int64
	step int64
}

func (g *barcodeStringGenerator) Generate() string {
	return g.GenerateID().String()
}

func (g *barcodeStringGenerator) GenerateInt64() int64 {
	return int64(g.GenerateID())
}
