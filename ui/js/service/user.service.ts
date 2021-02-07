import { Observable } from "rxjs";
import { map } from "rxjs/operators";
import { ENV } from "../env";
import { HTTPService } from "./http.service";

export class UserService extends HTTPService {
  private readonly api: string | undefined = ENV.api;

  login(
    request: UserService.LoginRequest
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
    request: UserService.RegisterRequest
  ): Observable<Response | { error: boolean; message: any }> {
    return this.post(this.api + "/register", request);
  }
}

export namespace UserService {
  export interface RegisterRequest {
    password: string;
    email: string;
  }

  export interface LoginRequest {
    email: string;
    password: string;
    updateCallback?: Function;
  }

  export interface Jwt {
    token: string;
  }

  export interface UserProfile {
    email: string;
  }
}
