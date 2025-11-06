package main

import (
	proto "Ricart_and_Agrawala_implementation/grpc" // adjust to your generated package path
	//"context"
	//"log"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net"
	"sync"
)

//TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>

type message struct {
	timestamp      int64
	answer_channel chan message
}

// used for each server for each node
// this is a server
type Node struct {
	proto.UnimplementedNodeServiceServer
	mutex     sync.Mutex // lock for lamport stuff
	state     int
	timestamp int64
	id        int
	port      int
	node_map  map[int]proto.NodeServiceClient
	queue     []chan message
}

var node_count = 5

func main() {
	var start_port = 8080
	for i := start_port; i < node_count+start_port; i++ {
		new_node := Node{
			timestamp: 0,
			port:      i,
			id:        i - (start_port - 1), // we want id to start at 1
			state:     0,
			node_map:  nil,
			queue:     nil,
		}
		go new_node.run()
	}
	for {

	}
}

func (n *Node) run() {
	// determine state
	// if state is...
	go n.receive()
	n.state = 0

	for {
		if n.state == 2 {
			n.setLeader()
		}
		n.enter()
	}
}

func (n *Node) enter() {
	// set state to 1, for wanted
	n.state = 1
	n.send()
}

func (n *Node) send() {

	msg := message{
		timestamp: n.timestamp,
		port:      n.port,
	}

	for _, c := range n.node_list {
		fmt.Printf("Sending message %v from %v to %v \n", msg.timestamp, n.id, p.id)
		// send msg to all other nodes
		c.send(msg)

		n.mutex.Lock()
		n.timestamp++
		n.mutex.Unlock()
	}

	for _, p := range nodes {
		if p.id == n.id {
			continue
		}

	}

	// wait to receive all "yes"#
	for i := 1; i < node_count; i++ {
		fmt.Printf("%v received %d yes \n", n.id, i)
		<-n.answer_channel
	}
	// Now has token
	n.state = 2
}

func (n *Node) receive() {
	for {
		rec_msg := <-n.req_channel
		if n.state == 2 || n.state == 1 && rec_msg.timestamp > n.timestamp {
			n.queue = append(n.queue, rec_msg.answer_channel)
		} else {
			rec_msg.answer_channel <- message{0, nil}
		}
	}
}

func (n *Node) exit() {
	for _, c := range n.queue {
		c <- message{0, nil}
	}
	n.queue = n.queue[:0]
}

func (n *Node) setLeader() {
	fmt.Printf("%v is the leader\n", n.id)
	// some critical process...
	n.state = 0
	n.exit()
}

// / Server stuff
// Should start server for a different address each time (8080, 8081...) one for each node
// Each "Node" is a server itself (As it is here it receives requests)
func (n *Node) startServer(port int) {
	addr := fmt.Sprintf(":%d", port) // so its not hardcoded anymore
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		// logging
	}

	// Makes new GRPC server
	grpc_server := grpc.NewServer()

	// Registers the grpc server with the System struct
	proto.RegisterNodeServiceServer(grpc_server, n)

	err = grpc_server.Serve(listener)
	if err != nil {
		// logging
	}

	// logging like, "server started at... for node..."
}

// connects to another client with their port
func (n *Node) connect(port int) {
	addr := fmt.Sprintf("localhost:%d", port) // so its not hardcoded anymore
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		// logging
	}
	defer conn.Close() // keeps connnection open

	// A mapping between port to an open connection on it. This is used to get a received port from send
	n.node_map[port] = proto.NewNodeServiceClient(conn)

	// Below is an example of how to call using the connection
	//client := proto.NewNodeServiceClient(conn)
	// message, err := client.Receive(context.Background())
	// if err != nil {
	// logging
	// }
}
