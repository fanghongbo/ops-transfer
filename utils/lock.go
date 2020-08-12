package utils

type Semaphore struct {
	bufSize int
	channel chan int8
}

func NewSemaphore(concurrencyNum int) *Semaphore {
	return &Semaphore{channel: make(chan int8, concurrencyNum), bufSize: concurrencyNum}
}

func (u *Semaphore) TryAcquire() bool {
	select {
	case u.channel <- int8(0):
		return true
	default:
		return false
	}
}

func (u *Semaphore) Acquire() {
	u.channel <- int8(0)
}

func (u *Semaphore) Release() {
	<-u.channel
}

func (u *Semaphore) AvailablePermits() int {
	return u.bufSize - len(u.channel)
}
