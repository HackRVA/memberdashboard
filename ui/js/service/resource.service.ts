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
  private readonly api: string = ENV.api;

  getResources(): Observable<Response | { error: boolean; message: any }> {
    return this.get(this.api + "/resource");
  }

  register(
    request: RegisterResourceRequest
  ): Observable<Response | { error: boolean; message: any }> {
    return this.post(this.api + "/resource/register", request);
  }

  deleteResource(
    request: RemoveResourceRequest
  ): Observable<Response | { error: boolean; message: any }> {
    return this.delete(this.api + "/resource", request);
  }

  updateACLs(): Observable<Response | { error: boolean; message: any }> {
    return this.post(this.api + "/resource/updateacls", {});
  }

  removeACLs(): Observable<Response | { error: boolean; message: any }> {
    return this.delete(this.api + "/resource/deleteacls", {});
  }

  bulkAddMembersToResource(
    request: BulkAddMembersToResourceRequest
  ): Observable<Response | { error: boolean; message: any }> {
    return this.post(this.api + "/resource/member/bulk", request);
  }

  removeMemberFromResource(
    request: RemoveMemberResourceRequest
  ): Observable<Response | { error: boolean; message: any }> {
    return this.delete(this.api + "/resource/member", request);
  }

  updateResource(
    request: UpdateResourceRequest
  ): Observable<Response | { error: boolean; message: any }> {
    return this.put(this.api + "/resource", request);
  }
}
