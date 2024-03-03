import {
  HttpErrorResponse,
  HttpEvent,
  HttpHandlerFn,
  HttpInterceptorFn,
  HttpRequest,
  HttpStatusCode,
} from '@angular/common/http';
import { AuthService } from '../services';
import { inject } from '@angular/core';
import { AuthUser } from '../types';
import { Observable, catchError, throwError } from 'rxjs';
import { Router } from '@angular/router';

export const authInterceptor: HttpInterceptorFn = (
  req: HttpRequest<unknown>,
  next: HttpHandlerFn
): Observable<HttpEvent<any>> => {
  // https://angular.io/guide/http-interceptor-use-cases
  const router: Router = inject(Router);
  const authService: AuthService = inject(AuthService);
  const auth: AuthUser = authService.user$.getValue();

  if (auth.isLogin) {
    const authRequest: HttpRequest<any> = req.clone({
      headers: req.headers.set('Authorization', 'Bearer ' + auth.token),
    });

    return next(authRequest).pipe(
      catchError((error: HttpErrorResponse) => {
        if (error.status === HttpStatusCode.Unauthorized) {
          authService.resetAndNext();
          router.navigate(['login']);
        }

        return throwError(() => error);
      })
    );
  }

  return next(req).pipe(
    catchError((error: HttpErrorResponse) => throwError(() => error))
  );
};
