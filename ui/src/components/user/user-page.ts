// lit element
import {
  LitElement,
  html,
  customElement,
  TemplateResult,
  CSSResult,
  property,
} from "lit-element";

// memberdashboard
import { UserService } from "../../service";
import { userPageStyles } from "./styles";
import "../shared/card-element";
import "./user-detail";
import { authUser } from "../../auth-user";

@customElement("user-page")
export class UserPage extends LitElement {
  @property({ type: String })
  email: string = "";

  userService: UserService = new UserService();

  static get styles(): CSSResult[] {
    return [userPageStyles];
  }

  firstUpdated(): void {
    this.email = authUser.getValue().email;
  }

  displayUserDetail(): TemplateResult | void {
    if (this.email) {
      return html` <user-detail .email=${this.email}> </user-detail> `;
    }
  }

  render(): TemplateResult {
    return html` <card-element> ${this.displayUserDetail()} </card-element> `;
  }
}
