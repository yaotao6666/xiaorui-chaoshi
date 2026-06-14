package database

import (
	"fmt"
	"log"
	"time"

	"chaoshi_api/internal/config"
	"chaoshi_api/internal/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// InitDB 初始化数据库连接
func InitDB(cfg *config.Database) error {
	// 配置日志
	newLogger := logger.Default.LogMode(logger.Info)

	// 连接数据库
	var err error
	DB, err = gorm.Open(mysql.Open(cfg.GetDSN()), &gorm.Config{
		Logger: newLogger,
	})

	if err != nil {
		return fmt.Errorf("数据库连接失败: %w", err)
	}

	if ensureErr := ensureMerchantStaffColumns(DB); ensureErr != nil {
		return fmt.Errorf("初始化商家员工扩展字段失败: %w", ensureErr)
	}

	if ensureErr := ensureUserBehaviorEventsTable(DB); ensureErr != nil {
		return fmt.Errorf("初始化用户行为事件表失败: %w", ensureErr)
	}

	if ensureErr := ensureUserTables(DB); ensureErr != nil {
		return fmt.Errorf("初始化用户表扩展字段失败: %w", ensureErr)
	}

	if ensureErr := ensureOrderColumns(DB); ensureErr != nil {
		return fmt.Errorf("初始化订单扩展字段失败: %w", ensureErr)
	}

	if ensureErr := ensureMerchantColumns(DB); ensureErr != nil {
		return fmt.Errorf("初始化商家扩展字段失败: %w", ensureErr)
	}

	if ensureErr := ensureMerchantBusinessTables(DB); ensureErr != nil {
		return fmt.Errorf("初始化商家业务表失败: %w", ensureErr)
	}

	// 获取底层 sql.DB
	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("获取数据库实例失败: %w", err)
	}

	// 设置连接池
	sqlDB.SetMaxIdleConns(cfg.MaxIdle)
	sqlDB.SetMaxOpenConns(cfg.MaxOpen)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.LifeTime) * time.Second)

	log.Println("数据库连接成功")
	return nil
}

func ensureMerchantStaffColumns(db *gorm.DB) error {
	if err := ensureMerchantStaffPushOpenID(db); err != nil {
		return err
	}

	if err := ensureMerchantStaffColumn(
		db,
		"notify_enabled",
		"ADD COLUMN notify_enabled TINYINT(1) NOT NULL DEFAULT 1 COMMENT '订单提示音开关' AFTER role",
	); err != nil {
		return err
	}

	if err := ensureMerchantStaffColumn(
		db,
		"browse_notify_enabled",
		"ADD COLUMN browse_notify_enabled TINYINT(1) NOT NULL DEFAULT 1 COMMENT '顾客浏览提示音开关' AFTER notify_enabled",
	); err != nil {
		return err
	}

	return ensureMerchantStaffIndex(
		db,
		"idx_merchant_staffs_push_openid",
		"CREATE INDEX idx_merchant_staffs_push_openid ON merchant_staffs (push_openid)",
	)
}

func ensureMerchantStaffPushOpenID(db *gorm.DB) error {
	pushOpenIDExists, err := tableColumnExists(db, "merchant_staffs", "push_openid")
	if err != nil {
		return err
	}
	if pushOpenIDExists {
		return nil
	}

	oldOpenIDExists, err := tableColumnExists(db, "merchant_staffs", "openid")
	if err != nil {
		return err
	}
	if oldOpenIDExists {
		return db.Exec(`
			ALTER TABLE merchant_staffs
			CHANGE COLUMN openid push_openid VARCHAR(64) DEFAULT NULL COMMENT '消息推送OpenID'
		`).Error
	}

	return ensureMerchantStaffColumn(
		db,
		"push_openid",
		"ADD COLUMN push_openid VARCHAR(64) DEFAULT NULL COMMENT '消息推送OpenID' AFTER phone",
	)
}

func ensureMerchantStaffColumn(db *gorm.DB, columnName string, addSQL string) error {
	exists, err := tableColumnExists(db, "merchant_staffs", columnName)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}

	return db.Exec("ALTER TABLE merchant_staffs " + addSQL).Error
}

func ensureMerchantStaffIndex(db *gorm.DB, indexName string, createSQL string) error {
	var count int64
	queryErr := db.Raw(`
		SELECT COUNT(*)
		FROM INFORMATION_SCHEMA.STATISTICS
		WHERE TABLE_SCHEMA = DATABASE()
		  AND TABLE_NAME = 'merchant_staffs'
		  AND INDEX_NAME = ?
	`, indexName).Scan(&count).Error
	if queryErr != nil {
		return queryErr
	}

	if count > 0 {
		return nil
	}

	return db.Exec(createSQL).Error
}

func tableColumnExists(db *gorm.DB, tableName string, columnName string) (bool, error) {
	var count int64
	queryErr := db.Raw(`
		SELECT COUNT(*)
		FROM INFORMATION_SCHEMA.COLUMNS
		WHERE TABLE_SCHEMA = DATABASE()
		  AND TABLE_NAME = ?
		  AND COLUMN_NAME = ?
	`, tableName, columnName).Scan(&count).Error
	if queryErr != nil {
		return false, queryErr
	}
	return count > 0, nil
}

func ensureUserBehaviorEventsTable(db *gorm.DB) error {
	return db.AutoMigrate(&models.UserBehaviorEvent{})
}

func ensureUserTables(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.User{},
		&models.UserVisit{},
	)
}

func ensureOrderColumns(db *gorm.DB) error {
	orderColumns := map[string]string{
		"completed_by_name":  "ADD COLUMN completed_by_name VARCHAR(64) DEFAULT NULL COMMENT '核销人' AFTER completed_at",
		"pay_notify_payload": "ADD COLUMN pay_notify_payload JSON DEFAULT NULL COMMENT '支付回调原始数据' AFTER paid_at",
	}

	for columnName, addSQL := range orderColumns {
		if err := ensureTableColumn(db, "orders", columnName, addSQL); err != nil {
			return err
		}
	}

	return nil
}

func ensureMerchantColumns(db *gorm.DB) error {
	type merchantColumn struct {
		name   string
		addSQL string
	}

	merchantColumns := []merchantColumn{
		{name: "cover_image", addSQL: "ADD COLUMN cover_image VARCHAR(512) DEFAULT NULL COMMENT '商家背景图' AFTER logo"},
		{name: "takeout_enabled", addSQL: "ADD COLUMN takeout_enabled TINYINT(1) NOT NULL DEFAULT 1 COMMENT '是否支持配送' AFTER min_order_amount"},
		{name: "dine_in_enabled", addSQL: "ADD COLUMN dine_in_enabled TINYINT(1) NOT NULL DEFAULT 1 COMMENT '是否支持堂食' AFTER takeout_enabled"},
		{name: "pickup_enabled", addSQL: "ADD COLUMN pickup_enabled TINYINT(1) NOT NULL DEFAULT 1 COMMENT '是否支持自提' AFTER dine_in_enabled"},
	}

	for _, column := range merchantColumns {
		if err := ensureTableColumn(db, "merchants", column.name, column.addSQL); err != nil {
			return err
		}
	}

	// 兼容历史库中已存在但值为空的场景，统一补成开启，避免下单方式被错误隐藏。
	if err := db.Exec(`
		UPDATE merchants
		SET takeout_enabled = COALESCE(takeout_enabled, 1),
			dine_in_enabled = COALESCE(dine_in_enabled, 1),
			pickup_enabled = COALESCE(pickup_enabled, 1)
	`).Error; err != nil {
		return err
	}

	return nil
}

func ensureMerchantBusinessTables(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.MerchantDeliverySettings{},
		&models.MerchantFullReductionRule{},
	)
}

func ensureTableColumn(db *gorm.DB, tableName, columnName, addSQL string) error {
	var count int64
	queryErr := db.Raw(`
		SELECT COUNT(*)
		FROM INFORMATION_SCHEMA.COLUMNS
		WHERE TABLE_SCHEMA = DATABASE()
		  AND TABLE_NAME = ?
		  AND COLUMN_NAME = ?
	`, tableName, columnName).Scan(&count).Error
	if queryErr != nil {
		return queryErr
	}

	if count > 0 {
		return nil
	}

	return db.Exec("ALTER TABLE " + tableName + " " + addSQL).Error
}

// GetDB 获取数据库实例
func GetDB() *gorm.DB {
	return DB
}

// CloseDB 关闭数据库连接
func CloseDB() error {
	if DB != nil {
		sqlDB, err := DB.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}
