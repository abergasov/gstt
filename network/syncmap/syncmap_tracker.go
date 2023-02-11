package syncmap

import (
	"container/list"
	"sync"

	"github.com/abergasov/gstt/network"
)

type Tracker struct {
	length   int
	listMU   *sync.RWMutex
	list     *list.List
	listMap  sync.Map
	messages []*network.Message
}

func NewMessageTracker(length int) network.MessageTracker {
	return &Tracker{
		length:  length,
		listMU:  &sync.RWMutex{},
		listMap: sync.Map{},
		list:    list.New(),
	}
}

// Add will add a message to the tracker, deleting the oldest message if necessary
func (t *Tracker) Add(message *network.Message) (err error) {
	// if message already exists, do nothing
	if _, ok := t.listMap.Load(message.ID); ok {
		return nil
	}

	t.listMU.Lock()
	defer t.listMU.Unlock()
	// if list is full, remove oldest message
	if t.list.Len() >= t.length {
		oldest := t.list.Front()
		t.list.Remove(oldest)
		// delete from map
		// as only this method work with list, so we can use unsafe cast
		t.listMap.Delete(oldest.Value.(*network.Message).ID)
	}

	t.list.PushBack(message)
	t.listMap.Store(message.ID, t.list.Back())
	t.messages = nil

	return nil
}

// Delete will delete message from tracker
func (t *Tracker) Delete(id string) (err error) {
	if e, ok := t.listMap.Load(id); ok {
		t.listMU.Lock()
		defer t.listMU.Unlock()
		t.list.Remove(e.(*list.Element))
		t.listMap.Delete(id)
		t.messages = nil
		return nil
	}
	return network.ErrMessageNotFound
}

// Message Get returns a message for a given ID.  Message is retained in tracker
func (t *Tracker) Message(id string) (message *network.Message, err error) {
	if e, ok := t.listMap.Load(id); ok {
		t.listMU.RLock()
		defer t.listMU.RUnlock()
		// double unsafe cast. but as only this method work with list, so we can use unsafe cast
		return e.(*list.Element).Value.(*network.Message), nil
	}

	return nil, network.ErrMessageNotFound
}

// Messages returns messages in FIFO order
func (t *Tracker) Messages() (messages []*network.Message) {
	t.listMU.RLock()
	defer t.listMU.RUnlock()
	if t.messages != nil {
		return t.messages
	}

	messages = make([]*network.Message, 0, t.list.Len())
	for e := t.list.Front(); e != nil; e = e.Next() {
		messages = append(messages, e.Value.(*network.Message))
	}
	t.messages = messages
	return messages
}
