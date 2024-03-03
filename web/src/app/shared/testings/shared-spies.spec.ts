import { HttpClient } from '@angular/common/http';
import { AuthService, VersionService } from '../services';

export class SharedSpies {
  static createVersionSpy(): jasmine.SpyObj<VersionService> {
    return jasmine.createSpyObj<VersionService>('VersionService', [
      'getVersion',
    ]);
  }

  static createHttpClientSpy(): jasmine.SpyObj<HttpClient> {
    return jasmine.createSpyObj<HttpClient>('HttpClient', [
      'get',
      'post',
      'put',
      'delete',
    ]);
  }

  static createAuthServiceSpy(): jasmine.SpyObj<AuthService> {
    return jasmine.createSpyObj<AuthService>('AuthService', [
      'logout',
      'login',
      'register',
    ]);
  }
}
