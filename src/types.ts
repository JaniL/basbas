export interface FetchResult {
  ok: boolean;
  status: number;
  body: string;
  error?: string;
}

export interface CachedValue<T> {
  value: T;
  expiresAt: number;
}

export interface DayAvailability {
  date: Date;
  markup: string | null;
  error?: string;
}
