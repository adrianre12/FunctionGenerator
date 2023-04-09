package ui

type Field interface {
	Write(text string)
	LineWidth() (innerWidth uint32, outboxWidth uint32)
	Invert(invert bool)
}
