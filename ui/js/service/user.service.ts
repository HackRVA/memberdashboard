// rxjs
import { Observable } from "rxjs";

// membership
import { ENV } from "../env";
import { HTTPService } from "./http.service";
import { UserProfile } from "./../components/user/types/user.interface";

export class UserService extends HTTPService {
  private readonly userUrlSegment: string = ENV.api + "/user";

  getUser(): Observable<UserProfile> {
    return this.get<UserProfile>(this.userUrlSegment);
  }
}
