import { Observable } from "rxjs";
import { HTTPService } from "./HTTPService";

export class ResourceService extends HTTPService {
  getResources(): Observable<Response | { error: boolean; message: any }> {
    return this.get("/api/resource");
  }
  register(
    updateRequest: ResourceService.ResourceRequest
  ): Observable<Response | { error: boolean; message: any }> {
    return this.post("/api/resource", updateRequest);
  }
  deleteResource(
    deleteRequest: ResourceService.ResourceRequest
  ): Observable<Response | { error: boolean; message: any }> {
    return this.delete("/api/resource", deleteRequest);
  }
}

export namespace ResourceService {
  export interface ResourceRequest {
    id?: number;
    name: string;
    address: string;
    // email is added to the request when attaching a member to a resource
    email?: string; 
  }
}
