Миграции  
`brew install goose`

Скачать либу  
`go get -u ...`

Запустить файл  
`go run ...`

Накатить миграшку  
`goose -dir ./migrations postgres "postgres://postgres:postgres@localhost:5436/postgres?sslmode=disable" up`  
Откатить миграшку  
`goose -dir ./migrations postgres "postgres://postgres:postgres@localhost:5436/postgres?sslmode=disable" down`
