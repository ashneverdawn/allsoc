package allsoc

import (
	"io"
)

// sockets stores a map of slices of *io.ReadWriter
var sockets map[string][]*io.ReadWriter
// SetupAllsoc sets up this library if not already setup.
func SetupAllsoc () {
    if sockets == nil { sockets = make(map[string][]*io.ReadWriter) }
}
// Socket allows socket operations.
type Socket struct {    
    rw *io.ReadWriter
    rooms map[string]int
}
// NewSocket creates a new Socket and returns a ptr to it.
func NewSocket(rw *io.ReadWriter) *Socket {
    return &Socket{rw : rw, rooms : make(map[string]int)}
}
// Join starts listening on a room
func (soc *Socket) Join(room string) {
    if _, ok := soc.rooms[room]; !ok {
        index := -1
        for i:= 0; i < len(sockets[room]); i++ {
            if sockets[room][i] == nil {
                sockets[room][i] = soc.rw
                index = i
                break
            }
        }
        if index == -1 {
            index = len(sockets[room])
            sockets[room] = append(sockets[room], soc.rw)
        }
        soc.rooms[room] = index
    }
}
// Leave stops listening on a room
func (soc *Socket) Leave(room string) {
    if index, ok := soc.rooms[room]; ok {
        soc.leave(room, index)
    }
}
// leave stops listening on a room. Called internally. Make sure to not call this on a room that was not joined
func (soc *Socket) leave(room string, index int) {
    sockets[room][index] = nil
    if sockets[room][len(sockets[room])-1] == nil {
        sockets[room] = sockets[room][:len(sockets[room])-1]
    }
}
// Read reads bytes from the socket
func (soc *Socket) Read(b []byte) (n int, err error) {
    n, err = (*soc.rw).Read(b)
    if err != nil {
        for room, index := range soc.rooms {
            soc.leave(room, index)
        }
    }
    return
}
// Emit send bytes to the client
func (soc *Socket) Write(b []byte) (n int, err error) {
    n, err = (*soc.rw).Write(b)
    return
}
// Broadcast sends bytes to listeners of the room
func (soc *Socket) Broadcast(room string, b []byte) (n int, err error) {
    index := soc.rooms[room]
    for i:= 0; i < len(sockets[room]); i++ {
        if sockets[room][i] != nil && i != index {
            n, err = (*sockets[room][i]).Write(b)
        }
    }
    return
}