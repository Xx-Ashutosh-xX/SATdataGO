package main

import (
    "database/sql"
    "encoding/json"
    "fmt"
    "log"
    "os"

    _ "github.com/go-sql-driver/mysql"
)

// SATResult represents a student's SAT data
type SATResult struct {
    Name    string `json:"name"`
    Address string `json:"address"`
    City    string `json:"city"`
    Country string `json:"country"`
    Pincode string `json:"pincode"`
    Score   int    `json:"score"`
    Passed  bool   `json:"passed"`
}

// Connect to MySQL database
func connectDB() *sql.DB {
    dbURL := "satuser:password@tcp(127.0.0.1:3306)/satdb"
    db, err := sql.Open("mysql", dbURL)
    if err != nil {
        log.Fatalf("Unable to connect to database: %v\n", err)
    }

    // Create table if it doesn't exist
    createTable := `CREATE TABLE IF NOT EXISTS SATResults (
        Name VARCHAR(100) PRIMARY KEY,
        Address VARCHAR(100),
        City VARCHAR(100),
        Country VARCHAR(100),
        Pincode VARCHAR(20),
        Score INT,
        Passed BOOLEAN
    );`

    _, err = db.Exec(createTable)
    if err != nil {
        log.Fatalf("Error creating table: %v\n", err)
    }

    return db
}

// Insert data into the database
func insertData(db *sql.DB) {
    var name, address, city, country, pincode string
    var score int

    fmt.Println("Enter Name:")
    fmt.Scanln(&name)
    fmt.Println("Enter Address:")
    fmt.Scanln(&address)
    fmt.Println("Enter City:")
    fmt.Scanln(&city)
    fmt.Println("Enter Country:")
    fmt.Scanln(&country)
    fmt.Println("Enter Pincode:")
    fmt.Scanln(&pincode)
    fmt.Println("Enter SAT Score:")
    fmt.Scanln(&score)

    passed := score > 1600*30/100

    query := `INSERT INTO SATResults (Name, Address, City, Country, Pincode, Score, Passed) VALUES (?, ?, ?, ?, ?, ?, ?)`
    _, err := db.Exec(query, name, address, city, country, pincode, score, passed)
    if err != nil {
        log.Println("Error inserting data:", err)
    } else {
        fmt.Println("Data inserted successfully")
    }
}

// View all data from the database
func viewAllData(db *sql.DB) {
    rows, err := db.Query("SELECT Name, Address, City, Country, Pincode, Score, Passed FROM SATResults")
    if err != nil {
        log.Println("Error fetching data:", err)
        return
    }
    defer rows.Close()

    var results []SATResult
    for rows.Next() {
        var result SATResult
        err := rows.Scan(&result.Name, &result.Address, &result.City, &result.Country, &result.Pincode, &result.Score, &result.Passed)
        if err != nil {
            log.Println("Error scanning row:", err)
            continue
        }
        results = append(results, result)
    }

    jsonData, err := json.MarshalIndent(results, "", "    ")
    if err != nil {
        log.Println("Error converting data to JSON:", err)
    } else {
        fmt.Println(string(jsonData))
    }
}

// Get the rank of a candidate based on SAT score
func getRank(db *sql.DB) {
    var name string
    fmt.Println("Enter Name to get rank:")
    fmt.Scanln(&name)

    var score int
    err := db.QueryRow("SELECT Score FROM SATResults WHERE Name = ?", name).Scan(&score)
    if err != nil {
        log.Println("Error fetching score:", err)
        return
    }

    var rank int
    err = db.QueryRow("SELECT COUNT(*) + 1 FROM SATResults WHERE Score > ?", score).Scan(&rank)
    if err != nil {
        log.Println("Error calculating rank:", err)
        return
    }

    fmt.Printf("%s's rank is: %d\n", name, rank)
}

// Update the SAT score of a candidate
func updateScore(db *sql.DB) {
    var name string
    var newScore int

    fmt.Println("Enter Name to update score:")
    fmt.Scanln(&name)
    fmt.Println("Enter new SAT Score:")
    fmt.Scanln(&newScore)

    passed := newScore > 30

    query := `UPDATE SATResults SET Score = ?, Passed = ? WHERE Name = ?`
    _, err := db.Exec(query, newScore, passed, name)
    if err != nil {
        log.Println("Error updating score:", err)
    } else {
        fmt.Println("Score updated successfully")
    }
}

// Delete a record from the database
func deleteRecord(db *sql.DB) {
    var name string
    fmt.Println("Enter Name to delete record:")
    fmt.Scanln(&name)

    query := `DELETE FROM SATResults WHERE Name = ?`
    _, err := db.Exec(query, name)
    if err != nil {
        log.Println("Error deleting record:", err)
    } else {
        fmt.Println("Record deleted successfully")
    }
}

// Main function - Menu for operations
func main() {
    db := connectDB()
    defer db.Close()

    for {
        fmt.Println("Choose an option:")
        fmt.Println("1. Insert data")
        fmt.Println("2. View all data")
        fmt.Println("3. Get rank")
        fmt.Println("4. Update score")
        fmt.Println("5. Delete one record")
        fmt.Println("6. Exit")

        var choice int
        fmt.Scanln(&choice)

        switch choice {
        case 1:
            insertData(db)
        case 2:
            viewAllData(db)
        case 3:
            getRank(db)
        case 4:
            updateScore(db)
        case 5:
            deleteRecord(db)
        case 6:
            fmt.Println("Exiting...")
            os.Exit(0)
        default:
            fmt.Println("Invalid option, try again.")
        }
    }
}
