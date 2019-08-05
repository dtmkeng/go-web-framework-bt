package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"time"

	"github.com/aerogo/aero"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/aofei/air"
	"github.com/astaxie/beego"
	co "github.com/astaxie/beego/context"
	"github.com/gin-gonic/gin"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo/render"
	"github.com/gobuffalo/envy"
	contenttype "github.com/gobuffalo/mw-contenttype"
	paramlogger "github.com/gobuffalo/mw-paramlogger"
	"github.com/gobuffalo/x/sessions"
	"github.com/gorilla/mux"
	"github.com/labstack/echo"
	"github.com/rs/cors"
)

var port = 8081
var sleepTime = 0
var cpuBound bool
var target = 15
var sleepTimeDuration time.Duration
var message = []byte("hello world")
var messageStr = "hello world"
var samplingPoint = 20 //seconds
func main() {

	args := os.Args
	argsLen := len(args)
	webFramework := "default"
	if argsLen > 1 {
		webFramework = args[1]
	}
	if argsLen > 2 {
		sleepTime, _ = strconv.Atoi(args[2])
		if sleepTime == -1 {
			cpuBound = true
			sleepTime = 0
		}
	}
	if argsLen > 3 {
		port, _ = strconv.Atoi(args[3])
	}
	if argsLen > 4 {
		samplingPoint, _ = strconv.Atoi(args[4])
	}
	sleepTimeDuration = time.Duration(sleepTime) * time.Millisecond
	samplingPointDuration := time.Duration(samplingPoint) * time.Second

	go func() {
		time.Sleep(samplingPointDuration)
		var mem runtime.MemStats
		runtime.ReadMemStats(&mem)
		var u uint64 = 1024 * 1024
		fmt.Printf("TotalAlloc: %d\n", mem.TotalAlloc/u) // ไบต์สูงสุดสะสมที่จัดสรรบนฮีป (จะไม่ลดลง)
		fmt.Printf("Alloc: %d\n", mem.Alloc/u)           // จำนวนไบต์ที่จัดสรรในปัจจุบันบนฮีป
		fmt.Printf("HeapAlloc: %d\n", mem.HeapAlloc/u)   // จำนวนไบต์ที่จัดสรรในปัจจุบันบนฮีป
		fmt.Printf("HeapSys: %d\n", mem.HeapSys/u)       // หน่วยความจำทั้งหมดที่ได้รับจากระบบปฏิบัติการ
	}()
	switch webFramework {
	case "default": // default
		startDefaultMux()
	case "beego": // beego
		startBeego()
	case "echo": // echo
		startEcho()
	case "aero": // aero
		startAero()
	case "air":
		startAir()
	case "gin":
		startGin()
	case "gorilamux":
		startMux()
	case "go-rest":
		startGoRest()
	case "buffalo":
		startBuffalo()
	}

}
func startDefaultMux() {
	http.HandleFunc("/hello", helloHandler)
	http.ListenAndServe(":"+strconv.Itoa(port), nil)
}
func helloHandler(w http.ResponseWriter, r *http.Request) {
	if cpuBound {
		pow(target)
	} else {
		if sleepTime > 0 {
			time.Sleep(sleepTimeDuration)
		} else {
			runtime.Gosched()
		}
	}
	w.Write(message)
}

//beego
func beegoHandler(ctx *co.Context) {
	if cpuBound {
		pow(target)
	} else {

		if sleepTime > 0 {
			time.Sleep(sleepTimeDuration)
		} else {
			runtime.Gosched()
		}
	}
	ctx.WriteString(messageStr)

}
func startBeego() {
	beego.BConfig.RunMode = beego.PROD
	// beego.BeeLogger.Close()
	mux := beego.NewControllerRegister()
	mux.Get("/hello", beegoHandler)
	http.ListenAndServe(":"+strconv.Itoa(port), mux)
}

// echo
func echoHandler(c echo.Context) error {
	if cpuBound {
		pow(target)
	} else {

		if sleepTime > 0 {
			time.Sleep(sleepTimeDuration)
		} else {
			runtime.Gosched()
		}
	}
	return c.String(http.StatusOK, "Hello, echo")
}
func startEcho() {
	e := echo.New()
	e.GET("/hello", echoHandler)
	e.Logger.Fatal(e.Start(":" + strconv.Itoa(port)))
}

// aero
func aeroHandler(ctx aero.Context) error {
	if cpuBound {
		pow(target)
	} else {

		if sleepTime > 0 {
			time.Sleep(sleepTimeDuration)
		} else {
			runtime.Gosched()
		}
	}
	return ctx.String("Hello World")
}
func startAero() {
	app := aero.New()
	app.Get("/hello", aeroHandler)
	app.Config.Ports.HTTP = 8081
	app.Run()
}

// air
func airHandler(req *air.Request, res *air.Response) error {
	if cpuBound {
		pow(target)
	} else {

		if sleepTime > 0 {
			time.Sleep(sleepTimeDuration)
		} else {
			runtime.Gosched()
		}
	}
	return res.WriteString("Hello, 世界")
}
func startAir() {
	a := air.New()
	a.Address = ":" + strconv.Itoa(port)
	a.GET("/hello", airHandler)
	a.Serve()
}

// Gin
func ginHandler(c *gin.Context) {
	if cpuBound {
		pow(target)
	} else {

		if sleepTime > 0 {
			time.Sleep(sleepTimeDuration)
		} else {
			runtime.Gosched()
		}
	}
	c.JSON(200, gin.H{
		"message": "Hello, 世界",
	})
}
func startGin() {
	gin.DisableConsoleColor()
	r := gin.Default()
	r.GET("/hello", ginHandler)
	r.Run(":" + strconv.Itoa(port))
}

// mux
func muxHandler(w http.ResponseWriter, r *http.Request) {
	if cpuBound {
		pow(target)
	} else {

		if sleepTime > 0 {
			time.Sleep(sleepTimeDuration)
		} else {
			runtime.Gosched()
		}
	}
	fmt.Fprintf(w, "Hello World")
}
func startMux() {
	r := mux.NewRouter()
	r.HandleFunc("/hello", muxHandler)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), r))
}

// go-rest-api
func goRestHandler(w rest.ResponseWriter, req *rest.Request) {
	if cpuBound {
		pow(target)
	} else {

		if sleepTime > 0 {
			time.Sleep(sleepTimeDuration)
		} else {
			runtime.Gosched()
		}
	}
	w.WriteJson(map[string]string{"Body": "Hello World!"})
}
func startGoRest() {
	api := rest.NewApi()
	api.Use(rest.DefaultDevStack...)
	router, err := rest.MakeRouter(
		rest.Get("/hello", goRestHandler),
	)
	if err != nil {
		log.Fatal(err)
	}
	api.SetApp(router)
	http.ListenAndServe(":"+strconv.Itoa(port), api.MakeHandler())
}

// buffalo
var r *render.Engine

// ENV ...
var ENV = envy.Get("GO_ENV", "development")

func buffaloHandler(c buffalo.Context) error {
	if cpuBound {
		pow(target)
	} else {

		if sleepTime > 0 {
			time.Sleep(sleepTimeDuration)
		} else {
			runtime.Gosched()
		}
	}
	return c.Render(200, r.JSON(map[string]string{"message": "Welcome to Buffalo!"}))
}
func startBuffalo() {
	app := buffalo.New(buffalo.Options{
		Env:          ENV,
		SessionStore: sessions.Null{},
		PreWares: []buffalo.PreWare{
			cors.Default().Handler,
		},
		SessionName: "_coke_session",
		Addr:        ":" + strconv.Itoa(port),
	})
	app.Use(paramlogger.ParameterLogger)
	app.Use(contenttype.Set("application/json"))

	app.GET("/hello", buffaloHandler)
	// app.Options.Addr := 8080
	app.Serve()
}
