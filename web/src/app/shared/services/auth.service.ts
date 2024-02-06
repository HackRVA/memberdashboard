import { Injectable } from '@angular/core';
import { BehaviorSubject, Observable } from 'rxjs';
import { map } from 'rxjs/operators';
import {
  AuthResponse,
  AuthUser,
  LoginRequest,
  RegisterRequest,
} from '../types';
// import { jwtDecode } from 'jwt-decode';
import { HttpClient, HttpHeaders } from '@angular/common/http';

@Injectable({
  providedIn: 'root',
})
export class AuthService {
  private readonly _authUrlSegment: string = '/api/auth';

  authUser$ = new BehaviorSubject<AuthUser>({
    token: null,
    isLogin: false,
    email: null,
  });
  constructor(private readonly http: HttpClient) {}

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

  logout(): Observable<null> {
    return this.http
      .delete<null>(this._authUrlSegment + '/logout')
      .pipe(map(() => null));
  }
}
