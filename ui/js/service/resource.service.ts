import { Observable } from "rxjs";
import { ENV } from "../env";
import { HTTPService } from "./http.service";
import {
  RegisterResourceRequest,
  RemoveResourceRequest,
  AddMemberResourceRequest,
  RemoveMemberResourceRequest,
  UpdateResourceRequest,
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
  addMemberToResource(
    request: AddMemberResourceRequest
  ): Observable<Response | { error: boolean; message: any }> {
    return this.post(this.api + "/resource/member", request);
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
