package presentation

import (
	"auth_backend/usecase"
	"encoding/json"
	"net/http"
)

type UserHandler struct {
	userUseCase usecase.UserUseCase
}

func NewUserHandler(userUseCase usecase.UserUseCase) *UserHandler {
	return &UserHandler{userUseCase: userUseCase}
}

func (userHandler *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	token, err := userHandler.userUseCase.Authenticate(r.Context(), body.Email, body.Password)

	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	response := map[string]string{"token": token}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)

	if err != nil {
		return
	}

}

func (userHandler *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := userHandler.userUseCase.Register(r.Context(), body.Email, body.Password); err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusCreated)
	err := json.NewEncoder(w).Encode(map[string]string{"message": "User created successfully!"})
	if err != nil {
		return
	}
}
