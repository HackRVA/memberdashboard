// vaadin
import { Router, Route, Context, Commands } from "@vaadin/router";

// memberdashboard
import { RoleGuard } from "./guard/role.guard";

const roleGuard = new RoleGuard();

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
        action: async (context: Context, commands: Commands) => {
          await import("./components/members/members-page");
          return await roleGuard.canActivate(context, commands);
        },
      },
      {
        path: "resources",
        component: "resources-page",
        action: async (context: Context, commands: Commands) => {
          await import("./components/resources/resources-page");
          return await roleGuard.canActivate(context, commands);
        },
      },
      {
        path: "reports",
        component: "payments-page",
        action: async (context: Context, commands: Commands) => {
          await import("./components/payments/payments-page");
          return await roleGuard.canActivate(context, commands);
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
