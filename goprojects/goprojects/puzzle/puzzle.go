package main

import (
	"fmt"
	"time"
)

const (
	second int64 = 1e9
)

type homes struct {
	red int
	green int
	blue int
	operator int
	check int
}
	
func lt() {
	fmt.Print(time.Nanoseconds(),":")
}
func checkRow(home homes, rotation int) bool {

	return true
}

var opstring string = "+"
func main() {
	lt();fmt.Println("Begin Program")
	success := false
	total := 0
	checkstring := ""
	for i := 0; i < 100; i ++ {// Red wheel home position
	for x := 0; x < 10 ; x ++ {// Operator wheel home position
	for j := 0; j < 100; j ++ {// Green wheel home position
	for y := 0; y < 10 ; y ++ {// Check wheel home position
	for k := 0; k < 100; k ++ {// Blue wheel home position
	for z := 0; z < 10 ; z ++ {// Barrel orientation	
		/*switch {// Confirm operator type and compute total
			case x == 0 || x == 2 || x == 6: // + operator
				total = i + j
				opstring = "+"
			case x == 1 || x == 3 || x == 7: // - operator
				total = i - j
				opstring = "-"
			case x == 4 || x == 8:		 // x operator
				total = i * j
				opstring = "x"
			case x == 5 || x == 9:		 // / operator
				total = i / j
				opstring = "/"
		}*/
		switch {// make the test
			case y == 0 || y == 1 || y == 2 || y == 5 || y == 6 || y == 7: // = check
				if total == k {success = true}
				checkstring = "="
			case y == 3: // > check
				if total > k {success = true}
				checkstring = ">"
			case y == 4: // < check
				if total < k {success = true}
				checkstring = "<"
			case y == 8: // <= check
				if total <= k {success = true}
				checkstring = "<="
			case y == 9: // >= check
				if total >= k {success = true}
				checkstring = ">="
		}
		fmt.Println(i,opstring,j,checkstring,k,success,total)
		success = false
		time.Sleep(1e8)
	}}}}}}
	time.Sleep(second)
	lt();fmt.Println("Exiting")
}
