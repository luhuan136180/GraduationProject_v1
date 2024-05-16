package contract

import (
	"context"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
	"v1/pkg/dao"
	"v1/pkg/model"
	"v1/pkg/utils"
)

const (
	SaveProjectType = 1
	SaveResumeType  = 2
	SaveInterviewee = 3
)

func CreateKey(savetype int) string {
	return utils.MD5Hex(time.Now().String() + string(savetype))
}

func getSaveTypeStr(saveType int) string {
	switch saveType {
	case SaveProjectType:
		return "项目"
	case SaveResumeType:
		return "简历"
	case SaveInterviewee:
		return "面试记录"
	}
	return ""
}

func (client *Client) SaveProject(ctx context.Context, db *gorm.DB) error {
	// 获取 未上链的列表
	projects, err := dao.FindProjectsUnContract(ctx, db)
	if err != nil {
		zap.L().Error("dao.FindProjectsUnContract ERR:", zap.Error(err))
		return err
	}

	if len(projects) == 0 {
		zap.L().Info("no new project ...")
		return nil
	}

	// 统计
	key := CreateKey(SaveProjectType)
	hashIDMap := make(map[int64]string) // id ； hash
	hashs := make([]string, 0)

	for i, project := range projects {
		project.Contract = model.Contract{} // 必须
		ahsStr, err := utils.CreeateSHA256ForProject(project)
		if err != nil {
			zap.L().Error("utils.CreeateSHA256ForProject ERR:", zap.Error(err))
			return err
		}

		hashIDMap[project.ID] = ahsStr
		hashs = append(hashs, ahsStr)
		projects[i].Contract = model.Contract{
			Flag:           true,
			ContractHashID: ahsStr,
			ContractKeyID:  key,
		}
	}

	//
	Transaction, err := client.SaveValue(key, hashs, SaveProjectType)
	if err != nil {
		zap.L().Error("client.SaveValue", zap.Error(err))
		return err
	}
	hash := Transaction.Hash()
	hashStr := hash.Hex()

	// 存数据库
	for _, project := range projects {
		project.BlockHash = hashStr
		err = db.WithContext(ctx).Model(&model.Project{}).Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			DoUpdates: clause.AssignmentColumns([]string{"flag", "contract_hash_id", "contract_key_id", "block_hash"}),
		}).Create(&project).Error
		if err != nil {
			zap.L().Error("update project.Contract config err:", zap.Error(err))
			return err
		}
	}

	_, err = dao.InsertBlockSaveLog(ctx, db, model.BlockSaveLog{
		BlockTXHash: hashStr,
		SaveType:    getSaveTypeStr(SaveProjectType),
		KeyHash:     key,
		Date:        Transaction.Time().String(),
	})
	if err != nil {
		zap.L().Error("dao.InsertBlockSaveLog", zap.Error(err))
		return err
	}

	return nil
}

func (client *Client) SaveResume(ctx context.Context, db *gorm.DB) error {
	// 获取 未上链的列表
	resumes, err := dao.FindResumesUnContract(ctx, db)
	if err != nil {
		zap.L().Error("dao.FindProjectsUnContract ERR:", zap.Error(err))
		return err
	}

	if len(resumes) == 0 {
		zap.L().Info("no new project ...")
		return nil
	}

	// 统计
	key := CreateKey(SaveResumeType)
	hashIDMap := make(map[int64]string) // id ； hash
	hashs := make([]string, 0)

	for i, resume := range resumes {
		resume.Contract = model.Contract{} // 必须
		ahsStr, err := utils.CreeateSHA256ForResume(resume)
		if err != nil {
			zap.L().Error("utils.CreeateSHA256ForProject ERR:", zap.Error(err))
			return err
		}

		hashIDMap[resume.ID] = ahsStr
		hashs = append(hashs, ahsStr)
		resumes[i].Contract = model.Contract{
			Flag:           true,
			ContractHashID: ahsStr,
			ContractKeyID:  key,
		}
	}

	//
	Transaction, err := client.SaveValue(key, hashs, SaveResumeType)
	if err != nil {
		zap.L().Error("client.SaveValue", zap.Error(err))
		return err
	}

	hash := Transaction.Hash()
	hashStr := hash.Hex()

	// 存数据库
	for _, resume := range resumes {
		resume.BlockHash = hashStr
		err = db.WithContext(ctx).Model(&model.Resume{}).Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			DoUpdates: clause.AssignmentColumns([]string{"flag", "contract_hash_id", "contract_key_id", "block_hash"}),
		}).Create(&resume).Error
		if err != nil {
			zap.L().Error("update project.Contract config err:", zap.Error(err))
			return err
		}
	}

	_, err = dao.InsertBlockSaveLog(ctx, db, model.BlockSaveLog{
		BlockTXHash: hashStr,
		SaveType:    getSaveTypeStr(SaveResumeType),
		KeyHash:     key,
		Date:        Transaction.Time().String(),
	})
	if err != nil {
		zap.L().Error("dao.InsertBlockSaveLog", zap.Error(err))
		return err
	}

	return nil
}

func (client *Client) SaveInterview(ctx context.Context, db *gorm.DB) error {
	// 获取 未上链的列表
	interviews, err := dao.FindInterviewsUnContract(ctx, db)
	if err != nil {
		zap.L().Error("dao.FindProjectsUnContract ERR:", zap.Error(err))
		return err
	}

	if len(interviews) == 0 {
		zap.L().Info("no new project ...")
		return nil
	}

	// 统计
	key := CreateKey(SaveInterviewee)
	hashIDMap := make(map[int64]string) // id ； hash
	hashs := make([]string, 0)

	for i, interview := range interviews {
		interview.Contract = model.Contract{} // 必须
		ahsStr, err := utils.CreeateSHA256ForInterview(interview)
		if err != nil {
			zap.L().Error("utils.CreeateSHA256ForProject ERR:", zap.Error(err))
			return err
		}

		hashIDMap[interview.ID] = ahsStr
		hashs = append(hashs, ahsStr)
		interviews[i].Contract = model.Contract{
			Flag:           true,
			ContractHashID: ahsStr,
			ContractKeyID:  key,
		}
	}

	//
	Transaction, err := client.SaveValue(key, hashs, SaveInterviewee)
	if err != nil {
		zap.L().Error("client.SaveValue", zap.Error(err))
		return err
	}

	hash := Transaction.Hash()
	hashStr := hash.Hex()

	// 存数据库
	for _, interview := range interviews {
		interview.BlockHash = hashStr
		err = db.WithContext(ctx).Model(&model.Interview{}).Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			DoUpdates: clause.AssignmentColumns([]string{"flag", "contract_hash_id", "contract_key_id", "block_hash"}),
		}).Create(&interview).Error
		if err != nil {
			zap.L().Error("update project.Contract config err:", zap.Error(err))
			return err
		}
	}

	_, err = dao.InsertBlockSaveLog(ctx, db, model.BlockSaveLog{
		BlockTXHash: hashStr,
		SaveType:    getSaveTypeStr(SaveInterviewee),
		KeyHash:     key,
		Date:        Transaction.Time().String(),
	})
	if err != nil {
		zap.L().Error("dao.InsertBlockSaveLog", zap.Error(err))
		return err
	}

	return nil
}

func (client *Client) Save(ctx context.Context, db *gorm.DB) error {
	err := client.SaveProject(ctx, db)
	if err != nil {
		zap.L().Error("save project failed:", zap.Error(err))
		return err
	}
	//
	err = client.SaveResume(ctx, db)
	if err != nil {
		zap.L().Error("save project failed:", zap.Error(err))
		return err
	}

	err = client.SaveInterview(ctx, db)
	if err != nil {
		zap.L().Error("save project failed:", zap.Error(err))
		return err
	}
	return nil
}
