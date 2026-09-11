import { ArrowUpRight } from 'lucide-react';
import type { Answer } from '../../types';
import { useSurvey } from '../../hooks/useSurvey';
import { dateText } from '../../lib/date';
import { AppLink, Button } from '../ui/Button';

export function MyAnswerGroup({
  surveyId,
  answers,
}: {
  surveyId: string;
  answers: Answer[];
}) {
  // One component (and resource request) per unique question, including across pages.
  const { survey, error, reload } = useSurvey(surveyId);
  const href = `/questions/${encodeURIComponent(surveyId)}`;
  return (
    <section
      className="my-answer-group"
      aria-label={survey?.title || '답변한 질문'}
    >
      <h2>
        <AppLink href={href}>
          {survey
            ? survey.title
            : error
              ? '질문 정보를 불러오지 못했습니다.'
              : '질문 제목을 불러오는 중…'}
          <ArrowUpRight size={18} aria-hidden />
        </AppLink>
      </h2>
      {error && (
        <div className="question-lookup-error">
          <p className="help" role="status">
            질문이 삭제되었거나 접근할 수 없을 수 있습니다. 내가 작성한 답변은
            아래에서 확인할 수 있습니다.
          </p>
          <Button variant="quiet" onClick={reload}>
            질문 제목 다시 불러오기
          </Button>
        </div>
      )}
      <ul className="my-answer-list">
        {answers.map((answer) => (
          <li key={answer.id} className="my-answer-item">
            <time dateTime={answer.created_at}>
              {dateText(answer.created_at)}
            </time>
            <p className="prose">{answer.content}</p>
            <AppLink
              className="text-link"
              href={`${href}#answer-${encodeURIComponent(answer.id)}`}
            >
              질문에서 이 답변 보기
              <ArrowUpRight size={15} aria-hidden />
            </AppLink>
          </li>
        ))}
      </ul>
    </section>
  );
}
