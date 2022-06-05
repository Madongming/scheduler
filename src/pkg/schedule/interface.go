package schedule

type Jober interface {
	Run() error
	Stop() error
}
