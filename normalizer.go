package seperno

import (
	"C"

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
	return options.NewFuncOption(func(options *options.NormalizerOptions) {
		options.ConvertHalfSpaceToSpace = true
	})
}

func WithSpaceCombiner() options.Options {
	return options.NewFuncOption(func(options *options.NormalizerOptions) {
		options.SpaceCombiner = true
	})
}

func WithOuterSpaceRemover() options.Options {
	return options.NewFuncOption(func(options *options.NormalizerOptions) {
		options.OuterSpaceRemover = true
	})
}

func WithURLRemover() options.Options {
	return options.NewFuncOption(func(options *options.NormalizerOptions) {
		options.URLRemover = true
	})
}

func WithNormalizePunctuations() options.Options {
	return options.NewFuncOption(func(options *options.NormalizerOptions) {
		options.NormalizePunctuations = true
	})
}

func WithEndsWithEndOfLineChar() options.Options {
	return options.NewFuncOption(func(options *options.NormalizerOptions) {
		options.EndsWithEndOfLineChar = true
	})
}

func WithIntToWord() options.Options {
	return options.NewFuncOption(func(options *options.NormalizerOptions) {
		options.IntToWord = true
	})
}

type Normalize interface {
	FindHalfSpace(input, halfSpace string) string
	BasicNormalizer(input string) string
	VariationSelectorsRemover(input []string) []string
	BasicNormalizerArray(input []string) []string
	BasicNormalizerSlice(input []string) []string
}
