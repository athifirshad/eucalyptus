# Internal Package

The `internal` package contains the core data handling and business logic of the application. It is not intended to be used by external packages and is kept private to ensure the integrity of the application's internal state and operations. Queries used here are manually created and not code generated by sqlc as in 'db'.

## Structure

The `internal` package is organized into several files, each with a specific purpose:

- `users.go`: Contains the `UserModel` struct and related functions for managing user data. This includes creating new users, retrieving users by email, and handling password hashing and verification.

- `models.go`: Provides a `Models` struct that encapsulates the `UserModel` and can be used to perform operations on user data. It is initialized with a database connection pool and provides a clean interface for interacting with the database.

## Key Components

- `UserModel`: A struct that represents the user model in the application. It includes fields for user information and methods for interacting with the database, such as creating new users and retrieving users by email.

- `User`: A struct that represents a user in the system. It includes fields for user details and a `password` field that is a custom type for handling password hashing and verification.

- `password`: A custom type that encapsulates the plaintext password and its hashed version. It provides methods for setting and verifying passwords.

