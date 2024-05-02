build: 
	@ echo ----------------- starting API build ----------------- 
	@ go build -o ./dist/app ./cmd/main/main.go 
	@ echo ----------------- echo build finished -----------------

start: build 
	@ ./dist/app 
