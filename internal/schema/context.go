package schema

type reqContextKey string

// DbContextKey is the GraphQL request context key for the database connection.
const DbContextKey = reqContextKey("db")

// UserKey is the GraphQL request context key for the user.
const UserKey = reqContextKey("user")
