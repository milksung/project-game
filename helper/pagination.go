package helper

import "errors"

func Pagination(page, limit *int) error {

	if *page == 0 {
		*page = 1
	}

	if *limit == 0 {
		*limit = 10
	}

	if *page < 1 {
		return errors.New("Page must be greater than 0")
	}

	if *limit < 1 {
		return errors.New("Limit must be greater than 0")
	}

	if *page > 0 {
		*page--
	}

	return nil
}

func UnlimitPagination(page, limit *int) error {

	if *page <= 0 {
		*page = 1
	}

	if *limit < 0 {
		*limit = 0
	}

	if *page > 0 {
		*page--
	}
	return nil
}
