package main

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var templates = template.Must(template.ParseGlob("templates/*.html"))
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Client struct {
	conn *websocket.Conn
	name string
	room string
}

var rooms = make(map[string][]*Client)

// ------------------- MongoDB Setup -------------------
var client *mongo.Client
var userCollection *mongo.Collection

type User struct {
	Name     string `bson:"name" json:"name"`
	Email    string `bson:"email" json:"email"`
	Password string `bson:"password" json:"password"`
}

func init() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var err error
	client, err = mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		panic(err)
	}

	userCollection = client.Database("codesync").Collection("users")
}

// ------------------- Main -------------------
func main() {

	http.Handle("/static/", http.StripPrefix("/static/",
		http.FileServer(http.Dir("static"))))

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/signup", signupHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/dashboard", dashboardHandler)
	http.HandleFunc("/create", createHandler)
	http.HandleFunc("/join", joinHandler)
	http.HandleFunc("/editor", editorHandler)
	http.HandleFunc("/ws", wsHandler)

	fmt.Println("Server running at http://localhost:8080")

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Server error:", err)
	}
}

// ------------------- Handlers -------------------
func indexHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "index.html", nil)
}

func signupHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {
		templates.ExecuteTemplate(w, "signup.html", nil)
		return
	}

	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var existing User
	err = userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&existing)

	w.Header().Set("Content-Type", "application/json")

	if err == nil {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Email already exists",
		})
		return
	}

	_, err = userCollection.InsertOne(ctx, user)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Signup failed",
		})
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"name":    user.Name,
	})
}

func loginHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {
		templates.ExecuteTemplate(w, "login.html", nil)
		return
	}

	var loginData User
	err := json.NewDecoder(r.Body).Decode(&loginData)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user User
	err = userCollection.FindOne(ctx, bson.M{"email": loginData.Email}).Decode(&user)

	w.Header().Set("Content-Type", "application/json")

	if err != nil || user.Password != loginData.Password {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Invalid email or password",
		})
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"name":    user.Name,
	})
}

func dashboardHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "dashboard.html", nil)
}

func createHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "create.html", nil)
}

func joinHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "join.html", nil)
}

func editorHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "editor.html", nil)
}

// ------------------- WebSocket -------------------
func wsHandler(w http.ResponseWriter, r *http.Request) {

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	name := r.URL.Query().Get("name")
	room := r.URL.Query().Get("room")

	client := &Client{
		conn: conn,
		name: name,
		room: room,
	}

	rooms[room] = append(rooms[room], client)

	broadcast(room, "JOIN:"+name)

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			break
		}

		broadcast(room, string(msg))
	}
}

func broadcast(room string, message string) {
	for _, client := range rooms[room] {
		err := client.conn.WriteMessage(websocket.TextMessage, []byte(message))
		if err != nil {
			fmt.Println(err)
		}
	}
}
