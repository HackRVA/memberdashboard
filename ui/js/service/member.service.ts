// rxjs
import { Observable } from "rxjs";

// membership
import { HTTPService } from "./http.service";
import { ENV } from "./../env";
import { AssignRFIDRequest } from "../components/members/types";

export class MemberService extends HTTPService {
  private readonly memberUrlSegment: string = ENV.api + "/member";

  getMembers(): Observable<Response | { error: boolean; message: any }> {
    return this.get(this.memberUrlSegment);
  }

  assignRFID(
    request: AssignRFIDRequest
  ): Observable<Response | { error: boolean; message: any }> {
    return this.post(this.memberUrlSegment + "/assignRFID", request);
  }
}
