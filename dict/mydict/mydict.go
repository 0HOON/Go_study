package mydict

import "errors"

// Dictionary type
type Dictionary map[string]string

var (
	errNotFount   = errors.New("not fount")
	errCantUpdate = errors.New("can't update non-existing word")
	errWordExists = errors.New("that word already exists")
	errCantDelete = errors.New("can't delete non-existing word")
)

// Search for a word
func (d Dictionary) Search(word string) (string, error) {
	value, exists := d[word]
	if exists {
		return value, nil
	}
	return "", errNotFount
}

// Add a word to the dictionary
func (d Dictionary) Add(word, def string) error {
	_, err := d.Search(word)
	switch err {
	case errNotFount:
		d[word] = def
	case nil:
		return errWordExists
	}
	return nil
}

// Update definition of a word
func (d Dictionary) Update(word, def string) error {
	_, err := d.Search(word)
	switch err {
	case errNotFount:
		return errCantUpdate
	case nil:
		d[word] = def
	}
	return nil
}

// Delete a word from dictionary
func (d Dictionary) Delete(word string) error {
	_, err := d.Search(word)
	switch err {
	case errNotFount:
		return errCantDelete
	case nil:
		delete(d, word)
	}
	return nil
}
