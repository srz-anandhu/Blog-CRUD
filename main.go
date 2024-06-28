package main

import (
	"database/sql"
	"fmt"
	"log"

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
	Db, err := sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatalf("check the connection string: %v", err)
	}

	defer Db.Close()

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
	query := `INSERT INTO blogs(title, authorid, content) VALUES ($1,$2,$3) RETURNING blogid`

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

	query := `SELECT title, authorid, content FROM blogs WHERE blogid=$1`

	if err := Db.QueryRow(query, id).Scan(&title, &authorid, &content); err != nil {
		return "", 0, "", err
	}
	return title, authorid, content, nil
}

// Read All Blogs
func ReadAllBlogs() ([]map[string]interface{}, error) {
	query := `SELECT * FROM blogs WHERE is_deleted=FALSE`

	rows, err := Db.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var blogs []map[string]interface{}

	for rows.Next() {
		var title string
		var authorId int
		var content string

		if err := rows.Scan(&title, &authorId, &content); err != nil {
			return nil, err
		}

		blog := map[string]interface{}{
			"title":    title,
			"authorid": authorId,
			"content":  content,
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
	query := `UPDATE blogs SET title=$1, content=$2 updated_at=NOW() WHERE blogid=$4`

	_, err := Db.Exec(query, title, content, id)
	return err

}

// Delete Blog : updating is_deleted field in blogs to TRUE
func DeleteBlog(id int) error {
	query := `UPDATE blogs SET is_deleted=TRUE WHERE blogid=$1`

	_, err := Db.Exec(query, id)
	return err
}

func main() {
	InitDB()
	defer Db.Close()
}
