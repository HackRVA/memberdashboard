// memberdashboard
import { ContainerProvider } from '../types/custom/container-provider';

export class DIContainer {
  private _providers: { [key: string]: any } = {};

  get providers(): { [key: string]: any } {
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

  private findProvider(token: string): any | null {
    for (const [key, value] of Object.entries(this._providers)) {
      if (key === token) {
        return value;
      }
    }

    return null;
  }
}

export const diContainer = new DIContainer();
