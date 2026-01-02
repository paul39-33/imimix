package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/paul39-33/imimix/internal/database"
)

func main() {
	//load environment variables from .env file
	godotenv.Load()
	//secret := os.Getenv("SECRET_KEY")
	db_url := os.Getenv("DATABASE_URL")
	secret := os.Getenv("SECRET_KEY")
	db, err := sql.Open("postgres", db_url)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	log.Print("db connected successfully")
	defer db.Close()

	dbQueries := database.New(db)

	//apiCfg
	apiCfg := apiConfig{
		dbQueries: dbQueries,
		secret:    secret,
	}

	//create Gin router
	r := gin.Default()

	//set frontend directory by url /app/... from .
	r.StaticFS("/app", gin.Dir(".", false))

	api := r.Group("/api")
	{
		api.POST("/create_user", apiCfg.CreateUser)
		api.POST("/login", apiCfg.UserLogin)
		api.POST("/add_mimix_obj", apiCfg.CreateObj)
		api.POST("/create_obj_req", apiCfg.CreateObjReq)
		api.DELETE("/delete_mimix_obj/:obj", apiCfg.RemoveObj)
		api.DELETE("/delete_obj_req/:reqid", apiCfg.RemoveMimixObjReq)
		api.GET("/get_mimix_obj_by_name/:name", apiCfg.GetObjByName)
		api.GET("/get_mimix_obj/:lib", apiCfg.GetObjByLib)
		api.GET("/get_mimix_obj_by_dev/:dev", apiCfg.GetObjByDev)
		api.PATCH("/update_mimix_obj_status/:obj", apiCfg.UpdateObjStatus)

		api.PATCH("/update_mimix_obj_info/:id", apiCfg.UpdateObjInfo) // handler expects :id
		api.POST("/add_obj_to_obj_req/:id", apiCfg.ObjtoObjReq)       // handler expects :id
		api.POST("/convert_obj_req/:reqid", apiCfg.ObjReqToObj)       // handler expects :reqid
	}

	//start server on port 8080
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

}
