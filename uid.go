package uid

// workid的获取
// id生成(时钟回绕 借用未来时间，sleep)
// id获取（ring buffer，double buffer， 正常获取）
// 获取id生成器
type IdGenerate interface {
	NextID() int64
	String() string
}

var idManager = make(map[string]IdGenerate)

func Register(ig ...IdGenerate) {
	for _, i := range ig {
		idManager[i.String()] = i
	}
}

func UnRegister(name string) {
	delete(idManager, name)
}

func IG(name string) IdGenerate {
	return idManager[name]
}

func NextID() int64 {
	return IG("default").NextID()
}

func init() {
	Register(NewDefaultID(defaultWorker(), defaultStartTimeNano))
}

// ring buffer 优化思路 避免false sharing， 当处于一个cache line时候禁止写入，
// option 生成id的方式（ringbuffer, double buffer），获取的方式
