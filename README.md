# Todo Watcher

## Getting started

Setup .env file in the root directory with the following content:

```
HOST (e.g. localhost)
PORT (e.g. 8080)
TODO_FOLDER (e.g. /path/to/todo/folder)
```

Run backend:

```
go run main.go
```

Setup the frontend/.env file with the following content:

```
VITE_API_URL (e.g. http://localhost:8080)
```

Run frontend:

```
cd frontend
npm install
npm run dev
```

The frontend will be available at `http://localhost:5173`. 

Note: The frontend is a Vue.js application (bootstrapped with `npm create vue@latest` ), and you need to have Node.js v18.3+ and npm installed to run it.
