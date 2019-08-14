package bus

import (
	"context"
	"log"
	"sync"

	// "log"
	"os"
	"os/signal"
	"time"

	"github.com/focusteach/gfworks/app"
)

var (
	singleton *Application
	once      sync.Once
)

func init() {

}

// Application Application
type Application struct {
	tasks []app.IAppTask
	quit  chan os.Signal
}

//GetInstance 用于获取单例模式对象
func GetInstance() *Application {
	once.Do(func() {
		singleton = &Application{}
	})

	return singleton
}

// AddTask add application level task
func (app *Application) AddTask(task app.IAppTask) {
	app.tasks = append(app.tasks, task)

	go func() {
		task.Start()
	}()
}

// Exec exec
func (app *Application) Exec() {
	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	app.quit = make(chan os.Signal)
	signal.Notify(app.quit, os.Interrupt)
	<-app.quit
	log.Println("Shutdown Servers ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for _, task := range app.tasks {
		if err := task.Shutdown(ctx); err != nil {
			log.Fatal("Server Shutdown:", err)
		}
	}

	log.Println("Server exiting")
}

// Quit force quit
func (app *Application) Quit() {
	app.quit <- os.Kill
}
