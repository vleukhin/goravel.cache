package cache

import (
	"strings"

	"github.com/mintance/go-uniqid"
)

// store is a collection of method to access to store
type TagSet struct {
	store Store
	names []string
}

func (ts *TagSet) reset() error {
	for i := 0; i != len(ts.names); i++ {
		_, err := ts.resetTag(ts.names[i])
		if err != nil {
			return err
		}
	}

	return nil
}

func (ts *TagSet) generateId() string {
	uniqId := uniqid.New(uniqid.Params{
		Prefix:      "",
		MoreEntropy: true,
	})
	id := strings.ReplaceAll(uniqId, ".", "")

	return id
}

func (ts *TagSet) resetTag(name string) (string, error) {
	tagKey := ts.tagKey(name)
	id := ts.generateId()

	_, err := ts.store.Forever(tagKey, []byte(id))
	if err != nil {
		return "", err
	}

	return id, nil
}

func (ts *TagSet) getNamespace() (string, error) {
	tagsIds, err := ts.tagIds()
	if err != nil {
		return "", err
	}

	namespace := strings.Join(tagsIds, "|")

	return namespace, nil
}

func (ts *TagSet) tagIds() ([]string, error) {
	tagsIDs := make([]string, len(ts.names))

	for i := 0; i != len(ts.names); i++ {
		tagID, err := ts.tagId(ts.names[i])
		if err != nil {
			return tagsIDs, err
		}

		tagsIDs[i] = tagID
	}

	return tagsIDs, nil
}

func (ts *TagSet) tagId(name string) (string, error) {
	tagName := ts.tagKey(name)

	idRaw, err := ts.store.Get(tagName)
	if err != nil && idRaw == nil {
		id, err := ts.resetTag(name)
		if err != nil {
			return "", err
		}

		return id, nil
	}

	return string(idRaw), nil
}

func (ts *TagSet) tagKey(name string) string {
	ident := "tag:" + name + ":key"

	return ident
}

func (ts *TagSet) getNames() []string {
	return ts.names
}

// NewTagSet instance of for tagged cache
func NewTagSet(store Store, names ...string) *TagSet {
	return &TagSet{
		store: store,
		names: names,
	}
}
