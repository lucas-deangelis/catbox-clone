package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type config struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (c *config) Authenticate(r *http.Request) bool {
	username, password, ok := r.BasicAuth()

	return ok && username == c.Username && password == c.Password
}

func loadConfig() (*config, error) {
	file, err := os.ReadFile("config.json")
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %v", err)
	}

	var cfg config
	err = json.Unmarshal(file, &cfg)
	if err != nil {
		return nil, fmt.Errorf("error parsing config file: %v", err)
	}

	if cfg.Username == "" || cfg.Password == "" {
		return nil, fmt.Errorf("username or password is empty in config file")
	}

	return &cfg, nil
}

func main() {
	cfg, err := loadConfig()
	if err != nil {
		fmt.Println("Error loading config:", err)
		os.Exit(1)
	}

	http.HandleFunc("POST /", authMiddleware(uploadHandler, cfg))
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("GET /{file}", fileHandler)
	fmt.Println("Server is running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("In home handler")
	html := `
		<html>
			<body>
				<h2>File Upload</h2>
				<form action="/" method="post" enctype="multipart/form-data">
					<input type="file" name="file">
					<input type="submit" value="Upload">
				</form>

				<p></p>
			</body>
		</html>
	`
	fmt.Fprintf(w, html)
}

func fileHandler(w http.ResponseWriter, r *http.Request) {
	fileName := r.PathValue("file")
	content, err := os.ReadFile(fileName)
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(content)
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("In upload handler")
	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Create a new file in the current directory
	dst, err := os.Create(header.Filename)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	// Copy the uploaded file to the destination file
	_, err = io.Copy(dst, file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "File uploaded successfully: %s", header.Filename)
}

func authMiddleware(next http.HandlerFunc, cfg *config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !cfg.Authenticate(r) {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	}
}
