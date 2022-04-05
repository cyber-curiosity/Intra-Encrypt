package node

import (
	"context"
	"encoding/binary"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	datastore "github.com/ipfs/go-datastore"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	dht "github.com/libp2p/go-libp2p-kad-dht"
)

const Proto = "/mirage/0.0.1"

func StreamHandler(stream network.Stream) {
	var packet = make([]byte, 1420)
	var packetsize = make([]byte, 2)

	// Read the incoming packet size
	_, err := stream.Read(packetsize)
	if err != nil {
		stream.Close()
		return
	}

	size := binary.LittleEndian.Uint16(packetsize)

	fmt.Println("Incoming Packet: Size:" + fmt.Sprint(size))

	// Read in the packet

	var plen uint16 = 0
	for plen < size {
		tmp, err := stream.Read(packet[plen:size])
		plen += uint16(tmp)
		if err != nil {
			stream.Close()
			return
		}
	}
	fmt.Println("Read incoming packet")
	// Add stream write LN: 358
}

func CraftNode(
	ctx context.Context,
	port int,
	handler network.StreamHandler) (node host.Host, dhtOut *dht.Ipfs, err error) {
	// Create ListAddrStrings
	tcpip4 := fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", port)
	node, err = libp2p.New(
		libp2p.ListenAddrStrings(tcpip4),
		libp2p.DefaultSecurity,
	)
	if err != nil {
		fmt.Println("[!] Error creating node!")
		return
	}

	// Print the node's peer info in multiaddr format
	peerinfo := peer.AddrInfo{
		ID:    node.ID(),
		Addrs: node.Addrs(),
	}
	addrs, err := peer.AddrInfoToP2pAddrs(&peerinfo)
	fmt.Println("libp2p node addr:", addrs[0])

	//peerLookupTable := make(map[string]peer.ID)

	// Setup Stream Handler
	node.SetStreamHandler(Proto, handler)

	// Create DHT Subsystem
	dhtOut = dht.NewDHTClient(ctx, node, datastore.NewMapDatastore())

	return node, dhtOut, nil
}

func ExitSignal(node host.Host) {
	// Create channel to wait for termination signal
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch

	// Close the node
	err := node.Close()
	if err != nil {
		fmt.Println("[!] Error shutting down node. May need manual termination...")
	}
	fmt.Println("Exit signal received. Shutting down...")

	// Close the app
	os.Exit(0)
}
