package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type JLaporanCRequest struct {
	UserId    int
	DateStart string
	DateEnd   string
	PageNow   int
	RowPage   int
}

type JLaporanCResponse struct {
	MerchantName string
}

func LaporanC(c *gin.Context) {
	db := DB()
	defer db.Close()

	var (
		reqBody     JLaporanCRequest
		resBody     JLaporanCResponse
		LaporanList []JLaporanCResponse
		bodyBytes   []byte
	)

	if c.Request.Body != nil {
		bodyBytes, _ = ioutil.ReadAll(c.Request.Body)
	}
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	bodyString := string(bodyBytes)

	// ------ Body Json Validation ------
	if string(bodyString) == "" {
		errorMessage := "Error, Body is empty"
		returnDataJsonLaporanC("1", errorMessage, LaporanList, 0, 0, 0, c)
		return
	}

	is_Json := isJSON(bodyString)
	if is_Json == false {
		errorMessage := "Error, Body - invalid json data"
		returnDataJsonLaporanC("1", errorMessage, LaporanList, 0, 0, 0, c)
		return
	}
	// ------ end of Body Json Validation ------

	// ------ Header Validation ------
	if ValidateHeader(bodyString, c) {
		if err := c.ShouldBindJSON(&reqBody); err != nil {
			errorMessage := "Error, Bind Json Data"
			returnDataJsonLaporanC("1", errorMessage, LaporanList, 0, 0, 0, c)
			return
		} else {
			UserId := reqBody.UserId
			DateStart := reqBody.DateStart
			DateEnd := reqBody.DateEnd
			PageNow := reqBody.PageNow
			RowPage := reqBody.RowPage
			errorMessage := ""

			// ------ Param Validation ------
			if UserId == 0 {
				errorMessage = errorMessage + "\n- " + "User Id cannot be null"
			}

			if DateStart == "" {
				errorMessage = errorMessage + "\n- " + "Date Start cannot be null"
			}

			if DateEnd == "" {
				errorMessage = errorMessage + "\n- " + "Date End cannot be null"
			}

			if errorMessage != "" {
				returnDataJsonLaporanC("1", errorMessage, LaporanList, 0, 0, 0, c)
				return
			}
			// ------ end of Param Validation ------

			TotalRecords := 0
			TotalPage := 0.0
			query := fmt.Sprintf("select count(1) as cnt from Transactions a, Merchants b, Users d where a.merchant_id = b.id and b.user_id = d.id and b.user_id = %d and a.created_at between '%s' and '%s';", UserId, DateStart, DateEnd)
			if err := db.QueryRow(query).Scan(&TotalRecords); err != nil {
				errorMessage := "Error query, " + err.Error()
				returnDataJsonLaporanC("1", errorMessage, LaporanList, 0, 0, 0, c)
				return
			}
			TotalPage = math.Ceil(float64(TotalRecords) / float64(RowPage))

			Omzet := 0
			queryOmzet := fmt.Sprintf("select IFNULL(sum(bill_total), 0) as omzet from Transactions a, Merchants b, Users d where a.merchant_id = b.id and b.user_id = d.id and b.user_id = %d and a.created_at between '%s' and '%s';", UserId, DateStart, DateEnd)
			if err := db.QueryRow(queryOmzet).Scan(&Omzet); err != nil {
				errorMessage := "Error query, " + err.Error()
				returnDataJsonLaporanC("1", errorMessage, LaporanList, 0, 0, 0, c)
				return
			}

			query1 := fmt.Sprintf("select b.merchant_name from Transactions a, Merchants b, Users d where a.merchant_id = b.id and b.user_id = d.id and b.user_id = %d and a.created_at between '%s' and '%s' LIMIT %d,%d", UserId, DateStart, DateEnd, PageNow, RowPage)
			rows, err := db.Query(query1)
			if err != nil {
				errorMessage := "Error query, " + err.Error()
				returnDataJsonLaporanC("1", errorMessage, LaporanList, 0, 0, 0, c)
				return
			}
			for rows.Next() {
				err = rows.Scan(
					&resBody.MerchantName,
				)
				LaporanList = append(LaporanList, resBody)
				if err != nil {
					errorMessage := "Error execute, " + err.Error()
					returnDataJsonLaporanC("1", errorMessage, LaporanList, 0, 0, 0, c)
					return
				}
			}
			defer rows.Close()
			if LaporanList != nil {
				returnDataJsonLaporanC("0", "", LaporanList, Omzet, int(TotalPage), TotalRecords, c)
				return
			} else {
				errorMessage = "Data not found"
				returnDataJsonLaporanC("1", errorMessage, LaporanList, 0, 0, 0, c)
				return
			}
		}
	}
}

func returnDataJsonLaporanC(ErrorCode string, ErrorMessage string, LaporanList []JLaporanCResponse, Omzet int, TotalPage int, TotalRecords int, c *gin.Context) {
	if strings.Contains(ErrorMessage, "Error running") {
		ErrorMessage = "Error Execute data"
	}

	if ErrorCode == "504" {
		c.String(http.StatusUnauthorized, "")
	} else {
		c.PureJSON(http.StatusOK, gin.H{
			"ErrCode":      ErrorCode,
			"ErrMessage":   ErrorMessage,
			"Result":       LaporanList,
			"Omzet":        Omzet,
			"TotalPages":   TotalPage,
			"TotalRecords": TotalRecords,
		})
	}
	return
}
