//go:generate -command asset go run asset.go
//go:generate asset -var netjs net.js

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	_ "net/http/pprof"
	"sync"

	"golang.org/x/net/websocket"
)

func main() {
	addr := flag.String("addr", "0.0.0.0:0", "address to bind to")

	flag.Parse()

	l, err := net.Listen("tcp", *addr)
	if err != nil {
		log.Fatalln("Error creating listener:", err)
	}

	defer l.Close()

	log.Println("Starting server at http://" + l.Addr().String())

	err = http.Serve(l, nil)
	log.Println("Fatal error:", err)
}

func init() {
	http.HandleFunc("/", IndexHandler)
	http.Handle("/sock", websocket.Handler(SocketHandler))
}

func js(a asset) asset { http.Handle("/"+a.Name, a); return a }

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	io.WriteString(w, `<script src="net.js"></script>`)
}

type Packet struct {
	Audio   string `json:",omitempty"`
	User    string `json:",omitempty"`
	Special string `json:",omitempty"`
}

func SocketReader(u *User, ws *websocket.Conn, input chan<- Packet) {
	defer close(input)

	dec := json.NewDecoder(ws)
	for {
		var p Packet
		err := dec.Decode(&p)
		if err != nil {
			u.Log("read:", err)
			return
		}
		input <- p
	}
}

func SocketWriter(u *User, ws *websocket.Conn, output <-chan Packet) {
	defer ws.Close()

	enc := json.NewEncoder(ws)
	for p := range output {
		err := enc.Encode(&p)
		if err != nil {
			u.Log("write:", err)
			return
		}
	}
}

func SocketHandler(ws *websocket.Conn) {
	u := NewUser(ws.Request().RemoteAddr)
	defer u.Remove()

	input := make(chan Packet)
	output := make(chan Packet)
	defer close(output)
	go SocketReader(u, ws, input)
	go SocketWriter(u, ws, output)

	u.Log("connected")

	for {
		select {
		case p, ok := <-input:
			if !ok {
				return
			}
			if p.Special != "" {
				continue
			}
			if p.User != "" {
				continue
			}
			if len([]rune(p.Audio)) != 4096*4 {
				continue
			}
			p.User = u.Addr
			go EachUser(func(o *User) {
				if o != u {
					o.Send <- p
				}
			})

		case p, ok := <-u.Send:
			if !ok {
				return
			}

			output <- p
		}
	}
}

type User struct {
	Addr string
	Send chan Packet
}

var users = make(map[*User]struct{})
var usersMtx sync.RWMutex

func NewUser(addr string) *User {
	u := &User{
		Addr: addr,
		Send: make(chan Packet),
	}

	p := Packet{
		User:    addr,
		Special: "connected",
	}
	go EachUser(func(o *User) {
		if u != o {
			o.Send <- p
			u.Send <- Packet{
				User:    o.Addr,
				Special: "connected",
			}
		}
	})

	usersMtx.Lock()
	users[u] = struct{}{}
	usersMtx.Unlock()

	return u
}

func (u *User) Remove() {
	usersMtx.Lock()
	delete(users, u)
	usersMtx.Unlock()

	p := Packet{
		User:    u.Addr,
		Special: "disconnected",
	}
	EachUser(func(o *User) {
		o.Send <- p
	})
}

func (u *User) Log(v ...interface{}) {
	log.Println("User/"+u.Addr, fmt.Sprintln(v...))
}

func (u *User) Logf(format string, v ...interface{}) {
	log.Println("User/"+u.Addr, fmt.Sprintf(format, v...))
}

func EachUser(f func(*User)) {
	usersMtx.RLock()
	for u := range users {
		usersMtx.RUnlock()
		f(u)
		usersMtx.RLock()
	}
	usersMtx.RUnlock()
}
