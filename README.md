# ProiectColectivBackEnd

## Install Go

### macOS
```bash
brew install go
```

### Linux/Windows
Download from [golang.org/dl](https://golang.org/dl/)

## Setup

1. Install dependencies:
```bash
  go mod tidy
```

2. Create `.env` file with Firebase configuration:
```
FIREBASE_DATABASE_URL=https://your-project-id-default-rtdb.region.firebasedatabase.app/
FIREBASE_CREDENTIALS_PATH=secret/your-firebase-adminsdk-key.json
```

3. Add your Firebase service account key to `secret/` directory

## Run Server

```bash
  go run main.go
```

Server runs on `http://localhost:8080`

## API Endpoints

- `POST /users/signup` - Create user
- `GET /users/:id` - Get user by ID
- `GET /users` - Get all users
- `PUT /users/:id` - Update user
- `DELETE /users/:id` - Delete user

- `POST/teams` - Create a team  (+ Json example: {"name": "nameTest", "description": "descTest", "ispublic": true})
- `POST/teams/addUserToTeam` - Add a user to a team (+Json example: {"userId":"id1", "teamId":"id2"})
- `GET/teams/:id` - Get team by ID
- `GET/teams` - Get all teams
- `GET/teams/search?prefix= &limit= ` - Get the first "limit" teams whose names start with "prefix"
- `GET/teams/by-name?name=` - Get team(s) by name
- `PUT/teams/:id` - Update team
- `DELETE/teams/:id`  - Delete team
