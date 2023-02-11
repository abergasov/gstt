package singlemap

import (
	"container/list"
	"sync"

	"github.com/abergasov/gstt/network"
)

type Tracker struct {
	length   int
	listMU   *sync.RWMutex
	list     *list.List
	listMap  map[string]*list.Element
	messages []*network.Message
}

func NewMessageTracker(length int) network.MessageTracker {
	return &Tracker{
		length:  length,
		listMU:  &sync.RWMutex{},
		listMap: make(map[string]*list.Element, length),
		list:    list.New(),
	}
}

// Add will add a message to the tracker, deleting the oldest message if necessary
func (t *Tracker) Add(message *network.Message) (err error) {
	t.listMU.Lock()
	defer t.listMU.Unlock()

	// if message already exists, do nothing
	if _, ok := t.listMap[message.ID]; ok {
		return nil
	}

	// if list is full, remove oldest message
	if t.list.Len() >= t.length {
		oldest := t.list.Front()
		t.list.Remove(oldest)
		// delete from map
		// as only this method work with list, so we can use unsafe cast
		delete(t.listMap, oldest.Value.(*network.Message).ID)
	}

	t.list.PushBack(message)
	t.listMap[message.ID] = t.list.Back()
	t.messages = nil

	return nil
}

// Delete will delete message from tracker
func (t *Tracker) Delete(id string) (err error) {
	t.listMU.Lock()
	defer t.listMU.Unlock()

	if e, ok := t.listMap[id]; ok {
		t.list.Remove(e)
		delete(t.listMap, id)
		t.messages = nil
		return nil
	}

	return network.ErrMessageNotFound
}

// Message Get returns a message for a given ID.  Message is retained in tracker
func (t *Tracker) Message(id string) (message *network.Message, err error) {
	t.listMU.RLock()
	defer t.listMU.RUnlock()

	if e, ok := t.listMap[id]; ok {
		return e.Value.(*network.Message), nil // as only this method work with list, so we can use unsafe cast
	}

	return nil, network.ErrMessageNotFound
}

// Messages returns all messages in the tracker.
func (t *Tracker) Messages() (messages []*network.Message) {
	t.listMU.Lock()
	defer t.listMU.Unlock()
	if t.messages != nil {
		return t.messages
	}

	messages = make([]*network.Message, 0, t.length)
	for e := t.list.Front(); e != nil; e = e.Next() {
		// unsafe cast is safe here, as only object methods work with list
		messages = append(messages, e.Value.(*network.Message))
	}
	t.messages = messages
	return messages
}
