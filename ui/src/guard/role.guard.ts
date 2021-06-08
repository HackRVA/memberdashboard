import { Commands, Context, RedirectResult } from "@vaadin/router";
import { isAdmin } from "../function";

export class RoleGuard {
  async canActivate(
    context: Context,
    commands: Commands
  ): Promise<RedirectResult | void> {
    const isAuthorized = await this.isAuthorized();

    if (!isAuthorized) {
      return commands.redirect("/");
    }
  }

  private isAuthorized(): Promise<boolean> {
    return new Promise((resolve: Function, reject: Function) => {
      resolve(isAdmin() === true);
    });
  }
}
