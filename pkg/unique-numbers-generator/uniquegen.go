package uniquenumbersgenerator

import (
	"math/rand"
	"time"

	"github.com/PoorMercymain/user-segmenter/errors"
)

func GenerateUniqueNonNegativeNumbers(amount int, rightLimit int) (map[int]struct{}, error) {
	r := rand.New(rand.NewSource(time.Now().Unix()))

	if rightLimit <= 0 {
		return nil, errors.ErrorInvalidRightLimit
	}

	if rightLimit < amount {
		return nil, errors.ErrorRightLimitIsTooLow
	}

	uniqueNumbers := make(map[int]struct{}, 0)

	for i := 0; i < amount; i++ {
		someInt := r.Intn(rightLimit)
		_, ok := uniqueNumbers[someInt]
		for ok {
			someInt = r.Intn(rightLimit)
			_, ok = uniqueNumbers[someInt]
		}
		uniqueNumbers[someInt] = struct{}{}
	}

	return uniqueNumbers, nil
}
