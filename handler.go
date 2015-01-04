package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"gossip/database"
	"gossip/response"

	"github.com/zenazn/goji/web"
)

var handler *Handler

type Handler struct {
	DB *database.DB
}

func init() {
	handler = &Handler{DB: database.New()}
}

func (h *Handler) SayHello(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "Hello, there.\n")
}

func (h *Handler) FindMessageById(c web.C, w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(c.URLParams["id"])
	if err != nil {
		http.Error(w, "You must submit an ID to lookup.", 400)
		return
	}

	message, err := h.DB.Find(id)

	if err != nil {
		http.Error(w, fmt.Sprintf("Could not fetch message id=%d: %s", id, err.Error()), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, response.New().WithMessage(message).Json())
}

func (h *Handler) FetchLatestMessages(c web.C, w http.ResponseWriter, r *http.Request) {
	limit := c.URLParams["limit"]
	if limit == "" { // no limit
		limit = "10"
	}
	fmt.Println("Fetching latest", limit, "messages")

	messages, err := h.DB.LatestMessages(limit)

	if err != nil {
		http.Error(w, fmt.Sprintf("Could not fetch latest messages with limit=%s: %s", limit, err.Error()), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, response.New().WithMessages(messages).WithLimit(limit).Json())
}

func (h *Handler) StoreMessage(w http.ResponseWriter, r *http.Request) {
	msg := database.Message{
		Room:    r.PostFormValue("room"),
		Author:  r.PostFormValue("author"),
		Message: r.PostFormValue("message"),
		At:      r.PostFormValue("time"),
	}

	log.Println("Storing the following message:", msg.String())

	message, err := h.DB.InsertMessage(msg)

	if err != nil {
		http.Error(w, fmt.Sprintf("Could not insert message %s: %s", msg.String(), err.Error()), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, response.New().WithMessage(message).Json())
}