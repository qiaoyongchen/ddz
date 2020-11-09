package poker

import (
	"fmt"
	"testing"
)

func TestSub(t *testing.T) {
	rst, rste := SubPokers([]IPoker{NewPoker(1, 7), NewPoker(0, 7), NewPoker(3, 8)}, []IPoker{NewPoker(1, 7), NewPoker(0, 7)})
	fmt.Println(rst)
	fmt.Println(rste)
}
