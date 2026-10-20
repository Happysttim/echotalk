import { ArrowUpRight } from 'lucide-react';
import type { SurveyAnswers } from '../../types';
import { AppLink } from '../ui/Button';
import { AuthorLine } from './QuestionRow';

export function MyAnswerGroup({ survey }: { survey: SurveyAnswers }) {
  const href = `/questions/${encodeURIComponent(survey.id)}`;
  return (
    <section
      className="my-answer-group"
      aria-label={survey?.title || '답변한 질문'}
    >
      <h2>
        <AppLink href={href}>
          {survey.title || '제목 없는 질문'}
          <ArrowUpRight size={18} aria-hidden />
        </AppLink>
      </h2>
      <ul className="my-answer-list">
        {survey.answers.map((answer) => (
          <li key={answer.id} className="my-answer-item">
            <AuthorLine
              author={answer.author}
              anonymous={
                answer.is_anonymous ? answer.anonymous || '익명' : undefined
              }
              date={answer.created_at}
            />
            <p className="prose">{answer.content}</p>
          </li>
        ))}
      </ul>
    </section>
  );
}
