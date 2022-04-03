// rxjs
import { Observable } from 'rxjs';

// memberdashboard
import { ENV } from '../../env';
import { UpdateResourceRequest } from '../types/api/update-resource-request';
import { OpenResourceRequest } from './../types/api/open-resource-request';
import { ResourceResponse } from '../types/api/resource-response';
import { RemoveResourceRequest } from '../types/api/remove-resource-request';
import { RemoveMemberResourceRequest } from '../types/api/remove-member-resource-request';
import { RegisterResourceRequest } from '../types/api/register-resource-request';
import { HTTPService } from '../../shared/services/http.service';
import { BulkAddMembersToResourceRequest } from '../types/api/bulk-add-members-to-resource-request';
import { Injectable } from '../../shared/di';

@Injectable('resource')
export class ResourceService extends HTTPService {
  private readonly resourceUrlSegment: string = ENV.api + '/resource';

  getResources(): Observable<ResourceResponse[]> {
    return this.get<ResourceResponse[]>(this.resourceUrlSegment);
  }

  register(request: RegisterResourceRequest): Observable<void> {
    return this.post<void>(this.resourceUrlSegment + '/register', request);
  }

  deleteResource(request: RemoveResourceRequest): Observable<void> {
    return this.delete<void>(this.resourceUrlSegment, request);
  }

  updateACLs(): Observable<void> {
    return this.post<void>(this.resourceUrlSegment + '/updateacls', {});
  }

  removeACLs(): Observable<void> {
    return this.delete<void>(this.resourceUrlSegment + '/deleteacls', {});
  }

  bulkAddMembersToResource(
    request: BulkAddMembersToResourceRequest
  ): Observable<void> {
    return this.post<void>(this.resourceUrlSegment + '/member/bulk', request);
  }

  removeMemberFromResource(
    request: RemoveMemberResourceRequest
  ): Observable<void> {
    return this.delete<void>(this.resourceUrlSegment + '/member', request);
  }

  updateResource(request: UpdateResourceRequest): Observable<void> {
    return this.put<void>(this.resourceUrlSegment, request);
  }

  openResource(request: OpenResourceRequest): Observable<void> {
    return this.post<void>(this.resourceUrlSegment + '/open', request);
  }
}
