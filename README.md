# Todo List

Basic todo backend powered by gorilla/mux 🦍

## Setup

1. Migrate sqlite3 database

    ```bash
    goose -dir migrations sqlite3 sqlite/todolist.db up
    ```

2. Run the application in production mode

    ```bash
    go build
    ./todolist
    ```
