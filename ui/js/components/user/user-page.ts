// lit element
import {
  LitElement,
  html,
  customElement,
  TemplateResult,
  CSSResult,
  property,
} from "lit-element";

// membership
import { UserProfile } from "./types";
import { UserService } from "../../service";
import { userPageStyles } from "./styles";
import "../shared/card-element";
import "./user-detail";

@customElement("user-page")
export class UserPage extends LitElement {
  @property({ type: String })
  email: string = "";

  userService: UserService = new UserService();

  static get styles(): CSSResult[] {
    return [userPageStyles];
  }

  firstUpdated(): void {
    this.getUser();
  }

  getUser(): void {
    this.userService.getUser().subscribe({
      next: (result: UserProfile) => {
        const { email } = result;
        this.email = email;
        this.requestUpdate();
      },
    });
  }

  displayUserDetail(): TemplateResult {
    if (this.email) {
      return html` <user-detail .email=${this.email}> </user-detail> `;
    } else {
      return html``;
    }
  }

  render(): TemplateResult {
    return html` <card-element> ${this.displayUserDetail()} </card-element> `;
  }
}
