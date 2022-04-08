package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	nodefunc "ice.com/initial/node_func"
	tun "ice.com/initial/tun_int"
	//"github.com/libp2p/go-libp2p-core/crypto"
)

var (
	tunDev *tun.TUN

	err error
	// RevLookup map[string]string

	// activeStreams map[string]network.Stream
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
		"tun_ice",
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

	// Start node with default settings
	// node, err := libp2p.New(
	// 	libp2p.ListenAddrStrings("/ip4/0.0.0.0/tcp/7788"),
	// 	libp2p.Ping(false),
	// 	nodefunc.StreamHandler,
	// )
	ctx := context.Background()

	livenode, err := nodefunc.CraftNode(
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

	// ch := make(chan os.Signal, 1)
	// signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	// <-ch

	// fmt.Println("Received signal! Shutting down...")

	// if err := livenode.Close(); err != nil {
	// 	panic(err)
	// }

	// Listen for packets on the created TUN interface

	go nodefunc.ExitSignal(livenode)

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
		fmt.Println("Decoding the packet")
		dest := net.IPv4(packet[16], packet[17], packet[18], packet[19]).String()

		fmt.Println(dest)

	}

}
