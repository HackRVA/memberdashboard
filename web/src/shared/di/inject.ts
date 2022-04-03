// memberdashboard
import { diContainer } from './di-container';

export function Inject(token: string) {
  return (targetProvider: any, key: string): any => {
    Object.defineProperty(targetProvider, key, {
      get: () => diContainer.resolve(token),
      enumerable: true,
      configurable: true,
    });
  };
}
