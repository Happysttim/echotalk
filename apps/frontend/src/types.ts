export type Survey = {
  id: string;
  title: string;
  content: string;
  is_public: boolean;
  closed: boolean;
  author_id: string;
  expires_at: string;
  created_at: string;
  updated_at: string;
};

export type Answer = {
  id: string;
  survey_id: string;
  author_id: string;
  is_anonymous: boolean;
  anonymous: string;
  content: string;
  rate_up: number;
  created_at: string;
  updated_at: string;
};

export type Page<T> = { items: T[]; nextCursor: string; hasNext: boolean };
export type AuthUser = { id: string } | null;
