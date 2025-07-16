# Todo Watcher

This is a simple todo watcher application that allows you to manage a folder filled with plain text todo files in a simple format. It provides a backend written in Go and a frontend built with Vue.js. Todos should be edited manually in a text editor of your choice, and the application will automatically update the view when changes are detected.

## Todo Format

The backend parses the following todo format:
- Status:
    - default: not started
    - `x` for done
    - `~` for in progress
    - `#` for dropped
- Priority:
    - default: medium priority
    - `!` for high priority
    - `.` for low priority
- Tentative todos can be marked with a question mark at the end of the line.
- Todos can have arbitrarily nested subtasks, which are either grouped by indentation or grouped as a section with at least 3 `-` or `=` characters at the beginning of the line (e.g. `--- Subtask` or `=== Subtask`).
    - Note: unlabelled section dividers (e.g. `---`) are not considered as task groups, but rather as visual separators and are ignored.
- Sections titled `done` or `dropped` are given the respective status.
- Special todo statuses and priorities are inherited from the parent task.

## Getting started

Setup .env file in the backend directory with the following content:

```
HOST (e.g. localhost)
PORT (e.g. 8080)
TODO_FOLDER (e.g. /path/to/todo/folder)
EDITOR (e.g. code for Visual Studio Code)
```

Run backend (written in go 1.19):

```
cd backend
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

## Backend API

The backend provides a simple REST API to interact with the todo files. The main endpoint is `/api/todos`, which returns a JSON representation of the todos in the specified folder. You can also open files directly using the `/api/open` endpoint, which will open the specified file in the editor defined in your `.env` file.
