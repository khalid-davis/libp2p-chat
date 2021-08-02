package main

import (
	"encoding/pem"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/multiformats/go-multiaddr"
	"io/ioutil"
	"log"

	"fmt"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-kad-dht/dual"
	"golang.org/x/net/context"
)

func main() {
	// openssl genrsa -out rsa_private.key 2048
	certBytes, err := ioutil.ReadFile("./relay/server.key")
	if err != nil {
		log.Println("unable to read client.pem, error: ", err)
		return
	}
	block, _ := pem.Decode(certBytes)

	priv, err := crypto.UnmarshalRsaPrivateKey(block.Bytes)
	if err != nil {
		panic(err)
	}


	ctx := context.Background()

	// libp2p.New constructs a new libp2p Host.
	// Other options can be added here.
	sourceMultiAddr, _ := multiaddr.NewMultiaddr("/ip4/0.0.0.0/tcp/5000")

	var externalMultiAddr multiaddr.Multiaddr
	externalMultiAddr, err = multiaddr.NewMultiaddr(fmt.Sprintf("/ip4/%s/tcp/%d", "119.8.58.38", 5000))
	if err != nil {
		fmt.Println("New multiaddr err: ", err)
	}
	addressFactory := func(addrs []multiaddr.Multiaddr) []multiaddr.Multiaddr {
		if externalMultiAddr != nil {
			addrs = append(addrs, externalMultiAddr)
		}
		return addrs
	}


	if err != nil {
		panic(err)
	}
	host, err := libp2p.New(
		ctx,
		libp2p.ListenAddrs(sourceMultiAddr),
		libp2p.Identity(priv),
		libp2p.EnableNATService(),
		libp2p.AddrsFactory(addressFactory),
	)
	if err != nil {
		panic(err)
	}

	fmt.Println("This node: ", host.ID().Pretty(), " ", host.Addrs())

	_, err = dual.New(ctx, host, dht.Mode(dht.ModeServer))
	if err != nil {
		panic(err)
	}


	select {}
}
