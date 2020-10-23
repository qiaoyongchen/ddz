package table

import (
	"ddz/message"
)

type processor interface {
	process(*Table, message.Message)
}

type processorFunc func(*Table, message.Message)

func (p processorFunc) process(tab *Table, msg message.Message) {
	p(tab, msg)
}

type processorMiddleware func(processor) processor
