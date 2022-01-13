# Overview & motivation

This project is the server for a chess-repertoire building app: Eski Chess. 

The app has been built to service intermediate / advanced chess players in their paths to becoming stronger players - specifically through the building and learning of opening repertoires. 

## Instructions
This should be done before the client is up and running. 

1. Install Go if you haven't: https://go.dev/dl/ 
2. Clone the project: `git clone https://github.com/Slayzur02/Eski-Server`
3. Create a PostgresQL db, and add a no-password user with all permissions.
4. Create a .env file with the following:

user=`username` (name of no-password user you created in previous step)

secretJwtKey=`dk01dk21239120dk` (random strong string for JWT encryption for the application)

hostEmail=`something@gmail.com` (your email which you want to use for the app to send email to users)

hostPassword=`cantbeguessed` (password for your above email)
1. Build the main.go file with: `go build cmd/main.go`
2. Run the final binary to start the server: `./main`

Finally, start the client. 

## Features

1. API-service for creating user profiles, including verification, setting changes, password resets. 
2. Websocket pools for validating board moves and board state for both analysis & games, PGN outputting
3. APIs for storing move sequences for books (repertoires) & games, querying of opening moves / lines, and line memorization with controlled randomization

## Tech stack

As with all go projects, the stack is simple: Go & PostgresQL & Redis. Uses Go-chi as the httprouter (the best one currently), sqlc as the database query generator (as I prefer thinking in SQL & speed, over thinking in ORMs and dealing with messy edge cases), goose for migrations, and numerous other small quality-of-life packages for security. Two communication protocols were used: Websockets & REST. 

