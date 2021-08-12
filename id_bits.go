package uid

import "fmt"

const (
	signBits = 1
)

// id bits 负责构造解析uid
/*
*
 * +------+----------------------+----------------+-----------+
 * | sign |     delta seconds    | worker node id | sequence  |
 * +------+----------------------+----------------+-----------+
 *   1bit          28bits              22bits         13bits
 * }
*/
type IdBits struct {
	/**
	 * Bits for [sign-> second-> workId-> sequence]
	 */
	TimestampBits int
	WorkerIdBits  int
	SequenceBits  int
	/**
	 * Max value for workId & sequence
	 */
	maxDeltaSeconds int64
	maxWorkerId     int64
	maxSequence     int64
	/**
	 * Shift for timestamp & workerId
	 */
	timestampShift int
	workerIdShift  int
}

func NewIdBits(timestampBits, workerIdBits, sequenceBits int) *IdBits {
	totalBits := signBits + timestampBits + workerIdBits + sequenceBits
	if totalBits != 64 {
		panic("id is not equal 64 bits")
	}

	return &IdBits{
		TimestampBits: timestampBits,
		WorkerIdBits:  workerIdBits,
		SequenceBits:  sequenceBits,

		maxDeltaSeconds: ^(int64(-1) << timestampBits),
		maxWorkerId:     ^(int64(-1) << workerIdBits),
		maxSequence:     ^(int64(-1) << sequenceBits),

		timestampShift: workerIdBits + sequenceBits,
		workerIdShift:  sequenceBits,
	}
}

func (i *IdBits) allocID(deltaSeconds int64, workerId int64, sequence int64) (id int64) {
	return (deltaSeconds << i.timestampShift) | (workerId << i.workerIdShift) | sequence
}

func (i *IdBits) parseID(uid int64) string {
	return fmt.Sprintf("deltaSeconds %v workerId %v sequence %v",
		i.second(uid),
		(uid<<(i.TimestampBits+signBits))>>(64-i.WorkerIdBits),
		i.sequence(uid),
	)
}

func (i *IdBits) second(uid int64) int64 {
	return uid >> i.timestampShift
}

func (i *IdBits) sequence(uid int64) int64 {
	return uid & i.maxSequence
}

func (i *IdBits) overSequence(seq int64) bool {
	return seq > int64(i.maxSequence)
}

var idBits = NewIdBits(40, 7, 16)

func SetIdBits(i *IdBits) {
	idBits = i
}

func AllocID(deltaSeconds int64, workerId int64, sequence int64) int64 {
	return idBits.allocID(deltaSeconds, workerId, sequence)
}

func ParseID(uid int64) string {
	return idBits.parseID(uid)
}

func Second(uid int64) int64 {
	return idBits.second(uid)
}

func MaxSecond() int64 {
	return int64(idBits.maxDeltaSeconds)
}

func Sequence(uid int64) int64 {
	return idBits.sequence(uid)
}

func MaxSequence() int64 {
	return int64(idBits.maxSequence)
}

func OverMaxSequence(seq int64) bool {
	return idBits.overSequence(seq)
}
