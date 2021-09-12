// vaadin
import { Router, Route, Context, Commands } from '@vaadin/router';

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
        action: async () => {
          await import('./report/components/report-page');
        },
      },
      {
        path: 'member',
        component: 'member-page',
        action: async () => {
          await import('./member/components/member-page');
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
