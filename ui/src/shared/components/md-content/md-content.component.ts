// vaadin
import { Router } from '@vaadin/router';

// memberdashboard
import { AuthService } from '../../../auth/services/auth.service';
import { coreStyle } from '../../styles/core.style';
import { authUser$ } from '../../../auth/auth-user';
import { mdContentStyle } from './md-content.style';
import { customElement, property } from 'lit/decorators.js';
import { CSSResult, html, LitElement, TemplateResult } from 'lit';
import { TabIndex } from '../../types/custom/tab-index';

@customElement('md-content')
export class MDContent extends LitElement {
  @property({ type: String })
  email: string;

  // version: VersionResponse;

  authService: AuthService = new AuthService();
  // versionService: VersionService = new VersionService();

  static get styles(): CSSResult[] {
    return [mdContentStyle, coreStyle];
  }

  firstUpdated(): void {
    // this.getVersion();
  }

  goToHome(): void {
    Router.go('/home');
  }

  goToUser(): void {
    Router.go('/user');
  }

  goToReports(): void {
    Router.go('/report');
  }

  goToMembers(): void {
    Router.go('/member');
  }

  goToResources(): void {
    Router.go('/resource');
  }

  getTabIndex(pathName: string): number {
    switch (pathName) {
      case '/home':
        return TabIndex.home;
      case '/user':
        return TabIndex.user;
      case '/report':
        return TabIndex.reports;
      case '/member':
        return TabIndex.members;
      case '/resource':
        return TabIndex.resources;
      default:
        return -1;
    }
  }

  handleLogout(): void {
    this.authService.logout().subscribe({
      next: (response: null) => {
        localStorage.removeItem('jwt');
        authUser$.next({ login: false, email: null });
        window.location.href = '/home';
      },
    });
  }

  // getVersion(): void {
  //   this.versionService.getVersion().subscribe({
  //     next: (response: VersionResponse) => {
  //       this.version = response;
  //       this.requestUpdate();
  //     },
  //   });
  // }

  // generateVersionNumber(version: VersionResponse): TemplateResult {
  //   return html` <span> ${version?.major}.<b>${version?.build}</b> </span> `;
  // }

  displayLogout(): TemplateResult {
    return html`
      <mwc-button
        class="logout"
        slot="actionItems"
        label="Log out"
        icon="logout"
        @click=${this.handleLogout}
      ></mwc-button>
    `;
  }

  renderAdminTabs(): TemplateResult | void {
    if (true) {
      return html`
        <mwc-tab
          label="Reports"
          icon="show_chart"
          @click=${this.goToReports}
        ></mwc-tab>
        <mwc-tab
          label="Members"
          icon="people"
          @click=${this.goToMembers}
        ></mwc-tab>
        <mwc-tab
          label="Resources"
          icon="devices"
          @click=${this.goToResources}
        ></mwc-tab>
      `;
    }
  }

  render(): TemplateResult {
    return html`
      <mwc-top-app-bar-fixed centerTitle>
        <div slot="title">Member Dashboard</div>
        <div slot="actionItems">${this.email}</div>
        ${this.displayLogout()}
      </mwc-top-app-bar-fixed>
      <mwc-tab-bar activeIndex=${this.getTabIndex(window.location.pathname)}>
        <mwc-tab label="Home" icon="home" @click=${this.goToHome}></mwc-tab>
        <mwc-tab
          label="User"
          icon="account_circle"
          @click=${this.goToUser}
        ></mwc-tab>
        ${this.renderAdminTabs()}
      </mwc-tab-bar>
      <slot> </slot>
    `;
  }
}
