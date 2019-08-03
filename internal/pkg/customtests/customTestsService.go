package customtests

import (
	"fmt"
	"vegeta-kubernetes/internal/pkg/utils"

	"vegeta-kubernetes/internal/pkg/vegeta"

	"github.com/pkg/errors"
)

var testMap = map[string]func(utils.AttackConf) error{
	"testExample": testExample,
}

func Run(testName string, ac utils.AttackConf) error {
	if val, ok := testMap[testName]; ok {
		if err := val(ac); err != nil {
			return errors.Wrapf(err, "Failed to run test name '%s'", testName)
		}
		return nil
	}

	return errors.Errorf("Test name '%s' doesn't exist", testName)
}

func testExample(ac utils.AttackConf) error {
	fmt.Println("testExample running")
	vegeta.Attack(ac)
	return nil
}
