// Formal Dinner Sorter
// Jacob Wu
// 01.28.2020

package main

import ( // Get packages we need
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
  "sort"
  "strconv"
	"time"
)

type person struct { // Every person has:
	name string // Name
	id int // ID
  timesServed int // How many times the student was waiter/kitchen crew
	haveMet []int // ID numbers of students the person has met
  currentAssignment string // Current status
}

type table struct { // Every table has properties:
  occupants []person // Who sits there
  disallow []int // This table's people have already met
}

type ByTimesServed []person // Sort in ascending order by number of times served as a kitchen/waiter staff

func (a ByTimesServed) Len() int           { return len(a) }
func (a ByTimesServed) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByTimesServed) Less(i, j int) bool { return a[i].timesServed < a[j].timesServed }

var ( // Global variables
	people []person // Everyone from the excel sheet
  numKitchenStaff = 7 // MUTABLE no. ppl needed at kitchen
  numTables = 31 // MUTABLE no. ppl needed for waiting
  printDetailsForDebug = false // Toggle for detailed info
  tables []table // Keeps track of all tables
  staffList []person // Keeps track of all staff
  kitchenNames []string // String array of kitchen staff
  waiterNames []string // String array of waiter staff
  data [][]string // data for CSV export
)

func resetVariables() { // Clear existing data for next round of assignments
  tables = tables[:0]
  staffList = make([]person, 0)
  kitchenNames = make([]string, 0)
  waiterNames = make([]string, 0)
  data = make([][]string, 0)
}

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
		people = append(people, person{line[1] + " " + line[0], i, 0, make([]int, 0), ""}) // Add the student
	}
}

func shuffleStudents(people []person) { // Randomize order of students
  rand.Seed(time.Now().UnixNano())
	for i := len(people) - 1; i > 0; i-- { // Fisher–Yates shuffle
		j := rand.Intn(i + 1)
		people[i], people[j] = people[j], people[i]
	}
}

func drawStaff(n int, role string) {
	sort.Sort(ByTimesServed(people)) // Sort by how many times served as Kitchen/Waiter staff
	for i := 1;  i<=n; i++ { // Until we get number of staff we need
    staffList = append(staffList, people[0]) // Take the first person
    if role == "Kitchen" {
      kitchenNames = append(kitchenNames, people[0].name)
    } else {
      waiterNames = append(waiterNames, people[0].name)
    } // And write down their name
    fmt.Println(role + " <- " + people[0].name)
    people[0].currentAssignment = role
    people[0].timesServed ++
    people = append(people, people[0]) // Moves first person to last, since he's been assigned role
    people = append(people[:0], people[1:]...)
  }
  if printDetailsForDebug { fmt.Print("\n", people, "\n\n") }
}

func sortTables() {
  tables = make([]table, 50) // Make empty table
  peoplePerTable := (len(people)-numKitchenStaff-numTables)/numTables // Calculate number of people per table by (Non-staff students)/(Number of tables)
  if printDetailsForDebug { fmt.Println(peoplePerTable) }
  tablesWithExtraPerson := (len(people)-numKitchenStaff-numTables)%numTables // Calculate tables with one extra person; remainder from (Non-staff students)/(Number of tables)
  if printDetailsForDebug { fmt.Println(tablesWithExtraPerson) }
  for i := 1; i<=numTables; i++ {
    if i <= tablesWithExtraPerson { // Accomodate the extra person
      tryStudent(i)
    }
    for j := 0; j<peoplePerTable; j++ { // Place number of people at table
      tryStudent(i)
    }
    tableIds := make([]int, 0)
    var tableNames []string // Keep a list of people at the table, in string format
    for _,occupant := range tables[i].occupants {
      tableNames = append(tableNames, occupant.name) // Add name to table's occupant property
      tableIds = append(tableIds, occupant.id) // Add name to table's id
    }
    for k := 1; k<=peoplePerTable; k++ {
      people[len(people)-k].haveMet = tableIds // Indicate everyone has met everyone at the table (id)
    }
    if i <= tablesWithExtraPerson { // If this is the table with one extra person
      people[len(people)-peoplePerTable-1].haveMet = tableIds // Indicate the extra person has met everyone at the table (id)
    }
    data = append(data,tableNames) // Add names to CSV data queue
    fmt.Print("Table ", i, " <- ")
    for _,occupant := range tables[i].occupants {
      fmt.Print(occupant.name)
      fmt.Print("(", occupant.id, ")  ")
    }
    fmt.Print("\n\n")
  }
}

func tryStudent(tableNo int) {
  trialNo := 0 // Start trying at the first person
  for contains(tables[tableNo].disallow, people[trialNo].id) {
    // The upcoming student has met one of the existing table members
    trialNo ++ // In that case, pick another student
  } // If this person have already met someone at the table, try someone else
  tables[tableNo].disallow = append(tables[tableNo].disallow, people[trialNo].haveMet...) // Incorporate this person's haveMet list into the table's haveMet list
  tables[tableNo].occupants = append(tables[tableNo].occupants, people[trialNo]) // Add student to table
  people[trialNo].currentAssignment = "Table " + strconv.Itoa(tableNo)
  people = append(people, people[trialNo]) // Moves that person to last, since he's been assigned role
    people = append(people[:trialNo], people[trialNo+1:]...)
}

func contains(array []int, key int) bool { // Check if there are overlap
	for _, a := range array {
		if a == key {
			return true
		}
	}
	return false
}

func writeStaffNames() { // Export kitchen crew at row 33 and waiters row 34
  data = append(data, make([]string, 0))
  data = append(data, kitchenNames)
  data = append(data, waiterNames)
}

func exportCSV(n string) {
  fileName := "result" + n + ".csv"
  file, err := os.Create(fileName)
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

  exportCSV("1")
  fmt.Print("Finished run without error\n\n")

  resetVariables()
////////////////////////
  shuffleStudents(people)
  if printDetailsForDebug { fmt.Print(people, "\n\n") }

  drawStaff(numKitchenStaff, "Kitchen")
  drawStaff(numTables, "Waiter")

  sortTables()
  writeStaffNames()

  exportCSV("2")
  fmt.Print("Finished run without error")

  resetVariables()
////////////////////////
  shuffleStudents(people)
  if printDetailsForDebug { fmt.Print(people, "\n\n") }

  drawStaff(numKitchenStaff, "Kitchen")
  drawStaff(numTables, "Waiter")

  sortTables()
  writeStaffNames()

  exportCSV("3")
  fmt.Print("Finished run without error")

  resetVariables()
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
// {9} https://stackoverflow.com/questions/10105935/how-to-convert-an-int-value-to-string-in-go