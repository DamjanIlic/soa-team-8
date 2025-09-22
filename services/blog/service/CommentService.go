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

type CommentService struct {
	CommentRepo       *repo.CommentRepository
	StakeholderClient *StakeholderClient
}

type UserDTO struct {
	ID           string `json:"id"`
	UserID       string `json:"user_id"`
	Username     string `json:"username"`
	Email        string `json:"email"`
	Role         string `json:"role"`
	Name         string `json:"name"`
	Surname      string `json:"surname"`
	ProfileImage string `json:"profile_image"`
	Biography    string `json:"biography"`
	Motto        string `json:"motto"`
	Blocked      bool   `json:"blocked"`
}

// Kreiranje komentara
func (s *CommentService) CreateComment(comment *model.Comment) error {
	comment.ID = uuid.New().String()
	comment.CreatedAt = time.Now()
	comment.UpdatedAt = time.Now()
	return s.CommentRepo.Create(comment)
}

// Dohvatanje korisnika direktno iz Stakeholders servisa
func GetUserFromStakeholders(userID string) (*UserDTO, error) {
	url := fmt.Sprintf("http://stakeholders-service:8080/api/stakeholders/user/%s", userID)
	fmt.Println("Requesting:", url)

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		fmt.Println("HTTP error:", err)
		return nil, err
	}
	defer resp.Body.Close()

	fmt.Println("HTTP status:", resp.Status)
	if resp.StatusCode != 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println("Response body:", string(body))
		return nil, fmt.Errorf("failed to get user, status: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Read body error:", err)
		return nil, err
	}

	var user UserDTO
	if err := json.Unmarshal(body, &user); err != nil {
		fmt.Println("Unmarshal error:", err, "Body:", string(body))
		return nil, err
	}

	fmt.Printf("Fetched user: %+v\n", user)
	return &user, nil
}

// Dohvatanje komentara za blog sa stvarnim user podacima
func (s *CommentService) GetComments(blogID string) ([]model.CommentResponse, error) {
	comments, err := s.CommentRepo.GetByBlogID(blogID)
	if err != nil {
		return nil, err
	}

	response := []model.CommentResponse{}
	userCache := make(map[string]*UserDTO) // keš za korisnike

	for _, c := range comments {
		var user *UserDTO
		if u, ok := userCache[c.UserID]; ok {
			user = u // koristi keš
		} else {
			user, err = GetUserFromStakeholders(c.UserID)
			if err != nil || user == nil {
				user = &UserDTO{
					Username:     "Nepoznat",
					Name:         "",
					Surname:      "",
					ProfileImage: "https://upload.wikimedia.org/wikipedia/commons/8/89/Portrait_Placeholder.png",
					Motto:        "",
				}
			}
			userCache[c.UserID] = user // sačuvaj u keš
		}

		response = append(response, model.CommentResponse{
			ID:        c.ID,
			BlogID:    c.BlogID,
			UserID:    c.UserID,
			Username:  user.Username,
			Name:      user.Name,
			Surname:   user.Surname,
			AvatarUrl: user.ProfileImage,
			Motto:     user.Motto,
			Text:      c.Text,
			CreatedAt: c.CreatedAt,
		})
	}

	return response, nil
}
