package jsonduplicatechecker

import (
	"encoding/json"
	"strconv"

	"github.com/PoorMercymain/user-segmenter/errors"
)

func CheckDuplicatesInJSON(d *json.Decoder, path []string) error {
	t, err := d.Token()
	if err != nil {
		return err
	}

	delim, ok := t.(json.Delim)

	if !ok {
		return nil
	}

	switch delim {
	case '{':
		keys := make(map[string]bool)
		for d.More() {
			t, err := d.Token()
			if err != nil {
				return err
			}
			key := t.(string)

			if keys[key] {
				return errors.ErrorDuplicateInJSON
			}
			keys[key] = true

			if err := CheckDuplicatesInJSON(d, append(path, key)); err != nil {
				return err
			}
		}

		if _, err := d.Token(); err != nil {
			return err
		}

	case '[':
		i := 0
		for d.More() {
			if err := CheckDuplicatesInJSON(d, append(path, strconv.Itoa(i))); err != nil {
				return err
			}
			i++
		}

		if _, err := d.Token(); err != nil {
			return err
		}

	}
	return nil
}
