package db

import (
	"context"
	"database/sql"
	"log"
	"strconv"
	"math/rand/v2"

	"github.com/DiscoDoggy/terabytes/go_backend/internal/store"
)

var titles = []string{
	"Revolutionizing Web Performance with Server-Side Rendering",
	"Kubernetes for Beginners: A Step-by-Step Guide",
	"The Art of Writing Clean Code: Best Practices",
	"Exploring the Future of Artificial Intelligence in DevOps",
	"Understanding Database Sharding and Its Use Cases",
	"Microservices vs. Monoliths: Choosing the Right Architecture",
	"Building Scalable APIs with GraphQL",
	"Debugging Complex Systems: Tips from the Trenches",
	"The Evolution of Frontend Frameworks: React, Vue, and Angular",
	"Automating CI/CD Pipelines for Maximum Efficiency",
	"Secrets of High-Performance SQL Queries",
	"A Deep Dive into Cloud-Native Application Design",
	"How We Optimized Our Mobile App for a 3x Speed Boost",
	"Real-Time Data Processing with Apache Kafka",
}

var descriptions = []string{
	"Discover how server-side rendering can drastically improve website load times and SEO.",
	"Learn Kubernetes concepts, from pods to clusters, in this beginner-friendly guide.",
	"Tips and techniques to write code that’s readable, maintainable, and scalable.",
	"An exploration of how AI is transforming DevOps practices for the better.",
	"Learn how database sharding works and when to use it for scalable systems.",
	"Compare microservices and monoliths to choose the right approach for your projects.",
	"Learn how to build robust APIs with GraphQL and avoid common pitfalls.",
	"Real-world strategies for debugging distributed systems and solving critical issues.",
	"Explore how React, Vue, and Angular have evolved and their future trajectories.",
	"Step up your CI/CD game with automation tips that save time and reduce errors.",
	"Learn the secrets to crafting SQL queries that are fast and efficient.",
	"An in-depth look at how to design applications for the cloud-native era.",
	"A case study on optimizing a mobile app to load faster and perform better.",
	"Learn the basics of processing real-time data streams with Apache Kafka.",
}

func Seed(store store.Storage, db *sql.DB) {
	users := createRandomUsers()
	
	ctx := context.Background()

	tx, _ := db.BeginTx(ctx, nil)

	for i := 0; i < len(users); i++ {
		err := store.Users.Create(ctx, tx, &users[i])
		if err != nil {
			log.Println("Failure inserting user", err.Error())
			continue
		}
	}

	allUserIds := `
		SELECT id FROM accounts
	`
	userIds := make([]string, 0)
	rows, err := db.QueryContext(ctx, allUserIds)
	if err != nil {
		log.Fatalln("Failure Reading user", err.Error())
	}
	for rows.Next() {
		var currId string
		err = rows.Scan(&currId)
		
		if err != nil {
			log.Println("Error reading user id into variable", err.Error())
		}

		userIds = append(userIds, currId)
	}

	for i:= 0; i < len(userIds); i++ {
		blogPosts := createRandomPosts(userIds[i])
		for _, value := range blogPosts {
			err = store.Posts.Create(ctx, &value)
			if err != nil {
				log.Println("ERROR INSERTING POST", err.Error())
				continue
			}
		}
	}

}

func createRandomUsers() []store.User {
	baseUsername := "test"
	mockUsers := make([]store.User, 0)

	for i := 1; i <= 10; i++ {
		tempMockUser := store.User{
			Username: baseUsername + strconv.Itoa(i),
			Email: baseUsername + strconv.Itoa(i) + "@email.com",
		}

		tempMockUser.Password.HashPassword("12345")

		mockUsers = append(mockUsers, tempMockUser)
	}

	return mockUsers
}

func createRandomPosts(userId string) []store.BlogPost {
	numUserPosts := rand.IntN(25) + 1
	blogPosts := make([]store.BlogPost, 0)

	for i:=0; i < numUserPosts; i++ {
		var bp store.BlogPost
		blogContent := make([]store.BlogPostContent, 0)
		var content store.BlogPostContent

		randIndex := rand.IntN(len(titles) - 1)

		bp.UserId = userId
		bp.Title = titles[randIndex]
		bp.Description = descriptions[randIndex]
		bp.Tags = []string{"Kubernetes", "Rust", "Go", "C++"}

		content.ContentType = "text"
		content.ContentData = "test blog content"
		content.ContentOrder = 1

		blogContent = append(blogContent, content)
		bp.Content = blogContent

		blogPosts = append(blogPosts, bp)

	}

	return blogPosts
}