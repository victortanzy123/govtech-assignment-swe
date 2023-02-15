package controller

import (
    "encoding/json"
    "fmt"
    "net/mail"
    "strings"
    "net/http"
    "database/sql"
 
   "github.com/victortanzy123/govtech-assignment-swe/model"
   "github.com/victortanzy123/govtech-assignment-swe/config"   
)

 /*///////////////////////////////////////////////////////////////
                            Main Functions
//////////////////////////////////////////////////////////////*/



// CommonStudents: Get All Common
// URL : /commonstudents
// Parameters: teacher
// Method: GET
// Output: JSON Encoded Array of student emails if found else JSON Encoded Exception.
func CommonStudents(w http.ResponseWriter, r *http.Request) {
    var teach model.Teach
    var commonStudentsList []string

    // Retrieve all teachers from query params
    teachers := r.URL.Query()["teacher"]
    if len(teachers) == 0 {
        ErrorResponse("No teacher specified.", w, http.StatusBadRequest)
        return
    }
 
    db := config.Connect()
    defer db.Close()

    // Build query string based on how many teachers has been specified:
    var commonStudentsQuery string
    commonStudentsQuery = getCommonStudentsQuery(teachers)
 
    rows, err := db.Query(commonStudentsQuery)
 
    if err != nil {
        ErrorResponse("Failed to get common students", w, http.StatusNotFound)
        return
    }
 
    // If successful, process the students into an array string
    for rows.Next() {
        err = rows.Scan(&teach.Student)
        if err != nil {
            ErrorResponse("Failed to get common students", w, http.StatusNotFound)
            return
        } else {
            commonStudentsList = append(commonStudentsList, teach.Student)
        }
    }

    // If no students, initialise an empty array as output
    if len(commonStudentsList) == 0 {
        commonStudentsList = make([]string, 0)
    }

 
    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    json.NewEncoder(w).Encode(commonStudentsList)
}

// RegisterStudents: Register students under a teacher
// URL : /register
// Parameters: teacher, students
// Method: POST
// Output: No content if successful, else error message.
func RegisterStudents(w http.ResponseWriter, r *http.Request) {
    var studentRegistration model.StudentRegistration
 
    db := config.Connect()
    defer db.Close()
 
    err := json.NewDecoder(r.Body).Decode(&studentRegistration)
    if err != nil {
        ErrorResponse("Failed request body format.", w, http.StatusBadRequest)
        return
    }

  
    var teacher string = studentRegistration.Teacher
    var students []string = studentRegistration.Students

    if len(teacher) == 0 {
        ErrorResponse("Missing teacher specified.", w, http.StatusBadRequest)
        return
    }


    // Insert students one by one
    for _, student := range students {

        rows, err := db.Query("SELECT * FROM Teach WHERE teacher  = ? AND student = ?",teacher, student)
        if err != nil {
            ErrorResponse("Failed to register students.", w, http.StatusBadRequest)
            return
        }
    
        count := 0
        for rows.Next() {
            count++
        }
        if count > 0 {
            ErrorResponse("Student has been registered previously.", w, http.StatusConflict)
            return
        }
        _, err = db.Exec("INSERT INTO Teach(teacher, student) VALUES(?, ?)", teacher,student)
        if err != nil {
           ErrorResponse("Failed to register students", w, http.StatusNotFound)
           return
        }
    }

 
    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.WriteHeader(http.StatusNoContent)
}


// SuspendStudent: Suspend a student
// URL : /suspend
// Parameters: student
// Method: POST
// Output: No content if successful, else error message.
func SuspendStudent(w http.ResponseWriter, r *http.Request) {
    var suspendStudent model.SuspendStudent
 
    db := config.Connect()
    defer db.Close()

    err := json.NewDecoder(r.Body).Decode(&suspendStudent)
    if err != nil {
        ErrorResponse("Invalid Request Body Format.", w, http.StatusBadRequest)
        return
    }

    var student string = suspendStudent.Student

    // Insert student to Suspend Table
    rows, err := db.Query("SELECT * FROM Suspend WHERE student = ?", student)
    if err != nil {
        ErrorResponse("Failed to suspend student.", w, http.StatusNotFound)
        return
    }

    count := 0
    for rows.Next() {
        count++
    }
    if count > 0 {
        ErrorResponse("Student has been suspended previously.", w, http.StatusConflict)
        return
    }

    _, err = db.Exec("INSERT INTO Suspend(student) VALUES(?)", student)
    if err != nil {
        ErrorResponse("Failed to suspend student", w, http.StatusNotFound)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.WriteHeader(http.StatusNoContent)
}


// RetrieveForNotification: Retrieve and send notifications to a list of registered, notified and non-suspended students by a teacher
// URL : /retrievefornotification
// Parameters: teacher, notification
// Method: POST
// Output: JSON Encoded Object of teacher, notification and list of students notified.
func RetrieveForNotification(w http.ResponseWriter, r *http.Request) {
    var requestBody model.RetrieveForNotificationBody
 
    db := config.Connect()
    defer db.Close()


    err := json.NewDecoder(r.Body).Decode(&requestBody)
    if err != nil {
        ErrorResponse("Invalid request body format.", w, http.StatusBadRequest)
        return
    }

    var teacher string = requestBody.Teacher
    var notification string = requestBody.Notification
    words := strings.Split(notification, " ")
    // If minimally no words in notification string
    if len(words[0]) == 0 {
        ErrorResponse("Empty teacher or notification format.", w, http.StatusBadRequest)
        return
    }

    var notificationWords []string
    var emails []string
    // Run a for loop until an '@' sign is found in the string, that marks the first email
    for _, word := range words {
        var isEmail bool = isEmailRegisterFormat(word)
        if !isEmail {
            notificationWords = append(notificationWords, word)
        } else {
            studentEmail := word[1:]
            emails = append(emails, studentEmail)
        }
    }
    
    // Add student emails to notification
   if len(emails) > 0 {
        for _, email := range emails {
            // Check if the email has been registered under the teacher:
            err = insertStudentIfNotExists(db, teacher, email)
            if err != nil {
                ErrorResponse("Failed to register student emails for notifications", w, http.StatusNotFound)
                return
            }
        }
   }

    var viewName string = "RegisteredStudentsForNotifications"
   err = dropViewIfExists(db, viewName)
   if err != nil {
        ErrorResponse("Failed to retrieve notifications.", w, http.StatusNotFound)
        return;
   }

   // 1. Create VIEW of all students under the specified teacher that are registered for notifications & also t side of :=gopreviously registered under teaching.
   var viewQuery string = fmt.Sprintf(`CREATE VIEW RegisteredStudentsForNotifications AS 
   SELECT Notification.student
   FROM Teach, Notification 
   WHERE Teach.teacher = Notification.teacher AND Teach.student = Notification.student AND Notification.teacher = '%s'`, teacher)
   _, err = db.Exec(viewQuery)
   if err != nil {
        ErrorResponse("Failed to retrieve notifications.", w, http.StatusNotFound)
        return
    }

   // 2. Retrieve all students from RegisteredStudentsForNotifcations - SuspendedStudents
   rows, err := db.Query("SELECT student FROM RegisteredStudentsForNotifications WHERE student NOT IN (SELECT student FROM Suspend)")
   if err != nil {
        ErrorResponse("Failed to retrieve students for notifications.", w, http.StatusNotFound)
        return
    }
    
    var studentStruct model.Student
    var notificationResponse model.RetrieveForNotificationResponse
    notificationResponse.Teacher = teacher
    notificationResponse.Notification = strings.Join(notificationWords, " ")
 
    for rows.Next() {
        err = rows.Scan(&studentStruct.StudentEmail)
        notificationResponse.Students = append(notificationResponse.Students, studentStruct.StudentEmail)
    }

    if len(notificationResponse.Students) == 0 {
        notificationResponse.Students =  make([]string, 0)// initialize to empty slice
    }


    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    json.NewEncoder(w).Encode(notificationResponse)
}



  /*///////////////////////////////////////////////////////////////
                        Error/Success Responses
    //////////////////////////////////////////////////////////////*/

func ErrorResponse( message string, w http.ResponseWriter, httpCode int) {
    var response model.MessageResponse

    response.Message = message
    w.WriteHeader(httpCode)
    json.NewEncoder(w).Encode(response)

}

func WriteSuccessResponse(message string, w http.ResponseWriter) {
    var response model.MessageResponse

    response.Message = message
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(response)
   
}

 /*///////////////////////////////////////////////////////////////
                        Helper Functions
    //////////////////////////////////////////////////////////////*/

func validEmailFormat(email string) bool {
    _, err := mail.ParseAddress(email)
    return err == nil
}

func isEmailRegisterFormat(text string) bool {
    if len(text) > 0 && text[0] =='@' {
        return true
    }
    return false
}

// @Desc:[CommonStudent] To dynamically build the query based on the teachers specified in `CommonStudent` Query
func getCommonStudentsQuery(teachers []string) string {
    var sqlPlaceholders []string
    for _, teacher := range teachers {
        sqlPlaceholders = append(sqlPlaceholders, fmt.Sprintf("'%s'", teacher))
    }
    // Concate the placeholder strings based on length of teachers in the input
    placeholdersStr := strings.Join(sqlPlaceholders, ", ")
    
    return fmt.Sprintf("SELECT student FROM Teach WHERE teacher IN (%s) GROUP BY student HAVING COUNT(DISTINCT teacher) = %d", placeholdersStr, len(teachers))
}

// @Desc: [RegisterStudents] To register a student into the table if only the teacher-student pair has not been registered yet.
func insertStudentIfNotExists(db *sql.DB, teacherEmail string,  studentEmail string) error {
    rows, err := db.Query("SELECT * FROM Notification WHERE teacher  = ? AND student = ?",teacherEmail, studentEmail)
    if err != nil {
        return err
    }

    count := 0
    for rows.Next() {
        count++
    }
    if count > 0 {
        return err
    }

    query := `INSERT INTO Notification (teacher, student) SELECT ?, ? WHERE NOT EXISTS (SELECT 1 FROM Notification WHERE teacher = ? AND student = ?)`
    _, err = db.Exec(query, teacherEmail, studentEmail, teacherEmail, studentEmail)
    return err
}


// @Desc: [RetrieveForNotification] Helper function to drop intermediate `RegisteredStudentsForNotifications` View table if exists.
func dropViewIfExists(db *sql.DB, viewName string) error {
    var exists bool
    query := fmt.Sprintf("SELECT COUNT(*) FROM information_schema.views WHERE table_name = '%s'", viewName)
    err := db.QueryRow(query).Scan(&exists)
    if err != nil {
        return err
    }
    if exists {
        _, err := db.Exec(fmt.Sprintf("DROP VIEW %s", viewName))
        if err != nil {
            return err
        }
    }
    return nil
}