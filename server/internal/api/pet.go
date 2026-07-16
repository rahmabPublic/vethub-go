package api

import (
	"database/sql"
	"net/http"
	"strconv"
)

type PetResponse struct {
	ID             int                    `json:"id"`
	Name           string                 `json:"name"`
	BirthDate      string                 `json:"birthDate"`
	Type           PetTypeResponse        `json:"type"`
	OwnerID        int                    `json:"ownerId"`
	OwnerFirstName string                 `json:"ownerFirstName"`
	OwnerLastName  string                 `json:"ownerLastName"`
	Visits         []VisitSummaryResponse `json:"visits"`
}

type VisitSummaryResponse struct {
	ID          int     `json:"id"`
	Date        string  `json:"date"`
	Description *string `json:"description"`
}

type PetTypeResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type CreatePetRequest struct {
	Name      string `json:"name"`
	BirthDate string `json:"birthDate"`
	TypeID    int    `json:"typeId"`
	OwnerID   *int   `json:"ownerId"`
}

type UpdatePetRequest struct {
	Name      string `json:"name"`
	BirthDate string `json:"birthDate"`
	TypeID    int    `json:"typeId"`
}

type CreatePetTypeRequest struct {
	Name string `json:"name"`
}

type UpdatePetTypeRequest = CreatePetTypeRequest

func RegisterPetRoutes(mux *http.ServeMux, db *sql.DB) {
	// Pet Types
	mux.HandleFunc("GET /api/v1/pet-types", listPetTypes(db))
	mux.HandleFunc("GET /api/v1/pet-types/{id}", getPetType(db))
	mux.HandleFunc("POST /api/v1/pet-types", createPetType(db))
	mux.HandleFunc("PUT /api/v1/pet-types/{id}", updatePetType(db))
	mux.HandleFunc("DELETE /api/v1/pet-types/{id}", deletePetType(db))

	// Pets nested under owner
	mux.HandleFunc("GET /api/v1/owners/{ownerId}/pets", listPetsByOwner(db))
	mux.HandleFunc("GET /api/v1/owners/{ownerId}/pets/{petId}", getPetByOwner(db))
	mux.HandleFunc("POST /api/v1/owners/{ownerId}/pets", createPetForOwner(db))
	mux.HandleFunc("PUT /api/v1/owners/{ownerId}/pets/{petId}", updatePetByOwner(db))
	mux.HandleFunc("DELETE /api/v1/owners/{ownerId}/pets/{petId}", deletePetByOwner(db))

	// Pets global
	mux.HandleFunc("GET /api/v1/pets", listAllPets(db))
	mux.HandleFunc("GET /api/v1/pets/{id}", getPetGlobal(db))
	mux.HandleFunc("POST /api/v1/pets", createPetGlobal(db))
	mux.HandleFunc("PUT /api/v1/pets/{id}", updatePetGlobal(db))
	mux.HandleFunc("DELETE /api/v1/pets/{id}", deletePetGlobal(db))
}

// --- Pet Type handlers ---

func listPetTypes(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT id, name FROM pet_types")
		if err != nil {
			Error(w, http.StatusInternalServerError, "ERR-0001", err.Error())
			return
		}
		defer rows.Close()
		types := []PetTypeResponse{}
		for rows.Next() {
			var t PetTypeResponse
			rows.Scan(&t.ID, &t.Name)
			types = append(types, t)
		}
		JSON(w, http.StatusOK, types)
	}
}

func getPetType(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(PathID(r, "id"))
		if err != nil {
			NotFound(w, "ERR-0004", "The requested pet type does not exist.")
			return
		}
		var t PetTypeResponse
		if err := db.QueryRow("SELECT id, name FROM pet_types WHERE id = ?", id).Scan(&t.ID, &t.Name); err != nil {
			NotFound(w, "ERR-0004", "The requested pet type does not exist.")
			return
		}
		JSON(w, http.StatusOK, t)
	}
}

func createPetType(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CreatePetTypeRequest
		if err := DecodeJSON(r, &req); err != nil {
			BadRequest(w, "ERR-0001", "Invalid request body")
			return
		}
		result, err := db.Exec("INSERT INTO pet_types (name) VALUES (?)", req.Name)
		if err != nil {
			Error(w, http.StatusInternalServerError, "ERR-0001", err.Error())
			return
		}
		id, _ := result.LastInsertId()
		JSON(w, http.StatusCreated, PetTypeResponse{ID: int(id), Name: req.Name})
	}
}

func updatePetType(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(PathID(r, "id"))
		if err != nil {
			NotFound(w, "ERR-0004", "The requested pet type does not exist.")
			return
		}
		var existing PetTypeResponse
		if err := db.QueryRow("SELECT id, name FROM pet_types WHERE id = ?", id).Scan(&existing.ID, &existing.Name); err != nil {
			NotFound(w, "ERR-0004", "The requested pet type does not exist.")
			return
		}
		var req UpdatePetTypeRequest
		if err := DecodeJSON(r, &req); err != nil {
			BadRequest(w, "ERR-0001", "Invalid request body")
			return
		}
		db.Exec("UPDATE pet_types SET name = ? WHERE id = ?", req.Name, id)
		JSON(w, http.StatusOK, PetTypeResponse{ID: id, Name: req.Name})
	}
}

func deletePetType(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(PathID(r, "id"))
		if err != nil {
			NotFound(w, "ERR-0004", "The requested pet type does not exist.")
			return
		}
		var exists int
		if err := db.QueryRow("SELECT id FROM pet_types WHERE id = ?", id).Scan(&exists); err != nil {
			NotFound(w, "ERR-0004", "The requested pet type does not exist.")
			return
		}
		if _, err := db.Exec("DELETE FROM pet_types WHERE id = ?", id); err != nil {
			Error(w, http.StatusConflict, "ERR-0001", "Cannot delete pet type: it is still in use.")
			return
		}
		NoContent(w)
	}
}

// --- Pet helpers ---

func findPetByID(db *sql.DB, id int) (*PetResponse, error) {
	var p PetResponse
	err := db.QueryRow(`
		SELECT p.id, p.name, p.birth_date, pt.id, pt.name, o.id, o.first_name, o.last_name
		FROM pets p
		JOIN pet_types pt ON p.type_id = pt.id
		JOIN owners o ON p.owner_id = o.id
		WHERE p.id = ?`, id,
	).Scan(&p.ID, &p.Name, &p.BirthDate, &p.Type.ID, &p.Type.Name, &p.OwnerID, &p.OwnerFirstName, &p.OwnerLastName)
	if err != nil {
		return nil, err
	}
	p.Visits = loadVisitSummaries(db, p.ID)
	return &p, nil
}

func findPetByIDAndOwner(db *sql.DB, petID, ownerID int) (*PetResponse, error) {
	pet, err := findPetByID(db, petID)
	if err != nil {
		return nil, err
	}
	if pet.OwnerID != ownerID {
		return nil, sql.ErrNoRows
	}
	return pet, nil
}

func loadVisitSummaries(db *sql.DB, petID int) []VisitSummaryResponse {
	rows, err := db.Query("SELECT id, date, description FROM visits WHERE pet_id = ?", petID)
	if err != nil {
		return []VisitSummaryResponse{}
	}
	defer rows.Close()
	visits := []VisitSummaryResponse{}
	for rows.Next() {
		var v VisitSummaryResponse
		var desc sql.NullString
		rows.Scan(&v.ID, &v.Date, &desc)
		if desc.Valid {
			v.Description = &desc.String
		}
		visits = append(visits, v)
	}
	return visits
}

func scanPets(db *sql.DB, rows *sql.Rows) []PetResponse {
	pets := []PetResponse{}
	for rows.Next() {
		var p PetResponse
		rows.Scan(&p.ID, &p.Name, &p.BirthDate, &p.Type.ID, &p.Type.Name, &p.OwnerID, &p.OwnerFirstName, &p.OwnerLastName)
		p.Visits = loadVisitSummaries(db, p.ID)
		pets = append(pets, p)
	}
	return pets
}

// --- Nested pet handlers ---

func listPetsByOwner(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ownerID, err := strconv.Atoi(PathID(r, "ownerId"))
		if err != nil {
			NotFound(w, "ERR-0002", "The requested owner does not exist.")
			return
		}
		var exists int
		if err := db.QueryRow("SELECT id FROM owners WHERE id = ?", ownerID).Scan(&exists); err != nil {
			NotFound(w, "ERR-0002", "The requested owner does not exist.")
			return
		}
		rows, err := db.Query(`
			SELECT p.id, p.name, p.birth_date, pt.id, pt.name, o.id, o.first_name, o.last_name
			FROM pets p JOIN pet_types pt ON p.type_id = pt.id JOIN owners o ON p.owner_id = o.id
			WHERE p.owner_id = ?`, ownerID)
		if err != nil {
			Error(w, http.StatusInternalServerError, "ERR-0001", err.Error())
			return
		}
		defer rows.Close()
		JSON(w, http.StatusOK, scanPets(db, rows))
	}
}

func getPetByOwner(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ownerID, _ := strconv.Atoi(PathID(r, "ownerId"))
		petID, _ := strconv.Atoi(PathID(r, "petId"))
		pet, err := findPetByIDAndOwner(db, petID, ownerID)
		if err != nil {
			NotFound(w, "ERR-0003", "The requested pet does not exist.")
			return
		}
		JSON(w, http.StatusOK, pet)
	}
}

func createPetForOwner(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ownerID, _ := strconv.Atoi(PathID(r, "ownerId"))
		var exists int
		if err := db.QueryRow("SELECT id FROM owners WHERE id = ?", ownerID).Scan(&exists); err != nil {
			NotFound(w, "ERR-0002", "The requested owner does not exist.")
			return
		}
		var req CreatePetRequest
		if err := DecodeJSON(r, &req); err != nil {
			BadRequest(w, "ERR-0001", "Invalid request body")
			return
		}
		if err := db.QueryRow("SELECT id FROM pet_types WHERE id = ?", req.TypeID).Scan(&exists); err != nil {
			NotFound(w, "ERR-0004", "The requested pet type does not exist.")
			return
		}
		result, err := db.Exec("INSERT INTO pets (name, birth_date, type_id, owner_id) VALUES (?, ?, ?, ?)",
			req.Name, req.BirthDate, req.TypeID, ownerID)
		if err != nil {
			Error(w, http.StatusInternalServerError, "ERR-0001", err.Error())
			return
		}
		id, _ := result.LastInsertId()
		pet, _ := findPetByID(db, int(id))
		JSON(w, http.StatusCreated, pet)
	}
}

func updatePetByOwner(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ownerID, _ := strconv.Atoi(PathID(r, "ownerId"))
		petID, _ := strconv.Atoi(PathID(r, "petId"))
		if _, err := findPetByIDAndOwner(db, petID, ownerID); err != nil {
			NotFound(w, "ERR-0003", "The requested pet does not exist.")
			return
		}
		var req UpdatePetRequest
		if err := DecodeJSON(r, &req); err != nil {
			BadRequest(w, "ERR-0001", "Invalid request body")
			return
		}
		var exists int
		if err := db.QueryRow("SELECT id FROM pet_types WHERE id = ?", req.TypeID).Scan(&exists); err != nil {
			NotFound(w, "ERR-0004", "The requested pet type does not exist.")
			return
		}
		db.Exec("UPDATE pets SET name=?, birth_date=?, type_id=? WHERE id=?", req.Name, req.BirthDate, req.TypeID, petID)
		pet, _ := findPetByID(db, petID)
		JSON(w, http.StatusOK, pet)
	}
}

func deletePetByOwner(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ownerID, _ := strconv.Atoi(PathID(r, "ownerId"))
		petID, _ := strconv.Atoi(PathID(r, "petId"))
		if _, err := findPetByIDAndOwner(db, petID, ownerID); err != nil {
			NotFound(w, "ERR-0003", "The requested pet does not exist.")
			return
		}
		db.Exec("DELETE FROM visits WHERE pet_id = ?", petID)
		db.Exec("DELETE FROM pets WHERE id = ?", petID)
		NoContent(w)
	}
}

// --- Global pet handlers ---

func listAllPets(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query(`
			SELECT p.id, p.name, p.birth_date, pt.id, pt.name, o.id, o.first_name, o.last_name
			FROM pets p JOIN pet_types pt ON p.type_id = pt.id JOIN owners o ON p.owner_id = o.id`)
		if err != nil {
			Error(w, http.StatusInternalServerError, "ERR-0001", err.Error())
			return
		}
		defer rows.Close()
		JSON(w, http.StatusOK, scanPets(db, rows))
	}
}

func getPetGlobal(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, _ := strconv.Atoi(PathID(r, "id"))
		pet, err := findPetByID(db, id)
		if err != nil {
			NotFound(w, "ERR-0003", "The requested pet does not exist.")
			return
		}
		JSON(w, http.StatusOK, pet)
	}
}

func createPetGlobal(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CreatePetRequest
		if err := DecodeJSON(r, &req); err != nil {
			BadRequest(w, "ERR-0001", "Invalid request body")
			return
		}
		if req.OwnerID == nil {
			BadRequest(w, "ERR-0001", "ownerId is required")
			return
		}
		var exists int
		if err := db.QueryRow("SELECT id FROM owners WHERE id = ?", *req.OwnerID).Scan(&exists); err != nil {
			NotFound(w, "ERR-0002", "The requested owner does not exist.")
			return
		}
		if err := db.QueryRow("SELECT id FROM pet_types WHERE id = ?", req.TypeID).Scan(&exists); err != nil {
			NotFound(w, "ERR-0004", "The requested pet type does not exist.")
			return
		}
		result, err := db.Exec("INSERT INTO pets (name, birth_date, type_id, owner_id) VALUES (?, ?, ?, ?)",
			req.Name, req.BirthDate, req.TypeID, *req.OwnerID)
		if err != nil {
			Error(w, http.StatusInternalServerError, "ERR-0001", err.Error())
			return
		}
		id, _ := result.LastInsertId()
		pet, _ := findPetByID(db, int(id))
		JSON(w, http.StatusCreated, pet)
	}
}

func updatePetGlobal(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, _ := strconv.Atoi(PathID(r, "id"))
		if _, err := findPetByID(db, id); err != nil {
			NotFound(w, "ERR-0003", "The requested pet does not exist.")
			return
		}
		var req UpdatePetRequest
		if err := DecodeJSON(r, &req); err != nil {
			BadRequest(w, "ERR-0001", "Invalid request body")
			return
		}
		var exists int
		if err := db.QueryRow("SELECT id FROM pet_types WHERE id = ?", req.TypeID).Scan(&exists); err != nil {
			NotFound(w, "ERR-0004", "The requested pet type does not exist.")
			return
		}
		db.Exec("UPDATE pets SET name=?, birth_date=?, type_id=? WHERE id=?", req.Name, req.BirthDate, req.TypeID, id)
		pet, _ := findPetByID(db, id)
		JSON(w, http.StatusOK, pet)
	}
}

func deletePetGlobal(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, _ := strconv.Atoi(PathID(r, "id"))
		if _, err := findPetByID(db, id); err != nil {
			NotFound(w, "ERR-0003", "The requested pet does not exist.")
			return
		}
		db.Exec("DELETE FROM visits WHERE pet_id = ?", id)
		db.Exec("DELETE FROM pets WHERE id = ?", id)
		NoContent(w)
	}
}
