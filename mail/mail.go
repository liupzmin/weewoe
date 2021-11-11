package mail

type Report struct {
	Title              string
	ReportName         string
	ReportDate         string
	Reporter           string
	NormalCount        int
	PlanStopCount      int
	PortExceptionCount int
	PIDExceptionCount  int
	Groups             []Group
}

type Group struct {
	Name      string
	Processes []Process
}

type Process struct {
	Name      string
	Host      string
	Ports     []Port
	State     int64
	StartTime string
	Suspend   bool
}

type Port struct {
	Number string
	State  int64
}

func SendMail() error {
	return nil
}
