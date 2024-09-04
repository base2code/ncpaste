package request_server_go

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"cloud.google.com/go/storage"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
)

func init() {
	functions.HTTP("RequestServer", RequestServer)
}

// HelloHTTP is an HTTP Cloud Function with a request parameter.
func RequestServer(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("RequestServer\n")
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
		randomStr := GenerateRandomString(10)

		for i := 0; i < 10; i++ {
			if _, err := bucket.Object(randomStr).Attrs(ctx); err == nil {
				randomStr = GenerateRandomString(10)
			} else {
				break
			}
		}

		obj := bucket.Object(randomStr)
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

		_, err = io.WriteString(w, randomStr)
		if err != nil {
			http.Error(w, "Error writing string", http.StatusInternalServerError)
			return
		}
	}
}

const charset = "abcdefghijklmnopqrstuvwxyz" +
	"1234567890"

func stringWithCharset(length int, charset string) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func GenerateRandomString(length int) string {
	return stringWithCharset(length, charset)
}
