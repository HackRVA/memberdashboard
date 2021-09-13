// vaadin
import { Router, Route, Context, Commands } from '@vaadin/router';
import { RoleGuard } from './auth/guard/role.guard';

const routes: Route[] = [
  {
    path: '/',
    component: 'member-dashboard',
    action: async () => {
      await import('./index');
    },
    children: [
      {
        path: 'home',
        component: 'home-page',
        action: async () => {
          await import('./home/components/home-page');
        },
      },
      {
        path: 'user',
        component: 'user-page',
        action: async () => {
          await import('./user/components/user-page');
        },
      },
      {
        path: 'report',
        component: 'report-page',
        action: async (context: Context, commands: Commands) => {
          await import('./report/components/report-page');
          return await new RoleGuard().canActivate(context, commands);
        },
      },
      {
        path: 'member',
        component: 'member-page',
        action: async (context: Context, commands: Commands) => {
          await import('./member/components/member-page');
          return await new RoleGuard().canActivate(context, commands);
        },
      },
      {
        path: 'resource',
        component: 'resource-page',
        action: async (context: Context, commands: Commands) => {
          await import('./resource/components/resource-page');
          return await new RoleGuard().canActivate(context, commands);
        },
      },
      {
        path: '(.*)',
        component: 'not-found',
        action: async () => {
          await import('./shared/components/not-found');
        },
      },
    ],
  },
];
const outlet = document.getElementById('outlet');
export const router = new Router(outlet);
router.setRoutes(routes);
