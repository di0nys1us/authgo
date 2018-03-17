package main

type rootMutation struct {
	repository repository
}

type identity struct {
	ID      string
	Version int
}

// CreateUser

func (m *rootMutation) CreateUser(args struct {
	Input userInput
}) (*userOutput, error) {
	user := &user{
		FirstName: args.Input.FirstName,
		LastName:  args.Input.LastName,
		Email:     args.Input.Email,
		Password:  args.Input.Password,
		Enabled:   args.Input.Enabled,
		Deleted:   args.Input.Deleted,
	}

	err := m.repository.saveUser(user)

	if err != nil {
		return nil, err
	}

	return &userOutput{&userResolver{m.repository, user}}, nil
}

type userInput struct {
	FirstName string
	LastName  string
	Email     string
	Password  string
	Enabled   bool
	Deleted   bool
}

type userOutput struct {
	user *userResolver
}

func (o *userOutput) User() *userResolver {
	return o.user
}

// UpdateUser

func (m *rootMutation) UpdateUser(args struct {
	Identity identity
	Input    userInput
}) (*userOutput, error) {
	return nil, nil
}

// CreateRole

func (m *rootMutation) CreateRole(args struct {
	Input roleInput
}) (*roleOutput, error) {
	return nil, nil
}

type roleInput struct {
	Name string
}

type roleOutput struct {
	role *roleResolver
}

func (o *roleOutput) Role() *roleResolver {
	return o.role
}

// UpdateRole

func (m *rootMutation) UpdateRole(args struct {
	Identity identity
	Input    roleInput
}) (*roleOutput, error) {
	return nil, nil
}

// CreateAuthority

func (m *rootMutation) CreateAuthority(args struct {
	Input authorityInput
}) (*authorityOutput, error) {
	return nil, nil
}

type authorityInput struct {
	Name string
}

type authorityOutput struct {
	authority *authorityResolver
}

func (o *authorityOutput) Authority() *authorityResolver {
	return o.authority
}

// UpdateAuthority

func (m *rootMutation) UpdateAuthority(args struct {
	Identity identity
	Input    authorityInput
}) (*authorityOutput, error) {
	return nil, nil
}
