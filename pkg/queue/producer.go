package queue

type Producer interface {
	AddListener(listener ProduceListener)
	Produce() (string,bool)
}

type ProduceListener interface {
	OnProducerPause()
	OnProducerResume()
}

type ProducerFactory func() (Producer, error)
