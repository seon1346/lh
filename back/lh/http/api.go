package http

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"lh/fabric"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

var doc_count uint64 = 0
var file_path = "/home/bstudent/doc"

func initial(c *gin.Context) {
	setCorsHeader(c)

	req := NewRequest(c.Request, nil)
	req.Params["name"] = c.Params.ByName("name")
	req.Params["id"] = c.Params.ByName("id")
	rsp := HandleinitialRequest(req)

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

func HandleinitialRequest(request *Request) *Response {

	name := request.Params["name"]
	id := request.Params["id"]

	problemDetails := initialProcedure(name, id)

	if problemDetails == nil {
		return NewResponse(http.StatusOK, nil, nil)
	} else if problemDetails != nil {
		return NewResponse(int(problemDetails.Status), nil, problemDetails)
	}

	problemDetails = &ProblemDetails{
		Status: http.StatusForbidden,
		Cause:  "UNSPECIFIED",
	}

	return NewResponse(http.StatusForbidden, nil, problemDetails)
}
func initialProcedure(name, id string) *ProblemDetails {
	fmt.Printf("Handle initialProcedure\n")
	err := fabric.InitUser(fmt.Sprintf("%s-%s", name, id), name, "initial", "User1_org1", "User1_org2")
	if err != nil {

		return &ProblemDetails{
			Status: 500,
			Cause:  err.Error(),
		}
	}
	fmt.Printf("Handle initialProcedure end\n")
	return nil
}

func uploadRightdoc(c *gin.Context) {
	setCorsHeader(c)

	req := NewRequest(c.Request, c.Request.Body)
	req.Params["name"] = c.Params.ByName("name")
	req.Params["id"] = c.Params.ByName("id")
	req.Params["docname"] = c.Params.ByName("docname")
	rsp := HandleuploadRightdocRequest(req)

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

func HandleuploadRightdocRequest(request *Request) *Response {

	log.Printf("HandleuploadRightdocRequest\n")
	name := request.Params["name"]
	id := request.Params["id"]
	docname := request.Params["docname"]
	problemDetails := uploadRightdocProcedure(name, id, docname, request.Body)

	if problemDetails == nil {
		return NewResponse(http.StatusOK, nil, nil)
	} else if problemDetails != nil {
		return NewResponse(int(problemDetails.Status), nil, problemDetails)
	}

	problemDetails = &ProblemDetails{
		Status: http.StatusForbidden,
		Cause:  "UNSPECIFIED",
	}

	return NewResponse(http.StatusForbidden, nil, problemDetails)
}
func uploadRightdocProcedure(name, id, docname string, Body interface{}) *ProblemDetails {
	fmt.Printf("Handle uploadRightdocProcedure\n")
	b, ok := Body.(io.ReadCloser)
	if !ok {
		return &ProblemDetails{
			Status: 500,
		}
	}
	filecontent, err := ioutil.ReadAll(b)
	if err != nil {

		return &ProblemDetails{
			Status: 500,
			Cause:  err.Error(),
		}
	}
	key := fmt.Sprintf("%s-%s", name, id)
	fmt.Println(key, "   ", id)
	savepath := fmt.Sprintf("%s/%s/%s", file_path, key, docname)
	file, err := os.OpenFile(savepath, os.O_CREATE|os.O_RDWR, os.FileMode(0777))
	if err != nil && strings.Contains(err.Error(), "no such file or directory") {
		os.MkdirAll(fmt.Sprintf("%s/%s", file_path, key), os.FileMode(0777))
		file, err = os.OpenFile(savepath, os.O_CREATE|os.O_RDWR, os.FileMode(0777))
	}

	if err != nil {
		return &ProblemDetails{
			Status: 500,
			Cause:  err.Error(),
		}
	}
	defer file.Close()
	file.Write(filecontent)
	h := sha256.New()
	h.Write(filecontent)
	hash := hex.EncodeToString(h.Sum(nil))
	err = fabric.UpdateRightDoc(key, name, docname, fmt.Sprintf("%d", doc_count), hash, "User1_org2")
	if err != nil {

		return &ProblemDetails{
			Status: 500,
			Cause:  err.Error(),
		}
	}
	doc_count++
	return nil
}

func getRight(c *gin.Context) {
	setCorsHeader(c)

	req := NewRequest(c.Request, c.Request.Body)
	req.Params["name"] = c.Params.ByName("name")
	req.Params["id"] = c.Params.ByName("id")
	rsp := HandlegetRightRequest(req)

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

func HandlegetRightRequest(request *Request) *Response {

	log.Printf("HandleuploadRightdocRequest\n")
	name := request.Params["name"]
	id := request.Params["id"]

	ret, problemDetails := getRightProcedure(name, id)

	if problemDetails == nil {
		return NewResponse(http.StatusOK, nil, ret)
	} else if problemDetails != nil {
		return NewResponse(int(problemDetails.Status), nil, problemDetails)
	}

	problemDetails = &ProblemDetails{
		Status: http.StatusForbidden,
		Cause:  "UNSPECIFIED",
	}

	return NewResponse(http.StatusForbidden, nil, problemDetails)
}
func getRightProcedure(name, id string) (*RightProcess, *ProblemDetails) {
	fmt.Printf("Handle uploadRightdocProcedure\n")
	key := fmt.Sprintf("%s-%s", name, id)
	bs, err := fabric.GetRight(key, "User1_org1")
	if err != nil {

		return nil, &ProblemDetails{
			Status: 500,
			Cause:  err.Error(),
		}
	}
	var ret *RightProcess

	json.Unmarshal(bs, &ret)
	return ret, nil
}

func changeRightstate(c *gin.Context) {
	setCorsHeader(c)

	req := NewRequest(c.Request, c.Request.Body)
	req.Params["name"] = c.Params.ByName("name")
	req.Params["id"] = c.Params.ByName("id")
	req.Params["statevalue"] = c.Params.ByName("statevalue")
	rsp := HandlechangeRightstate(req)

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

func HandlechangeRightstate(request *Request) *Response {

	log.Printf("Handle changeRightstate\n")
	name := request.Params["name"]
	id := request.Params["id"]
	statevalue := request.Params["statevalue"]
	problemDetails := changeRightstateProcedure(name, id, statevalue)

	if problemDetails == nil {
		return NewResponse(http.StatusOK, nil, nil)
	} else if problemDetails != nil {
		return NewResponse(int(problemDetails.Status), nil, problemDetails)
	}

	problemDetails = &ProblemDetails{
		Status: http.StatusForbidden,
		Cause:  "UNSPECIFIED",
	}

	return NewResponse(http.StatusForbidden, nil, problemDetails)
}
func changeRightstateProcedure(name, id, statevalue string) *ProblemDetails {
	fmt.Printf("changeRightstateProcedure\n")
	key := fmt.Sprintf("%s-%s", name, id)
	err := fabric.UpdateRightState(key, statevalue, "User1_org1")
	if err != nil {

		return &ProblemDetails{
			Status: 500,
			Cause:  err.Error(),
		}
	}

	return nil
}
