import { CachedValue } from "./types";

export class TtlCache<T> {
  private readonly values = new Map<string, CachedValue<T>>();

  constructor(private readonly ttlMs: number) {}

  get(key: string): T | undefined {
    const cached = this.values.get(key);
    if (!cached) {
      return undefined;
    }
    if (cached.expiresAt <= Date.now()) {
      this.values.delete(key);
      return undefined;
    }
    return cached.value;
  }

  set(key: string, value: T): void {
    this.values.set(key, {
      value,
      expiresAt: Date.now() + this.ttlMs
    });
  }
}
