package test

import (
	assetsDao "seccloud/cspm/modules/assets/dao"
	"seccloud/cspm/modules/assets/modules/kubernetes/dao"
	assetsSvc "seccloud/cspm/modules/assets/service"
	groupDao "seccloud/cspm/modules/policies/modules/groups/dao"
	sharedTenantSvc "seccloud/cspm/modules/shared/service/tenant"
	userService "seccloud/cspm/modules/users/service"
	"seccloud/cspm/pkg/container"
	"seccloud/cspm/pkg/utils/opt"
	"seccloud/cspm/providers/cache"
	"seccloud/cspm/providers/global"
	"seccloud/cspm/providers/grpcs"
	"seccloud/cspm/providers/uuid"
)

func Provide(opts ...opt.Option) error {
	if err := container.Container.Provide(func(
		accountDao *dao.AccountDao,
		accountMetaDao *assetsDao.AccountMetaDao,
		accountMetaSvc *assetsSvc.AccountMetaService,
		cache *cache.Cache,
		clusterDao *dao.KubeClusterDao,
		clusterSvc *KubeClusterService,
		global *global.Config,
		grpcs *grpcs.Grpc,
		imageDao *dao.KubeImageDao,
		kubeGroupDao *groupDao.KubeGroupDao,
		kubeGroupMonitorDao *groupDao.KubeGroupMonitorFileDao,
		kubeGroupProcessDao *groupDao.KubeGroupProcessDao,
		namespaceDao *dao.KubeNamespaceDao,
		nodeDao *dao.KubeNodeDao,
		podContainersDao *dao.KubePodContainerDao,
		podDao *dao.KubePodDao,
		serviceDao *dao.KubeServiceDao,
		sharedTenantSvc *sharedTenantSvc.TenantService,
		userService *userService.UserService,
		uuid *uuid.Generator,
		workloadDao *dao.KubeWorkloadDao,
	) (*AccountService, error) {
		obj := &AccountService{
			accountDao:          accountDao,
			accountMetaDao:      accountMetaDao,
			accountMetaSvc:      accountMetaSvc,
			cache:               cache,
			clusterDao:          clusterDao,
			clusterSvc:          clusterSvc,
			global:              global,
			grpcs:               grpcs,
			imageDao:            imageDao,
			kubeGroupDao:        kubeGroupDao,
			kubeGroupMonitorDao: kubeGroupMonitorDao,
			kubeGroupProcessDao: kubeGroupProcessDao,
			namespaceDao:        namespaceDao,
			nodeDao:             nodeDao,
			podContainersDao:    podContainersDao,
			podDao:              podDao,
			serviceDao:          serviceDao,
			sharedTenantSvc:     sharedTenantSvc,
			userService:         userService,
			uuid:                uuid,
			workloadDao:         workloadDao,
		}
		return obj, nil
	}); err != nil {
		return err
	}

	return nil
}
