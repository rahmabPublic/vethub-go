package api

import (
	"database/sql"
	"net/http"
	"strconv"
)

type VisitResponse struct {
	ID          int     `json:"id"`
	Date        string  `json:"date"`
	Description *string `json:"description"`
	PetID       int     `json:"petId"`
}

type CreateVisitRequest struct {
	Date        string  `json:"date"`
	Description *string `json:"description"`
	PetID       *int    `json:"petId"`
}

type UpdateVisitRequest struct {
	Date        string  `json:"date"`
	Description *string `json:"description"`
}

func RegisterVisitRoutes(mux *http.ServeMux, db *sql.DB) {
	// Nested under owner/pet
	mux.HandleFunc("GET /api/v1/owners/{ownerId}/pets/{petId}/visits", listVisitsByPet(db))
	mux.HandleFunc("GET /api/v1/owners/{ownerId}/pets/{petId}/visits/{visitId}", getVisitByPet(db))
	mux.HandleFunc("POST /api/v1/owners/{ownerId}/pets/{petId}/visits", createVisitForPet(db))
	mux.HandleFunc("PUT /api/v1/owners/{ownerId}/pets/{petId}/visits/{visitId}", updateVisitByPet(db))
	mux.HandleFunc("DELETE /api/v1/owners/{ownerId}/pets/{petId}/visits/{visitId}", deleteVisitByPet(db))

	// Global
	mux.HandleFunc("GET /api/v1/visits", listAllVisits(db))
	mux.HandleFunc("GET /api/v1/visits/{id}", getVisitGlobal(db))
	mux.HandleFunc("POST /api/v1/visits", createVisitGlobal(db))
	mux.HandleFunc("PUT /api/v1/visits/{id}", updateVisitGlobal(db))
	mux.HandleFunc("DELETE /api/v1/visits/{id}", deleteVisitGlobal(db))
}

func findVisitByID(db *sql.DB, id int) (*VisitResponse, error) {
	var v VisitResponse
	var desc sql.NullString
	err := db.QueryRow("SELECT id, date, description, pet_id FROM visits WHERE id = ?", id).
		Scan(&v.ID, &v.Date, &desc, &v.PetID)
	if err != nil {
		return nil, err
	}
	if desc.Valid {
		v.Description = &desc.String
	}
	return &v, nil
}

func scanVisits(rows *sql.Rows) []VisitResponse {
	visits := []VisitResponse{}
	for rows.Next() {
		var v VisitResponse
		var desc sql.NullString
		rows.Scan(&v.ID, &v.Date, &desc, &v.PetID)
		if desc.Valid {
			v.Description = &desc.String
		}
		visits = append(visits, v)
	}
	return visits
}

// --- Nested visit handlers ---

func verifyOwnerPet(db *sql.DB, w http.ResponseWriter, r *http.Request) (int, int, bool) {
	ownerID, _ := strconv.Atoi(PathID(r, "ownerId"))
	petID, _ := strconv.Atoi(PathID(r, "petId"))
	if _, err := findPetByIDAndOwner(db, petID, ownerID); err != nil {
		NotFound(w, "ERR-0003", "The requested pet does not exist.")
		return 0, 0, false
	}
	return ownerID, petID, true
}

func listVisitsByPet(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, petID, ok := verifyOwnerPet(db, w, r)
		if !ok {
			return
		}
		rows, err := db.Query("SELECT id, date, description, pet_id FROM visits WHERE pet_id = ?", petID)
		if err != nil {
			Error(w, http.StatusInternalServerError, "ERR-0001", err.Error())
			return
		}
		defer rows.Close()
		JSON(w, http.StatusOK, scanVisits(rows))
	}
}

func getVisitByPet(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, petID, ok := verifyOwnerPet(db, w, r)
		if !ok {
			return
		}
		visitID, _ := strconv.Atoi(PathID(r, "visitId"))
		visit, err := findVisitByID(db, visitID)
		if err != nil || visit.PetID != petID {
			NotFound(w, "ERR-0007", "The requested visit does not exist.")
			return
		}
		JSON(w, http.StatusOK, visit)
	}
}

func createVisitForPet(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, petID, ok := verifyOwnerPet(db, w, r)
		if !ok {
			return
		}
		var req CreateVisitRequest
		if err := DecodeJSON(r, &req); err != nil {
			BadRequest(w, "ERR-0001", "Invalid request body")
			return
		}
		result, err := db.Exec("INSERT INTO visits (date, description, pet_id) VALUES (?, ?, ?)",
			req.Date, nilStrPtr(req.Description), petID)
		if err != nil {
			Error(w, http.StatusInternalServerError, "ERR-0001", err.Error())
			return
		}
		id, _ := result.LastInsertId()
		visit, _ := findVisitByID(db, int(id))
		JSON(w, http.StatusCreated, visit)
	}
}

func updateVisitByPet(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, petID, ok := verifyOwnerPet(db, w, r)
		if !ok {
			return
		}
		visitID, _ := strconv.Atoi(PathID(r, "visitId"))
		existing, err := findVisitByID(db, visitID)
		if err != nil || existing.PetID != petID {
			NotFound(w, "ERR-0007", "The requested visit does not exist.")
			return
		}
		var req UpdateVisitRequest
		if err := DecodeJSON(r, &req); err != nil {
			BadRequest(w, "ERR-0001", "Invalid request body")
			return
		}
		db.Exec("UPDATE visits SET date=?, description=? WHERE id=?", req.Date, nilStrPtr(req.Description), visitID)
		visit, _ := findVisitByID(db, visitID)
		JSON(w, http.StatusOK, visit)
	}
}

func deleteVisitByPet(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, petID, ok := verifyOwnerPet(db, w, r)
		if !ok {
			return
		}
		visitID, _ := strconv.Atoi(PathID(r, "visitId"))
		existing, err := findVisitByID(db, visitID)
		if err != nil || existing.PetID != petID {
			NotFound(w, "ERR-0007", "The requested visit does not exist.")
			return
		}
		db.Exec("DELETE FROM visits WHERE id = ?", visitID)
		NoContent(w)
	}
}

// --- Global visit handlers ---

func listAllVisits(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT id, date, description, pet_id FROM visits")
		if err != nil {
			Error(w, http.StatusInternalServerError, "ERR-0001", err.Error())
			return
		}
		defer rows.Close()
		JSON(w, http.StatusOK, scanVisits(rows))
	}
}

func getVisitGlobal(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, _ := strconv.Atoi(PathID(r, "id"))
		visit, err := findVisitByID(db, id)
		if err != nil {
			NotFound(w, "ERR-0007", "The requested visit does not exist.")
			return
		}
		JSON(w, http.StatusOK, visit)
	}
}

func createVisitGlobal(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CreateVisitRequest
		if err := DecodeJSON(r, &req); err != nil {
			BadRequest(w, "ERR-0001", "Invalid request body")
			return
		}
		if req.PetID == nil {
			BadRequest(w, "ERR-0001", "petId is required")
			return
		}
		var exists int
		if err := db.QueryRow("SELECT id FROM pets WHERE id = ?", *req.PetID).Scan(&exists); err != nil {
			NotFound(w, "ERR-0003", "The requested pet does not exist.")
			return
		}
		result, err := db.Exec("INSERT INTO visits (date, description, pet_id) VALUES (?, ?, ?)",
			req.Date, nilStrPtr(req.Description), *req.PetID)
		if err != nil {
			Error(w, http.StatusInternalServerError, "ERR-0001", err.Error())
			return
		}
		id, _ := result.LastInsertId()
		visit, _ := findVisitByID(db, int(id))
		JSON(w, http.StatusCreated, visit)
	}
}

func updateVisitGlobal(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, _ := strconv.Atoi(PathID(r, "id"))
		if _, err := findVisitByID(db, id); err != nil {
			NotFound(w, "ERR-0007", "The requested visit does not exist.")
			return
		}
		var req UpdateVisitRequest
		if err := DecodeJSON(r, &req); err != nil {
			BadRequest(w, "ERR-0001", "Invalid request body")
			return
		}
		db.Exec("UPDATE visits SET date=?, description=? WHERE id=?", req.Date, nilStrPtr(req.Description), id)
		visit, _ := findVisitByID(db, id)
		JSON(w, http.StatusOK, visit)
	}
}

func deleteVisitGlobal(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, _ := strconv.Atoi(PathID(r, "id"))
		if _, err := findVisitByID(db, id); err != nil {
			NotFound(w, "ERR-0007", "The requested visit does not exist.")
			return
		}
		db.Exec("DELETE FROM visits WHERE id = ?", id)
		NoContent(w)
	}
}

func nilStrPtr(s *string) any {
	if s == nil {
		return nil
	}
	return *s
}
