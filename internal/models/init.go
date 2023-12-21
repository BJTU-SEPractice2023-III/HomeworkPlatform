package models

import (
	"fmt"
	"homework_platform/internal/bootstrap"
	"homework_platform/internal/utils"
	"log"
	"os"
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm/logger"

	// "gocloud.dev/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func sqliteDB(dsn string, config *gorm.Config) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(dsn+"?_pragma=foreign_keys(1)"), config)
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	sqlDB.SetMaxOpenConns(1)
	if err != nil {
		return nil, err
	}
	// if res := db.Exec("PRAGMA foreign_keys = ON", nil); res.Error != nil {
	// 	return nil, err
	// }

	return db, nil
}

func InitDB() {
	// db, err := gorm.Open(sqlite.Open("ach.db"), &gorm.Config{})
	// db, err := gorm.Open(mysql.Open(bootstrap.Config.SQLDSN), &gorm.Config{})
	var db *gorm.DB
	var err error
	if bootstrap.Sqlite {
		db, err = sqliteDB("homework-platform.db", &gorm.Config{})
	} else if bootstrap.SqliteInMemEmpty {
		db, err = sqliteDB(":memory:", &gorm.Config{})
	} else if bootstrap.Mysql {
		db, err = gorm.Open(mysql.Open(bootstrap.Config.SQLDSN), &gorm.Config{})
	} else if bootstrap.Test {
		db, err = sqliteDB(":memory:", &gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		})
	} else {
		db, err = gorm.Open(postgres.Open(bootstrap.Config.SQLDSN), &gorm.Config{})
	}

	if err != nil {
		log.Panicf("无法连接数据库，%s", err)
	}

	DB = db
	Migrate()

	if bootstrap.GenDataOverwrite {
		deleteData()
	}
	if bootstrap.GenData || bootstrap.GenDataOverwrite {
		generateData()
	}

	// 创建初始管理员账户
	// addDefaultUser()
}

func deleteData() {
	log.Println("正在删库🥳...")
	DB.Where("1 = 1").Delete(&User{})
	DB.Where("1 = 1").Delete(&File{})
	DB.Where("1 = 1").Delete(&Course{})
	DB.Where("1 = 1").Delete(&Homework{})
	DB.Where("1 = 1").Delete(&HomeworkSubmission{})
	DB.Where("1 = 1").Delete(&Comment{})
	DB.Where("1 = 1").Delete(&Complaint{})
}

func Migrate() {
	DB.AutoMigrate(&User{})
	DB.AutoMigrate(&File{})
	DB.AutoMigrate(&Course{})
	DB.AutoMigrate(&Homework{})
	DB.AutoMigrate(&HomeworkSubmission{})
	DB.AutoMigrate(&Comment{})
	DB.AutoMigrate(&Complaint{})
}

func generateData() {
	xyh, _ := CreateUser("xyh", "xyh")
	tjw, _ := CreateUser("tjw", "tjw")
	xb, _ := CreateUser("xb", "xb")
	xeh, _ := CreateUser("xeh", "xyh")
	xsh, _ := CreateUser("xsh", "xyh")

	arknights, _ := tjw.CreateCourse(
		"明日方舟",
		time.Date(2019, 5, 1, 3, 59, 0, 0, time.Local), time.Now().Add(24*time.Hour),
		`《明日方舟》（Arknights）是由鹰角网络所开发的一款策略塔防手机游戏，也是鹰角网络开发的第一款游戏。游戏在2017年9月8日发布概念预告，经历三次删档测试[5]，2019年5月1日全平台公开测试[6]。

		游戏玩法
		为二次元塔防抽卡类型游戏，有许多角色（“干员”）作为塔。许多干员拥有不同的数值、技能和攻击方式，因此玩家应在关卡开始前进行合理的编队（最多13名干员：12名玩家干员，一名助战干员）并对这些被选中的干员选择合理的技能以发挥最大效果。每名干员不仅性格各异而且各有所长。因此干员们共被分为八大职业：先锋、近卫、重装、狙击、术师、医疗、辅助和特种。每个职业还细分为多个分支。近战干员可以放置在地面上，远程干员可以放置在高台上。大部分近战干员能阻儅一定数量的敌人前进，远程干员造成远程伤害、治疗或以其他方式支持其他干员。玩家必须将干员在合理的时机以合理的顺序放置在合理的位置上，以防止种类繁多的敌人渗透到玩家的基地。

		由于同一关卡的通关方式有限，尤其是在较高难度下，明日方舟也被认为是一款益智游戏[7]。游玩时并不需要快速的反应力（游戏是可暂停的，进行操作时游戏时间可调整），而是需要现场的战术分析和预见[8]。

		该游戏还具有基地建设方面玩法，允许玩家建造设施生产资源并为其分配干员、布置家具。这允许玩家不在线时也能获取资源[9]。它具有常见的免费游戏、扭蛋游戏，例如每日登录奖励和通过虚拟货币获得的随机干员，这些虚拟货币可以通过玩游戏、限时活动或可选的使用真实货币购买获得。

		背景故事
		游戏设定在一个天灾肆虐的世界中，经由天灾席卷过后的土地上出现了一种矿石——源石。源石蕴含的能量使文明迅速发展，但同时过多与源石接触也使部分人患上了不治之症——“矿石病”。感染者理论死亡率100%，且死亡时会发生二次扩散，导致感染者被世界各国政府大规模隔离和驱逐，并引发了感染者和非感染者之间的紧张关系。[10]
		
		玩家扮演失忆后的“博士”，指挥着医疗机构“罗德岛”。罗德岛的目的是寻找治疗矿石病的方法，同时保护自己免受“整合运动”等威胁。整合运动是一支受感染的无政府主义军队，一心要推翻泰拉各国政府，以报复对感染者的迫害。玩家将与罗德岛的干员一同面对天灾，在各个势力间游走，发掘不为人知的内幕。[9]`,
	)

	xyh.SelectCourse(arknights.ID)
	xb.SelectCourse(arknights.ID)
	xeh.SelectCourse(arknights.ID)
	xsh.SelectCourse(arknights.ID)

	contingency_contract_pyrolysis, _ := arknights.CreateHomework(
		"危机合约#1 灼燃作战",
		`罗德岛成立之初，物质与商业资源的匮乏带来了无数的危机，在这片充满危险的大地上立足成了一种奢望。
		而“危机合约”的存在，为罗德岛这样的企业提供了赖以生存的机会。
		
		传闻“危机合约”的前身最早由一群天灾信使所创立。
		为了对抗当时几乎无法预测的天灾、以及其他紧随其后的危机，在数位天灾信使的主导之下，天灾信使们建立了一套独特的情报交换机制，以期为各自服务的城邦，以及无数遭到天灾威胁的市民谋夺一线生机。
		
		然而，各个势力的不断发展，并未给这片大地带去友善的交流与可贵的和平。对资源的争夺，各势力间种种力量的冲突，以及在各处不断滋生、不断被惨剧验证的阴谋诡计，使得各个城邦对于信息的封锁一度达到了最高峰。即使遭遇了极为可怕的灾难，各势力也从未互相伸出援手，毋论做出哪怕一丁点微小的牺牲。
		
		天灾信使们动摇了，他们服务着的理念与实体相互割裂，现实逼迫他们面对自己的选择。
		而其中最无私的那些，决定结合各自的力量与知识，为他们承载的理想与他们热爱的人们，建立统一的“危机合约”。
		
		无论出身，无论种族，无论善恶，只要你有足够的实力——活下来，处理目标，获得报酬。或是，处理危险的目标，获得巨额的报酬，以及，活下来。
		
		最终，在这片支离破碎的大地上，一个特殊的防范机制被秘密地建立了起来。不管接下来发生的将是怎样前所未有的灾难，“危机合约”都将在阴暗处继续撑起一张无人知晓的网。
		这一切，都是为了更多的生命。
		`,
		time.Date(2023, 11, 21, 16, 0, 0, 0, time.Local),
		time.Date(2023, 12, 5, 3, 59, 0, 0, time.Local),
		time.Date(2023, 12, 10, 3, 59, 0, 0, time.Local),
	)

	filepath := fmt.Sprintf("./data/%d/%s-%s", tjw.ID, utils.GetTimeStamp(), "危机合约.txt")
	err := os.MkdirAll(fmt.Sprintf("./data/%d", tjw.ID), 0777)
	if err != nil {
		log.Println(err)
	}
	err = os.WriteFile(filepath, []byte("我超，好难"), 0666)
	if err != nil {
		log.Println(err)
	}
	file, _ := createFile(tjw.ID, "危机合约.txt", 666, filepath)
	file.Attach(contingency_contract_pyrolysis.ID, TargetTypeHomework)

	arknights.CreateHomework(
		"tjw快给我写前端（作业进行中）",
		"狠狠地写",
		time.Now(),
		time.Now().Add(7*24*time.Hour),
		time.Now().Add(14*24*time.Hour),
	)
	arknights.CreateHomework(
		"tjw快给我写前端（作业已截止，正在互评）",
		"狠狠地写",
		time.Now().Add(-7*24*time.Hour),
		time.Now(),
		time.Now().Add(7*24*time.Hour),
	)
	arknights.CreateHomework(
		"tjw快给我写前端（作业已截止，互评已截止）",
		"狠狠地写",
		time.Now().Add(-14*24*time.Hour),
		time.Now().Add(-7*24*time.Hour),
		time.Now(),
	)

	ongoing, _ := GetHomeworkByID(2)
	ending, _ := GetHomeworkByID(3)
	commenting, _ := GetHomeworkByID(4)

	ongoing.AddSubmission(1, "xyh提交")
	ongoing.AddSubmission(3, "xb提交")

	ending.AddSubmission(1, "xyh提交")
	ending.AddSubmission(3, "xb提交")
	ending.AddSubmission(4, "xeh提交")
	// ending.AddSubmission(5, "xsh提交")

	xyh_3, _ := commenting.AddSubmission(1, "xyh提交")
	xb_3, _ := commenting.AddSubmission(3, "xb提交")

	AssignComment(4)
	xyh_to_comment, _ := xb_3.GetCommentByUserId(1)
	xyh_to_comment.UpdateSelf("xyh的批阅", 99)

	xb_to_comment, _ := xyh_3.GetCommentByUserId(3)
	xb_to_comment.UpdateSelf("xb的批阅", 99)

	course2, _ := tjw.CreateCourse("课程2", time.Date(2019, 5, 1, 3, 59, 0, 0, time.Local), time.Now().Add(24*time.Hour), "kksk")
	course2.CreateHomework(
		"123",
		"321",
		time.Date(2023, 11, 21, 16, 0, 0, 0, time.Local),
		time.Date(2023, 12, 5, 3, 59, 0, 0, time.Local),
		time.Date(2023, 12, 10, 3, 59, 0, 0, time.Local),
	)
}

// func addDefaultUser() {
// 	_, err := GetUserByID(1)
// 	password := utils.RandStringRunes(8)

// 	if err == gorm.ErrRecordNotFound {
// 		defaultUser := &User{}

// 		defaultUser.ID = 1
// 		defaultUser.Username = "Admin"
// 		defaultUser.Password = utils.EncodePassword(password, utils.RandStringRunes(16))
// 		defaultUser.IsAdmin = true

// 		if err := DB.Create(&defaultUser).Error; err != nil {
// 			log.Panicf("创建初始管理员账户失败: %s\n", err)
// 		}

// 		log.Println("初始管理员账户创建完成")
// 		log.Printf("用户名: %s\n", "Admin")
// 		log.Printf("密码: %s\n", password)
// 	}
// }
