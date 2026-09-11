# Echotalk — 반응형 화면 설계 및 프론트엔드 명세

작성일: 2026-09-10  
기준: 사용자가 마지막으로 제공한 Go 모델·API 명세. 특히 `/auth/register`의 `nickname`과 마지막 반환 항목 기준 커서 규칙을 반영했다. 2026-09-11 사용자 확인에 따라 메일 링크 기반 인증 흐름을 갱신했다.  
구현 대상: React + Tailwind CSS + Radix UI. 이 시안은 실제 서비스에 연결하지 않은 검토용 프로토타입이다. 모든 질문·작성자·날짜·추천 수는 예시 데이터다.

## 1. 디자인 원칙

- 흰 바탕, 짙은 녹색 Primary, 차콜 본문. 그라데이션을 사용하지 않는다.
- 반복 콘텐츠는 카드가 아니라 구분선으로 나눈 행이다. 패널은 인증 소개·작성 맥락 등 의미 있는 묶음에만 사용한다.
- 서비스 문구는 제공된 기능만 설명한다. 사용자 수, 답변 속도, 보안 수준, 성과, 인기·추천 지표를 추정하지 않는다.
- 사용자에게는 ‘질문’이라는 용어로 통일한다. `Survey`는 제목과 단일 텍스트 본문이므로 객관식 문항, 복수 문항, 설문 통계, 응답 진행률을 추가하지 않는다.
- 없는 기능의 입력창을 만들지 않는다. 검색, 태그, 북마크, 알림, 이미지 업로드, 신고, 답변 채택, 자동 저장은 현재 시안 범위에 없다. 내 질문·내 답변 목록은 추가된 /me API에 맞춰 제공한다.
- 답변 추천은 `RateUp`에 대응한다. 추천 취소 API가 없으므로 토글 동작을 확정하지 않는다.

## 2. 레퍼런스 수집과 적용

확인일: 2026-09-10. 아래는 실제 확인한 페이지다. 외부 이미지·로고·원문 콘텐츠는 복제하지 않았다. 적용은 Echotalk의 데이터와 API에 맞춘 설계 판단이다.

| 레퍼런스                                                                      | 참고 요소                                | Echotalk 적용 / 제외                                       |
| ----------------------------------------------------------------------------- | ---------------------------------------- | ---------------------------------------------------------- |
| [Stack Overflow Questions](https://stackoverflow.com/questions)               | 제목, 본문 요약, 작성자 중심의 질문 목록 | 선으로 나눈 행 구조. 조회 수·답변 수·태그는 제외           |
| [GitHub Community Discussions](https://github.com/orgs/community/discussions) | 목록·상세에서 토론 콘텐츠의 정보 위계    | 질문 뒤 답변이 이어지는 읽기 흐름. 카테고리·해결 표시 제외 |
| [Dribbble Forum](https://dribbble.com/tags/forum)                             | 포럼·커뮤니티 시각 레퍼런스 컬렉션       | 밀도와 여백의 보조 탐색 자료. 특정 작품 복제 아님          |
| [Dribbble Question Answer](https://dribbble.com/tags/question-answer)         | 질문·답변 관련 화면 사례 컬렉션          | 텍스트 위계 검토. FAQ 아코디언·선택 문항 제외              |

기술 참고: [Tailwind responsive design](https://tailwindcss.com/docs/responsive-design), [Radix Dialog](https://www.radix-ui.com/primitives/docs/components/dialog).

## 3. 디자인 토큰

| 토큰              | 값                                 | 용도                    |
| ----------------- | ---------------------------------- | ----------------------- |
| primary           | #116149                            | 등록, 로그인, 핵심 동작 |
| primary-hover     | #0D4E3B                            | 주요 버튼 hover         |
| foreground        | #202824                            | 제목·본문               |
| muted-foreground  | #626D66                            | 보조 본문               |
| background        | #FFFFFF                            | 페이지                  |
| border            | #DCE3DE                            | 콘텐츠 구분선           |
| input-border      | #C8D2CB                            | 입력 경계               |
| status-background | #EAF4ED                            | 진행 상태               |
| destructive       | #B52F34                            | 오류·파괴 동작          |
| focus             | #116149, 2px, offset 4px           | 키보드 포커스           |
| radius            | 4 / 6 / 8px                        | 배지 / 입력 / 패널      |
| spacing           | 4, 8, 12, 16, 20, 24, 32, 48, 64px | 공통 여백 체계          |

폰트: Pretendard 우선, Noto Sans KR, 시스템 sans-serif 순. 시안은 외부 폰트 파일을 포함하지 않아 설치 상태에 따라 대체 폰트를 사용한다. 실서비스에서는 허용 라이선스 확인 후 폰트 파일을 함께 배포할 수 있다.

본문 기본 16px / 줄간격 1.65~1.9. 읽기 본문 데스크톱 17px / 모바일 16px. 폼 입력은 모바일에서도 16px. 데스크톱 상세 제목 36px, 모바일 27px. 주요 버튼 48px, 입력 50px. 작은 조작의 터치 영역은 44px 이상으로 확장한다. 검토용 캔버스의 축소 배율은 실제 서비스 폰트 크기와 다르다.

## 4. 공통 반응형 규칙

| 구간        | 서비스 화면 배치                                   |
| ----------- | -------------------------------------------------- |
| 360~767px   | 한 열, 좌우 20px. 헤더 64px. 보조 열 숨김          |
| 768~1100px  | 축소된 두 열. 최소 본문 공간 확보, 바깥 여백 48px  |
| 1101px 이상 | 최대 본문 폭 1120px. 헤더 내부 최대 1280px         |
| 인증 폼     | 최대 440px. 데스크톱은 안내 패널+폼, 모바일은 폼만 |
| 질문 상세   | 데스크톱 본문 744px + 간격 64px + 우측 280px       |
| 작성 화면   | 데스크톱 본문 736px + 간격 72px + 우측 280px       |

- 기기 종류가 아니라 CSS 가용 폭을 기준으로 배치한다.
- 768px 미만에 보조 열을 숨겨도 필수 동작은 없어지지 않아야 한다. 질문 작성자 관리는 본문 아래에 다시 배치한다.
- 모바일 질문 상세의 답변 진입 버튼은 하단 고정. `env(safe-area-inset-bottom)`과 콘텐츠 하단 여백을 확보한다.
- 작성 화면의 등록 버튼은 문서 흐름에 둔다. 키보드 위 강제 고정으로 입력창을 가리지 않는다.
- 본문·닉네임·ID·긴 문자열은 줄바꿈을 허용한다. 피드 요약만 2줄로 줄인다. 전체 본문은 상세에서 읽는다.
- 200% 확대와 키보드만 사용하는 상황을 실서비스 인수 검사에 포함한다.

## 5. 페이지별 화면 명세

검토 시안은 페이지 선택 후 ‘나란히 / 데스크톱 / 모바일’로 전환한다. 화면 상태 선택에서 정상·오류·빈 상태 등을 바꾼다. 실제 프론트엔드의 경로는 아래 제안이며 검토 사이트에서는 쿼리 파라미터로 선택한다.

### 01. 메인 — `/`

데스크톱: 헤더 → 왼쪽 서비스 설명/질문 작성/질문 둘러보기, 오른쪽 짧은 보조 설명 → 최근 질문 목록 → 푸터.

모바일: 보조 설명 패널을 숨긴다. 36px 제목, CTA 두 개, 최근 질문을 한 열로 배치한다. 첫 화면에서 서비스 목적과 질문 탐색 진입을 확인할 수 있다.

문구: ‘질문을 올리고, 답변을 나눕니다.’ / ‘궁금한 내용을 텍스트로 작성하고 다른 사용자의 답변을 확인하세요.’

연결: `GET /surveys?type=_id&limit=10`, 반환된 질문 중 최대 3개 표시. 실제 정렬이 내림차순인지 검증한 뒤 ‘최근’ 라벨을 확정한다. 비로그인 질문 작성은 `/login?next=/questions/new`로 이동한다.

상태: 기본, 로딩 행, ‘등록된 질문이 없습니다.’, ‘목록을 불러오지 못했습니다.’ + 다시 시도.

### 02. 로그인 — `/login`

데스크톱: 안내 패널 440px + 폼 440px. Google 버튼 → 구분선 → 이메일 → 비밀번호(표시/숨김) → 로그인 → 비밀번호 변경/회원가입.

모바일: 안내 패널을 숨긴 최대 440px 단일 폼. 이메일과 비밀번호 레이블을 항상 표시한다.

연결: `POST /auth/local {email,password}`. Google 인가 코드 수신 후 `GET /auth/google?code=...`. 두 요청 모두 Authorization 헤더를 보내지 않는다. Google은 프론트엔드 `/auth/google/callback`으로 돌아온다. 개발 Redirect URI는 `http://localhost:5173/auth/google/callback`이며 Google Console 및 백엔드 GOOGLE_REDIRECT_URL과 일치시킨다. 공개 Client ID는 환경변수로 주입하고 Client Secret은 백엔드에만 둔다.

상태: 기본, 전송 중, ‘이메일 또는 비밀번호를 확인해 주세요.’, 네트워크 오류, Google 취소·실패. 전송 중 중복 클릭 방지. 성공 시 검증된 내부 `next` 경로로 복귀한다.

HTML 검토 시안의 버튼은 연결 전 안내만 제공하지만, React 구현의 Google 버튼은 OAuth authorization code 흐름에 연결되어 있다. state 검증, 취소/오류, 일회용 코드 중복 교환 방지와 안전한 next 복귀를 처리한다. 실제 계정으로 끝까지 로그인한 검증과는 구분한다.

### 03. 회원가입 — `/register`

Google 진입과 이메일 가입을 함께 제시한다. 이메일 가입은 세 단계다.

| 단계        | 입력·버튼                                                                                         | 연결                                                           |
| ----------- | ------------------------------------------------------------------------------------------------- | -------------------------------------------------------------- |
| 1 이메일    | 이메일 / 인증 메일 받기                                                                           | `GET /verify?email=...&verifyType=register`                    |
| 2 인증      | 메일 발송 안내 / 다시 받기 / 이메일 변경. 메일의 `/register?verify={token}` 링크를 열면 자동 확인 | `GET /code?code={token}&verifyType=register`                   |
| 3 가입 정보 | 닉네임 / 비밀번호 / 확인 / 회원가입                                                               | `POST /auth/register {email,password,nickname}` + VERIFY_TOKEN |

최신 Payload의 `nickname`을 반영했다. HashNumber를 사용자가 입력하지 않는다. 비밀번호 확인 값은 클라이언트 검증용이며 Payload에서 제외한다.

모바일은 동일한 진행 순서와 단계 표시를 유지한다. 세 번째 단계의 닉네임과 비밀번호는 세로로 배치한다.

인증 코드를 직접 입력하지 않는다. 메일 링크의 `verify` 쿼리 값을 그대로 `code` 쿼리로 전달한다. 검증 중에는 다음 단계 입력을 노출하지 않고, 성공 응답 후에만 가입 정보 단계로 이동한다. 실패·빈 토큰은 새 메일 요청으로 안내한다. 일회용 링크이므로 개발 모드의 effect 재실행에서도 같은 요청을 중복 전송하지 않는다.

메일 요청 이메일을 변경하면 발송 대기 상태를 초기화한다. 새 인증 성공 전 다음 단계로 진행하지 않는다. 성공 시 서버가 설정하는 HttpOnly VERIFY_TOKEN 쿠키를 브라우저가 저장하며 프론트엔드에서 읽거나 생성하지 않는다. 사용자 추가 계약인 GET /verify/me로 인증 이메일을 조회해 최종 단계에 수정 불가능하게 표시한다. 서버가 Payload 이메일과 쿠키의 이메일 일치를 검사한다. 표시값 자체는 인증의 증거가 아니다.

완료 문구: ‘회원가입이 완료되었습니다.’ + 질문 둘러보기. 현재 서버의 가입 응답은 accessToken과 세션 쿠키를 반환하므로 구현은 로그인 상태를 유지한다. 검토 시안의 완료 화면 자체는 실제 계정 생성 결과가 아니다.

현재 백엔드 LoginWithGoogle은 신규 Google 사용자를 생성하고 세션을 반환한다. Google 회원가입 후 닉네임 수집 경로는 제공되지 않았으므로 별도 온보딩 성공을 가정하지 않는다.

### 04. 비밀번호 변경 — `/change` (기존 `/password/reset`도 지원)

이 API는 기존 비밀번호로 변경하는 로그인 후 설정 화면이 아니라 이메일 인증 기반 재설정 흐름이다.

이메일 → 메일 링크 확인 → 새 비밀번호/확인 → 완료. 메일 링크는 `/change?verify={token}`이다. 가입과 같은 레이아웃을 사용하지만 Google 버튼과 닉네임은 없다. 링크 확인 성공 전에는 새 비밀번호 입력을 노출하지 않는다.

연결: `/verify?email=...&verifyType=password` → 메일의 `/change?verify={token}` → `/code?code={token}&verifyType=password` → `POST /auth/change {email,password}` + VERIFY_TOKEN. 경로는 change지만 verifyType 값은 password다. 모든 인증 요청은 AuthUnrequired이므로 Authorization 헤더를 제거한다. 인증 이메일은 /verify/me로 조회하며, 권한은 서버의 쿠키 검증으로 확인한다.

문구: ‘이메일로 가입한 계정의 비밀번호를 변경합니다.’ / ‘Google 계정은 Google에서 비밀번호를 관리합니다.’

로그인 중 진입 정책은 UI에서 정해야 한다. 최소 구현은 이메일 계정 복구 화면으로 제공하고 해당 요청에 Bearer 헤더를 자동 주입하지 않는 것이다. 기존 세션 전체 무효화·RefreshToken 폐기 정책은 서버 확인 대상이다.

### 05. 질문 목록 — `/questions`

데스크톱: 좌측 220px 탐색 안내, 우측 피드. 상단 제목 → 등록순/수정순 → 행 목록 → 질문 더 보기.

모바일: 좌측 영역 제거. 제목 오른쪽 질문 작성, 아래 정렬, 한 열 피드. 각 행은 상태·날짜 → 제목 → 본문 2줄 → 닉네임#해시번호. 답변 수 집계 API가 없어 답변 수를 표시하지 않는다.

연결: `GET /surveys?type=_id|updated_at&limit=10`. 10 또는 30만 허용하며 화면은 10을 사용한다. 정렬 변경 시 커서·누적 목록을 초기화한다. 서버 정렬 방향은 확인하고 라벨을 최종화한다.

최근 수정은 새로운 답변이 달렸다는 의미로 해석하지 않는다. `Survey.updated_at`과 답변 작성의 관계가 제공되지 않았기 때문이다.

상태: 로딩, 빈 목록, 초기 요청 오류, 더 보기 오류(기존 행 유지), 마지막 페이지. 페이지 번호·총 페이지 수는 없다.

### 06. 질문 뷰어와 답변 피드 — `/questions/:id`

질문: 상태(답변 받는 중/마감), 공개 라벨, 제목, 작성자, 날짜, 전체 텍스트, 마감일.

답변: 별도 섹션 제목 ‘답변’ → 작성자·날짜 → 본문 → 추천 → 권한별 수정/삭제. 누적 개수를 전체 답변 수처럼 표시하지 않는다. 답변 정렬 옵션 API가 없어 피드에 정렬 선택을 추가하지 않는다.

데스크톱: 우측에 로그인하고 답변 / 익명으로 답변. 로그인 상태는 답변 작성 / 익명으로 답변. 작성자에게 질문 수정·마감·삭제를 노출한다.

모바일: 질문과 답변을 한 열로 읽는다. 답변 진입은 하단 고정. 작성자 동작은 질문 본문 아래에 배치하여 숨겨지지 않도록 한다.

마감: UI 설계상 `closed || expires_at <= 현재시각`이면 답변 작성 비활성. 서버에서도 마감 조건을 검사해야 한다. 마감 상태에서도 기존 답변은 읽는다. 일반 마감화면과 ‘내가 마감한 질문’의 재개 정책은 추가 확인한다.

연결: `GET /surveys/{id}`, `GET /answers?cursor=...`, `POST /answers/rateup`. 질문 수정·마감은 PATCH에 전체 수정 필드와 기존 값을 보존한다. 삭제는 `DELETE /surveys {survey_id}`. 질문 삭제 시 답변의 연쇄 삭제/보존 방식은 확인 대상이다.

추천: 비로그인 클릭 시 로그인 안내. 로그인은 요청 진행 중 잠그고 서버 응답으로 수를 갱신한다. 임의로 추천 취소 요청을 만들지 않는다. 재조회 후 추천 상태를 재현하려면 `has_rated_up`에 해당하는 응답 계약이 필요하다.

일반 답변 수정·삭제: 현재 사용자와 Answer.author_id 일치 시 노출한다. 질문 소유자라는 이유만으로 다른 사용자의 답변을 수정·삭제하도록 설계하지 않는다. 시안 작성자 상태의 첫 답변은 동일 사용자의 답변 예시다.

익명 답변 삭제: 비밀번호 입력 확인창. 일반 답변 삭제: 작성자 세션 확인 후 삭제 확인창. 확인창에는 취소를 항상 둔다. 익명 비밀번호 검증 실패 시 모달과 입력 맥락을 유지한다.

### 07. 일반 답변 작성 — `/questions/:id/answer`

질문으로 돌아가기 → 답변 작성 제목 → 대상 질문 요약/전체 보기 → 현재 닉네임 → 답변 textarea → 취소/등록.

데스크톱은 오른쪽 안내 열, 모바일은 한 열. 비어 있거나 공백뿐이면 등록 비활성. 임의 최대 글자 수를 두지 않고 현재 입력 글자 수만 보여준다. 서버 제한을 확인한 뒤 일치시킨다.

Payload:

```json
{
  "survey_id": "<id>",
  "content": "<text>",
  "is_anonymous": false,
  "anonymous": "",
  "answer_password": ""
}
```

Authorization 필수. AuthorID는 서버가 설정하므로 클라이언트 Payload에 추가하지 않는다. 등록 실패 시 본문 유지. 성공 후 질문으로 복귀하고 새 답변을 조회한다. 응답에 새 answer ID가 있으면 해당 답변에 포커스 또는 앵커 이동.

수정 모드는 같은 입력 화면을 재사용한다. `PATCH /answers {survey_id,answer_id,content}`. 일반 답변에 한해 소유권이 확인된 동작으로 시작하고 익명 답변 수정은 정책 확인 전 노출하지 않는다.

### 08. 익명 답변 작성 — `/questions/:id/answer?mode=anonymous`

일반 답변과 같은 틀에 표시 이름·답변 비밀번호를 추가한다. 데스크톱은 두 입력을 2열, 모바일은 세로로 쌓는다. 로그인 없이 접근 가능하다.

Payload:

```json
{
  "survey_id": "<id>",
  "content": "<text>",
  "is_anonymous": true,
  "anonymous": "<display name>",
  "answer_password": "<password>"
}
```

비밀번호는 필수. 표시 이름 필수 여부는 미제공이므로 시안은 선택 입력이다. 서버 기본 표시 이름이 없다면 필수 여부를 확정해야 한다.

문구: ‘익명 답변을 삭제할 때 필요합니다.’ 로그인한 익명 답변도 삭제는 answer_password를 사용한다. ‘완전 익명’, ‘추적 불가’, ‘암호화 저장’ 등 검증되지 않은 보장은 사용하지 않는다.

비밀번호는 로컬 저장·주소·로그에 남기지 않는다. AnswerPassword는 JSON 응답에서 제외되어 있으므로 조회해서 복구할 수 있다고 안내하지 않는다.

### 09. 질문 작성·수정 — `/questions/new`, `/questions/:id/edit`

제목 → 질문 내용 → 공개 Switch → 답변 마감일 → 취소/질문 등록. 텍스트 기반이며 서식 툴바·문항 추가·이미지 첨부·설문 통계를 넣지 않는다.

데스크톱 본문 736px + 안내 280px. 모바일은 한 열, 모든 설정이 등록 전에 보인다.

POST Payload: `title, content, is_public, expires_at`. PATCH Payload: `survey_id, title, content, is_public, closed, expires_at`. 수정 시 기존 closed 값을 보존하고, 마감 동작에서만 명시적으로 바꾼다.

시안은 공개 기본값 true, 마감일 선택 필수를 제안한다. 백엔드의 만료일 생략/무기한 정책은 미제공이므로 확인해야 한다. 예시 날짜를 실제 기본값으로 고정하지 않는다.

시간: UI에서 KST로 명시한다. `2026-09-30 23:59 KST`는 `2026-09-30T14:59:00Z`다. 브라우저 로컬 시간대가 항상 한국이라고 가정하지 않는다. 실제 구현에서 입력값에 +09:00을 명시해 UTC로 변환하거나 시간대 라이브러리를 사용한다.

공개/비공개: is_public의 의미가 확정되지 않아 ‘링크로만 접근’, ‘작성자만 읽기’라고 적지 않는다. 공개 목록에서 비공개 항목 제외 및 상세 접근 제한은 서버가 책임져야 한다. 현재 GET 라우터에 Auth middleware가 없다는 사실만으로 안전한 비공개 접근을 보장할 수 없다.

### 추가 화면: 내 활동

계정 DropdownMenu에서 ‘내가 작성한 질문’과 ‘내가 작성한 답변’으로 진입한다. 두 페이지는 AuthRequired이며, 세션 복원 후 비로그인이면 검증된 next를 포함해 로그인으로 이동한다. 로그인 후 원래 내 활동 페이지로 복귀한다. 데스크톱은 기존 질문 목록의 안내 열+피드, 모바일은 한 열을 사용한다. 두 페이지를 오가는 내 활동 링크는 모바일에서도 숨기지 않는다.

#### 내가 작성한 질문 — `/my/questions`

기존 질문 목록의 QuestionRow, 등록순/수정순 Select, 로딩/빈/오류/더 보기 양식을 재사용한다. `GET /surveys/me?type=_id|updated_at&cursor=...&limit=10`을 Bearer 인증으로 호출한다. 질문 제목은 상세로 연결한다. 정렬 변경 시 이전 커서와 항목을 초기화한다. 현재 백엔드 구현은 `is_public: true` 필터를 유지하므로 이 페이지도 작성한 **공개 질문만** 반환한다. 비공개 질문까지 포함하려면 서버 필터 변경이 필요하다. 프론트엔드에서 전체 목록을 조회해 사용자 ID로 필터링하지 않는다.

#### 내가 작성한 답변 — `/my/answers`

`GET /answers/me?cursor=...`을 Bearer 인증으로 호출한다. 반환 답변을 survey_id별로 묶어 **질문 제목 → 들여쓴 내 답변 목록**으로 표시한다. 같은 질문의 답변이 다음 페이지에서 추가되면 기존 그룹에 합치고 ID 중복은 제거한다. 그룹은 받은 답변에서 질문이 처음 등장한 순서, 그룹 내부는 서버가 반환한 순서를 유지한다. 정렬 UI를 임의로 추가하지 않는다.

현재 답변 DTO에는 질문 제목이 없어 그룹당 한 번 `GET /surveys/{surveyId}`로 보완한다. 추가 페이지의 동일 그룹은 재조회하지 않는다. 제목 조회 실패 시에도 답변 내용·날짜·질문 이동 링크를 유지하고 제목 재조회 버튼을 제공한다. 질문 제목은 `/questions/{surveyId}`, 개별 ‘질문에서 이 답변 보기’는 `#answer-{answerId}`로 연결한다. 삭제/접근 불가 질문은 조회 실패 문구로 명확히 표시한다. 계정으로 식별되지 않는 익명 답변을 클라이언트에서 추정해 포함하지 않는다.

#### 내 활동 응답과 커서

- 내 질문: `{status:'ok', data:{surveys:Survey[]|null, skip:string, hasNext:boolean}}`. 첫 cursor는 빈 문자열, 다음은 서버 skip을 수정 없이 사용한다.
- 내 답변: `{status:'ok', data:{answers:Answer[]|null, skip:string, hasNext:boolean}}`. 첫 cursor는 빈 문자열, 다음은 서버가 반환한 **답변 ID 원문**이다. 공개 답변 피드의 base64 JSON 커서와 혼용하지 않는다.
- 서버가 null 배열을 반환하면 빈 목록으로 처리한다. 더 보기 실패는 기존 그룹과 커서를 보존한다.
- 서버 확인 사항: 현재 내 답변은 updated_at/\_id 내림차순 정렬인데 다음 페이지 조건은 `_id < cursor`만 적용한다. 수정 시각과 ID 순서가 다르면 누락 가능성이 있어 서버의 정렬·커서 정합성 검증이 필요하다. 클라이언트에서 skip을 임의로 보정하지 않는다.

#### 마감일 입력 변경

질문 작성·수정의 datetime-local에 KST 기준 동적인 min과 step=60을 적용한다. 분 단위 입력이므로 현재 분 다음 분부터 선택할 수 있다. 시간이 흐르면 최솟값을 갱신한다. 과거 값 직접 입력·기존 만료된 마감일은 필드 오류와 저장 버튼 disabled로 차단한다. 제출 직전에도 실제 현재 시각과 비교해 대기 중 지난 마감일의 전송을 막는다. 기존 만료된 질문을 수정할 때는 미래 마감일을 선택해야 저장 가능하며 closed 값은 임의로 바꾸지 않는다.

추가 인수 검사: 내 활동 직접 진입/로그인 복귀, 인증 헤더, 두 커서 형식, null/로딩/오류/더 보기, 페이지 간 동일 질문 병합, 제목 조회 실패/재시도, 계정 변경 시 이전 목록 제거, KST 날짜 경계와 시간 경과 후 저장 차단. 디자인 시안의 내 질문·내 답변 데스크톱/모바일 및 기존 작성 화면의 날짜 제한도 함께 갱신한다.

## 6. 커서 구현 — 현재 API에 맞춘 코드

아래는 공개 피드의 커서 구조를 설명하는 예시다. 실제 구현은 서버 응답의 `data.answers` / `data.surveys` 배열과 `skip`, `hasNext`를 사용한다. `/answers/me`의 원문 ID 커서와 혼용하지 않는다.

```ts
type AnswerCursor = { survey_id: string; answer_id: string };
type SurveyCursor = { survey_id: string; updated_at: string };

// 일반 JSON도 UTF-8로 인코딩. 서버의 base64url 디코더와 padding 정책 확인.
function encodeCursor(value: AnswerCursor | SurveyCursor): string {
  const bytes = new TextEncoder().encode(JSON.stringify(value));
  let binary = '';
  for (const byte of bytes) binary += String.fromCharCode(byte);
  return btoa(binary)
    .replace(/\+/g, '-')
    .replace(/\//g, '_')
    .replace(/=+$/, '');
}

// 첫 페이지도 cursor가 필수. cursor 생략은 현재 400.
const initialAnswers = new URLSearchParams({
  cursor: encodeCursor({ survey_id: surveyId, answer_id: '' }),
});
// GET /answers?${initialAnswers}

// 다음 커서는 페이지의 10번째 answer 기준.
const nextAnswerCursor =
  answers.length === 10
    ? encodeCursor({ survey_id: surveyId, answer_id: answers[9].id })
    : null;

// 질문 페이지는 마지막 반환 항목 기준.
const lastSurvey = surveys.at(-1);
const nextSurveyCursor = lastSurvey
  ? encodeCursor({
      survey_id: lastSurvey.id,
      updated_at: lastSurvey.updated_at,
    })
  : null;
// GET /surveys?type=updated_at&limit=10&cursor=...
```

- 질문 커서의 `updated_at`은 RFC3339Nano 문자열을 그대로 유지한다. `new Date(...).toISOString()`으로 재직렬화하면 나노초 정보가 밀리초로 잘릴 수 있다.
- 질문의 다음 커서는 마지막 항목 기준이다. 커서가 있다는 것만으로 다음 항목 존재를 보장하지 않는다. has_more/next_cursor 응답 유무와 마지막 페이지 판정 계약을 확인한다.
- 답변이 10개 미만이면 10번째 기준 커서를 만들지 않는다. 정확히 10개라면 다음 요청이 빈 목록일 수 있다. 페이지 크기·서버 끝 판정과 일치시킨다.
- 동일 ID를 중복 append하지 않는다. 수정순 목록은 페이지 사이 변경으로 항목 이동 가능성이 있으므로 응답 계약과 정렬 tie-breaker를 확인한다.
- 정렬 변경 시 오래된 응답은 AbortController 또는 요청 세대값으로 버린다. 목록과 커서를 새 정렬 기준으로 초기화한다.
- 클라이언트는 type에 `_id`, `updated_at`만, limit에 `10`, `30`만 보낸다. 50/100은 사용하지 않는다.

## 7. API 매핑 전체

| 요청                    | UI                                 | 헤더/쿠키 및 Payload                                                                                     |
| ----------------------- | ---------------------------------- | -------------------------------------------------------------------------------------------------------- |
| GET /auth/google        | Google 로그인·가입 진입            | Authorization 없음, code query                                                                           |
| POST /auth/local        | 이메일 로그인                      | Authorization 없음, email/password                                                                       |
| POST /auth/register     | 이메일 가입 최종                   | Authorization 없음, VERIFY_TOKEN, email/password/nickname                                                |
| POST /auth/change       | 이메일 비밀번호 재설정             | Authorization 없음, VERIFY_TOKEN, email/password                                                         |
| GET /auth/refresh       | AccessToken 재발급 내부 흐름       | REFRESH_TOKEN                                                                                            |
| GET /auth/logout        | 계정 메뉴 로그아웃                 | AccessToken + REFRESH_TOKEN                                                                              |
| DELETE /auth            | 현재 요청 페이지 범위 밖 계정 삭제 | AccessToken. 설정 화면 임의 추가 안 함                                                                   |
| POST /answers           | 일반·익명 답변 등록                | 선택 인증, survey_id/content/is_anonymous/anonymous/answer_password                                      |
| PATCH /answers          | 본인 일반 답변 수정                | 인증, survey_id/answer_id/content                                                                        |
| DELETE /answers         | 답변 삭제 확인창                   | 선택 인증, answer_id/answer_password                                                                     |
| POST /answers/rateup    | 답변 추천                          | 인증, answer_id                                                                                          |
| GET /answers/{answerId} | 특정 답변 조회·수정 로딩           | 인증 middleware 없음                                                                                     |
| GET /answers            | 답변 피드                          | cursor 필수, 최대 10개                                                                                   |
| GET /answers/me         | 내가 작성한 답변                   | AuthRequired, cursor는 원문 답변 ID                                                                      |
| POST /surveys           | 질문 작성                          | 인증, title/content/is_public/expires_at                                                                 |
| PATCH /surveys          | 질문 수정·마감                     | 인증, survey_id/title/content/is_public/closed/expires_at                                                |
| DELETE /surveys         | 질문 삭제 확인창                   | 인증, survey_id                                                                                          |
| GET /surveys/{surveyId} | 질문 상세                          | 인증 middleware 없음; 비공개 정책 별도 확인                                                              |
| GET /surveys            | 메인·질문 목록                     | type/cursor/limit                                                                                        |
| GET /surveys/me         | 내가 작성한 질문                   | AuthRequired, type/cursor/limit, 현재 공개 질문만 반환                                                   |
| GET /verify             | 가입/변경 인증 메일 요청           | Authorization 없음, email/verifyType                                                                     |
| GET /code               | 메일 링크 토큰 확인                | Authorization 없음, code=메일 링크의 verify 값 / verifyType=register 또는 password. 성공 시 VERIFY_TOKEN |

## 8. 인증·요청 처리

- AuthRequired: 유효 AccessToken 없으면 로그인으로 안내한다. 익명 흐름에 강제로 로그인 게이트를 씌우지 않는다.
- AuthNoRequired: 토큰이 없으면 헤더를 아예 생략한다. `Bearer undefined`, 빈 Bearer, 만료 토큰을 보내지 않는다. 만료 토큰이 있다면 적절한 갱신 실패 처리 후 비로그인 상태를 명확히 한다.
- AuthUnrequired: 전역 Axios/fetch 인터셉터가 헤더를 주입하지 않도록 제외한다. 최신 명세상 Authorization 헤더가 있으면 400이다.
- VERIFY_TOKEN/REFRESH_TOKEN은 쿠키로 유지한다. 프론트/백엔드가 다른 origin이면 credentials와 서버 CORS·SameSite·Secure 설정을 함께 맞춘다. 제공된 정보만으로 실제 쿠키 옵션을 확정하지 않는다.
- 동시 401에 refresh 요청을 하나로 합치고 무한 재시도하지 않는다. 실패하면 비로그인 상태로 전환한다. 토큰 응답 구조·회전 규칙은 확인 필요.
- 로그인 후 `next`는 허용된 내부 경로만 받는다. 외부 URL로의 오픈 리다이렉트를 허용하지 않는다.
- 비밀번호·인증 토큰은 분석 로그나 브라우저 저장소에 남기지 않는다. 메일 링크의 verify 값은 API 계약상 code query로 교환하고 완료 후 주소에서 제거한다. no-referrer 정책을 적용하고 인증 확인 요청은 no-store로 보낸다. 서버·프록시 로그의 민감 query 처리도 확인한다.
- 서버 응답이 불명확한 POST를 자동 반복하지 않는다. 중복 질문·답변 생성 방지와 재시도 전략은 idempotency 지원 여부와 함께 정한다.

## 9. 상태와 컴포넌트

| 상태         | 표현                          | 보존/후속 처리                      |
| ------------ | ----------------------------- | ----------------------------------- |
| 초기 조회    | 제목·행 Skeleton              | 실제 데이터처럼 카운트 표시 안 함   |
| 빈 질문      | 등록된 질문이 없습니다.       | 질문 작성 진입                      |
| 빈 답변      | 아직 답변이 없습니다.         | 답변 작성 진입                      |
| 조회 오류    | 불러오지 못했습니다.          | 다시 시도                           |
| 더 보기 오류 | 기존 목록 아래 오류           | 기존 행·커서 보존                   |
| 인증 오류    | 필드 인접 텍스트              | 비밀번호 노출 금지                  |
| 등록 중      | 진행 문구, 버튼 잠금          | 입력값 보존                         |
| 등록 오류    | 인라인 오류                   | 작성한 본문 보존                    |
| 마감         | 마감 라벨, 등록 비활성        | 기존 답변 조회 유지                 |
| 삭제         | AlertDialog                   | 취소 가능, 성공 후 목록 갱신        |
| 페이지 없음  | 질문을 찾을 수 없습니다.      | 목록으로 이동                       |
| 접근 불가    | 이 질문에 접근할 수 없습니다. | 비공개 정책에 맞춰 로그인/목록 진입 |

시안에서 직접 전환 가능한 상태: 메인·피드 기본/빈/로딩/오류; 가입·비밀번호 단계/완료/요청오류; 로그인 오류; 상세 비로그인/작성자/마감/답변없음/로딩/오류; 일반 답변 작성/수정/오류; 익명 답변 작성/오류; 질문 작성/수정/오류. 404, 접근 불가, 만료되는 세션, 실제 네트워크 전송 중 상태는 이 문서의 구현 명세이며 서버 연결까지 구현했다는 뜻은 아니다.

컴포넌트: ProductHeader, QuestionRow, AnswerRow, StatusBadge, AuthorLine, TextField, PasswordField, AnswerComposer, SurveyEditor, PaginationMore, LoadingRows, InlineError, EmptyState.

Radix: Select(정렬), Switch(공개), DropdownMenu(관리), AlertDialog(마감·삭제), Tabs(검토 도구). 레이아웃은 semantic HTML/React, 스타일은 Tailwind 토큰과 CSS로 구성한다. 카드 컴포넌트로 목록을 둘러싸지 않는다.

## 10. 구현 전 확인 목록

| 우선순위 | 확인 사항                                 | 디자인 영향                       |
| -------- | ----------------------------------------- | --------------------------------- |
| 필수     | 실제 성공/오류 응답 DTO, 배열 wrapper     | 정확한 화면 바인딩·끝 페이지 판정 |
| 필수     | 현재 로그인 사용자 반환 또는 /me          | 권한 버튼·계정 표시               |
| 필수     | author_id에 대응하는 nickname/hash DTO    | 작성자 표시, N+1 조회 방지        |
| 필수     | 비공개 조회 정책, 목록 제외               | 공개 설정 도움말·접근 제한        |
| 필수     | Google 신규 생성·닉네임 처리              | Google 가입 후 분기               |
| 필수     | 메일 링크 TTL/재발송 쿨다운/시도 수       | 링크 만료·재발송 안내             |
| 필수     | 비밀번호·닉네임·텍스트 유효성             | 실제 field 제약·오류 문구         |
| 필수     | 추천 중복·취소·내 추천 상태               | 추천 버튼 재조회 상태             |
| 필수     | PATCH/DELETE 소유권 검사                  | 본인 수정·삭제                    |
| 필수     | closed/만료 시 서버 답변 차단             | 마감 UI와 실제 쓰기 일치          |
| 확인     | 익명 이름 기본값·수정 가능 여부           | 필수 입력과 수정 노출             |
| 확인     | 만료일 생략·무기한·최소/최대 기한         | 날짜 입력 제약                    |
| 확인     | 질문 삭제 시 기존 답변 처리               | 삭제 확인 문구                    |
| 확인     | refresh 회전·비밀번호 변경 후 세션 무효화 | 재인증·로그아웃                   |
| 확인     | 정렬 방향·RFC3339Nano 비교·중복 방지      | 피드 연속 페이지 정확성           |

모델 검토 메모: User의 예시에는 json 태그가 없고 PasswordHash가 포함되어 있다. 모델을 그대로 직렬화하지 말고 안전한 응답 DTO를 사용해야 한다. RateUp.UserID의 json 태그는 `user_ud`이므로 의도한 `user_id`인지 확인한다. Session에는 사용자 식별 필드가 제시되지 않아 모든 기기의 세션 관리 기능을 UI에 추가하지 않는다.

## 11. 프론트엔드 인수 검사

- 360/390/768/1024/1440px에서 모든 요청 페이지 확인. 가로 넘침, 버튼 겹침, 긴 제목·닉네임 줄바꿈 검사.
- Tab 이동·Enter/Space 조작·Escape 닫기. 모달 포커스 잠금·원위치 복귀. 삭제 버튼을 색상만으로 구분하지 않기.
- input에 label, 오류에 aria-describedby/aria-invalid, 상태에 aria-live. 아이콘 단독 버튼에 이름 제공.
- 모바일 입력은 16px, 조작 영역 최소 44px. 하단 고정 답변 CTA가 콘텐츠와 safe area를 침범하지 않는지 확인.
- 회원가입 nickname Payload 포함, password confirm 미포함. /register?verify 및 /change?verify 직접 진입, 실패/빈 토큰/중복 요청 차단, 인증 성공 후에만 최종 폼 진입 확인. 이메일 입력만으로 권한을 부여하지 않음.
- AuthUnrequired 요청에서 Authorization 없음 확인. 익명 요청에 잘못된 Bearer 값 없는지 확인.
- 첫 answers 요청의 cursor 확인. 질문 limit 10/30, 수정순 RFC3339Nano 보존 확인.
- 마감 직전/직후, 비공개, 비작성자 수정, 잘못된 익명 비밀번호, 중복 추천은 API 통합 검사로 확인.
- 입력 중 뒤로가기/취소에는 변경 내용이 있을 경우 이탈 확인을 적용한다. 현재 시안에서는 데이터가 저장되지 않으며 이탈 확인은 실서비스 구현 항목이다.

검증 범위: 제공된 모델/API와 화면 설계 정합성 검토, TypeScript/배포 빌드 검증. 실제 백엔드 연동, Google 인증, 이메일 발송, 브라우저 시각 QA 및 스크린리더 검증을 완료한 결과물로 간주하지 않는다.
