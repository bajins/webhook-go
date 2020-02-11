package main

import (
	"flag"
	"log"
	"os"
	"webhook-go/utils"
)

// 获取地址，如果没传默认值为0.0.0.0
func Host() (host string) {
	flag.StringVar(&host, "h", "0.0.0.0", "默认地址:0.0.0.0")
	flag.Parse()
	return host
}

// 获取传入参数的端口，如果没传默认值为8000
func Port() (port string) {
	flag.StringVar(&port, "p", "8000", "默认端口:8000")
	flag.Parse()
	return port

	//if len(os.Args[1:]) == 0 {
	//	return ":8000"
	//}
	//return ":" + os.Args[1]
}

func main() {
	if err := LoadConfig(); err != nil {
		utils.Log2file(err.Error(), "")
		os.Exit(1)
	}

	listenErr := StartService(Host(), Port())
	if listenErr != nil {
		log.Fatal("ListenAndServer error: ", listenErr)
	}

}
