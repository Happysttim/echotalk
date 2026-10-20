# Echotalk

Echotalk은 궁금한 내용을 질문으로 올리고, 다른 사용자와 답변을 나누는 웹 서비스입니다. 데스크톱과 모바일에서 질문을 읽고 일반 또는 익명으로 답변할 수 있습니다.

## 주요 기능

- 이메일 인증 기반 회원가입·로그인, Google OAuth 로그인
- 질문 작성·수정·마감·삭제와 답변 마감일 설정
- 일반·익명 답변 작성, 답변 수정·삭제·추천
- 내가 작성한 질문과 질문별로 모아 보는 내 답변 목록

## 프로젝트 구성

- `apps/frontend`: React, TypeScript, Vite, TailwindCSS, Radix UI, Zustand
- `apps/backend`: Go, Gin, MongoDB, Redis
- `docker-compose.local.yml`: 로컬 MongoDB·Redis·백엔드 컨테이너 설정

pnpm workspace와 Turborepo로 프론트엔드·백엔드 명령을 함께 관리합니다.

## 로컬 실행

저장소 설정 기준으로 Node.js **22.19.0**, pnpm **9.4.0**, Go **1.26.1**을 준비합니다. 아래 방법은 MongoDB와 Redis만 Docker로 실행하고, 프론트엔드와 백엔드는 로컬에서 실행합니다. Docker Compose도 필요합니다.

### 1. 의존성 설치

저장소 루트에서 실행합니다.

```sh
pnpm install
go -C apps/backend mod download
```

### 2. 환경변수 설정

`apps/backend/.env.sample`을 참고하여 `apps/backend/.env`를 만듭니다. 기존 파일이 있으면 필요한 값만 수정합니다. 로컬 연결값의 예시는 다음과 같습니다.

```dotenv
WEB_URL=http://localhost:5173/
PORT=8080
MONGODB_URI=mongodb://localhost:27017
MONGODB=echotalk
REDIS_ADDR=localhost:6379
REDIS_PASSWORD=
REDIS_DB=0
```

샘플의 나머지 항목도 설정하세요. `SECRET_KEY`에는 직접 생성한 충분히 긴 비밀값을 사용하고, 메일 인증에는 실제 사용 가능한 `SMTP_HOST`, 숫자형 `SMTP_PORT`, `MAIL_ID`, `MAIL_PASSWORD`, `MAIL_FROM`이 필요합니다. 특히 `SMTP_PORT`는 서버 시작 시에도 숫자로 설정되어 있어야 합니다.

`apps/frontend/.env.example`을 참고하여 `apps/frontend/.env.development.local`을 만듭니다.

```dotenv
VITE_API_BASE_URL=/api
API_PROXY_TARGET=http://localhost:8080
```

Google 로그인을 사용할 경우 프론트엔드의 `VITE_GOOGLE_CLIENT_ID`와 백엔드의 `GOOGLE_CLIENT_ID`를 동일하게 설정합니다. Google Console의 승인된 Redirect URI, 백엔드 `GOOGLE_REDIRECT_URL`, 프론트엔드 `VITE_GOOGLE_REDIRECT_URI`는 모두 `http://localhost:5173/auth/google/callback`으로 맞춥니다. `GOOGLE_CLIENT_SECRET`은 백엔드에만 설정하며 `.env`의 비밀값을 Git에 올리지 않습니다.

### 3. 데이터 저장소와 앱 실행

```sh
docker compose -f docker-compose.local.yml up -d mongo redis
pnpm dev
```

- 프론트엔드: [http://localhost:5173](http://localhost:5173)
- 백엔드: [http://localhost:8080](http://localhost:8080)

`pnpm dev`는 프론트엔드와 백엔드를 함께 실행합니다. 위 명령은 Docker의 backend 서비스를 실행하지 않으므로 포트가 중복되지 않습니다. Windows에서 실행 정책 오류가 나면 `pnpm` 대신 `pnpm.cmd`를 사용하세요.

앱은 실행 터미널에서 `Ctrl+C`로 종료합니다. MongoDB와 Redis도 중지하려면 다음을 실행합니다. 저장된 볼륨 데이터는 유지됩니다.

```sh
docker compose -f docker-compose.local.yml stop mongo redis
```

## 사용방법

1. **질문 둘러보기**: 질문 목록에서 제목을 선택하여 본문과 답변을 확인합니다.
2. **가입·로그인**: 이메일 가입은 인증 메일의 링크를 열어 완료하거나, 설정된 Google 로그인으로 접속합니다.
3. **질문 작성**: 로그인 후 제목·내용·공개 여부·미래의 마감일(KST)을 설정합니다. 본인 질문은 상세 화면에서 관리합니다.
4. **답변 작성**: 일반 답변은 로그인 후 작성합니다. 익명 답변은 로그인 없이 작성할 수 있으며 삭제용 비밀번호를 설정합니다.
5. **내 활동 확인**: 계정 메뉴에서 내가 작성한 질문·답변을 확인합니다. 내 답변 목록의 질문 제목으로 상세 화면에 이동합니다.

## 빌드 및 검사

저장소 루트에서 실행합니다. 백엔드 테스트는 테스트 구성에 따라 데이터 저장소와 환경변수가 필요할 수 있습니다.

```sh
pnpm build
pnpm lint
pnpm typecheck
pnpm test
```

프론트엔드 빌드 결과는 `apps/frontend/dist`, 백엔드 빌드 결과는 `apps/backend/dist`에 생성됩니다. 운영 환경에서는 실제 API 주소, OAuth 설정, HTTPS·쿠키 정책과 프론트엔드 라우팅 fallback을 별도로 구성해야 합니다.

## 상세 문서

- [백엔드 실행·API 엔드포인트 안내](apps/backend/README.md)
- [프론트엔드 실행·환경변수 안내](apps/frontend/README.md)
- [프론트엔드 화면 명세](apps/frontend/docs/Echotalk-Frontend-Handoff.md)
- [반응형 디자인 시안](apps/frontend/examples/Echotalk-Responsive-Design.html)
