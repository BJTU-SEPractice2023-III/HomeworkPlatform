package service_test

// func TestMain(m *testing.M) {
// 	server.InitRouter()
// 	bootstrap.Sqlite = true
// 	var err error
// 	models.DB, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
// 		Logger: logger.Default.LogMode(logger.Silent),
// 	})
// 	bootstrap.Sqlite = true
// 	if err != nil {
// 		panic(err)
// 	}
// 	models.DB.AutoMigrate(&models.User{})
// 	models.DB.AutoMigrate(&models.Course{})
// 	models.DB.AutoMigrate(&models.Homework{})
// 	models.DB.AutoMigrate(&models.HomeworkSubmission{})
// 	models.DB.AutoMigrate(&models.Comment{})
// 	models.DB.AutoMigrate(&models.Complaint{})
// 	os.Exit(m.Run())
// }
