package test

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"seccloud/cspm/database/models"
	"seccloud/cspm/modules/assets/modules/kubernetes/dao"
	"seccloud/cspm/modules/assets/modules/kubernetes/dto"
	"seccloud/cspm/pkg/common"
	"seccloud/cspm/pkg/common/pg"
	"seccloud/cspm/pkg/consts"
	"seccloud/cspm/pkg/errorx"
	"seccloud/cspm/pkg/utils"
	"seccloud/cspm/providers/cache"
	"seccloud/cspm/providers/global"
	"seccloud/cspm/providers/grpcs"
	"seccloud/cspm/providers/log"
	"seccloud/cspm/providers/uuid"
	"strconv"
	"strings"

	assetsDao "seccloud/cspm/modules/assets/dao"
	assetsSvc "seccloud/cspm/modules/assets/service"
	groupDao "seccloud/cspm/modules/policies/modules/groups/dao"
	sharedDto "seccloud/cspm/modules/shared/dto"
	sharedTenantSvc "seccloud/cspm/modules/shared/service/tenant"
	userService "seccloud/cspm/modules/users/service"

	"git-biz.qianxin-inc.cn/seccloud/cspm-proto/gen/proto/platform/v1alpha2"
	"github.com/golang-module/carbon/v2"
	"github.com/pkg/errors"
	"github.com/samber/lo"
)

// @provider:except
type AccountService struct {
	uuid   *uuid.Generator
	global *global.Config
	grpcs  *grpcs.Grpc
	cache  *cache.Cache

	accountDao          *dao.AccountDao
	accountMetaSvc      *assetsSvc.AccountMetaService
	accountMetaDao      *assetsDao.AccountMetaDao
	clusterSvc          *KubeClusterService
	clusterDao          *dao.KubeClusterDao
	nodeDao             *dao.KubeNodeDao
	podDao              *dao.KubePodDao
	podContainersDao    *dao.KubePodContainerDao
	imageDao            *dao.KubeImageDao
	namespaceDao        *dao.KubeNamespaceDao
	serviceDao          *dao.KubeServiceDao
	workloadDao         *dao.KubeWorkloadDao
	userService         *userService.UserService
	kubeGroupDao        *groupDao.KubeGroupDao
	kubeGroupMonitorDao *groupDao.KubeGroupMonitorFileDao
	kubeGroupProcessDao *groupDao.KubeGroupProcessDao

	sharedTenantSvc *sharedTenantSvc.TenantService
}

func (svc *AccountService) GetByID(ctx context.Context, id int32) (*models.Account, error) {
	return svc.accountDao.GetByID(ctx, id)
}

func (svc *AccountService) PageByQueryFilter(
	ctx context.Context,
	queryFilter *dto.AccountListQueryFilter,
	pageFilter *common.PageQueryFilter,
	sortFilter *common.SortQueryFilter,
) ([]*models.Account, int64, error) {
	if pageFilter != nil {
		pageFilter = pageFilter.Format()
	}
	queryFilter.PlatformType = consts.PlatformTypeKubernetes.String()
	if queryFilter.BenchmarkCheckedAt != nil {
		benchCheckedArr := strings.Split(*queryFilter.BenchmarkCheckedAt, ",")
		if len(benchCheckedArr) == 2 {
			accountIds := svc.accountMetaDao.GetAccountIdsByCheckedAtRange(ctx, benchCheckedArr[0], benchCheckedArr[1])
			queryFilter.AccountIDs = accountIds
		}
	}
	return svc.accountDao.PageByQueryFilter(ctx, queryFilter, pageFilter, sortFilter)
}

func (svc *AccountService) DecorateItem(model *models.Account, _ int) *dto.Cluster {
	if model == nil {
		return nil
	}
	ctx := context.Background()

	resp := &dto.Cluster{
		ID:                  model.ID,
		AccountID:           model.ID,
		TenantID:            model.TenantID,
		TenantName:          "",
		Name:                model.Name,
		PlatformID:          model.CloudType,
		PlatformName:        consts.CloudType(model.CloudType).String(),
		Version:             model.Version,
		Statistics:          dto.ClusterStatistics{},
		BenchmarkStatistics: sharedDto.BenchmarkStatistics{},
		BenchmarkCheckedAt:  "",
		CreatedAt:           model.CreatedAt.Format(utils.TIME_RFCFE),
		Status:              model.Status,
		TrafficolWorkmode:   "",
		StatusName:          consts.PlatformStatus(model.Status).String(),
		LastSyncAt:          "",
	}

	meta, err := svc.accountMetaSvc.GetByAccountID(ctx, model.ID, consts.AccountMetaKeyCheckedAt)
	if err == nil {
		timestamp, _ := strconv.Atoi(meta.Value)
		resp.BenchmarkCheckedAt = carbon.CreateFromTimestamp(int64(timestamp)).ToDateTimeString()
	}

	tenantName, err := svc.sharedTenantSvc.GetNameByID(ctx, model.TenantID)
	if err == nil {
		resp.TenantName = tenantName
	}

	statistics, err := svc.GetStatistics(ctx, model.ID)
	if err == nil {
		resp.Statistics = statistics
	}

	benchmarkStatistics, err := svc.GetBenchmarkStatistics(ctx, model.ID)
	if err == nil {
		resp.BenchmarkStatistics = benchmarkStatistics
	}

	trafficWorkmode, err := svc.clusterSvc.GetTrafficolWorkmodeByAccountID(ctx, model.ID)
	if err == nil {
		resp.TrafficolWorkmode = trafficWorkmode.String()
	}

	lastSyncAt, err := svc.accountMetaSvc.GetLastSyncAt(ctx, model.ID)
	if err == nil {
		resp.LastSyncAt = lastSyncAt.ToDateTimeString()
	}
	// log.Infof("id:%d",model.ID)

	return resp
}

// Statistics
func (svc *AccountService) GetStatistics(ctx context.Context, id int32) (dto.ClusterStatistics, error) {
	statistics := dto.ClusterStatistics{
		NodeMaster: 0,
		NodeWorker: 0,
		Pod:        0,
		Image:      0,
		Namespace:  0,
		Service:    0,
		Workload:   0,
	}

	nodeFilter := &dto.KubeNodeListQueryFilter{}
	nodeFilter.AccountID = &id
	nodeFilter.Roles = lo.ToPtr("control-plane")

	c, err := svc.nodeDao.GetAmountByQueryFilter(ctx, nodeFilter)
	if err == nil {
		statistics.NodeMaster = int32(c)
	}

	nodeFilter.Roles = lo.ToPtr("node")
	c, err = svc.nodeDao.GetAmountByQueryFilter(ctx, nodeFilter)
	if err == nil {
		statistics.NodeWorker = int32(c)
	}

	podFilter := &dto.KubePodListQueryFilter{}
	podFilter.AccountID = &id
	c, err = svc.podDao.GetPodAmountByFilter(ctx, podFilter)
	if err == nil {
		statistics.Pod = int32(c)
	}

	c, err = svc.imageDao.GetClusterImageAmount(ctx, id)
	if err == nil {
		statistics.Image = int32(c)
	}

	c, err = svc.namespaceDao.GetAmountByAccountID(ctx, id)
	if err == nil {
		statistics.Namespace = int32(c)
	}

	c, err = svc.serviceDao.GetAmountByAccountID(ctx, id)
	if err == nil {
		statistics.Service = int32(c)
	}

	c, err = svc.workloadDao.GetAmountByAccountID(ctx, id)
	if err == nil {
		statistics.Workload = int32(c)
	}

	return statistics, nil
}

// Statistics
func (svc *AccountService) GetBenchmarkStatistics(ctx context.Context, accountId int32) (sharedDto.BenchmarkStatistics, error) {
	statistics := sharedDto.BenchmarkStatistics{
		All:     0,
		High:    0,
		Medium:  0,
		Low:     0,
		Safe:    0,
		Ignore:  0,
		Warning: 0,
	}
	k8sList, err := svc.accountDao.GroupByAccountId(ctx, accountId)
	if err != nil {
		return statistics, err
	}
	return sharedDto.NewBenchmarkStatistics(k8sList), nil
}

// Create
func (svc *AccountService) Create(ctx context.Context, userID, tenantID int32, body *dto.AccountForm) error {
	if !svc.accountDao.IsNameUnique(ctx, body.Name, nil) {
		return errors.New("集群名称重复")
	}

	if body.ClusterID == "" {
		return errors.New("集群ID不能为空")
	}

	return svc.accountDao.Transaction(func() error {
		model := &models.Account{
			TenantID:     tenantID,
			UserID:       userID,
			CloudType:    body.PlatformID,
			PlatformType: string(consts.PlatformTypeKubernetes),
			Config:       pg.AccountConfigFromString(body.KubeConfig).ToField(),
			Name:         body.Name,
			Status:       consts.PlatformStatusInit,
		}

		if err := svc.accountDao.Create(ctx, model); err != nil {
			return err
		}

		// create cluster
		clusterModel := &models.KubeCluster{
			UUID:              body.ClusterID,
			AccountID:         model.ID,
			TenantID:          tenantID,
			TrafficolWorkmode: body.TrafficolWorkmode.String(),
		}
		if err := svc.clusterDao.Create(ctx, clusterModel); err != nil {
			return err
		}
		return nil
	})
}

func (svc *AccountService) UpdateFromModel(ctx context.Context, id int32, model *models.Account) error {
	return svc.accountDao.Transaction(func() error {
		if err := svc.accountDao.Update(ctx, model); err != nil {
			return err
		}

		cluster, err := svc.clusterSvc.GetByID(ctx, id)
		if err != nil {
			return err
		}
		// update cluster
		clusterModel := &models.KubeCluster{
			TrafficolWorkmode: cluster.TrafficolWorkmode,
		}
		if err := svc.clusterDao.UpdateByAccountID(ctx, id, clusterModel); err != nil {
			return err
		}
		return nil
	})
}

func (svc *AccountService) UpdateVersionStatusAndTrafficMode(ctx context.Context, id int32, version string, status consts.PlatformStatus, trafficMode string) error {
	return svc.accountDao.Transaction(func() error {
		model, err := svc.GetByID(ctx, id)
		if err != nil {
			return err
		}
		model.Status = status
		model.Version = version
		if err := svc.accountDao.Update(ctx, model); err != nil {
			return err
		}

		// update cluster
		clusterModel := &models.KubeCluster{
			TrafficolWorkmode: trafficMode,
		}
		if err := svc.clusterDao.UpdateByAccountID(ctx, id, clusterModel); err != nil {
			return err
		}
		return nil
	})
}

// Update
func (svc *AccountService) Update(ctx context.Context, id int32, body *dto.AccountForm) error {
	if !svc.accountDao.IsNameUnique(ctx, body.Name, []int32{id}) {
		return errors.New("集群名称重复")
	}

	return svc.accountDao.Transaction(func() error {
		model, err := svc.GetByID(ctx, id)
		if err != nil {
			return err
		}

		m := &models.Account{
			TenantID:     model.TenantID,
			UserID:       model.UserID,
			PlatformType: body.PlatformID,
			Config:       pg.AccountConfigFromString(body.KubeConfig).ToField(),
			Name:         body.Name,
			Status:       consts.PlatformStatusInit,
		}

		if err := svc.accountDao.Update(ctx, m); err != nil {
			return err
		}

		// update cluster
		clusterModel := &models.KubeCluster{
			TrafficolWorkmode: body.TrafficolWorkmode.String(),
		}
		if err := svc.clusterDao.UpdateByAccountID(ctx, id, clusterModel); err != nil {
			return err
		}
		return nil
	})
}

// Delete
func (svc *AccountService) Delete(ctx context.Context, loginUid, id int32) error {
	m, err := svc.GetByID(ctx, id)
	if err != nil {
		return err
	}
	isSuper, isAdmin := svc.userService.IsSuperOrAdmin(ctx, loginUid)
	// 删除自己创建的
	if !isSuper && !isAdmin && m.UserID != loginUid {
		return errorx.ErrForbidden
	}

	deleteFuncs := []func(ctx context.Context, id int32) error{
		svc.accountDao.Delete,
		svc.accountMetaDao.DeleteByAccountID,
		svc.clusterDao.DeleteByAccountID,
		svc.nodeDao.DeleteByAccountID,
		svc.podDao.DeleteByAccountID,
		svc.podContainersDao.DeleteByAccountID,
		svc.imageDao.DeleteByAccountID,
		svc.namespaceDao.DeleteByAccountID,
		svc.serviceDao.DeleteByAccountID,
		svc.workloadDao.DeleteByAccountID,
		svc.kubeGroupDao.DeleteByAccountID,
		svc.kubeGroupMonitorDao.DeleteByAccountID,
		svc.kubeGroupProcessDao.DeleteByAccountID,
	}

	return svc.accountDao.Transaction(func() error {
		for _, f := range deleteFuncs {
			if err := f(ctx, m.ID); err != nil {
				return err
			}
		}

		return nil
	})
}

// GetPlatforms
func (svc *AccountService) GetPlatforms(ctx context.Context) []dto.PlatformTypeItem {
	platformTypes := []dto.PlatformTypeItem{}

	for _, v := range consts.KubeTypeNames() {
		kt, err := consts.ParseKubeType(v)
		if err != nil {
			continue
		}

		cn, err := kt.CnName()
		if err != nil {
			continue
		}

		platformTypes = append(platformTypes, dto.PlatformTypeItem{
			ID:   v,
			Name: cn,
		})
	}
	return platformTypes
}

// GetVersions
func (svc *AccountService) GetVersions(ctx context.Context) ([]string, error) {
	return svc.accountDao.GetVersions(ctx)
}

// CreateClusterID
func (svc *AccountService) CreateClusterID(ctx context.Context, body *dto.CreateClusterIDForm) (*dto.CreateClusterIDResponse, error) {
	var err error
	clusterModel := &models.KubeCluster{
		UUID:              svc.uuid.MustGenerate(),
		TrafficolWorkmode: body.TrafficolWorkmode.String(),
	}
	if body.ID != nil {
		clusterModel, err = svc.clusterDao.GetByAccountID(ctx, *body.ID)
		if err != nil {
			return nil, err
		}
	}

	uuidSections := strings.Split(clusterModel.UUID, "-")
	encryptUUID := uuidSections[len(uuidSections)-1]

	message := &dto.ClusterOperationMessageVars{
		ClusterId:   clusterModel.UUID,
		JoinToken:   encryptUUID,
		GrpcAddress: svc.grpcs.Address(),
	}

	msg, err := message.Render()
	if err != nil {
		return nil, err
	}

	return &dto.CreateClusterIDResponse{
		ClusterID:        clusterModel.UUID,
		OperationMessage: msg,
	}, nil
}

func (svc *AccountService) ClusterNodeOption(ctx context.Context, queryFilter *dto.ClusterOptionQueryFilter) ([]dto.ClusterSection, error) {
	list, err := svc.clusterDao.ClusterNodeOption(ctx, queryFilter)
	if err != nil {
		return nil, err
	}
	accountIds := lo.Map(list, func(item *models.Account, index int) int32 {
		return item.ID
	})

	kubeNodeList, err := svc.nodeDao.GetNodesSectionByAccountIds(ctx, accountIds)
	if err != nil {
		return nil, err
	}
	result := []dto.ClusterSection{}
	for _, account := range list {
		cluster := svc.FormatClusterSection(account)
		tenantName, err := svc.sharedTenantSvc.GetNameByID(ctx, account.TenantID)
		if err != nil {
			continue
		}
		cluster.TenantName = tenantName
		for _, node := range kubeNodeList {
			r := dto.KubeNodeEntity{}
			r.KubeNodeSection = node
			r.EntityId = fmt.Sprintf("%s%d", consts.TASK_NODE_PREFIX, node.Id)

			if node.AccountId == account.ID {
				cluster.Nodes = append(cluster.Nodes, r)
			}
		}

		result = append(result, cluster)
	}
	return result, nil
}

// 格式化给前端显示 FormatClusterSection
func (svc *AccountService) FormatClusterSection(account *models.Account) dto.ClusterSection {
	statusName, err := account.Status.CnName()
	if err != nil {
		log.Error("platformStatus.CnName err:%v", err)
		return dto.ClusterSection{}
	}

	clusterSection := dto.ClusterSection{
		ID:           account.ID,
		EntityId:     fmt.Sprintf("%s%d", consts.TASK_ACCOUNT_PREFIX, account.ID),
		TenantId:     account.TenantID,
		TenantName:   "",
		Name:         account.Name,
		PlatformId:   account.CloudType,
		PlatformName: account.PlatformType,
		Version:      account.Version,
		CreatedAt:    account.CreatedAt.String(),
		Status:       account.Status,
		StatusName:   statusName,
		Nodes:        nil,
	}

	return clusterSection
}

// GetIDsByName
func (svc *AccountService) GetIDsByName(ctx context.Context, name string) []int32 {
	return svc.accountDao.GetIDsByName(ctx, name)
}

func (svc *AccountService) Join(ctx context.Context, token string) (*v1alpha2.JoinResp, error) {
	cluster, err := svc.clusterDao.GetByClusterIDLastSection(ctx, token)
	if err != nil {
		return nil, errors.Wrapf(err, "cluster not found")
	}

	// if carbon.Now().Carbon2Time().After(cluster.CreatedAt) {
	// 	return nil, errors.New("cluster token expired")
	// }

	account, err := svc.GetByID(ctx, cluster.AccountID)
	if err != nil {
		return nil, errors.Wrapf(err, "account not found")
	}
	if account.Status != consts.PlatformStatusInit {
		return nil, errors.New("cluster status is not init")
	}

	cert, key, ca, err := svc.getCertInfo()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get cert info")
	}

	regServer, regUser, regPass := svc.getRegistryInfo()

	resp := v1alpha2.JoinResp{
		ClusterId:         cluster.UUID,
		Cert:              string(cert),
		Key:               string(key),
		CaCert:            string(ca),
		RegistryServer:    regServer,
		RegistryUsername:  regUser,
		RegistryPassword:  regPass,
		TrafficolWorkmode: cluster.TrafficolWorkmode,
	}

	// 先放开，稳定后再改状态
	// account.Status = consts.PlatformStatusJoin.String()
	// if err := svc.accountDao.Update(ctx, account); err != nil {
	//	return nil, err
	// }

	log.Infof("Cluster %s join ok", cluster.UUID)

	return &resp, nil
}

func (svc *AccountService) getCertInfo() (cert, key, ca []byte, err error) {
	ca, err = os.ReadFile(svc.global.CaCert)
	if err != nil {
		return
	}

	cert, err = os.ReadFile(svc.global.SslCert)
	if err != nil {
		return
	}

	key, err = os.ReadFile(svc.global.SslKey)
	if err != nil {
		return
	}

	return
}

func (svc *AccountService) getRegistryInfo() (server, username, password string) {
	var dockerConfig struct {
		Auths map[string]struct {
			Auth string `json:"auth"`
		} `json:"auths"`
	}

	fp, err := os.Open(consts.DockerConfig)
	if err != nil {
		return
	}

	defer fp.Close()

	jd := json.NewDecoder(fp)
	if err = jd.Decode(&dockerConfig); err != nil {
		log.Errorf("Failed to decode docker config:%s %v", consts.DockerConfig, err)
		return
	}

	if len(dockerConfig.Auths) == 0 {
		log.Errorf("No auths in docker config:%s", consts.DockerConfig)
		return
	}

	server = lo.Keys(dockerConfig.Auths)[0]
	auth := dockerConfig.Auths[server]
	cred, err := base64.StdEncoding.DecodeString(auth.Auth)
	if err != nil {
		log.Errorf("Base64 decode auth string:%s", auth.Auth)
		return
	}
	userPass := strings.SplitN(string(cred), ":", 2)
	if len(userPass) != 2 {
		log.Errorf("Invalid credential %s", string(cred))
		return
	}

	username, password = userPass[0], userPass[1]
	log.Infof("Success got registry %s username:%s", server, username)

	return
}
