package domain

import (
	"time"

	"github.com/stackus/errors"
	"github.com/jfelipeforero/iparking/internal/ddd"
)

const BookingAggregate = "booking.Reserve"

var (
	ErrCustomerIDCannotBeBlank = errors.Wrap(errors.ErrBadRequest, "the customer id cannot be blank.") 	
	ErrPaymentIDCannotBeBlank = errors.Wrap(errors.ErrBadRequest, "the payment id cannot be blank")
	ErrStartDateMustBeAfterCurrentTime = errors.Wrap(errors.ErrBadRequest, "the start date of the reserve must be after current time.")
	ErrEndDateMustBeAfterStartTime = errors.Wrap(errors.ErrBadRequest, "the end date of the reserve must be after start date ")
	ErrLocationMustValid = errors.Wrap(errors.ErrBadRequest, "the location must be an existing parking location")
)

type Reserve struct {
	//es.Aggregate
	CustomerID string
	PaymentID  string
	InvoiceID  string
	Location   string
	StartDate  time.Time
	EndDate    time.Time
	Status	   string
}

// Check interface implementations
var _ interface {

}

func NewReserve(id string) *Reserve {
	return &Reserve{
		
	}
}

func (r *Reserve) CreateReserve(id, customerID, paymentID, invoiceID, location, startDate, endDate, Status) error {
	if r.CustomerID == "" {
		return ErrCustomerIDCannotBeBlank
	}
	if r.PaymentID == "" {
		return ErrPaymentIDCannotBeBlank
	}
	if r.Location == "" {
		return ErrLocationMustValid
	}
}
