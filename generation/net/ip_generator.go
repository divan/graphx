package net

import "net"

// IPGenerator generates sequental and non-repeating IP
// addresses.
type IPGenerator struct {
	last net.IP
}

// NewIPGenerator creates new ip generator and sets last IP to the
// given value.
func NewIPGenerator(ip string) *IPGenerator {
	last := net.ParseIP(ip)
	return &IPGenerator{
		last: last,
	}
}

// NextAddress returns next generated IP address.
func (i *IPGenerator) NextAddress() string {
	p := i.last
	a, b, c, d := p[12], p[13], p[14], p[15]

	if d < 255 {
		d++
	} else {
		d = 0
		if c < 255 {
			c++
		} else {
			c = 0
			if b < 255 {
				b++
			} else {
				b = 0
				if a > 255 {
					panic("Seriously?")
					// a = 0 :)
				}

				a++
			}
		}
	}

	ip := net.IPv4(a, b, c, d)
	i.last = ip
	return ip.String()
}
