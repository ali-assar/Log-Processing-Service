set shell := ["powershell.exe", "-c"]


mock:
	go run .\mock-log-generator\cmd\main.go --url :8080 --interval-ms 500

processor:
    go run .\log-processor\cmd\main.go --urls "ws://localhost:8080/ws/logs,ws://localhost:9090/ws/logs"
