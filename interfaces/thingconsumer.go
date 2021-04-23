package interfaces

type ThingConsumer interface {
	ConsumeThings() ([]string, error)
}
