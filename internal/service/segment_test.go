package service

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"github.com/PoorMercymain/user-segmenter/internal/domain/mocks"
	"github.com/PoorMercymain/user-segmenter/pkg/logger"
)

func TestNewSegment(t *testing.T) {
	seg := NewSegment(nil)
	require.Empty(t, seg)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockSegmentRepository(ctrl)

	seg = NewSegment(mockRepo)
	require.NotEmpty(t, seg)
}

func TestCreateSegment(t *testing.T) {
	logger.InitLogger()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockSegmentRepository(ctrl)

	seg := NewSegment(mockRepo)

	mockRepo.EXPECT().CreateSegment(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

	err := seg.CreateSegment(context.Background(), "~not~a~slug~")
	require.Error(t, err)

	err = seg.CreateSegment(context.Background(), "a-slug")
	require.NoError(t, err)
}

func TestDeleteSegment(t *testing.T) {
	logger.InitLogger()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockSegmentRepository(ctrl)

	seg := NewSegment(mockRepo)

	mockRepo.EXPECT().DeleteSegment(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

	err := seg.DeleteSegment(context.Background(), "a-slug")
	require.NoError(t, err)
}
