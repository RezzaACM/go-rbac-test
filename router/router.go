package router

import (
	"example.com/go-crud/internal/config"
	"example.com/go-crud/internal/handlers"
	"example.com/go-crud/internal/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRouter(cfg *config.Config) *gin.Engine {
	router := gin.Default()

	// Handler
	userHandler := handlers.NewUserHandler(cfg)
	roleHandler := handlers.NewRoleHandler(cfg)
	permissionHandler := handlers.NewPermissionHandler(cfg)
	userRoleHandler := handlers.NewUserRoleHandler(cfg)
	rolePermissionHandler := handlers.NewRolePermissionHandler(cfg)
	groupHandler := handlers.NewGroupHandler(cfg)
	userGroupHandler := handlers.NewUserGroupHandler(cfg)
	groupRoleHandler := handlers.NewGroupRoleHandler(cfg)
	actionLogHandler := handlers.NewActionLogHandler(cfg)

	// Middleware
	router.Use(middleware.CORSMiddleware())

	// Routes

	userRouter := router.Group("/api/v1/users")
	{
		userRouter.GET("", middleware.AuthMiddleware(cfg), userHandler.GetUsers)
		userRouter.POST("", userHandler.CreateUser)
		userRouter.POST("/login", userHandler.LoginUser)
	}

	roleRouter := router.Group("/api/v1/roles")
	{
		roleRouter.GET("", middleware.AuthMiddleware(cfg), roleHandler.GetRoles)
		roleRouter.GET("/:id", middleware.AuthMiddleware(cfg), roleHandler.GetRole)
		roleRouter.POST("", middleware.AuthMiddleware(cfg), middleware.RBACMiddleware("admin", "manager", "user"), roleHandler.CreateRole)
		roleRouter.PUT("/:id", middleware.AuthMiddleware(cfg), middleware.RBACMiddleware("admin", "manager", "user"), roleHandler.UpdateRole)
		roleRouter.DELETE("/:id", middleware.AuthMiddleware(cfg), middleware.RBACMiddleware("admin", "manager", "user"), roleHandler.DeleteRole)
	}

	permissionRouter := router.Group("/api/v1/permissions")
	{
		permissionRouter.GET("", middleware.AuthMiddleware(cfg), permissionHandler.GetPermissions)
		permissionRouter.GET("/:id", middleware.AuthMiddleware(cfg), permissionHandler.GetPermission)
		permissionRouter.POST("", middleware.AuthMiddleware(cfg), middleware.RBACMiddleware("admin", "manager", "user"), permissionHandler.CreatePermission)
		permissionRouter.PUT("/:id", middleware.AuthMiddleware(cfg), middleware.RBACMiddleware("admin", "manager", "user"), permissionHandler.UpdatePermission)
		permissionRouter.DELETE("/:id", middleware.AuthMiddleware(cfg), middleware.RBACMiddleware("admin", "manager", "user"), permissionHandler.DeletePermission)
	}

	userRoleRouter := router.Group("/api/v1/user-roles")
	{
		userRoleRouter.GET("", middleware.AuthMiddleware(cfg), userRoleHandler.GetUserRoles)
		userRoleRouter.POST("", middleware.AuthMiddleware(cfg), middleware.RBACMiddleware("admin", "manager", "user"), userRoleHandler.CreateUserRoles)
		userRoleRouter.PUT("/:id", middleware.AuthMiddleware(cfg), middleware.RBACMiddleware("admin", "manager", "user"), userRoleHandler.UpdateUserRoles)
		userRoleRouter.DELETE("/:id", middleware.AuthMiddleware(cfg), middleware.RBACMiddleware("admin", "manager", "user"), userRoleHandler.DeleteUserRoles)
		userRoleRouter.GET("/:id", middleware.AuthMiddleware(cfg), userRoleHandler.GetUserRole)
	}

	rolePermissionRouter := router.Group("/api/v1/role-permissions")
	{
		rolePermissionRouter.GET("", middleware.AuthMiddleware(cfg), rolePermissionHandler.GetRolePermissions)
		rolePermissionRouter.POST("", middleware.AuthMiddleware(cfg), middleware.RBACMiddleware("admin", "manager", "user"), rolePermissionHandler.CreateRolePermission)
		rolePermissionRouter.PUT("/:id", middleware.AuthMiddleware(cfg), middleware.RBACMiddleware("admin", "manager", "user"), rolePermissionHandler.UpdateRolePermission)
		rolePermissionRouter.DELETE("/:id", middleware.AuthMiddleware(cfg), middleware.RBACMiddleware("admin", "manager", "user"), rolePermissionHandler.DeleteRolePermission)
		rolePermissionRouter.GET("/:id", middleware.AuthMiddleware(cfg), rolePermissionHandler.GetRolePermission)
	}
	groupRouter := router.Group("/api/v1/groups")
	{
		groupRouter.GET("", middleware.AuthMiddleware(cfg), groupHandler.GetGroups)
		groupRouter.POST("", middleware.AuthMiddleware(cfg), middleware.RBACMiddleware("admin", "manager", "user"), groupHandler.CreateGroup)
		groupRouter.PUT("/:id", middleware.AuthMiddleware(cfg), middleware.RBACMiddleware("admin", "manager", "user"), groupHandler.UpdateGroup)
		groupRouter.DELETE("/:id", middleware.AuthMiddleware(cfg), middleware.RBACMiddleware("admin", "manager", "user"), groupHandler.DeleteGroup)
		groupRouter.GET("/:id", middleware.AuthMiddleware(cfg), groupHandler.GetGroup)
	}

	userGroupRouter := router.Group("/api/v1/user-groups")
	{
		userGroupRouter.GET("", middleware.AuthMiddleware(cfg), userGroupHandler.GetUserGroups)
		userGroupRouter.GET("/:id", middleware.AuthMiddleware(cfg), userGroupHandler.GetUserGroup)
		userGroupRouter.POST("", middleware.AuthMiddleware(cfg), middleware.RBACMiddleware("admin", "manager", "user"), userGroupHandler.CreateUserGroup)
		userGroupRouter.PUT("/:id", middleware.AuthMiddleware(cfg), middleware.RBACMiddleware("admin", "manager", "user"), userGroupHandler.UpdateUserGroup)
		userGroupRouter.DELETE("/:id", middleware.AuthMiddleware(cfg), middleware.RBACMiddleware("admin", "manager", "user"), userGroupHandler.DeleteUserGroup)
	}

	groupRoleRouter := router.Group("/api/v1/group-roles")
	{
		groupRoleRouter.GET("", middleware.AuthMiddleware(cfg), groupRoleHandler.GetGroupRoles)
		groupRoleRouter.GET("/:id", middleware.AuthMiddleware(cfg), groupRoleHandler.GetGroupRole)
		groupRoleRouter.POST("", middleware.AuthMiddleware(cfg), middleware.RBACMiddleware("admin", "manager", "user"), groupRoleHandler.CreateGroupRole)
		groupRoleRouter.PUT("/:id", middleware.AuthMiddleware(cfg), middleware.RBACMiddleware("admin", "manager", "user"), groupRoleHandler.UpdateGroupRole)
		groupRoleRouter.DELETE("/:id", middleware.AuthMiddleware(cfg), middleware.RBACMiddleware("admin", "manager", "user"), groupRoleHandler.DeleteGroupRole)
	}

	actionLogRouter := router.Group("/api/v1/audit-logs")
	{
		actionLogRouter.GET("", middleware.AuthMiddleware(cfg), actionLogHandler.GetActionLogs)
		actionLogRouter.GET("/:id", middleware.AuthMiddleware(cfg), actionLogHandler.GetActionLog)
		actionLogRouter.DELETE("/:id", middleware.AuthMiddleware(cfg), middleware.RBACMiddleware("admin", "manager", "user"), actionLogHandler.DeleteActionLog)
		actionLogRouter.POST("", middleware.AuthMiddleware(cfg), middleware.RBACMiddleware("admin", "manager", "user"), actionLogHandler.CreateActionLog)
		actionLogRouter.PUT("/:id", middleware.AuthMiddleware(cfg), middleware.RBACMiddleware("admin", "manager", "user"), actionLogHandler.UpdateActionLog)
	}

	productRouter := router.Group("/api/v1/products")
	{
		productRouter.GET("", handlers.GetProducts)
	}

	return router
}
