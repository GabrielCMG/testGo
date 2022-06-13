package main

import (
	"fmt"
	"math/rand"
	"time"
)

func triFusionGo(liste []int, c chan []int) {
	if len(liste) <= 1 {
		c <- liste
	} else {
		c1 := make(chan []int)
		go triFusionGo(liste[:len(liste)/2], c1)
		go triFusionGo(liste[len(liste)/2:], c1)
		c2 := make(chan []int)
		go fusionGo(<-c1, <-c1, c2)
		c <- <-c2
		//c <- fusion(<-c1, <-c2)
	}
}

func triFusion(liste []int) []int {
	if len(liste) == 1 {
		return liste
	} else {
		return fusion(triFusion(liste[:len(liste)/2]), triFusion(liste[len(liste)/2:]))
	}
}

func fusionGo(liste1, liste2 []int, c chan []int) {
	if len(liste1) == 0 {
		c <- liste2
	} else if len(liste2) == 0 {
		c <- liste1
	} else if liste1[0] <= liste2[0] {
		c1 := make(chan []int)
		go fusionGo(liste1[1:], liste2, c1)
		c <- append([]int{liste1[0]}, <-c1...)
	} else {
		c1 := make(chan []int)
		go fusionGo(liste1, liste2[1:], c1)
		c <- append([]int{liste2[0]}, <-c1...)
	}

}

func fusion(liste1, liste2 []int) []int {
	if len(liste1) == 0 {
		return liste2
	} else if len(liste2) == 0 {
		return liste1
	} else if liste1[0] <= liste2[0] {
		return append([]int{liste1[0]}, fusion(liste1[1:], liste2)...)
	} else {
		return append([]int{liste2[0]}, fusion(liste1, liste2[1:])...)
	}

}

func main() {
	r1, r2, r3, r4, r5 := rand.Perm(50000), rand.Perm(50000), rand.Perm(50000), rand.Perm(50000), rand.Perm(50000)

	c := make(chan []int, 5)
	t1 := time.Now()
	go triFusionGo(r1, c)
	go triFusionGo(r2, c)
	go triFusionGo(r3, c)
	go triFusionGo(r4, c)
	go triFusionGo(r5, c)
	l11 := <-c
	l12 := <-c
	l13 := <-c
	l14 := <-c
	l15 := <-c
	fmt.Print(time.Since(t1))

	t2 := time.Now()
	l21 := triFusion(r1)
	l22 := triFusion(r2)
	l23 := triFusion(r3)
	l24 := triFusion(r4)
	l25 := triFusion(r5)
	fmt.Print(time.Since(t2))

	fmt.Printf("\n%v %v %v %v %v", l11[:10], l12[:10], l13[:10], l14[:10], l15[:10])
	fmt.Printf("\n%v %v %v %v %v", l21[:10], l22[:10], l23[:10], l24[:10], l25[:10])
	fmt.Printf("\n%v", r1[:10])
}
