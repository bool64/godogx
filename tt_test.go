package godogx_test

import (
	"testing"

	"github.com/dailymotion/allure-go"
)

func TestParameterized(t *testing.T) {
	for i := 0; i < 5; i++ {
		t.Run("", func(t *testing.T) {
			allure.Test(t,
				allure.Description("Test with parameters"),
				allure.Parameter("counter", i),
				allure.Action(func() {}))
		})
	}
}
