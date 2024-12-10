# Real-Time Forum

Welcome to the Real-Time Forum project! This application is an enhanced version of a previous forum, featuring real-time interactions, private messaging, and a single-page application design. The forum is built using a combination of Go, JavaScript, HTML, CSS, and SQLite.

## Table of Contents

- [Objectives](#objectives)
- [Features](#features)
  - [Registration and Login](#registration-and-login)
  - [Posts and Comments](#posts-and-comments)
  - [Private Messages](#private-messages)
- [Technologies Used](#technologies-used)
- [Setup and Installation](#setup-and-installation)
- [Usage](#usage)


## Objectives

The main objectives of this project are to:

- Implement a real-time forum with enhanced features.
- Utilize WebSockets for real-time communication.
- Create a single-page application using JavaScript for dynamic content updates.

## Features

### Registration and Login

- Users can register by providing:
  - Nickname
  - Age
  - Gender
  - First Name
  - Last Name
  - E-mail
  - Password
- Users can log in using either their nickname or e-mail combined with their password.
- Users can log out from any page within the forum.

### Posts and Comments

- Users can create posts with categories.
- Users can comment on posts.
- Posts are displayed in a feed, with comments visible upon clicking a post.

### Private Messages

- Users can send private messages to each other.
- The chat interface includes:
  - An online/offline user list, organized by last message sent or alphabetically for new users.
  - A chat history that loads the last 10 messages, with more messages loaded as the user scrolls up.
- Messages include:
  - A date indicating when the message was sent.
  - The username of the sender.
- Real-time message notifications using WebSockets.

## Technologies Used

- **Backend**: Go (Golang) with Gorilla WebSocket
- **Frontend**: JavaScript, HTML, CSS
- **Database**: SQLite
- **Security**: bcrypt for password hashing
- **Unique Identifiers**: UUID for user and session management

## Setup and Installation

1. **Clone the repository**:

2. **Install dependencies**:
   - Ensure you have Go installed. Follow the [official guide](https://golang.org/doc/install) if needed.
   - Install SQLite3 and ensure it's accessible from your command line.

3. **Run the application after navigating to the project directory**:
   ```bash
   go run .
   ```

4. **Access the application**:
   - Open your web browser and navigate to `http://localhost:8081`.

## Usage

- **Register**: Create a new account to access the forum.
- **Login**: Use your credentials to log in and start interacting.
- **Create Posts**: Share your thoughts and engage with the community.
- **Send Messages**: Communicate privately with other users in real-time.
