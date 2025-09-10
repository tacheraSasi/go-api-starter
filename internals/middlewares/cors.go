package middlewares

import "github.com/gin-gonic/gin"

func CORSMiddleware(origins ...string) gin.HandlerFunc {
	// TODO: Figure out a way to manage multiple origins with allowing all
    return func(c *gin.Context) {
        origin := "*"
        if len(origins) > 0 && origins[0] != "" {
            origin = origins[0]
        }
        c.Header("Access-Control-Allow-Origin", origin)
        c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")
        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }
        c.Next()
    }
}