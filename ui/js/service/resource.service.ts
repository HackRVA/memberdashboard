// rxjs
import { Observable } from "rxjs";

// membership
import { ENV } from "../env";
import { HTTPService } from "./http.service";
import {
  RegisterResourceRequest,
  RemoveResourceRequest,
  RemoveMemberResourceRequest,
  UpdateResourceRequest,
  BulkAddMembersToResourceRequest,
} from "../components/resources/types";

export class ResourceService extends HTTPService {
  private readonly resourceUrlSegment: string = ENV.api + "/resource";

  getResources(): Observable<Response | { error: boolean; message: any }> {
    return this.get(this.resourceUrlSegment);
  }

  register(
    request: RegisterResourceRequest
  ): Observable<Response | { error: boolean; message: any }> {
    return this.post(this.resourceUrlSegment + "/register", request);
  }

  deleteResource(
    request: RemoveResourceRequest
  ): Observable<Response | { error: boolean; message: any }> {
    return this.delete(this.resourceUrlSegment, request);
  }

  updateACLs(): Observable<Response | { error: boolean; message: any }> {
    return this.post(this.resourceUrlSegment + "/updateacls", {});
  }

  removeACLs(): Observable<Response | { error: boolean; message: any }> {
    return this.delete(this.resourceUrlSegment + "/deleteacls", {});
  }

  bulkAddMembersToResource(
    request: BulkAddMembersToResourceRequest
  ): Observable<Response | { error: boolean; message: any }> {
    return this.post(this.resourceUrlSegment + "/member/bulk", request);
  }

  removeMemberFromResource(
    request: RemoveMemberResourceRequest
  ): Observable<Response | { error: boolean; message: any }> {
    return this.delete(this.resourceUrlSegment + "/member", request);
  }

  updateResource(
    request: UpdateResourceRequest
  ): Observable<Response | { error: boolean; message: any }> {
    return this.put(this.resourceUrlSegment, request);
  }
}
