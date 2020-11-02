package proc1

import (
	"ddz/game/message"
)

// Processor 处理
type Processor interface {
	Process(message.Message)
}

// ProcessorFunc 处理函数
type ProcessorFunc func(message.Message)

// Process Process
func (p ProcessorFunc) Process(msg message.Message) {
	p(msg)
}

// ProcessorMiddleware 处理器中间件
type ProcessorMiddleware func(Processor) Processor
