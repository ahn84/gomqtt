package events

import (
	"net"

	"github.com/ahn84/gomqtt/server/internal/packets"
)

// Events provides callback handlers for different event hooks.
type Events struct {
	OnProcessMessage // published message receieved before evaluation.
	OnMessage        // published message receieved.
	OnError          // server error.
	OnConnect        // client connected.
	OnDisconnect     // client disconnected.
	OnSubscribe      // topic subscription created.
	OnUnsubscribe    // topic subscription removed.
}

// Packets is an alias for packets.Packet.
type Packet packets.Packet

// ClientLike contains limited information about a connected client.
type ClientLike interface {
	GetID() string
	GetRemote() string
	GetListener() string
	GetUsername() []byte
	SetUsername([]byte)
	GetCleanSession() bool
	GetConn() net.Conn
}

// OnProcessMessage is called when a publish message is received, allowing modification
// of the packet data after ACL checking has occurred but before any data is evaluated
// for processing - e.g. for changing the Retain flag. Note, this hook is ONLY called
// by connected client publishers, it is not triggered when using the direct
// s.Publish method. The function receives the sent message and the
// data of the client who published it, and allows the packet to be modified
// before it is dispatched to subscribers. If no modification is required, return
// the original packet data. If an error occurs, the original packet will
// be dispatched as if the event hook had not been triggered.
// This function will block message dispatching until it returns. To minimise this,
// have the function open a new goroutine on the embedding side.
// The `mqtt.ErrRejectPacket` error can be returned to reject and abandon any further
// processing of the packet.
type OnProcessMessage func(ClientLike, Packet) (Packet, error)

// OnMessage function is called when a publish message is received. Note,
// this hook is ONLY called by connected client publishers, it is not triggered when
// using the direct s.Publish method. The function receives the sent message and the
// data of the client who published it, and allows the packet to be modified
// before it is dispatched to subscribers. If no modification is required, return
// the original packet data. If an error occurs, the original packet will
// be dispatched as if the event hook had not been triggered.
// This function will block message dispatching until it returns. To minimise this,
// have the function open a new goroutine on the embedding side.
type OnMessage func(ClientLike, Packet) (Packet, error)

// OnConnect is called when a client successfully connects to the broker.
type OnConnect func(ClientLike, Packet)

// OnDisconnect is called when a client disconnects to the broker. An error value
// is passed to the function if the client disconnected abnormally, otherwise it
// will be nil on a normal disconnect.
type OnDisconnect func(ClientLike, error)

// OnError is called when errors that will not be passed to
// OnDisconnect are handled by the server.
type OnError func(ClientLike, error)

// OnSubscribe is called when a new subscription filter for a client is created.
type OnSubscribe func(filter string, cl ClientLike, qos byte)

// OnUnsubscribe is called when an existing subscription filter for a client is removed.
type OnUnsubscribe func(filter string, cl ClientLike)
