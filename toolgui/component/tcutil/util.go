package tcutil

import (
	"crypto/md5"
	"fmt"
	"math/rand"
)

func HashedID(componentName string, bs []byte) string {
	return fmt.Sprintf("%s_%x", componentName, md5.Sum(bs))
}

func NormalID(componentName, label string) string {
	return fmt.Sprintf("%s_%s", componentName, label)
}

func RandID(componentName string) string {
	return fmt.Sprintf("%s_%d", componentName, rand.Int())
}
