package model
 



 /*///////////////////////////////////////////////////////////////
                Database Queries Related Struct
    //////////////////////////////////////////////////////////////*/

type Teach struct {
    Teacher string `form:"teacher" json:"teacher"`
	Student string `form:"student" json:"student"`
}

type Suspend struct {
    Teacher string `form:"teacher" json:"teacher"`
	Student string `form:"student" json:"student"`
}

type Notification struct {
    Teacher string `form:"teacher" json:"teacher"`
	Student string `form:"student" json:"student"`
}

type StudentRegistration struct {
    Teacher string `json:"teacher"`
    Students []string `json:"students"`
}

type CommonStudents struct {
    Students []string `json:"students"`
}

type SuspendStudent struct {
    Student string `json:"student"`
}

type Student struct {
    StudentEmail string `form:"student" json:"student"`
}

 /*///////////////////////////////////////////////////////////////
                API Response related structs
    //////////////////////////////////////////////////////////////*/

type RetrieveForNotificationBody struct {
    Teacher string `json:"teacher"`
    Notification string `json:"notification"`
}

type RetrieveForNotificationResponse struct {
    Teacher string `json:"teacher"`
    Notification string `json:"notification"`
    Students []string `json:"students"`
}


type MessageResponse struct {
    Message string `json:"message"`
}


