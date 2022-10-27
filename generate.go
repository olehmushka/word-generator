package wordgenerator

func Generate(opts GenerateOpts) (string, error) {
	return New().Generate(opts)
}
