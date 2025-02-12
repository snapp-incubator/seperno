package seperno

import (
	"github.com/snapp-incubator/seperno/internal"
	"github.com/snapp-incubator/seperno/pkg/options"
)

func NewNormalize(ops ...options.Options) Normalize {
	opts := options.DefaultOptions
	for _, config := range ops {
		config.Apply(&opts)
	}
	return internal.NewNormalizer(opts)
}

func WithConvertHalfSpaceToSpace() options.Options {
	return options.NewFuncOption(func(option *options.NormalizerOptions) {
		option.ConvertHalfSpaceToSpace = true
	})
}

func WithSpaceCombiner() options.Options {
	return options.NewFuncOption(func(option *options.NormalizerOptions) {
		option.SpaceCombiner = true
	})
}

func WithOuterSpaceRemover() options.Options {
	return options.NewFuncOption(func(option *options.NormalizerOptions) {
		option.OuterSpaceRemover = true
	})
}

func WithURLRemover() options.Options {
	return options.NewFuncOption(func(option *options.NormalizerOptions) {
		option.URLRemover = true
	})
}

func WithNormalizePunctuations() options.Options {
	return options.NewFuncOption(func(option *options.NormalizerOptions) {
		option.NormalizePunctuations = true
	})
}

func WithEndsWithEndOfLineChar() options.Options {
	return options.NewFuncOption(func(option *options.NormalizerOptions) {
		option.EndsWithEndOfLineChar = true
	})
}

// WithIntToWord do not use WithConvertNumberToLanguage after use this option
func WithIntToWord() options.Options {
	return options.NewFuncOption(func(option *options.NormalizerOptions) {
		option.IntToWord = true
		option.ConvertNumberLang = options.LanguageEn
	})
}

// WithConvertNumberToLanguage default language is "en" , options are : "en" , "fa" , "ar"
func WithConvertNumberToLanguage(language options.Language) options.Options {
	return options.NewFuncOption(func(options *options.NormalizerOptions) {
		options.ConvertNumberLang = language
	})
}

type Normalize interface {
	FindHalfSpace(input, halfSpace string) string
	BasicNormalizer(input string) string
	VariationSelectorsRemover(input []string) []string
	BasicNormalizerArray(input []string) []string
	BasicNormalizerSlice(input []string) []string
}
