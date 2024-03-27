package apiserver

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"errors"
	"strconv"
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
	errs = append(errs, initConfig(ctx, s.RDBClient))
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

func initConfig(ctx context.Context, db *gorm.DB) error {
	// configs/college.csv
	data, err := readCsv("./configs/college.csv")
	if err != nil {
		zap.L().Fatal("read file compliance fails", zap.Error(err))
	}
	cmap, err := writeColleges(ctx, db, data)
	if err != nil {
		zap.L().Fatal("writeCompliances fails", zap.Error(err))
	}

	// configs/profession.csv
	data, err = readCsv("./configs/profession.csv")
	if err != nil {
		zap.L().Fatal("read file policy fails", zap.Error(err))
	}
	err = writeProfession(ctx, db, data, cmap)
	if err != nil {
		zap.L().Fatal("writePolicies fails", zap.Error(err))
	}

	// config/class.csv
	data, err = readCsv("./configs/class.csv")
	if err != nil {
		zap.L().Fatal("read file policy fails", zap.Error(err))
	}
	err = writeClass(ctx, db, data, cmap)
	if err != nil {
		zap.L().Fatal("writePolicies fails", zap.Error(err))
	}

	return nil
}

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

func writeColleges(ctx context.Context, db *gorm.DB, data [][]string) (map[string]struct{}, error) {
	CollegeCatalogM := make(map[string]struct{})
	var err error
	err = db.WithContext(ctx).Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&model.College{}).Error
	if err != nil {
		zap.L().Info("delete compliances fails", zap.Error(err))
		return CollegeCatalogM, err
	}

	err = db.WithContext(ctx).Exec("ALTER TABLE `colleges` AUTO_INCREMENT = 1;").Error
	if err != nil {
		zap.L().Info("reset compliances auto_increment fails", zap.Error(err))
		return CollegeCatalogM, err
	}

	colleges := make([]model.College, 0)
	for i := 1; i < len(data); i++ {

		college, _ := json.Marshal(model.CollegeInfo{
			Info: data[i][1],
		})
		colleges = append(colleges, model.College{
			HashID:      utils.HashCollegeID(data[i][0]),
			CollegeName: data[i][0],
			CollegeInfo: string(college),
			Creator:     model.SystemUsername,
			CreatedAt:   time.Now().UnixMilli(),
			Updater:     model.SystemUsername,
			UpdatedAt:   time.Now().UnixMilli(),
		})
		CollegeCatalogM[data[i][0]] = struct{}{}
	}
	err = db.WithContext(ctx).CreateInBatches(colleges, 100).Error
	if err != nil {
		zap.L().Info("create compliances fails", zap.Error(err))
		return CollegeCatalogM, err
	}
	return CollegeCatalogM, nil
}

func writeProfession(ctx context.Context, db *gorm.DB, data [][]string, cmap map[string]struct{}) error {
	var err error

	dbProfessions := make([]model.Profession, 0)
	err = db.WithContext(ctx).Unscoped().Where("1 = 1").Find(&dbProfessions).Error
	if err != nil {
		zap.L().Info("find db policies fails", zap.Error(err))
		return err
	}
	dbpmap := make(map[string]model.Profession)
	for i := range dbProfessions {
		dbpmap[dbProfessions[i].HashID] = dbProfessions[i]
	}

	pmap := make(map[string]struct{})
	for i := 1; i < len(data); i++ {
		if _, ok := cmap[data[i][0]]; !ok && len(data[i][0]) != 0 {
			zap.L().Error("invalid profession college_name", zap.String("college_name : ", data[i][0]))
			return errors.New("invalid compliance catalog_id")
		}

		hashID := utils.HashProfessionID(utils.HashCollegeID(data[i][0]), data[i][1])
		_, ok := pmap[hashID]
		if ok {
			continue
		}
		pmap[hashID] = struct{}{}

		profession, _ := json.Marshal(model.ProfessionInfo{
			Info: data[i][2],
		})

		p := model.Profession{
			HashID:         hashID,
			CollegeHashID:  utils.HashCollegeID(data[i][0]),
			CollegeName:    data[i][0],
			ProfessionName: data[i][1],
			ProfessionInfo: string(profession),

			Creator:   model.SystemUsername,
			CreatedAt: time.Now().UnixMilli(),
			Updater:   model.SystemUsername,
			UpdatedAt: time.Now().UnixMilli(),
		}
		if _, ok := dbpmap[hashID]; ok {
			// 已经存在的规则，更新
			err = db.WithContext(ctx).Where("hash_id = ?", hashID).Updates(p).Error
			if err != nil {
				zap.L().Info("update profession fails", zap.Error(err))
				continue
			}
		} else {
			// 新增的规则，插入
			err = db.WithContext(ctx).Create(p).Error
			if err != nil {
				zap.L().Info("create policy fails", zap.Error(err))
				continue
			}
		}
		// 遍历过后删除，遍历完成后map里还留下的policy就是被删除的policy，统一软删除
		delete(dbpmap, hashID)
	}
	// 统一软删除
	deletePids := make([]string, 0)
	for pid := range dbpmap {
		deletePids = append(deletePids, pid)
	}

	err = db.WithContext(ctx).Where("hash_id in ?", deletePids).Delete(&model.Profession{}).Error
	if err != nil {
		zap.L().Info("soft delete policies fails", zap.Error(err))
		return err
	}

	return nil
}

func writeClass(ctx context.Context, db *gorm.DB, data [][]string, cmap map[string]struct{}) error {
	var err error

	dbClasses := make([]model.Class, 0)
	err = db.WithContext(ctx).Unscoped().Where("1 = 1").Find(&dbClasses).Error
	if err != nil {
		zap.L().Info("find db policies fails", zap.Error(err))
		return err
	}
	dbcmap := make(map[string]model.Class)
	for i := range dbClasses {
		dbcmap[dbClasses[i].ClassHashID] = dbClasses[i]
	}

	professionNames := make([]string, 0)
	for i := 1; i < len(data); i++ {
		professionNames = append(professionNames, data[i][0])
	}

	professions, err := dao.GetProfessionByProfessionName(ctx, db, professionNames)
	if err != nil {
		zap.L().Info("find profession by name fails", zap.Error(err))
		return err
	}

	pmap := make(map[string]model.Profession)
	for i, _ := range professions {
		pmap[professions[i].ProfessionName] = professions[i]
	}

	classmap := make(map[string]struct{})
	for i := 1; i < len(data); i++ {
		if _, ok := pmap[data[i][0]]; !ok && len(data[i][0]) != 0 {
			zap.L().Error("invalid profession college_name", zap.String("college_name : ", data[i][0]))
			return errors.New("invalid compliance catalog_id")
		}

		classID, _ := strconv.Atoi(data[i][2])
		hashID := utils.HashClassID(pmap[data[i][0]].HashID, data[i][1], classID)
		_, ok := classmap[hashID]
		if ok {
			continue
		}
		classmap[hashID] = struct{}{}

		p := model.Class{
			ProfessionHashID: pmap[data[i][0]].HashID,
			ClassHashID:      hashID,
			ClassName:        data[i][1],
			ClassID:          classID,

			Creator:   model.SystemUsername,
			CreatedAt: time.Now().UnixMilli(),
			Updater:   model.SystemUsername,
			UpdatedAt: time.Now().UnixMilli(),
		}
		if _, ok := dbcmap[hashID]; ok {
			// 已经存在的规则，更新
			err = db.WithContext(ctx).Where("class_hash_id = ?", hashID).Updates(p).Error
			if err != nil {
				zap.L().Info("update profession fails", zap.Error(err))
				continue
			}
		} else {
			// 新增的规则，插入
			err = db.WithContext(ctx).Create(p).Error
			if err != nil {
				zap.L().Info("create policy fails", zap.Error(err))
				continue
			}
		}
		// 遍历过后删除，遍历完成后map里还留下的policy就是被删除的policy，统一软删除
		delete(dbcmap, hashID)
	}
	// 统一软删除
	deletePids := make([]string, 0)
	for pid := range dbcmap {
		deletePids = append(deletePids, pid)
	}

	err = db.WithContext(ctx).Where("class_hash_id in ?", deletePids).Delete(&model.Profession{}).Error
	if err != nil {
		zap.L().Info("soft delete policies fails", zap.Error(err))
		return err
	}

	return nil
}
