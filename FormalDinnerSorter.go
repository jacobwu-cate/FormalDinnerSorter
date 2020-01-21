package main
 
//import (
//	"strings"
//	"fmt"
//	"os"
//)

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"math/rand"
	"time"
)

type Person struct {
	Firstname string
	Lastname string
}

func main() {
	csvFile, _ := os.Open("Dinner Seating - Student List 2018-19.csv")
	reader := csv.NewReader(bufio.NewReader(csvFile))
	var people []Person
	for {
		line, error := reader.Read()
		if error == io.EOF {
			break // exit if finished (for loop not while)
		} else if error != nil {
			log.Fatal(error)
		}
//		fmt.Println (line[1], line[0])
		people = append (people, Person{line[1], line[0]})
	}
//	fmt.Print(people)
	rand.Seed(time.Now().UnixNano())
	for i := len(people) - 1; i > 0; i-- { // Fisher–Yates shuffle
		j := rand.Intn(i + 1)
		people[i], people[j] = people[j], people[i]
	}
	fmt.Print(people)
	fmt.Print("\n\nKitchen Staff (7):")
	fmt.Print(people[:7])
	fmt.Print("\n\nWaiters (31):")
	fmt.Print(people[7:38])
}
 




//package main
//
//import (
//	"bufio"
//	"encoding/csv"
//	"encoding/json"
//	"fmt"
//	"io"
//	"log"
//	"os"
//)
//
//type Person struct {
//	Firstname string `json:"firstname`
//	Lastname string `json:"lastname`
//}
//
//func main() {
//	csvFile, _ := os.Open("Dinner Seating - Student List 2018-19.csv")
//	reader := csv.NewReader(bufio.NewReader(csvFile))
//	var people []Person
//	for {
//		line, error := reader.Read()
//		if error == io.EOF {
//			break
//		} else if error != nil {
//			log.Fatal(error)
//		}
//		people = append(people, Person{
//			Firstname: line[1],
//			Lastname: line[0],
//		})
//	}
//	peopleJson, _ := json.Marshal(people)
//	fmt.Println(string(peopleJson))
//}
//

// References
// {1} https://www.thepolyglotdeveloper.com/2017/03/parse-csv-data-go-programming-language/
// {2} https://yourbasic.org/golang/shuffle-slice-array/