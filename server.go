package main

import (
	"html/template"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Templates struct {
	templates *template.Template
}

type Participant struct {
	Name  string
	Email string
}

type Question struct {
	Question    string
	Votes       []Vote
	FinalAnswer int
}

type Vote struct {
	Vote  int
	Voter Participant
}

type Game struct {
	Participants []Participant
	Questions    []Question
	Leader       Participant
}

type Session struct {
	Game   Game
	Player Participant
}

func (t *Templates) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func newTemplate() *Templates {
	return &Templates{
		templates: template.Must(template.ParseGlob("views/*.html")),
	}
}

type Count struct {
	Count int
}

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Static("/static", "static")

	// Setting up init game data for testing
	session := Session{}
	game := Game{}
	participant1 := Participant{Name: "John", Email: "joe@gmail.com"}
	participant2 := Participant{Name: "Jane", Email: "jane@gmail.com"}
	game.Participants = append(game.Participants, participant1, participant2)
	question1 := Question{Question: "Building a toolshed"}
	question2 := Question{Question: "Painting the toolshed"}
	question1.Votes = append(question1.Votes, Vote{Vote: 5, Voter: participant1}, Vote{Vote: 3, Voter: participant2})
	question2.Votes = append(question2.Votes, Vote{Vote: 1, Voter: participant1}, Vote{Vote: 2, Voter: participant2})
	game.Questions = append(game.Questions, question1, question2)
	session.Game = game

	e.Renderer = newTemplate()

	e.POST("/register_player", func(c echo.Context) error {
		participant := Participant{}
		participant.Name = c.FormValue("name")
		participant.Email = c.FormValue("email")
		game.Participants = append(game.Participants, participant)
		return c.Render(http.StatusOK, "participants", game.Participants)
	})

	e.GET("/", func(c echo.Context) error {
		return c.Render(http.StatusOK, "index", session)
	})

	e.Logger.Fatal(e.Start(":1323"))
}
