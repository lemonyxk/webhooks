/**
* @program: webhooks
*
* @description:
*
* @author: lemo
*
* @create: 2022-03-09 17:45
**/

package main

// https://developer.github.com/webhooks/
import (
	"os"
	"time"

	"github.com/lemoyxk/console"
	"github.com/lemoyxk/kitty"
	http2 "github.com/lemoyxk/kitty/http"
	"github.com/lemoyxk/kitty/http/server"
	"github.com/lemoyxk/utils"
)

func main() {

	var httpServer = kitty.NewHttpServer("0.0.0.0:8667")

	httpServer.ReadTimeout = 10 * time.Second
	httpServer.WriteTimeout = 10 * time.Second

	var httpServerRouter = kitty.NewHttpServerRouter()

	httpServer.Use(func(next server.Middle) server.Middle {
		return func(stream *http2.Stream) {
			next(stream)
		}
	})

	httpServer.OnSuccess = func() {
		console.Info(httpServer.LocalAddr())
	}

	httpServerRouter.Group().Before(GithubBefore).Handler(func(handler *server.RouteHandler) {
		handler.Post("/github").Handler(func(stream *http2.Stream) error {
			var github = stream.Context.Value("github").(*GitHub)
			var event = stream.Request.Header.Get("X-GitHub-Event")

			console.Info("Repository:", github.Repository.FullName, "Event:", event)

			var repo Repo
			for i := 0; i < len(Config.Repositories); i++ {
				if Config.Repositories[i].FullName == github.Repository.FullName {
					repo = Config.Repositories[i]
					break
				}
			}

			console.Info("Start Command:", repo.Script.Start, "Dir:", repo.Script.Dir)

			var cmd = newCmd(repo.Script.Start)

			cmd.Dir = repo.Script.Dir
			cmd.Stderr = os.Stderr
			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout

			err := cmd.Start()
			if err != nil {
				panic(err)
			}

			_, err = cmd.Process.Wait()
			if err != nil {
				console.Error(err)
			}

			console.Info("End Command:", repo.Script.Start, "Dir:", repo.Script.Dir)

			return stream.JsonFormat("SUCCESS", 200, nil)
		})
	})

	go httpServer.SetRouter(httpServerRouter).Start()

	utils.Signal.ListenKill().Done(func(sig os.Signal) {
		console.Info("server stop")
	})
}
