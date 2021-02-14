// rxjs
import { Observable } from "rxjs";

// membership
import { HTTPService } from "./http.service";
import { ENV } from "./../env";
import { AssignRFIDRequest } from "../components/members/types";

export class MemberService extends HTTPService {
  private readonly api: string = ENV.api;

  getMembers(): Observable<Response | { error: boolean; message: any }> {
    return this.get(this.api + "/member");
  }

  assignRFID(
    request: AssignRFIDRequest
  ): Observable<Response | { error: boolean; message: any }> {
    return this.post(this.api + "/assignRFID", request);
  }
}
