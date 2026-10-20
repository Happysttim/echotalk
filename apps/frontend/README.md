# Echotalk frontend

Echotalk은 사용자가 질문을 올리고 서로 답변을 나누는 서비스입니다. 이 디렉터리는 데스크톱과 모바일에서 사용할 수 있는 반응형 웹 프론트엔드를 담당합니다.

이메일·Google 로그인, 질문 작성 및 관리, 일반·익명 답변, 답변 추천, 내가 작성한 질문·답변 목록을 제공합니다. React / TypeScript / Vite / TailwindCSS / Radix UI / Zustand를 사용합니다.

전체 서비스 구성과 백엔드 실행 방법은 [프로젝트 README](../../README.md)를 참고하세요. 화면 기준은 [AGENTS.md](AGENTS.md), [프론트엔드 명세](docs/Echotalk-Frontend-Handoff.md), [반응형 디자인 시안](examples/Echotalk-Responsive-Design.html)입니다.

## 빠른 시작

저장소 설정 기준으로 Node.js **22.19.0**과 pnpm **9.4.0**을 준비합니다. 로그인과 데이터 조회·저장에는 실행 중인 Echotalk 백엔드가 필요합니다.

1. 저장소 루트(`echotalk`)에서 의존성을 설치한 뒤 프론트엔드 디렉터리로 이동합니다.

   ```sh
   pnpm install
   cd apps/frontend
   ```

2. [.env.example](.env.example)을 참고해 현재 디렉터리에 `.env.development.local`을 만듭니다. 기존 파일이 있다면 필요한 값만 수정합니다.

   ```dotenv
   VITE_API_BASE_URL=/api
   API_PROXY_TARGET=http://localhost:8080
   ```

   `API_PROXY_TARGET`은 실제 백엔드 주소로 지정하세요. Google 로그인은 아래의 [Google OAuth 설정](#google-oauth-설정)도 필요합니다.

3. 개발 서버를 실행하고 [http://localhost:5173](http://localhost:5173)에 접속합니다.

   ```sh
   pnpm dev
   ```

   이 디렉터리에서 실행하는 `pnpm dev`는 프론트엔드만 실행합니다. 프로젝트 루트에서 이미 전체 서비스를 실행했다면 중복 실행하지 않습니다. 개발 서버는 실행 터미널에서 `Ctrl+C`로 종료합니다.
