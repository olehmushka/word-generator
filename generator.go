package wordgenerator

type Generator interface {
	Generate(opts GenerateOpts) (string, error)
}
