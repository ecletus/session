package manager

import (
	"net/http"

	"github.com/moisespsena-go/xroute"
	"github.com/ecletus/core"
	"github.com/ecletus/session"
	"github.com/ecletus/session/gorilla"
)

type RequestSessionManager struct {
	Manager session.ManagerInterface
	writer  http.ResponseWriter
	request *http.Request
}

func (sm *RequestSessionManager) ResponseWriter() http.ResponseWriter {
	return sm.writer
}

func (sm *RequestSessionManager) Request() *http.Request {
	return sm.request
}

// Add value to session data, if value is not string, will marshal it into JSON encoding and save it into session data.
func (sm *RequestSessionManager) Add(key string, value interface{}) error {
	return sm.Manager.Add(sm.writer, sm.request, key, value)
}

// Get value from session data
func (sm *RequestSessionManager) Get(key string) string {
	return sm.Manager.Get(sm.request, key)
}

// Pop value from session data
func (sm *RequestSessionManager) Pop(key string) string {
	return sm.Manager.Pop(sm.writer, sm.request, key)
}

// Flash add flash message to session data
func (sm *RequestSessionManager) Flash(message session.Message) error {
	return sm.Manager.Flash(sm.writer, sm.request, message)
}

// Flashes returns a slice of flash messages from session data
func (sm *RequestSessionManager) Flashes() []session.Message {
	messages := sm.Manager.Flashes(sm.writer, sm.request)
	return messages
}

// Load get value from session data and unmarshal it into result
func (sm *RequestSessionManager) Load(key string, result interface{}) error {
	return sm.Manager.Load(sm.request, key, result)
}

// PopLoad pop value from session data and unmarshal it into result
func (sm *RequestSessionManager) PopLoad(key string, result interface{}) error {
	return sm.Manager.PopLoad(sm.writer, sm.request, key, result)
}

func (sm *RequestSessionManager) Middleware(handler http.Handler) http.Handler {
	return sm.Manager.Middleware(handler)
}

func Middleware(setupConfig *core.SetupConfig) *xroute.Middleware {
	return &xroute.Middleware{
		Name: "qor:session",
		Handler: func(chain *xroute.ChainHandler) {
			context := core.ContextFromRequest(chain.Request())
			rsm := context.SessionManager()
			if rsm == nil {
				cookieStore := setupConfig.CookieStoreFactory()(context, nil, nil)
				sm := gorilla.New("_session", cookieStore)
				rsm = &RequestSessionManager{sm, context.Writer, context.Request}
				context.SetSessionManager(rsm)
			}

			handler := rsm.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				context.Request = r
				chain.SetRequest(r)
				chain.Pass()
			}))
			handler.ServeHTTP(chain.Writer, chain.Request())
		},
	}
}
