package main

import (
	"fmt"
	"sqlfun/backup"
	"sqlfun/internal"
	"sync"
)

func main() {
	originalDB := "../database.db"
    backupDirectory := "backup/"
	
	var wg sync.WaitGroup
    wg.Add(1)
    go func() {
        defer wg.Done()
        if err := backup.Sqlite(originalDB, backupDirectory); err != nil {
            fmt.Println("Error during backup:", err)
        }
    }()

    // Continue with the rest of your application
    wg.Wait() // Wait for the backup goroutine to finish before exiting, if necessary

	s := internal.Server{}
	s.Init()
}
