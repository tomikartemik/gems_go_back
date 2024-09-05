package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type ReplInput struct {
	Amount float64 `json:"amount"`
	Promo  string  `json:"promo"`
	I      int     `json:"i"`
	Ip     string  `json:"ip"`
}

func (h *Handler) NewReplenishment(c *gin.Context) {
	var input ReplInput
	userID := c.GetString("user_id")
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	fmt.Printf("input: %+v\n", input)
	location, err := h.services.Replenishment.NewReplenishment(userID, input.Amount, input.Promo, input.I, input.Ip)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, location)
}

func (h *Handler) RedirectAccepted(c *gin.Context) {
	queryParams := c.Request.URL.Query()
	fmt.Printf("queryParams: %+v\n", queryParams)

	replenishmentIdStr := c.Query("MERCHANT_ORDER_ID")
	replenishmentId, _ := strconv.Atoi(replenishmentIdStr)
	go h.services.AcceptReplenishment(replenishmentId)
	fmt.Printf("replenishmentId: %d\n", replenishmentId)
	c.Redirect(http.StatusFound, "https://dododrop.ru")
}

func (h *Handler) RedirectDenied(c *gin.Context) {
	queryParams := c.Request.URL.Query()

	for key, values := range queryParams {
		for _, value := range values {
			c.String(http.StatusOK, "Key: %s, Value: %s\n", key, value)
		}
	}
	fmt.Printf("queryParams: %+v\n", queryParams)
	c.Redirect(http.StatusFound, "https://dododrop.ru")
}

func (h *Handler) MSGFromFreekassa(c *gin.Context) {
	merchantID := c.PostForm("MERCHANT_ID")
	amount := c.PostForm("AMOUNT")
	intid := c.PostForm("intid")
	orderID := c.PostForm("MERCHANT_ORDER_ID")
	email := c.PostForm("P_EMAIL")
	phone := c.PostForm("P_PHONE")
	curID := c.PostForm("CUR_ID")
	sign := c.PostForm("SIGN")
	usKey := c.PostForm("us_key")
	payerAccount := c.PostForm("payer_account")
	commission := c.PostForm("commission")

	// Вывод полученных данных на консоль
	fmt.Printf("MERCHANT_ID: %s\n", merchantID)
	fmt.Printf("AMOUNT: %s\n", amount)
	fmt.Printf("intid: %s\n", intid)
	fmt.Printf("MERCHANT_ORDER_ID: %s\n", orderID)
	fmt.Printf("P_EMAIL: %s\n", email)
	fmt.Printf("P_PHONE: %s\n", phone)
	fmt.Printf("CUR_ID: %s\n", curID)
	fmt.Printf("SIGN: %s\n", sign)
	fmt.Printf("us_key: %s\n", usKey)
	fmt.Printf("payer_account: %s\n", payerAccount)
	fmt.Printf("commission: %s\n", commission)

	c.JSON(http.StatusOK, "OK")
	//c.Redirect(http.StatusFound, "https://dododrop.ru")
}
