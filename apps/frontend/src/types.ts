export type Author = { id: string; nickname: string; hash_number: number };

export type Survey = {
  id: string;
  title: string;
  content: string;
  is_public: boolean;
  closed: boolean;
  author_id: string;
  author?: Author | null;
  expires_at: string;
  created_at: string;
  updated_at: string;
};

export type Answer = {
  id: string;
  survey_id: string;
  author_id: string;
  author?: Author | null;
  is_anonymous: boolean;
  anonymous: string;
  content: string;
  rate_up: number;
  created_at: string;
  updated_at: string;
};

export type Page<T> = { items: T[]; nextCursor: string; hasNext: boolean };
// The personal answer feed embeds question metadata, not its content/author.
export type SurveyAnswers = Omit<Survey, 'content' | 'author'> & {
  answers: Answer[];
};
export type MyAnswerDetail = {
  survey_answers: Omit<SurveyAnswers, 'answers'> & { answers: Answer[] | null };
};
export type AuthUser = { id: string } | null;
