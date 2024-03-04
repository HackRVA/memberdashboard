import { inject } from '@angular/core';
import {
  ActivatedRouteSnapshot,
  CanActivateFn,
  Router,
  RouterStateSnapshot,
  UrlTree,
} from '@angular/router';
import { AuthService } from '../services';
import { Observable, map } from 'rxjs';
import { AuthUser } from '../types/auth.types';

export const anonymousGuard: CanActivateFn = (
  route: ActivatedRouteSnapshot,
  state: RouterStateSnapshot
): Observable<boolean | UrlTree> => {
  const authService: AuthService = inject(AuthService);
  const router: Router = inject(Router);

  return authService.user$.pipe(
    map((user: AuthUser) => {
      return user.isLogin ? router.createUrlTree(['home']) : true;
    })
  );
};
