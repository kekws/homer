@echo off
echo make sure you've changed the IP and PORT in both homer/main.go and homer/client/client.go
echo -----------------------------------------------------------------------------------------
echo building homer c2 server and client...

go build -o build/server.exe .
go build -o build/client.exe client/client.go