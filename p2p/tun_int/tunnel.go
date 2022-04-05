package tun

import (
	"errors"

	"github.com/songgao/water"
	"github.com/vishvananda/netlink"
)

type TUN struct {
	Inter *water.Interface
	MTU   int
	Src   string
	Dst   string
}

// Configure the specified options for the TUN interface
func (t *TUN) Apply(opts ...Option) error {
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		if err := opt(t); err != nil {
			return err
		}
	}
	return nil
}

func New(name string, opts ...Option) (*TUN, error) {
	// Configuration for the TUN interface
	config := water.Config{
		DeviceType: water.TUN,
	}

	config.Name = name

	// Create Water interface
	Inter, err := water.New(config)
	if err != nil {
		return nil, err
	}

	// TUN result struct
	result := TUN{
		Inter: Inter,
	}

	// Apply the provided config options
	err = result.Apply(opts...)
	return &result, err
}

func (t *TUN) setMTU(mtu int) error {
	link, err := netlink.LinkByName(t.Inter.Name())
	if err != nil {
		return err
	}
	return netlink.LinkSetMTU(link, mtu)
}

func (t *TUN) setAddress(address string) error {
	addr, err := netlink.ParseAddr(address)
	if err != nil {
		return err
	}
	link, err := netlink.LinkByName(t.Inter.Name())
	if err != nil {
		return err
	}
	return netlink.AddrAdd(link, addr)
}

func (t *TUN) setDestAddress(address string) error {
	return errors.New("Error: Destination addresses are not supported by Linux")
}

func (t *TUN) Up() error {
	link, err := netlink.LinkByName(t.Inter.Name())
	if err != nil {
		return err
	}
	return netlink.LinkSetUp(link)
}

func (t *TUN) Down() error {
	link, err := netlink.LinkByName(t.Inter.Name())
	if err != nil {
		return err
	}
	return netlink.LinkSetDown(link)
}

func (t *TUN) Delete() error {
	link, err := netlink.LinkByName(t.Inter.Name())
	if err != nil {
		return err
	}
	return netlink.LinkDel(link)
}
