import {
  HttpErrorResponse,
  HttpEvent,
  HttpHandlerFn,
  HttpInterceptorFn,
  HttpRequest,
} from '@angular/common/http';
import { AuthService } from '../services';
import { inject } from '@angular/core';
import { AuthUser } from '../types';
import { Observable, catchError, throwError } from 'rxjs';

export const authInterceptor: HttpInterceptorFn = (
  req: HttpRequest<unknown>,
  next: HttpHandlerFn
): Observable<HttpEvent<any>> => {
  // https://angular.io/guide/http-interceptor-use-cases
  const authService: AuthService = inject(AuthService);
  const auth: AuthUser = authService.user$.getValue();

  if (auth.isLogin) {
    const authRequest: HttpRequest<any> = req.clone({
      headers: req.headers.set('Authorization', 'Bearer ' + auth.token),
    });

    return next(authRequest);
  }

  return next(req).pipe(
    catchError((error: HttpErrorResponse) => throwError(() => error))
  );
};
