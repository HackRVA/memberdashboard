import { Observable } from "rxjs";
import { HTTPService } from "./http.service";
import { ENV } from "./../env";

export class MemberService extends HTTPService {
  private readonly api: string | undefined = ENV.api;

  getMembers(): Observable<Response | { error: boolean; message: any }> {
    return this.get(this.api + "/member");
  }

  assignRFID(
    request: MemberService.AssignRFIDRequest
  ): Observable<Response | { error: boolean; message: any }> {
    return this.post(this.api + "/assignRFID", request);
  }
}

export namespace MemberService {
  export interface AssignRFIDRequest {
    email: string;
    rfid: string;
  }
  export interface MemberResponse {
    id: string;
    name: string;
    email: string;
    memberLevel: MemberLevel;
    resources: Array<MemberResource>;
  }

  export interface MemberResource {
    resourceID: string;
    name: string;
  }

  export enum MemberLevel {
    inactive = 1,
    student = 2,
    classic = 3,
    standard = 4,
    premium = 5,
  }
}
