# Echotalk frontend

React / TypeScript / Vite / TailwindCSS / Radix UI / Zustand 프론트엔드입니다. 화면 기준은 `AGENTS.md`, `docs/Echotalk-Frontend-Handoff.md`, `examples/Echotalk-Responsive-Design.html`입니다.

## 실행과 검증

프론트엔드 디렉터리에서 실행합니다. Windows에서는 `pnpm.cmd`를 사용해도 됩니다.

```sh
pnpm dev
pnpm lint
pnpm typecheck
pnpm test
pnpm build:development
pnpm build
pnpm preview
```

`dev`는 development 모드, `build`는 production 모드입니다. `preview`는 빌드 결과 확인용이며 운영 서버가 아닙니다. OAuth Redirect URI 일치를 위해 개발 서버는 `http://localhost:5173`을 사용하며 포트가 사용 중이면 임의 포트로 변경하지 않고 종료합니다.

## 환경변수

| 파일 / 변수                                        | 역할                                                                       |
| -------------------------------------------------- | -------------------------------------------------------------------------- |
| `.env`                                             | 공통 기본값                                                                |
| `.env.development`                                 | 개발 모드 설정                                                             |
| `.env.production`                                  | 운영 빌드 설정                                                             |
| `.env.example`                                     | 설정 예시                                                                  |
| `.env.development.local` / `.env.production.local` | Git에서 제외되는 환경별 개인 설정                                          |
| `VITE_API_BASE_URL`                                | 브라우저의 API 기본 주소. 기본 `/api`                                      |
| `API_PROXY_TARGET`                                 | 개발 서버가 `/api` 요청을 전달할 백엔드 주소. 기본 `http://localhost:8080` |
| `VITE_GOOGLE_CLIENT_ID`                            | Google OAuth의 공개 Web application Client ID                              |
| `VITE_GOOGLE_REDIRECT_URI`                         | Google에서 돌아올 프론트엔드 callback. 기본 `/auth/google/callback`        |

개발 백엔드 포트를 바꾸려면 `.env.development.local`에 다음을 지정합니다.

```dotenv
VITE_API_BASE_URL=/api
API_PROXY_TARGET=http://localhost:8080
```

개발 프록시는 `/api/surveys`를 백엔드 `/surveys`로 전달합니다. 현재 백엔드는 `/api` 접두사 없는 라우트를 사용합니다.

운영 기본값도 `/api`입니다. 운영 웹 서버에서 `/api/*`를 백엔드 `/*`로 전달하는 reverse proxy와, 프론트엔드 경로의 `index.html` fallback을 설정해야 합니다. `API_PROXY_TARGET`은 운영 번들에 포함되지 않으며 운영 reverse proxy를 설정하지 않습니다.

별도 API 도메인을 사용할 경우 `.env.production.local` 또는 CI 빌드 환경에 실제 주소를 지정합니다.

```dotenv
VITE_API_BASE_URL=https://api.example.com
```

위 도메인은 예시이며 실제 운영 주소로 바꿔야 합니다. 별도 origin을 쓰면 백엔드에서 정확한 프론트엔드 origin에 대한 CORS와 credentials를 허용하고, 쿠키의 Secure/SameSite/Domain 정책을 배포 환경에 맞춰야 합니다. 모든 API 요청은 `credentials: include`를 사용합니다.

Vite 환경변수는 **빌드 시점**에 반영됩니다. 변경 후 개발 서버를 재시작하거나 운영 번들을 다시 빌드하세요. 모드별 파일은 공통 파일보다 우선하고, 빌드 프로세스에 지정된 환경변수가 우선합니다. `VITE_*` 값은 사용자에게 공개되므로 비밀키, DB 비밀번호, OAuth client secret을 넣지 않습니다.

## 모듈 구조

- `src/pages`: 인증, 질문 목록/상세/편집, 답변 작성/수정 화면
- `src/components`: 공통 버튼/입력/확인창, 레이아웃, 질문/답변 행
- `src/api`: 공통 요청·오류·세션 갱신과 기능별 API
- `src/stores`: Zustand 메모리 세션과 경로/이탈 확인 상태
- `src/hooks`: 요청 취소·커서 피드·리소스 조회·작성 내용 변경 감지
- `src/routing`: 경로 해석, 로그인 복귀, 작성 중 이탈 방지
- `src/config`, `src/lib`: 환경변수 검증, KST 날짜 변환, 안전한 내부 이동
- `src/styles`: 시안 기반 색상/간격, 화면별 스타일, 반응형/접근성 보완
- `src/test`: API 계약과 실패/재시도/포커스/인증 단계 회귀 테스트

Radix UI의 Select, Switch, DropdownMenu, AlertDialog와 Lucide 아이콘을 사용합니다. Access token은 Zustand 메모리에만 유지하고 localStorage/sessionStorage에 저장하지 않습니다. JWT의 사용자 ID는 UI 노출에만 사용하며 실제 권한 검증은 서버 책임입니다. 만료된 세션의 읽기 요청만 한 번 재시도하고, 쓰기 요청을 자동 재전송하지 않습니다.

## 구현 상태와 연동 확인 사항

메일 인증은 코드 입력이 아닌 링크 확인 방식입니다.

- 회원가입: `/register?verify={token}` → `/code?code={token}&verifyType=register` → 가입 정보 입력
- 비밀번호 변경: `/change?verify={token}` → `/code?code={token}&verifyType=password` → 새 비밀번호 입력
- 기존 `/password/reset` 경로도 지원합니다. 메일 발송 후에는 메일 링크를 열어야 다음 단계로 진행합니다.
- API 성공 후에만 최종 폼을 열며 VERIFY_TOKEN은 HttpOnly 쿠키로 브라우저가 처리합니다. 링크 토큰은 처리 후 URL에서 제거하고 저장소에 보관하지 않습니다.
- 현재 `/code` 응답은 이메일을 제공하지 않습니다. 사용자 수정본의 `authApi.verifyMe()`는 `/verify/me`에서 인증 이메일을 조회해 입력 불가능한 필드에 표시합니다. 서버는 쿠키의 인증 이메일과 Payload 이메일을 비교합니다.
- 처리 결과가 불확실하거나 링크가 만료되면 새 메일을 요청합니다. 새로고침 후 클라이언트 상태만으로 인증 성공을 복원하지 않습니다.

메인, 이메일 로그인/가입/비밀번호 재설정, 질문 목록/상세/작성/수정/마감/삭제, 일반·익명 답변, 일반 답변 수정, 답변 삭제·추천의 화면과 API 호출을 구현했습니다. 입력 보존, 중복 제출 방지, 로딩/빈/오류/재시도, 모바일 고정 답변 버튼, 키보드 포커스 복원과 이탈 확인을 포함합니다.

현재 서버 코드의 응답 wrapper와 커서를 반영했습니다. 회원가입 응답에는 accessToken이 있어 가입 완료 시 로그인 상태를 유지합니다. 이 부분은 응답 계약 미제공을 전제로 한 handoff의 로그인 이동안에서 조정한 사항입니다.

다음은 완료로 간주하지 않습니다.

- Google OAuth 실제 계정 통합: 프론트엔드 구현과 Google 로그인 화면 진입은 확인했습니다. 실제 계정의 동의 및 백엔드 세션 생성까지 검증한 결과는 아닙니다.
- 작성자 닉네임/해시: 현재 질문·답변 DTO에 없어 일반 작성자는 기본 작성자 표시를 사용합니다. 임의 해시를 생성하지 않습니다.
- 실제 백엔드 통합: 이메일 발송, 인증 쿠키 회전, 비공개 정책, 서버 측 소유권/마감 차단, 추천 중복 정책을 실제 서비스에서 검증해야 합니다. 프론트엔드는 미제공 정책을 임의로 확정하지 않습니다.
- 운영 주소/배포: 실제 운영 API 주소와 reverse proxy 설정은 제공된 값으로 확정해야 합니다.

## 화면 검증용 fixture

`node scripts/preview-api.mjs`는 `127.0.0.1:3100`에서 읽기 전용 샘플 데이터를 제공합니다. `--signed-out`을 붙이면 비로그인 상태를 확인할 수 있습니다. 모든 데이터 변경 요청은 503으로 거절합니다.

별도 터미널에서 개발 서버의 `API_PROXY_TARGET`만 `http://127.0.0.1:3100`으로 지정해 사용합니다. 이 서버는 개발 환경 기본값이나 운영 번들에 연결되지 않습니다. 화면 검증 결과는 실제 백엔드 연동 성공을 의미하지 않습니다.

### 2026-09-11 검증 기록

- `lint`, `typecheck`, 테스트 14개, development 빌드, production 빌드 통과.
- 360 / 390 / 768 / 1024 / 1440px에서 주요 11개 경로의 가로 넘침 검사 55건과 비로그인 로그인 화면 5건 확인. 검사 대상 입력 글꼴은 16px.
- 데스크톱 로그인·가입·홈·상세, 태블릿 가입, 모바일 로그인·익명 작성·상세·이탈 확인창을 스크린샷으로 확인.
- 작성 중 이탈 취소 시 입력 보존 및 포커스 복원, 정렬 변경 후 선택값·포커스 반영 확인.
- 위 브라우저 검사는 read-only fixture 기반입니다. 실제 이메일/OAuth/데이터 변경의 end-to-end 검증과 스크린리더 검증은 미실시입니다.

메일 링크 변경 후 추가 회귀 테스트: register/change 직접 진입, 기존 reset 경로, StrictMode 중복 요청 방지, 빈/만료 토큰, 이전 요청의 늦은 응답 무시, 쿠키 포함·Bearer 제외, 명시적인 false 응답 거부를 검사합니다.

OAuth 추가 후 테스트 28개 통과. OAuth 검증은 state 위조/만료/재사용/취소, 코드 교환 시 쿠키와 헤더, StrictMode 중복 교환 방지, 초기 세션 갱신과 새 로그인 간 순서, 안전한 복귀 경로를 포함합니다. 사용자 수정 중인 이메일 조회 API는 테스트에서 격리했으며 해당 구현을 덮어쓰지 않았습니다.

## 2026-09-12 내 활동 추가

- `/my/questions`: 인증된 내 질문 목록, 기존 질문 행·정렬·더 보기 재사용.
- `/my/answers`: 내 답변을 질문별로 묶고 질문 제목과 개별 답변에서 상세로 이동. 추가 페이지 병합, 중복 제거, 제목 조회 실패 시 답변 보존 및 재시도.
- 계정 메뉴와 모바일 내 활동 탭으로 접근. 계정 변경 시 이전 목록 제거.
- 질문 작성·수정 마감일은 KST의 다음 분부터 선택 가능. 과거 입력·시간 경과 시 저장 비활성화와 제출 직전 재검증.
- 현재 서버의 `/surveys/me`는 공개 질문만 반환한다. `/answers/me`는 정렬과 커서 조건이 달라 누락 가능성이 있으므로 백엔드 보완이 필요하다. 상세 계약은 handoff에 기록했다.

개인 페이지 전용 읽기 fixture는 `node scripts/preview-personal-api.mjs`로 실행한다. 기존 fixture와 같은 3100 포트를 사용하므로 동시에 실행하지 않는다. `API_PROXY_TARGET=http://127.0.0.1:3100` 설정은 이 검증에만 사용하며, fixture가 발급하는 세션은 실제 계정이 아니다.

검증: 테스트 37개, `typecheck`, `lint`, production 빌드 통과. 360 / 390 / 768 / 1024 / 1440px에서 내 질문·내 답변·질문 작성의 15개 조합을 검사해 가로 넘침이 없음을 확인했다. 데스크톱·모바일 답변 그룹 화면도 확인했다. 이 결과는 fixture 기반이며 실제 백엔드 E2E 검증을 의미하지 않는다.

## Google OAuth 설정

1. Google Console의 Web application OAuth Client에서 승인된 Redirect URI에 `http://localhost:5173/auth/google/callback`을 등록합니다. 사용자 확인에 따라 개발 Console 등록은 완료되었습니다.
2. 프론트엔드 `.env.development.local`의 공개 Client ID와 백엔드 `GOOGLE_CLIENT_ID`를 동일하게 설정합니다. 로컬 설정은 Git에서 제외됩니다.
3. 백엔드 `.env`의 `GOOGLE_REDIRECT_URL`도 `http://localhost:5173/auth/google/callback`이어야 합니다. 로컬 값은 승인에 따라 변경했습니다. 실행 중인 백엔드는 재시작해야 반영됩니다.
4. 코드 교환 API는 기존 `http://localhost:8080/auth/google`을 유지합니다. 프론트엔드는 `/api/auth/google?code=...`를 통해 호출합니다.
5. 운영은 별도 OAuth Client 또는 승인된 운영 URI를 구성하고, `.env.production.local`/CI에 공개 Client ID와 운영 프론트엔드 callback 주소를 지정합니다. 운영값을 임의로 추정하지 않습니다.

Google 버튼은 `openid email profile` 범위의 authorization code 흐름을 시작합니다. 임의의 state와 안전한 내부 복귀 경로를 해당 탭의 sessionStorage에 잠시 저장하고 callback에서 일치·만료 여부를 검사한 뒤 삭제합니다. 인증 코드와 AccessToken, Client Secret은 저장소에 넣지 않습니다. Client Secret과 Google 토큰 교환은 백엔드가 담당합니다. callback 오류/취소 시에는 자동 교환 재시도 없이 로그인 화면으로 돌아갈 수 있습니다.

Redirect URI는 브라우저에서 연 origin과도 일치해야 하므로 `127.0.0.1:5173` 대신 `localhost:5173`으로 접근하세요. 설정 근거: [Google의 authorization code 및 state 문서](https://developers.google.com/identity/protocols/oauth2/web-server).
