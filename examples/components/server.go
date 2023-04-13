package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mbertschler/blocks/html"
)

type Server struct {
	engine      *gin.Engine
	guiapi      *Handler
	sessions    SessionStorage
	withSession *gin.RouterGroup
}

type SessionStorage interface {
	GetSession(id string) (*Session, error)
	SetSession(s *Session) error
}

func NewServer(storage SessionStorage) *Server {
	gin.SetMode(gin.ReleaseMode)

	engine := gin.New()
	engine.Use(gin.Recovery())
	guiapi := NewGuiapi()

	s := &Server{
		engine:   engine,
		guiapi:   guiapi,
		sessions: storage,
	}

	withSession := engine.Group("")
	withSession.Use(s.sessionMiddleware)
	withSession.POST("/guiapi", guiapi.Handle)

	s.withSession = withSession
	return s
}

// Static serves static files from the given directory.
func (s *Server) Static(path, dir string) {
	s.engine.Static(path, dir)
}

func (s *Server) SetFunc(name string, fn Callable) {
	s.guiapi.SetFunc(name, fn)
}

type PageFunc func(*gin.Context) (html.Block, error)

func (s *Server) Page(path string, page PageFunc) {
	s.withSession.GET(path, func(c *gin.Context) {
		pageBlock, err := page(c)
		if err != nil {
			log.Println("Page error:", err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		err = html.RenderMinified(c.Writer, pageBlock)
		if err != nil {
			log.Println("RenderMinified error:", err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
	})
}

func (s *Server) Handler() http.Handler {
	return s.engine.Handler()
}

const sessionCookie = "session"

func sessionFromContext(c *gin.Context) *Session {
	sess, ok := c.Keys["session"].(*Session)
	if ok && sess != nil {
		return sess
	}
	if !ok {
		log.Printf("bad session %#v", c.Keys["session"])
	}
	if sess == nil {
		sess = &Session{}
	}
	return sess
}

func (s *Server) sessionMiddleware(c *gin.Context) {
	cookie, err := c.Request.Cookie(sessionCookie)
	if err != nil && err != http.ErrNoCookie {
		log.Println("sessionMiddleware.Cookie:", err)
	}
	id := ""
	if cookie != nil {
		id = cookie.Value
	}
	sess, err := s.sessions.GetSession(id)
	if err != nil {
		log.Println("sessionMiddleware.GetSession:", err)
	}
	if sess.New {
		http.SetCookie(c.Writer, &http.Cookie{
			Name:     sessionCookie,
			Value:    sess.ID,
			Path:     "/",
			HttpOnly: true,
			Expires:  time.Now().Add(30 * 24 * time.Hour),
		})
	}
	c.Keys = map[string]interface{}{
		"session": sess,
	}
	c.Next()
	err = s.sessions.SetSession(sess)
	if err != nil {
		log.Println("sessionMiddleware.SetSession:", err)
	}
}

type Component interface {
	Component() *ComponentConfig
}

type ComponentConfig struct {
	Name    string
	Actions map[string]Callable
}

func (s *Server) RegisterComponent(c Component) {
	config := c.Component()
	for name, fn := range config.Actions {
		s.SetFunc(config.Name+"."+name, fn)
	}
}

func (s *Server) RegisterPage(path string, fn PageFunc) {
	s.Page(path, fn)
}
