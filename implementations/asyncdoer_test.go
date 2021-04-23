package implementations

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/mikerott/gomock-bad/mocks"
)

func TestIt(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockThingConsumer := mocks.NewMockThingConsumer(ctrl)

	asyncDoer := Processor{
		ThingConsumer: mockThingConsumer,
		Jobs:          3,
	}

	mockThingConsumer.EXPECT().ConsumeThings().Return([]string{"thing1", "thing2", "thing3"}, nil)
	mockThingConsumer.EXPECT().ConsumeThings().Return(nil, fmt.Errorf("kablooie")).Times(2)

	asyncDoer.Process()

}
