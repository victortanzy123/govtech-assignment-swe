package controller

import (
	"bytes"
	"strings"
	"net/http"
	"net/http/httptest"
	"testing"
    _ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"log"
)

 /*///////////////////////////////////////////////////////////////
                Student Registration By Teacher
    //////////////////////////////////////////////////////////////*/

// @Desc: [VALID] Registering 3 students by a teacher with correct input format and should succeed with HTTP Code 204.
func TestRegisterStudentValid1(t *testing.T) {

	var jsonBody = []byte(`{"teacher": "t1@gmail.com","students":["s1@gmail.com","s2@gmail.com","s3@gmail.com"]}`)
	req, err := http.NewRequest("POST", "/api/register", bytes.NewBuffer(jsonBody))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(RegisterStudents)
	handler.ServeHTTP(rr, req)
	status := rr.Code

	// Check the response body is what we expect.
	expected := ``
	actual := strings.TrimRight(rr.Body.String(), "\n")


	assert.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, status, "Status code should be 204")
	assert.Equal(t, expected, actual, "Response should be empty.")
	log.Println("SUCCESS: TestRegisterStudentValid1")
}


// @Desc: [VALID] Registering the same 3 students by different teacher with correct input format and should succeed with HTTP Code 204.
func TestRegisterStudentValid2(t *testing.T) {

	var jsonBody = []byte(`{"teacher": "t2@gmail.com","students":["s1@gmail.com","s2@gmail.com","s3@gmail.com"]}`)
	req, err := http.NewRequest("POST", "/api/register", bytes.NewBuffer(jsonBody))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(RegisterStudents)
	handler.ServeHTTP(rr, req)
	status := rr.Code

	// Check the response body is what we expect.
	expected := ``
	actual := strings.TrimRight(rr.Body.String(), "\n")


	assert.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, status, "Status code should be 204")
	assert.Equal(t, expected, actual, "Response should be empty.")
	log.Println("SUCCESS: TestRegisterStudentValid2")
}



// @Desc: [FAIL] Registering a student without specifying a teacher field resulting in incorrect input format and should fail with HTTP Code 400.
func TestRegisterStudentWithInvalidBodyFormat(t *testing.T) {

	var jsonBody = []byte(`{"students":["s1@gmail.com"]}`)
	req, err := http.NewRequest("POST", "/api/register", bytes.NewBuffer(jsonBody))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(RegisterStudents)
	handler.ServeHTTP(rr, req)
	status := rr.Code

	// Check the response body is what we expect.
	expected := `{"message":"Missing teacher specified."}`
	actual := strings.TrimRight(rr.Body.String(), "\n")


	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, status, "Status code should be 400")
	assert.Equal(t, expected, actual, "Response should be empty.")
	log.Println("SUCCESS: TestRegisterStudentWithInvalidBodyFormat")
}


// @Desc: [FAIL] Registering the a student by a teacher that has been previously registered with correct input format and should fail with HTTP Code 409 (Status Conflict).
func TestRegisterSameStudentUnderSameTeacher(t *testing.T) {

	var jsonBody = []byte(`{"teacher": "t1@gmail.com","students":["s1@gmail.com"]}`)
	req, err := http.NewRequest("POST", "/api/register", bytes.NewBuffer(jsonBody))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(RegisterStudents)
	handler.ServeHTTP(rr, req)
	status := rr.Code

	// Check the response body is what we expect.
	expected := `{"message":"Student has been registered previously."}`
	actual := strings.TrimRight(rr.Body.String(), "\n")


	assert.NoError(t, err)
	assert.Equal(t, http.StatusConflict, status, "Status code should be 409")
	assert.Equal(t, expected, actual, "Response should be empty.")
	log.Println("SUCCESS: TestRegisterSameStudentUnderSameTeacher")
}


 /*///////////////////////////////////////////////////////////////
                	Fetching Common Students
    //////////////////////////////////////////////////////////////*/

// @Desc: [VALID] Retrieving common students under one teacher with correct input format and should succeed with HTTP Code 200.
func TestGetCommonStudents1Teacher(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/commonstudents", nil)
	if err != nil {
		t.Fatal(err)
	}

	q := req.URL.Query()
	q.Add("teacher", "t1@gmail.com")
	req.URL.RawQuery = q.Encode()
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(CommonStudents)
	handler.ServeHTTP(rr, req)
	status := rr.Code

	// Check the response body is what we expect.
	expected := `["s1@gmail.com","s2@gmail.com","s3@gmail.com"]`
    actual := strings.TrimRight(rr.Body.String(), "\n")


	assert.NoError(t, err)
	assert.Equal(t, status, http.StatusOK, "Status code should be 200")
	assert.Equal(t, expected, actual, "Response should be the same as expected.")

	log.Println("SUCCESS: TestGetCommonStudents1Teacher")
}

// @Desc: [VALID] Retrieving common students under two teachers with correct input format and should succeed with HTTP Code 200.
func TestGetCommonStudents2Teachers(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/commonstudents", nil)
	if err != nil {
		t.Fatal(err)
	}

	q := req.URL.Query()
	q.Add("teacher", "t1@gmail.com")
	q.Add("teacher", "t2@gmail.com")
	req.URL.RawQuery = q.Encode()
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(CommonStudents)
	handler.ServeHTTP(rr, req)
	status := rr.Code

	// Check the response body is what we expect.
	expected := `["s1@gmail.com","s2@gmail.com","s3@gmail.com"]`
    actual := strings.TrimRight(rr.Body.String(), "\n")


	assert.NoError(t, err)
	assert.Equal(t, status, http.StatusOK, "Status code should be 200")
	assert.Equal(t, expected, actual, "Response should be the same as expected.")

	log.Println("SUCCESS: TestGetCommonStudents2Teachers")
}

// @Desc: [FAIL] Retrieving common students without specifying teacher in the query parameters, resulting in an incorrect input format and should fail with HTTP Code 400.
func TestGetCommonStudentsInvalidTeacherQuery(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/commonstudents", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(CommonStudents)
	handler.ServeHTTP(rr, req)
	status := rr.Code

	// Check the response body is what we expect.
	expected := `{"message":"No teacher specified."}`
    actual := strings.TrimRight(rr.Body.String(), "\n")


	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, status, "Status code should be 400")
	assert.Equal(t, expected, actual, "Response should be the same as expected.")

	log.Println("SUCCESS: TestGetCommonStudentsInvalidTeacherQuery")
}


 /*///////////////////////////////////////////////////////////////
                	Fetching Common Students
    //////////////////////////////////////////////////////////////*/

// @Desc: [VALID] Suspending a student with the correct input format, which should result in HTTP code 204.
func TestSuspendStudent(t *testing.T) {

	var jsonBody = []byte(`{"student": "s1@gmail.com"}`)
	req, err := http.NewRequest("POST", "/api/suspend", bytes.NewBuffer(jsonBody))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(SuspendStudent)
	handler.ServeHTTP(rr, req)
	status := rr.Code

	// Check the response body is what we expect.
	expected := ``
	actual := strings.TrimRight(rr.Body.String(), "\n")


	assert.NoError(t, err)
	assert.Equal(t, status, http.StatusNoContent, "Status code should be 204")
	assert.Equal(t, expected, actual, "Response should be empty.")
	log.Println("SUCCESS: TestSuspendStudent")
}

// @Desc: [FAIL] Suspending a student with incorrect input format, which should fail and result in HTTP code 400 (Status Bad Request).
func TestSuspendStudentWithInvalidBodyFormat(t *testing.T) {

	var jsonBody = []byte(``)
	req, err := http.NewRequest("POST", "/api/suspend", bytes.NewBuffer(jsonBody))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(SuspendStudent)
	handler.ServeHTTP(rr, req)
	status := rr.Code

	// Check the response body is what we expect.
	expected := `{"message":"Invalid Request Body Format."}`
	actual := strings.TrimRight(rr.Body.String(), "\n")


	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, status, "Status code should be 400")
	assert.Equal(t, expected, actual, "Response should be the same.")
	log.Println("SUCCESS: TestSuspendStudentWithInvalidBodyFormat")
}

// @Desc: [FAIL] Suspending an already suspended student with the correct input format, which should fail and result in HTTP code 409 (Status Conflict).
func TestSuspendAnExistingSuspendedStudent(t *testing.T) {

	var jsonBody = []byte(`{"student": "s1@gmail.com"}`)
	req, err := http.NewRequest("POST", "/api/suspend", bytes.NewBuffer(jsonBody))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(SuspendStudent)
	handler.ServeHTTP(rr, req)
	status := rr.Code

	// Check the response body is what we expect.
	expected := `{"message":"Student has been suspended previously."}`
	actual := strings.TrimRight(rr.Body.String(), "\n")


	assert.NoError(t, err)
	assert.Equal(t, http.StatusConflict, status, "Status code should be 409")
	assert.Equal(t, expected, actual, "Response should be empty.")
	log.Println("SUCCESS: TestSuspendAnExistingSuspendedStudent")
}



 /*///////////////////////////////////////////////////////////////
                	Retrieve For Notifications
    //////////////////////////////////////////////////////////////*/

// @Desc: [VALID] Retrieving students for notifications with correct input format by a teacher, which should result students who are registered previously, and not suspended to appear in the student list output. This query should succeed with HTTP code 200.
func TestRetrieveForNotications(t *testing.T) {
	var jsonBody = []byte(`{
		"teacher": "t1@gmail.com",
		"notification": "hello world bye @s1@gmail.com @s2@gmail.com @s3@gmail.com"
	}`)
	req, err := http.NewRequest("POST", "/api/retrievefornotifications", bytes.NewBuffer(jsonBody))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(RetrieveForNotification)
	handler.ServeHTTP(rr, req)
	status := rr.Code

	// Check the response body is what we expect.
	expected := `{"teacher":"t1@gmail.com","notification":"hello world bye","students":["s2@gmail.com","s3@gmail.com"]}`
	actual := strings.TrimRight(rr.Body.String(), "\n")


	assert.NoError(t, err)
	assert.Equal(t, status, http.StatusOK, "Status code should be 204")
	assert.Equal(t, expected, actual, "Response should be empty.")
	log.Println("SUCCESS: TestRetrieveForNotification")
}


// @Desc: [FAIL] Retrieving students for notifications with incorrect input format with a missing request body. This query should fail with HTTP code 400.
func TestRetrieveForNoticationsEmptyBody(t *testing.T) {
	// Empty request body
	var jsonBody = []byte(``)
	req, err := http.NewRequest("POST", "/api/retrievefornotifications", bytes.NewBuffer(jsonBody))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(RetrieveForNotification)
	handler.ServeHTTP(rr, req)
	status := rr.Code

	// Check the response body is what we expect.
	expected := `{"message":"Invalid request body format."}`
	actual := strings.TrimRight(rr.Body.String(), "\n")


	assert.NoError(t, err)
	assert.Equal(t, status, http.StatusBadRequest, "Status code should be 400")
	assert.Equal(t, expected, actual, "Response should be empty.")
	log.Println("SUCCESS: TestRetrieveForNotificationEmptyBody")
}


// @Desc: [FAIL] Retrieving students for notifications with incorrect input format with an incomplete request body, missing the notification field. This query should fail with HTTP code 400.
func TestRetrieveForNoticationsMissingNotificationField(t *testing.T) {
	// Empty request body
	var jsonBody = []byte(`{"teacher": "t5@gmail.com"}`)
	req, err := http.NewRequest("POST", "/api/retrievefornotifications", bytes.NewBuffer(jsonBody))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(RetrieveForNotification)
	handler.ServeHTTP(rr, req)
	status := rr.Code

	// Check the response body is what we expect.
	expected := `{"message":"Empty teacher or notification format."}`
	actual := strings.TrimRight(rr.Body.String(), "\n")


	assert.NoError(t, err)
	assert.Equal(t, status, http.StatusBadRequest, "Status code should be 400")
	assert.Equal(t, expected, actual, "Response should be empty.")
	log.Println("SUCCESS: TestRetrieveForNotificationMissingNotificationField")
}




