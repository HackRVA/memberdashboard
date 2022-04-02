// lit element
import { html } from 'lit';

// testing
import { fixture, expect, assert } from '@open-wc/testing';

// memberdashboard
import { MemberDashboard } from '../src/';
import '../src/member-dashboard';

describe('MemberDashboard', () => {
  let element: MemberDashboard;
  beforeEach(async () => {
    element = await fixture(html`<member-dashboard></member-dashboard>`);
  });

  it('is defined', () => {
    assert.instanceOf(element, MemberDashboard);
  });

  it('should not be signed in', () => {
    // ARRANGE
    const loadingContent = element.shadowRoot.querySelector('loading-content');
    const loginPage = loadingContent.querySelector('login-page');
    const header = loginPage.shadowRoot.querySelector('h1');

    // ASSERT
    expect(loginPage).not.be.null;
    expect(header.innerText).equal('Login');
  });
});
