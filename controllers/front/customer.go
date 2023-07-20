package controllers

import (
	"context"
	"errors"
	"github/abbgo/yenil_yol/backend/config"
	"github/abbgo/yenil_yol/backend/helpers"
	"github/abbgo/yenil_yol/backend/models"
	"net/http"
	"strconv"

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
	var customer models.Customer
	if err := c.BindJSON(&customer); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	if err := models.ValidateCustomer(customer.PhoneNumber, "", true); err != nil {
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
	var customer models.ShopOwnerLogin
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
	if err := db.QueryRow(context.Background(), "SELECT id,password FROM customers WHERE phone_number = $1 AND deleted_at IS NULL", customer.PhoneNumber).Scan(&id, &password); err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

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
	accessTokenString, refreshTokenString, err := helpers.GenerateAccessTokenForAdmin(customer.PhoneNumber, id, false)
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
		"access_token":  accessTokenString,
		"refresh_token": refreshTokenString,
		"admin":         adm,
	})

}

func GetCustomerByID(id string) (models.Customer, error) {
	db, err := config.ConnDB()
	if err != nil {
		return models.Customer{}, err
	}
	defer db.Close()

	// parametrler edilip berilen id - boyunca database - den customer - in maglumatlary cekilyar
	var customer models.Customer
	rowCustomer, err := db.Query(context.Background(), "SELECT full_name,phone_number FROM customers WHERE deleted_at IS NULL AND id = $1", id)
	if err != nil {
		return models.Customer{}, err
	}
	defer rowCustomer.Close()

	for rowCustomer.Next() {
		if err := rowCustomer.Scan(&customer.FullName, &customer.PhoneNumber); err != nil {
			return models.Customer{}, err
		}
	}

	// eger parametrler edilip berilen id boyunca database - de maglumat yok bolsa error return edilyar
	if customer.PhoneNumber == "" {
		return models.Customer{}, errors.New("customer not found")
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
	var customer models.CustomerUpdate
	if err := c.BindJSON(&customer); err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	if err := models.ValidateCustomer(customer.PhoneNumber, customer.ID, false); err != nil {
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
	var customer models.CustomerUpdatePassword
	if err := c.BindJSON(&customer); err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	// gelen id den bolan maglumat database - de barmy ya yok sol barlanyar
	var customer_id string
	if err := db.QueryRow(context.Background(), "SELECT id FROM customers WHERE id = $1 AND deleted_at IS NULL", customer.ID).Scan(&customer_id); err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	// eger gelen id den bolan maglumat database - de yok bolsa error return edilyar
	if customer_id == "" {
		helpers.HandleError(c, 404, "customer not found")
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

	customerID, hasID := c.Get("customer_id")
	if !hasID {
		helpers.HandleError(c, 400, "customerID is required")
		return
	}

	var ok bool
	customer_id, ok := customerID.(string)
	if !ok {
		helpers.HandleError(c, 400, "customerID must be string")
		return
	}

	adm, err := GetCustomerByID(customer_id)
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"customer": adm,
	})

}

func GetCustomers(c *gin.Context) {

	db, err := config.ConnDB()
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	defer db.Close()

	// request parametr - den limit alynyar
	limitStr := c.Query("limit")
	if limitStr == "" {
		helpers.HandleError(c, 400, "limit is required")
		return
	}
	limit, err := strconv.ParseUint(limitStr, 10, 32)
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	// // request parametr - den page alynyar
	pageStr := c.Query("page")
	if pageStr == "" {
		helpers.HandleError(c, 400, "page is required")
		return
	}
	page, err := strconv.ParseUint(pageStr, 10, 32)
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	// limit we page boyunca offset hasaplanyar
	offset := limit * (page - 1)

	// database - den admin - lerin sany alynyar
	countOfCustomers := 0
	if err := db.QueryRow(context.Background(), "SELECT COUNT(id) FROM customers WHERE deleted_at IS NULL").Scan(&countOfCustomers); err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	// databae - den request - den gelen limit we page boyunca limitlap customer - ler alynyar
	var customers []models.Customer
	rowsCustomer, err := db.Query(context.Background(), "SELECT full_name,phone_number FROM customers WHERE deleted_at IS NULL LIMIT $1 OFFSET $2", limit, offset)
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	defer rowsCustomer.Close()

	for rowsCustomer.Next() {
		var customer models.Customer
		if err := rowsCustomer.Scan(&customer.FullName, &customer.PhoneNumber); err != nil {
			helpers.HandleError(c, 400, err.Error())
			return
		}
		customers = append(customers, customer)
	}

	c.JSON(http.StatusOK, gin.H{
		"status":    true,
		"customers": customers,
		"total":     countOfCustomers,
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

	// gelen id den bolan maglumat database - de barmy sol barlanyar
	var id string
	if err := db.QueryRow(context.Background(), "SELECT id FROM customers WHERE id = $1 AND deleted_at IS NULL", ID).Scan(&id); err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	// eger database - de gelen id degisli maglumat yok bolsa error return edilyar
	if id == "" {
		helpers.HandleError(c, 404, "record not found")
		return
	}

	// hemme zat dogry bolsa customer - in  deleted_at - ine current_time goyulyar
	_, err = db.Exec(context.Background(), "UPDATE customers SET deleted_at=NOW() WHERE id = $1", ID)
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

	// alynan id den bolan maglumat database - de barmy ya yok sol barlanyar
	var id string
	if err := db.QueryRow(context.Background(), "SELECT id FROM customers WHERE id = $1 AND deleted_at IS NOT NULL", ID).Scan(&id); err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	// eger database sol id degisli maglumat yok bolsa error return edilyar
	if id == "" {
		helpers.HandleError(c, 404, "record not found")
		return
	}

	// hemme zat dogry bolsa maglumat restore edilyar
	_, err = db.Exec(context.Background(), "UPDATE customers SET deleted_at=NULL WHERE id = $1", ID)
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "data successfully restored",
	})

}
