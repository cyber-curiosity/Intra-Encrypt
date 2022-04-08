package main

import (
	"context"
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	discover "ice.com/initial/discover_func"
	nodefunc "ice.com/initial/node_func"
	tun "ice.com/initial/tun_int"
	//"github.com/libp2p/go-libp2p-core/crypto"
)

var (
	tunDev *tun.TUN

	err error
	// RevLookup map[string]string

	activeStreams map[string]network.Stream
)

func main() {

	// var Up = cmd.Sub{
	// 	Name: 	"up",
	// 	Alias:	"up",
	// 	Short:	"Create and bring up interface",
	// 	Args:	&UpArgs{},
	// 	Flags:	&UpFlags{},
	// 	Run:	UpRun,
	// }

	//Create the TUN interface
	fmt.Println("[*]Creating TUN Interface...")

	tunDev, err = tun.New(
		"tun_",
		tun.Address("192.168.72.1/24"),
		tun.MTU(1420),
	)
	if err != nil {
		fmt.Println("[!] Error creating TUN interface...")
		os.Exit(1)
	}
	fmt.Println("[+] Successfully created the TUN interface")

	// Establsih system context
	//ctx := context.Background()

	fmt.Println("Setting up Node...")

	if err != nil {
		fmt.Println(err)
	}

	// Create a peer table
	peerTable := make(map[string]peer.ID)

	// Start node with default settings
	// node, err := libp2p.New(
	// 	libp2p.ListenAddrStrings("/ip4/0.0.0.0/tcp/7788"),
	// 	libp2p.Ping(false),
	// 	nodefunc.StreamHandler,
	// )
	ctx := context.Background()

	livenode, dht, err := nodefunc.CraftNode(
		ctx,
		7788,
		nodefunc.StreamHandler,
	)

	if err != nil {
		panic(err)
	}

	err = tunDev.Up()
	if err != nil {
		fmt.Println("[!] Error bringing up TUN device")
	}

	// Setup peer discovery via DHT
	fmt.Println("[+] Setting up Node Discovery...")
	go discover.Discover(ctx, livenode, dht, peerTable)

	// ch := make(chan os.Signal, 1)
	// signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	// <-ch

	// fmt.Println("Received signal! Shutting down...")

	// if err := livenode.Close(); err != nil {
	// 	panic(err)
	// }

	// Listen for packets on the created TUN interface

	go nodefunc.ExitSignal(livenode)

	// Create a map of active streams (connections to other peers)
	activeStreams = make(map[string]network.Stream)

	var packet = make([]byte, 1420)
	for {
		// Read a packet
		plen, err := tunDev.Inter.Read(packet)
		if err != nil {
			log.Println(err)
			continue
		}

		fmt.Println(plen)

		// Decode the destination address
		//fmt.Println("Decoding the packet")
		dest := net.IPv4(packet[16], packet[17], packet[18], packet[19]).String()

		fmt.Println("Packet Received on TUN - Dest:", dest)

		// Check if there is already an open stream to the destination peer
		stream, active := activeStreams[dest]
		if active {
			// Send the length of the packet - this is to verify full delivery
			err = binary.Write(stream, binary.LittleEndian, uint16(plen))
			if err == nil {
				// As long is there is no error writing the length - write the packet
				_, err = stream.Write(packet[:plen])
				if err == nil {
					continue
				}

			}
			// Handle an error writing the length
			stream.Close()
			delete(activeStreams, dest)
		}

		// See if the destination peer is a known peer
		if peer, known := peerTable[dest]; known {
			stream, err = livenode.NewStream(ctx, peer, nodefunc.Proto)
			if err != nil {
				continue
			}
			// Write packet length
			err = binary.Write(stream, binary.LittleEndian, uint16(plen))
			if err != nil {
				stream.Close()
				continue
			}
			// Write the packet
			_, err = stream.Write(packet[:plen])
			if err != nil {
				stream.Close()
				continue
			}

			// If all succeeds when writing the packet to the stream
			// we should reuse this stream by adding it active streams map.
			activeStreams[dest] = stream
		}

	}

}
