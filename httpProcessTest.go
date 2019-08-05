package main

import (
	"fmt"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"time"

	"github.com/aerogo/aero"
	"github.com/aofei/air"
	"github.com/astaxie/beego"
	co "github.com/astaxie/beego/context"
	"github.com/labstack/echo"
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
	return c.String(http.StatusOK, "Hello, echo")
}
func startEcho() {
	e := echo.New()
	e.GET("/hello", echoHandler)
	e.Logger.Fatal(e.Start(":" + strconv.Itoa(port)))
}

// aero
func aeroHandler(ctx aero.Context) error {
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
	return res.WriteString("Hello, 世界")
}
func startAir() {
	a := air.New()
	a.Address = ":" + strconv.Itoa(port)
	a.GET("/hello", airHandler)
	a.Serve()
}
