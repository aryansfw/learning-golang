package main

import (
	"fmt"
	"strconv"
	"strings"
)

type BlogService struct {
	Repository BlogRepository
}

func NewBlogService(r BlogRepository) *BlogService {
	return &BlogService{Repository: r}
}

func (s *BlogService) GetBlogs() []Blog {
	return s.Repository.GetAll()
}

func (s *BlogService) GetBlogById(id string) Blog {
	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return Blog{}
	}
	return s.Repository.GetById(idInt)
}

func (s *BlogService) CreateBlog(blog Blog) Blog {
	return s.Repository.Create(blog)
}

func (s *BlogService) UpdateBlog(id string, blog Blog) (*Blog, error) {
	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, err
	}
	blog.Id = idInt
	return s.Repository.Update(blog)
}

func (s *BlogService) ReplaceBlog(id string, blog Blog) (*Blog, error) {
	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, err
	}
	blog.Id = idInt

	var missingFields []string

	if blog.Title == nil {
		missingFields = append(missingFields, "title")
	}
	if blog.Content == nil {
		missingFields = append(missingFields, "content")
	}

	if len(missingFields) > 0 {
		return nil, fmt.Errorf("fields missing: %s", strings.Join(missingFields, ", "))
	}
	return &blog, nil
}

func (s *BlogService) DeleteBlog(id string) {
	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return
	}

	s.Repository.Delete(idInt)
}
