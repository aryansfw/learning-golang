package main

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
)

type Blog struct {
	Id      int64   `json:"id"`
	Title   *string `json:"title"`
	Content *string `json:"content"`
}

type BlogRepository interface {
	GetAll() []Blog
	GetById(id int64) Blog
	Create(b Blog) Blog
	Update(blog Blog) (*Blog, error)
	Delete(id int64)
}

type BlogRepositoryDB struct {
	DB *sql.DB
}

func NewBlogRepository(db *sql.DB) BlogRepository {
	return &BlogRepositoryDB{DB: db} // only one
}

func (r *BlogRepositoryDB) GetAll() []Blog {
	var blogs []Blog

	rows, err := r.DB.Query("SELECT * FROM blogs")
	if err != nil {
		log.Println(err.Error())
		return []Blog{}
	}
	defer rows.Close()

	for rows.Next() {
		var blog = Blog{}
		if err := rows.Scan(&blog.Id, &blog.Title, &blog.Content); err != nil {
			log.Println(err.Error())
			return []Blog{}
		}

		blogs = append(blogs, blog)
	}

	if err = rows.Err(); err != nil {
		log.Println(err.Error())
		return []Blog{}
	}

	return blogs
}

func (r *BlogRepositoryDB) GetById(id int64) Blog {
	var blog Blog
	if err := r.DB.QueryRow("SELECT * FROM blogs WHERE id = $1", id).Scan(&blog.Id, &blog.Title, &blog.Content); err != nil {
		log.Println(err.Error())
		return blog
	}
	return blog
}

func (r *BlogRepositoryDB) Create(b Blog) Blog {
	var id int64
	if err := r.DB.QueryRow(
		"INSERT INTO blogs (title, content) VALUES ($1, $2) RETURNING id",
		b.Title, b.Content,
	).Scan(&id); err != nil {
		log.Println(err.Error())
		return b
	}

	b.Id = id
	return b
}

func (r *BlogRepositoryDB) Update(b Blog) (*Blog, error) {
	var columns []string
	var data []any

	if b.Title != nil {
		columns = append(columns, "title")
		data = append(data, *b.Title)
	}

	if b.Content != nil {
		columns = append(columns, "content")
		data = append(data, *b.Content)
	}

	data = append(data, b.Id)

	for i := range columns {
		columns[i] = fmt.Sprintf("%s = $%d", columns[i], i+1)
	}

	var query = fmt.Sprintf("UPDATE blogs SET %s WHERE id = $%d RETURNING id,title,content", strings.Join(columns, ","), len(columns)+1)

	var blog Blog
	if err := r.DB.QueryRow(
		query,
		data...,
	).Scan(&blog.Id, &blog.Title, &blog.Content); err != nil {
		return nil, err
	}

	return &blog, nil
}

func (r *BlogRepositoryDB) Delete(id int64) {
	res, err := r.DB.Exec("DELETE FROM blogs WHERE id = $1", id)
	if err != nil {
		log.Println(err.Error())
		return
	}

	rowsAffected, err := res.RowsAffected()

	if err != nil {
		log.Println(err.Error())
		return
	}
	log.Printf("%d row(s) deleted", rowsAffected)
}
