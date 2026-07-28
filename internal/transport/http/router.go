package http

import (
	"github.com/arrase21/mobileapi/internal/domain"
	"github.com/arrase21/mobileapi/internal/service"
	"github.com/gin-gonic/gin"
)

func NewRouter(
	userSvc *service.UserService,
	roleSvc *service.RoleService,
	assessmentSvc *service.AssessmentService,
	skinfoldSvc *service.SkinfoldService,
	tenantSvc *service.TenantService,
	permissionSvc *service.PermissionService,
	permissionActionSvc *service.PermissionActionService,
	userRepo domain.UserRepo,
	roleRepo domain.RoleRepo,
	projectID string,
) *gin.Engine {
	r := gin.Default()
	r.Use(CORSMiddleware())
	r.Use(LoggingMiddleware())
	r.Use(NewIPRateLimiter(1, 60).Middleware())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "ok",
			"service": "user",
		})
	})

	fbAuth := NewFirebaseAuth(projectID)

	tenants := r.Group("/api/v1/tenants")
	tenants.Use(fbAuth.Middleware())
	{
		tenantHandler := NewTenantHandler(tenantSvc)
		tenants.POST("", SuperAdminMiddleware(userRepo), tenantHandler.Create)
		tenants.GET("", SuperAdminMiddleware(userRepo), tenantHandler.List)
		tenants.GET("/:id", SuperAdminMiddleware(userRepo), tenantHandler.GetByID)
		tenants.PUT("/:id", SuperAdminMiddleware(userRepo), tenantHandler.Update)
		tenants.DELETE("/:id", SuperAdminMiddleware(userRepo), tenantHandler.Delete)
	}

	v1 := r.Group("/api/v1")
	{
		usrs := v1.Group("/users")
		usrs.Use(fbAuth.Middleware())
		{
			userHandler := NewUserHandler(userSvc)
			usrs.POST("", RequirePermission(userRepo, roleRepo, "user.create"), userHandler.Create)
			usrs.GET("", RequirePermission(userRepo, roleRepo, "user.read"), userHandler.List)
			usrs.GET("/deleted", RequirePermission(userRepo, roleRepo, "user.read"), userHandler.ListDeleted)
			usrs.GET("/:dni", RequirePermission(userRepo, roleRepo, "user.read"), userHandler.GetByDni)
			usrs.PUT("/:user", RequirePermission(userRepo, roleRepo, "user.update"), userHandler.Update)
			usrs.DELETE("/:user", RequirePermission(userRepo, roleRepo, "user.delete"), userHandler.SoftDelete)
			usrs.PATCH("/:user/restore", RequirePermission(userRepo, roleRepo, "user.restore"), userHandler.Restore)
			usrs.DELETE("/:user/permanent", RequirePermission(userRepo, roleRepo, "user.delete"), userHandler.Delete)
		}

		roles := v1.Group("/roles")
		roles.Use(fbAuth.Middleware())
		{
			roleHandler := NewRoleHandler(roleSvc)
			roles.POST("", RequirePermission(userRepo, roleRepo, "role.create"), roleHandler.CreateRole)
			roles.GET("", RequirePermission(userRepo, roleRepo, "role.read"), roleHandler.List)
			roles.GET("/:id", RequirePermission(userRepo, roleRepo, "role.read"), roleHandler.GetByID)
			roles.PUT("/:id", RequirePermission(userRepo, roleRepo, "role.update"), roleHandler.Update)
			roles.DELETE("/:id", RequirePermission(userRepo, roleRepo, "role.delete"), roleHandler.SoftDelete)
			roles.PATCH("/:id/restore", RequirePermission(userRepo, roleRepo, "role.restore"), roleHandler.Restore)
		}

		assessments := v1.Group("/assessments")
		assessments.Use(fbAuth.Middleware())
		{
			assessmentHandler := NewAssessmentHandler(assessmentSvc)
			assessments.POST("", RequirePermission(userRepo, roleRepo, "assessment.create"), assessmentHandler.Create)
			assessments.GET("", RequirePermission(userRepo, roleRepo, "assessment.read"), assessmentHandler.ListByTenant)
			assessments.GET("/user/:user_id", RequirePermission(userRepo, roleRepo, "assessment.read"), assessmentHandler.ListByUser)
			assessments.GET("/:id", RequirePermission(userRepo, roleRepo, "assessment.read"), assessmentHandler.GetByID)
			assessments.PUT("/:id", RequirePermission(userRepo, roleRepo, "assessment.update"), assessmentHandler.Update)
			assessments.DELETE("/:id", RequirePermission(userRepo, roleRepo, "assessment.delete"), assessmentHandler.SoftDelete)
			assessments.PATCH("/:id/restore", RequirePermission(userRepo, roleRepo, "assessment.update"), assessmentHandler.Restore)
			assessments.DELETE("/:id/permanent", RequirePermission(userRepo, roleRepo, "assessment.delete"), assessmentHandler.Delete)
		}

		skinfolds := v1.Group("/skinfolds")
		skinfolds.Use(fbAuth.Middleware())
		{
			skinfoldHandler := NewSkinfoldHandler(skinfoldSvc)
			skinfolds.POST("", RequirePermission(userRepo, roleRepo, "assessment.create"), skinfoldHandler.Create)
			skinfolds.GET("/assessment/:assessment_id", RequirePermission(userRepo, roleRepo, "assessment.read"), skinfoldHandler.ListByAssessment)
			skinfolds.GET("/:id", RequirePermission(userRepo, roleRepo, "assessment.read"), skinfoldHandler.GetByID)
			skinfolds.PUT("/:id", RequirePermission(userRepo, roleRepo, "assessment.update"), skinfoldHandler.Update)
			skinfolds.DELETE("/:id", RequirePermission(userRepo, roleRepo, "assessment.delete"), skinfoldHandler.SoftDelete)
			skinfolds.PATCH("/:id/restore", RequirePermission(userRepo, roleRepo, "assessment.update"), skinfoldHandler.Restore)
			skinfolds.DELETE("/:id/permanent", RequirePermission(userRepo, roleRepo, "assessment.delete"), skinfoldHandler.Delete)
		}

		permissions := v1.Group("/permissions")
		permissions.Use(fbAuth.Middleware())
		{
			permissionHandler := NewPermissionHandler(permissionSvc)
			permissions.POST("", RequirePermission(userRepo, roleRepo, "permission.create"), permissionHandler.Create)
			permissions.GET("", RequirePermission(userRepo, roleRepo, "permission.read"), permissionHandler.List)
			permissions.GET("/:id", RequirePermission(userRepo, roleRepo, "permission.read"), permissionHandler.GetByID)
			permissions.PUT("/:id", RequirePermission(userRepo, roleRepo, "permission.update"), permissionHandler.Update)
			permissions.DELETE("/:id", RequirePermission(userRepo, roleRepo, "permission.delete"), permissionHandler.SoftDelete)
			permissions.PATCH("/:id/restore", RequirePermission(userRepo, roleRepo, "permission.restore"), permissionHandler.Restore)
		}

		permissionActions := v1.Group("/permission-actions")
		permissionActions.Use(fbAuth.Middleware())
		{
			actionHandler := NewPermissionActionHandler(permissionActionSvc)
			permissionActions.POST("", RequirePermission(userRepo, roleRepo, "permission.create"), actionHandler.Create)
			permissionActions.GET("", RequirePermission(userRepo, roleRepo, "permission.read"), actionHandler.List)
			permissionActions.GET("/:id", RequirePermission(userRepo, roleRepo, "permission.read"), actionHandler.GetByID)
			permissionActions.PUT("/:id", RequirePermission(userRepo, roleRepo, "permission.update"), actionHandler.Update)
			permissionActions.DELETE("/:id", RequirePermission(userRepo, roleRepo, "permission.delete"), actionHandler.SoftDelete)
			permissionActions.PATCH("/:id/restore", RequirePermission(userRepo, roleRepo, "permission.restore"), actionHandler.Restore)
		}
	}
	return r
}
