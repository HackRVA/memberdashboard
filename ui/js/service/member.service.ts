import { Observable } from "rxjs";
import { HTTPService } from "./http.service";
export class MemberService extends HTTPService {
  getMembers(): Observable<Response | { error: boolean; message: any }> {
    return this.get("/api/members");
  }
}
