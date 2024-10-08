package controllers

import (
	"context"
	"errors"
	"github/abbgo/yenil_yol/backend/config"
	"github/abbgo/yenil_yol/backend/helpers"
	"github/abbgo/yenil_yol/backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterCustomer(c *gin.Context) {
	db, err := config.ConnDB()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": err.Error(),
		})
		return
	}
	defer db.Close()

	// request - den gelen maglumatlar alynyar
	var customer models.Admin
	if err := c.BindJSON(&customer); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	if err := models.ValidateCustomer(customer, true); err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	// parol hashlenyan
	hashPassword, err := helpers.HashPassword(customer.Password)
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	// hemme zat yerbe yer bolsa maglumatlar customers tablisa gosulyar
	_, err = db.Exec(context.Background(), "INSERT INTO customers (full_name,phone_number,password) VALUES ($1,$2,$3)", customer.FullName, customer.PhoneNumber, hashPassword)
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":       true,
		"phone_number": customer.PhoneNumber,
		"full_name":    customer.FullName,
	})
}

func LoginCustomer(c *gin.Context) {
	db, err := config.ConnDB()
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	defer db.Close()

	// request - den maglumatlar alynyar
	var customer models.Login
	if err := c.BindJSON(&customer); err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	if !helpers.ValidatePhoneNumber(customer.PhoneNumber) {
		helpers.HandleError(c, 400, "invalid phone number")
		return
	}

	// database - den telefon belgisi request - den gelen telefon belga den bolan maglumat cekilyar
	var id, password string
	db.QueryRow(context.Background(), "SELECT id,password FROM customers WHERE phone_number = $1 AND deleted_at IS NULL", customer.PhoneNumber).Scan(&id, &password)

	// eger request - den gelen telefon belgili admin database - de yok bolsa onda error response edilyar
	if id == "" {
		helpers.HandleError(c, 404, "this customer does not exist")
		return
	}

	// eger customer bar bolsa onda paroly dogry yazypdyrmy yazmandyrmy sol barlanyar
	credentialError := helpers.CheckPassword(customer.Password, password)
	if credentialError != nil {
		helpers.HandleError(c, 400, "invalid credentials")
		return
	}

	// maglumatlar dogry bolsa auth ucin access_toke bilen resfresh_token generate edilyar
	accessTokenString /* refreshTokenString, */, err := helpers.GenerateAccessTokenForAdmin(customer.PhoneNumber, id, false)
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	// front tarapa ugratmak ucin admin - in id - si boyunca maglumatlary get edilyar
	adm, err := GetCustomerByID(id)
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token": accessTokenString,
		// "refresh_token": refreshTokenString,
		"admin": adm,
	})
}

func GetCustomerByID(id string) (models.Admin, error) {
	db, err := config.ConnDB()
	if err != nil {
		return models.Admin{}, err
	}
	defer db.Close()

	// parametrler edilip berilen id - boyunca database - den customer - in maglumatlary cekilyar
	var customer models.Admin
	db.QueryRow(context.Background(), "SELECT full_name,phone_number FROM customers WHERE deleted_at IS NULL AND id = $1", id).Scan(&customer.FullName, &customer.PhoneNumber)

	// eger parametrler edilip berilen id boyunca database - de maglumat yok bolsa error return edilyar
	if customer.PhoneNumber == "" {
		return models.Admin{}, errors.New("customer not found")
	}

	// hemme zat dogry bolsa admin - in maglumatlary return edilyar
	return customer, nil
}

func UpdateCustomer(c *gin.Context) {
	db, err := config.ConnDB()
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	defer db.Close()

	// request body - den admin - in maglumatlary alynyar
	var customer models.Admin
	if err := c.BindJSON(&customer); err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	if err := models.ValidateCustomer(customer, false); err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	// eger customer database - de bar bolsa onda onun maglumatlary request body - dan gelen maglumatlar bilen update edilyar
	_, err = db.Exec(context.Background(), "UPDATE customers SET full_name = $1 , phone_number = $2 WHERE id = $3", customer.FullName, customer.PhoneNumber, customer.ID)
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "data successfully updated",
	})
}

func UpdateCustomerPassword(c *gin.Context) {
	db, err := config.ConnDB()
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	defer db.Close()

	// request body - den maglumatlar alynyar
	var customer models.UpdatePassword
	if err := c.BindJSON(&customer); err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	if err := helpers.ValidateRecordByID("customers", customer.ID, "NULL", db); err != nil {
		helpers.HandleError(c, 404, err.Error())
		return
	}

	// maglumat bar bolsa admin - in taze paroly hashlenyar
	hashPassword, err := helpers.HashPassword(customer.Password)
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	// taze parol kone parol bilen calsylyar
	_, err = db.Exec(context.Background(), "UPDATE customers SET password = $1 WHERE id = $2", hashPassword, customer.ID)
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "password of customer successfuly updated",
	})
}

func GetCustomer(c *gin.Context) {
	adm, err := GetCustomerByID(c.Param("id"))
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"customer": adm,
	})
}

func GetCustomers(c *gin.Context) {
	var requestQuery helpers.StandartQuery
	var count uint
	var customers []models.Admin

	db, err := config.ConnDB()
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	defer db.Close()

	// request query - den maglumatlar bind edilyar
	if err := c.Bind(&requestQuery); err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	// request query - den maglumatlar validate edilyar
	if err := helpers.ValidateStructData(&requestQuery); err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	// limit we page boyunca offset hasaplanyar
	offset := requestQuery.Limit * (requestQuery.Page - 1)

	// database - den customer - lerin sany alynyar
	queryCount := `SELECT COUNT(id) FROM customers WHERE deleted_at IS NULL`
	if requestQuery.IsDeleted {
		queryCount = `SELECT COUNT(id) FROM customers WHERE deleted_at IS NOT NULL`
	}
	if err := db.QueryRow(context.Background(), queryCount).Scan(&count); err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	// request query - den status - a gora customer - lary almak ucin query yazylyar
	rowQuery := `SELECT full_name,phone_number FROM customers WHERE deleted_at IS NULL ORDER BY created_at DESC LIMIT $1 OFFSET $2`
	if requestQuery.IsDeleted {
		rowQuery = `SELECT full_name,phone_number FROM customers WHERE deleted_at IS NOT NULL ORDER BY created_at DESC LIMIT $1 OFFSET $2`
	}
	// databae - den request - den gelen limit we page boyunca limitlap customer - ler alynyar
	rowsCustomer, err := db.Query(context.Background(), rowQuery, requestQuery.Limit, offset)
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	defer rowsCustomer.Close()

	for rowsCustomer.Next() {
		var customer models.Admin
		if err := rowsCustomer.Scan(&customer.FullName, &customer.PhoneNumber); err != nil {
			helpers.HandleError(c, 400, err.Error())
			return
		}
		customers = append(customers, customer)
	}

	c.JSON(http.StatusOK, gin.H{
		"status":    true,
		"customers": customers,
		"total":     count,
	})
}

func DeleteCustomerByID(c *gin.Context) {
	// initialize database connection
	db, err := config.ConnDB()
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	defer db.Close()

	// request parametr - den customer id alynyar
	ID := c.Param("id")
	if err := helpers.ValidateRecordByID("customers", ID, "NULL", db); err != nil {
		helpers.HandleError(c, 404, err.Error())
		return
	}

	// hemme zat dogry bolsa customer - in  deleted_at - ine current_time goyulyar
	_, err = db.Exec(context.Background(), "CALL delete_customer($1)", ID)
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "data successfully deleted",
	})
}

func RestoreCustomerByID(c *gin.Context) {
	// initialize database connection
	db, err := config.ConnDB()
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	defer db.Close()

	// request parametr - den customer id alynyar
	ID := c.Param("id")
	if err := helpers.ValidateRecordByID("customers", ID, "NOT NULL", db); err != nil {
		helpers.HandleError(c, 404, err.Error())
		return
	}

	// hemme zat dogry bolsa maglumat restore edilyar
	_, err = db.Exec(context.Background(), "CALL restore_customer($1)", ID)
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "data successfully restored",
	})
}

func DeletePermanentlyCustomerByID(c *gin.Context) {
	// initialize database connection
	db, err := config.ConnDB()
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	defer db.Close()

	// request parametr - den customer id alynyar
	ID := c.Param("id")
	if err := helpers.ValidateRecordByID("customers", ID, "NOT NULL", db); err != nil {
		helpers.HandleError(c, 404, err.Error())
		return
	}

	// sonra customer database - den pozulyar
	_, err = db.Exec(context.Background(), "DELETE FROM customers WHERE id = $1", ID)
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "data successfully deleted",
	})
}
