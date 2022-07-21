package http

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func setCorsHeader(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
	c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, PATCH, DELETE")
}

func test(c *gin.Context) {
	setCorsHeader(c)

	req := NewRequest(c.Request, nil)

	rsp := HandletestRequest(req)

	responseBody, err := json.Marshal(rsp.Body)
	if err != nil {
		log.Println(err)
		problemDetails := ProblemDetails{
			Status: http.StatusInternalServerError,
			Cause:  "SYSTEM_FAILURE",
			Detail: err.Error(),
		}
		c.JSON(http.StatusInternalServerError, problemDetails)
	} else {
		c.Data(rsp.Status, "application/json", responseBody)
	}
}

func HandletestRequest(request *Request) *Response {
	// step 1: log
	log.Printf("test")

	// step 2: retrieve request
	supi := request.Params["supi"]

	// step 3: handle the message
	response, problemDetails := testProcedure(supi)

	// step 4: process the return value from step 3
	if response != nil {
		// status code is based on SPEC, and option headers
		return NewResponse(http.StatusOK, nil, response)
	} else if problemDetails != nil {
		return NewResponse(int(problemDetails.Status), nil, problemDetails)
	}

	problemDetails = &ProblemDetails{
		Status: http.StatusForbidden,
		Cause:  "UNSPECIFIED",
	}

	return NewResponse(http.StatusForbidden, nil, problemDetails)
}
func testProcedure(supi string) (response *int, problemDetails *ProblemDetails) {
	log.Printf("Handle OAM Get Am Policy")
	//Convert
	//Diameter send
	//return
	return
}
