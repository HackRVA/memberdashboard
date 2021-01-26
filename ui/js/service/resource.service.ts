import { Observable } from "rxjs";
import { HTTPService } from "./http.service";

export class ResourceService extends HTTPService {
  getResources(): Observable<Response | { error: boolean; message: any }> {
    return this.get("/api/resource");
  }
  register(
    registerRequest: ResourceService.RegisterResourceRequest
  ): Observable<Response | { error: boolean; message: any }> {
    return this.post("/api/resource/register", registerRequest);
  }
  deleteResource(
    deleteRequest: ResourceService.RemoveResourceRequest
  ): Observable<Response | { error: boolean; message: any }> {
    return this.delete("/api/resource", deleteRequest);
  }
  updateResource(
    updateRequest: ResourceService.UpdateResourceRequest
  ): Observable<Response | { error: boolean; message: any }> {
    return this.put("/api/resource", updateRequest);
  }
}

export namespace ResourceService {
  export interface RegisterResourceRequest {
    address: string;
    name: string;
  }
  export interface UpdateResourceRequest {
    address: string;
    id: number;
    name: string;
  }
  export interface RemoveResourceRequest {
    id: number;
  }
  export interface ResourceResponse {
    address: string;
    id: number;
    name: string;
  }
  export enum ResourceStatus {
    good = 0,
    outOfDate = 1,
    offline = 2,
  }
}
