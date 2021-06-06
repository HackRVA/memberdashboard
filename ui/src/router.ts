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
          await import("./components/home/home-page");
        },
      },
      {
        path: "user",
        component: "user-page",
        action: async () => {
          await import("./components/user/user-page");
        },
      },
      {
        path: "members",
        component: "members-page",
        action: async () => {
          await import("./components/members/members-page");
        },
      },
      {
        path: "resources",
        component: "resources-page",
        action: async () => {
          await import("./components/resources/resources-page");
        },
      },
      {
        path: "reports",
        component: "payments-page",
        action: async () => {
          await import("./components/payments/payments-page");
        },
      },
      {
        path: "(.*)",
        component: "not-found",
        action: async () => {
          await import("./components/shared/not-found");
        },
      },
    ],
  },
];
const outlet = document.getElementById("outlet");
export const router = new Router(outlet);
router.setRoutes(routes);
