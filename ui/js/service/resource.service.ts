import { Observable } from "rxjs";
import { HTTPService } from "./http.service";

export class ResourceService extends HTTPService {
  getResources(): Observable<Response | { error: boolean; message: any }> {
    return this.get("/edge/api/resource");
  }
  register(
    addRequest: ResourceService.AddResourceRequest
  ): Observable<Response | { error: boolean; message: any }> {
    return this.post("/edge/api/resource/register", addRequest);
  }

  updateResource(
    updateRequest: ResourceService.UpdateResourceRequest
  ): Observable<Response | { error: boolean; message: any }> {
    return this.post("/edge/api/resource", updateRequest);
  }
  deleteResource(
    deleteRequest: ResourceService.RemoveResourceRequest
  ): Observable<Response | { error: boolean; message: any }> {
    return this.delete("/edge/api/resource", deleteRequest);
  }
  addMember(
    addMemberRequest: ResourceService.AddMemberRequest
  ): Observable<Response | { error: boolean; message: any }> {
    return this.post("/edge/api/resource/member/add", addMemberRequest);
  }
  removeMember(
    removeMemberRequest: ResourceService.RemoveMemberRequest
  ): Observable<Response | { error: boolean; message: any }> {
    return this.post("/edge/api/resource/member/remove", removeMemberRequest);
  }
}

export namespace ResourceService {
  export interface AddResourceRequest {
    name: string;
    address: string;
  }

  export interface UpdateResourceRequest {
    id: number;
    name: string;
    address: string;
  }
  export interface RemoveResourceRequest {
    id: number;
    name: string;
    address: string;
  }

  export interface ResourceResponse {
    id: number;
    address: string;
    name: string;
  }

  export interface AddMemberRequest {
    email: string;
    resourceID: number;
  }

  export interface RemoveMemberRequest {
    email: string;
    resourceID: number;
  }
}
