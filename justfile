set shell := ["powershell.exe", "-c"]


mock:
	go run .\mock-log-generator\cmd\main.go --url :9090 --interval-ms 500