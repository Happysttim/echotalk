import type { SurveyAnswers } from '../types';
import { appendUnique } from '../hooks/useCursorFeed';

// Keep question/answer order stable when a question spans multiple pages.
export function mergeSurveyAnswers(
  old: SurveyAnswers[],
  incoming: SurveyAnswers[],
) {
  const groups = new Map(old.map((group) => [group.id, group]));
  for (const group of incoming) {
    const previous = groups.get(group.id);
    groups.set(group.id, {
      ...group,
      answers: appendUnique(previous?.answers ?? [], group.answers),
    });
  }
  return [...groups.values()];
}
