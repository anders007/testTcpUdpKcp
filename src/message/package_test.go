package message

import (
	"testing"
	"fmt"
)

func TestMessageHead(t *testing.T) {
	head := MessageHead{
		id: 2,
		len: 32,
		compress: 0,
	}

	buf := head.Encode()
	fmt.Println("head len : ", len(buf))
	fmt.Println(buf)

	headNew := MessageHead{}

	fmt.Println("new head before: ", headNew)

	headNew.Decode(buf)

	fmt.Println("source head : ", head)
	fmt.Println("new head after: ", headNew)
}
