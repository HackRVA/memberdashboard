// testing
import { expect, assert } from '@open-wc/testing';

// rxjs
import { of } from 'rxjs';

// memberdashboard
import { LoginRequest } from './../src/auth/types/api/login-request';
import { AuthService } from './../src/auth/services/auth.service';
import { AuthResponse } from '../src/auth/types/api/auth-response';

describe('AuthService', () => {
  let mockService: AuthService;

  before(async () => {
    mockService = setupMockService();
  });

  it('should be defined', () => {
    assert.instanceOf(mockService, AuthService);
  });

  it('should be able to login', () => {
    // ARRANGE
    let actualResponse: string = null;
    const request: LoginRequest = {
      email: 'babyyoda@gmail.com',
      password: 'Do@or@do@not.There@is@no@try',
    };
    const expectedResponse: string = request.email + request.password;

    // ACT
    mockService.login(request).subscribe((response: AuthResponse) => {
      actualResponse = response.token;
    });

    // ASSERT
    expect(actualResponse).not.be.null;
    expect(actualResponse).equal(expectedResponse);
  });
});

function setupMockService(): AuthService {
  let service = new AuthService();

  service.login = (request: LoginRequest) => {
    return of({
      token: request.email + request.password,
    });
  };

  return service;
}
