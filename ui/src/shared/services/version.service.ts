// rxjs
import { Observable } from 'rxjs';
import { ENV } from '../../env';
import { VersionResponse } from '../types/api/version-response';

// memberdashboard
import { HTTPService } from './http.service';

export class VersionService extends HTTPService {
  private readonly versionUrlSegment: string = ENV.api + '/version';

  getVersion(): Observable<VersionResponse> {
    return this.get<VersionResponse>(this.versionUrlSegment);
  }
}
