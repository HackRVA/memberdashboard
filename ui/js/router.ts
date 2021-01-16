import { Router, Route } from "@vaadin/router";

const routes: Route[] = [
  {
    path: "/",
    component: "member-dashboard",
    action: async () => {
      await import("./index");
    },
    children: [
      {
        path: "home",
        component: "home-page",
        action: async () => {
          await import("./components/pages/home-page");
        },
      },
      {
        path: "users",
        component: "users-page",
        action: async () => {
          await import("./components/pages/users-page");
        },
      },
      {
        path: "members",
        component: "members-page",
        action: async () => {
          await import("./components/pages/members-page");
        },
      },
      {
        path: "resources",
        component: "resources-page",
        action: async () => {
          await import("./components/pages/resources-page");
        },
      },
      {
        path: "status",
        component: "status-page",
        action: async () => {
          await import("./components/pages/status-page");
        },
      },
    ],
  },
];
const outlet = document.getElementById("outlet");
export const router = new Router(outlet);
router.setRoutes(routes);
