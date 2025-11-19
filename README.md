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
- `DELETE/teams/deleteUserFromTeam` - Delete a user from a team (+Json example: {"userId":"id1", "teamId":"id2"})
- `GET/teams/:id` - Get team by ID
- `GET/teams` - Get all teams
- `GET/teams/search?prefix= &limit= ` - Get the first "limit" teams whose names start with "prefix"
- `GET/teams/by-name?name=` - Get team(s) by name
- `PUT/teams/:id` - Update team
- `DELETE/teams/:id`  - Delete team

- `POST /quizzes` - Create a quiz (protected - requires Bearer token)
  + JSON example:
  {
    "quiz_name": "Sample Quiz",
    "user_id": "123",
    "team_id": "team123",
    "questions": [
      {"type": "multiple_choice", "question": "What is 2+2?", "options": ["1", "2", "4"], "answers": ["4"]}
    ]
  }
- `GET /quizzes/:id` - Get a quiz with answers (protected - requires Bearer token)
- `GET /quizzes/:id/test` - Get a quiz without answers for taking the test (protected - requires Bearer token)
- `POST /quizzes/:id/test` - Submit quiz answers and get results (protected - requires Bearer token)
  + JSON example:
  {
    "quiz_id": "quiz123",
    "attempts": [
      {"quiz_question_id": "q1", "answer": ["4"]}
    ]
  }

## WebSockets

### Real-time messaging

`GET /messages/connect?token=<JWT>`: Connect to real-time messaging

The WebSocket then sends messages of type:

```json
{
  "type": "string",
  "payload": {
    "senderId": "string",
    "receiverId": "string",
    "textContent": "string"
  }
}
```

- **type**: `"direct_message"` or `"team_message"`
- **senderId**: the sending `userID`
- **receiverId**: either a `userID` or `teamID` depending on `type`
- **textContent**: the actual message content

## Swagger Support

Swagger UI runs on `http://localhost:8080/swagger/index.html`

### How to use

- Annotate the controller functions with comments as seen in the documentation [here](https://github.com/swaggo/swag/blob/master/README.md#declarative-comments-format).
- Run `swag fmt -g main.go` in the project root to format the annotation comments.
- Generate the Swagger files using `swag init -g main.go` in the project root. The generated files are located in `docs/`.
