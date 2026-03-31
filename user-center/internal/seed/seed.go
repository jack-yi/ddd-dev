package seed

import (
	"context"
	"log"

	"github.com/yangboyi/ddd-dev/user-center/internal/domain/user/entity"
	"github.com/yangboyi/ddd-dev/user-center/internal/domain/user/repository"
	"golang.org/x/crypto/bcrypt"
)

func SeedRoles(ctx context.Context, roleRepo repository.RoleRepository) {
	roles := []entity.Role{
		{Name: entity.RoleSuperAdmin, Description: "超级管理员", IsDefault: false},
		{Name: entity.RoleAdmin, Description: "管理员", IsDefault: false},
		{Name: entity.RoleOperator, Description: "运营人员", IsDefault: false},
		{Name: entity.RoleViewer, Description: "只读用户", IsDefault: true},
	}
	for _, role := range roles {
		if err := roleRepo.SaveIfNotExist(ctx, &role); err != nil {
			log.Printf("warn: seed role %s failed: %v", role.Name, err)
		}
	}
	log.Println("roles seeded successfully")
}

func SeedSuperAdmin(ctx context.Context, userRepo repository.UserRepository, roleRepo repository.RoleRepository) {
	has, err := roleRepo.HasSuperAdmin(ctx)
	if err != nil {
		log.Printf("warn: check super admin failed: %v", err)
		return
	}
	if has {
		return
	}

	// Create super admin user with password
	passwordHash, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("warn: hash password failed: %v", err)
		return
	}

	admin := entity.NewAdminUser("admin", string(passwordHash), "admin@localhost")
	if err := userRepo.Save(ctx, admin); err != nil {
		log.Printf("warn: create super admin failed: %v", err)
		return
	}

	// Assign super_admin role
	role, err := roleRepo.FindByName(ctx, entity.RoleSuperAdmin)
	if err != nil {
		log.Printf("warn: find super_admin role failed: %v", err)
		return
	}
	if err := roleRepo.AssignRole(ctx, admin.ID, role.ID); err != nil {
		log.Printf("warn: assign super_admin role failed: %v", err)
		return
	}

	log.Println("========================================")
	log.Println("Super admin created:")
	log.Println("  Username: admin")
	log.Println("  Password: admin123")
	log.Println("  Please change password after first login!")
	log.Println("========================================")
}
