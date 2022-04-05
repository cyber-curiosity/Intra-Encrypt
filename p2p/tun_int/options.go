package tun

// Define modifiers for the TUN device
type Option func(tun *TUN) error

func Address(address string) Option {
	return func(tun *TUN) error {
		return tun.setAddress(address)
	}
}

func MTU(mtu int) Option {
	return func(tun *TUN) error {
		return tun.setMTU(mtu)
	}
}

func DestAddress(address string) Option {
	return func(tun *TUN) error {
		return tun.setDestAddress(address)
	}
}
