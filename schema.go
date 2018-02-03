package authgo

const (
	schema = `
		schema {
			query: Query
		}

		type Query {
			users: [User!]!
		}

		type User {
			id: ID!
			version: Int!
			firstName: String!
			lastName: String!
			email: String!
			password: String!
			enabled: Boolean!
			deleted: Boolean!
		}
	`
)
