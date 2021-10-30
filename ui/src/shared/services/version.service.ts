// rxjs
import { Observable } from 'rxjs';

// memberdashboard
import { HTTPService } from './http.service';
import { ENV } from '../../env';
import { VersionResponse } from '../types/api/version-response';
import { Injectable } from '../di';

@Injectable('version')
export class VersionService extends HTTPService {
  private readonly versionUrlSegment: string = ENV.api + '/version';

  getVersion(): Observable<VersionResponse> {
    return this.get<VersionResponse>(this.versionUrlSegment);
  }
}
