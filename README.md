# Todo List

Basic todo backend powered by gorilla/mux 🦍

## Production

1. Build a docker image

    ```bash
    docker build -t todolist .
    ```

2. Run the docker image you just built

    ```bash
    docker run -p 3000:80 todolist
    ```

## Development

1. Migrate the database

    ```bash
    mkdir -p sqlite
    goose -dir migrations sqlite3 sqlite/todolist.db up
    ```

2. Create a `.env.local` file with the following content

    ```bash
    DB_PATH=sqlite/todolist.db

    BE_HOST=0.0.0.0
    BE_PORT=3000
    ```

3. Run the app passing the `.env.local` file

    ```bash
    go run main.go -env .env.local
    ```
