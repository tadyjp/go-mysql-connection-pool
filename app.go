package main

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"os/exec"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

const (
	timeFormat = "2006-01-02 15:04:05"
)

var (
	dbDriver   = "mysql"
	dbuser     = "username"
	dbpass     = "password"
	dbhost     = "db"
	dbport     = "3306"
	dbname     = "app"
	dataSource string

	dbNotFoundError = errors.New("dbNotFoundError")

	sharedDB *sql.DB
)

type Item struct {
	ID        int    `json:"id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	Name      string `json:"name"`
	Hash      string
}

type hashFunc func(str string) string

type RequestData struct {
	Name string `json:"name" binding:"required"`
}

func handleWithHashProcessWithoutPool(c *gin.Context) {
	db, err := sql.Open(dbDriver, dataSource)
	if err != nil {
		panic("Unable to connect the DB.")
	}
	defer db.Close()

	process(c, db, makeHashWithProcess)
}

func handleWithHashLibWithoutPool(c *gin.Context) {
	db, err := sql.Open(dbDriver, dataSource)
	if err != nil {
		panic("Unable to connect the DB.")
	}
	defer db.Close()

	process(c, db, makeHashWithLib)
}

func handleWithHashLibWithPool(c *gin.Context) {
	process(c, sharedDB, makeHashWithLib)
}

func process(c *gin.Context, db *sql.DB, hashFn hashFunc) {
	var data RequestData
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item, err := selectItemByName(db, data.Name)
	if err != nil {
		switch err {
		case dbNotFoundError:
			hash := hashFn(data.Name)
			err := insertItem(db, data.Name, hash)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			}
			item, err = selectItemByName(db, data.Name)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			}
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
	}

	if item.Hash == hashFn(data.Name) {
		c.JSON(http.StatusOK, gin.H{"status": "OK", "item": item})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "hash is invalid"})
	}
}

func makeHashWithProcess(str string) string {
	result := str

	for i := 0; i < 1; i++ {
		out, err := exec.Command("sh", "-c", "echo -n '"+result+"' | openssl sha256").Output()
		if err != nil {
			panic("Failed to generate the hash.")
		}

		outputRaw := string(out)
		output := strings.Split(outputRaw, " ")
		result = strings.TrimRight(output[1], "\n")
	}
	return result
}

func makeHashWithLib(str string) string {
	result := str

	h := sha256.New()
	for i := 0; i < 1; i++ {
		h.Write([]byte(result))
		result = hex.EncodeToString(h.Sum(nil))
		h.Reset()
	}

	return result
}

func selectItemByName(db *sql.DB, name string) (*Item, error) {
	item := new(Item)

	query := "SELECT id, created_at, updated_at, name, hash FROM items WHERE name=(?)"
	rows, err := db.Query(query, name)
	if err != nil {
		return item, fmt.Errorf("Query Error: %w", err)
	}
	defer rows.Close()

	if rows.Next() == false {
		return item, dbNotFoundError
	}

	err = rows.Scan(&item.ID, &item.CreatedAt, &item.UpdatedAt, &item.Name, &item.Hash)
	if err != nil {
		return item, fmt.Errorf("DB Scan Error: %w", err)
	}

	return item, nil
}

func insertItem(db *sql.DB, name string, hash string) error {
	t := time.Now()
	query := "INSERT INTO items (created_at, updated_at, name, hash) VALUES ((?), (?), (?), (?));"
	_, err := db.Exec(query, t.Format(timeFormat), t.Format(timeFormat), name, hash)
	if err != nil {
		return fmt.Errorf("Insert Error: %w", err)
	}

	return nil
}

func main() {
	dataSource = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&interpolateParams=true", dbuser, dbpass, dbhost, dbport, dbname)

	db, err := sql.Open(dbDriver, dataSource)
	if err != nil {
		panic("Unable to connect the DB.")
	}
	sharedDB = db
	defer sharedDB.Close()

	router := gin.Default()
	router.POST("/withHashProcessWithoutPool", handleWithHashProcessWithoutPool)
	router.POST("/withHashLibWithoutPool", handleWithHashLibWithoutPool)
	router.POST("/withHashLibWithPool", handleWithHashLibWithPool)
	router.Run(":8080")
}
