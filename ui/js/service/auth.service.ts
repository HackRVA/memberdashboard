// rxjs
import { Observable } from "rxjs";
import { map } from "rxjs/operators";

// membership
import { LoginRequest, RegisterRequest } from "../components/user/types";
import { ENV } from "../env";
import { HTTPService } from "./http.service";

export class AuthService extends HTTPService {
  private readonly authUrlSegment: string = ENV.api + "/auth";

  registerUser(
    request: RegisterRequest
  ): Observable<Response | { error: boolean; message: any }> {
    return this.post(this.authUrlSegment + "/register", request);
  }

  login(
    request: LoginRequest
  ): Observable<Response | { error: boolean; message: any }> {
    return this.post(this.authUrlSegment + "/login", request);
  }

  logout(): Observable<null> {
    return this.post(this.authUrlSegment + "/logout").pipe(map(() => null));
  }
}
