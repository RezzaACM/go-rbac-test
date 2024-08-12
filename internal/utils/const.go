package utils

// CONSTANT STATUS RESPONSE
const (
	StatusSuccessfully           = "Action is Successfully"
	StatusUnauthorized           = "Email or Password is invalid"
	StatusParameterIsRequired    = "Parameter is required can not be null"
	StatusLoginIsSuccessfully    = "Login is Successfully"
	StatusEmailAlreadyUsed       = "Email is already used"
	StatusUsernameAlreadyUsed    = "Username is already used"
	StatusLoginHeaderRequired    = "Authorization header is required"
	StatusLoginInvalidToken      = "Login failed. Token has invalid"
	StatusRoleAlreadyUsed        = "Role with name :name already used"
	StatusDataNotFound           = ":name not found"
	StatusDeletedDataSucessfully = ":name has been successfully deleted."
	StatusMismatchRole           = "You are not allowed to perform this action"
	StatusForbidden              = "You don't have permission to perform this action"
	StatusSomethingWrong         = "Something went wrong"
)

// CONSTANT ROLE NAME
const (
	Admin   = "admin"
	User    = "user"
	Manager = "manager"
	Guest   = "guest"
)

// CONSTANT ACTION LOG
const (
	CreateGroup          = "CREATE_GROUP"
	UpdateGroup          = "UPDATE_GROUP"
	DeleteGroup          = "DELETE_GROUP"
	CreateUser           = "CREATE_USER"
	UpdateUser           = "UPDATE_USER"
	DeleteUser           = "DELETE_USER"
	CreateRole           = "CREATE_ROLE"
	UpdateRole           = "UPDATE_ROLE"
	DeleteRole           = "DELETE_ROLE"
	CreateGroupRole      = "CREATE_GROUP_ROLE"
	UpdateGroupRole      = "UPDATE_GROUP_ROLE"
	DeleteGroupRole      = "DELETE_GROUP_ROLE"
	CreatePermission     = "CREATE_PERMISSION"
	UpdatePermission     = "UPDATE_PERMISSION"
	DeletePermission     = "DELETE_PERMISSION"
	CreateUserRole       = "CREATE_USER_ROLE"
	UpdateUserRole       = "UPDATE_USER_ROLE"
	DeleteUserRole       = "DELETE_USER_ROLE"
	CreateRolePermission = "CREATE_ROLE_PERMISSION"
	UpdateRolePermission = "UPDATE_ROLE_PERMISSION"
	DeleteRolePermission = "DELETE_ROLE_PERMISSION"
	CreateUserGroup      = "CREATE_USER_GROUP"
	UpdateUserGroup      = "UPDATE_USER_GROUP"
	DeleteUserGroup      = "DELETE_USER_GROUP"
	CreateActionLog      = "CREATE_ACTION_LOG"
	Login                = "LOGIN"
)
