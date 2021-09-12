import { html } from 'lit';
import { fixture, expect } from '@open-wc/testing';

import { MemberDashboard } from '../src/member-dashboard';
import '../src/member-dashboard.js';

describe('MemberDashboard', () => {
  let element: MemberDashboard;
  beforeEach(async () => {
    element = await fixture(html`<member-dashboard></member-dashboard>`);
  });

  it('renders a h1', () => {
    const h1 = element.shadowRoot!.querySelector('h1')!;
    expect(h1).to.exist;
    expect(h1.textContent).to.equal('My app');
  });

  it('passes the a11y audit', async () => {
    await expect(element).shadowDom.to.be.accessible();
  });
});
