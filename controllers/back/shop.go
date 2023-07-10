package controllers

// func CreateShop(c *gin.Context) {

// 	// initialize database connection
// 	db, err := config.ConnDB()
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"status":  false,
// 			"message": err.Error(),
// 		})
// 		return
// 	}
// 	defer func() {
// 		if err := db.Close(); err != nil {
// 			c.JSON(http.StatusBadRequest, gin.H{
// 				"status":  false,
// 				"message": err.Error(),
// 			})
// 			return
// 		}
// 	}()

// 	var brend models.Brend
// 	if err := c.BindJSON(&brend); err != nil {
// 		c.JSON(http.StatusBadRequest, err.Error())
// 		return
// 	}

// 	// CREATE BREND
// 	result, err := db.Query("INSERT INTO brends (name,image,slug) VALUES ($1,$2)", brend.Name, brend.Image, slug.MakeLang(brend.Name, "en"))
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"status":  false,
// 			"message": err.Error(),
// 		})
// 		return
// 	}
// 	defer func() {
// 		if err := result.Close(); err != nil {
// 			c.JSON(http.StatusBadRequest, gin.H{
// 				"status":  false,
// 				"message": err.Error(),
// 			})
// 			return
// 		}
// 	}()

// 	c.JSON(http.StatusOK, gin.H{
// 		"status":  true,
// 		"message": "data successfully added",
// 	})

// }
