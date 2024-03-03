import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { OpenResourceRequest } from '../types';
import { Observable } from 'rxjs';

@Injectable()
export class HomeService {
  private readonly _resourceUrlSegment: string = '/api/resource';

  constructor(private readonly http: HttpClient) {}

  openResource(request: OpenResourceRequest): Observable<void> {
    return this.http.post<void>(this._resourceUrlSegment + '/open', request);
  }
}
