// Miscellaneous utilities.
package utils

import (
	"github.com/blckur/blckur/errortypes"
	"github.com/dropbox/godropbox/errors"
	"math/rand"
	"net"
	"net/mail"
	"os"
	"os/exec"
	"path"
	"strings"
)

var (
	chars = []rune(
		"abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
)

func RandStr(n int) (str string) {
	strList := make([]rune, n)
	for i := range strList {
		strList[i] = chars[rand.Intn(len(chars))]
	}
	str = string(strList)
	return
}

func ParseEmail(input string) (email string, err error) {
	address, err := mail.ParseAddress(input)
	if err != nil {
		err = &InvalidEmailError{
			errors.Wrap(err, "utils: Invalid email"),
		}
		return
	}
	email = address.Address
	return
}

func ShuffleStrings(data []string) {
	n := len(data)
	for i := n - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		data[i], data[j] = data[j], data[i]
	}
}

func ShuffleStringsNew(data []string) (randData []string) {
	randData = make([]string, len(data))
	copy(randData, data)
	ShuffleStrings(randData)
	return
}

func GetLocalAddress() (addr string, err error) {
	name, err := os.Hostname()
	if err != nil {
		err = &errortypes.UnknownError{
			errors.Wrap(err, "beanstalkd: Get ip"),
		}
		return
	}

	addrs, err := net.LookupHost(name)
	if err != nil {
		err = &errortypes.UnknownError{
			errors.Wrap(err, "beanstalkd: Get ip"),
		}
		return
	}

	addr = addrs[0]

	return
}

func InfToBool(obj interface{}) (val bool) {
	if obj == nil {
		val = false
	} else {
		val = obj.(bool)
	}
	return
}

func ReverseStrings(data []string, start int, end int) {
	for start < end {
		tmp := data[start]
		data[start] = data[end]
		data[end] = tmp
		start += 1
		end -= 1
	}
}

func RotateStrings(data []string, d int) {
	n := len(data) - 1
	ReverseStrings(data, 0, n)
	ReverseStrings(data, 0, d-1)
	ReverseStrings(data, d, n)
}

func MinInt(x int, y int) int {
	if x < y {
		return x
	}
	return y
}

func MaxInt(x int, y int) int {
	if x > y {
		return x
	}
	return y
}

func isEscapable(b byte) bool {
	return !('A' <= b && b <= 'Z' || 'a' <= b && b <= 'z' ||
		'0' <= b && b <= '9' || b == '-' || b == '.' || b == '_' || b == '~')
}

func Escape(s string) string {
	t := make([]byte, 0, 3*len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if isEscapable(c) {
			t = append(t, '%')
			t = append(t, "0123456789ABCDEF"[c>>4])
			t = append(t, "0123456789ABCDEF"[c&15])
		} else {
			t = append(t, s[i])
		}
	}
	return string(t)
}

func GetVersion() (ver string) {
	pkgPath := path.Join(os.Getenv("GOPATH"), "src/github.com/blckur/blckur")

	output, err := exec.Command("git", "-C", pkgPath,
		"rev-parse", "HEAD").Output()
	if err != nil {
		ver = "unknown"
		return
	}

	ver = strings.TrimSpace(string(output))

	if len(ver) > 8 {
		ver = ver[:8]
	}

	return
}
