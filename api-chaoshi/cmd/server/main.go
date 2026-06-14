package main

import (
	"log"
	"os"

	"chaoshi_api/internal/config"
	"chaoshi_api/internal/handlers/merchant"
	"chaoshi_api/internal/handlers/sp"
	"chaoshi_api/internal/handlers/upload"
	"chaoshi_api/internal/handlers/user"
	"chaoshi_api/internal/middleware"
	"chaoshi_api/internal/storage"
	"chaoshi_api/pkg/database"

	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化配置
	if err := config.InitConfig(); err != nil {
		log.Fatalf("配置初始化失败: %v", err)
	}
	// 初始化数据库
	if err := database.InitDB(&config.Config.Database); err != nil {
		log.Fatalf("数据库初始化失败: %v", err)
	}
	defer database.CloseDB()

	// 初始化本地文件存储
	if err := storage.InitLocalStorage(); err != nil {
		log.Fatalf("本地文件存储初始化失败: %v", err)
	}

	// 设置Gin模式
	if config.Config.App.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	// 创建Gin实例
	r := gin.Default()

	// 全局中间件
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.CORS())
	r.Static(storage.GetService().URLPrefix(), storage.GetService().RootDir())

	// 路由设置
	setupRoutes(r)

	// 启动服务器
	addr := config.Config.App.GetAddr()
	log.Printf("服务启动中，监听地址: %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("服务启动失败: %v", err)
		os.Exit(1)
	}
}

// setupRoutes 设置路由
func setupRoutes(r *gin.Engine) {
	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// API v1 路由组
	v1 := r.Group("/api/v1")
	{
		// 认证相关
		authGroup := v1.Group("/auth")
		{
			authGroup.POST("/merchant/login", merchant.Login)
			authGroup.POST("/user/wechat-login", user.WechatLogin)
		}

		// 文件上传接口
		uploadHandler := upload.NewUploadHandler()
		uploadGroup := v1.Group("/upload")
		uploadGroup.Use(middleware.JWTAuth())
		{
			uploadGroup.POST("/file", uploadHandler.UploadFile)
		}

		// C端店铺接口
		storeGroup := v1.Group("/store/:merchant_id")
		storeGroup.Use(middleware.OptionalJWTAuth())
		{
			storeGroup.GET("/home", user.GetStoreHome)
			storeGroup.GET("/products", user.GetProducts)
			storeGroup.GET("/products/:product_id", user.GetProductDetail)
			storeGroup.GET("/delivery-rules", user.GetDeliveryRules)
			storeGroup.GET("/pickup-points", user.GetPickupPoints)
			storeGroup.GET("/full-reduction-rules", user.GetStoreFullReductionRules)
			storeGroup.POST("/visit", user.RecordUserVisit)
			storeGroup.POST("/event", user.RecordBehaviorEvent)

			storeAuthedGroup := storeGroup.Group("")
			storeAuthedGroup.Use(middleware.JWTAuth(), middleware.UserAuth())
			{
				storeAuthedGroup.POST("/orders", user.CreateOrder)
			}
		}

		// C端用户接口（需要登录）
		userGroup := v1.Group("/user")
		userGroup.Use(middleware.JWTAuth(), middleware.UserAuth())
		{
			userGroup.GET("/orders", user.GetOrders)
			userGroup.GET("/orders/:order_id", user.GetOrderDetail)
			userGroup.POST("/orders/:order_id/cancel", user.CancelOrder)
			userGroup.POST("/orders/:order_id/refund", user.ApplyRefund)
			userGroup.GET("/addresses", user.GetAddresses)
			userGroup.POST("/addresses", user.CreateAddress)
			userGroup.PUT("/addresses/:id", user.UpdateAddress)
			userGroup.DELETE("/addresses/:id", user.DeleteAddress)
		}

		// 商家管理员接口
		merchantGroup := v1.Group("/merchant")
		merchantGroup.Use(middleware.JWTAuth(), middleware.MerchantAuth())
		{
			// 商家信息
			merchantGroup.GET("/profile", merchant.GetProfile)
			merchantGroup.PUT("/profile", merchant.UpdateProfile)
			merchantGroup.GET("/settings", merchant.GetSettings)
			merchantGroup.PUT("/settings", merchant.UpdateSettings)
			merchantGroup.POST("/account/change-password", merchant.ChangePassword)
			merchantGroup.POST("/status", merchant.UpdateStatus)
			merchantGroup.GET("/qrcode", merchant.GetQRCode)
			merchantGroup.GET("/delivery-settings", merchant.GetDeliverySettings)
			merchantGroup.PUT("/delivery-settings", merchant.UpdateDeliverySettings)
			merchantGroup.GET("/pickup-points", merchant.GetPickupPoints)
			merchantGroup.POST("/pickup-points", merchant.CreatePickupPoint)
			merchantGroup.PUT("/pickup-points/:id", merchant.UpdatePickupPoint)
			merchantGroup.DELETE("/pickup-points/:id", merchant.DeletePickupPoint)
			merchantGroup.GET("/full-reduction-rules", merchant.GetFullReductionRules)
			merchantGroup.PUT("/full-reduction-rules", merchant.UpdateFullReductionRules)
			merchantGroup.GET("/subscriptions", merchant.GetSubscriptions)
			merchantGroup.PUT("/subscriptions", merchant.UpdateSubscriptions)

			// 员工管理
			merchantGroup.GET("/staff", merchant.GetStaffList)
			merchantGroup.POST("/staff", merchant.CreateStaff)
			merchantGroup.PUT("/staff/:id", merchant.UpdateStaff)
			merchantGroup.DELETE("/staff/:id", merchant.DeleteStaff)
			merchantGroup.POST("/staff/:id/reset-password", merchant.ResetStaffPassword)

			// 系统公告（商家查看）
			merchantGroup.GET("/announcements", merchant.GetAnnouncements)
			merchantGroup.GET("/announcements/:id", merchant.GetAnnouncementDetail)

			// 商品分类
			merchantGroup.GET("/categories", merchant.GetCategories)
			merchantGroup.POST("/categories", merchant.CreateCategory)
			merchantGroup.PUT("/categories/:category_id", merchant.UpdateCategory)
			merchantGroup.DELETE("/categories/:category_id", merchant.DeleteCategory)
			merchantGroup.POST("/categories/sort", merchant.SortCategories)

			// 商品管理
			merchantGroup.GET("/products", merchant.GetProducts)
			merchantGroup.GET("/products/:product_id", merchant.GetProduct)
			merchantGroup.POST("/products", merchant.CreateProduct)
			merchantGroup.PUT("/products/:product_id", merchant.UpdateProduct)
			merchantGroup.POST("/products/:product_id/on-sale", merchant.ProductOnSale)
			merchantGroup.POST("/products/:product_id/off-sale", merchant.ProductOffSale)
			merchantGroup.POST("/products/batch-status", merchant.BatchUpdateProductStatus)
			merchantGroup.DELETE("/products/:product_id", merchant.DeleteProduct)
			merchantGroup.PUT("/products/:product_id/stock", merchant.UpdateStock)
			merchantGroup.GET("/products/:product_id/specs", merchant.GetProductSpecs)
			merchantGroup.PUT("/products/:product_id/specs", merchant.UpdateProductSpecs)
			merchantGroup.DELETE("/products/:product_id/specs", merchant.DeleteProductSpecs)

			// 订单管理
			merchantGroup.GET("/orders", merchant.GetOrders)
			merchantGroup.GET("/orders/:order_id", merchant.GetOrderDetail)
			merchantGroup.POST("/orders/quick-complete", merchant.QuickCompleteOrder)
			merchantGroup.POST("/orders/:order_id/complete", merchant.CompleteOrder)
			merchantGroup.POST("/orders/:order_id/refund", merchant.RefundOrder)
			merchantGroup.GET("/orders/statistics", merchant.GetOrderStatistics)

			// 数据分析
			merchantGroup.GET("/analytics/overview", merchant.GetAnalyticsOverview)
			merchantGroup.GET("/analytics/sales-trend", merchant.GetSalesTrend)
			merchantGroup.GET("/analytics/product-ranking", merchant.GetProductRanking)
			merchantGroup.GET("/analytics/hourly", merchant.GetHourlyAnalysis)
			merchantGroup.GET("/analytics/stock-alert", merchant.GetStockAlert)
			merchantGroup.GET("/analytics/customers", merchant.GetCustomerAnalysis)
			merchantGroup.GET("/analytics/customer-trend", merchant.GetCustomerTrend)

			merchantGroup.GET("/printers", merchant.GetPrinters)
			merchantGroup.POST("/printers", merchant.CreatePrinter)
			merchantGroup.PUT("/printers/:printer_id", merchant.UpdatePrinter)
			merchantGroup.DELETE("/printers/:printer_id", merchant.DeletePrinter)
			merchantGroup.POST("/printers/:printer_id/test", merchant.TestPrinter)
			merchantGroup.GET("/print-logs", merchant.GetPrintLogs)
		}

		// 总部后台接口
		adminPublicGroup := v1.Group("/admin")
		{
			adminPublicGroup.POST("/auth/login", sp.Login)

			adminGroup := adminPublicGroup.Group("")
			adminGroup.Use(middleware.JWTAuth(), middleware.AdminAuth())
			{
				adminGroup.POST("/auth/logout", sp.Logout)
				adminGroup.GET("/dashboard", sp.GetDashboard)
				adminGroup.POST("/stores", sp.CreateMerchant)
				adminGroup.PUT("/stores/:merchant_id", sp.UpdateMerchant)
				adminGroup.GET("/stores/:merchant_id", sp.GetMerchantDetail)
				adminGroup.POST("/stores/:merchant_id/admin/reset-password", sp.ResetMerchantAdminPassword)
				adminGroup.PUT("/stores/:merchant_id/assets", sp.UpdateMerchantAssets)
				adminGroup.GET("/stores/analytics/distribution", sp.GetMerchantDistribution)
				adminGroup.GET("/stores/list", sp.GetMerchantList)
				adminGroup.GET("/stores/:merchant_id/qrcode", sp.GetMerchantQRCode)
				adminGroup.GET("/stores/:merchant_id/categories", sp.GetMerchantCategories)
				adminGroup.POST("/stores/:merchant_id/categories", sp.CreateMerchantCategory)
				adminGroup.PUT("/stores/:merchant_id/categories/:category_id", sp.UpdateMerchantCategory)
				adminGroup.DELETE("/stores/:merchant_id/categories/:category_id", sp.DeleteMerchantCategory)
				adminGroup.POST("/stores/:merchant_id/categories/sort", sp.SortMerchantCategories)
				adminGroup.GET("/stores/:merchant_id/products", sp.GetMerchantProducts)
				adminGroup.GET("/stores/:merchant_id/products/:product_id", sp.GetMerchantProduct)
				adminGroup.POST("/stores/:merchant_id/products", sp.CreateMerchantProduct)
				adminGroup.PUT("/stores/:merchant_id/products/:product_id", sp.UpdateMerchantProduct)
				adminGroup.POST("/stores/:merchant_id/products/:product_id/on-sale", sp.MerchantProductOnSale)
				adminGroup.POST("/stores/:merchant_id/products/:product_id/off-sale", sp.MerchantProductOffSale)
				adminGroup.POST("/stores/:merchant_id/products/batch-status", sp.BatchUpdateMerchantProductStatus)
				adminGroup.DELETE("/stores/:merchant_id/products/:product_id", sp.DeleteMerchantProduct)
				adminGroup.PUT("/stores/:merchant_id/products/:product_id/stock", sp.UpdateMerchantProductStock)
				adminGroup.GET("/stores/:merchant_id/products/:product_id/specs", sp.GetMerchantProductSpecs)
				adminGroup.PUT("/stores/:merchant_id/products/:product_id/specs", sp.UpdateMerchantProductSpecs)
				adminGroup.DELETE("/stores/:merchant_id/products/:product_id/specs", sp.DeleteMerchantProductSpecs)
				adminGroup.GET("/stores/:merchant_id/pickup-points", sp.GetMerchantPickupPoints)
				adminGroup.POST("/stores/:merchant_id/pickup-points", sp.CreateMerchantPickupPoint)
				adminGroup.PUT("/stores/:merchant_id/pickup-points/:id", sp.UpdateMerchantPickupPoint)
				adminGroup.DELETE("/stores/:merchant_id/pickup-points/:id", sp.DeleteMerchantPickupPoint)
				adminGroup.GET("/orders/analytics", sp.GetOrderAnalytics)
				adminGroup.GET("/orders", sp.GetOrders)
				adminGroup.GET("/orders/:order_id", sp.GetOrderDetail)
				adminGroup.GET("/amount/analytics", sp.GetAmountAnalytics)
				adminGroup.GET("/amount/top-stores", sp.GetTopMerchants)
				adminGroup.POST("/account/change-password", sp.ChangePassword)
				adminGroup.GET("/announcements", sp.GetAnnouncements)
				adminGroup.GET("/announcements/:id", sp.GetAnnouncementDetail)
				adminGroup.POST("/announcements", sp.CreateAnnouncement)
				adminGroup.PUT("/announcements/:id", sp.UpdateAnnouncement)
				adminGroup.DELETE("/announcements/:id", sp.DeleteAnnouncement)

				// 活动管理
				adminGroup.GET("/activities", sp.GetActivities)
				adminGroup.POST("/activities", sp.CreateActivity)
				adminGroup.PUT("/activities/:id", sp.UpdateActivity)
				adminGroup.DELETE("/activities/:id", sp.DeleteActivity)
			}
		}

	}
}
