package service

import (
	"blog/model"
	"blog/repo"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type BlogService struct {
	BlogRepo *repo.BlogRepository
	LikeRepo *repo.LikeRepository
}

// Kreiranje novog bloga
func (s *BlogService) Create(blog *model.Blog) error {
	blog.ID = uuid.New().String()
	blog.CreatedAt = time.Now()
	blog.UpdatedAt = time.Now()
	blog.Likes = 0
	return s.BlogRepo.Create(blog)
}

// Dohvatanje svih blogova sa ažuriranim brojem lajkova
func (s *BlogService) GetAll() ([]model.BlogResponse, error) {
	blogs, err := s.BlogRepo.GetAll()
	if err != nil {
		return nil, err
	}

	response := []model.BlogResponse{}
	userCache := make(map[string]string) // keš userID -> username

	for _, b := range blogs {
		// Dohvati username iz Stakeholders servisa
		username, ok := userCache[b.UserID]
		if !ok {
			user, err := GetUserFromStakeholders(b.UserID)
			if err != nil || user == nil {
				username = "Nepoznat"
			} else {
				username = user.Username
			}
			userCache[b.UserID] = username
		}

		// Dodaj broj lajkova
		count, _ := s.LikeRepo.CountByBlogID(b.ID)
		b.Likes = int(count)

		response = append(response, model.BlogResponse{
			ID:        b.ID,
			Title:     b.Title,
			Content:   b.Content,
			UserID:    b.UserID,
			Username:  username,
			Likes:     b.Likes,
			ImageURL:  b.ImageURL,
			CreatedAt: b.CreatedAt,
			UpdatedAt: b.UpdatedAt,
		})
	}

	return response, nil
}

// Dohvatanje bloga po ID-u sa ažuriranim brojem lajkova
func (s *BlogService) Get(id string) (*model.BlogResponse, error) {
	blog, err := s.BlogRepo.Get(id)
	if err != nil {
		return nil, err
	}

	count, _ := s.LikeRepo.CountByBlogID(id)
	blog.Likes = int(count)

	// Dohvati username autora iz Stakeholders servisa
	user, err := GetUserFromStakeholders(blog.UserID)
	username := "Nepoznat"
	if err == nil && user != nil {
		username = user.Username
	}

	resp := &model.BlogResponse{
		ID:        blog.ID,
		Title:     blog.Title,
		Content:   blog.Content,
		UserID:    blog.UserID,
		Username:  username,
		CreatedAt: blog.CreatedAt,
		UpdatedAt: blog.UpdatedAt,
		Likes:     blog.Likes,
		ImageURL:  blog.ImageURL,
	}

	return resp, nil
}

func (s *BlogService) Like(blogID, userID string) (int, error) {
	// Proveri da li blog postoji
	_, err := s.BlogRepo.Get(blogID)
	if err != nil {
		return 0, fmt.Errorf("blog not found")
	}

	// Proveri da li korisnik već lajkovao
	exists, err := s.LikeRepo.Exists(userID, blogID)
	if err != nil {
		return 0, err
	}
	if exists {
		count, _ := s.LikeRepo.CountByBlogID(blogID)
		return int(count), nil
	}

	like := model.NewLike(userID, blogID)
	if err := s.LikeRepo.Create(like); err != nil {
		return 0, err
	}

	count, _ := s.LikeRepo.CountByBlogID(blogID)
	return int(count), nil
}

func (s *BlogService) Unlike(blogID, userID string) (int, error) {
	// Proveri da li blog postoji
	if _, err := s.BlogRepo.Get(blogID); err != nil {
		return 0, fmt.Errorf("blog not found")
	}

	if err := s.LikeRepo.Delete(userID, blogID); err != nil {
		return 0, err
	}

	count, _ := s.LikeRepo.CountByBlogID(blogID)
	return int(count), nil
}

func (s *BlogService) GetForUser(userId string) ([]model.Blog, error) {
	// 1. Pozovi Follow mikroservis da dobiješ listu userId-eva koje prati
	url := fmt.Sprintf("http://follow-service:8000/api/follow/following/%s", userId)
	client := &http.Client{Timeout: 3 * time.Second}

	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get following: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("follow service returned %d: %s", resp.StatusCode, string(body))
	}

	var followingResp struct {
		Following []string `json:"following"`
	}
	body, _ := ioutil.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &followingResp); err != nil {
		return nil, fmt.Errorf("failed to parse follow response: %w", err)
	}

	// 2. Dohvati sve blogove iz repozitorijuma
	blogs, err := s.BlogRepo.GetAll()
	if err != nil {
		return nil, err
	}

	// 3. Filtriraj samo blogove autora koje korisnik prati
	filtered := []model.Blog{}
	for _, blog := range blogs {
		for _, followedId := range followingResp.Following {
			if blog.UserID == followedId {
				// Dodaj broj lajkova
				count, _ := s.LikeRepo.CountByBlogID(blog.ID)
				blog.Likes = int(count)
				filtered = append(filtered, blog)
				break
			}
		}
	}

	return filtered, nil
}

func (s *BlogService) GetMyBlogs(userID string) ([]model.Blog, error) {
	return s.BlogRepo.GetByUserID(userID)
}
