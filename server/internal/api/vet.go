package api

import (
	"database/sql"
	"net/http"
	"sort"
	"strconv"
)

type VetResponse struct {
	ID          int                 `json:"id"`
	FirstName   string              `json:"firstName"`
	LastName    string              `json:"lastName"`
	Specialties []SpecialtyResponse `json:"specialties"`
}

type SpecialtyResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type CreateVetRequest struct {
	FirstName    string `json:"firstName"`
	LastName     string `json:"lastName"`
	SpecialtyIds []int  `json:"specialtyIds"`
}

type UpdateVetRequest = CreateVetRequest

type CreateSpecialtyRequest struct {
	Name string `json:"name"`
}

type UpdateSpecialtyRequest = CreateSpecialtyRequest

func RegisterVetRoutes(mux *http.ServeMux, db *sql.DB) {
	// Specialties
	mux.HandleFunc("GET /api/v1/specialties", listSpecialties(db))
	mux.HandleFunc("GET /api/v1/specialties/{id}", getSpecialty(db))
	mux.HandleFunc("POST /api/v1/specialties", createSpecialty(db))
	mux.HandleFunc("PUT /api/v1/specialties/{id}", updateSpecialty(db))
	mux.HandleFunc("DELETE /api/v1/specialties/{id}", deleteSpecialty(db))

	// Vets
	mux.HandleFunc("GET /api/v1/vets", listVets(db))
	mux.HandleFunc("GET /api/v1/vets/{id}", getVet(db))
	mux.HandleFunc("POST /api/v1/vets", createVet(db))
	mux.HandleFunc("PUT /api/v1/vets/{id}", updateVet(db))
	mux.HandleFunc("DELETE /api/v1/vets/{id}", deleteVet(db))
}

// --- Specialty handlers ---

func listSpecialties(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT id, name FROM specialties")
		if err != nil {
			Error(w, http.StatusInternalServerError, "ERR-0001", err.Error())
			return
		}
		defer rows.Close()
		specs := []SpecialtyResponse{}
		for rows.Next() {
			var s SpecialtyResponse
			rows.Scan(&s.ID, &s.Name)
			specs = append(specs, s)
		}
		JSON(w, http.StatusOK, specs)
	}
}

func getSpecialty(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(PathID(r, "id"))
		if err != nil {
			NotFound(w, "ERR-0006", "The requested specialty does not exist.")
			return
		}
		var s SpecialtyResponse
		if err := db.QueryRow("SELECT id, name FROM specialties WHERE id = ?", id).Scan(&s.ID, &s.Name); err != nil {
			NotFound(w, "ERR-0006", "The requested specialty does not exist.")
			return
		}
		JSON(w, http.StatusOK, s)
	}
}

func createSpecialty(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CreateSpecialtyRequest
		if err := DecodeJSON(r, &req); err != nil {
			BadRequest(w, "ERR-0001", "Invalid request body")
			return
		}
		result, err := db.Exec("INSERT INTO specialties (name) VALUES (?)", req.Name)
		if err != nil {
			Error(w, http.StatusInternalServerError, "ERR-0001", err.Error())
			return
		}
		id, _ := result.LastInsertId()
		JSON(w, http.StatusCreated, SpecialtyResponse{ID: int(id), Name: req.Name})
	}
}

func updateSpecialty(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(PathID(r, "id"))
		if err != nil {
			NotFound(w, "ERR-0006", "The requested specialty does not exist.")
			return
		}
		var exists int
		if err := db.QueryRow("SELECT id FROM specialties WHERE id = ?", id).Scan(&exists); err != nil {
			NotFound(w, "ERR-0006", "The requested specialty does not exist.")
			return
		}
		var req UpdateSpecialtyRequest
		if err := DecodeJSON(r, &req); err != nil {
			BadRequest(w, "ERR-0001", "Invalid request body")
			return
		}
		db.Exec("UPDATE specialties SET name = ? WHERE id = ?", req.Name, id)
		JSON(w, http.StatusOK, SpecialtyResponse{ID: id, Name: req.Name})
	}
}

func deleteSpecialty(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(PathID(r, "id"))
		if err != nil {
			NotFound(w, "ERR-0006", "The requested specialty does not exist.")
			return
		}
		var exists int
		if err := db.QueryRow("SELECT id FROM specialties WHERE id = ?", id).Scan(&exists); err != nil {
			NotFound(w, "ERR-0006", "The requested specialty does not exist.")
			return
		}
		db.Exec("DELETE FROM vet_specialties WHERE specialty_id = ?", id)
		db.Exec("DELETE FROM specialties WHERE id = ?", id)
		NoContent(w)
	}
}

// --- Vet helpers ---

func findVetByID(db *sql.DB, id int) (*VetResponse, error) {
	var v VetResponse
	err := db.QueryRow("SELECT id, first_name, last_name FROM vets WHERE id = ?", id).
		Scan(&v.ID, &v.FirstName, &v.LastName)
	if err != nil {
		return nil, err
	}
	v.Specialties = loadVetSpecialties(db, v.ID)
	return &v, nil
}

func loadVetSpecialties(db *sql.DB, vetID int) []SpecialtyResponse {
	rows, err := db.Query(`
		SELECT s.id, s.name FROM specialties s
		JOIN vet_specialties vs ON s.id = vs.specialty_id
		WHERE vs.vet_id = ?`, vetID)
	if err != nil {
		return []SpecialtyResponse{}
	}
	defer rows.Close()
	specs := []SpecialtyResponse{}
	for rows.Next() {
		var s SpecialtyResponse
		rows.Scan(&s.ID, &s.Name)
		specs = append(specs, s)
	}
	sort.Slice(specs, func(i, j int) bool { return specs[i].ID < specs[j].ID })
	return specs
}

func setVetSpecialties(db *sql.DB, vetID int, specialtyIds []int) {
	db.Exec("DELETE FROM vet_specialties WHERE vet_id = ?", vetID)
	for _, sID := range specialtyIds {
		var exists int
		if err := db.QueryRow("SELECT id FROM specialties WHERE id = ?", sID).Scan(&exists); err == nil {
			db.Exec("INSERT INTO vet_specialties (vet_id, specialty_id) VALUES (?, ?)", vetID, sID)
		}
	}
}

// --- Vet handlers ---

func listVets(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT id, first_name, last_name FROM vets")
		if err != nil {
			Error(w, http.StatusInternalServerError, "ERR-0001", err.Error())
			return
		}
		defer rows.Close()
		vets := []VetResponse{}
		for rows.Next() {
			var v VetResponse
			rows.Scan(&v.ID, &v.FirstName, &v.LastName)
			v.Specialties = loadVetSpecialties(db, v.ID)
			vets = append(vets, v)
		}
		JSON(w, http.StatusOK, vets)
	}
}

func getVet(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(PathID(r, "id"))
		if err != nil {
			NotFound(w, "ERR-0005", "The requested vet does not exist.")
			return
		}
		vet, err := findVetByID(db, id)
		if err != nil {
			NotFound(w, "ERR-0005", "The requested vet does not exist.")
			return
		}
		JSON(w, http.StatusOK, vet)
	}
}

func createVet(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CreateVetRequest
		if err := DecodeJSON(r, &req); err != nil {
			BadRequest(w, "ERR-0001", "Invalid request body")
			return
		}
		result, err := db.Exec("INSERT INTO vets (first_name, last_name) VALUES (?, ?)", req.FirstName, req.LastName)
		if err != nil {
			Error(w, http.StatusInternalServerError, "ERR-0001", err.Error())
			return
		}
		id, _ := result.LastInsertId()
		if len(req.SpecialtyIds) > 0 {
			setVetSpecialties(db, int(id), req.SpecialtyIds)
		}
		vet, _ := findVetByID(db, int(id))
		JSON(w, http.StatusCreated, vet)
	}
}

func updateVet(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(PathID(r, "id"))
		if err != nil {
			NotFound(w, "ERR-0005", "The requested vet does not exist.")
			return
		}
		if _, err := findVetByID(db, id); err != nil {
			NotFound(w, "ERR-0005", "The requested vet does not exist.")
			return
		}
		var req UpdateVetRequest
		if err := DecodeJSON(r, &req); err != nil {
			BadRequest(w, "ERR-0001", "Invalid request body")
			return
		}
		db.Exec("UPDATE vets SET first_name=?, last_name=? WHERE id=?", req.FirstName, req.LastName, id)
		if req.SpecialtyIds != nil {
			setVetSpecialties(db, id, req.SpecialtyIds)
		}
		vet, _ := findVetByID(db, id)
		JSON(w, http.StatusOK, vet)
	}
}

func deleteVet(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(PathID(r, "id"))
		if err != nil {
			NotFound(w, "ERR-0005", "The requested vet does not exist.")
			return
		}
		if _, err := findVetByID(db, id); err != nil {
			NotFound(w, "ERR-0005", "The requested vet does not exist.")
			return
		}
		db.Exec("DELETE FROM vet_specialties WHERE vet_id = ?", id)
		db.Exec("DELETE FROM vets WHERE id = ?", id)
		NoContent(w)
	}
}
