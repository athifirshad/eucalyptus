package main

// import (
// 	"fmt"
// 	_ "image/jpeg"
// 	"io"
// 	"net/http"
// 	"os"
// )

// func (app *application) qrCodeHandler(w http.ResponseWriter, r *http.Request) {
// 	// Check if the request method is POST
// 	if r.Method != http.MethodPost {
// 		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
// 		return
// 	}

// 	// Parse the multipart form in the request
// 	err := r.ParseMultipartForm(10 << 20) // 10 MB max memory
// 	if err != nil {
// 		http.Error(w, "Could not parse multipart form", http.StatusBadRequest)
// 		return
// 	}

// 	// Retrieve the file from form data
// 	file, handler, err := r.FormFile("image")
// 	if err != nil {
// 		http.Error(w, "Error retrieving the file", http.StatusBadRequest)
// 		return
// 	}
// 	defer file.Close()

// 	// Create a temporary file to store the uploaded file
// 	tempFile, err := os.CreateTemp("", "upload-*.jpg")
// 	if err != nil {
// 		http.Error(w, "Could not create temporary file", http.StatusInternalServerError)
// 		return
// 	}
// 	defer os.Remove(tempFile.Name()) // Clean up

// 	// Copy the uploaded file to the temporary file
// 	_, err = io.Copy(tempFile, file)
// 	if err != nil {
// 		http.Error(w, "Error copying file", http.StatusInternalServerError)
// 		return
// 	}

// 	// Call the ReadQr function to read the QR code from the image
// 	decodedText := ReadQr(tempFile.Name())

// 	// Respond with the decoded QR code text
// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusOK)
// 	fmt.Fprintf(w, `{"qr_code_text": "%s"}`, decodedText)
// }
