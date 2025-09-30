package seperno

import "github.com/snapp-incubator/seperno/pkg/lfd"

func NewPersianNumberDetector() lfd.NumberDetector {
	return &lfd.PersianNumberDetector{}
}
