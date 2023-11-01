package request_server

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"cloud.google.com/go/storage"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
)

func init() {
	functions.HTTP("RequestServer", RequestServer)
}

// HelloHTTP is an HTTP Cloud Function with a request parameter.
func RequestServer(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	// Get the file contents from the request body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}

	// Get path
	path := r.URL.Path
	fmt.Printf("Path: %s\n", path)

	// Create a new client and bucket object
	client, err := storage.NewClient(ctx)
	if err != nil {
		http.Error(w, "Error creating storage client", http.StatusInternalServerError)
		return
	}
	bucket := client.Bucket("testbucket-niklas")

	if len(body) == 0 && path == "/" {
		fmt.Printf("Show index.html\n")
		// Show index.html
		_, err := fmt.Fprint(w, "Hello, World!\n")
		if err != nil {
			http.Error(w, "Error writing string", http.StatusInternalServerError)
			return
		}
	} else if len(body) == 0 && path != "/" {
		// Show file
		fmt.Printf("Show file\n")
		path = strings.Replace(path, "/", "", 1)
		obj := bucket.Object(path)
		reader, err := obj.NewReader(ctx)
		if err != nil {
			http.Error(w, "Error reading file from bucket", http.StatusInternalServerError)
			return
		}

		data, err := ioutil.ReadAll(reader)

		if err != nil {
			http.Error(w, "Error reading file from bucket", http.StatusInternalServerError)
			return
		}
		_, err = io.WriteString(w, string(data))
		if err != nil {
			http.Error(w, "Error writing string", http.StatusInternalServerError)
			return
		}
		// return 200

	} else if len(body) != 0 && path == "/" {
		fmt.Printf("Write file\n")
		// Create a new object and write the file contents to it
		// Generate random string
		randomstr, err := uuid.NewRandom()
		if err != nil {
			http.Error(w, "Error generating random string", http.StatusInternalServerError)
			return
		}
		obj := bucket.Object(randomstr.String())
		writer := obj.NewWriter(ctx)
		_, err = writer.Write(body)
		if err != nil {
			http.Error(w, "Error writing file to bucket", http.StatusInternalServerError)
			return
		}
		if err := writer.Close(); err != nil {
			http.Error(w, "Error closing writer", http.StatusInternalServerError)
			return
		}

		_, err = io.WriteString(w, randomstr.String())
		if err != nil {
			http.Error(w, "Error writing string", http.StatusInternalServerError)
			return
		}
	}
}
