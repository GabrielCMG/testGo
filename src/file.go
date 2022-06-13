package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"
)

//this type represnts a record with three fields
type payload struct {
	One   float32
	Two   float64
	Three uint32
}

//func main() {
//	writeFile()
//	readFile()
//}

func readFile() {

	file, err := os.Open("test.bin")
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}

	var m payload

	for i := 0; i < 30; i++ {
		data := readNextBytes(file, binary.Size(m))
		buffer := bytes.NewBuffer(data)
		err = binary.Read(buffer, binary.BigEndian, &m)
		if err != nil {
			log.Fatal("binary.Read failed", err)
		}

		fmt.Println(m)
	}

}

func readNextBytes(file *os.File, number int) []byte {
	b := make([]byte, number)

	_, err := file.Read(b)
	if err != nil {
		log.Fatal(err)
	}

	return b
}

func writeFile() {
	file, err := os.OpenFile("test.bin", os.O_RDWR, 0644)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < 3; i++ {

		s := payload{
			r.Float32(),
			r.Float64(),
			r.Uint32(),
		}

		fmt.Printf("%v\n", s)

		var binBuf bytes.Buffer
		err := binary.Write(&binBuf, binary.BigEndian, s)
		if err != nil {
			fmt.Print("oups")
			return
		}
		writeNextBytes(file, binBuf.Bytes())
	}
}

func writeNextBytes(file *os.File, bytes []byte) {
	_, err := file.WriteAt(bytes, 0)

	if err != nil {
		log.Fatal(err)
	}

}
