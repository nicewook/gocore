package usecase

import (
	"context"

	"github.com/nicewook/gocore/internal/domain"
)

type userUseCase struct {
	userRepo domain.UserRepository
}

func NewUserUseCase(userRepo domain.UserRepository) domain.UserUseCase {
	return &userUseCase{userRepo: userRepo}
}

func (uc *userUseCase) GetByID(ctx context.Context, id int64) (*domain.User, error) {
	return uc.userRepo.GetByID(ctx, id)
}

func (uc *userUseCase) GetAll(ctx context.Context, req *domain.GetAllUsersRequest) (*domain.GetAllResponse, error) {
	// usecase 레이어에서 추가적인 비즈니스 로직이 필요할 경우 여기에 구현
	// 예: 특정 권한에 따른 필터링, 데이터 검증, 복잡한 비즈니스 룰 적용 등

	// 리포지토리 호출 결과를 확인하고 처리
	response, err := uc.userRepo.GetAll(ctx, req)
	if err != nil {
		// 에러 발생 시 nil과 에러를 반환
		return nil, err
	}

	// 성공 시 응답 반환
	return response, nil
}
