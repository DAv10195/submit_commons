package http

const (
	limit					= "limit"
	afterId					= "after_id"

	base					= 10
	bitSize					= 64

	ElementsLeftToProcess	= "X-Elements-Left-To-Process"
	ForSubmitUser			= "X-Submit-User"
	ForSubmitCourse			= "X-Submit-Course"
	ForSubmitAss			= "X-Submit-Ass"
	SubmitState				= "X-Submit-State"

	SubmitSessionRoles			= "X-Submit-Session-Roles"
	SubmitSessionStaffCourses 	= "X-Submit-Session-Staff-Courses"
	SubmitSessionStudentCourses = "X-Submit-Session-Student-Courses"
	SubmitSessionUser			= ForSubmitUser

	AppealStateOpen			= "open"
	AppealStateClosed		= "closed"

	FsPlaceHolderFileName	= ".submit"
)
