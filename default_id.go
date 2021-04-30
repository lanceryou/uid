package uid

import (
	"sync/atomic"
	"time"
)

const (
	defaultStartTimeStr        = "2021-04-01T00:00:00.000Z"
	defaultStartTimeNano int64 = 1617206400000000000
)

// 默认的uid生成器
// cas 无锁分配
type DefaultUid struct {
	lastID     int64
	workerID   int64
	startEpoch int64
}

// 获取下一个id
func (i *DefaultUid) NextID() int64 {
	for {
		localLastID := atomic.LoadInt64(&i.lastID)
		lastSecond := Second(localLastID)
		seq := Sequence(localLastID)

		now := epoch(i.startEpoch)
		// at the same second, increase sequence
		if now == lastSecond {
			seq++
		} else if now > lastSecond {
			seq = 0
		}

		// Clock moved backwards or exceed max seq, wait the next second to generate uid
		if OverMaxSequence(seq) || now < lastSecond {
			time.Sleep(time.Duration(0xFFFFF - (time.Now().Unix() & 0xFFFFF)))
			continue
		}

		newID := AllocID(now, i.workerID, seq)
		if atomic.CompareAndSwapInt64(&i.lastID, localLastID, newID) {
			return newID
		}
		// cas fail, wait
		time.Sleep(time.Duration(20))
	}
}

func (i *DefaultUid) String() string {
	return "default"
}

func epoch(startEpoch int64) int64 {
	return (time.Now().UnixNano() - startEpoch) >> 20 & MaxSecond()
}

func NewDefaultID(workerID int64, startEpoch int64) *DefaultUid {
	return &DefaultUid{
		workerID:   workerID,
		startEpoch: startEpoch,
		lastID:     0,
	}
}
