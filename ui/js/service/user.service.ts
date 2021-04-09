// rxjs
import { Observable } from "rxjs";

// membership
import { ENV } from "../env";
import { HTTPService } from "./http.service";

export class UserService extends HTTPService {
  private readonly userUrlSegment: string = ENV.api + "/user";

  getUser(): Observable<Response | { error: boolean; message: any }> {
    return this.get(this.userUrlSegment);
  }
}
