// vaadin
import { Router, Route } from '@vaadin/router';

const routes: Route[] = [
  {
    path: '/',
    component: 'member-dashboard',
    action: async () => {
      await import('./index');
    },
  },
];
const outlet = document.getElementById('outlet');
export const router = new Router(outlet);
router.setRoutes(routes);
