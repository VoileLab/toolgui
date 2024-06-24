package component

import (
	"crypto/md5"
	"fmt"
	"math/rand"
)

func hashedID(componentName string, bs []byte) string {
	return fmt.Sprintf("%s_%x", componentName, md5.Sum(bs))
}

func normalID(componentName, label string) string {
	return fmt.Sprintf("%s_%s", componentName, label)
}

func randID(componentName string) string {
	return fmt.Sprintf("%s_%d", componentName, rand.Int())
}
