// Package main provides a script to seed the database with mock data.
package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
	"fuck_boss/backend/internal/infrastructure/config"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig("config/config.yaml")
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Connect to database
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.DBName,
		cfg.Database.SSLMode,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	fmt.Println("Connected to database successfully")

	// Clear existing posts
	fmt.Println("Clearing existing posts...")
	_, err = db.ExecContext(ctx, "TRUNCATE TABLE posts CASCADE")
	if err != nil {
		log.Fatalf("Failed to truncate posts table: %v", err)
	}
	fmt.Println("Posts table cleared")

	// Insert mock data
	fmt.Println("Inserting mock data...")
	mockPosts := []struct {
		company    string
		cityCode   string
		cityName   string
		content    string
		occurredAt *time.Time
	}{
		{
			company:  "某互联网公司",
			cityCode: "beijing",
			cityName: "北京",
			content:  "公司强制996，没有加班费，还要求员工24小时待命。管理层态度恶劣，经常PUA员工。福利待遇差，五险一金按最低标准缴纳。建议求职者慎重考虑。",
		},
		{
			company:  "某电商平台",
			cityCode: "shanghai",
			cityName: "上海",
			content:  "面试时承诺的薪资和实际发放的不一致，试用期无故延长。工作强度大，经常需要加班到深夜。团队氛围不好，同事之间竞争激烈，缺乏合作精神。",
		},
		{
			company:  "某金融科技公司",
			cityCode: "shenzhen",
			cityName: "深圳",
			content:  "公司管理混乱，经常临时安排任务，打乱工作计划。绩效考核不透明，晋升机会少。技术栈老旧，不利于个人成长。离职率高，人员流动频繁。",
		},
		{
			company:  "某游戏公司",
			cityCode: "hangzhou",
			cityName: "杭州",
			content:  "项目周期紧张，经常需要通宵加班。公司对员工健康不重视，没有合理的休息时间。薪资水平低于行业平均，但工作强度却很大。建议有更好的选择时不要来。",
		},
		{
			company:  "某教育科技公司",
			cityCode: "beijing",
			cityName: "北京",
			content:  "公司业务不稳定，经常调整方向，导致员工工作内容频繁变更。管理层决策不透明，员工缺乏参与感。福利待遇一般，没有额外的激励措施。",
		},
		{
			company:  "某物流公司",
			cityCode: "guangzhou",
			cityName: "广州",
			content:  "工作环境差，办公设施老旧。管理制度僵化，缺乏灵活性。员工培训机会少，职业发展受限。薪资增长缓慢，难以跟上生活成本上涨。",
		},
		{
			company:  "某制造业公司",
			cityCode: "chengdu",
			cityName: "成都",
			content:  "公司文化保守，不接受新思想。工作时间长，经常需要周末加班。管理层与基层员工沟通不畅，问题得不到及时解决。工作压力大，但薪资待遇一般。",
		},
		{
			company:  "某咨询公司",
			cityCode: "shanghai",
			cityName: "上海",
			content:  "项目压力大，经常需要出差。客户要求苛刻，工作强度高。公司对员工关怀不足，缺乏人性化管理。虽然薪资较高，但性价比不高，工作生活难以平衡。",
		},
		{
			company:  "某房地产公司",
			cityCode: "nanjing",
			cityName: "南京",
			content:  "行业不景气，公司业务下滑。裁员频繁，员工缺乏安全感。工作内容单一，缺乏挑战性。晋升通道狭窄，职业发展前景不明朗。",
		},
		{
			company:  "某医疗科技公司",
			cityCode: "wuhan",
			cityName: "武汉",
			content:  "公司规模小，资源有限。技术团队能力参差不齐，缺乏技术大牛指导。项目质量要求不高，不利于技术提升。薪资水平低于一线城市，但生活成本也在上涨。",
		},
		{
			company:  "某广告公司",
			cityCode: "xian",
			cityName: "西安",
			content:  "客户需求变化快，经常需要修改方案。工作时间不规律，经常需要熬夜赶项目。创意被客户随意修改，缺乏成就感。薪资待遇一般，但工作强度大。",
		},
		{
			company:  "某餐饮连锁公司",
			cityCode: "chongqing",
			cityName: "重庆",
			content:  "工作环境嘈杂，工作时间长。员工流动性大，团队不稳定。管理制度不完善，缺乏标准化流程。薪资水平低，但工作强度不小。建议有更好的选择时不要考虑。",
		},
	}

	// Insert posts with different timestamps to make them look more realistic
	now := time.Now()
	for i, post := range mockPosts {
		// Create varied timestamps (some older, some newer)
		createdAt := now.Add(-time.Duration(12-i) * 24 * time.Hour) // Spread over 12 days
		
		// Some posts have occurred_at, some don't
		var occurredAt *time.Time
		if i%3 == 0 { // Every 3rd post has occurred_at
			occurred := createdAt.Add(-time.Duration(i+1) * 7 * 24 * time.Hour)
			occurredAt = &occurred
		}

		query := `
			INSERT INTO posts (company_name, city_code, city_name, content, occurred_at, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
		`

		_, err := db.ExecContext(ctx, query,
			post.company,
			post.cityCode,
			post.cityName,
			post.content,
			occurredAt,
			createdAt,
			createdAt,
		)
		if err != nil {
			log.Fatalf("Failed to insert post %d: %v", i+1, err)
		}
		fmt.Printf("Inserted post %d: %s (%s)\n", i+1, post.company, post.cityName)
	}

	fmt.Printf("\nSuccessfully inserted %d mock posts\n", len(mockPosts))
}

