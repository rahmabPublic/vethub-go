package api

import (
	"database/sql"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

type OwnerResponse struct {
	ID        int                  `json:"id"`
	FirstName string               `json:"firstName"`
	LastName  string               `json:"lastName"`
	Address   string               `json:"address"`
	City      string               `json:"city"`
	Telephone string               `json:"telephone"`
	Email     *string              `json:"email"`
	Pets      []PetSummaryResponse `json:"pets"`
}

type PetSummaryResponse struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	BirthDate string `json:"birthDate"`
	TypeName  string `json:"typeName"`
}

type CreateOwnerRequest struct {
	FirstName string  `json:"firstName"`
	LastName  string  `json:"lastName"`
	Address   *string `json:"address"`
	City      *string `json:"city"`
	Telephone *string `json:"telephone"`
	Email     *string `json:"email"`
}

type UpdateOwnerRequest = CreateOwnerRequest

var (
	digitRegex = regexp.MustCompile(`^\d+$`)
	emailRegex = regexp.MustCompile(`^[A-Za-z0-9+_.-]+@[A-Za-z0-9.-]+$`)
)

const ownerPictureMaxBytes int64 = 5 * 1024 * 1024

func validateOwnerPicture(contentType string, sizeBytes int64) *string {
	normalizedType := strings.TrimSpace(strings.Split(strings.ToLower(contentType), ";")[0])
	switch normalizedType {
	case "image/jpeg", "image/png", "image/webp":
	default:
		msg := "picture must be a JPG, PNG, or WebP image"
		return &msg
	}

	if sizeBytes <= 0 {
		msg := "picture file must not be empty"
		return &msg
	}
	if sizeBytes > ownerPictureMaxBytes {
		msg := "picture file size must be at most 5 MB"
		return &msg
	}

	return nil
}

func validateOwner(req *CreateOwnerRequest) *string {
	if strings.TrimSpace(req.FirstName) == "" {
		msg := "firstName is required"
		return &msg
	}
	if len(req.FirstName) > 255 {
		msg := "firstName must be at most 255 characters"
		return &msg
	}
	if strings.TrimSpace(req.LastName) == "" {
		msg := "lastName is required"
		return &msg
	}
	if len(req.LastName) > 255 {
		msg := "lastName must be at most 255 characters"
		return &msg
	}
	if req.Address != nil && len(*req.Address) > 255 {
		msg := "address must be at most 255 characters"
		return &msg
	}
	if req.City != nil && len(*req.City) > 255 {
		msg := "city must be at most 255 characters"
		return &msg
	}
	if req.Telephone != nil && *req.Telephone != "" {
		if len(*req.Telephone) > 255 {
			msg := "telephone must be at most 255 characters"
			return &msg
		}
		if !digitRegex.MatchString(*req.Telephone) {
			msg := "telephone must contain only digits"
			return &msg
		}
	}
	if req.Email != nil && *req.Email != "" {
		if len(*req.Email) > 255 {
			msg := "email must be at most 255 characters"
			return &msg
		}
		if !emailRegex.MatchString(*req.Email) {
			msg := "email must be a valid email address"
			return &msg
		}
	}
	return nil
}

func RegisterOwnerRoutes(mux *http.ServeMux, db *sql.DB) {
	mux.HandleFunc("GET /api/v1/owners", listOwners(db))
	mux.HandleFunc("GET /api/v1/owners/{id}", getOwner(db))
	mux.HandleFunc("POST /api/v1/owners", createOwner(db))
	mux.HandleFunc("PUT /api/v1/owners/{id}", updateOwner(db))
	mux.HandleFunc("DELETE /api/v1/owners/{id}", deleteOwner(db))
}

func listOwners(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		lastName := r.URL.Query().Get("lastName")
		var rows *sql.Rows
		var err error
		if lastName != "" {
			rows, err = db.Query("SELECT id, first_name, last_name, address, city, telephone, email FROM owners WHERE LOWER(last_name) LIKE LOWER(? || '%')", lastName)
		} else {
			rows, err = db.Query("SELECT id, first_name, last_name, address, city, telephone, email FROM owners")
		}
		if err != nil {
			Error(w, http.StatusInternalServerError, "ERR-0001", err.Error())
			return
		}
		defer rows.Close()

		owners := []OwnerResponse{}
		for rows.Next() {
			var o OwnerResponse
			var address, city, telephone, email sql.NullString
			if err := rows.Scan(&o.ID, &o.FirstName, &o.LastName, &address, &city, &telephone, &email); err != nil {
				Error(w, http.StatusInternalServerError, "ERR-0001", err.Error())
				return
			}
			o.Address = address.String
			o.City = city.String
			o.Telephone = telephone.String
			if email.Valid {
				o.Email = &email.String
			}
			o.Pets = loadPetSummaries(db, o.ID)
			owners = append(owners, o)
		}
		JSON(w, http.StatusOK, owners)
	}
}

func getOwner(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(PathID(r, "id"))
		if err != nil {
			NotFound(w, "ERR-0002", "The requested owner does not exist.")
			return
		}
		owner, err := findOwnerByID(db, id)
		if err != nil {
			NotFound(w, "ERR-0002", "The requested owner does not exist.")
			return
		}
		JSON(w, http.StatusOK, owner)
	}
}

func createOwner(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CreateOwnerRequest
		if err := DecodeJSON(r, &req); err != nil {
			BadRequest(w, "ERR-0001", "Invalid request body")
			return
		}
		if msg := validateOwner(&req); msg != nil {
			BadRequest(w, "ERR-0001", *msg)
			return
		}
		result, err := db.Exec(
			"INSERT INTO owners (first_name, last_name, address, city, telephone, email) VALUES (?, ?, ?, ?, ?, ?)",
			req.FirstName, req.LastName, nilStr(req.Address), nilStr(req.City), nilStr(req.Telephone), nilStr(req.Email),
		)
		if err != nil {
			Error(w, http.StatusInternalServerError, "ERR-0001", err.Error())
			return
		}
		id, _ := result.LastInsertId()
		owner, _ := findOwnerByID(db, int(id))
		JSON(w, http.StatusCreated, owner)
	}
}

func updateOwner(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(PathID(r, "id"))
		if err != nil {
			NotFound(w, "ERR-0002", "The requested owner does not exist.")
			return
		}
		if _, err := findOwnerByID(db, id); err != nil {
			NotFound(w, "ERR-0002", "The requested owner does not exist.")
			return
		}
		var req UpdateOwnerRequest
		if err := DecodeJSON(r, &req); err != nil {
			BadRequest(w, "ERR-0001", "Invalid request body")
			return
		}
		if msg := validateOwner(&req); msg != nil {
			BadRequest(w, "ERR-0001", *msg)
			return
		}
		db.Exec(
			"UPDATE owners SET first_name=?, last_name=?, address=?, city=?, telephone=?, email=? WHERE id=?",
			req.FirstName, req.LastName, nilStr(req.Address), nilStr(req.City), nilStr(req.Telephone), nilStr(req.Email), id,
		)
		owner, _ := findOwnerByID(db, id)
		JSON(w, http.StatusOK, owner)
	}
}

func deleteOwner(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(PathID(r, "id"))
		if err != nil {
			NotFound(w, "ERR-0002", "The requested owner does not exist.")
			return
		}
		if _, err := findOwnerByID(db, id); err != nil {
			NotFound(w, "ERR-0002", "The requested owner does not exist.")
			return
		}
		db.Exec("DELETE FROM visits WHERE pet_id IN (SELECT id FROM pets WHERE owner_id = ?)", id)
		db.Exec("DELETE FROM pets WHERE owner_id = ?", id)
		db.Exec("DELETE FROM owners WHERE id = ?", id)
		NoContent(w)
	}
}

func findOwnerByID(db *sql.DB, id int) (*OwnerResponse, error) {
	var o OwnerResponse
	var address, city, telephone, email sql.NullString
	err := db.QueryRow(
		"SELECT id, first_name, last_name, address, city, telephone, email FROM owners WHERE id = ?", id,
	).Scan(&o.ID, &o.FirstName, &o.LastName, &address, &city, &telephone, &email)
	if err != nil {
		return nil, err
	}
	o.Address = address.String
	o.City = city.String
	o.Telephone = telephone.String
	if email.Valid {
		o.Email = &email.String
	}
	o.Pets = loadPetSummaries(db, o.ID)
	return &o, nil
}

func loadPetSummaries(db *sql.DB, ownerID int) []PetSummaryResponse {
	rows, err := db.Query(
		"SELECT p.id, p.name, p.birth_date, pt.name FROM pets p JOIN pet_types pt ON p.type_id = pt.id WHERE p.owner_id = ?",
		ownerID,
	)
	if err != nil {
		return []PetSummaryResponse{}
	}
	defer rows.Close()
	pets := []PetSummaryResponse{}
	for rows.Next() {
		var p PetSummaryResponse
		rows.Scan(&p.ID, &p.Name, &p.BirthDate, &p.TypeName)
		pets = append(pets, p)
	}
	return pets
}

func nilStr(s *string) any {
	if s == nil {
		return nil
	}
	return *s
}
