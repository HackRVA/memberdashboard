// memberdashboard
import { diContainer } from './di-container';

export function Injectable(token: string): Function {
  return (targetProvider: { new () }): void => {
    diContainer.providers[token] = new targetProvider();
  };
}
