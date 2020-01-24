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
	var people []string
	var tables [31][10]string
	for {
		line, error := reader.Read()
		if error == io.EOF {
			break // exit if finished (for loop not while)
		} else if error != nil {
			log.Fatal(error)
		}
//		fmt.Println (line[1], line[0])
		people = append (people, line[1]+" "+line[0])
	}
//	fmt.Print(people)
	rand.Seed(time.Now().UnixNano())
	for i := len(people) - 1; i > 0; i-- { // Fisher–Yates shuffle
		j := rand.Intn(i + 1)
		people[i], people[j] = people[j], people[i]
	}
	fmt.Print("Kitchen Staff (7):")
	fmt.Print(people[:7])
	for i := 1;  i<=7; i++ { // For i in range(7):
		people = append(people[:0], people[1:]...) // Remove first element of the list
	}
	fmt.Print("\n\nWaiters (31):")
	fmt.Print(people[:31])
	for i := 1;  i<=31; i++ { // For i in range(31):
		people = append(people[:0], people[1:]...) // Remove first element of the list
	}
//	fmt.Print("\n\n", people)
	peoplePerTable := len(people)/31
	tablesWithExtraPerson := len(people)%31
//	fmt.Print("\n\nNumber of people per table:", peoplePerTable)
//	fmt.Print("\n\nNumber of tables with one more person:", tablesWithExtraPerson, "\n\n")
	for i := 0;  i<31; i++ { // For every table
		for j := 0; j < peoplePerTable; j++ {
			tables[i][j] = people[0]
			people = append(people[:0], people[1:]...) // Remove first element of the list
		}
		if i < tablesWithExtraPerson {
			tables[i][peoplePerTable] = people[0]
			people = append(people[:0], people[1:]...) // Remove first element of the list
		}
	}
	fmt.Print("\n\n", tables)
	
	file, err := os.Create("result.csv")
	if err != nil {
		log.Fatal("Error returned from filepath.Walk:", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, value := range tables {
		err := writer.Write(value[:])
		if err != nil {
			log.Fatal("Error returned from filepath.Walk:", err)
		}
	}
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
// {3} https://stackoverflow.com/questions/37334119/how-to-delete-an-element-from-a-slice-in-golang
// {4} https://www.cyberciti.biz/faq/golang-for-loop-examples/
// {5} https://www.digitalocean.com/community/tutorials/how-to-do-math-in-go-with-operators
// {6} https://golangcode.com/write-data-to-a-csv-file/