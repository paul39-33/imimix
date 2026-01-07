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
	Keterangan  string    `json:"keterangan"`
	UpdatedAt   time.Time `json:"updated_at"`
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
	Developer   string    `json:"developer"`
}

type ObjRequest struct {
	ID            uuid.UUID `json:"id"`
	ObjName       string    `json:"obj_name"`
	Requester     string    `json:"requester"`
	Developer     string    `json:"developer,omitempty"`
	ReqStatus     string    `json:"req_status"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	Lib           string    `json:"lib"`
	ObjVer        string    `json:"obj_ver"`
	ObjType       string    `json:"obj_type"`
	PromoteDate   time.Time `json:"promote_date"`
	SourceObjID   uuid.UUID `json:"source_obj_id,omitempty"`
	PromoteStatus string    `json:"promote_status,omitempty"`
}

type ObjStatus struct {
	Obj         string               `json:"obj"`
	MimixStatus database.MimixStatus `json:"mimix_status"`
}

var allowedUserJobs = map[string]database.UserJob{
	"cmt":  database.UserJobCmt,
	"dev":  database.UserJobDev,
	"dc":   database.UserJobDc,
	"user": database.UserJobUser,
}

var allowedMimixStatus = map[string]database.MimixStatus{
	"unset":              database.MimixStatusUnset,
	"done":               database.MimixStatusDone,
	"daftarkan":          database.MimixStatusDaftarkan,
	"tidak perlu daftar": database.MimixStatusTidakperludaftar,
	"on progress":        database.MimixStatusOnprogress,
}

var allowedReqStatus = map[string]database.ReqStatus{
	"pending":   database.ReqStatusPending,
	"completed": database.ReqStatusCompleted,
}

var allowedPromoteStatus = map[string]database.PromoteStatus{
	"in_progress": database.PromoteStatusInProgress,
	"deployed":    database.PromoteStatusDeployed,
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

func NullStringToString(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return ""
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
		Username        string `json:"username" binding:"required"`
		Password        string `json:"password" binding:"required"`
		ConfirmPassword string `json:"confirm_password" binding:"required"`
		Job             string `json:"job" binding:"required"`
	}

	var params parameters
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if params.Password != params.ConfirmPassword {
		c.JSON(http.StatusBadRequest, gin.H{"error": "passwords do not match"})
		return
	}

	// normalize job string
	jobStr := strings.ToLower(strings.TrimSpace(params.Job))
	// normalize username
	params.Username = strings.ToLower(strings.TrimSpace(params.Username))

	// allowed enum values (replace/add values if your DB enum has more)
	job, ok := allowedUserJobs[jobStr]
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid job type"})
		return
	}

	// check if user already exists
	exists, err := cfg.dbQueries.CheckUserExists(c.Request.Context(), params.Username)
	if err != nil {
		log.Printf("error checking if user exists: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not check user existence"})
		return
	}
	if exists {
		c.JSON(http.StatusConflict, gin.H{"error": "username already exists"})
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
		Job:            job,
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


	// normalize username
	input.Username = strings.ToLower(strings.TrimSpace(input.Username))

	user, err := cfg.dbQueries.GetUserByUsername(c.Request.Context(), input.Username)
	if err != nil {
		log.Printf("error getting user by username: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid username or password"})
		return
	}

	//check password
	if !auth.CheckPasswordHash(input.Pass, user.HashedPassword) {
		log.Printf("Password verification failed for user: %s", input.Username)
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

	// validate mimix status
	statusKey := strings.ToLower(strings.TrimSpace(string(params.MimixStatus)))
	statusVal, ok := allowedMimixStatus[statusKey]
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid mimix_status"})
		return
	}

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
		MimixStatus: statusVal,
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
		MimixStatus: string(obj.MimixStatus),
		Developer:   obj.Developer,
		UpdatedAt:   obj.UpdatedAt,
	}

	c.JSON(http.StatusOK, createdObj)
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

	// validate incoming status string and convert to enum
	statusKey := strings.ToLower(strings.TrimSpace(params.MimixStatus))
	statusVal, ok := allowedMimixStatus[statusKey]
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid mimix_status"})
		return
	}

	err = cfg.dbQueries.UpdateObjStatus(c.Request.Context(), database.UpdateObjStatusParams{
		Obj:         objName,
		MimixStatus: statusVal,
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
		MimixStatus: statusVal,
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
		Developer: sql.NullString{
			String: input.Developer,
			Valid:  strings.TrimSpace(input.Developer) != "",
		},
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
		ReqStatus:   string(ObjReqRow.ReqStatus),
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



func (cfg *apiConfig) UpdateObjInfo(c *gin.Context) {
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

	//check if user job is "cmt" or "dev"
	if user.Job != "cmt" && user.Job != "dev" {
		log.Printf("user unauthorized job: %v", user.Job)
		c.JSON(http.StatusForbidden, gin.H{
			"error": "forbidden: insufficient permissions",
		})
		return
	}

	//get obj by id
	objID := c.Param("id")

	objUUID, err := uuid.Parse(objID)
	if err != nil {
		log.Printf("error parsing obj id: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid obj id",
		})
		return
	}

	type parameters struct {
		Obj         string    `json:"obj"`
		ObjType     string    `json:"obj_type"`
		Lib         string    `json:"lib"`
		PromoteDate time.Time `json:"promote_date"`
		ObjVer      string    `json:"obj_ver"`
		Developer   string    `json:"developer"`
		MimixStatus string    `json:"mimix_status"`
		Keterangan  string    `json:"keterangan"`
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

	// validate mimix status
	statusKey := strings.ToLower(strings.TrimSpace(params.MimixStatus))
	statusVal, ok := allowedMimixStatus[statusKey]
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid mimix_status"})
		return
	}

	//fix promote date null issue
	promoteDate := ToNullTime(params.PromoteDate)

	keteranganNull := sql.NullString{
		String: params.Keterangan,
		Valid:  strings.TrimSpace(params.Keterangan) != "",
	}

	updatedObj, err := cfg.dbQueries.UpdateObjInfo(c.Request.Context(), database.UpdateObjInfoParams{
		ID:          objUUID,
		Obj:         params.Obj,
		ObjType:     params.ObjType,
		PromoteDate: promoteDate,
		ObjVer:      params.ObjVer,
		Developer:   params.Developer,
		MimixStatus: statusVal,
		Lib:         params.Lib,
		Keterangan:  keteranganNull,
	})
	if err != nil {
		log.Printf("error updating obj info: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "could not update obj info",
		})
		return
	}

	// map to api struct
	respObj := MimixObj{
		ID:          updatedObj.ID,
		Obj:         updatedObj.Obj,
		ObjType:     updatedObj.ObjType,
		PromoteDate: NullTimeToTime(updatedObj.PromoteDate),
		Lib:         updatedObj.Lib,
		LibID:       updatedObj.LibID,
		ObjVer:      updatedObj.ObjVer,
		MimixStatus: string(updatedObj.MimixStatus),
		Developer:   updatedObj.Developer,
		Keterangan:  NullStringToString(updatedObj.Keterangan),
		UpdatedAt:   updatedObj.UpdatedAt,
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "obj info updated successfully",
		"data":    respObj,
	})
}

func (cfg *apiConfig) ObjtoObjReq(c *gin.Context) {
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
	if user.Job != "dev" && user.Job != "cmt" {
		log.Printf("user unauthorized job: %v", user.Job)
		c.JSON(http.StatusForbidden, gin.H{
			"error": "forbidden: insufficient permissions",
		})
		return
	}

	//get obj by id
	objID := c.Param("id")

	objUUID, err := uuid.Parse(objID)
	if err != nil {
		log.Printf("error parsing obj id: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid obj id",
		})
		return
	}

	obj, err := cfg.dbQueries.GetObjByID(c.Request.Context(), objUUID)
	if err != nil {
		log.Printf("error getting mimix object: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "could not get mimix object",
		})
		return
	}

	err = cfg.dbQueries.AddObjToObjReq(c.Request.Context(), database.AddObjToObjReqParams{
		ID:        obj.ID,
		Requester: user.Username,
		ReqStatus: "pending",
	})
	if err != nil {
		log.Printf("error adding obj to obj request: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "could not add obj to obj request",
		})
		return
	}

	//update obj mimix status to "on progress"
	err = cfg.dbQueries.UpdateObjStatus(c.Request.Context(), database.UpdateObjStatusParams{
		Obj:         obj.Obj,
		MimixStatus: database.MimixStatusOnprogress,
	})
	if err != nil {
		log.Printf("error updating mimix object status: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "could not update mimix object status",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "obj added to obj request successfully",
	})
}

func (cfg *apiConfig) ObjReqToObj(c *gin.Context) {
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

	//check if user job is "dc"
	if userData.Job != "dc" {
		log.Printf("user unauthorized job: %v", userData.Job)
		c.JSON(http.StatusForbidden, gin.H{
			"error": "forbidden: insufficient permissions",
		})
		return
	}

	//get obj req by id
	objReqID := c.Param("reqid")

	objReqUUID, err := uuid.Parse(objReqID)
	if err != nil {
		log.Printf("error parsing obj req id: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid obj req id",
		})
		return
	}

	objReq, err := cfg.dbQueries.GetMimixObjReqByID(c.Request.Context(), objReqUUID)
	if err != nil {
		log.Printf("error getting mimix object request: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "could not get mimix object request",
		})
		return
	}

	//change obj req status to "completed"
	err = cfg.dbQueries.CompleteMimixObjReq(c.Request.Context(), objReq.ID)
	if err != nil {
		log.Printf("error updating mimix object request status: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "could not update mimix object request status",
		})
		return
	}

	//check if obj req already exists as obj
	var sourceObj database.MimixObj
	var objExists bool

	if objReq.SourceObjID.Valid {
		sourceObjID := objReq.SourceObjID.UUID
		var err error
		sourceObj, err = cfg.dbQueries.GetObjByID(c.Request.Context(), sourceObjID)
		if err == nil {
			objExists = true
		} else if !errors.Is(err, sql.ErrNoRows) {
			log.Printf("error getting mimix object by source obj id: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "could not get mimix object",
			})
			return
		}
	}

	if objExists {
		//change obj mimix status to "completed"
		err = cfg.dbQueries.UpdateObjStatus(c.Request.Context(), database.UpdateObjStatusParams{
			Obj:         sourceObj.Obj,
			MimixStatus: database.MimixStatusDone,
		})
		if err != nil {
			log.Printf("error updating mimix object status: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "could not update mimix object status",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "obj request already exists as obj, status updated to completed",
			"obj_id":  sourceObj.ID,
		})
		return
	}

	{
		//check if new obj lib exists (create if not)
		var libID uuid.UUID
		libRow, err := cfg.dbQueries.GetMimixLibByName(c.Request.Context(), objReq.Lib)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				createdLib, err := cfg.dbQueries.CreateMimixLib(c.Request.Context(), objReq.Lib)
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

		//create new obj from obj req
		newObj, err := cfg.dbQueries.AddObj(c.Request.Context(), database.AddObjParams{
			Obj:         objReq.ObjName,
			ObjType:     objReq.ObjType,
			PromoteDate: ToNullTime(objReq.PromoteDate),
			Lib:         objReq.Lib,
			LibID:       libID,
			ObjVer:      objReq.ObjVer,
			MimixStatus: database.MimixStatusDone,
			Developer:   NullStringToString(objReq.Developer),
		})
		if err != nil {
			log.Printf("error creating mimix object from obj request: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "could not create mimix object from obj request",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "obj request converted to obj successfully",
			"obj_id":  newObj.ID,
		})
	}
}

func (cfg *apiConfig) UpdateObjReqInfo(c *gin.Context) {
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

	//get obj req by id
	objReqID := c.Param("id")

	objReqUUID, err := uuid.Parse(objReqID)
	if err != nil {
		log.Printf("error parsing obj req id: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid obj req id",
		})
		return
	}

	type parameters struct {
		ObjName       string    `json:"obj_name"`
		Lib           string    `json:"lib"`
		PromoteDate   time.Time `json:"promote_date"`
		ObjVer        string    `json:"obj_ver"`
		ObjType       string    `json:"obj_type"`
		Developer     string    `json:"developer"`
		PromoteStatus string    `json:"promote_status"`
		ReqStatus     string    `json:"req_status"`
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



	// prepare developer as sql.NullString
	devNull := sql.NullString{
		String: params.Developer,
		Valid:  strings.TrimSpace(params.Developer) != "",
	}

	// validate req_status against allowed values
	reqKey := strings.ToLower(strings.TrimSpace(params.ReqStatus))
	reqVal, ok := allowedReqStatus[reqKey]
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid req_status"})
		return
	}

	// prepare promote_status as nullable enum
	var promoteStatus database.NullPromoteStatus
	if strings.TrimSpace(params.PromoteStatus) != "" {
		psKey := strings.ToLower(strings.TrimSpace(params.PromoteStatus))
		psVal, ok := allowedPromoteStatus[psKey]
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid promote_status"})
			return
		}
		promoteStatus = database.NullPromoteStatus{
			PromoteStatus: psVal,
			Valid:         true,
		}
	} else {
		promoteStatus = database.NullPromoteStatus{Valid: false}
	}

	updatedMimixObjReq, err := cfg.dbQueries.UpdateMimixObjReqInfo(c.Request.Context(), database.UpdateMimixObjReqInfoParams{
		ID:            objReqUUID,
		ObjName:       params.ObjName,
		Lib:           params.Lib,
		PromoteDate:   params.PromoteDate,
		ObjVer:        params.ObjVer,
		ObjType:       params.ObjType,
		Developer:     devNull,
		PromoteStatus: promoteStatus,
		ReqStatus:     reqVal,
	})
	if err != nil {
		log.Printf("error updating obj req info: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "could not update obj req info",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "obj req info updated successfully",
		"data":    updatedMimixObjReq,
	})
}

func (cfg *apiConfig) SearchObj(c *gin.Context) {
	// get user token
	token, err := auth.GetBearerToken(c.Request.Header)
	if err != nil {
		log.Printf("error getting bearer token: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid token",
		})
		return
	}

	// validate user token
	_, err = auth.ValidateJWT(token, cfg.secret)
	if err != nil {
		log.Printf("error validating token: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}

	// get search input and clean it
	query := c.Param("query")
	query = strings.ToLower(strings.TrimSpace(query))
	search := sql.NullString{
		String: query,
		Valid:  true,
	}

	// search objs by obj / lib / developer
	objs, err := cfg.dbQueries.SearchMimixObj(
		c.Request.Context(),
		search,
	)
	if err != nil {
		log.Printf("error searching mimix objects: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "could not search mimix objects",
		})
		return
	}

	// map DB models → API models
	var resultObjs []MimixObj
	for _, obj := range objs {
		resultObjs = append(resultObjs, MimixObj{
			ID:          obj.ID,
			Obj:         obj.Obj,
			ObjType:     obj.ObjType,
			PromoteDate: NullTimeToTime(obj.PromoteDate),
			Lib:         obj.Lib,
			ObjVer:      obj.ObjVer,
			MimixStatus: string(obj.MimixStatus),
			Developer:   obj.Developer,
			Keterangan:  NullStringToString(obj.Keterangan),
			UpdatedAt:   obj.UpdatedAt,
		})
	}

	c.JSON(http.StatusOK, resultObjs)
}

type MimixObjReq struct {
	ID            uuid.UUID `json:"id"`
	ObjName       string    `json:"obj_name"`
	Requester     string    `json:"requester"`
	ReqStatus     string    `json:"req_status"`
	Lib           string    `json:"lib"`
	ObjVer        string    `json:"obj_ver"`
	ObjType       string    `json:"obj_type"`
	PromoteDate   time.Time `json:"promote_date"`
	Developer     string    `json:"developer"`
	PromoteStatus string    `json:"promote_status"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func (cfg *apiConfig) SearchObjReq(c *gin.Context) {
	// get user token
	token, err := auth.GetBearerToken(c.Request.Header)
	if err != nil {
		log.Printf("error getting bearer token: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid token",
		})
		return
	}

	// validate user token
	_, err = auth.ValidateJWT(token, cfg.secret)
	if err != nil {
		log.Printf("error validating token: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}

	// get search input and clean it
	query := c.Param("query")
	query = strings.ToLower(strings.TrimSpace(query))
	search := sql.NullString{
		String: query,
		Valid:  true,
	}

	// search obj requests
	reqs, err := cfg.dbQueries.SearchMimixObjReq(
		c.Request.Context(),
		search,
	)
	if err != nil {
		log.Printf("error searching mimix object requests: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "could not search mimix object requests",
		})
		return
	}

	// map DB models → API models
	var resultReqs []MimixObjReq
	for _, req := range reqs {
		// handle nullable promote_status
		var ps string
		if req.PromoteStatus.Valid {
			ps = string(req.PromoteStatus.PromoteStatus)
		} else {
			ps = ""
		}

		resultReqs = append(resultReqs, MimixObjReq{
			ID:            req.ID,
			ObjName:       req.ObjName,
			Requester:     req.Requester,
			ReqStatus:     string(req.ReqStatus),
			Lib:           req.Lib,
			ObjVer:        req.ObjVer,
			ObjType:       req.ObjType,
			PromoteDate:   req.PromoteDate,
			Developer:     NullStringToString(req.Developer),
			PromoteStatus: ps,
			CreatedAt:     req.CreatedAt,
			UpdatedAt:     req.UpdatedAt,
		})
	}

	c.JSON(http.StatusOK, resultReqs)
}
