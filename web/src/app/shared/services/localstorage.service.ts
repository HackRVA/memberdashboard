import { Injectable } from '@angular/core';
import { withSafeHandler } from '../functions';

@Injectable({
  providedIn: 'root',
})
export class LocalStorageService {
  get<T>(key: string): T {
    return withSafeHandler<T>('LocalStorageService_Get', () => {
      const item: string = localStorage.getItem(key);
      return item ? JSON.parse(item) : null;
    });
  }

  add<T>(key: string, value: T): void {
    withSafeHandler('LocalStorageService_Add', () =>
      localStorage.setItem(key, JSON.stringify(value))
    );
  }

  update<T>(key: string, value: T): void {
    withSafeHandler('LocalStorageService_Update', () => {
      const item: T = this.get(key);

      if (item) {
        this.delete(key);
        this.add(key, value);
        return;
      }

      this.add(key, value);
    });
  }

  delete(key: string): void {
    withSafeHandler('LocalStorageService_Delete', () =>
      localStorage.removeItem(key)
    );
  }
}
