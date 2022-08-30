package main

import (
	"fmt"
	"strconv"
)
 
func main() {
	DeretPertama := 2
	DeretKedua := 4
	Value := 5

	Jarak := DeretKedua - DeretPertama
	Deret := DeretPertama - Jarak
	Max := 0
	Koma := ""
	for index := 0; index < Value; index++ {		
		Deret = Deret + Jarak
		Max = Value - 1
		if Max == index {
			Koma = ""
		} else {
			Koma = ","
		}
		fmt.Print(strconv.Itoa(Deret) + Koma)
	}
}