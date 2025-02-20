# gocore

## 개요

Go 백엔드 서비스를 만들때의 기본이 되는 부분을 하나씩 만들어본다.
완성된 전체를 보여주는 것이 아닌 만들어가는 과정을 따라가 보는 것을 목표로 한다. 
- 각 단계는 별도의 브랜치를 딴다.
- 다음 단계에서는 기능들을 하니씩 추가한다. 

## 구현 프로세스
### 1. 클린 아키텍처 기본
기초적인 클린 아키텍처 구조의 백엔드 서버를 만들어본다.  
- 브랜치: 1_clean-architecture-basic
- 블로그 포스팅: https://jusths.tistory.com/442

### 2. 설정과 데이터베이스 
환경설정을 위한 설정 파일, 그리고 PostgreSQL을 사용하여 개선한다.
- 브랜치: 2_config-and-db
- 블로그 포스팅: 
  - https://jusths.tistory.com/443
  - https://jusths.tistory.com/444
- docker compose 사용법
  - 로컬 PC에서 PostgreSQL을 사용하기 위해 docker compose를 사용한다.
  - Docker Desktop이 아닌 OrbStack을 사용하였다([관련 블로그 링크](https://velog.io/@nchime/OrbStack)).
  ```shell
  $ docker-compose up -d
  $ docker-compose down
  ```
  
### 3. 유닛 테스트
Mockery를 사용하여 Mock을 만들어 테스트에 사용한다.
- 브랜치: 3_unit-test
- 블로그 포스팅:- https://jusths.tistory.com/445
- Mockery로 Mock 생성하기
  ```shell
  $ mockery
  ```
- Makefile을 사용하여 테스트 실행
  ```shell
  $ make test # mockery 부터 실행하고 test를 실행한다.
  ```
- pre-commit을 사용하여 커밋 전에 테스트 `make test`를 실행하도록 설정
  - test를 실행하고 성공하면 커밋을 진행한다.

### 4. 의존성 주입
[fx](https://github.com/uber-go/fx) 라이브러리를 사용하여 의존성 주입을 구현한다.
- 브랜치: 4_fx-begin, 4_fx-end
- 블로그 포스팅:- https://jusths.tistory.com/446

### 5. 미들웨어
다양한 미들웨어를 적용해본다.
- 브랜치: 5_middleware
- 블로그 포스팅:- https://jusths.tistory.com/447