import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import {
  BulkAddMembersToResourceRequest,
  RegisterResourceRequest,
  RemoveMemberResourceRequest,
  RemoveResourceRequest,
  ResourceResponse,
  UpdateResourceRequest,
} from '../types/resource.types';

@Injectable({
  providedIn: 'root',
})
export class ResourceService {
  private readonly _resourceUrlSegment: string = '/api/resource';
  constructor(private readonly http: HttpClient) {}

  getResources(): Observable<ResourceResponse[]> {
    return this.http.get<ResourceResponse[]>(this._resourceUrlSegment);
  }

  deleteResource(request: RemoveResourceRequest): Observable<void> {
    return this.http.request<void>('DELETE', this._resourceUrlSegment, {
      body: request,
    });
  }

  updateResource(request: UpdateResourceRequest): Observable<void> {
    return this.http.put<void>(this._resourceUrlSegment, request);
  }

  registerResource(request: RegisterResourceRequest): Observable<void> {
    return this.http.post<void>(
      this._resourceUrlSegment + '/register',
      request
    );
  }

  updateACLs(): Observable<void> {
    return this.http.post<void>(this._resourceUrlSegment + '/updateacls', {});
  }

  removeACLs(): Observable<void> {
    return this.http.delete<void>(this._resourceUrlSegment + '/deleteacls');
  }

  bulkAddMembersToResource(
    request: BulkAddMembersToResourceRequest
  ): Observable<void> {
    return this.http.post<void>(
      this._resourceUrlSegment + '/member/bulk',
      request
    );
  }

  removeMemberFromResource(
    request: RemoveMemberResourceRequest
  ): Observable<void> {
    return this.http.request<void>(
      'DELETE',
      this._resourceUrlSegment + '/member',
      {
        body: request,
      }
    );
  }
}
