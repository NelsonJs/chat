package objpool

import (
	"sync"

	"github.com/NelsonJs/chat/internal/model"
)

type objPool struct {
	msgPool sync.Pool
}

var (
	pool *objPool
	once sync.Once
)

func Pool() *objPool {
	once.Do(func() {
		pool = &objPool{
			msgPool: sync.Pool{
				New: func() any {
					return &model.Message{}
				},
			},
		}
	})
	return pool
}

func (p *objPool) GetMessage() *model.Message {
	return p.msgPool.Get().(*model.Message)
}
