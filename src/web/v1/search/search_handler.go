package web

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/anikethz/HertzDB/src/core/index"
	"github.com/anikethz/HertzDB/src/core/utils"
	hertzTypes "github.com/anikethz/HertzDB/src/web/types"
)

func (apiConfig *ApiConfig) SearchHandler(w http.ResponseWriter, r *http.Request) {

	body := hertzTypes.SearchRequest{}
	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&body)
	if err != nil {
		utils.ResponseWithError(w, 400, fmt.Sprintf("Error parsing JSON: %s", err))
		return
	}

	var res [][2]int64
	for _, v := range body.Field.Values {
		_res, _ := index.SearchTerm(apiConfig.Filename, body.Field.Name, v)
		res = append(res, _res...)
	}

	docs, _ := index.GetDocument(apiConfig.Json_Filename, res)

	response := ""

	for _, doc := range docs {
		response = response + "," + doc
	}

	if len(response) > 1 {
		response = "[" + response[1:len(response)-1] + "]"
	} else {
		response = "[]"
	}

	utils.RespondWithJson(w, 200, response)

}
