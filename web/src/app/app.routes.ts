import { Routes } from '@angular/router';
import { anonymousGuard, authGuard } from './shared/guards';

export const routes: Routes = [
  {
    path: 'login',
    canActivate: [anonymousGuard],
    loadComponent: () => import('./pages/login').then((m) => m.LoginComponent),
  },
  {
    path: '',
    canActivateChild: [authGuard],
    children: [
      {
        path: '',
        redirectTo: 'home',
        pathMatch: 'full',
      },
      {
        path: 'home',
        loadComponent: () =>
          import('./pages/home').then((m) => m.HomeComponent),
      },
      {
        path: 'user',
        loadComponent: () =>
          import('./pages/user').then((m) => m.UserComponent),
      },
      {
        path: 'report',
        loadComponent: () =>
          import('./pages/report').then((m) => m.ReportComponent),
      },
      {
        path: 'member',
        loadComponent: () =>
          import('./pages/member').then((m) => m.MemberComponent),
      },
      {
        path: 'resource',
        loadComponent: () =>
          import('./pages/resource').then((m) => m.ResourceComponent),
      },
      {
        path: '404',
        loadComponent: () =>
          import('./pages/not-found').then((m) => m.NotFoundComponent),
      },
      {
        path: '**',
        redirectTo: '404',
      },
    ],
  },
];
