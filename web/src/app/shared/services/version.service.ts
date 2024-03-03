import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { VersionResponse } from '../types/version-response';
import { HttpClient } from '@angular/common/http';

@Injectable({
  providedIn: 'root',
})
export class VersionService {
  private readonly _versionUrlSegment: string = '/api/version';

  constructor(private readonly http: HttpClient) {}

  getVersion(): Observable<VersionResponse> {
    return this.http.get<VersionResponse>(this._versionUrlSegment);
  }
}
