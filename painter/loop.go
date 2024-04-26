package painter

import (
	"image"
	"sync"

	"golang.org/x/exp/shiny/screen"
)

// Receiver отримує текстуру, яка була підготовлена в результаті виконання команд у циклі подій.
type Receiver interface {
	Update(t screen.Texture)
}

// Loop реалізує цикл подій для формування текстури отриманої через виконання операцій отриманих з внутрішньої черги.
type Loop struct {
	Receiver Receiver

	next screen.Texture // текстура, яка зараз формується
	prev screen.Texture // текстура, яка була відправлення останнього разу у Receiver

	mq messageQueue

	stop    chan struct{}
	stopReq bool
}

var size = image.Pt(400, 400)

// Start запускає цикл подій. Цей метод потрібно запустити до того, як викликати на ньому будь-які інші методи.
func (l *Loop) Start(s screen.Screen) {
	l.next, _ = s.NewTexture(size)
	l.prev, _ = s.NewTexture(size)

	go func() {
		for {
			op := l.mq.pull()
			if update := op.Do(l.next); update {
				l.Receiver.Update(l.next)
				l.next, l.prev = l.prev, l.next
			}

			if l.stopReq {
				break
			}
		}
	}()
}

// Post додає нову операцію у внутрішню чергу.
func (l *Loop) Post(op Operation) {
	if update := op.Do(l.next); update {
		l.Receiver.Update(l.next)
		l.next, l.prev = l.prev, l.next
	}
}

// StopAndWait сигналізує про необхідність завершити цикл та блокується до моменту його повної зупинки.
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

func (mq *messageQueue) push(op Operation) {
	mq.mutex.Lock()
	defer mq.mutex.Unlock()

	mq.data = append(mq.data, op)

	if mq.pushSignal != nil {
		close(mq.pushSignal)
		mq.pushSignal = nil
	}
}

func (mq *messageQueue) pull() Operation {
	mq.mutex.Lock()
	defer mq.mutex.Unlock()

	for len(mq.data) == 0 {
		mq.pushSignal = make(chan struct{})
		mq.mutex.Unlock()
		<-mq.pushSignal
		mq.mutex.Lock()
	}

	res := mq.data[0]
	mq.data[0] = nil
	mq.data = mq.data[1:]
	return res
}

func (mq *messageQueue) empty() bool {
	mq.mutex.Lock()
	defer mq.mutex.Unlock()

	return len(mq.data) == 0
}
