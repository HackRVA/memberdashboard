import { Observable } from "rxjs";
import { HTTPService } from "./http.service";

export class UserService extends HTTPService {
  login(
    loginRequest: UserService.LoginRequest
  ): Observable<UserService.ILoginResponse> {
    return this.post("/edge/api/login", loginRequest);
  }

  logout(): Observable<Response | { error: boolean; message: any }> {
    return this.get("/edge/api/logout");
  }

  getUser(): Observable<Response | { error: boolean; message: any }> {
    return this.get("/edge/api/user");
  }

  registerUser(
    registerRequest: UserService.RegisterRequest
  ): Observable<Response | { error: boolean; message: any }> {
    return this.post("/edge/api/register", registerRequest);
  }
}

export namespace UserService {
  export interface RegisterRequest {
    username: string;
    password: string;
    email: string;
  }

  export interface LoginRequest {
    username: string;
    password: string;
    updateCallback?: Function;
  }

  export interface ILoginResponse {
    token: string;
  }

  export interface UserProfile {
    username: string;
    email: string;
  }
}
