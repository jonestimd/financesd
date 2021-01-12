package graphql

type reqContextKey string

// DbContextKey is the GraphQL request context key for the database connection.
const DbContextKey = reqContextKey("db")
