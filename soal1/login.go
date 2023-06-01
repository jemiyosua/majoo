package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type JLogin struct {
	Username string
	Password string
}

func Login(c *gin.Context) {
	db := DB()
	defer db.Close()

	var (
		reqBody   JLogin
		bodyBytes []byte
	)

	if c.Request.Body != nil {
		bodyBytes, _ = ioutil.ReadAll(c.Request.Body)
	}
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	bodyString := string(bodyBytes)

	// ------ Body Json Validation ------
	if string(bodyString) == "" {
		errorMessage := "Error, Body is empty"
		returnDataJsonLogin("1", errorMessage, c)
		return
	}

	is_Json := isJSON(bodyString)
	if is_Json == false {
		errorMessage := "Error, Body - invalid json data"
		returnDataJsonLogin("1", errorMessage, c)
		return
	}
	// ------ end of Body Json Validation ------

	// ------ Header Validation ------
	if ValidateHeader(bodyString, c) {
		if err := c.ShouldBindJSON(&reqBody); err != nil {
			errorMessageReturn := "Error, Bind Json Data"
			errorMessage := err.Error()
			returnDataJsonLogin("1", errorMessageReturn+" | "+errorMessage, c)
			return
		} else {
			Username := reqBody.Username
			Password := reqBody.Password
			errorMessage := ""

			// ------ Param Validation ------
			if Username == "" {
				errorMessage = errorMessage + "\n- " + "Username cannot be null"
			}

			if Password == "" {
				errorMessage = errorMessage + "\n- " + "Password cannot be null"
			}

			if errorMessage != "" {
				returnDataJsonLogin("1", errorMessage, c)
				return
			}
			// ------ end of Param Validation ------

			Cnt := 0
			query := "SELECT count(1) AS cnt FROM users where user_name = '" + Username + "' and password = '" + GetMD5Hash(Password) + "' "
			if err := db.QueryRow(query).Scan(&Cnt); err != nil {
				errorMessage := "Error query, " + err.Error()
				returnDataJsonLogin("1", errorMessage, c)
				return
			}

			if Cnt > 0 {
				returnDataJsonLogin("0", "", c)
			} else {
				errorMessage = "Username or Password invalid, please try again!"
				returnDataJsonLogin("1", errorMessage, c)
			}
		}
	}
}

func returnDataJsonLogin(ErrorCode string, ErrorMessage string, c *gin.Context) {

	if strings.Contains(ErrorMessage, "Error running") {
		ErrorMessage = "Error Execute data"
	}

	if ErrorCode == "504" {
		c.String(http.StatusUnauthorized, "")
	} else {
		c.PureJSON(http.StatusOK, gin.H{
			"ErrorCode":    ErrorCode,
			"ErrorMessage": ErrorMessage,
		})
	}

	return
}
