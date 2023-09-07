package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/trieudh58/ethereum-parser/parser"
)

type ApiServer struct {
	parser *parser.EthereumParser
	server *http.Server
}

func NewApiServer(parser *parser.EthereumParser) *ApiServer {
	return &ApiServer{
		parser: parser,
	}
}

func (s *ApiServer) Stop(ctx context.Context) {
	if err := s.server.Shutdown(ctx); err != nil {
		log.Println(err)
	}
}

func (s *ApiServer) ListenAndServe(port int) error {
	s.server = &http.Server{
		Addr: fmt.Sprintf(":%d", apiPort),
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.URL.Path {
			case "/subscribe":
				// POST /subscribe
				s.handleSubsribe(w, r)
			case "/transactions":
				// GET /transactions?address=0x...
				s.handleGetTransactions(w, r)
			default:
				http.NotFound(w, r)
			}
		}),
	}

	if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

func (s *ApiServer) handleGetTransactions(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			log.Println(err)
		}
	}()
	if r.Method != "GET" {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	address := r.URL.Query().Get("address")
	txs := s.parser.GetTransactions(address)

	response := map[string]interface{}{
		"txs": txs,
	}
	jsonData, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Failed to marshal JSON", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func (s *ApiServer) handleSubsribe(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			log.Println(err)
		}
	}()
	if r.Method != "POST" {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	var data struct {
		Address string `json:"address"`
	}
	err = json.Unmarshal(body, &data)
	if err != nil {
		http.Error(w, "Failed to unmarshal JSON", http.StatusBadRequest)
		return
	}

	success := s.parser.Subscribe(data.Address)

	response := map[string]interface{}{
		"success": success,
	}
	jsonData, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Failed to marshal JSON", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}
