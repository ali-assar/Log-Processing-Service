set shell := ["powershell.exe", "-c"]


mock:
	go run .\mock-log-generator\cmd\main.go --url :8082 --interval-ms 500

processor:
    go run .\log-processor\cmd\main.go --urls "ws://localhost:8080/ws/logs,ws://localhost:8081/ws/logs,ws://localhost:8082/ws/logs" --http-addr ":9090"
