import { Router, Route } from '@vaadin/router';

const routes: Route[] = [
    {
        path: '/build',
        component: 'member-dashboard',
        action: async() => {
            await import ('./index')
        },
        children: [
            {
                path: 'users',
                component: 'users-page',
                action: async() => {
                    await import ('./components/users-page')
                }
            },
            {
                path: 'members',
                component: 'members-page',
                action: async() => {
                    await import ('./components/members-page')
                }
            },
            {
                path: 'resources',
                component: 'resources-page',
                action: async() => {
                    await import ('./components/resources-page')
                }
            },
            {
                path: 'status',
                component: 'status-page',
                action: async() => {
                    await import ('./components/status-page')
                }
            }
        ]
    }
]
const outlet = document.getElementById('outlet');
export const router = new Router(outlet);
router.setRoutes(routes);