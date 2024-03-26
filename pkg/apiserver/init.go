package apiserver

import (
	"context"
	"encoding/csv"
	"v1/pkg/utils"

	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	"os"

	"time"
	"v1/pkg/dao"
	"v1/pkg/model"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

func (s *APIServer) initSystem() (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*10)
	defer cancel()

	zap.L().Info("initSystem....")
	defer func(t1 time.Time) {
		if err != nil {
			zap.L().Info("initSystem exception！！！", zap.Duration("since", time.Since(t1)), zap.Error(err))
		} else {
			zap.L().Info("initSystem done!", zap.Duration("since", time.Since(t1)), zap.Error(err))
		}
	}(time.Now())

	var errs []error

	errs = append(errs, initSuperAdmin(ctx, s.RDBClient))
	// errs = append(errs, initConfig(ctx, s.RDBClient))
	// errs = append(errs, initDefaultBenchmark(ctx, s.RDBClient))
	// errs = append(errs, initRiskScanTask(ctx, s.RDBClient))

	// scan.AssetsScanBootUpRepair(ctx, s.RDBClient, s.Sched)
	// scan.RiskScanBootUpRepair(ctx, s.RDBClient, s.Sched)

	if len(errs) != 0 {
		err = utilerrors.NewAggregate(errs)
	}

	return err
}

func initSuperAdmin(ctx context.Context, db *gorm.DB) error {
	_, err := dao.GetUserByUsername(ctx, db, model.SuperAdminUsername)
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			return err
		}
	} else {
		zap.L().Debug("superadmin already exists，skip")
		return nil
	}

	// 初始化超管
	_, err = dao.InsertUser(ctx, db, model.User{
		UID:      utils.NextID(),
		Username: model.SuperAdminUsername,
		Password: utils.MD5Hex(model.SuperAdminDefaultPassword),
		Name:     "默认管理员",
		Creator:  model.SystemUsername,
		Updater:  model.SystemUsername,
		Status:   model.UserStatusNormal,
		Phone:    "",
		Emial:    "",
		Role:     model.RoleTypeSuperAdmin,
	})
	if err != nil {
		zap.L().Error("initSuperAdmin error", zap.Error(err))
	}

	return err
}

// func initConfig(ctx context.Context, db *gorm.DB) error {
// 	// configs/compliance.csv
// 	data, err := readCsv("./configs/compliance.csv")
// 	if err != nil {
// 		zap.L().Fatal("read file compliance fails", zap.Error(err))
// 	}
// 	cmap, err := writeCompliances(ctx, db, data)
// 	if err != nil {
// 		zap.L().Fatal("writeCompliances fails", zap.Error(err))
// 	}
//
// 	// configs/policy.csv
// 	data, err = readCsv("./configs/policy.csv")
// 	if err != nil {
// 		zap.L().Fatal("read file policy fails", zap.Error(err))
// 	}
// 	err = writePolicies(ctx, db, data, cmap)
// 	if err != nil {
// 		zap.L().Fatal("writePolicies fails", zap.Error(err))
// 	}
// 	err = rebuildByPolicies(ctx, db)
// 	if err != nil {
// 		zap.L().Fatal("rebuildByPolicies fails", zap.Error(err))
// 	}
//
// 	// configs/region.csv
// 	data, err = readCsv("./configs/region.csv")
// 	if err != nil {
// 		zap.L().Fatal("read file region fails", zap.Error(err))
// 	}
// 	err = writeRegions(ctx, db, data)
// 	if err != nil {
// 		zap.L().Fatal("writeRegions fails", zap.Error(err))
// 	}
//
// 	// configs/service.csv
// 	data, err = readCsv("./configs/service.csv")
// 	if err != nil {
// 		zap.L().Fatal("read file service fails", zap.Error(err))
// 	}
// 	err = writeServices(ctx, db, data)
// 	if err != nil {
// 		zap.L().Fatal("writeServices fails", zap.Error(err))
// 	}
//
// 	// configs/ai_prompt_template.csv
// 	data, err = readCsv("./configs/ai_prompt_template.csv")
// 	if err != nil {
// 		zap.L().Fatal("read file ai_prompt_template fails", zap.Error(err))
// 	}
// 	err = writeAiPromptTemplates(ctx, db, data)
// 	if err != nil {
// 		zap.L().Fatal("writeAiPromptTemplates fails", zap.Error(err))
// 	}
// 	return nil
// }

func readCsv(path string) ([][]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	csvReader := csv.NewReader(f)
	data, err := csvReader.ReadAll()
	if err != nil {
		return nil, err
	}
	return data, nil
}

// func writeCompliances(ctx context.Context, db *gorm.DB, data [][]string) (map[string]struct{}, error) {
// 	complianceCatalogM := make(map[string]struct{})
// 	var err error
// 	err = db.WithContext(ctx).Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&model.Compliance{}).Error
// 	if err != nil {
// 		zap.L().Info("delete compliances fails", zap.Error(err))
// 		return complianceCatalogM, err
// 	}
//
// 	err = db.WithContext(ctx).Exec("ALTER TABLE `compliances` AUTO_INCREMENT = 1;").Error
// 	if err != nil {
// 		zap.L().Info("reset compliances auto_increment fails", zap.Error(err))
// 		return complianceCatalogM, err
// 	}
//
// 	compliances := make([]model.Compliance, 0)
// 	for i := 1; i < len(data); i++ {
// 		compliances = append(compliances, model.Compliance{
// 			CatalogID:        data[i][3],
// 			Key:              data[i][0],
// 			ComplianceType:   data[i][2],
// 			ComplianceTypeEn: data[i][1],
// 			Catalog1:         data[i][5],
// 			Catalog1En:       data[i][4],
// 			Catalog2:         data[i][7],
// 			Catalog2En:       data[i][6],
// 			Catalog3:         data[i][9],
// 			Catalog3En:       data[i][8],
// 			Requirement:      utils.TrimLineBreaks(data[i][11]),
// 			RequirementEn:    utils.TrimLineBreaks(data[i][10]),
// 			RiskScene:        utils.TrimLineBreaks(data[i][13]),
// 			RiskSceneEn:      utils.TrimLineBreaks(data[i][12]),
// 			Suggestion:       utils.TrimLineBreaks(data[i][15]),
// 			SuggestionEn:     utils.TrimLineBreaks(data[i][14]),
// 			CreatedAt:        time.Now().UnixMilli(),
// 		})
// 		complianceCatalogM[data[i][0]+" - "+data[i][3]] = struct{}{}
// 	}
// 	err = db.WithContext(ctx).CreateInBatches(compliances, 100).Error
// 	if err != nil {
// 		zap.L().Info("create compliances fails", zap.Error(err))
// 		return complianceCatalogM, err
// 	}
// 	return complianceCatalogM, nil
// }
//
// func writePolicies(ctx context.Context, db *gorm.DB, data [][]string, cmap map[string]struct{}) error {
// 	var err error
//
// 	dbPolicies := make([]model.BenchmarkPolicy, 0)
// 	err = db.WithContext(ctx).Unscoped().Where("1 = 1").Find(&dbPolicies).Error
// 	if err != nil {
// 		zap.L().Info("find db policies fails", zap.Error(err))
// 		return err
// 	}
// 	dbpmap := make(map[string]model.BenchmarkPolicy)
// 	for i := range dbPolicies {
// 		dbpmap[dbPolicies[i].ID] = dbPolicies[i]
// 	}
//
// 	pmap := make(map[string]struct{})
// 	for i := 1; i < len(data); i++ {
// 		key := data[i][0]
// 		platform := model.CloudPlatform(data[i][3])
// 		if key == "" || platform == "" {
// 			continue
// 		}
// 		apis := make([]string, 0)
// 		if len(data[i][13]) != 0 {
// 			err = json.Unmarshal([]byte(strings.ReplaceAll(data[i][13], `'`, `"`)), &apis)
// 			if err != nil {
// 				// apis使用空值继续写入
// 				zap.L().Info("unmarshal apis fails", zap.Error(err), zap.String("data", data[i][13]))
// 			}
// 		}
//
// 		references := make([]string, 0)
// 		if len(data[i][14]) != 0 {
// 			if !strings.Contains(data[i][14], "[") { // fixme: 暂时兼容处理csv中的两种格式
// 				references = []string{data[i][14]}
// 			} else {
// 				err = json.Unmarshal([]byte(strings.ReplaceAll(data[i][14], `'`, `"`)), &references)
// 				if err != nil {
// 					// references使用空值继续写入
// 					zap.L().Info("unmarshal references fails", zap.Error(err), zap.String("data", data[i][14]))
// 				}
// 			}
// 		}
//
// 		// 简单检查合规的目录映射是否合法
// 		// if _, ok := cmap[model.ComplianceKeyHipaa+" - "+data[i][8]]; !ok && len(data[i][8]) != 0 {
// 		//	zap.L().Error("invalid compliance catalog_id", zap.String("catalog_id", model.ComplianceKeyHipaa+" - "+data[i][8]))
// 		//	return errors.New("invalid compliance catalog_id")
// 		// }
// 		// if _, ok := cmap[model.ComplianceKeyPci+" - "+data[i][9]]; !ok && len(data[i][9]) != 0 {
// 		//	zap.L().Error("invalid compliance catalog_id", zap.String("catalog_id", model.ComplianceKeyPci+" - "+data[i][9]))
// 		//	return errors.New("invalid compliance catalog_id")
// 		// }
// 		if _, ok := cmap[model.ComplianceKeyMlps2+" - "+data[i][10]]; !ok && len(data[i][10]) != 0 {
// 			zap.L().Error("invalid compliance catalog_id", zap.String("catalog_id", model.ComplianceKeyMlps2+" - "+data[i][10]))
// 			return errors.New("invalid compliance catalog_id")
// 		}
// 		if _, ok := cmap[model.ComplianceKeyMlps3+" - "+data[i][11]]; !ok && len(data[i][11]) != 0 {
// 			zap.L().Error("invalid compliance catalog_id", zap.String("catalog_id", model.ComplianceKeyMlps3+" - "+data[i][11]))
// 			return errors.New("invalid compliance catalog_id")
// 		}
//
// 		hashID := utils.HashPolicyID(string(platform), data[i][4], key)
// 		_, ok := pmap[hashID]
// 		if ok {
// 			zap.L().Error("duplicated policy key", zap.String("platform", string(platform)), zap.String("service", data[i][4]), zap.String("key", key))
// 			continue
// 		}
// 		pmap[hashID] = struct{}{}
//
// 		p := model.BenchmarkPolicy{
// 			ID:            hashID,
// 			Key:           key,
// 			Title:         data[i][1],
// 			TitleEn:       data[i][2],
// 			Desc:          utils.TrimLineBreaks(data[i][5]),
// 			DescEn:        utils.TrimLineBreaks(data[i][6]),
// 			CloudPlatform: platform,
// 			Service:       data[i][4],
// 			Mitigation:    utils.TrimLineBreaks(data[i][7]),
// 			MitigationEn:  utils.TrimLineBreaks(data[i][8]),
// 			Dengbao:       data[i][9],
// 			Mlps2:         data[i][10],
// 			Mlps3:         data[i][11],
// 			Severity:      model.SeverityConvert(data[i][12]),
// 			Apis:          apis,
// 			References:    references,
// 		}
// 		if _, ok := dbpmap[hashID]; ok {
// 			// 已经存在的规则，更新
// 			err = db.WithContext(ctx).Where("id = ?", hashID).Updates(p).Error
// 			if err != nil {
// 				zap.L().Info("update policy fails", zap.Error(err))
// 				continue
// 			}
// 		} else {
// 			// 新增的规则，插入
// 			err = db.WithContext(ctx).Create(p).Error
// 			if err != nil {
// 				zap.L().Info("create policy fails", zap.Error(err))
// 				continue
// 			}
// 		}
// 		// 遍历过后删除，遍历完成后map里还留下的policy就是被删除的policy，统一软删除
// 		delete(dbpmap, hashID)
// 	}
// 	// 统一软删除
// 	deletePids := make([]string, 0)
// 	for pid := range dbpmap {
// 		deletePids = append(deletePids, pid)
// 	}
//
// 	err = db.WithContext(ctx).Where("id in ?", deletePids).Delete(&model.BenchmarkPolicy{}).Error
// 	if err != nil {
// 		zap.L().Info("soft delete policies fails", zap.Error(err))
// 		return err
// 	}
//
// 	return nil
// }
//
// func writeRegions(ctx context.Context, db *gorm.DB, data [][]string) error {
// 	var err error
// 	err = db.WithContext(ctx).Exec("TRUNCATE TABLE " + model.Region{}.TableName()).Error
// 	if err != nil {
// 		zap.L().Info("truncate region table fails", zap.Error(err))
// 		return err
// 	}
//
// 	err = db.WithContext(ctx).Exec("TRUNCATE TABLE " + model.ServiceRegion{}.TableName()).Error
// 	if err != nil {
// 		zap.L().Info("truncate service_region table fails", zap.Error(err))
// 		return err
// 	}
//
// 	regionsM := make(map[string]struct{})
// 	regions := make([]model.Region, 0)
// 	serviceRegions := make([]*model.ServiceRegion, 0)
// 	for i := 1; i < len(data); i++ {
// 		key := data[i][2]
// 		platform := model.CloudPlatform(data[i][0])
// 		if key == "" || platform == "" {
// 			continue
// 		}
//
// 		service := data[i][1]
// 		regionID := utils.HashRegionID(string(platform), key)
// 		serviceID := utils.HashServiceID(string(platform), service)
//
// 		serviceRegions = append(serviceRegions, &model.ServiceRegion{
// 			CloudPlatform: platform,
// 			Service:       model.ServiceType(service),
// 			ServiceID:     serviceID,
// 			Region:        model.RegionType(key),
// 			RegionID:      regionID,
// 		})
//
// 		if _, ok := regionsM[string(platform)+"$"+key]; ok {
// 			continue
// 		} else {
// 			regionsM[string(platform)+"$"+key] = struct{}{}
// 		}
// 		regions = append(regions, model.Region{
// 			ID:            regionID,
// 			CloudPlatform: platform,
// 			Region:        model.RegionType(key),
// 			Name:          data[i][3],
// 			NameEn:        data[i][4],
// 			Remark:        data[i][5],
// 		})
// 	}
// 	err = db.WithContext(ctx).CreateInBatches(regions, 100).Error
// 	if err != nil {
// 		zap.L().Info("create regions fails", zap.Error(err))
// 		return err
// 	}
//
// 	err = db.WithContext(ctx).CreateInBatches(serviceRegions, 100).Error
// 	if err != nil {
// 		zap.L().Info("create service_regions fails", zap.Error(err))
// 		return err
// 	}
//
// 	return nil
// }
//
// func writeServices(ctx context.Context, db *gorm.DB, data [][]string) error {
// 	var err error
// 	err = db.WithContext(ctx).Exec("TRUNCATE TABLE " + model.Service{}.TableName()).Error
// 	if err != nil {
// 		zap.L().Info("truncate service tabled fails", zap.Error(err))
// 		return err
// 	}
//
// 	services := make([]model.Service, 0)
// 	for i := 1; i < len(data); i++ {
// 		key := data[i][0]
// 		platform := model.CloudPlatform(data[i][1])
// 		if key == "" || platform == "" {
// 			continue
// 		}
// 		services = append(services, model.Service{
// 			ID:            utils.HashServiceID(string(platform), key),
// 			CloudPlatform: platform,
// 			Service:       model.ServiceType(key),
// 			Name:          data[i][2],
// 			NameEn:        data[i][3],
// 		})
//
// 		// 构造service名称的映射数据，第一版现在内存里做，后面看要不要挪到数据库上
// 		model.CSN2TSN[data[i][3]] = []string{data[i][4], data[i][5]}
// 		if hids, ok := model.TSN2CSHIDs[data[i][5]]; ok {
// 			model.TSN2CSHIDs[data[i][5]] = append(hids, utils.HashServiceID(string(platform), key))
// 		} else {
// 			model.TSN2CSHIDs[data[i][5]] = []string{utils.HashServiceID(string(platform), key)}
// 		}
// 	}
// 	err = db.WithContext(ctx).CreateInBatches(services, 100).Error
// 	if err != nil {
// 		zap.L().Info("create services fails", zap.Error(err))
// 		return err
// 	}
// 	return nil
// }
//
// func writeAiPromptTemplates(ctx context.Context, db *gorm.DB, data [][]string) error {
// 	var err error
// 	err = db.WithContext(ctx).Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&model.PromptTemplate{}).Error
// 	if err != nil {
// 		zap.L().Info("delete prompt template fails", zap.Error(err))
// 		return err
// 	}
//
// 	err = db.WithContext(ctx).Exec("ALTER TABLE `prompt_templates` AUTO_INCREMENT = 1;").Error
// 	if err != nil {
// 		zap.L().Info("reset prompt_templates auto_increment fails", zap.Error(err))
// 		return err
// 	}
//
// 	promptTemplate := make([]model.PromptTemplate, 0)
// 	for i := 1; i < len(data); i++ {
// 		promptTemplate = append(promptTemplate, model.PromptTemplate{
// 			Key:          data[i][0],
// 			Confirm:      data[i][1],
// 			ConfirmEn:    data[i][2],
// 			UIQuestion:   data[i][3],
// 			UIQuestionEn: data[i][4],
// 			Prompt:       data[i][5],
// 			PromptEn:     data[i][6],
// 			CreatedAt:    time.Now().UnixMilli(),
// 		})
// 	}
// 	err = db.WithContext(ctx).CreateInBatches(promptTemplate, 100).Error
// 	if err != nil {
// 		zap.L().Info("create prompt_templates fails", zap.Error(err))
// 		return err
// 	}
// 	return nil
// }
//
// func rebuildByPolicies(ctx context.Context, db *gorm.DB) error {
// 	// 查询被删除的policy
// 	policies := make([]model.BenchmarkPolicy, 0)
// 	err := db.WithContext(ctx).Unscoped().Find(&policies).Error
// 	if err != nil {
// 		zap.L().Error("find deleted policies fails", zap.Error(err))
// 		return err
// 	}
//
// 	workingpmap := make(map[string]struct{})
// 	deletedPids := make([]string, 0)
// 	for i := range policies {
// 		if policies[i].Deleted.Valid {
// 			deletedPids = append(deletedPids, policies[i].ID)
// 		} else {
// 			workingpmap[policies[i].ID] = struct{}{}
// 		}
// 	}
//
// 	// 删除失效policy的risk
// 	// 再过一遍risks，确保没有遗漏的
// 	risks := make([]model.Risk, 0)
// 	err = db.WithContext(ctx).Where("1 = 1").Distinct("policy_id").Find(&risks).Error
// 	if err != nil {
// 		zap.L().Error("distinct risks fails", zap.Error(err))
// 		return err
// 	}
// 	rpids := make([]string, 0)
// 	for i := range risks {
// 		if _, ok := workingpmap[risks[i].PolicyID]; !ok {
// 			rpids = append(rpids, risks[i].PolicyID)
// 		}
// 	}
// 	err = db.WithContext(ctx).Where("policy_id in ?", append(deletedPids, rpids...)).Delete(&model.Risk{}).Error
// 	if err != nil {
// 		zap.L().Error("delete risks with deleted policies fails", zap.Error(err))
// 		return err
// 	}
//
// 	// 更新资产的风险统计
// 	type groupResult struct {
// 		AssetHashID string `gorm:"column:assets_hash_id"`
// 		Severity    int    `json:"column:severity"`
// 		Count       int    `json:"column:count"`
// 	}
// 	results := make([]groupResult, 0)
// 	err = db.WithContext(ctx).Where("1 = 1").Model(&model.Risk{}).
// 		Select("assets_hash_id, severity, COUNT(*) as count").
// 		Group("assets_hash_id, severity").
// 		Find(&results).Error
// 	if err != nil {
// 		zap.L().Error("group risks by asset & severity fails", zap.Error(err))
// 		return err
// 	}
//
// 	ramap := make(map[string]struct{})
// 	for i := range results {
// 		ramap[results[i].AssetHashID] = struct{}{}
// 		severityKey := ""
// 		if results[i].Severity == int(model.BenchmarkSeverityHigh) {
// 			severityKey = "severity_high"
// 		} else if results[i].Severity == int(model.BenchmarkSeverityMedium) {
// 			severityKey = "severity_medium"
// 		} else if results[i].Severity == int(model.BenchmarkSeverityLow) {
// 			severityKey = "severity_low"
// 		} else {
// 			zap.L().Error("invalid severity_level", zap.Int("severity_level", results[i].Severity))
// 			continue
// 		}
// 		err = db.WithContext(ctx).Model(&model.Assets{}).Omit("updated_at", "updater").
// 			Where("hash_id = ?", results[i].AssetHashID).Update(severityKey, results[i].Count).Error
// 		if err != nil {
// 			zap.L().Error("update asset severity fails", zap.Error(err), zap.Int("severity_level", results[i].Severity))
// 			continue
// 		}
// 	}
//
// 	assets := make([]model.Assets, 0)
// 	err = db.WithContext(ctx).Where("1 = 1").Find(&assets).Error
// 	if err != nil {
// 		zap.L().Error("find assets fails", zap.Error(err))
// 		return err
// 	}
// 	clearAids := make([]string, 0)
// 	for i := range assets {
// 		if _, ok := ramap[assets[i].HashID]; !ok {
// 			clearAids = append(clearAids, assets[i].HashID)
// 		}
// 	}
// 	err = db.WithContext(ctx).Model(&model.Assets{}).Omit("updated_at", "updater").
// 		Where("hash_id in ?", clearAids).
// 		Updates(map[string]interface{}{"severity_high": 0, "severity_medium": 0, "severity_low": 0}).Error
// 	if err != nil {
// 		zap.L().Error("clear assets severity count fails", zap.Error(err))
// 		return err
// 	}
//
// 	return nil
// }
//
// func initDefaultBenchmark(ctx context.Context, db *gorm.DB) error {
// 	var count int64
// 	err := db.WithContext(ctx).Table(model.Benchmark{}.TableName()).Where("is_default = 1").Count(&count).Error
// 	if err != nil {
// 		zap.L().Error("count default benchmark fails", zap.Error(err))
// 		return err
// 	}
//
// 	policies := make([]model.BenchmarkPolicy, 0)
// 	err = db.WithContext(ctx).Find(&policies).Error
// 	if err != nil {
// 		zap.L().Error("find all policies fails", zap.Error(err))
// 		return err
// 	}
// 	pids := make([]string, len(policies))
// 	for i := range policies {
// 		pids[i] = policies[i].ID
// 	}
//
// 	if count == 0 {
// 		now := time.Now().UnixMilli()
// 		err = db.WithContext(ctx).Create(&model.Benchmark{
// 			Name:        model.BenchmarkDefaultNameEN,
// 			Description: "",
// 			IsDefault:   true,
// 			Platforms:   model.ALLCloudPlatforms,
// 			Policies:    pids,
// 			CreatedAt:   now,
// 			Creator:     "",
// 			UpdatedAt:   now,
// 			Updater:     "",
// 		}).Error
// 		if err != nil {
// 			zap.L().Error("create default benchmark fails", zap.Error(err))
// 			return err
// 		}
// 	} else {
// 		spids := ""
// 		for i := range pids {
// 			spids = fmt.Sprintf("%s,\"%s\"", spids, pids[i])
// 		}
// 		spids = "[" + spids[1:] + "]"
// 		sps := ""
// 		for i := range model.ALLCloudPlatforms {
// 			sps = fmt.Sprintf("%s,\"%s\"", sps, model.ALLCloudPlatforms[i])
// 		}
// 		sps = "[" + sps[1:] + "]"
// 		err = db.WithContext(ctx).Model(&model.Benchmark{}).Where("is_default = 1").Updates(map[string]interface{}{
// 			"policies":  spids,
// 			"platforms": sps,
// 		}).Error
// 		if err != nil {
// 			zap.L().Error("update default benchmark fails", zap.Error(err))
// 			return err
// 		}
// 	}
//
// 	return nil
// }
//
// // initRiskScanTask 将未完成扫描的任务都设置为失败
// func initRiskScanTask(ctx context.Context, db *gorm.DB) error {
// 	err := db.WithContext(ctx).Model(&model.RiskScanTask{}).
// 		Where("`status` = ?", model.RiskScanTaskStatusRunning).
// 		Updates(map[string]any{
// 			"status":      model.RiskScanTaskStatusFailed,
// 			"finished_at": time.Now().UnixMilli(),
// 			"reason":      "restart fail",
// 		}).Error
// 	if err != nil {
// 		zap.L().Error("", zap.Error(err))
// 	}
//
// 	return err
// }
//
// // initEncryption 转换以前的老数据，下个版本移除
// func initEncryption(ctx context.Context, db *gorm.DB) error {
// 	type item struct {
// 		ID       int64
// 		Platform model.CloudPlatform
// 		Auth     []byte `gorm:"column:secret"`
// 	}
//
// 	type aliCredential struct {
// 		AccessKeyId     string `json:"access_key_id"`
// 		AccessKeySecret string `json:"access_key_secret"`
// 		Token           string `json:"token"`
// 	}
//
// 	type awsCredential struct {
// 		AccessKeyId     string `json:"access_key_id"`
// 		SecretAccessKey string `json:"secret_access_key"`
// 		Token           string `json:"token"`
// 	}
//
// 	type huaweiCredential struct {
// 		AccessKey string `json:"access_key"`
// 		SecretKey string `json:"secret_key"`
// 		Token     string `json:"token"`
// 	}
//
// 	type tencentCredential struct {
// 		SecretId  string `json:"secret_id"`
// 		SecretKey string `json:"secret_key"`
// 		Token     string `json:"token"`
// 	}
//
// 	list := make([]*item, 0)
//
// 	err := db.WithContext(ctx).Model(&model.Credential{}).
// 		Find(&list).Error
// 	if err != nil {
// 		zap.L().Error("", zap.Error(err))
// 		panic(err)
// 	}
//
// 	for _, v := range list {
// 		if v.Auth[0] == '{' {
// 			auth := model.CredentialAuth{}
// 			if v.Platform == model.CloudPlatformAli {
// 				c := aliCredential{}
// 				if err = json.Unmarshal(v.Auth, &c); err != nil {
// 					panic(err)
// 				}
//
// 				auth.Access = c.AccessKeyId
// 				auth.Secret = c.AccessKeySecret
// 				auth.Token = c.Token
// 			} else if v.Platform == model.CloudPlatformAWS {
// 				c := awsCredential{}
// 				if err = json.Unmarshal(v.Auth, &c); err != nil {
// 					panic(err)
// 				}
//
// 				auth.Access = c.AccessKeyId
// 				auth.Secret = c.SecretAccessKey
// 				auth.Token = c.Token
// 			} else if v.Platform == model.CloudPlatformHuawei {
// 				c := huaweiCredential{}
// 				if err = json.Unmarshal(v.Auth, &c); err != nil {
// 					panic(err)
// 				}
//
// 				auth.Access = c.AccessKey
// 				auth.Secret = c.SecretKey
// 				auth.Token = c.Token
// 			} else if v.Platform == model.CloudPlatformTencent {
// 				c := tencentCredential{}
// 				if err = json.Unmarshal(v.Auth, &c); err != nil {
// 					panic(err)
// 				}
//
// 				auth.Access = c.SecretId
// 				auth.Secret = c.SecretKey
// 				auth.Token = c.Token
// 			} else {
// 				panic("platform error" + v.Platform.String())
// 			}
//
// 			err = db.Model(&model.Credential{}).Where("id = ?", v.ID).
// 				UpdateColumn("secret", auth).Error
// 			if err != nil {
// 				panic(err)
// 			}
// 		}
// 	}
//
// 	return nil
// }
