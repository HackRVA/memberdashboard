// rxjs
import { Observable } from 'rxjs';
import { map } from 'rxjs/operators';

// memberdashboard
import { AuthResponse } from './../types/api/auth-response';
import { HTTPService } from '../../shared/services/http.service';
import { RegisterRequest } from '../types/api/register-request';
import { LoginRequest } from '../types/api/login-request';
import { ENV } from '../../env';
import { Injectable } from '../../shared/di';

@Injectable('auth')
export class AuthService extends HTTPService {
  private readonly authUrlSegment: string = ENV.api + '/auth';

  registerUser(request: RegisterRequest): Observable<void> {
    return this.post<void>(this.authUrlSegment + '/register', request);
  }

  login(request: LoginRequest): Observable<AuthResponse> {
    return this.post<AuthResponse>(
      this.authUrlSegment + '/login',
      request,
      'Basic'
    );
  }

  logout(): Observable<null> {
    return this.delete<null>(this.authUrlSegment + '/logout').pipe(
      map(() => null)
    );
  }
}
