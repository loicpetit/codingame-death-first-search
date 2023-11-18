# codingame-death-first-search

## Episode-1
https://www.codingame.com/ide/puzzle/death-first-search-episode-1

## Episode-2
https://www.codingame.com/ide/puzzle/death-first-search-episode-2

## Build unique file
Install the bundle tool (once)

go install golang.org/x/tools/cmd/bundle@latest

Bundle the main package

& ($env:USERPROFILE + "/go/bin/bundle") -o dist/main.go -dst ./main -prefix '""'  github.com/loicpetit/codingame-death-first-search/main

## Test
In project root "go test -v ./main"

## Benchmark
Debug, execute "go test -v -run nothing -benchtime 1000x -bench Debug ./main"

# Generate documention
In project root "go doc -cmd -u -all main > dist/main.txt"
