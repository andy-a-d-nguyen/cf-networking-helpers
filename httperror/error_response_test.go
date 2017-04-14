package httperror_test

import (
	"errors"
	"net/http"
	"net/http/httptest"

	"code.cloudfoundry.org/go-db-helpers/fakes"
	"code.cloudfoundry.org/go-db-helpers/httperror"
	"code.cloudfoundry.org/lager/lagertest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
)

var _ = Describe("ErrorResponse", func() {
	var (
		errorResponse     *httperror.ErrorResponse
		logger            *lagertest.TestLogger
		fakeMetricsSender *fakes.MetricsSender
		resp              *httptest.ResponseRecorder
		err               error
	)

	BeforeEach(func() {
		logger = lagertest.NewTestLogger("test")
		fakeMetricsSender = &fakes.MetricsSender{}
		errorResponse = &httperror.ErrorResponse{
			Logger:        logger,
			MetricsSender: fakeMetricsSender,
		}
		resp = httptest.NewRecorder()
		err = errors.New("potato")
	})

	Describe("Internal Server Error", func() {
		It("Logs the error", func() {
			errorResponse.InternalServerError(resp, err, "message", "description")
			Expect(logger).To(gbytes.Say("message: description.*potato"))
		})
		It("responds with an error body and status code 500", func() {
			errorResponse.InternalServerError(resp, err, "message", "description")
			Expect(resp.Code).To(Equal(http.StatusInternalServerError))
			Expect(resp.Body.String()).To(MatchJSON(`{"error": "message: description"}`))
		})
		It("increments the counter", func() {
			errorResponse.InternalServerError(resp, err, "message", "description")
			Expect(fakeMetricsSender.IncrementCounterCallCount()).To(Equal(1))
			Expect(fakeMetricsSender.IncrementCounterArgsForCall(0)).To(Equal("message"))
		})
	})

	Describe("Bad Request", func() {
		It("Logs the error", func() {
			errorResponse.BadRequest(resp, err, "message", "description")
			Expect(logger).To(gbytes.Say("message: description.*potato"))
		})
		It("responds with an error body and status code 400", func() {
			errorResponse.BadRequest(resp, err, "message", "description")
			Expect(resp.Code).To(Equal(http.StatusBadRequest))
			Expect(resp.Body.String()).To(MatchJSON(`{"error": "message: description"}`))
		})
	})

	Describe("Forbidden", func() {
		It("Logs the error", func() {
			errorResponse.Forbidden(resp, err, "message", "description")
			Expect(logger).To(gbytes.Say("message: description.*potato"))
		})
		It("responds with an error body and status code 403", func() {
			errorResponse.Forbidden(resp, err, "message", "description")
			Expect(resp.Code).To(Equal(http.StatusForbidden))
			Expect(resp.Body.String()).To(MatchJSON(`{"error": "message: description"}`))
		})
	})

	Describe("Unauthorized", func() {
		It("Logs the error", func() {
			errorResponse.Unauthorized(resp, err, "message", "description")
			Expect(logger).To(gbytes.Say("message: description.*potato"))
		})
		It("responds with an error body and status code 401", func() {
			errorResponse.Unauthorized(resp, err, "message", "description")
			Expect(resp.Code).To(Equal(http.StatusUnauthorized))
			Expect(resp.Body.String()).To(MatchJSON(`{"error": "message: description"}`))
		})
	})

	Describe("Conflict", func() {
		It("Logs the error", func() {
			errorResponse.Conflict(resp, err, "message", "description")
			Expect(logger).To(gbytes.Say("message: description.*potato"))
		})
		It("responds with an error body and status code 409", func() {
			errorResponse.Conflict(resp, err, "message", "description")
			Expect(resp.Code).To(Equal(http.StatusConflict))
			Expect(resp.Body.String()).To(MatchJSON(`{"error": "message: description"}`))
		})
	})
})
