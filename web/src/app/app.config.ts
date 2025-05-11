import {
  ApplicationConfig,
  inject,
  provideAppInitializer,
} from '@angular/core';
import { provideRouter } from '@angular/router';

import { routes } from './app.routes';
import { provideAnimations } from '@angular/platform-browser/animations';
import { provideHttpClient, withInterceptors } from '@angular/common/http';
import { authInterceptor } from './shared/interceptors';
import { AuthService } from './shared/services';

function appInitSession(authService: AuthService): () => Promise<any> {
  return () => authService.validateSession();
}

export const appConfig: ApplicationConfig = {
  providers: [
    provideRouter(routes),
    provideAnimations(),
    provideHttpClient(withInterceptors([authInterceptor])),
    provideAppInitializer(() => {
      const initializerFn = appInitSession(inject(AuthService));
      return initializerFn();
    }),
  ],
};
