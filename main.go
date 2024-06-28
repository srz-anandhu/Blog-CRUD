package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

// DB credentials
const (
	user     = "postgres"
	password = "password"
	host     = "localhost"
	port     = 5432
	dbname   = "BlogDB"
)

// DB Initialization
var Db *sql.DB

func InitDB() {
	connectionString := fmt.Sprintf("user=%s password=%s host=%s port=%d dbname=%s sslmode=disable", user, password, host, port, dbname)
	var err error
	Db, err = sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatalf("check the connection string: %v", err)
	}

	if err := Db.Ping(); err != nil {
		log.Fatalf("unable to connect to the database : %v", err)
	}

	fmt.Println("successfully connected to db...")
}

// Create author
func CreateAuthor(username, password string) (int, error) {
	var authorId int
	query := `INSERT INTO authors(author_name, author_password) VALUES ($1,$2) RETURNING authorid`

	if err := Db.QueryRow(query, username, password).Scan(&authorId); err != nil {
		return 0, err
	}
	return authorId, nil
}

// Create Blog
func CreateBlog(title string, id int, content string) (int, error) {
	var blogId int
	query := `INSERT INTO blog(title, authorid, content) VALUES ($1,$2,$3) RETURNING blogid`

	if err := Db.QueryRow(query, title, id, content).Scan(&blogId); err != nil {
		return 0, err
	}
	return blogId, nil
}

// Read Blog
func ReadBlog(id int) (string, int, string, error) {
	var title string
	var authorid int
	var content string

	query := `SELECT title, authorid, content FROM blog WHERE blogid=$1`

	if err := Db.QueryRow(query, id).Scan(&title, &authorid, &content); err != nil {
		return "", 0, "", err
	}
	return title, authorid, content, nil
}

// Read All Blogs
func ReadAllBlogs() ([]map[string]interface{}, error) {
	query := `SELECT title, blogid, authorid, content, created_at, modified_at FROM blog WHERE is_deleted=FALSE`

	rows, err := Db.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var blogs []map[string]interface{}

	for rows.Next() {
		var title string
		var blogId int
		var authorId int
		var content string
		var created_at time.Time
		var modified_at time.Time

		if err := rows.Scan(&title, &blogId, &authorId, &content, &created_at, &modified_at); err != nil {
			return nil, err
		}

		blog := map[string]interface{}{
			"title":       title,
			"blogid":      blogId,
			"authorid":    authorId,
			"content":     content,
			"created_at":  created_at,
			"modified_at": modified_at,
		}

		blogs = append(blogs, blog)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return blogs, nil
}

// Update Blog
func UpdateBlog(id int, title, content string) error {
	query := `UPDATE blog SET title=$1, content=$2, modified_at=NOW() WHERE blogid=$3`

	_, err := Db.Exec(query, title, content, id)
	return err

}

// Delete Blog : updating is_deleted field in blogs to TRUE
func DeleteBlog(id int) error {
	query := `UPDATE blog SET is_deleted=TRUE WHERE blogid=$1`

	_, err := Db.Exec(query, id)
	return err
}

func main() {
	InitDB()
	defer Db.Close()

	// Creating an author
	autherId, err := CreateAuthor("user1", "user123")
	if err != nil {
		log.Fatalf("error while creating author: %v", err)
	}
	fmt.Printf("New author created with ID: %d\n", autherId)

	// Creating a blog with auther id
	blogId, err := CreateBlog("Blog title", 1, "its a demo blog content")
	if err != nil {
		log.Fatalf("error while creating a new blog: %v", err)
	}
	fmt.Printf("New blog created with ID: %d\n", blogId)

	title, id, content, err := ReadBlog(1)
	if err != nil {
		log.Fatalf("error while getting a blog with id : %v", err)
	}
	fmt.Printf("Title:%s, Author ID:%d, Blog Content: %s", title, id, content)

	// Getting All Blogs
	blogs, err := ReadAllBlogs()
	if err != nil {
		log.Fatalf("error while getting all blogs: %v", err)
	}
	fmt.Println("Blogs: ", blogs)

	if err := UpdateBlog(7, "updated Title", "this is updated content"); err != nil {
		log.Fatalf("error while updating a blog : %v", err)
	}

	// Delete a Blog (Soft Delete)
	if err := DeleteBlog(5); err != nil {
		log.Fatalf("error while deleting a blog: %v", err)
	}
}
