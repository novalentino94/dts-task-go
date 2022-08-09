package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

func main() {
	// Open database connection
	// db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/dts-go")
	db, err := sql.Open("postgres", "postgresql://lgrkhsrnqvheta:60f66899372dc0342a5f04cb1c44a6274ad7b8ce4a7808901cdc4c8c899dc455@ec2-54-225-234-165.compute-1.amazonaws.com/d9dqhb00i76kr5")

	// if there is an error opening the connection, handle it
	if err != nil {
		panic(err.Error())
	}

	defer db.Close()

	type Task struct {
		ID       int       `json:"id"`
		Detail   string    `json:"detail"`
		Assignee string    `json:"assignee"`
		DueDate  time.Time `json:"dueDate"`
		IsDone   bool      `json:"isDone"`
	}

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"PUT", "PATCH"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return origin == "*"
		},
		MaxAge: 12 * time.Hour,
	}))

	r.GET("/tasks", func(c *gin.Context) {
		// Execute the query
		results, err := db.Query("SELECT * FROM tasks")
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}

		tasks := []Task{}
		for results.Next() {
			var task Task
			// for each row, scan the result into our tag composite object
			err = results.Scan(&task.ID, &task.Detail, &task.Assignee, &task.DueDate, &task.IsDone)
			if err != nil {
				panic(err.Error()) // proper error handling instead of panic in your app
			}

			tasks = append(tasks, task)
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Success",
			"data":    tasks,
		})
	})

	r.GET("/tasks/:id", func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))

		var task Task
		err := db.QueryRow("SELECT * FROM tasks where id = ?", id).Scan(&task.ID, &task.Detail, &task.Assignee, &task.DueDate, &task.IsDone)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Success",
			"data":    task,
		})
	})

	// CREATE
	r.POST("/task", func(c *gin.Context) {
		detail := c.PostForm("detail")
		assignee := c.PostForm("assignee")
		dueDate := c.PostForm("dueDate")
		fmt.Println("detail")
		fmt.Println(detail)
		theSql := fmt.Sprintf("INSERT INTO tasks ( detail, assignee, dueDate ) VALUES ('%s', '%s', '%s')", detail, assignee, dueDate)
		task, err := db.Query(theSql)
		// if there is an error inserting, handle it
		if err != nil {
			panic(err.Error())
		}
		// be careful deferring Queries if you are using transactions
		defer task.Close()

		c.JSON(http.StatusCreated, gin.H{
			"message": "Success",
		})
	})

	// UPDATE
	r.PUT("/task/:id", func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))

		detail := c.PostForm("detail")
		assignee := c.PostForm("assignee")
		dueDate := c.PostForm("dueDate")

		theSql := fmt.Sprintf("UPDATE tasks SET detail = '%s', assignee = '%s', dueDate = '%s' WHERE id = %d", detail, assignee, dueDate, id)
		task, err := db.Query(theSql)
		// if there is an error inserting, handle it
		if err != nil {
			panic(err.Error())
		}
		// be careful deferring Queries if you are using transactions
		defer task.Close()

		c.JSON(http.StatusCreated, gin.H{
			"message": "Success",
		})
	})

	// UPDATE STATUS
	r.PATCH("/task/:id", func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))

		theSql := fmt.Sprintf("UPDATE tasks SET isDone = true WHERE id = %d", id)
		task, err := db.Query(theSql)
		// if there is an error inserting, handle it
		if err != nil {
			panic(err.Error())
		}
		// be careful deferring Queries if you are using transactions
		defer task.Close()

		c.JSON(http.StatusCreated, gin.H{
			"message": "Success",
		})
	})

	r.SetFuncMap(template.FuncMap{
		"upper": strings.ToUpper,
		"inc": func(i int) int {
			return i + 1
		},
		"formatDate": func(date time.Time) string {
			return date.Format("2006-01-02")
		},
	})

	r.Static("/assets", "./assets")
	r.LoadHTMLGlob("templates/*.html")

	r.GET("/", func(c *gin.Context) {
		results, err := db.Query("SELECT * FROM tasks")
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}

		tasks := []Task{}
		for results.Next() {
			var task Task
			// for each row, scan the result into our tag composite object
			err = results.Scan(&task.ID, &task.Detail, &task.Assignee, &task.DueDate, &task.IsDone)
			if err != nil {
				panic(err.Error()) // proper error handling instead of panic in your app
			}

			tasks = append(tasks, task)
		}

		c.HTML(http.StatusOK, "index.html", gin.H{
			"tasks": tasks,
		})
	})

	r.GET("/create", func(c *gin.Context) {
		c.HTML(http.StatusOK, "create.html", gin.H{})
	})

	r.GET("/edit/:id", func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))

		var task Task
		sqlStatement := fmt.Sprintf("SELECT * FROM tasks where id = %d", id)
		err := db.QueryRow(sqlStatement).Scan(&task.ID, &task.Detail, &task.Assignee, &task.DueDate, &task.IsDone)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}

		c.HTML(http.StatusOK, "edit.html", gin.H{
			"id":       task.ID,
			"detail":   task.Detail,
			"assignee": task.Assignee,
			"dueDate":  task.DueDate.Format("2006-01-02"),
		})
	})

	// r.Run(":3000")
	r.Run()

}
