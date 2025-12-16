package main

import (
	"github.com/google/uuid"
	"github.com/paul39-33/imimix/internal/database"
)

type apiConfig struct {
	dbQueries *database.Queries
	secret    string
}

type UserLogin struct {
	ID       uuid.UUID `json:"id"`
	Email    string    `json:"email"`
	Password string    `json:"password"`
}

type User struct {
	ID	   uuid.UUID `json:"id"`
	Email    string    `json:"email"`
	Pass string    `json:"pass"`
}

type MimixObj struct {
	ID      uuid.UUID `json:"id"`
	Obj   	string    `json:"obj"`
	ObjType string    `json:"obj_type"`
	PromoteDate time.Time `json:"promote_date"`
	ObjVer  string    `json:"obj_ver"`
	Lib 	string   `json:"lib"`
	MimixStatus string    `json:"mimix_status"`
}

type MimixLib struct {
	ID      uuid.UUID `json:"id"`
	Lib	 string    `json:"lib"`
}

func (cfg *apiConfig) CreateObj(c *gin.Context) {
	//get user token
	token, err := auth.GetBearerToken(c.Request.Header)
	if err != nil {
		log.Printf("error getting bearer token: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid token",
		})
		return
	}

	//validate user token
	userID, err := auth.ValidateToken(token, cfg.secret)
	if err != nil {
		log.Printf("error validating token: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}

	type parameters struct {
		Obj string `json:"obj" binding:"required"`
		ObjType string `json:"obj_type" binding:"required"`
		PromoteDate time.Time `json:"promote_date" binding:"required"`
		Lib string `json:"lib" binding:"required"`
		ObjVer string `json:"obj_ver" binding:"required"`
		MimixStatus string `json:"mimix_status" binding:"required"`
	}

	var params parameters

	//bind json parameters
	if err := c.ShouldBindJSON(&params);  err != nil {
		log.Printf("error binding json: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid parameters",
		})
		return
	}

	//create mimix object
	obj, err := cfg.dbQueries.CreateMimixObj(c.Request.Context(), database.CreateMimixObjParams{
		Obj:         params.Obj,
		ObjType:     params.ObjType,
		PromoteDate: params.PromoteDate,
		Lib:         params.Lib,
		ObjVer:      params.ObjVer,
		MimixStatus: params.MimixStatus,
	})

	if err != nil {
		log.Printf("error creating mimix object: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "could not create mimix object",
		})
		return
	}

	createdObj := MimixObj{
		ID:		  obj.ID,
		Obj:         obj.Obj,
		ObjType:     obj.ObjType,
		PromoteDate: obj.PromoteDate,
		Lib:         obj.Lib,
		ObjVer:      obj.ObjVer,
		MimixStatus: obj.MimixStatus,
	}

	c.JSON(http.StatusOK, createdObj)
}

func (cfg *apiConfig) GetObj(c *gin.Context) {
	
}
