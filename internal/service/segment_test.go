package service

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"github.com/PoorMercymain/user-segmenter/errors"
	"github.com/PoorMercymain/user-segmenter/internal/domain"
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

func TestUpdateUserSegments(t *testing.T) {
	logger.InitLogger()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockSegmentRepository(ctrl)

	seg := NewSegment(mockRepo)

	mockRepo.EXPECT().UpdateUserSegments(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.ErrorLoggerNotInitialized).MaxTimes(1)
	mockRepo.EXPECT().UpdateUserSegments(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

	err := seg.UpdateUserSegments(context.Background(), "1", []string{"a"}, []string{"b"})
	require.Error(t, err)

	err = seg.UpdateUserSegments(context.Background(), "1", []string{"a"}, []string{"b"})
	require.NoError(t, err)
}

func TestReadUserSegments(t *testing.T) {
	logger.InitLogger()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockSegmentRepository(ctrl)

	seg := NewSegment(mockRepo)

	mockRepo.EXPECT().ReadUserSegments(gomock.Any(), gomock.Any()).Return(nil, errors.ErrorNoRows).MaxTimes(1)
	mockRepo.EXPECT().ReadUserSegments(gomock.Any(), gomock.Any()).Return([]string{"a", "b"}, nil).AnyTimes()

	userSegments, err := seg.ReadUserSegments(context.Background(), "1")
	require.Error(t, err)
	require.Empty(t, userSegments)

	userSegments, err = seg.ReadUserSegments(context.Background(), "1")
	require.NoError(t, err)
	require.Len(t, userSegments, 2)
}

func TestReadUserSegmentsHistory(t *testing.T) {
	logger.InitLogger()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockSegmentRepository(ctrl)

	seg := NewSegment(mockRepo)

	mockRepo.EXPECT().ReadUserSegmentsHistory(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.ErrorNoRows).MaxTimes(1)
	mockRepo.EXPECT().ReadUserSegmentsHistory(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return([]domain.HistoryElem{{UserID: "1", Slug: "a", Operation: "addition", DateTime: time.Now()}}, nil).AnyTimes()

	history, err := seg.ReadUserSegmentsHistory(context.Background(), "1", time.Now(), time.Now())
	require.Error(t, err)
	require.Empty(t, history)

	history, err = seg.ReadUserSegmentsHistory(context.Background(), "1", time.Now(), time.Now())
	require.NoError(t, err)
	require.Len(t, history, 1)
}

func TestCreateDeletionTime(t *testing.T) {
	logger.InitLogger()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockSegmentRepository(ctrl)

	seg := NewSegment(mockRepo)

	mockRepo.EXPECT().CreateDeletionTime(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.ErrorLoggerNotInitialized).MaxTimes(1)
	mockRepo.EXPECT().CreateDeletionTime(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

	err := seg.CreateDeletionTime(context.Background(), "1", "a", time.Now())
	require.Error(t, err)

	err = seg.CreateDeletionTime(context.Background(), "1", "a", time.Now())
	require.NoError(t, err)
}

func TestAddSegmentToPercentOfUsers(t *testing.T) {
	logger.InitLogger()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockSegmentRepository(ctrl)

	seg := NewSegment(mockRepo)

	mockRepo.EXPECT().AddSegmentToPercentOfUsers(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.ErrorLoggerNotInitialized).MaxTimes(1)
	mockRepo.EXPECT().AddSegmentToPercentOfUsers(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

	err := seg.AddSegmentToPercentOfUsers(context.Background(), "a", 10)
	require.Error(t, err)

	err = seg.AddSegmentToPercentOfUsers(context.Background(), "a", 10)
	require.NoError(t, err)
}
