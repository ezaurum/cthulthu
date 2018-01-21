package identity

import (
	"github.com/ezaurum/cthulthu/database"
	"github.com/ezaurum/cthulthu/route"
	"github.com/ezaurum/cthulthu/session"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"sort"
)

type Score struct {
	database.Model
	Rank          int
	TeamNo        int `form:"teamNo" binding:"required"`
	TeamName      string `form:"teamName" binding:"required"`
	TotalScore    int
	Curling       int `form:"curling"`
	FigureSkating int `form:"figureSkating"`
	Luge          int `form:"luge"`
	SkiJumping    int `form:"skiJumping"`
	IceHockey     int `form:"iceHockey"`
	Biathlon      int `form:"biathlon"`
}

func (s *Score) calculateTotal() *Score {
	s.TotalScore = s.Curling + s.FigureSkating + s.SkiJumping + s.IceHockey + s.Biathlon + s.Luge
	return s
}
// ByAge implements sort.Interface for []Person based on
// the Age field.
type ByScore []Score

func (a ByScore) Len() int           { return len(a) }
func (a ByScore) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByScore) Less(i, j int) bool { return a[i].TotalScore > a[j].TotalScore }

func ScoreRoute() route.Routes {

	rt := make(route.Routes)
	rt.GET("/scores", func(c *gin.Context) {
		var list []Score
		database.GetDatatbase(c).Find(&list)

		for i := range list {
			list[i].calculateTotal()
		}
		sort.Sort(ByScore(list))

		rank := 0
		prevScore := 0
		for i := range list {
			if prevScore != list[i].TotalScore {
				rank++
			}
			prevScore = list[i].TotalScore
			list[i].Rank = rank
		}

		c.HTML(http.StatusOK,	"score/list", gin.H{"DataList": list, "Actions": true})
	}).
		GET("/scores/:id", func(c *gin.Context) {
			id := c.Param("id")
			var score Score
			if "new" != id {
				dbm := database.GetDatatbase(c)
				i, _ := strconv.ParseInt(id, 10, 64)
				dbm.Find(&score, i)
			}
			c.HTML(http.StatusOK, "score/form", gin.H{"Score": score})

		}).POST("/scores", route.GetProcess("/",
		func(c *gin.Context, s session.Session, m *database.Manager) (int, interface{}) {

			if c.IsAborted() {
				panic("WTF")
			}
			if c.Writer.Written() {
				panic("WTF? wri")
			}

			var score Score
			err := c.Bind(&score)
			if nil != err {
				panic(err)
			}
			idString  := c.PostForm("ID")
			if c.Writer.Written() {
				panic("WTFfdsffffffffffff 01?")
			}

			i, err := strconv.ParseInt(idString, 10, 64)
			if nil != err {
				panic(err)
			}

			if i == 0 {
				m.Create(&score)
			} else {
				score.ID = i

				if c.Writer.Written() {
					panic("WTFfdsfff 02 fffffffff?")
				}
				var sc Score
				m.Find(&sc, i)
				score.CreatedAt = sc.CreatedAt
				if c.Writer.Written() {
					panic("WTFfdsfffffff 03 fffff?")
				}
				m.Save(&score)
				if c.Writer.Written() {
					panic("WTFfdsfffffffff 04 fff?")
				}
			}


			return http.StatusFound, "/scores"
		}))
	return rt
}
