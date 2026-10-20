# Echotalk Backend

Echotalk은 사용자가 질문을 올리고 서로 답변을 나누는 서비스입니다. 백엔드는 회원 인증, 질문·답변 관리, 추천, 이메일 인증과 내 활동 조회를 위한 HTTP API를 제공합니다.

Go와 Gin을 사용하며, 서비스 데이터는 MongoDB에, 이메일 인증 링크 정보는 Redis에 저장합니다. Google OAuth 로그인과 SMTP 인증 메일 발송을 지원합니다.

## 목차

- [로컬 실행](#로컬-실행)
- [환경변수](#환경변수)
- [API 공통 규칙](#api-공통-규칙)
- [API 엔드포인트](#api-엔드포인트)
- [응답 데이터와 페이지네이션](#응답-데이터와-페이지네이션)
- [사용 흐름과 요청 예시](#사용-흐름과-요청-예시)
- [프로젝트 구조 및 개발 명령](#프로젝트-구조-및-개발-명령)
- [현재 구현의 주의사항](#현재-구현의-주의사항)

## 로컬 실행

필요한 도구는 Go **1.26.1**(`go.mod` 기준), MongoDB, Redis입니다. 아래 예시는 Docker Compose로 MongoDB 7과 Redis 7을 실행하고 백엔드는 로컬에서 실행하는 방법입니다.

### 1. 환경 설정

[.env.sample](.env.sample)을 참고해 `apps/backend/.env`를 만듭니다. 기존 파일이 있으면 필요한 값만 수정하고, 실제 비밀값은 Git에 올리지 마세요.

```dotenv
WEB_URL=http://localhost:5173/
HOST=localhost
PORT=8080
MONGODB_URI=mongodb://localhost:27017
MONGODB=echotalk
REDIS_ADDR=localhost:6379
REDIS_PASSWORD=
REDIS_DB=0

SMTP_HOST=smtp.example.com
SMTP_PORT=587
MAIL_ID=your-mail-account
MAIL_PASSWORD=replace-with-mail-credential
MAIL_FROM=noreply@example.com

GOOGLE_CLIENT_ID=
GOOGLE_CLIENT_SECRET=
GOOGLE_REDIRECT_URL=http://localhost:5173/auth/google/callback
SECRET_KEY=replace-with-a-long-random-secret
```

SMTP 값은 형식 예시이므로 실제 메일 제공자의 설정으로 교체하세요. `SECRET_KEY`도 직접 생성한 충분히 긴 무작위 값으로 교체해야 합니다. Google 관련 비밀값은 Google 로그인을 사용할 때 설정합니다.

### 2. 데이터 저장소 실행

**저장소 루트**에서 실행합니다.

```sh
docker compose -f docker-compose.local.yml up -d mongo redis
```

이 명령은 MongoDB와 Redis만 실행합니다. Docker의 `backend` 서비스와 아래 로컬 서버를 동시에 실행하면 8080 포트가 충돌할 수 있습니다.

### 3. 백엔드 실행

저장소 루트에서 다음을 실행합니다.

```sh
cd apps/backend
go mod download
go run ./cmd/
```

기본 예시 주소는 [http://localhost:8080](http://localhost:8080)입니다. 브라우저에서 `/`를 열면 별도 홈 화면은 없으며, 아래 API 경로로 요청해야 합니다.

백엔드는 **현재 작업 디렉터리의 `.env`를 읽으므로 `apps/backend`에서 실행**하세요. 저장소 의존성을 설치했다면 이 디렉터리에서 `pnpm dev`로도 실행할 수 있습니다. 저장소 루트의 `pnpm dev`는 프론트엔드와 백엔드를 함께 실행합니다.

종료는 실행 터미널에서 `Ctrl+C`를 누릅니다. 데이터 저장소도 중지하려면 저장소 루트에서 다음을 실행합니다. 볼륨 데이터는 유지됩니다.

```sh
docker compose -f docker-compose.local.yml stop mongo redis
```

## 환경변수

| 변수 | 용도 / 설정 시 주의사항 |
| --- | --- |
| `WEB_URL` | 메일 인증 링크의 프론트엔드 주소. 현재 문자열 결합 방식이므로 **마지막 `/` 필수** |
| `HOST` | 설정에 로드되지만 현재 HTTP 서버의 바인딩 주소에는 사용하지 않음 |
| `PORT` | HTTP 포트. 예: `8080`. 서버는 `:PORT`로 수신 |
| `MONGODB_URI` | MongoDB 접속 URI |
| `MONGODB` | 데이터베이스 이름. `DB_NAME`이 아님 |
| `REDIS_ADDR` | Redis 주소. 예: `localhost:6379` |
| `REDIS_PASSWORD` | Redis 비밀번호. 로컬 무인증 설정에서는 빈 값 |
| `REDIS_DB` | Redis DB 번호. 예: `0` |
| `SMTP_HOST`, `SMTP_PORT` | 메일 제공자의 SMTP 호스트·포트 |
| `MAIL_ID`, `MAIL_PASSWORD`, `MAIL_FROM` | SMTP 인증 계정·비밀번호·발신 주소 |
| `GOOGLE_CLIENT_ID`, `GOOGLE_CLIENT_SECRET` | Google Web application OAuth 클라이언트 설정 |
| `GOOGLE_REDIRECT_URL` | Google Console과 일치하는 프론트엔드 callback URI |
| `SECRET_KEY` | JWT 서명용 비밀값 |

`.env` 파일이 없거나 `PORT`, `SMTP_PORT`, `REDIS_DB`를 정수로 해석할 수 없으면 초기화에 실패합니다. MongoDB는 초기화 시 실제 연결을 확인하므로 먼저 실행되어 있어야 합니다.

프론트엔드의 공개 `VITE_GOOGLE_CLIENT_ID`는 백엔드 Client ID와 같아야 합니다. Google Console의 승인된 Redirect URI도 `GOOGLE_REDIRECT_URL`과 일치시킵니다. Client Secret을 프론트엔드 환경변수에 넣지 마세요.

## API 공통 규칙

- 기본 주소: `http://localhost:8080`
- 백엔드 경로에는 `/api` 접두사가 없습니다. 프론트엔드 개발 프록시가 `/api/*`를 백엔드 `/*`로 전달합니다.
- JSON 요청 본문은 `Content-Type: application/json`을 사용합니다. DELETE 요청의 ID도 아래 표처럼 JSON Body로 전달합니다.
- `:surveyId`, `:answerId`와 본문의 ID는 MongoDB ObjectID의 24자리 16진수 문자열입니다.
- 날짜는 시간대가 포함된 RFC3339 문자열을 사용합니다. 예: `2099-12-31T14:59:00Z`.
- 아래 응답 표의 `data`는 공통 응답 객체 안의 `data`를 의미합니다.

### 인증 표기

| 표기 | 요구사항 |
| --- | --- |
| 없음 | 인증 미들웨어 없음 |
| 비로그인 | `AuthUnrequired`. Authorization 헤더를 보내지 않음 |
| 필수 | `AuthRequired`. `Authorization: Bearer <accessToken>` 필요 |
| 선택 | `AuthNoRequired`. 헤더가 없으면 진행하며, 보내면 유효 토큰이어야 함 |
| 인증 쿠키 | `VerifyRequired`. `VERIFY_TOKEN` 쿠키 필요 |

로그인·가입 성공 시 `data.accessToken`과 HttpOnly `REFRESH_TOKEN` 쿠키를 받습니다. `/auth/refresh`는 쿠키로 토큰을 갱신하며, `/code` 성공 시에는 HttpOnly `VERIFY_TOKEN` 쿠키를 받습니다. 브라우저에서는 쿠키가 필요한 API에 `credentials: 'include'`를 사용하세요.

### 공통 응답

```json
{
  "status": "ok",
  "data": {}
}
```

데이터가 없는 성공 응답은 `{"status":"ok"}`입니다. 실패는 HTTP 오류 상태와 함께 아래 형식으로 반환합니다. `reason`은 선택 항목입니다.

```json
{
  "status": "failed",
  "error": "오류 메시지",
  "reason": "추가 오류 정보"
}
```

HTTP 상태와 `status`를 함께 확인하세요. 생성 API는 성공 시 201, 비밀번호 변경은 202, 그 외 아래 성공 응답은 200을 사용합니다. 실제 실패 코드·메시지는 핸들러별로 다릅니다.

## API 엔드포인트

기준: [router.go](internal/router/router.go)와 현재 핸들러 구현, 2026-10-20. `Query`는 URL 쿼리, `Body`는 JSON 본문입니다.

### 인증·계정

| Method | Endpoint | 인증 | 입력 | 성공 응답 / 동작 |
| --- | --- | --- | --- | --- |
| GET | `/auth/google` | 비로그인 | Query: `code` | Google 인가 코드 교환, `data.accessToken` + refresh 쿠키 |
| POST | `/auth/local` | 비로그인 | Body: `email`, `password` | 이메일 로그인, `data.accessToken` + refresh 쿠키 |
| POST | `/auth/register` | 비로그인 + 인증 쿠키 | Body: `email`, `password`, `nickname` | 가입, `data.accessToken` + refresh 쿠키, 인증 쿠키 제거 |
| POST | `/auth/change` | 비로그인 + 인증 쿠키 | Body: `email`, `password` | 이메일 계정 비밀번호 변경, 인증 쿠키 제거 |
| GET | `/auth/refresh` | refresh 쿠키 | `REFRESH_TOKEN` 쿠키 | `data.accessToken` + 새 refresh 쿠키 |
| GET | `/auth/logout` | 필수 + refresh 쿠키 | `REFRESH_TOKEN` 쿠키 | 로그아웃 및 refresh 쿠키 제거 |
| DELETE | `/auth` | 필수 | Body 없음 | 로그인한 계정 삭제 및 refresh 쿠키 제거 |

가입·비밀번호 변경의 `email`은 인증 쿠키의 이메일과 같아야 합니다. Google 신규 사용자의 닉네임은 Google `Claims.Name`을 사용하며 별도 닉네임 설정 API는 없습니다.

### 질문

| Method | Endpoint | 인증 | 입력 | 성공 응답 / 동작 |
| --- | --- | --- | --- | --- |
| GET | `/surveys` | 없음 | Query: `type`, `cursor`, `limit` | `data: {skip, hasNext, surveys: SurveyDetail[]}` |
| GET | `/surveys/me` | 필수 | Query: `type`, `cursor`, `limit` | 내가 작성한 공개 질문, 동일한 목록 응답 |
| GET | `/surveys/:surveyId` | 없음 | Path: `surveyId` | `data: SurveyDetail` |
| POST | `/surveys` | 필수 | Body: `title`, `content`, `is_public`, `expires_at` | `data: Survey` |
| PATCH | `/surveys` | 필수 | Body: `survey_id`, `title`, `content`, `is_public`, `closed`, `expires_at` | 본인 질문 수정·마감 |
| DELETE | `/surveys` | 필수 | Body: `survey_id` | 본인 질문 및 연결된 답변 삭제 |

`type`은 `_id` 또는 `updated_at`, 기본값은 `_id`입니다. `limit`은 10 또는 30이며 누락·허용되지 않은 값은 30으로 처리합니다. 첫 요청은 `cursor`를 생략하거나 빈 값으로 보냅니다.

질문 수정은 일부 필드만 보내는 방식으로 가정하지 마세요. 현재 요청 모델에 필요한 전체 필드를 보내며, 마감은 같은 PATCH 요청에서 `closed: true`로 설정합니다.

### 답변

| Method | Endpoint | 인증 | 입력 | 성공 응답 / 동작 |
| --- | --- | --- | --- | --- |
| GET | `/answers` | 없음 | Query: `cursor` **필수** | `data: {skip, hasNext, answers: AnswerDetail[]}` |
| GET | `/answers/me` | 필수 | Query: `cursor` | `data: {skip, hasNext, answers: MyAnswerDetail[]}` |
| GET | `/answers/:answerId` | 없음 | Path: `answerId` | `data: AnswerDetail` |
| POST | `/answers` | 선택 | Body: `survey_id`, `content`, `is_anonymous`, `anonymous`, `answer_password` | `data: Answer` |
| PATCH | `/answers` | 필수 | Body: `survey_id`, `answer_id`, `content` | 본인 답변 수정 |
| DELETE | `/answers` | 선택 | Body: `answer_id`, `answer_password` | 일반 답변은 본인 인증, 익명 답변은 비밀번호로 삭제 |
| POST | `/answers/rateup` | 필수 | Body: `answer_id` | 답변 추천. 추천 취소 엔드포인트는 없음 |

일반 답변 작성에는 로그인 토큰이 필요합니다. 익명 답변은 `is_anonymous: true`와 삭제용 `answer_password`가 필요하며, `anonymous`가 빈 값이면 서버 기본 이름을 사용합니다. 일반 답변의 `anonymous`, `answer_password`는 빈 문자열로 보낼 수 있습니다. 마감되었거나 만료된 질문에는 새 답변을 등록할 수 없습니다.

### 이메일 인증

| Method | Endpoint | 인증 | 입력 | 성공 응답 / 동작 |
| --- | --- | --- | --- | --- |
| GET | `/verify` | 비로그인 | Query: `email`, `verifyType` | 인증 메일 발송 |
| GET | `/code` | 비로그인 | Query: `code`, `verifyType` | 링크 토큰 확인, `VERIFY_TOKEN` 쿠키 설정 |
| GET | `/verify/me` | 인증 쿠키 | Body 없음 | `data`: 인증된 이메일 문자열 |

`verifyType`은 `register`(가입) 또는 `password`(비밀번호 변경)입니다. `/code` 성공은 JSON의 `true`가 아니라 `{"status":"ok"}`와 인증 쿠키로 확인합니다.

## 응답 데이터와 페이지네이션

### 주요 데이터 형식

| 타입 | 주요 필드 |
| --- | --- |
| `UserSummary` | `id`, `nickname`, `hash_number` |
| `Survey` | `id`, `title`, `content`, `is_public`, `closed`, `author_id`, `expires_at`, `created_at`, `updated_at` |
| `SurveyDetail` | `Survey` 필드 + `author: UserSummary` |
| `Answer` | `id`, `survey_id`, `author_id`, `is_anonymous`, `anonymous`, `content`, `rate_up`, `created_at`, `updated_at` |
| `AnswerDetail` | `Answer` 필드 + `author: UserSummary` |
| `MyAnswerDetail` | `survey_answers` 안에 질문 메타데이터와 `answers: AnswerDetail[]` |

생성 응답의 `Survey`·`Answer`와 조회 응답의 `SurveyDetail`·`AnswerDetail`을 구분하세요. 답변 비밀번호는 응답에 포함하지 않습니다. `hash_number`는 int64 숫자이므로 JavaScript에서 안전 정수 범위를 넘어가는 값은 별도 정밀도 처리가 필요합니다.

`/answers/me`의 구조는 다음과 같습니다. 읽기 편하도록 일부 필드만 표시했습니다.

```json
{
  "status": "ok",
  "data": {
    "skip": "",
    "hasNext": false,
    "answers": [
      {
        "survey_answers": {
          "id": "507f1f77bcf86cd799439011",
          "title": "어떤 책을 좋아하시나요?",
          "answers": [
            {
              "id": "507f1f77bcf86cd799439012",
              "survey_id": "507f1f77bcf86cd799439011",
              "content": "여행 에세이를 좋아합니다.",
              "author": {
                "id": "507f1f77bcf86cd799439013",
                "nickname": "서윤",
                "hash_number": 10000001
              }
            }
          ]
        }
      }
    ]
  }
}
```

`survey_answers`에는 질문의 `is_public`, `closed`, `author_id`, `expires_at`, `created_at`, `updated_at`도 포함됩니다. 질문의 `content`와 질문 작성자의 `author` 요약은 포함하지 않습니다. 내 답변의 질문 제목은 이 응답을 그대로 사용할 수 있어 추가 질문 조회가 필요하지 않습니다.

### 커서 사용법

| API | 첫 요청 | 다음 요청 |
| --- | --- | --- |
| `/surveys`, `/surveys/me` | cursor 생략 또는 빈 값 | 응답 `skip` 그대로 전달 |
| `/answers` | `{"survey_id":"질문 ID","answer_id":""}`를 padding 없는 base64url로 인코딩 | 응답 `skip` 그대로 전달 |
| `/answers/me` | cursor 생략 또는 빈 값 | 응답 `skip` 그대로 전달. **원문 답변 ID**이며 base64url이 아님 |

질문 cursor의 내부 구조는 `{survey_id, updated_at}`입니다. `updated_at`은 RFC3339Nano 정밀도를 유지해야 하므로 클라이언트에서 날짜를 재직렬화하거나 커서를 재조립하지 않는 것이 좋습니다.

`hasNext`로 다음 페이지 존재 여부를 판단하고, 다음 요청에 `skip`을 URL 인코딩해 전달합니다. 목록이 `null`이면 빈 배열로 처리합니다. 내 답변은 페이지 간 같은 질문이 반복될 수 있으므로 질문 ID로 그룹을 병합하고 답변 ID 중복을 제거합니다.

## 사용 흐름과 요청 예시

### 이메일 가입·비밀번호 변경

1. `GET /verify?email=...&verifyType=register` 또는 `verifyType=password`로 메일을 요청합니다.
2. 메일의 프론트엔드 링크를 엽니다. 가입은 `/register?verify={token}`, 비밀번호 변경은 `/change?verify={token}`입니다.
3. 프론트엔드가 `GET /code?code={token}&verifyType=...`를 호출합니다.
4. 성공 후 `VERIFY_TOKEN` 쿠키를 유지한 채 `/verify/me`로 인증 이메일을 확인합니다.
5. 가입은 `POST /auth/register`, 비밀번호 변경은 `POST /auth/change`로 완료합니다.

링크 토큰은 일회용입니다. 쿠키를 지원하는 클라이언트를 사용하고, 토큰·비밀번호를 로그나 공유 문서에 남기지 마세요.

### Google 로그인

프론트엔드에서 Google 인증을 시작하고 등록된 callback으로 인가 코드를 받습니다. 이후 `GET /auth/google?code=...`로 교환합니다. 이 API는 Google 로그인 화면으로 리다이렉트하는 진입 URL이 아니라 **인가 코드 교환 API**입니다. Google client secret은 서버에서만 사용합니다.

### 공개 질문 조회

```sh
curl "http://localhost:8080/surveys?type=_id&limit=10"
```

### 로그인

```http
POST /auth/local
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "<YOUR_PASSWORD>"
}
```

### 질문 작성

```http
POST /surveys
Authorization: Bearer <ACCESS_TOKEN>
Content-Type: application/json

{
  "title": "어떤 책을 좋아하시나요?",
  "content": "좋아하는 책과 그 이유를 알려주세요.",
  "is_public": true,
  "expires_at": "2099-12-31T14:59:00Z"
}
```

예시 이메일·토큰·날짜는 실제 요청에 맞게 교체하세요. 생성한 질문의 ID로 답변을 등록할 수 있습니다. Windows PowerShell에서 `curl`이 다른 명령의 별칭이면 `curl.exe`를 사용하세요.

## 프로젝트 구조 및 개발 명령

```text
cmd/                   서버 진입점과 의존성 연결
internal/
  auth/                JWT, 인증 미들웨어, Google OAuth
  config/              환경변수 로드
  database/            MongoDB 연결과 인덱스
  handler/             요청 검증과 HTTP 응답
  model/               저장 모델과 요청·응답 DTO
  repositories/        데이터 접근
  service/             서비스 로직
  redis/               인증 링크 정보
  mail/                SMTP 발송
  router/              API 경로 등록
  response/            공통 JSON 응답
  server/              HTTP 서버와 종료 처리
  timer/               주기적 작업
```

`apps/backend`에서 실행합니다.

```sh
go run ./cmd/
go build -o dist/server ./cmd
go test ./...
```

pnpm을 사용한다면 각각 `pnpm dev`, `pnpm build`, `pnpm test`에 대응합니다. Windows용 실행 파일을 명시적으로 만들려면 `go build -o dist/server.exe ./cmd`를 사용하세요. 실행 파일도 `.env`를 읽을 수 있는 작업 디렉터리에서 시작해야 합니다.

테스트가 설정·데이터베이스 패키지를 로드하면 `.env`와 MongoDB 등 외부 의존성이 필요할 수 있습니다. Go 테스트는 패키지 디렉터리에서 실행되므로 테스트별 설정 로딩 방식도 확인하세요. 이 문서는 코드 확인을 기반으로 작성했으며 실제 SMTP·Google OAuth 연동 성공을 보장하는 테스트 결과는 아닙니다.

## 현재 구현의 주의사항

- `/surveys/me`도 `is_public: true` 필터를 적용하므로 **내 공개 질문만** 반환합니다.
- 질문·답변 상세 GET에는 인증 미들웨어가 없습니다. `is_public: false`만으로 상세 접근이 차단된다고 가정하지 마세요.
- 질문 생성의 `is_public`은 현재 `bool`에 `binding:"required"`가 적용되어 `false`가 검증에서 거부될 수 있습니다. 비공개 생성 사용 전 서버 검증 보완이 필요합니다.
- `/answers/me`는 정렬과 cursor 필터·그룹 단위 제한의 정합성을 추가 확인해야 합니다. 클라이언트에서는 서버의 `skip`을 그대로 사용하며 누락 문제를 임의 커서 보정으로 해결하지 않습니다.
- 익명 답변에 `author`가 함께 반환되더라도 화면에는 계정 닉네임·해시번호를 표시하지 않습니다. 응답 자체에서의 개인정보 보호는 서버 정책도 별도로 검증해야 합니다.
- 현재 로그인·refresh 응답은 `accessToken`만 제공하며, 별도 로그인 사용자 프로필 조회 API는 없습니다.
- 현재 쿠키는 개발용으로 `Secure=false`입니다. 운영 배포 전 HTTPS, 쿠키 속성, CORS와 허용 origin, 비밀값 관리 정책을 검토하세요.
- Docker Compose의 `DB_NAME`과 애플리케이션이 읽는 `MONGODB`는 다릅니다. 컨테이너로 백엔드를 실행할 때도 `.env`에 `MONGODB`를 지정해야 합니다.

## 관련 문서

- [프로젝트 README](../../README.md)
- [프론트엔드 README](../frontend/README.md)
- [API 라우터](internal/router/router.go)
- [요청·응답 모델](internal/model)
