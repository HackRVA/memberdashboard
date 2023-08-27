// lit element
import { CSSResult, html, LitElement, TemplateResult } from 'lit';
import { customElement } from 'lit/decorators.js';

// memberdashboard
import { coreStyle } from './../../../shared/styles/core.style';
import { homeDetailStyle } from './home-detail.style';
import { Inject } from '../../../shared/di';
import { ResourceService } from './../../../resource/services/resource.service';
import { OpenResourceRequest } from './../../../resource/types/api/open-resource-request';
import { showComponent } from '../../../shared/functions';
import { ToastMessage } from '../../../shared/types/custom/toast-msg';
import '../../../shared/components/gigachad-boss';

@customElement('home-detail')
export class HomeDetail extends LitElement {
  @Inject('resource')
  private resourceService: ResourceService;

  toastMsg: ToastMessage;

  isProcessing: boolean = false;

  static get styles(): CSSResult[] {
    return [coreStyle, homeDetailStyle];
  }

  openDoor(): void {
    this.activateProcessState(true);
    this.tryToOpenDoor();
  }

  tryToOpenDoor(): void {
    const request: OpenResourceRequest = {
      name: 'frontdoor',
    };

    this.resourceService.openResource(request).subscribe(
      () => {
        this.displayToastMsg('Front door is open (probably)');
        this.activateProcessState(false);
      },
      (error: null) => {
        this.displayToastMsg('Hrmmmmmm, it failed');
        this.activateProcessState(false);
      }
    );
  }

  activateProcessState(isProcessing: boolean): void {
    this.isProcessing = isProcessing;
    this.requestUpdate();
  }

  displayToastMsg(message: string): void {
    this.toastMsg = Object.assign({}, { message: message, duration: 4000 });
    this.requestUpdate();
    showComponent('#toast-msg', this.shadowRoot);
  }

  showProcessState(): TemplateResult | void {
    if (this.isProcessing) {
      return html`
        <mwc-circular-progress indeterminate></mwc-circular-progress>
      `;
    }
  }

  render(): TemplateResult {
    return html`
      <div id="home-detail" class="center">
        <section>
          <gigachad-boss />
        </section>
        <section>
          <mwc-button
            unelevated
            dense
            .disabled=${this.isProcessing}
            label="Open Front Door"
            @click=${this.openDoor}
          >
            ${this.showProcessState()}
          </mwc-button>
        </section>
      </div>
      <toast-msg id="toast-msg" .toastMsg=${this.toastMsg}> </toast-msg>
    `;
  }
}
