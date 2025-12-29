package main

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/paul39-33/imimix/internal/auth"
	"github.com/paul39-33/imimix/internal/database"
)

type apiConfig struct {
	dbQueries *database.Queries
	secret    string
}

type UserLogin struct {
	Username string `json:"username"`
	Pass     string `json:"pass"`
}

type User struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	Job      string    `json:"job"`
}

type MimixObj struct {
	ID          uuid.UUID `json:"id"`
	Obj         string    `json:"obj"`
	ObjType     string    `json:"obj_type"`
	PromoteDate time.Time `json:"promote_date"`
	ObjVer      string    `json:"obj_ver"`
	Lib         string    `json:"lib"`
	LibID       uuid.UUID `json:"lib_id"`
	MimixStatus string    `json:"mimix_status"`
	Developer   string    `json:"developer"`
}

type MimixLib struct {
	ID  uuid.UUID `json:"id"`
	Lib string    `json:"lib"`
}

type CreateObjReqInput struct {
	ObjName     string    `json:"obj_name" binding:"required"`
	Lib         string    `json:"lib" binding:"required"`
	ObjVer      string    `json:"obj_ver"`
	ObjType     string    `json:"obj_type"`
	PromoteDate time.Time `json:"promote_date"`
}

type ObjRequest struct {
	ID          uuid.UUID `json:"id"`
	ObjName     string    `json:"obj_name"`
	Requester   string    `json:"requester"`
	ReqStatus   string    `json:"req_status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Lib         string    `json:"lib"`
	ObjVer      string    `json:"obj_ver"`
	ObjType     string    `json:"obj_type"`
	PromoteDate time.Time `json:"promote_date"`
}

type ObjStatus struct {
	Obj         string `json:"obj"`
	MimixStatus string `json:"mimix_status"`
}

func ToNullTime(t time.Time) sql.NullTime {
	if t.IsZero() {
		return sql.NullTime{Valid: false}
	}
	return sql.NullTime{Time: t, Valid: true}
}

func NullTimeToTime(nt sql.NullTime) time.Time {
	if nt.Valid {
		return nt.Time
	}
	return time.Time{}
}

// automate middleware for authentication
func AuthMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := auth.GetBearerToken(c.Request.Header)
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		userID, err := auth.ValidateJWT(token, secret)
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Set("userID", userID)
		c.Next()
	}
}

func (cfg *apiConfig) CreateUser(c *gin.Context) {
	type parameters struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
		Job      string `json:"job" binding:"required"`
	}

	var params parameters
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// normalize job string
	jobStr := strings.ToLower(strings.TrimSpace(params.Job))

	// allowed enum values (replace/add values if your DB enum has more)
	allowedJobs := map[database.UserJob]struct{}{
		database.UserJob("user"): {},
		database.UserJob("dev"):  {},
		database.UserJob("cmt"):  {},
		database.UserJob("dc"):   {},
	}

	jobVal := database.UserJob(jobStr)
	if _, ok := allowedJobs[jobVal]; !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid job value"})
		return
	}

	//hash password
	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		log.Printf("error hashing password: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create user"})
		return
	}

	// Create user in the database
	user, err := cfg.dbQueries.CreateUser(c.Request.Context(), database.CreateUserParams{
		Username:       params.Username,
		HashedPassword: hashedPassword,
		Job:            jobVal,
	})
	if err != nil {
		log.Printf("error creating user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"user": user})
}

func (cfg *apiConfig) UserLogin(c *gin.Context) {
	var input UserLogin
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := cfg.dbQueries.GetUserByUsername(c.Request.Context(), input.Username)
	if err != nil {
		log.Printf("error getting user by username: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid username or password"})
		return
	}

	//check password
	if !auth.CheckPasswordHash(input.Pass, user.HashedPassword) {
		log.Printf("error checking password hash: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid username or password"})
		return
	}

	//access token exp duration
	accessTokenExp := 1 * time.Hour

	//generate JWT token
	token, err := auth.MakeJWT(user.ID, cfg.secret, accessTokenExp)
	if err != nil {
		log.Printf("error generating JWT token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not generate token"})
		return
	}

	userInfo := User{
		ID:       user.ID,
		Username: user.Username,
		Job:      string(user.Job),
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token": token,
		"user":         userInfo,
	})
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
	user, err := auth.ValidateJWT(token, cfg.secret)
	if err != nil {
		log.Printf("error validating token: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}

	//get user job
	userData, err := cfg.dbQueries.GetUserByID(c.Request.Context(), user)
	if err != nil {
		log.Printf("error getting user by Username: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "could not get user",
		})
		return
	}

	//check if user job is "cmt" or "dc"
	if userData.Job != "cmt" && userData.Job != "dc" {
		log.Printf("user unauthorized job: %v", userData.Job)
		c.JSON(http.StatusForbidden, gin.H{
			"error": "forbidden: insufficient permissions",
		})
		return
	}

	var params MimixObj
	//bind json parameters
	if err = c.ShouldBindJSON(&params); err != nil {
		log.Printf("error binding json: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid parameters",
		})
		return
	}

	// ensure lib exists (create if not)
	var libID uuid.UUID
	libRow, err := cfg.dbQueries.GetMimixLibByName(c.Request.Context(), params.Lib)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			createdLib, err := cfg.dbQueries.CreateMimixLib(c.Request.Context(), params.Lib)
			if err != nil {
				log.Printf("error creating lib: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create lib"})
				return
			}
			libID = createdLib.ID
		} else {
			log.Printf("error fetching lib: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not get lib"})
			return
		}
	} else {
		libID = libRow.ID
	}

	//clean obj name input
	params.Obj = strings.ToLower(strings.TrimSpace(params.Obj))
	params.Lib = strings.ToLower(strings.TrimSpace(params.Lib))
	params.Developer = strings.ToLower(strings.TrimSpace(params.Developer))

	//fix promote date null issue
	promoteDate := ToNullTime(params.PromoteDate)

	//create mimix object
	obj, err := cfg.dbQueries.AddObj(c.Request.Context(), database.AddObjParams{
		Obj:         params.Obj,
		ObjType:     params.ObjType,
		PromoteDate: promoteDate,
		Lib:         params.Lib,
		LibID:       libID,
		ObjVer:      params.ObjVer,
		MimixStatus: params.MimixStatus,
		Developer:   params.Developer,
	})

	if err != nil {
		log.Printf("error creating mimix object: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "could not create mimix object",
		})
		return
	}

	// associate the created obj with the lib id
	if err := cfg.dbQueries.UpdateObjLibID(c.Request.Context(), database.UpdateObjLibIDParams{
		ID:    obj.ID,
		LibID: libID,
	}); err != nil {
		log.Printf("error updating obj lib_id: %v", err)
		// optionally: delete created obj or return partial success
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not associate obj with lib"})
		return
	}

	createdObj := MimixObj{
		ID:          obj.ID,
		Obj:         obj.Obj,
		ObjType:     obj.ObjType,
		PromoteDate: NullTimeToTime(promoteDate),
		Lib:         obj.Lib,
		LibID:       libID,
		ObjVer:      obj.ObjVer,
		MimixStatus: obj.MimixStatus,
		Developer:   obj.Developer,
	}

	c.JSON(http.StatusOK, createdObj)
}

func (cfg *apiConfig) GetObjByName(c *gin.Context) {
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
	_, err = auth.ValidateJWT(token, cfg.secret)
	if err != nil {
		log.Printf("error validating token: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}

	//get obj name input and clean it
	objName := c.Param("obj")
	objName = strings.ToLower(strings.TrimSpace(objName))

	obj, err := cfg.dbQueries.GetObjByName(c.Request.Context(), objName)
	if err != nil {
		log.Printf("error getting mimix object: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "could not get mimix object",
		})
		return
	}

	//fix promote date null issue
	promoteDate := NullTimeToTime(obj.PromoteDate)

	resultObj := MimixObj{
		ID:          obj.ID,
		Obj:         obj.Obj,
		ObjType:     obj.ObjType,
		PromoteDate: promoteDate,
		Lib:         obj.Lib,
		ObjVer:      obj.ObjVer,
		MimixStatus: obj.MimixStatus,
	}

	c.JSON(http.StatusOK, resultObj)
}

func (cfg *apiConfig) GetObjByLib(c *gin.Context) {
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
	_, err = auth.ValidateJWT(token, cfg.secret)
	if err != nil {
		log.Printf("error validating token: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}

	//get lib name input and clean it
	libName := c.Param("lib")
	libName = strings.ToLower(strings.TrimSpace(libName))

	objs, err := cfg.dbQueries.GetObjByLib(c.Request.Context(), libName)
	if err != nil {
		log.Printf("error getting mimix objects by lib: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "could not get mimix objects",
		})
		return
	}

	var resultObjs []MimixObj
	for _, obj := range objs {
		resultObj := MimixObj{
			ID:          obj.ID,
			Obj:         obj.Obj,
			ObjType:     obj.ObjType,
			PromoteDate: NullTimeToTime(obj.PromoteDate),
			Lib:         obj.Lib,
			ObjVer:      obj.ObjVer,
			MimixStatus: obj.MimixStatus,
		}
		resultObjs = append(resultObjs, resultObj)
	}

	c.JSON(http.StatusOK, resultObjs)
}

func (cfg *apiConfig) RemoveObj(c *gin.Context) {
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
	user, err := auth.ValidateJWT(token, cfg.secret)
	if err != nil {
		log.Printf("error validating token: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}

	//get user job
	userData, err := cfg.dbQueries.GetUserByID(c.Request.Context(), user)
	if err != nil {
		log.Printf("error getting user by Username: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "could not get user",
		})
		return
	}

	//check if user job is "cmt" or "dc"
	if userData.Job != "cmt" && userData.Job != "dc" {
		log.Printf("user unauthorized job: %v", userData.Job)
		c.JSON(http.StatusForbidden, gin.H{
			"error": "forbidden: insufficient permissions",
		})
		return
	}

	//get obj id input and clean it
	objID := c.Param("obj")
	objID = strings.ToLower(strings.TrimSpace(objID))
	//parse id to uuid
	objUUID, err := uuid.Parse(objID)
	if err != nil {
		log.Printf("error parsing obj id: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid obj id",
		})
		return
	}

	err = cfg.dbQueries.RemoveObjByID(c.Request.Context(), objUUID)
	//if no obj is found
	if errors.Is(err, sql.ErrNoRows) {
		log.Printf("No matching obj found: %v", err)
		c.JSON(http.StatusNotFound, gin.H{
			"error": "no matching obj found",
		})
		return
	}
	if err != nil {
		log.Printf("error deleting mimix object: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "could not delete mimix object",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "mimix object deleted successfully",
		"obj_id":  objUUID,
	})
}

func (cfg *apiConfig) UpdateObjStatus(c *gin.Context) {
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
	id, err := auth.ValidateJWT(token, cfg.secret)
	if err != nil {
		log.Printf("error validating token: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}

	//get user job
	user, err := cfg.dbQueries.GetUserByID(c.Request.Context(), id)
	if err != nil {
		log.Printf("error getting user by Username: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "could not get user",
		})
		return
	}

	//check if user job is "dev" or "cmt"
	if user.Job != "dev" && user.Job != "cmt" && user.Job != "dc" {
		log.Printf("user unauthorized job: %v", user.Job)
		c.JSON(http.StatusForbidden, gin.H{
			"error": "forbidden: insufficient permissions",
		})
		return
	}

	//get obj name input and clean it
	objName := c.Param("obj")
	objName = strings.ToLower(strings.TrimSpace(objName))

	_, err = cfg.dbQueries.GetObjByName(c.Request.Context(), objName)
	if err != nil {
		log.Printf("error getting mimix object: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "could not get mimix object",
		})
		return
	}

	type parameters struct {
		MimixStatus string `json:"mimix_status" binding:"required"`
	}

	var params parameters

	//bind json parameters
	if err = c.ShouldBindJSON(&params); err != nil {
		log.Printf("error binding json: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid parameters",
		})
		return
	}

	err = cfg.dbQueries.UpdateObjStatus(c.Request.Context(), database.UpdateObjStatusParams{
		Obj:         objName,
		MimixStatus: params.MimixStatus,
	})

	if err != nil {
		log.Printf("error updating mimix object status: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "could not update mimix object status",
		})
		return
	}

	MimixStatus := ObjStatus{
		Obj:         objName,
		MimixStatus: params.MimixStatus,
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "mimix object status updated successfully",
		"data":    MimixStatus,
	})
}

func (cfg *apiConfig) CreateObjReq(c *gin.Context) {
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
	id, err := auth.ValidateJWT(token, cfg.secret)
	if err != nil {
		log.Printf("error validating token: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid token",
		})
		return
	}

	//get user job
	user, err := cfg.dbQueries.GetUserByID(c.Request.Context(), id)
	if err != nil {
		log.Printf("error getting user by Username: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "could not get user",
		})
		return
	}

	//check if user job is "dev" or "cmt"
	if user.Job != "dev" && user.Job != "cmt" {
		log.Printf("user unauthorized job: %v", user.Job)
		c.JSON(http.StatusForbidden, gin.H{
			"error": "forbidden: insufficient permissions",
		})
		return
	}

	var input CreateObjReqInput
	if err = c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	objReq := database.CreateMimixObjReqParams{
		ObjName:     input.ObjName,
		Requester:   user.Username,
		ReqStatus:   "pending",
		Lib:         input.Lib,
		ObjVer:      input.ObjVer,
		ObjType:     input.ObjType,
		PromoteDate: input.PromoteDate,
	}

	ObjReqRow, err := cfg.dbQueries.CreateMimixObjReq(c.Request.Context(), objReq)
	if err != nil {
		log.Printf("error creating mimix object request: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "could not create mimix object request",
		})
		return
	}

	CreatedObjReq := ObjRequest{
		ID:          ObjReqRow.ID,
		ObjName:     ObjReqRow.ObjName,
		Requester:   ObjReqRow.Requester,
		ReqStatus:   ObjReqRow.ReqStatus,
		CreatedAt:   ObjReqRow.CreatedAt,
		UpdatedAt:   ObjReqRow.UpdatedAt,
		Lib:         ObjReqRow.Lib,
		ObjVer:      ObjReqRow.ObjVer,
		ObjType:     ObjReqRow.ObjType,
		PromoteDate: ObjReqRow.PromoteDate,
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "obj request created successfully",
		"data":    CreatedObjReq,
	})
}

func (cfg *apiConfig) RemoveMimixObjReq(c *gin.Context) {
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
	id, err := auth.ValidateJWT(token, cfg.secret)
	if err != nil {
		log.Printf("error validating token: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}

	//get user job
	user, err := cfg.dbQueries.GetUserByID(c.Request.Context(), id)
	if err != nil {
		log.Printf("error getting user by Username: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "could not get user",
		})
		return
	}

	//check if user job is "dev" or "cmt" or "dc"
	if user.Job != "dev" && user.Job != "cmt" && user.Job != "dc" {
		log.Printf("user unauthorized job: %v", user.Job)
		c.JSON(http.StatusForbidden, gin.H{
			"error": "forbidden: insufficient permissions",
		})
		return
	}

	//get obj req id input
	objReqID := c.Param("reqid")
	if objReqID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing request id"})
		return
	}

	//parse id to uuid
	objReqUUID, err := uuid.Parse(objReqID)
	if err != nil {
		log.Printf("error parsing obj req id: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid obj req id",
		})
		return
	}

	//get obj name before deleting
	objNameRow, err := cfg.dbQueries.GetMimixObjReqByID(c.Request.Context(), objReqUUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "obj request not found"})
			return
		}
		log.Printf("error getting mimix object request by id: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "could not get mimix object request",
		})
		return
	}
	objName := objNameRow.ObjName

	//delete obj request

	err = cfg.dbQueries.RemoveMimixObjReq(c.Request.Context(), objReqUUID)
	if err != nil {
		log.Printf("error removing mimix object request: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "could not remove mimix object request",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "mimix object request removed successfully",
		"obj_name": objName,
	})
}

func (cfg *apiConfig) GetObjByDev(c *gin.Context) {
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
	_, err = auth.ValidateJWT(token, cfg.secret)
	if err != nil {
		log.Printf("error validating token: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}

	//get dev name input and clean it
	devName := c.Param("dev")
	devName = strings.ToLower(strings.TrimSpace(devName))

	objs, err := cfg.dbQueries.GetObjByDev(c.Request.Context(), devName)
	if err != nil {
		log.Printf("error getting mimix objects by dev: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "could not get mimix objects",
		})
		return
	}

	var resultObjs []MimixObj
	for _, obj := range objs {
		resultObjs = append(resultObjs, MimixObj{
			ID:          obj.ID,
			Obj:         obj.Obj,
			ObjType:     obj.ObjType,
			PromoteDate: NullTimeToTime(obj.PromoteDate),
			Lib:         obj.Lib,
			ObjVer:      obj.ObjVer,
			MimixStatus: obj.MimixStatus,
			Developer:   obj.Developer,
		})
	}
	c.JSON(http.StatusOK, resultObjs)
}
