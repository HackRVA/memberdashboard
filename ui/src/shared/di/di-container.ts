// memberdashboard
import { DIProviders } from './../types/custom/di-providers';
import { ContainerProvider } from '../types/custom/container-provider';

export class DIContainer {
  private _providers: DIProviders = {};

  get providers(): DIProviders {
    return this._providers;
  }

  public resolve(token: string): any | Error {
    const matchedProvider = this.findProvider(token);

    if (matchedProvider) {
      return matchedProvider;
    } else {
      throw new Error(`No provider found for ${token}!`);
    }
  }

  public setProvider(provider: ContainerProvider): void {
    this._providers[provider.token] = provider.useValue;
  }

  private findProvider(token: string): any | undefined {
    return this.providers[token];
  }
}

export const diContainer = new DIContainer();
