// rxjs
import { Observable } from "rxjs";
import { map } from "rxjs/operators";

// membership
import { ENV } from "../env";
import { HTTPService } from "./http.service";
import { LoginRequest, RegisterRequest } from "./../components/user/types";

export class UserService extends HTTPService {
  private readonly api: string = ENV.api;

  login(
    request: LoginRequest
  ): Observable<Response | { error: boolean; message: any }> {
    return this.post(this.api + "/login", request);
  }

  logout(): Observable<null> {
    return this.post(this.api + "/logout").pipe(map(() => null));
  }

  getUser(): Observable<Response | { error: boolean; message: any }> {
    return this.get(this.api + "/user");
  }

  registerUser(
    request: RegisterRequest
  ): Observable<Response | { error: boolean; message: any }> {
    return this.post(this.api + "/register", request);
  }
}
