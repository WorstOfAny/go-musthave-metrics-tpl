package resetable_pool

import(
	"sync"
)

type Resetable interface {
	Reset()
}

type Pool[T Resetable] struct {
	current *PoolItem[T]
	mu sync.Mutex
	newObj func() T
}

type PoolItem[T Resetable] struct {
	next *PoolItem[T]
	value T
}


func New[T Resetable](objInitF func() T) *Pool[T] {
	p := &Pool[T]{newObj: objInitF}
	return p
}

func (p *Pool[T]) Get() T {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.current != nil {
		i := p.current
		p.current = i.next
		return i.value
	}

	return p.newObj()
}

func (p *Pool[T]) Put(item T) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.current != nil {
		i := p.current.next
		for {
			if i.next != nil {
				i = i.next
				continue
			}
			i.next = &PoolItem[T]{ value: item }
			break
		}
	}

	p.current = &PoolItem[T]{ value: item }
}
