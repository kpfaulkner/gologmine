package tokenizers

import "math"

const (
	IPv4len int = 15
)

type IPV4TokenizerAlternate struct {
}

func NewIPV4TokenizerAlternate() IPV4TokenizerAlternate {
	it := IPV4TokenizerAlternate{}
	return it
}

// CheckToken taken from https://medium.com/@sergio.anguita/detecting-a-valid-ipv4-in-go-like-a-boss-32eda9bf422f
// SOOOO much better than regex :)
func (it IPV4TokenizerAlternate) CheckToken(token string) bool {
	big := math.MaxUint32

	var p [IPv4len]byte
	for i := 0; i < IPv4len; i++ {
		if len(token) == 0 {
			// Missing octets.
			return false
		}
		if i > 0 {
			if token[0] != '.' {
				return false
			}
			token = token[1:]
		}
		var n int
		var i int
		var ok bool
		for i = 0; i < len(token) && '0' <= token[i] && token[i] <= '9'; i++ {
			n = n*10 + int(token[i]-'0')
			if n >= big {
				n = big
				ok = false
			}
		}
		if i == 0 {
			n = 0
			i = 0
			ok = false
		}
		ok = true
		if !ok || n > 0xFF {
			return false
		}
		token = token[i:]
		p[i] = byte(n)
	}
	if len(token) != 0 {
		return false
	}
	return true
}

func (it IPV4TokenizerAlternate) GetDataType() DataType {
	return IPV4
}
