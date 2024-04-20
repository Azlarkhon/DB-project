package main

import (
    "database/sql"
    "fmt"
    "log"
    "net/http"
    "os"

    "github.com/gin-gonic/gin"
    _ "github.com/lib/pq"
    "github.com/joho/godotenv"
)

type Train struct {
    ID    uint   `json:"train_id"`
    Name  string `json:"train_name"`
    Price uint   `json:"train_price"`
}

type Plane struct {
    ID    uint   `json:"plane_id"`
    Name  string `json:"plane_name"`
    Price uint   `json:"plane_price"`
}

var db *sql.DB

func init() {
    if err := godotenv.Load(); err != nil {
        log.Fatal("Error loading .env file")
    }
}

func main() {
    dbUsername := os.Getenv("DB_USERNAME")
    dbPassword := os.Getenv("DB_PASSWORD")
    dbHost := os.Getenv("DB_HOST")
    dbPort := os.Getenv("DB_PORT")
    dbName := os.Getenv("DB_NAME")

    dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUsername, dbPassword, dbHost, dbPort, dbName)

    var err error
    db, err = sql.Open("postgres", dsn)
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }
    defer db.Close()

    router := gin.Default()
    router.GET("/", homePage)
    router.GET("/trains", getAllTrains)
    router.GET("/planes", getAllPlanes)
    router.POST("/trains/add", insertTrain)
    router.POST("/planes/add", insertPlane)

    port := os.Getenv("PORT")
    router.Run(":" + port)
}

func homePage(c *gin.Context) {
    c.String(http.StatusOK, "Welcome to my application!")
}

func getAllTrains(c *gin.Context) {
    rows, err := db.Query("SELECT * FROM trains")
    if err != nil {
        handleDBError(c, err)
        return
    }
    defer rows.Close()

    var trains []Train
    for rows.Next() {
        var train Train
        err := rows.Scan(&train.ID, &train.Name, &train.Price)
        if err != nil {
            handleDBError(c, err)
            return
        }
        trains = append(trains, train)
    }
    c.JSON(http.StatusOK, trains)
}

func getAllPlanes(c *gin.Context) {
    rows, err := db.Query("SELECT * FROM planes")
    if err != nil {
        handleDBError(c, err)
        return
    }
    defer rows.Close()

    var planes []Plane
    for rows.Next() {
        var plane Plane
        err := rows.Scan(&plane.ID, &plane.Name, &plane.Price)
        if err != nil {
            handleDBError(c, err)
            return
        }
        planes = append(planes, plane)
    }
    c.JSON(http.StatusOK, planes)
}

func insertTrain(c *gin.Context) {
    var newTrain Train
    if err := c.BindJSON(&newTrain); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    _, err := db.Exec("INSERT INTO trains (train_name, train_price) VALUES (?, ?)", newTrain.Name, newTrain.Price)
    if err != nil {
        handleDBError(c, err)
        return
    }

    c.JSON(http.StatusCreated, gin.H{"message": "Train created successfully"})
}

func insertPlane(c *gin.Context) {
    var newPlane Plane
    if err := c.BindJSON(&newPlane); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    _, err := db.Exec("INSERT INTO planes (plane_name, plane_price) VALUES (?, ?)", newPlane.Name, newPlane.Price)
    if err != nil {
        handleDBError(c, err)
        return
    }

    c.JSON(http.StatusCreated, gin.H{"message": "Plane created successfully"})
}

func handleDBError(c *gin.Context, err error) {
    log.Printf("Database error: %v", err)
    c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
}