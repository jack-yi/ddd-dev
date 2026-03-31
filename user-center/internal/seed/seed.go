package seed

import (
	"context"
	"log"

	"github.com/yangboyi/ddd-dev/user-center/internal/domain/user/entity"
	"github.com/yangboyi/ddd-dev/user-center/internal/domain/user/repository"
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
