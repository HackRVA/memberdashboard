import { Observable } from "rxjs";
import { HTTPService } from "./http.service";
import { ENV } from "./../env";

export class MemberService extends HTTPService {
  private readonly api: string | undefined = ENV.api;

  getMembers(): Observable<Response | { error: boolean; message: any }> {
    return this.get(this.api + "/member");
  }
}

export namespace MemberService {
  export interface MemberResponse {
    id: number;
    name: string;
    email: string;
    memberLevel: MemberLevel;
    resources: Array<MemberResource>;
  }

  export interface MemberResource {
    resourceID: number;
    name: string;
  }

  export enum MemberLevel {
    inactive = 1,
    student = 2,
    standard = 3,
    premium = 4,
  }
}
