package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	_ "modernc.org/sqlite"
)

type MagicEvent struct {
	ID        string
	Name      string
	CreatedAt time.Time
}

type Ticket struct {
	ID          int
	EventID     string
	Origin      string
	Destination string
	Status      string
	Courier     string
	CreatedAt   time.Time
}

var db *sql.DB
var tmpl *template.Template

func main() {
	var err error

	db, err = sql.Open("sqlite", "database.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	createTables()

	tmpl = template.Must(template.ParseGlob("templates/*.html"))

	http.HandleFunc("/", handleIndex)
	http.HandleFunc("/new-event", handleNewEvent)
	http.HandleFunc("/e/", handleEventDispatch)
	http.HandleFunc("/ticket/add", handleAddTicket)
	http.HandleFunc("/ticket/claim", handleClaimTicket)
	http.HandleFunc("/ticket/update", handleUpdateTicket)

	log.Println("Server started on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func createTables() {
	query := `
	CREATE TABLE IF NOT EXISTS events (
		id TEXT PRIMARY KEY,
		name TEXT,
		created_at DATETIME
	);
	CREATE TABLE IF NOT EXISTS tickets (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		event_id TEXT,
		origin TEXT,
		destination TEXT,
		status TEXT DEFAULT 'Pending',
		courier TEXT DEFAULT '',
		created_at DATETIME,
		FOREIGN KEY(event_id) REFERENCES events(id)
	);`
	_, err := db.Exec(query)
	if err != nil {
		log.Fatal("DB Init Error:", err)
	}
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "index.html", nil)
}

func handleNewEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	eventName := r.FormValue("name")
	magicID := uuid.New().String()

	_, err := db.Exec("INSERT INTO events (id, name, created_at) VALUES (?, ?, ?)", magicID, eventName, time.Now())
	if err != nil {
		http.Error(w, "Failed to create event", 500)
		return
	}

	http.Redirect(w, r, "/e/"+magicID, http.StatusSeeOther)
}

func handleEventDispatch(w http.ResponseWriter, r *http.Request) {
	eventID := r.URL.Path[len("/e/"):]

	var event MagicEvent

	err := db.QueryRow("SELECT id, name, created_at FROM events WHERE id = ?", eventID).Scan(&event.ID, &event.Name, &event.CreatedAt)
	if err != nil {
		http.Error(w, "Event not found", 404)
		return
	}

	rows, _ := db.Query("SELECT id, origin, destination, status, courier, created_at FROM tickets WHERE event_id = ? ORDER BY created_at DESC", eventID)
	defer rows.Close()

	var tickets []Ticket

	for rows.Next() {
		var t Ticket
		rows.Scan(&t.ID, &t.Origin, &t.Destination, &t.Status, &t.Courier, &t.CreatedAt)
		t.EventID = eventID
		tickets = append(tickets, t)
	}

	data := struct {
		Event   MagicEvent
		Tickets []Ticket
	}{event, tickets}

	tmpl.ExecuteTemplate(w, "dispatch.html", data)
}

func handleAddTicket(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		eventID := r.FormValue("event_id")
		origin := r.FormValue("origin")
		dest := r.FormValue("destination")

		db.Exec("INSERT INTO tickets (event_id, origin, destination, created_at) VALUES (?, ?, ?, ?)", eventID, origin, dest, time.Now())

		http.Redirect(w, r, "/e/"+eventID, http.StatusSeeOther)
	}
}

func handleClaimTicket(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		id := r.FormValue("ticket_id")
		eventID := r.FormValue("event_id")
		courier := r.FormValue("courier_name")

		db.Exec("UPDATE tickets SET courier = ?, status = 'In Transit' WHERE id = ?", courier, id)

		http.Redirect(w, r, "/e/"+eventID, http.StatusSeeOther)
	}
}

func handleUpdateTicket(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		id := r.FormValue("ticket_id")
		eventID := r.FormValue("event_id")
		status := r.FormValue("status")

		db.Exec("UPDATE tickets SET status = ? WHERE id = ?", status, id)

		http.Redirect(w, r, "/e/"+eventID, http.StatusSeeOther)
	}
}
