package main

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
)

// DomainContext holds a redis client
// reference for use by API routes
type DomainContext struct {
	RedisClient *redis.Client
}

func main() {

	// Set up redis client
	client := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
	})

	// Set up the context for passing the redis
	// client to routes
	domContext := DomainContext{
		RedisClient: client,
	}

	// Initialize the router
	r := gin.Default()
	r.GET("/majestic/:domain", domContext.SearchMajestic)
	r.GET("/dynamicdns/:domain", domContext.SearchDynamicDNS)

	// Run the server
	r.Run(":8080")
}

// SearchMajestic checks if a given domain name is
// on the Majestic Million
func (ctx *DomainContext) SearchMajestic(c *gin.Context) {
	domain := c.Param("domain")

	if ctx.RedisClient.SIsMember("majestic", domain).Val() {
		c.JSON(200, gin.H{domain: true})
		return
	}

	c.JSON(200, gin.H{domain: false})
	return

}

// SearchDynamicDNS checks if a given domain name is a known
// DynamicDNS provider
func (ctx *DomainContext) SearchDynamicDNS(c *gin.Context) {
	domain := c.Param("domain")

	if ctx.RedisClient.SIsMember("dynamicdns", domain).Val() {
		c.JSON(200, gin.H{domain: true})
		return
	}

	c.JSON(200, gin.H{domain: false})
	return
}
