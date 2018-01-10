package identity

import (
	"github.com/gin-gonic/gin"
)

func SetupRoute(r *gin.Engine) {

	// define your router, and use the Casbin authz middleware.
	// the access that is denied by authz will return HTTP 403 error.
	/*
	   	r.GET("/", func(c *gin.Context) {
	   		//TODO session := cthulthu.GetSession(c)
	   		//TODO user := session.Get("IdentityName").(string)
	   		//if nil != user {
	   		if len(user) > 0 {
	   			c.HTML(http.StatusOK, "index", gin.H{"Title": "테스 하하하", "UserName": user})
	   		} else {
	   			c.HTML(http.StatusOK, "index", gin.H{"Title": "테스트"})
	   		}
	   	})
	   	r.GET("/login", func(c *gin.Context) {
	   		err := c.Query("err")
	   		if len(err) > 0 {
	   			c.HTML(http.StatusOK, "login", gin.H{"Title": "Login", "Error":err})
	   		} else {
	   			c.HTML(http.StatusOK, "login", gin.H{"Title": "Login" })
	   			//c.HTML(http.StatusOK, "test", gin.H{"Title": "Login" })
	   		}
	   	})
	   	r.POST("/login", checkLogin)
	   	r.GET("/t", func(c *gin.Context) {
	   		template.InitRenderer(r)
	   	})


	   	r.GET("/admins", func(c *gin.Context) {

	   		db := GetDatabase(c)

	   		var admins []Admin
	   		db.Find(&admins)

	   		c.HTML(http.StatusOK, "list", gin.H{ "UserList": admins, })
	   	})

	   	r.GET("/admins/:id/form", func(c *gin.Context) {

	   		idParam := c.Param("id")
	   		if "new" == idParam {
	   			c.HTML(http.StatusOK, "create", nil)
	   			return
	   		}

	   		db := GetDatabase(c)
	   		id, err := strconv.ParseInt(idParam, 10, 64)
	   		if nil != err {
	   			panic(err)
	   		}
	   		var admin Admin
	   		db.First(&admin, id)

	   		c.HTML(http.StatusOK, "create", admin)

	   	})

	   	r.GET("/admins/:id", func(c *gin.Context) {

	   		idParam := c.Param("id")

	   		db := GetDatabase(c)
	   		id, err := strconv.ParseInt(idParam, 10, 64)
	   		if nil != err {
	   			panic(err)
	   		}
	   		var admin Admin
	   		db.First(&admin, id)

	   		c.HTML(http.StatusOK, "single", admin)
	   	})

	   	r.POST("/admins", func(c *gin.Context) {

	   		db := GetDatabase(c)

	   		var userForm struct {
	   			InstanceID	 string `form:"InstanceID"`
	   			UserID       string `form:"userID" binding:"required"`
	   			UserName       string `form:"userName" binding:"required"`
	   			UserPassword string `form:"userPassword" binding:"required"`
	   			UserRole string `form:"userRole" binding:"required"`
	   		}

	   		if err := c.Bind(&userForm); nil != err {
	   			// TODO 로그 및 안 꺼지도록
	   			panic(err)
	   		}

	   		model := template.Model{
	   			UpdatedAt: time.Now(),
	   			CreatedAt: time.Now(),
	   			//TODO 인스턴스 ID
	   			ID:        int64(time.Now().Second()),
	   		}
	   		admin := Admin {
	   			Identity:Identity{
	   				Model:        model,
	   				UserName:     userForm.UserName,
	   				UserID:      userForm.UserID,
	   				UserPassword: userForm.UserPassword,
	   				RoleType: userForm.UserRole,
	   			},
	   		}

	   		db.Create(&admin.Identity).Create(&admin)

	   		c.Redirect(http.StatusFound, fmt.Sprintf("/admins/%v",admin.ID))
	   	})
	   }

	   func GetDatabase(c *gin.Context) (* gorm.DB) {
	   	db := c.MustGet("dbConnection").(*gorm.DB)
	   	return db
	   }

	   func checkLogin(c *gin.Context) {
	   	var loginForm struct {
	   		UserID       string `form:"userID" binding:"required"`
	   		UserPassword string `form:"userPassword" binding:"required"`
	   		UserRemember bool   `form:"userRemember"`
	   	}
	   	if err := c.Bind(&loginForm); nil != err {
	   		// TODO 로그 및 안 꺼지도록
	   		panic(err)
	   	}

	   	db := GetDatabase(c)
	   	var user Identity
	   	noRecord := db.Where("deleted_at IS NULL AND user_name = ? AND user_password = ?",
	   		loginForm.UserID, loginForm.UserPassword).First(&user).RecordNotFound()

	   	if noRecord {
	   		c.Redirect(http.StatusFound, "/login?err=Error")
	   		return
	   	}
	   	session := template.GetSessionInstance(c)

	   	session.Set("IdentityName", user.UserName)

	   	template.SetUserRole(session, user.RoleType)

	   	redirect := c.Query("redirect")
	   	if len(redirect) > 0 {
	   		c.Redirect(http.StatusFound, redirect)
	   	} else {
	   		c.Redirect(http.StatusFound, "/")
	   	}

	   }

	   func SetupDB(r *gin.Engine) (disconnect func()) {
	   	db := template.ConnectDB(&Admin{}, &Organizer{}, &Notice{}, &Identity{})
	   	db.AutoMigrate(&Admin{}, &Organizer{}, &Notice{})

	   	// Init default values
	   	var superAdmin Admin
	   	if db.Find(&superAdmin, SuperAdminID).RecordNotFound() {
	   		db = createSuperuser(db)
	   	}

	   	r.Use(func(context *gin.Context) {
	   		context.Set("dbConnection",db)
	   		context.Next()
	   	})

	   	return func() { db.Close() }
	   }

	   func createSuperuser(db *gorm.DB) *gorm.DB {
	   	model := template.Model{
	   		UpdatedAt: time.Now(),
	   		CreatedAt: time.Now(),
	   		ID:        SuperAdminID,
	   	}
	   	admin := Admin {
	   		Identity:Identity{
	   			Model:        model,
	   			UserName:     "SuperAdmin",
	   			UserID:       "SuperAdmin",
	   			UserPassword: "SuperAdmin",
	   			RoleType:SuperAdminType,
	   			},
	   	}

	   	return db.Create(&admin.Identity).Create(&admin)*/
}
