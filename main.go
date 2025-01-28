package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type book struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Author   string `json:"author"`
	Quantity int    `json:"quantity"`
}

func getBooksFromDB(books *[]book) error {
	rows, err := DB.Query("SELECT id, title, author, quantity FROM books")
	if err != nil {
		return err
	}
	defer rows.Close()
	*books = []book{}			// bad for performance (will work on that later)
	for rows.Next() {
		var b book
		err := rows.Scan(&b.ID, &b.Title, &b.Author, &b.Quantity)
		if err != nil {
			return err
		}
		*books = append(*books, b)
	}
	return nil
}

func getBooks(c *gin.Context) {
	var books []book
    err := getBooksFromDB(&books)
    if err != nil {
        c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch books"})
        return
    }
	c.HTML(http.StatusOK, "books.html", gin.H{"books": books})
}

func addBook(c *gin.Context) {
	var newBook book
    if err := c.BindJSON(&newBook); err != nil {
        c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid request"})
        return
    }
	_, err := DB.Exec("INSERT INTO books (id, title, author, quantity) VALUES ($1, $2, $3, $4)", newBook.ID, newBook.Title, newBook.Author, newBook.Quantity)
    if err != nil {
        c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Failed to add book"})
        return
    }
	c.IndentedJSON(http.StatusCreated, gin.H{"message": "Book added successfully"})
}

func deleteBook(c *gin.Context) {
	id := c.Param("id")
	_, err := DB.Exec("DELETE FROM books WHERE id = $1", id)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Failed to delete book"})
        return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"message": "Book deleted successfully"})
}

func bookByID(c *gin.Context) {
	id := c.Param("id") 
	var b book
	err := DB.QueryRow("SELECT id, title, author, quantity FROM books WHERE id = $1", id).Scan(&b.ID, &b.Title, &b.Author, &b.Quantity)
	if err == sql.ErrNoRows {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not found"})
		return
	} else if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Error retrieving book"})
		return
	}
	c.HTML(http.StatusOK, "book.html", b)
}

func checkoutBook(c *gin.Context) {
	id := c.Param("id")
	var b book
	err := DB.QueryRow("SELECT id, title, author, quantity FROM books WHERE id = $1", id).Scan(&b.ID, &b.Title, &b.Author, &b.Quantity)
	if err == sql.ErrNoRows {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not found"})
		return
	} else if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Error retrieving book"})
		return
	}

	if b.Quantity <= 0 {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Book not available"})
		return
	}
	_, err = DB.Exec("UPDATE books SET quantity = quantity - 1 WHERE id = $1", id)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Failed to checkout book"})
		return
	}
	b.Quantity-- 
	c.HTML(http.StatusOK, "checkout.html", b)
}	

func returnBook(c *gin.Context) {
	id := c.Param("id")
	var b book
	err := DB.QueryRow("SELECT id, title, author, quantity FROM books WHERE id = $1", id).Scan(&b.ID, &b.Title, &b.Author, &b.Quantity)
	if err == sql.ErrNoRows {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not found"})
		return
	} else if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Error retrieving book"})
		return
	}

	_, err = DB.Exec("UPDATE books SET quantity = quantity + 1 WHERE id = $1", id)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Failed to return book"})
		return
	}
	b.Quantity++ 
	c.HTML(http.StatusOK, "return.html", b)
}

func main() {
	if err := openDB(); err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer closeDB()
	
	router := gin.Default()
	router.Use(CORSMiddleware())

	// Serve static files
	router.Static("/static", "./static")

	// Load HTML templates
	router.LoadHTMLGlob("templates/*")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	router.GET("/books", getBooks)
	router.POST("/books", addBook)
	router.DELETE("/books/:id", deleteBook)
	router.GET("/books/:id", bookByID)
	router.PUT("/books/:id/checkout", checkoutBook)
	router.PUT("/books/:id/return", returnBook)


	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	if err := router.Run(":" + port); err != nil {
		log.Panicf("error: %s", err)
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
