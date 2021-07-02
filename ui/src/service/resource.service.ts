// rxjs
import { Observable } from "rxjs";

// memberdashboard
import { ENV } from "../env";
import { HTTPService } from "./http.service";
import {
  RegisterResourceRequest,
  RemoveResourceRequest,
  RemoveMemberResourceRequest,
  UpdateResourceRequest,
  BulkAddMembersToResourceRequest,
  ResourceResponse,
} from "../components/resources/types";

export class ResourceService extends HTTPService {
  private readonly resourceUrlSegment: string = ENV.api + "/resource";

  getResources(): Observable<ResourceResponse[]> {
    return this.get<ResourceResponse[]>(this.resourceUrlSegment);
  }

  register(request: RegisterResourceRequest): Observable<void> {
    return this.post<void>(this.resourceUrlSegment + "/register", request);
  }

  deleteResource(request: RemoveResourceRequest): Observable<void> {
    return this.delete<void>(this.resourceUrlSegment, request);
  }

  updateACLs(): Observable<void> {
    return this.post<void>(this.resourceUrlSegment + "/updateacls", {});
  }

  removeACLs(): Observable<void> {
    return this.delete<void>(this.resourceUrlSegment + "/deleteacls", {});
  }

  bulkAddMembersToResource(
    request: BulkAddMembersToResourceRequest
  ): Observable<void> {
    return this.post<void>(this.resourceUrlSegment + "/member/bulk", request);
  }

  removeMemberFromResource(
    request: RemoveMemberResourceRequest
  ): Observable<void> {
    return this.delete<void>(this.resourceUrlSegment + "/member", request);
  }

  updateResource(request: UpdateResourceRequest): Observable<void> {
    return this.put<void>(this.resourceUrlSegment, request);
  }
}
