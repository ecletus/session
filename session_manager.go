package session

import (
	"fmt"
	"net/http"

	"github.com/moisespsena-go/i18n-modular/i18nmod"
	"github.com/moisespsena/template/html/template"
)

// ManagerInterface session manager interface
type ManagerInterface interface {
	// Add value to session data, if value is not string, will marshal it into JSON encoding and save it into session data.
	Add(w http.ResponseWriter, req *http.Request, key string, value interface{}) error
	// Get value from session data
	Get(req *http.Request, key string) string
	// Pop value from session data
	Pop(w http.ResponseWriter, req *http.Request, key string) string

	// Flash add flash message to session data
	Flash(w http.ResponseWriter, req *http.Request, message Message) error
	// Flashes returns a slice of flash messages from session data
	Flashes(w http.ResponseWriter, req *http.Request) []Message

	// Load get value from session data and unmarshal it into result
	Load(req *http.Request, key string, result interface{}) error
	// PopLoad pop value from session data and unmarshal it into result
	PopLoad(w http.ResponseWriter, req *http.Request, key string, result interface{}) error

	// Middleware returns a new session manager middleware instance.
	Middleware(http.Handler) http.Handler
}

// Message message struct
type Message struct {
	Message template.HTML
	Type    string
}

type RequestSessionManager interface {
	// Add value to session data, if value is not string, will marshal it into JSON encoding and save it into session data.
	Add(key string, value interface{}) error
	// Get value from session data
	Get(key string) string
	// Pop value from session data
	Pop(key string) string

	// Flash add flash message to session data
	Flash(message Message) error
	// Flashes returns a slice of flash messages from session data
	Flashes() []Message

	// Load get value from session data and unmarshal it into result
	Load(key string, result interface{}) error
	// PopLoad pop value from session data and unmarshal it into result
	PopLoad(key string, result interface{}) error

	Middleware(http.Handler) http.Handler

	ResponseWriter() http.ResponseWriter

	Request() *http.Request
}

func TranslatedMessage(ctx i18nmod.Context, msg interface{}, typ string) Message {
	switch t := msg.(type) {
	case i18nmod.Translater:
		return Message{template.HTML(t.Translate(ctx)), typ}
	case error:
		return Message{ctx.T(t.Error()).GetHtml(), typ}
	case string:
		return Message{ctx.T(t).GetHtml(), typ}
	default:
		return Message{ctx.T(fmt.Sprint(t)).GetHtml(), typ}
	}
}

func TranslatedMessageE(ctx i18nmod.Context, msg interface{}) Message {
	return TranslatedMessage(ctx, msg, "error")
}
