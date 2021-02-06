import { Observable } from "rxjs";
import { HTTPService } from "./http.service";
import { ENV } from "./../env";
import { ResourceService } from "./resource.service";

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
  export interface RemoveMemberResourceModalData {
    email: string;
    memberResources: MemberResource[];
    handleResourceChange: Function;
    handleSubmitRemoveMemberResource: Function;
    emptyFormValuesOnClosed: Function;
  }

  export interface AddMemberResourceModalData {
    email: string;
    resources: ResourceService.ResourceResponse[];
    handleResourceChange: Function;
    handleSubmitAddMemberResource: Function;
    emptyFormValuesOnClosed: Function;
  }

  export interface RFIDModalData {
    email: string;
    rfid: string;
    handleEmailChange: Function;
    handleRFIDChange: Function;
    handleSubmitForAssigningMemberToRFID: Function;
    emptyFormValuesOnClosed: Function;
  }

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
