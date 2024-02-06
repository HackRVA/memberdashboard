import { Routes } from '@angular/router';
import { LoginComponent } from './pages/login';
import { NotFoundComponent } from './pages/not-found';
import { HomeComponent } from './pages/home';

export const routes: Routes = [
  {
    path: 'login',
    component: LoginComponent,
  },
  {
    path: '',
    children: [
      { path: 'home', component: HomeComponent },
      { path: '404', component: NotFoundComponent },
      { path: '**', redirectTo: '404' },
    ],
  },
];
