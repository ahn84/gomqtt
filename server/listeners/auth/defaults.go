package auth

import "github.com/ahn84/gomqtt/server/events"

// Allow is an auth controller which allows access to all connections and topics.
type Allow struct{}

// Authenticate returns true if a username and password are acceptable. Allow always
// returns true.
func (a *Allow) Authenticate(cl events.ClientLike, password []byte) (*[]byte, bool) {
	return nil, true
}

// ACL returns true if a user has access permissions to read or write on a topic.
// Allow always returns true.
func (a *Allow) ACL(cl events.ClientLike, topic string, write bool) bool {
	return true
}

// Disallow is an auth controller which disallows access to all connections and topics.
type Disallow struct{}

// Authenticate returns true if a username and password are acceptable. Disallow always
// returns false.
func (d *Disallow) Authenticate(cl events.ClientLike, password []byte) (*[]byte, bool) {
	return nil, false
}

// ACL returns true if a user has access permissions to read or write on a topic.
// Disallow always returns false.
func (d *Disallow) ACL(cl events.ClientLike, topic string, write bool) bool {
	return false
}
