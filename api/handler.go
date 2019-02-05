package api

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/globalsign/mgo/bson"
	_ "payment-service/docs"
	"payment-service/logger"
	"payment-service/model"
	"payment-service/repository"
	"payment-service/service"
	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"mime"
	"net/http"
)

const (
	ContentType    = "Content-Type"
	DatabaseName   = "PaymentDB"
	CollectionName = "Payment"
	ID             = "id"
)

// PaymentHandler the card payment handler
type PaymentHandler struct {
	repo repository.Repository
	fx   service.FXService
	ch   service.ChargesService
}

// NewPaymentHandler creates a type of CardPaymentHandler
func NewPaymentHandler(repo repository.Repository, fxUrl string, chUrl string) *PaymentHandler {
	return &PaymentHandler{repo, service.NewFxService(fxUrl), service.NewChargesService(chUrl)}
}

//----------------------------------------------------------------------------------------
//							Endpoint handlers
//----------------------------------------------------------------------------------------

// Health the health endpoint handler
func (h *PaymentHandler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, model.HealthResponse{Message: "Up and running!"})
}

// @Summary Creates new payment
// @ID create-payment
// @Description Creates new payment
// @Accept  json
// @Produce  json
// @Param new-tag body model.CreatePaymentRequest true "New tag"
// @Success 201 {object} model.CreatePaymentResponse "Tag created"
// @Failure 400 {object} model.ErrorResponse "Bad request"
// @Failure 500 {object} model.ErrorResponse "Internal server error"
// @Router /payment [post]
func (h *PaymentHandler) CreatePayment(c *gin.Context) {
	var req model.CreatePaymentRequest

	// This will infer what binder to use depending on the content-type header.
	if errB := c.ShouldBindWith(&req, binding.JSON); errB != nil {
		logger.Error.Println(errB.Error())
		setErrorResponse("Failed to parse payment request", http.StatusBadRequest, c)
		return
	}

	logger.Info.Printf("Received request to create payment for organisationId: %s", req.OrganisationID)

	// Get exchange rate from the mock service
	fx := model.ForeignExchange{ExchangeRate: 1.0}

	if foreignExchangeRequired(req) {
		_, fx = h.fx.GetExchangeRate(req.BeneficiaryParty.Currency, req.DebtorParty.Currency, req.Amount)
	}

	// calculate the new amount based on the exchange rate
	amount := getAmount(req.Amount, fx.ExchangeRate)

	// Get charges information from the mock service
	_, charges := h.ch.GetCharges(fx.ExchangeRate, req.BearerCode, req.BeneficiaryParty.Currency, req.DebtorParty.Currency)

	// persisting payment into database
	payment := buildPayment(amount, fx, charges, req)
	logger.Info.Printf("Storing payment with ID %s", payment.ID.Hex())
	err := h.repo.Insert(DatabaseName, CollectionName, payment)

	if err != nil {
		logger.Error.Println(err.Error())
		setErrorResponse("Failed to create payment", http.StatusInternalServerError, c)
		return
	}

	// if all good create success response
	c.Writer.Header().Set(ContentType, mime.TypeByExtension("json"))
	c.JSON(http.StatusCreated, model.CreatePaymentResponse{ID: payment.ID.Hex(), OrganisationId: payment.OrganisationId})
}

// @Summary Get all payments
// @ID get-payments
// @Accept  json
// @Produce  json
// @Success 200 {object} model.PaymentResponse	"ok"
// @Failure 400 {object} model.ErrorResponse "Bad request"
// @Failure 500 {object} model.ErrorResponse "Internal server error"
// @Router /payment [get]
func (h *PaymentHandler) FindAllPayments(c *gin.Context) {
	logger.Info.Println("Received request to query all payments")
	resp, err := h.repo.FindAll(DatabaseName, CollectionName)

	if err != nil {
		logger.Error.Println(err.Error())
		setErrorResponse("Failed to query payments", http.StatusInternalServerError, c)
		return
	}

	// if all good create success response
	c.Writer.Header().Set(ContentType, mime.TypeByExtension("json"))
	c.JSON(http.StatusOK, resp)
}

// @Summary Get a payment for given ID
// @ID get-payment
// @Accept  json
// @Produce  json
// @Success 200 {object} model.PaymentResponse	"ok"
// @Failure 400 {object} model.ErrorResponse "Bad request"
// @Failure 404 {object} model.ErrorResponse "Not found"
// @Failure 500 {object} model.ErrorResponse "Internal server error"
// @Router /payment/{id} [get]
func (h *PaymentHandler) FindPayment(c *gin.Context) {

	id := c.Params.ByName(ID)
	logger.Info.Printf("Received request to query a payment for a given ID %s", id)
	resp, err := h.repo.Find(DatabaseName, CollectionName, bson.ObjectIdHex(id))

	if err != nil {
		c.JSON(http.StatusNotFound, model.EmptyBody{})
		return
	}

	// if all good create success response
	c.Writer.Header().Set(ContentType, mime.TypeByExtension("json"))
	c.JSON(http.StatusOK, resp)
}

// @Summary Delete a payment for given ID
// @ID delete-payment
// @Accept  json
// @Produce  json
// @Success 204 "Payment deleted"
// @Failure 400 {object} model.ErrorResponse "Bad request"
// @Failure 404 {object} model.ErrorResponse "Not found"
// @Failure 500 {object} model.ErrorResponse "Internal server error"
// @Router /payment/{id} [delete]
func (h *PaymentHandler) DeletePayment(c *gin.Context) {
	id := c.Params.ByName(ID)
	logger.Info.Printf("Received request to delete a payment for a given ID %s", id)

	// query the payment first
	_, errQ := h.repo.Find(DatabaseName, CollectionName, bson.ObjectIdHex(id))
	if errQ != nil {
		c.JSON(http.StatusNotFound, model.EmptyBody{})
		return
	}

	err := h.repo.Delete(DatabaseName, CollectionName, bson.ObjectIdHex(id))

	if err != nil {
		setErrorResponse("Failed to delete payment", http.StatusInternalServerError, c)
		return
	}

	// if all good create success response
	logger.Info.Printf("Payment with id [%s] successfully deleted", id)
	c.Status(http.StatusNoContent)
}

// @Summary Update a payment for given ID - partial payment is not supported
// @ID update-payment
// @Accept  json
// @Produce  json
// @Success 204 "Payment updated"
// @Failure 400 {object} model.ErrorResponse "Bad request"
// @Failure 404 {object} model.ErrorResponse "Not found"
// @Failure 500 {object} model.ErrorResponse "Internal server error"
// @Router /payment/{id} [put]
func (h *PaymentHandler) UpdatePayment(c *gin.Context) {
	var req model.CreatePaymentRequest

	id := c.Params.ByName(ID)

	// This will infer what binder to use depending on the content-type header.
	if errB := c.ShouldBindWith(&req, binding.JSON); errB != nil {
		logger.Error.Println(errB.Error())
		setErrorResponse("Failed to parse update payment request", http.StatusBadRequest, c)
		return
	}

	logger.Info.Printf("Received request to update payment for payment ID: %s", id)

	// Get exchange rate from the mock service
	fx := model.ForeignExchange{ExchangeRate: 1.0}

	if foreignExchangeRequired(req) {
		_, fx = h.fx.GetExchangeRate(req.BeneficiaryParty.Currency, req.DebtorParty.Currency, req.Amount)
	}

	// calculate the new amount based on the exchange rate
	amount := getAmount(req.Amount, fx.ExchangeRate)

	// Get charges information from the mock service
	_, charges := h.ch.GetCharges(fx.ExchangeRate, req.BearerCode, req.BeneficiaryParty.Currency, req.DebtorParty.Currency)

	// persisting payment into database
	payment := updatePayment(amount, fx, charges, req, bson.ObjectIdHex(id))
	logger.Info.Printf("Updating payment with ID %s", payment.ID.Hex())
	err := h.repo.Update(DatabaseName, CollectionName, bson.ObjectIdHex(id), payment)
	if err != nil {
		logger.Error.Println(err.Error())
		setErrorResponse("Failed to update payment", http.StatusNotFound, c)
		return
	}

	// if all good create success response
	c.Status(http.StatusNoContent)
}

//----------------------------------------------------------------------------------------
//							Initialise the router
//----------------------------------------------------------------------------------------

// NewRouter creates an instance of the router
func (h *PaymentHandler) NewRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// configure all the route
	router.GET("/health", h.Health)
	router.POST("/payment", h.CreatePayment)
	router.GET("/payment", h.FindAllPayments)
	router.GET("/payment/:id", h.FindPayment)
	router.DELETE("/payment/:id", h.DeletePayment)
	router.PUT("/payment/:id", h.UpdatePayment)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return router
}

// helper function
func setErrorResponse(msg string, status int, c *gin.Context) {
	c.JSON(status, model.ErrorResponse{Message: msg, Code: status})
}

// Helper function to calculate the new amount based on the given exchange rate
func getAmount(amount float64, rate float64) float64 {
	return amount / rate
}

// Helper function to build payment instance
func buildPayment(amount float64, fx model.ForeignExchange, charges model.ChargesInformation, req model.CreatePaymentRequest) model.Payment {

	attr := buildAttr(amount, req, charges, fx)

	return model.Payment{Type: "Payment", ID: bson.NewObjectId(), OrganisationId: req.OrganisationID, Attributes: attr, Version: 0}
}

// Helper function to build payment instance
func updatePayment(amount float64, fx model.ForeignExchange, charges model.ChargesInformation, req model.CreatePaymentRequest, oid bson.ObjectId) model.Payment {
	attr := buildAttr(amount, req, charges, fx)
	return model.Payment{Type: "Payment", ID: oid, OrganisationId: req.OrganisationID, Attributes: attr, Version: 0}
}

// Helper function to determine if foreign exchange to this payment is relevant
func foreignExchangeRequired(req model.CreatePaymentRequest) bool {
	return req.DebtorParty.Currency != req.BeneficiaryParty.Currency
}

// Helper function to build payment attributes
func buildAttr(amount float64, req model.CreatePaymentRequest, charges model.ChargesInformation, fx model.ForeignExchange) model.Attributes {
	attr := model.Attributes{Amount: amount, BeneficiaryParty: req.BeneficiaryParty, DebtorParty: req.DebtorParty,
		ChargesInformation: charges, Currency: req.BeneficiaryParty.Currency, EndToEndReference: req.EndToEndReference,
		Fx: fx, NumericReference: req.NumericReference, PaymentID: req.PaymentID, PaymentPurpose: req.PaymentPurpose,
		PaymentScheme: req.PaymentScheme, PaymentType: req.PaymentType, ProcessingDate: req.ProcessingDate, Reference: req.Reference,
		SchemePaymentSubType: req.SchemePaymentSubType, SchemePaymentType: req.SchemePaymentType, SponsorParty: req.SponsorParty}
	return attr
}
