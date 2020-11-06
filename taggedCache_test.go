package cache

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vleukhin/goravel.cache/errs"
	cacheStore "github.com/vleukhin/goravel.cache/store"
	"go.uber.org/zap"
)

func TestTaggedCacheService_Get(t *testing.T) {
	type TestCase struct {
		TestName        string
		Key             string
		Tags            []string
		TaggedItemKey   string
		TaggedItemValue string
		Value           string
		Result          string
		HashedItemKey   string
		Error           error
	}
	testCases := make([]TestCase, 0)
	testCases = append(testCases,
		TestCase{
			TestName:        "Simply single key test",
			Key:             "goodwin.com/qui-accusamus-saepe-et-laborum-velit-id-vel_ad-units",
			Value:           "testValue",
			Result:          "testValue",
			Tags:            []string{"adUnits"},
			TaggedItemKey:   "tag:adUnits:key",
			TaggedItemValue: "5e4e50b381b08354922998",
			HashedItemKey:   "bff7bf5e588590a06e98f9109993699721b495f2:goodwin.com/qui-accusamus-saepe-et-laborum-velit-id-vel_ad-units",
			Error:           nil,
		},
		TestCase{
			TestName:        "Simply single MISSED key test",
			Value:           "testValue",
			Result:          "",
			Key:             "goodwin.com/qui-accusamus-saepe-et-laborum-velit-id-vel_ad-units",
			Tags:            []string{"adUnits"},
			TaggedItemKey:   "tag:adUnits:key",
			TaggedItemValue: "5e4e59648ef3e386705450",
			HashedItemKey:   "WRONG_KEY_497a5809edea4ee6e35e1c36add0bb:goodwin.com/qui-accusamus-saepe-et-laborum-velit-id-vel_ad-units",
			Error:           errs.ErrCacheMiss,
		},
	)

	lCfg := zap.NewProductionConfig()
	l, err := lCfg.Build()
	if err != nil {
		log.Fatal(err.Error())
	}

	for i, tCase := range testCases {
		t.Run(fmt.Sprintf("test %s", tCase.TestName), func(t *testing.T) {
			store, _ := cacheStore.NewInMemoryStore(context.Background(), "cn")
			cacheService, _ := NewTaggedCacheService(store, tCase.Tags...)

			val := []byte(tCase.Value)

			store.Forever(tCase.TaggedItemKey, []byte(tCase.TaggedItemValue))
			store.Forever(tCase.HashedItemKey, val)

			data, err := cacheService.Get(tCase.Key)
			res := string(data)

			expected := tCase.Result

			assert.Equal(t, err, tCase.Error,
				"they should be %s ,but %s \n Test name: %d", errs.ErrCacheMiss.Error(), err, i)

			assert.Equal(t, expected, res,
				"they should be %s ,but %s \n Test name: %d", tCase.Value, res, i)
		})
	}

}
