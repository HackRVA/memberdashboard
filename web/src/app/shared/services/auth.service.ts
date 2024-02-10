import { Injectable } from '@angular/core';
import { BehaviorSubject, Observable, lastValueFrom, of } from 'rxjs';
import { catchError, map } from 'rxjs/operators';
import {
  AuthResponse,
  AuthUser,
  LoginRequest,
  RegisterRequest,
} from '../types';
import { JwtPayload, jwtDecode } from 'jwt-decode';
import {
  HttpClient,
  HttpErrorResponse,
  HttpHeaders,
} from '@angular/common/http';
import { LocalStorageService } from './localstorage.service';

@Injectable({
  providedIn: 'root',
})
export class AuthService {
  private readonly _authUrlSegment: string = '/api/auth';
  private readonly _userUrlSegment: string = '/api/user';
  private readonly _sessionAuthKey: string = 'memberdashboard';

  user$ = new BehaviorSubject<AuthUser>({
    token: null,
    isLogin: false,
    email: null,
    isAdmin: false,
  });
  constructor(
    private readonly http: HttpClient,
    private readonly localStorageService: LocalStorageService
  ) {}

  register(request: RegisterRequest): Observable<void> {
    return this.http.post<void>(this._authUrlSegment + '/register', request);
  }

  login(request: LoginRequest): Observable<AuthResponse> {
    return this.http.post<AuthResponse>(
      this._authUrlSegment + '/login',
      {},
      {
        headers: new HttpHeaders({
          'Content-Type': 'application/json',
          Authorization:
            'Basic ' + btoa(`${request.email}:${request.password}`),
        }),
      }
    );
  }

  validateSession(): Promise<any> {
    const jwt: string = this.localStorageService.get(this._sessionAuthKey);

    if (jwt) {
      this.rehydrateAuthUser(jwt);
      return lastValueFrom(
        this.http.get(this._userUrlSegment).pipe(
          catchError((_: HttpErrorResponse) => {
            this.rehydrateAuthUser(null);
            return of(null);
          })
        )
      );
    }

    return lastValueFrom(of(null));
  }

  logout(): Observable<void> {
    return this.http.delete<null>(this._authUrlSegment + '/logout').pipe(
      map((_) => {
        this.localStorageService.delete(this._sessionAuthKey);
      })
    );
  }

  getSessionKey(): string {
    return this._sessionAuthKey;
  }

  private rehydrateAuthUser(token: string): void {
    this.localStorageService.update(this._sessionAuthKey, token);

    try {
      const claims: JwtPayload & {
        ID: string;
        Name: string;
        Groups: string[];
      } = jwtDecode(token);
      this.user$.next({
        isLogin: true,
        email: claims.ID,
        token: token,
        isAdmin: claims.Groups?.includes('admin'),
      });
    } catch (error) {
      // fail quietly. no need to pollute the console
    }
  }
}
