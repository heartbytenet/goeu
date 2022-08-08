package goeu

import "sync/atomic"

type Promise struct {
	Res      *ApiExecuteResult
	Creation int64
	done     int32
	Callback chan byte
}

func (p *Promise) Init(res *ApiExecuteResult, ts int64, callback chan byte) *Promise {
	p.Res      = res
	p.Creation = ts
	p.Callback = callback
	return p
}

func (p *Promise) DoneSet(v bool) {
	var i int32
	if v {
		i = 1
	} else {
		i = 0
	}
	atomic.StoreInt32(&p.done, i)
}

func (p *Promise) Done() bool {
	return atomic.LoadInt32(&p.done) == 1
}

func (p *Promise) Ended(ts int64, limit int64) bool {
	return p.Done() || ((ts - p.Creation) >= limit)
}