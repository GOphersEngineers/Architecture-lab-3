package painter

import (
	"image"
	"sync"

	"golang.org/x/exp/shiny/screen"
)

type Receiver interface {
	Update(t screen.Texture)
}

type Loop struct {
	Receiver Receiver
	next, prev screen.Texture
	mq messageQueue
	stop chan struct{}
	stopReq bool
}

var size = image.Pt(400, 400)

func (l *Loop) Start(s screen.Screen) {
	l.next, _ = s.NewTexture(size)
	l.prev, _ = s.NewTexture(size)

	go l.runEventLoop()
}

func (l *Loop) runEventLoop() {
	for {
		op := l.mq.pull()
		if update := op.Do(l.next); update {
			l.swapTextures()
		}

		if l.stopReq {
			break
		}
	}
}

func (l *Loop) swapTextures() {
	l.Receiver.Update(l.next)
	l.next, l.prev = l.prev, l.next
}

func (l *Loop) Post(op Operation) {
	if update := op.Do(l.next); update {
		l.swapTextures()
	}
}

func (l *Loop) StopAndWait() {
	l.Post(OperationFunc(func(t screen.Texture) {
		l.stopReq = true
	}))
	<-l.stop
}

type messageQueue struct {
	pushSignal chan struct{}
	mutex sync.Mutex
	data  []Operation
}

func (mq *messageQueue) signalPush() {
	if mq.pushSignal != nil {
		close(mq.pushSignal)
		mq.pushSignal = nil
	}
}

func (mq *messageQueue) push(op Operation) {
	mq.mutex.Lock()
	defer mq.mutex.Unlock()

	mq.data = append(mq.data, op)
	mq.signalPush()
}

func (mq *messageQueue) pull() Operation {
	mq.mutex.Lock()
	defer mq.mutex.Unlock()

	mq.waitUntilNotEmpty()
	res := mq.pop()

	return res
}

func (mq *messageQueue) waitUntilNotEmpty() {
	for mq.isEmpty() {
		mq.pushSignal = make(chan struct{})
		mq.mutex.Unlock()
		<-mq.pushSignal
		mq.mutex.Lock()
	}
}

func (mq *messageQueue) pop() Operation {
	res := mq.data[0]
	mq.data[0] = nil
	mq.data = mq.data[1:]
	return res
}

func (mq *messageQueue) isEmpty() bool {
	return len(mq.data) == 0
}
