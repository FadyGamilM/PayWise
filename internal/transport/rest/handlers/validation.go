package handlers

import (
	"log"
	"paywise/internal/core"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type InvalidReqBody struct {
	Field string `json:"field"`
	Value string `json:"value"`
	Tag   string `json:"tag"`
	Param string `json:"param"`
}

/*
@Params

	c : context from the gin request to have an access to the req body
	req : the variable (dto type) where we will bind the request body into it

@Returns

	bool value, true if the req.body is valid and false if its not valid
*/
func BindReqBody(c *gin.Context, req interface{}) bool {
	// 1. check the content type
	if !checkReqContentType(c) {
		err := core.NewBadRequestError()
		c.JSON(err.StatusCode(), gin.H{
			"error": err,
		})
		return false
	}

	// 2. try to bind the request body to our passed req dto instance
	// i will use shoouldBind not bindjson because i need to check the content type first  and stop if any thing happend while binding
	if err := c.ShouldBind(req); err != nil {
		errs, bound := err.(validator.ValidationErrors)
		// if bound, thats means we have a validation error
		if bound {
			var invalidReqArgs []InvalidReqBody
			for _, e := range errs {
				invalidReqArgs = append(invalidReqArgs, InvalidReqBody{
					Field: e.Field(),
					Value: e.Value().(string),
					Tag:   e.Tag(),
					Param: e.Param(),
				})
			}
			// now return the proper error response
			AppErr := core.NewBadRequestError()
			c.JSON(AppErr.StatusCode(), gin.H{
				"error":        AppErr,
				"invalid_args": invalidReqArgs,
			})
			return false
		}
		// if it can't be bound, so its internal server error
		log.Printf("validation error in request body")
		err := core.NewInternalServerError()
		c.JSON(err.StatusCode(), gin.H{
			"error": err,
		})
		return false
	}
	return true
}

func checkReqContentType(c *gin.Context) bool {
	// check the content type, if its not json return bad request error
	return c.ContentType() == "application/json"
}
