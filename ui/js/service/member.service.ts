import { Observable } from "rxjs";
import { HTTPService } from "./http.service";

export class MemberService extends HTTPService {
  getMembers(): Observable<Response | { error: boolean; message: any }> {
    return this.get("/api/member");
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
