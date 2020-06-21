package main

import (
	"encoding/json"
	"flag"
	"github.com/jinzhu/gorm"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/websocket"
	//"fmt"
	"grizzled/database"
	"grizzled/dotenv"
	. "grizzled/game"
	"log"
	"os"
	"strconv"
)

var (
	initdb bool
	err    error
)

func init() {
	flag.BoolVar(&initdb, "initdb", false, "Init the database")
}

type Request struct {
	Order  string `json:"order"`
	Choose string `json:"choose"`
}

func main() {
	flag.Parse()
	{ //基本載入設定
		if err := dotenv.Config(); err != nil {
			return
		} //設定檔載入

		database.DB, err = gorm.Open("sqlite3", os.Getenv("dbroute"))
		if err != nil {
			log.Fatal(err)
			return
		}
		defer database.DB.Close()

		if initdb == true {
			database.InitCard()
		}
	}

	ws := websocket.New(websocket.DefaultGorillaUpgrader, websocket.Events{
		websocket.OnNativeMessage: func(nsConn *websocket.NSConn, msg websocket.Message) error {
			log.Printf("Server got: %s from [%s]", msg.Body, nsConn.Conn.ID())
			
			var request Request
			json.Unmarshal(msg.Body, &request)

			receive := request.Order
			switch receive {
			case "draw":
				Game.Draw(nsConn.Conn.ID(),6)
				msg.Body = Players[GetOrder(nsConn.Conn.ID())].Render()
				nsConn.Conn.Write(msg)

			case "restart":
				Game.InitGame()

				msg.Body = Players[GetOrder(nsConn.Conn.ID())].Render()
				nsConn.Conn.Write(msg)
				msg.Body = Game.Render()
				nsConn.Conn.Write(msg)
				nsConn.Conn.Server().Broadcast(nsConn, msg)


			case "heroUse":
				if Players[0].PlayHero() {
					msg.Body = Players[GetOrder(nsConn.Conn.ID())].Render()
					nsConn.Conn.Write(msg)
				}

			case "speech":
				if Players[GetOrder(nsConn.Conn.ID())].Speech() {
					msg.Body = Players[GetOrder(nsConn.Conn.ID())].Render()
					nsConn.Conn.Write(msg)
					msg.Body = Game.Render()
					nsConn.Conn.Write(msg)
				}

			case "luckyClover":
				choose, err := strconv.Atoi(request.Choose)
				if err != nil {
					log.Println("Something wrong!")
				}
				if Players[GetOrder(nsConn.Conn.ID())].HeroPower(choose) {
					msg.Body = Players[GetOrder(nsConn.Conn.ID())].Render()
					nsConn.Conn.Write(msg)
					msg.Body = Game.Render()
					nsConn.Conn.Write(msg)
				}

			case "playCard":
				num, _ := strconv.Atoi(request.Choose)
				Players[GetOrder(nsConn.Conn.ID())].PlayCard(num)
				Game.GameNext()

				msg.Body = Players[GetOrder(nsConn.Conn.ID())].Render()
				nsConn.Conn.Write(msg)
				msg.Body = Game.Render()
				nsConn.Conn.Write(msg)

				nsConn.Conn.Server().Broadcast(nsConn, msg)

			case "login":
				Game.NewPlayer(nsConn.Conn.ID(), request.Choose)
				Game.NoMansLand = []database.Card{}
				Game.Stage = "Waiting"
				msg.Body = Game.Render()
				nsConn.Conn.Write(msg)
				nsConn.Conn.Server().Broadcast(nsConn, msg)

			default:

			}

			return nil
		},
	})

	ws.OnConnect = func(c *websocket.Conn) error {
		log.Printf("[%s] Connected to server!", c.ID())
		//c.Write(game.Board.NoMansLand)
		return nil
	}

	ws.OnDisconnect = func(c *websocket.Conn) {
		log.Printf("[%s] Disconnected from server", c.ID())
	}

	app := iris.New()
	app.RegisterView(iris.HTML("./v", ".html"))
	app.HandleDir("/js", "./v/js") // serve our custom javascript code.
	app.HandleDir("/", "./v")
	app.Get("/my_endpoint", websocket.Handler(ws))

	//主頁
	app.Get("/", func(ctx iris.Context) {
		//綁定數據
		ctx.ViewData("Host", "localhost:"+os.Getenv("port"))
		// 渲染視圖文件: ./v/index.html
		ctx.View("index.html")

	})

	app.Run(iris.Addr(":" + os.Getenv("port")))
}
