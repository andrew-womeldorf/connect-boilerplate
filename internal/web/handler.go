package web

import (
	_ "embed"
	"fmt"
	"html/template"
	"net/http"

	pb "github.com/andrew-womeldorf/connect-boilerplate/gen/user/v1"
	"github.com/andrew-womeldorf/connect-boilerplate/internal/services/user"
)

//go:embed templates/index.html
var indexTemplate string

type Handler struct {
	service *user.Service
}

func NewHandler(service *user.Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) IndexHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	listResp, err := h.service.ListUsers(ctx, &pb.ListUsersRequest{})
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to list users: %v", err), http.StatusInternalServerError)
		return
	}

	data := struct {
		Users []*pb.User
	}{
		Users: listResp.Users,
	}

	tmpl, err := template.New("index").Parse(indexTemplate)
	if err != nil {
		http.Error(w, fmt.Sprintf("Template error: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, fmt.Sprintf("Template execution error: %v", err), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()

	name := r.FormValue("name")
	email := r.FormValue("email")

	if name == "" || email == "" {
		http.Error(w, "Name and email are required", http.StatusBadRequest)
		return
	}

	createReq := &pb.CreateUserRequest{
		Name:  name,
		Email: email,
	}

	_, err := h.service.CreateUser(ctx, createReq)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create user: %v", err), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
