package controllers

import (
	"context"
	"github/abbgo/yenil_yol/backend/config"
	"github/abbgo/yenil_yol/backend/helpers"
	"github/abbgo/yenil_yol/backend/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ProductID struct {
	IDS []string `json:"product_ids"`
}

func AddOrRemoveLike(c *gin.Context) {
	db, err := config.ConnDB()
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}
	defer db.Close()

	// bu yerde api query - den status ady bilen status alynyar
	// eger status = true bolsa onda halanlarym sahypa haryt gosulyar
	// eger status = false bolsa onda halanlarym sahypadan haryt ayrylyar
	statusStr := c.DefaultQuery("status", "true")
	status, err := strconv.ParseBool(statusStr)
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	// bu yerde middleware - den gelen musderinin id - si alynyar
	custID, hasCustomer := c.Get("customer_id")
	if !hasCustomer {
		helpers.HandleError(c, 400, "customer_id is required")
		return
	}
	customerID, ok := custID.(string)
	if !ok {
		helpers.HandleError(c, 400, "customer_id must be string")
		return
	}

	// front - dan gelen maglumatlar bind edilyar
	var productIds ProductID
	if err := c.BindJSON(&productIds); err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	if status { // eger status = true bolsa halanlarym sahyp haryt gosulyar
		if len(productIds.IDS) != 0 { // eger front - dan gelen id bar bolsa onda halanlarym sahypa haryt gosulyar
			for _, v := range productIds.IDS {
				// front - dan gelen id - lere den gelyan bazada haryt barmy yokmy sol barlanyar
				var product_id string
				if err := db.QueryRow(context.Background(), "SELECT id FROM products WHERE id = $1 AND deleted_at IS NULL", v).Scan(&product_id); err != nil {
					helpers.HandleError(c, 400, err.Error())
					return
				}

				if product_id != "" { // eger haryt products tablida bar bolsa onda sol haryt on gelen musderinin
					// halanlarynyn arasynda barmy ya-da yokmy sol barlanyar
					var product string
					db.QueryRow(context.Background(), "SELECT product_id FROM likes WHERE customer_id = $1 AND product_id = $2 AND deleted_at IS NULL", customerID, v).Scan(&product)

					if product == "" { // eger haryt musderinin halanlarym harytlarynyn arasynda yok bolsa
						// onda haryt sol musderinin halanlarym tablisasyna gosulyar
						_, err := db.Exec(context.Background(), "INSERT INTO likes (customer_id,product_id) VALUES ($1,$2)", customerID, v)
						if err != nil {
							helpers.HandleError(c, 400, err.Error())
							return
						}
					}
				}
			}

			// front - dan gelen harytlar halanlarym sahypa gosulandan son
			// yzyna sol harytlar ddoly maglumatlary bilen berilyar
			products, err := GetLikes(customerID)
			if err != nil {
				helpers.HandleError(c, 400, err.Error())
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"status":   true,
				"products": products,
			})
		} else { // eger front hic hile id gelmese onda musderinin onki bazadaky halan harytlaryny fronta bermeli
			products, err := GetLikes(customerID)
			if err != nil {
				helpers.HandleError(c, 400, err.Error())
				return
			}

			if len(products) == 0 {
				helpers.HandleError(c, 400, "like empty")
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"status":   true,
				"products": products,
			})
		}
	} else { // eger status = false gelse onda front - dan gele id - li harydy sol musderinin halanlarym harytlaryndan pozmaly
		if len(productIds.IDS) != 0 { // front - dan maglumat gelyarmi ya-da gelenokmy sony barlayas
			// front - dan gelen id - ler halanlarym tablisada barmy ya-da yokmy sony barlayas
			var product_id string
			if err := db.QueryRow(context.Background(), "SELECT product_id FROM likes WHERE customer_id = $1 AND product_id = $2 AND deleted_at IS NULL", customerID, productIds.IDS[0]).Scan(&product_id); err != nil {
				helpers.HandleError(c, 400, err.Error())
				return
			}

			// eger haryt halanlarym tablisada yok bolsa
			// yzyna yalnyslyk goyberyas
			if product_id == "" {
				c.JSON(http.StatusNotFound, gin.H{
					"status":  false,
					"message": "this product not found in this customer",
				})
				return
			}

			// haryt halanlarym tablisada bar bolsa onda ony pozyas
			_, err = db.Exec(context.Background(), "DELETE FROM likes WHERE customer_id = $1 AND product_id = $2 AND deleted_at IS NULL", customerID, productIds.IDS[0])
			if err != nil {
				helpers.HandleError(c, 400, err.Error())
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"status":  true,
				"message": "like successfull deleted",
			})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  false,
				"message": "product id is required",
			})
		}

	}
}

func GetLikes(customerID string) ([]models.Product, error) {
	db, err := config.ConnDB()
	if err != nil {
		return []models.Product{}, err
	}
	defer db.Close()

	rowsProduct, err := db.Query(context.Background(),
		"SELECT p.id,p.name_tm,p.name_ru,p.image,p.price,p.old_price,p.status,p.color_name_tm,p.color_name_ru,p.gender_name_tm,p.gender_name_ru,p.code,p.shop_id,p.category_id,p.brend_id FROM products p INNER JOIN likes l ON l.product_id = p.id WHERE l.customer_id = $1 AND l.deleted_at IS NULL AND p.deleted_at IS NULL",
		customerID)
	if err != nil {
		return []models.Product{}, err
	}
	defer rowsProduct.Close()

	var products []models.Product
	for rowsProduct.Next() {
		var product models.Product
		if err := rowsProduct.Scan(&product.ID,
			&product.NameTM,
			&product.NameRU,
			&product.Image,
			&product.Price,
			&product.OldPrice,
			&product.Status,
			&product.ColorNameTM,
			&product.ColorNameRU,
			&product.GenderNameTM,
			&product.ColorNameRU,
			&product.Code,
			&product.ShopID,
			&product.CategoryID,
			&product.BrendID,
		); err != nil {
			return []models.Product{}, err
		}

		products = append(products, product)
	}

	return products, nil
}

func GetCustomerLikes(c *gin.Context) {
	custID, hasCustomer := c.Get("customer_id")
	if !hasCustomer {
		helpers.HandleError(c, 400, "customer_id is required")
		return
	}
	customerID, ok := custID.(string)
	if !ok {
		helpers.HandleError(c, 400, "customer_id must be string")
		return
	}

	products, err := GetLikes(customerID)
	if err != nil {
		helpers.HandleError(c, 400, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":   true,
		"products": products,
	})
}
