import { Routes } from '@angular/router';
import { LoginComponent } from './pages/login';
import { NotFoundComponent } from './pages/not-found';

export const routes: Routes = [
  {
    path: 'login',
    component: LoginComponent,
  },
  {
    path: '404',
    component: NotFoundComponent,
  },
  {
    path: '**',
    redirectTo: '404',
  },
];
