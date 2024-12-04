package ingest

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/anikethz/HertzDB/src/core/index"
	"github.com/anikethz/HertzDB/src/core/utils"
	"github.com/anikethz/HertzDB/src/web/types"
)

func IngestionHandler(w http.ResponseWriter, r *http.Request, apiConfig *types.ApiConfig) {

	// r.Body = http.MaxBytesReader(w, r.Body, 10<<20) // 10 MB limit

	mr, err := r.MultipartReader()
	if err != nil {
		utils.ResponseWithError(w, 400, "Invalid multipart request")
		return
	}

	for {
		part, err := mr.NextPart()
		if err == io.EOF {
			break // No more parts
		}
		if err != nil {
			utils.ResponseWithError(w, 400, "Error reading multipart data")
			return
		}

		if part.FileName() != "" {
			fmt.Printf("Processing file: %s\n", part.FileName())

			dst, err := os.Create(fmt.Sprintf("./uploads/%s", apiConfig.Json_Filename))
			if err != nil {
				utils.ResponseWithError(w, 500, "Unable to create destination file")
				return
			}
			defer dst.Close()

			// Stream the file in chunks to the destination
			buf := make([]byte, 1024) // 1 KB buffer
			for {
				n, err := part.Read(buf)
				if err != nil && err != io.EOF {
					utils.ResponseWithError(w, 500, "Error reading file stream")
					return
				}
				if n == 0 {
					break // End of file
				}

				if _, err := dst.Write(buf[:n]); err != nil {
					utils.ResponseWithError(w, 500, "Error writing file to disk")
					return
				}
			}

			log.Printf("File %s uploaded successfully\n", part.FileName())
		}
	}

	indexDocument, _ := index.NewIndexDocument(apiConfig.Filename, apiConfig.Json_Filename)
	indexDocument.ParseEntireFile([]string{"title", "cast"})

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("File(s) uploaded successfully"))

}
