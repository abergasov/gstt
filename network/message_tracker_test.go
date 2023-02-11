package network_test

import (
	"fmt"
	"testing"

	"github.com/abergasov/gstt/network"
	"github.com/abergasov/gstt/network/singlemap"
	"github.com/abergasov/gstt/network/syncmap"
	"github.com/stretchr/testify/assert"
)

const (
	length = 5
)

func generateMessage(n int) *network.Message {
	return &network.Message{
		ID:     fmt.Sprintf("someID%d", n),
		PeerID: fmt.Sprintf("somePeerID%d", n),
		Data:   []byte{0, 1, 1},
	}
}

func getTrackerMethods() []func(int) network.MessageTracker {
	return []func(int) network.MessageTracker{
		singlemap.NewMessageTracker,
		syncmap.NewMessageTracker,
	}
}

func TestMessageTracker_Add(t *testing.T) {
	for _, trackerCreator := range getTrackerMethods() {
		t.Run("add, get, then all messages", func(t *testing.T) {
			mt := trackerCreator(length)

			for i := 0; i < length; i++ {
				err := mt.Add(generateMessage(i))
				assert.NoError(t, err)

				msg, err := mt.Message(generateMessage(i).ID)
				assert.NoError(t, err)
				assert.NotNil(t, msg)
			}

			msgs := mt.Messages()
			assert.Equal(t, []*network.Message{
				generateMessage(0),
				generateMessage(1),
				generateMessage(2),
				generateMessage(3),
				generateMessage(4),
			}, msgs)
		})

		t.Run("add, get, then all messages, delete some", func(t *testing.T) {
			mt := trackerCreator(length)

			for i := 0; i < length; i++ {
				err := mt.Add(generateMessage(i))
				assert.NoError(t, err)

				msg, err := mt.Message(generateMessage(i).ID)
				assert.NoError(t, err)
				assert.NotNil(t, msg)
			}

			msgs := mt.Messages()
			assert.Equal(t, []*network.Message{
				generateMessage(0),
				generateMessage(1),
				generateMessage(2),
				generateMessage(3),
				generateMessage(4),
			}, msgs)

			for i := 0; i < length-2; i++ {
				err := mt.Delete(generateMessage(i).ID)
				assert.NoError(t, err)
			}

			msgs = mt.Messages()
			assert.Equal(t, []*network.Message{
				generateMessage(3),
				generateMessage(4),
			}, msgs)
		})

		t.Run("not full, with duplicates", func(t *testing.T) {
			mt := trackerCreator(length)

			for i := 0; i < length-1; i++ {
				err := mt.Add(generateMessage(i))
				assert.NoError(t, err)
			}
			for i := 0; i < length-1; i++ {
				err := mt.Add(generateMessage(length - 2))
				assert.NoError(t, err)
			}

			msgs := mt.Messages()
			assert.Equal(t, []*network.Message{
				generateMessage(0),
				generateMessage(1),
				generateMessage(2),
				generateMessage(3),
			}, msgs)
		})

		t.Run("not full, with duplicates from other peers", func(t *testing.T) {
			mt := trackerCreator(length)

			for i := 0; i < length-1; i++ {
				err := mt.Add(generateMessage(i))
				assert.NoError(t, err)
			}
			for i := 0; i < length-1; i++ {
				msg := generateMessage(length - 2)
				msg.PeerID = "somePeerID0"
				err := mt.Add(msg)
				assert.NoError(t, err)
			}

			msgs := mt.Messages()
			assert.Equal(t, []*network.Message{
				generateMessage(0),
				generateMessage(1),
				generateMessage(2),
				generateMessage(3),
			}, msgs)
		})
	}
}

func TestMessageTracker_Cleanup(t *testing.T) {
	for _, trackerCreator := range getTrackerMethods() {
		t.Run("overflow and cleanup", func(t *testing.T) {
			mt := trackerCreator(length)

			for i := 0; i < length*2; i++ {
				err := mt.Add(generateMessage(i))
				assert.NoError(t, err)
			}

			msgs := mt.Messages()
			assert.Equal(t, []*network.Message{
				generateMessage(5),
				generateMessage(6),
				generateMessage(7),
				generateMessage(8),
				generateMessage(9),
			}, msgs)
		})

		t.Run("overflow and cleanup with duplicate", func(t *testing.T) {
			mt := trackerCreator(length)

			for i := 0; i < length*2; i++ {
				err := mt.Add(generateMessage(i))
				assert.NoError(t, err)
			}

			for i := length; i < length*2; i++ {
				err := mt.Add(generateMessage(i))
				assert.NoError(t, err)
			}

			msgs := mt.Messages()
			assert.Equal(t, []*network.Message{
				generateMessage(5),
				generateMessage(6),
				generateMessage(7),
				generateMessage(8),
				generateMessage(9),
			}, msgs)
		})
	}
}

func TestMessageTracker_Delete(t *testing.T) {
	for _, trackerCreator := range getTrackerMethods() {
		t.Run("empty tracker", func(t *testing.T) {
			mt := trackerCreator(length)
			err := mt.Delete("bleh")
			assert.ErrorIs(t, err, network.ErrMessageNotFound)
		})
	}
}

func TestMessageTracker_Message(t *testing.T) {
	for _, trackerCreator := range getTrackerMethods() {
		t.Run("empty tracker", func(t *testing.T) {
			mt := trackerCreator(length)
			msg, err := mt.Message("bleh")
			assert.ErrorIs(t, err, network.ErrMessageNotFound)
			assert.Nil(t, msg)
		})
	}
}