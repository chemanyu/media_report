@echo off
echo Starting Media Report API Service...
cd service\api
go run media.go -f etc\media-api.yaml