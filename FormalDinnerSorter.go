// Formal Dinner Sorter
// Jacob Wu
// 01.28.2020

package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
  "sort"
	"time"
)

type person struct {
	name string
	id int
  timesServed int // How many times the student was waiter/kitchen crew
	haveMet []int
  currentAssignment string
}

type table struct {
  occupants []person
  disallow []int
}

type ByTimesServed []person

func (a ByTimesServed) Len() int           { return len(a) }
func (a ByTimesServed) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByTimesServed) Less(i, j int) bool { return a[i].timesServed < a[j].timesServed }

var ( // Global variables
	people []person // Everyone from the excel sheet
  numKitchenStaff = 7
  numTables = 31
  printDetailsForDebug = false
  tables [50]table
  staffList []person
  kitchenNames []string
  waiterNames []string
  data [][]string
)

func getStudentNames() { // This function populates student data
  csvFile, _ := os.Open("Dinner Seating - Student List 2018-19.csv") // Opens csv
	reader := csv.NewReader(bufio.NewReader(csvFile)) // Reads csv
  i := 0 // Counter used to assign ID numbers
	for { // For every student name
    i ++ // increment counter
		line, error := reader.Read()
		if error == io.EOF { // Deal with error, if there is one
			break // exit if finished (for loop not while)
		} else if error != nil {
			log.Fatal(error)
		}
		people = append (people, person{line[1] + " " + line[0], i, 0, make([]int, 0), ""}) // Add the student
	}
}

func shuffleStudents(people []person) {
  rand.Seed(time.Now().UnixNano())
	for i := len(people) - 1; i > 0; i-- { // Fisher–Yates shuffle
		j := rand.Intn(i + 1)
		people[i], people[j] = people[j], people[i]
	}
}

func drawStaff(n int, role string) {
	sort.Sort(ByTimesServed(people)) // Sort by how many times on Kitchen/Waiter
	for i := 1;  i<=n; i++ {
    staffList = append(staffList, people[0])
    if role == "Kitchen" {
      kitchenNames = append(kitchenNames, people[0].name)
    } else {
      waiterNames = append(waiterNames, people[0].name)
    }
    fmt.Println(role + " <- " + people[0].name)
    people = append(people, people[0]) // Moves first person to last, since he's been assigned role
    people = append(people[:0], people[1:]...)
  }
  if printDetailsForDebug { fmt.Print("\n", people, "\n\n") }
}

func sortTables() {
  peoplePerTable := (len(people)-numKitchenStaff-numTables)/numTables
  if printDetailsForDebug { fmt.Println(peoplePerTable) }
  tablesWithExtraPerson := (len(people)-numKitchenStaff-numTables)%numTables
  if printDetailsForDebug { fmt.Println(tablesWithExtraPerson) }
  for i := 1; i<=numTables; i++ {
    if i <= tablesWithExtraPerson { // Accomodate the extra person
      tryStudent(i)
    }
    for j := 0; j<peoplePerTable; j++ {
      tryStudent(i)
    }
    fmt.Print("Table ", i, tables[i].occupants, "\n\n")
    var tableNames []string
    for _,occupant := range tables[i].occupants {
      tableNames = append(tableNames, occupant.name)
    }
    data = append(data,tableNames)
  }
}

func tryStudent(tableNo int) {
  trialNo := 0
  for contains(tables[tableNo].disallow, people[trialNo].id) {
    // The upcoming student has met one of the existing table members
    trialNo ++ // In that case, pick another student
  }
  // fmt.Print(people[trialNo].name)
  if people[trialNo].haveMet != nil {
    tables[tableNo].disallow = append(tables[tableNo].disallow, people[trialNo].haveMet...)
  }
  tables[tableNo].occupants = append(tables[tableNo].occupants, people[trialNo]) // Add student to table
  people = append(people, people[trialNo]) // Moves that person to last, since he's been assigned role
    people = append(people[:trialNo], people[trialNo+1:]...)
}

func contains(array []int, key int) bool {
	for _, a := range array {
		if a == key {
			return true
		}
	}
	return false
}

func writeStaffNames() {
  data = append(data, make([]string, 0))
  data = append(data, kitchenNames)
  data = append(data, waiterNames)
}

func exportCSV() {
  file, err := os.Create("result.csv")
  if err != nil {
    log.Fatal(err)
  }
  defer file.Close()

  writer := csv.NewWriter(file)
  defer writer.Flush()

  for _, value := range data {
    err := writer.Write(value)
    if err != nil {
      log.Fatal(err)
    }
  }
}

func main() {
	getStudentNames()
	if printDetailsForDebug { fmt.Print(people, "\n\n") }

  shuffleStudents(people)
  if printDetailsForDebug { fmt.Print(people, "\n\n") }

  drawStaff(numKitchenStaff, "Kitchen")
  drawStaff(numTables, "Waiter")

  sortTables()
  writeStaffNames()

  exportCSV()
}

// References
// {1} https://www.thepolyglotdeveloper.com/2017/03/parse-csv-data-go-programming-language/
// {2} https://yourbasic.org/golang/shuffle-slice-array/
// {3} https://stackoverflow.com/questions/37334119/how-to-delete-an-element-from-a-slice-in-golang
// {4} https://www.cyberciti.biz/faq/golang-for-loop-examples/
// {5} https://www.digitalocean.com/community/tutorials/how-to-do-math-in-go-with-operators
// {6} https://golangcode.com/write-data-to-a-csv-file/
// {7} https://golang.org/pkg/sort/
// {8} https://ispycode.com/GO/Collections/Arrays/Check-if-item-is-in-array