package options

var DefaultOptions = NormalizerOptions{
	ConvertHalfSpaceToSpace: false,
	URLRemover:              false,
	OuterSpaceRemover:       false,
	SpaceCombiner:           false,
	NormalizePunctuations:   false,
	EndsWithEndOfLineChar:   false,
	IntToWord:               false,
	ConvertNumberLang:       LanguageEn,
}

type Language string

const (
	LanguageFa Language = "fa"
	LanguageAr Language = "ar"
	LanguageEn Language = "en"
)

type NormalizerOptions struct {
	ConvertHalfSpaceToSpace bool
	URLRemover              bool
	OuterSpaceRemover       bool
	SpaceCombiner           bool
	NormalizePunctuations   bool
	EndsWithEndOfLineChar   bool
	IntToWord               bool
	ConvertNumberLang       Language
}

type Options interface {
	Apply(options *NormalizerOptions)
}

type FuncConfig struct {
	ops func(options *NormalizerOptions)
}

func (w FuncConfig) Apply(conf *NormalizerOptions) {
	w.ops(conf)
}

func NewFuncOption(f func(options *NormalizerOptions)) *FuncConfig {
	return &FuncConfig{ops: f}
}
