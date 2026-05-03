package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

/*
type validationResponse struct {
	Valid bool `json:"valid"`
}
*/

type cleanString struct {
	CleanedBody string `json:"cleaned_body"`
}

type errorResponse struct {
	Error string `json:"error"`
}

type requestParams struct {
	Body string `json:"body"`
}

const maxChars int = 140

func (cfg *apiConfig) handlerValidateChirp(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	params := requestParams{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		respondWithError(w, http.StatusInternalServerError, "unable to decode parameters", err)
		return
	}
	if len(params.Body) > maxChars {
		log.Printf("Error, body too long")
		respondWithError(w, http.StatusBadRequest, "length must be 140 chars or less", err)
		return
	}
	censoredString := censorString(params.Body)

	respondWithJSON(w, http.StatusOK, cleanString{CleanedBody: censoredString})
}

func respondWithError(w http.ResponseWriter, code int, msg string, err error) {
	if err != nil {
		log.Println(err)
	}
	if code > 499 {
		log.Printf("Responding with 5XX error: %s", msg)
	}
	respondWithJSON(w, code, errorResponse{
		Error: msg,
	})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("error building response: %s\n", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data)
}

func censorString(s string) string {
	var profanity = []string{"kerfuffle", "sharbert", "fornax"}
	words := strings.Split(s, " ")

	newWords := make([]string, 0, len(words))
	for _, word := range words {
		newWord := word
		for _, badWord := range profanity {
			if strings.ToLower(word) == badWord {
				newWord = "****"
				break
			}
		}
		newWords = append(newWords, newWord)
	}

	return strings.Join(newWords, " ")
}
